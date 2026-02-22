// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package channel_telephony

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"

	"github.com/rapidaai/api/assistant-api/config"
	callcontext "github.com/rapidaai/api/assistant-api/internal/callcontext"
	internal_services "github.com/rapidaai/api/assistant-api/internal/services"
	web_client "github.com/rapidaai/pkg/clients/web"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/types"
	type_enums "github.com/rapidaai/pkg/types/enums"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

// InboundDispatcher handles inbound call processing across all telephony
// channels (SIP, Asterisk, Twilio, Exotel, Vonage). It encapsulates the
// common business logic: provider resolution, call reception, conversation
// creation, call-context persistence, telemetry application, and session resolution.
type InboundDispatcher struct {
	cfg                 *config.AssistantConfig
	store               callcontext.Store
	logger              commons.Logger
	vaultClient         web_client.VaultClient
	assistantService    internal_services.AssistantService
	conversationService internal_services.AssistantConversationService
	telephonyOpt        TelephonyOption
}

// NewInboundDispatcher creates a new inbound call dispatcher.
func NewInboundDispatcher(deps TelephonyDispatcherDeps) *InboundDispatcher {
	return &InboundDispatcher{
		cfg:                 deps.Cfg,
		store:               deps.Store,
		logger:              deps.Logger,
		vaultClient:         deps.VaultClient,
		assistantService:    deps.AssistantService,
		conversationService: deps.ConversationService,
		telephonyOpt:        deps.TelephonyOpt,
	}
}

// HandleStatusCallback resolves the telephony provider and processes a status callback
// webhook. It builds telemetry (metric + event) from the StatusInfo returned by the provider.
func (d *InboundDispatcher) HandleStatusCallback(c *gin.Context, provider string, auth types.SimplePrinciple, assistantId, conversationId uint64) error {
	tel, err := GetTelephony(Telephony(provider), d.cfg, d.logger, d.telephonyOpt)
	if err != nil {
		return fmt.Errorf("invalid telephony provider %s: %w", provider, err)
	}

	statusInfo, err := tel.StatusCallback(c, auth, assistantId, conversationId)
	if err != nil {
		return fmt.Errorf("status callback failed: %w", err)
	}
	if statusInfo == nil {
		return nil
	}

	// Build telemetry from StatusInfo — the dispatcher owns telemetry construction.
	metric := types.NewMetric("STATUS", statusInfo.Event, utils.Ptr("Status of conversation"))
	if _, err := d.conversationService.ApplyConversationMetrics(c, auth, assistantId, conversationId, []*types.Metric{metric}); err != nil {
		d.logger.Errorf("failed to apply conversation metrics in callback: %v", err)
		return fmt.Errorf("failed to process metrics: %w", err)
	}

	event := types.NewEvent(statusInfo.Event, statusInfo.Payload)
	if _, err := d.conversationService.ApplyConversationTelephonyEvent(c, auth, provider, assistantId, conversationId, []*types.Event{event}); err != nil {
		d.logger.Errorf("failed to apply telephony events in callback: %v", err)
		return fmt.Errorf("failed to process events: %w", err)
	}
	return nil
}

// HandleStatusCallbackByContext resolves a call context from Postgres using the contextId and
// processes the status callback. Unlike ResolveCallSessionByContext, this reads the context
// without changing its status, since status callbacks can fire multiple times during a call
// and may arrive asynchronously even after the call has ended (status="completed").
func (d *InboundDispatcher) HandleStatusCallbackByContext(c *gin.Context, contextID string) error {
	cc, err := d.store.Get(c, contextID)
	if err != nil {
		d.logger.Errorf("failed to resolve call context %s for event callback: %v", contextID, err)
		return fmt.Errorf("call context not found or expired: %w", err)
	}

	auth := cc.ToAuth()
	return d.HandleStatusCallback(c, cc.Provider, auth, cc.AssistantID, cc.ConversationID)
}

// HandleReceiveCall processes an inbound call webhook. It resolves the telephony provider,
// receives the call, creates a conversation, saves a CallContext in Postgres, applies telemetry,
// and instructs the provider to answer the call.
// Returns the contextID for AudioSocket/WebSocket resolution.
func (d *InboundDispatcher) HandleReceiveCall(c *gin.Context, provider string, auth types.SimplePrinciple, assistantId uint64) (string, error) {
	tel, err := GetTelephony(Telephony(provider), d.cfg, d.logger, d.telephonyOpt)
	if err != nil {
		return "", fmt.Errorf("telephony provider %s not connected: %w", provider, err)
	}

	callInfo, err := tel.ReceiveCall(c)
	if err != nil {
		return "", fmt.Errorf("receive call failed: %w", err)
	}

	assistant, err := d.assistantService.Get(c, auth, assistantId, utils.GetVersionDefinition("latest"), &internal_services.GetAssistantOption{InjectPhoneDeployment: true})
	if err != nil {
		d.logger.Debugf("unable to find assistant %v", err)
		return "", fmt.Errorf("unable to find assistant: %w", err)
	}

	conversation, err := d.conversationService.CreateConversation(c, auth, callInfo.CallerNumber, assistant.Id, assistant.AssistantProviderId, type_enums.DIRECTION_INBOUND, utils.PhoneCall)
	if err != nil {
		return "", fmt.Errorf("unable to create conversation: %w", err)
	}

	// Build and apply telemetry from CallInfo — the dispatcher owns telemetry construction.
	var wg errgroup.Group

	// Apply metadata from CallInfo.Extra (provider-specific fields)
	wg.Go(func() error {
		var metadatas []*types.Metadata
		for k, v := range callInfo.Extra {
			metadatas = append(metadatas, types.NewMetadata(k, v))
		}
		if len(metadatas) > 0 {
			mtdas, err := d.conversationService.ApplyConversationMetadata(c, auth, assistant.Id, conversation.Id, metadatas)
			if err != nil {
				d.logger.Errorf("failed to apply conversation metadata: %v", err)
				return err
			}
			conversation.Metadatas = mtdas
		}
		return nil
	})

	// Apply metric from CallInfo.Status
	wg.Go(func() error {
		metric := types.NewMetric("STATUS", callInfo.Status, utils.Ptr("Status of telephony api"))
		metrics, err := d.conversationService.ApplyConversationMetrics(c, auth, assistant.Id, conversation.Id, []*types.Metric{metric})
		if err != nil {
			d.logger.Errorf("failed to apply conversation metrics: %v", err)
			return err
		}
		conversation.Metrics = append(conversation.Metrics, metrics...)
		return nil
	})

	// Apply telephony event from CallInfo.StatusInfo
	wg.Go(func() error {
		event := types.NewEvent(callInfo.StatusInfo.Event, callInfo.StatusInfo.Payload)
		evts, err := d.conversationService.ApplyConversationTelephonyEvent(c, auth, assistant.AssistantPhoneDeployment.TelephonyProvider, assistant.Id, conversation.Id, []*types.Event{event})
		if err != nil {
			d.logger.Errorf("failed to apply telephony events: %v", err)
			return err
		}
		conversation.TelephonyEvents = append(conversation.TelephonyEvents, evts...)
		return nil
	})

	if err := wg.Wait(); err != nil {
		d.logger.Errorf("failed to process telemetry for inbound call: %v", err)
		return "", fmt.Errorf("failed to process call telemetry: %w", err)
	}

	// Store call context in Postgres for AudioSocket/WebSocket resolution.
	// ChannelUUID comes directly from CallInfo — no need to scan metadata.
	cc := &callcontext.CallContext{
		AssistantID:         assistant.Id,
		ConversationID:      conversation.Id,
		AssistantProviderId: assistant.AssistantProviderId,
		AuthToken:           auth.GetCurrentToken(),
		AuthType:            auth.Type(),
		Direction:           "inbound",
		CallerNumber:        callInfo.CallerNumber,
		Provider:            provider,
		ChannelUUID:         callInfo.ChannelUUID,
	}
	if auth.GetCurrentProjectId() != nil {
		cc.ProjectID = *auth.GetCurrentProjectId()
	}
	if auth.GetCurrentOrganizationId() != nil {
		cc.OrganizationID = *auth.GetCurrentOrganizationId()
	}
	contextID, err := d.store.Save(c, cc)
	if err != nil {
		d.logger.Errorf("failed to save call context: %v", err)
		return "", fmt.Errorf("failed to create call context: %w", err)
	}

	// Pass contextId to the telephony provider for inbound call setup
	// For Asterisk: the contextId is returned as plain text — dialplan uses it as the AudioSocket UUID
	// For WebSocket providers: the contextId is embedded in the WebSocket URL path
	c.Set("contextId", contextID)

	if err := tel.InboundCall(c, auth, assistant.Id, callInfo.CallerNumber, conversation.Id); err != nil {
		d.logger.Errorf("failed to initiate inbound call: %v", err)
		return "", fmt.Errorf("unable to initiate inbound call: %w", err)
	}

	return contextID, nil
}

// ResolveVaultCredential fetches the vault credential for the given assistant.
// This is the only DB round-trip needed — call IDs (assistant, conversation,
// provider) are already in the CallContext from Redis.
func (d *InboundDispatcher) ResolveVaultCredential(ctx context.Context, auth types.SimplePrinciple, assistantId, conversationId uint64) (*protos.VaultCredential, error) {
	assistant, err := d.assistantService.Get(ctx, auth, assistantId, nil, &internal_services.GetAssistantOption{InjectPhoneDeployment: true})
	if err != nil {
		return nil, err
	}
	if !assistant.IsPhoneDeploymentEnable() {
		return nil, fmt.Errorf("phone deployment not enabled for assistant %d", assistantId)
	}
	credentialID, err := assistant.AssistantPhoneDeployment.GetOptions().GetUint64("rapida.credential_id")
	if err != nil {
		return nil, err
	}
	vltC, err := d.vaultClient.GetCredential(ctx, auth, credentialID)
	if err != nil {
		d.conversationService.ApplyConversationMetrics(ctx, auth, assistantId, conversationId, []*types.Metric{{Name: type_enums.STATUS.String(), Value: type_enums.RECORD_FAILED.String(), Description: "Failed to resolve vault credential"}})
		return nil, fmt.Errorf("failed to resolve vault credential: %w", err)
	}
	return vltC, nil
}

// ResolveCallSessionByContext resolves a call context and vault credential using
// a contextId stored in Postgres. The call context is atomically claimed by
// transitioning its status from "pending" to "claimed". Only one media connection
// can claim a given context — subsequent callers get an error.
// The context remains in Postgres so that event/status callbacks can still read it.
// Returns the CallContext (which contains all IDs and auth info) plus the vault
// credential needed for the streamer.
func (d *InboundDispatcher) ResolveCallSessionByContext(ctx context.Context, contextID string) (*callcontext.CallContext, *protos.VaultCredential, error) {
	cc, err := d.store.Claim(ctx, contextID)
	if err != nil {
		d.logger.Errorf("failed to resolve call context %s: %v", contextID, err)
		return nil, nil, fmt.Errorf("call context not found or already claimed: %w", err)
	}

	auth := cc.ToAuth()
	vaultCred, err := d.ResolveVaultCredential(ctx, auth, cc.AssistantID, cc.ConversationID)
	if err != nil {
		return nil, nil, err
	}
	return cc, vaultCred, nil
}

// CompleteCallSession marks a call context as completed. Should be called
// when the call/session ends (talker exits).
func (d *InboundDispatcher) CompleteCallSession(ctx context.Context, contextID string) {
	if err := d.store.Complete(ctx, contextID); err != nil {
		d.logger.Warnf("failed to complete call context %s: %v", contextID, err)
	}
}
