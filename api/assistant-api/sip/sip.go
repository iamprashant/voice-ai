// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package assistant_sip

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/rapidaai/api/assistant-api/config"
	internal_adapter "github.com/rapidaai/api/assistant-api/internal/adapters"
	internal_telephony "github.com/rapidaai/api/assistant-api/internal/channel/telephony"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_services "github.com/rapidaai/api/assistant-api/internal/services"
	internal_assistant_service "github.com/rapidaai/api/assistant-api/internal/services/assistant"
	sip_infra "github.com/rapidaai/api/assistant-api/sip/infra"
	web_client "github.com/rapidaai/pkg/clients/web"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
	"github.com/rapidaai/pkg/storages"
	storage_files "github.com/rapidaai/pkg/storages/file-storage"
	"github.com/rapidaai/pkg/types"
	type_enums "github.com/rapidaai/pkg/types/enums"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

// SIPManager manages SIP connections for voice conversations
// SIP uses native signaling (UDP/TCP/TLS) and RTP for audio - no WebSocket
// Multi-tenant: Single shared server, config resolved per-call from assistant's deployment
type SIPManager struct {
	mu       sync.RWMutex
	cfg      *config.AssistantConfig
	logger   commons.Logger
	server   *sip_infra.Server
	sessions map[string]*sip_infra.SIPSession

	ctx    context.Context
	cancel context.CancelFunc

	postgres   connectors.PostgresConnector
	redis      connectors.RedisConnector
	opensearch connectors.OpenSearchConnector
	storage    storages.Storage

	assistantConversationService internal_services.AssistantConversationService
	assistantService             internal_services.AssistantService
	vaultClient                  web_client.VaultClient
	authClient                   web_client.AuthClient
}

// NewSIPManager creates a new SIP manager
// Multi-tenant: single server, config resolved per-call via ConfigResolver
func NewSIPManager(config *config.AssistantConfig, logger commons.Logger,
	postgres connectors.PostgresConnector,
	redis connectors.RedisConnector,
	opensearch connectors.OpenSearchConnector,
	vectordb connectors.VectorConnector) *SIPManager {
	return &SIPManager{
		cfg:                          config,
		logger:                       logger,
		postgres:                     postgres,
		redis:                        redis,
		opensearch:                   opensearch,
		assistantConversationService: internal_assistant_service.NewAssistantConversationService(logger, postgres, storage_files.NewStorage(config.AssetStoreConfig, logger)),
		assistantService:             internal_assistant_service.NewAssistantService(config, logger, postgres, opensearch),
		storage:                      storage_files.NewStorage(config.AssetStoreConfig, logger),
		vaultClient:                  web_client.NewVaultClientGRPC(&config.AppConfig, logger, redis),
		authClient:                   web_client.NewAuthenticator(&config.AppConfig, logger, redis),
		sessions:                     make(map[string]*sip_infra.SIPSession),
	}
}

// NewSIPListenConfig creates a ListenConfig for the shared SIP server
// Multi-tenant: Server listens on this address, tenant config resolved per-call
func (m *SIPManager) listenConfig() *sip_infra.ListenConfig {
	transportType := sip_infra.TransportUDP
	switch m.cfg.SIPConfig.Transport {
	case "tcp":
		transportType = sip_infra.TransportTCP
	case "tls":
		transportType = sip_infra.TransportTLS
	}
	lc := &sip_infra.ListenConfig{
		Address:    m.cfg.SIPConfig.Server,
		ExternalIP: m.cfg.SIPConfig.ExternalIP,
		Port:       m.cfg.SIPConfig.Port,
		Transport:  transportType,
	}
	m.logger.Infow("SIP ListenConfig from app config",
		"address", lc.Address,
		"external_ip", lc.ExternalIP,
		"port", lc.Port,
		"transport", lc.Transport,
		"raw_sip_config_external_ip", m.cfg.SIPConfig.ExternalIP,
		"raw_sip_config_server", m.cfg.SIPConfig.Server)
	return lc
}

// Start initializes the shared SIP server with per-call middleware-based authentication.
// The middleware chain authenticates every SIP request (not just INVITE):
//
//	CredentialMiddleware → AuthMiddleware → AssistantMiddleware → VaultConfigMiddleware
//
// URI format: sip:{assistantID}:{apiKey}@aws.ap-south-east-01.rapida.ai
func (m *SIPManager) Start(ctx context.Context) error {
	m.ctx, m.cancel = context.WithCancel(ctx)

	server, err := sip_infra.NewServer(m.ctx, &sip_infra.ServerConfig{
		ListenConfig:      m.listenConfig(),
		Logger:            m.logger,
		RedisClient:       m.redis.GetConnection(),
		RTPPortRangeStart: m.cfg.SIPConfig.RTPPortRangeStart,
		RTPPortRangeEnd:   m.cfg.SIPConfig.RTPPortRangeEnd,
	})
	if err != nil {
		return fmt.Errorf("failed to create SIP server: %w", err)
	}

	// Register middleware chain for SIP request authentication.
	// Each middleware enriches the SIPRequestContext; the final handler
	// returns the InviteResult with the resolved SIP config.
	server.SetMiddlewares(
		[]sip_infra.Middleware{
			sip_infra.CredentialMiddleware, // Parse assistantID:apiKey from URI
			m.authMiddleware,               // Validate API key → set auth principal
			m.assistantMiddleware,          // Load assistant → set assistant entity
		},
		m.vaultConfigResolver, // Fetch SIP config from vault (final handler)
	)

	// Set up handlers for incoming calls
	server.SetOnInvite(m.handleInvite)
	server.SetOnBye(m.handleBye)
	server.SetOnCancel(m.handleCancel)

	// Start the server
	if err := server.Start(); err != nil {
		return fmt.Errorf("failed to start SIP server: %w", err)
	}
	m.server = server
	return nil
}

// GetServer returns the shared SIP server instance.
// Used for dependency injection into outbound call providers.
func (m *SIPManager) GetServer() *sip_infra.Server {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.server
}

// authMiddleware validates the API key extracted by CredentialMiddleware.
// Sets Extra["auth"] with the authenticated principal for downstream middlewares.
//
// URI format: sip:{assistantID}:{apiKey}@aws.ap-south-east-01.rapida.ai
//   - apiKey is a project-scoped token (rpd-prj-xxx or raw key)
func (m *SIPManager) authMiddleware(ctx *sip_infra.SIPRequestContext, next func() (*sip_infra.InviteResult, error)) (*sip_infra.InviteResult, error) {
	if ctx.APIKey == "" {
		m.logger.Warnw("SIP: missing API key", "call_id", ctx.CallID, "method", ctx.Method, "from", ctx.FromURI)
		return sip_infra.Reject(401, "Missing credentials. Use sip:{assistantID}:{apiKey}@host"), nil
	}

	auth, err := m.validateAPIKey(ctx.APIKey)
	if err != nil {
		m.logger.Warnw("SIP: invalid API key", "call_id", ctx.CallID, "method", ctx.Method, "error", err)
		return sip_infra.Reject(403, "Invalid API key"), nil
	}

	ctx.Set("auth", auth)
	return next()
}

// assistantMiddleware loads the assistant entity and verifies access.
// Requires Extra["auth"] to be set by authMiddleware.
// Sets Extra["assistant"] for downstream middlewares.
func (m *SIPManager) assistantMiddleware(ctx *sip_infra.SIPRequestContext, next func() (*sip_infra.InviteResult, error)) (*sip_infra.InviteResult, error) {
	authVal, _ := ctx.Get("auth")
	auth, _ := authVal.(types.SimplePrinciple)
	if auth == nil {
		return sip_infra.Reject(401, "Authentication required"), nil
	}

	if ctx.AssistantID == "" {
		return sip_infra.Reject(404, "Invalid SIP URI format, expected: sip:{assistantID}:{apiKey}@host"), nil
	}
	assistantID, err := strconv.ParseUint(ctx.AssistantID, 10, 64)
	if err != nil {
		m.logger.Warnw("SIP: invalid assistant ID", "call_id", ctx.CallID, "method", ctx.Method, "assistant_id", ctx.AssistantID)
		return sip_infra.Reject(404, "Invalid assistant ID format"), nil
	}

	assistant, err := m.assistantService.Get(m.ctx, auth, assistantID, utils.GetVersionDefinition("latest"),
		&internal_services.GetAssistantOption{InjectPhoneDeployment: true})
	if err != nil {
		m.logger.Error("SIP: assistant not found", "call_id", ctx.CallID, "method", ctx.Method, "assistant_id", assistantID, "error", err)
		return sip_infra.Reject(404, "Assistant not found"), nil
	}

	if !m.hasAccessToAssistant(auth, assistant) {
		return sip_infra.Reject(403, "API key does not have access to this assistant"), nil
	}

	ctx.Set("assistant", assistant)
	return next()
}

// vaultConfigResolver is the final handler in the middleware chain.
// It fetches the SIP provider config from vault and returns the InviteResult
// with the resolved config and all middleware-enriched metadata.
func (m *SIPManager) vaultConfigResolver(ctx *sip_infra.SIPRequestContext) (*sip_infra.InviteResult, error) {
	authVal, _ := ctx.Get("auth")
	auth, _ := authVal.(types.SimplePrinciple)
	assistantVal, _ := ctx.Get("assistant")
	assistant, _ := assistantVal.(*internal_assistant_entity.Assistant)

	if auth == nil || assistant == nil {
		return sip_infra.Reject(500, "Middleware chain incomplete"), nil
	}

	// Fetch both SIP config and vault credential from vault
	sipConfig, vaultCred, err := m.fetchSIPConfigAndVaultCredential(auth, assistant)
	if err != nil {
		m.logger.Error("SIP: failed to resolve config", "call_id", ctx.CallID, "method", ctx.Method, "error", err)
		return sip_infra.Reject(500, "Failed to resolve SIP configuration"), nil
	}

	m.logger.Infow("SIP request authenticated",
		"call_id", ctx.CallID,
		"method", ctx.Method,
		"assistant_id", assistant.Id,
		"org_id", *auth.GetCurrentOrganizationId())

	// Pass auth/assistant/config to session via Extra
	return sip_infra.AllowWithExtra(sipConfig, map[string]interface{}{
		"auth":             auth,
		"assistant":        assistant,
		"sip_config":       sipConfig,
		"vault_credential": vaultCred,
	}), nil
}

// validateAPIKey validates the API key as a project-scoped authentication token.
// It strips the "rpd-prj-" prefix (if present) and calls the auth service to
// resolve the project and organization context — exactly like the HTTP/gRPC
// project authenticator middleware does.
func (m *SIPManager) validateAPIKey(apiKey string) (types.SimplePrinciple, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("empty API key")
	}

	// Strip the project key prefix (same as grpc project middleware)
	cleanKey := strings.Replace(apiKey, types.PROJECT_KEY_PREFIX, "", 1)

	// Resolve via auth service (cached, gRPC call to web-api)
	scoped, err := m.authClient.ScopeAuthorize(m.ctx, cleanKey, "project")
	if err != nil {
		return nil, fmt.Errorf("project auth failed: %w", err)
	}

	projectScope := &types.ProjectScope{
		ProjectId:      &scoped.ProjectId,
		OrganizationId: &scoped.OrganizationId,
		Status:         scoped.GetStatus(),
		CurrentToken:   apiKey,
	}

	if !projectScope.IsAuthenticated() {
		return nil, fmt.Errorf("API key is not active (status: %s)", scoped.GetStatus())
	}

	return projectScope, nil
}

// hasAccessToAssistant checks if the auth context has access to the assistant
func (m *SIPManager) hasAccessToAssistant(auth types.SimplePrinciple, assistant *internal_assistant_entity.Assistant) bool {
	// Check if the assistant belongs to the same project/organization
	if auth.GetCurrentProjectId() == nil || assistant.ProjectId == 0 {
		return true // Skip check if project info not available
	}
	return *auth.GetCurrentProjectId() == assistant.ProjectId
}

// handleInvite processes incoming SIP INVITE requests
func (m *SIPManager) handleInvite(session *sip_infra.Session, fromURI, toURI string) error {
	info := session.GetInfo()
	callID := info.CallID

	m.logger.Infow("Incoming SIP INVITE", "from", fromURI, "to", toURI, "call_id", callID, "direction", info.Direction)

	// For outbound calls (answered), use the pre-stored context from the
	// original OutboundCall flow instead of re-resolving and creating a duplicate conversation.
	if info.Direction == sip_infra.CallDirectionOutbound {
		if err := m.handleOutboundAnswered(session, fromURI, toURI); err != nil {
			return err
		}
		return nil
	}

	// Check if session is still alive (may have been terminated by a race with BYE)
	if session.IsEnded() {
		m.logger.Warnw("Session already ended, skipping INVITE handling", "call_id", callID)
		return fmt.Errorf("session already ended")
	}

	// Retrieve middleware-resolved auth & assistant from session metadata
	// (set by the SIP middleware chain and propagated by infra server.go)
	authVal, _ := session.GetMetadata("auth")
	auth, _ := authVal.(types.SimplePrinciple)
	assistantVal, _ := session.GetMetadata("assistant")
	assistant, _ := assistantVal.(*internal_assistant_entity.Assistant)
	sipConfigVal, _ := session.GetMetadata("sip_config")
	sipConfig, _ := sipConfigVal.(*sip_infra.Config)

	if auth == nil || assistant == nil {
		m.logger.Error("SIP session missing auth/assistant metadata (middleware chain issue)",
			"call_id", callID, "has_auth", auth != nil, "has_assistant", assistant != nil)
		return fmt.Errorf("missing auth context on session")
	}

	// Create conversation for inbound call
	callerID := fromURI
	conversation, err := m.assistantConversationService.CreateConversation(
		m.ctx, auth,
		callerID,
		assistant.Id, assistant.AssistantProviderId,
		type_enums.DIRECTION_INBOUND, utils.SIP,
	)
	if err != nil {
		m.logger.Error("Failed to create conversation", "error", err, "call_id", callID)
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	_, _ = m.assistantConversationService.ApplyConversationMetadata(m.ctx, auth, assistant.Id, conversation.Id,
		[]*types.Metadata{types.NewMetadata("sip.caller_uri", fromURI)})

	// Start the call in a goroutine with tenant-specific config
	go m.startCall(m.ctx, session, auth, assistant, conversation, callerID, sipConfig, utils.SIP)

	return nil
}

// handleOutboundAnswered processes the onInvite callback for outbound calls that have been answered.
// It retrieves the assistant and conversation from metadata set during OutboundCall
// instead of re-resolving everything and creating a duplicate conversation.
func (m *SIPManager) handleOutboundAnswered(session *sip_infra.Session, fromURI, toURI string) error {
	callID := session.GetInfo().CallID

	// Retrieve outbound call context from session metadata.
	// All these were stored by telephony.OutboundCall to avoid expensive re-lookups.
	assistantIDVal, ok := session.GetMetadata("assistant_id")
	if !ok {
		return fmt.Errorf("outbound session missing assistant_id metadata")
	}
	conversationIDVal, ok := session.GetMetadata("conversation_id")
	if !ok {
		return fmt.Errorf("outbound session missing conversation_id metadata")
	}
	authVal, ok := session.GetMetadata("auth")
	if !ok {
		return fmt.Errorf("outbound session missing auth metadata")
	}
	sipConfigVal, ok := session.GetMetadata("sip_config")
	if !ok {
		return fmt.Errorf("outbound session missing sip_config metadata")
	}

	assistantID, _ := assistantIDVal.(uint64)
	conversationID, _ := conversationIDVal.(uint64)
	auth, _ := authVal.(types.SimplePrinciple)
	sipConfig, _ := sipConfigVal.(*sip_infra.Config)

	toPhone := ""
	if v, ok := session.GetMetadata("to_phone"); ok {
		toPhone, _ = v.(string)
	}

	m.logger.Infow("Outbound call answered, resolving context",
		"call_id", callID,
		"assistant_id", assistantID,
		"conversation_id", conversationID)

	// Load assistant — still needed because startCall requires the full object.
	// GetConversation runs in parallel to minimise wall-clock time.
	type assistantResult struct {
		assistant *internal_assistant_entity.Assistant
		err       error
	}
	type conversationResult struct {
		conversation *internal_conversation_entity.AssistantConversation
		err          error
	}

	assistantCh := make(chan assistantResult, 1)
	conversationCh := make(chan conversationResult, 1)

	go func() {
		a, err := m.assistantService.Get(m.ctx, auth, assistantID, utils.GetVersionDefinition("latest"),
			&internal_services.GetAssistantOption{InjectPhoneDeployment: true})
		assistantCh <- assistantResult{a, err}
	}()

	go func() {
		c, err := m.assistantConversationService.GetConversation(
			m.ctx, auth, assistantID, conversationID,
			&internal_services.GetConversationOption{InjectMetadata: true},
		)
		conversationCh <- conversationResult{c, err}
	}()

	aRes := <-assistantCh
	if aRes.err != nil {
		return fmt.Errorf("failed to get assistant for outbound call: %w", aRes.err)
	}
	assistant := aRes.assistant

	cRes := <-conversationCh
	conversation := cRes.conversation
	if cRes.err != nil {
		m.logger.Warnw("Could not retrieve existing conversation, creating new one",
			"call_id", callID,
			"conversation_id", conversationID,
			"error", cRes.err)
		var err error
		conversation, err = m.assistantConversationService.CreateConversation(
			m.ctx, auth, toPhone, assistant.Id, assistant.AssistantProviderId,
			type_enums.DIRECTION_OUTBOUND, utils.PhoneCall,
		)
		if err != nil {
			return fmt.Errorf("failed to create fallback conversation: %w", err)
		}
	}

	_, _ = m.assistantConversationService.ApplyConversationMetadata(m.ctx, auth, assistant.Id, conversation.Id,
		[]*types.Metadata{types.NewMetadata("sip.callee_uri", toURI)})

	// Run startCall synchronously for outbound calls. This is called from
	// handleOutboundDialog (which runs in its own goroutine), so blocking here is
	// safe and intentional. Running synchronously prevents the race condition where
	// handleOutboundDialog tears down the session (on BYE) before startCall has
	// registered its context cancellation. onInvite returns only after the call
	// fully ends, keeping the dialog alive in sipgo's cache for the entire duration.
	m.startCall(m.ctx, session, auth, assistant, conversation, toPhone, sipConfig, utils.PhoneCall)

	return nil
}

// startCall starts the SIP conversation with the assistant.
// Multi-tenant: receives config specific to this call/tenant.
//
// State-machine approach: instead of waiting for setup to complete before
// allowing BYE, each step checks session.IsEnded() or session.ByeReceived().
// For outbound calls, startCall runs synchronously from onInvite (inside
// handleOutboundDialog's goroutine), so it owns the session lifecycle.
// When startCall returns, it calls session.End() to signal handleOutboundDialog
// that cleanup can proceed.
func (m *SIPManager) startCall(ctx context.Context, session *sip_infra.Session, auth types.SimplePrinciple, assistant *internal_assistant_entity.Assistant, conversation *internal_conversation_entity.AssistantConversation, callerID string, sipConfig *sip_infra.Config, source utils.RapidaSource) {
	callID := session.GetInfo().CallID
	isOutbound := session.GetInfo().Direction == sip_infra.CallDirectionOutbound

	// For outbound calls, we own the session lifecycle — ensure session.End() is
	// called when we return so handleOutboundDialog can proceed with cleanup.
	if isOutbound {
		defer func() {
			if !session.IsEnded() {
				session.End()
			}
		}()
	}

	// Bail out early if session already ended (e.g., setup failure).
	if session.IsEnded() {
		m.logger.Warnw("Session already ended before startCall", "call_id", callID)
		return
	}

	callCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	tenantID := fmt.Sprintf("%d", *auth.GetCurrentOrganizationId())

	// Register session so handleBye can find it and cancel the context.
	m.mu.Lock()
	m.sessions[callID] = &sip_infra.SIPSession{
		CallID:      callID,
		AssistantID: assistant.Id,
		TenantID:    tenantID,
		Auth:        auth,
		Config:      sipConfig,
		Cancel:      cancel,
	}
	m.mu.Unlock()

	defer func() {
		m.mu.Lock()
		delete(m.sessions, callID)
		m.mu.Unlock()
	}()

	// Check if BYE arrived before we registered. If so, handleBye (sip.go) couldn't
	// find us in m.sessions to cancel callCtx. Cancel it now ourselves.
	select {
	case <-session.ByeReceived():
		m.logger.Infow("BYE was received before startCall registered — cancelling call context",
			"call_id", callID)
		cancel()
	default:
	}

	// Checkpoint: BYE may have arrived while we registered.
	if session.IsEnded() {
		m.logger.Warnw("Session ended during setup (post-register)", "call_id", callID)
		return
	}

	// Create SIP streamer (uses existing session's RTP handler)
	streamer, err := internal_telephony.Telephony(internal_telephony.SIP).
		SipStreamer(
			callCtx, m.logger, session, sipConfig, assistant, conversation, session.GetVaultCredential(),
		)
	if err != nil {
		m.logger.Error("Failed to create SIP streamer", "error", err, "call_id", callID)
		return
	}

	// Checkpoint: BYE may have arrived during streamer creation.
	if session.IsEnded() {
		if closeable, ok := streamer.(io.Closer); ok {
			closeable.Close()
		}
		m.logger.Warnw("Session ended during setup (post-streamer)", "call_id", callID)
		return
	}

	talker, err := internal_adapter.GetTalker(
		utils.PhoneCall,
		callCtx,
		m.cfg,
		m.logger,
		m.postgres,
		m.opensearch,
		m.redis,
		m.storage,
		streamer,
	)
	if err != nil {
		if closeable, ok := streamer.(io.Closer); ok {
			closeable.Close()
		}
		m.logger.Error("Failed to create SIP talker", "error", err, "call_id", callID)
		return
	}

	m.logger.Infow("SIP call started",
		"call_id", callID,
		"assistant_id", assistant.Id,
		"conversation_id", conversation.Id,
		"caller", callerID)

	// talker.Talk blocks for the call duration. If BYE cancels callCtx,
	// Talk will observe it and return.
	if err := talker.Talk(callCtx, auth); err != nil {
		m.logger.Warnw("SIP talker exited", "error", err, "call_id", callID)
	}

	m.logger.Infow("SIP call ended", "call_id", callID)
}

// handleBye processes SIP BYE requests
func (m *SIPManager) handleBye(session *sip_infra.Session) error {
	callID := session.GetInfo().CallID
	m.logger.Infow("SIP BYE received", "call_id", callID)

	// Cancel the call context
	m.mu.Lock()
	if sipSession, exists := m.sessions[callID]; exists {
		if sipSession.Cancel != nil {
			sipSession.Cancel()
		}
		delete(m.sessions, callID)
	}
	m.mu.Unlock()

	return nil
}

// handleCancel processes SIP CANCEL requests
func (m *SIPManager) handleCancel(session *sip_infra.Session) error {
	callID := session.GetInfo().CallID
	m.logger.Infow("SIP CANCEL received", "call_id", callID)

	// Cancel the call context
	m.mu.Lock()
	if sipSession, exists := m.sessions[callID]; exists {
		if sipSession.Cancel != nil {
			sipSession.Cancel()
		}
		delete(m.sessions, callID)
	}
	m.mu.Unlock()

	return nil
}

// HandleIncomingCall processes an incoming SIP INVITE
// This is called when a SIP call arrives for an assistant
func (m *SIPManager) HandleIncomingCall(
	ctx context.Context,
	auth types.SimplePrinciple,
	assistantID uint64,
	callerID string,
	sipConfig *sip_infra.Config,
) error {
	m.logger.Infow("Incoming SIP call",
		"assistant", assistantID,
		"caller", callerID)

	// Load assistant
	assistant, err := m.assistantService.Get(ctx, auth, assistantID, nil, nil)
	if err != nil {
		m.logger.Errorf("Failed to load assistant for SIP call: %v", err)
		return fmt.Errorf("failed to load assistant: %w", err)
	}

	// Create identifier for the conversation

	// Create new conversation for SIP session
	conversation, err := m.assistantConversationService.
		CreateConversation(
			ctx, auth, callerID, assistantID, assistant.AssistantProviderId,
			type_enums.DIRECTION_INBOUND, utils.SIP,
		)
	if err != nil {
		m.logger.Errorf("Failed to create conversation for SIP call: %v", err)
		return fmt.Errorf("failed to create conversation: %w", err)
	}

	// Create SIP streamer
	// TODO: HandleIncomingCall needs a session from SIP server — this path is for webhook-based calls
	// TODO: Also need to pass vault credential here — currently nil for webhook-based calls
	sipCtx, cancel := context.WithCancel(ctx)
	streamer, err := internal_telephony.Telephony(internal_telephony.SIP).SipStreamer(
		sipCtx, m.logger, nil, sipConfig, assistant, conversation, nil,
	)
	if err != nil {
		cancel()
		m.assistantConversationService.ApplyConversationMetrics(ctx, auth, assistantID, conversation.Id,
			[]*types.Metric{{Name: type_enums.STATUS.String(), Value: type_enums.RECORD_FAILED.String(), Description: "SIP setup failed"}})
		m.logger.Errorf("Failed to create SIP streamer: %v", err)
		return fmt.Errorf("failed to create SIP streamer: %w", err)
	}

	// Create talker with SIP source
	talker, err := internal_adapter.GetTalker(
		utils.SIP,
		ctx,
		m.cfg,
		m.logger,
		m.postgres,
		m.opensearch,
		m.redis,
		m.storage,
		streamer,
	)
	if err != nil {
		if closeable, ok := streamer.(io.Closer); ok {
			closeable.Close()
		}
		cancel()
		m.assistantConversationService.ApplyConversationMetrics(ctx, auth, assistantID, conversation.Id,
			[]*types.Metric{{Name: type_enums.STATUS.String(), Value: type_enums.RECORD_FAILED.String(), Description: "Talker creation failed"}})
		m.logger.Errorf("Failed to create SIP talker: %v", err)
		return fmt.Errorf("failed to create talker: %w", err)
	}

	// Store session with tenant-specific config
	tenantID := fmt.Sprintf("%d", *auth.GetCurrentOrganizationId())
	callID := fmt.Sprintf("sip-%s-%d-%d", tenantID, assistantID, conversation.Id)
	m.mu.Lock()
	m.sessions[callID] = &sip_infra.SIPSession{
		CallID:      callID,
		AssistantID: assistantID,
		TenantID:    tenantID,
		Auth:        auth,
		Config:      sipConfig,
		Cancel:      cancel,
	}
	m.mu.Unlock()

	m.logger.Infof("SIP session started: assistant=%d, conversation=%d, caller=%s",
		assistantID, conversation.Id, callerID)

	// Start the conversation in a goroutine
	go func() {
		defer func() {
			m.mu.Lock()
			delete(m.sessions, callID)
			m.mu.Unlock()
			cancel()
		}()

		if err := talker.Talk(sipCtx, auth); err != nil {
			m.logger.Errorf("SIP conversation error: %v", err)
		}

		m.logger.Infof("SIP session ended: assistant=%d, conversation=%d",
			assistantID, conversation.Id)
	}()

	return nil
}

// EndCall terminates an active SIP call
func (m *SIPManager) EndCall(callID string) error {
	m.mu.Lock()
	session, exists := m.sessions[callID]
	if !exists {
		m.mu.Unlock()
		return fmt.Errorf("call not found: %s", callID)
	}
	delete(m.sessions, callID)
	m.mu.Unlock()

	if session.Cancel != nil {
		session.Cancel()
	}

	m.logger.Infow("SIP call ended", "callID", callID)
	return nil
}

// GetActiveCalls returns the number of active SIP calls
func (m *SIPManager) GetActiveCalls() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.sessions)
}

// Stop stops the SIP manager and terminates all active calls
func (m *SIPManager) Stop() {
	if m.cancel != nil {
		m.cancel()
	}

	m.mu.Lock()
	// Stop the shared server
	if m.server != nil {
		m.server.Stop()
		m.server = nil
	}

	// Cancel all active sessions
	for callID, session := range m.sessions {
		if session.Cancel != nil {
			session.Cancel()
		}
		delete(m.sessions, callID)
	}
	m.mu.Unlock()

	m.logger.Infow("SIP Manager stopped")
}

// Close implements the closeable interface for graceful shutdown
func (m *SIPManager) Close(ctx context.Context) error {
	m.Stop()
	return nil
}

// fetchSIPConfigFromVault fetches SIP provider credentials from vault, then overlays
// platform operational settings (port, transport, RTP range) from app config.
// Twilio/providers give: sip_uri, sip_username, sip_password
// Our platform provides: port, transport, rtp_port_range from app config
func (m *SIPManager) fetchSIPConfigFromVault(auth types.SimplePrinciple, assistant *internal_assistant_entity.Assistant) (*sip_infra.Config, error) {
	if assistant.AssistantPhoneDeployment == nil {
		return nil, fmt.Errorf("assistant has no phone deployment configured")
	}

	opts := assistant.AssistantPhoneDeployment.GetOptions()
	credentialID, err := opts.GetUint64("rapida.credential_id")
	if err != nil {
		return nil, fmt.Errorf("no credential_id in phone deployment: %w", err)
	}

	vaultCred, err := m.vaultClient.GetCredential(m.ctx, auth, credentialID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vault credential %d: %w", credentialID, err)
	}

	// Parse provider credentials from vault
	sipConfig, err := GetSIPConfigFromVault(vaultCred)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SIP config from vault: %w", err)
	}

	// Overlay platform operational settings from app config
	if m.cfg.SIPConfig != nil {
		sipConfig.ApplyOperationalDefaults(
			m.cfg.SIPConfig.Port,
			sip_infra.Transport(m.cfg.SIPConfig.Transport),
			m.cfg.SIPConfig.RTPPortRangeStart,
			m.cfg.SIPConfig.RTPPortRangeEnd,
		)
	}

	return sipConfig, nil
}

// fetchSIPConfigAndVaultCredential fetches both the SIP config and the raw vault credential.
// Returns (*sip_infra.Config, *protos.VaultCredential, error)
func (m *SIPManager) fetchSIPConfigAndVaultCredential(auth types.SimplePrinciple, assistant *internal_assistant_entity.Assistant) (*sip_infra.Config, *protos.VaultCredential, error) {
	if assistant.AssistantPhoneDeployment == nil {
		return nil, nil, fmt.Errorf("assistant has no phone deployment configured")
	}

	opts := assistant.AssistantPhoneDeployment.GetOptions()
	credentialID, err := opts.GetUint64("rapida.credential_id")
	if err != nil {
		return nil, nil, fmt.Errorf("no credential_id in phone deployment: %w", err)
	}

	vaultCred, err := m.vaultClient.GetCredential(m.ctx, auth, credentialID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch vault credential %d: %w", credentialID, err)
	}

	// Parse provider credentials from vault
	sipConfig, err := GetSIPConfigFromVault(vaultCred)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse SIP config from vault: %w", err)
	}

	// Overlay platform operational settings from app config
	if m.cfg.SIPConfig != nil {
		sipConfig.ApplyOperationalDefaults(
			m.cfg.SIPConfig.Port,
			sip_infra.Transport(m.cfg.SIPConfig.Transport),
			m.cfg.SIPConfig.RTPPortRangeStart,
			m.cfg.SIPConfig.RTPPortRangeEnd,
		)
	}

	return sipConfig, vaultCred, nil
}

// GetSIPConfigFromVault extracts SIP provider credentials from vault.
// Only parses what a provider (Twilio, etc.) gives:
//
//	sip_uri      - SIP URI (e.g. "sip:pstn.twilio.com") → parsed into Server (and Port if present)
//	sip_username - SIP username for registration
//	sip_password - SIP password for registration
//	sip_server   - (optional) explicit server address, overrides sip_uri
//	sip_realm    - (optional) SIP realm for auth
//	sip_domain   - (optional) SIP domain
//
// Does NOT set operational fields (port, transport, RTP range) — those come from app config.
func GetSIPConfigFromVault(vaultCredential *protos.VaultCredential) (*sip_infra.Config, error) {
	if vaultCredential == nil || vaultCredential.GetValue() == nil {
		return nil, fmt.Errorf("vault credential is required")
	}

	credMap := vaultCredential.GetValue().AsMap()
	cfg := &sip_infra.Config{}

	// Parse sip_uri to extract server and port (e.g. "sip:192.168.1.5:5060")
	if sipURI, ok := credMap["sip_uri"].(string); ok && sipURI != "" {
		server, port, err := parseSIPURI(sipURI)
		if err == nil {
			cfg.Server = server
			if port > 0 {
				cfg.Port = port
			}
		}
	}

	// Explicit sip_server overrides sip_uri
	if server, ok := credMap["sip_server"].(string); ok && server != "" {
		cfg.Server = server
	}

	// Provider credentials
	if username, ok := credMap["sip_username"].(string); ok {
		cfg.Username = username
	}
	if password, ok := credMap["sip_password"].(string); ok {
		cfg.Password = password
	}
	if realm, ok := credMap["sip_realm"].(string); ok {
		cfg.Realm = realm
	}
	if domain, ok := credMap["sip_domain"].(string); ok {
		cfg.Domain = domain
	}

	return cfg, nil
}

// parseSIPURI parses a SIP URI into host and port
// Supports formats: "sip:host:port", "sip:host", "host:port", "host"
func parseSIPURI(uri string) (string, int, error) {
	// Strip sip: or sips: scheme
	raw := uri
	raw = strings.TrimPrefix(raw, "sips:")
	raw = strings.TrimPrefix(raw, "sip:")

	host, portStr, err := net.SplitHostPort(raw)
	if err != nil {
		// No port in URI, treat whole string as host
		return raw, 0, nil
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return host, 0, fmt.Errorf("invalid port in SIP URI %q: %w", uri, err)
	}

	return host, port, nil
}

// GetSIPConfigFromDeployment extracts SIP provider credentials from assistant deployment options.
// Only parses credential fields — operational settings must be applied separately via ApplyOperationalDefaults.
func GetSIPConfigFromDeployment(opts map[string]interface{}) (*sip_infra.Config, error) {
	cfg := &sip_infra.Config{}

	if server, ok := opts["sip_server"].(string); ok {
		cfg.Server = server
	}
	if username, ok := opts["sip_username"].(string); ok {
		cfg.Username = username
	}
	if password, ok := opts["sip_password"].(string); ok {
		cfg.Password = password
	}
	if realm, ok := opts["sip_realm"].(string); ok {
		cfg.Realm = realm
	}
	if domain, ok := opts["sip_domain"].(string); ok {
		cfg.Domain = domain
	}

	return cfg, nil
}

// // SIPCallReceiver handles incoming SIP call webhooks (for SIP trunks that support webhooks)
// // This is similar to telephony providers like Twilio that use webhooks for call events
// // Multi-tenant: Config is passed per-call, resolved from assistant's deployment settings
// // Note: This doesn't start a SIP server - it handles webhook-based calls where the provider manages SIP
// func (cApi *ConversationApi) SIPCallReceiver(ctx context.Context, auth types.SimplePrinciple, assistantID uint64, callerID string, sipConfig *sip_infra.Config) error {
// 	manager := NewSIPManager(cApi)
// 	manager.ctx, manager.cancel = context.WithCancel(ctx)
// 	return manager.HandleIncomingCall(ctx, auth, assistantID, callerID, sipConfig)
// }

// // SIPCallWebhookRequest represents an incoming SIP call webhook
// type SIPCallWebhookRequest struct {
// 	CallID    string                 `json:"call_id"`
// 	From      string                 `json:"from"`
// 	To        string                 `json:"to"`
// 	Direction string                 `json:"direction"`
// 	SIPConfig map[string]interface{} `json:"sip_config,omitempty"`
// }

// // SIPEventWebhookRequest represents a SIP event webhook
// type SIPEventWebhookRequest struct {
// 	CallID    string                 `json:"call_id"`
// 	EventType string                 `json:"event_type"` // answered, hangup, dtmf, etc.
// 	Timestamp string                 `json:"timestamp"`
// 	Data      map[string]interface{} `json:"data,omitempty"`
// }

// // SIPCallWebhook handles incoming SIP call webhooks from SIP trunks
// // POST /v1/talk/sip/call/:assistantId
// // This endpoint is called by SIP providers (Telnyx, SignalWire, etc.) when a call arrives
// func (cApi *ConversationApi) SIPCallWebhook(c *gin.Context) {
// 	auth, isAuthenticated := types.GetAuthPrinciple(c)
// 	if !isAuthenticated {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated request"})
// 		return
// 	}

// 	assistantIdStr := c.Param("assistantId")
// 	assistantId, err := strconv.ParseUint(assistantIdStr, 10, 64)
// 	if err != nil {
// 		logger.Errorf("Invalid assistantId: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assistantId"})
// 		return
// 	}

// 	var req SIPCallWebhookRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		logger.Errorf("Invalid SIP webhook request: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
// 		return
// 	}

// 	// Extract SIP config from request or use defaults
// 	sipConfig, err := GetSIPConfigFromDeployment(req.SIPConfig)
// 	if err != nil {
// 		logger.Errorf("Invalid SIP config: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid SIP configuration"})
// 		return
// 	}

// 	// Handle the incoming call
// 	if err := SIPCallReceiver(c.Request.Context(), auth, assistantId, req.From, sipConfig); err != nil {
// 		logger.Errorf("Failed to handle SIP call: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle call"})
// 		return
// 	}

// 	logger.Infof("SIP call webhook received: assistant=%d, from=%s, callId=%s",
// 		assistantId, req.From, req.CallID)

// 	c.JSON(http.StatusOK, gin.H{
// 		"status":  "accepted",
// 		"call_id": req.CallID,
// 	})
// }

// // SIPEventWebhook handles SIP event webhooks (hangup, dtmf, etc.)
// // POST /v1/talk/sip/event/:assistantId/:conversationId
// func (cApi *ConversationApi) SIPEventWebhook(c *gin.Context) {
// 	auth, isAuthenticated := types.GetAuthPrinciple(c)
// 	if !isAuthenticated {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated request"})
// 		return
// 	}

// 	assistantIdStr := c.Param("assistantId")
// 	assistantId, err := strconv.ParseUint(assistantIdStr, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assistantId"})
// 		return
// 	}

// 	conversationIdStr := c.Param("conversationId")
// 	conversationId, err := strconv.ParseUint(conversationIdStr, 10, 64)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversationId"})
// 		return
// 	}

// 	var req SIPEventWebhookRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		logger.Errorf("Invalid SIP event webhook: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
// 		return
// 	}

// 	logger.Infof("SIP event webhook: assistant=%d, conversation=%d, event=%s",
// 		assistantId, conversationId, req.EventType)

// 	// Process event based on type
// 	switch req.EventType {
// 	case "hangup", "bye":
// 		// Apply end metrics
// 		assistantConversationService.ApplyConversationMetrics(c, auth, assistantId, conversationId,
// 			[]*types.Metric{{Name: type_enums.STATUS.String(), Value: type_enums.RECORD_COMPLETE.String(), Description: "SIP call ended"}})
// 	case "answered":
// 		// Apply connected metrics
// 		assistantConversationService.ApplyConversationMetrics(c, auth, assistantId, conversationId,
// 			[]*types.Metric{{Name: "sip_answered", Value: "true", Description: "SIP call answered"}})
// 	}

// 	c.JSON(http.StatusOK, gin.H{"status": "processed"})
// }
