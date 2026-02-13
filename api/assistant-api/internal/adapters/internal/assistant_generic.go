// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package adapter_internal

import (
	"context"
	"errors"
	"time"

	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_services "github.com/rapidaai/api/assistant-api/internal/services"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

func (dm *genericRequestor) Assistant() *internal_assistant_entity.Assistant {
	return dm.assistant
}

func (gr *genericRequestor) Conversation() *internal_conversation_entity.AssistantConversation {
	return gr.assistantConversation
}

func (gr *genericRequestor) GetSpeechToTextTransformer() (
	*internal_assistant_entity.AssistantDeploymentAudio,
	error,
) {
	switch gr.source {
	case utils.PhoneCall:
		if a := gr.assistant; a != nil && a.AssistantPhoneDeployment != nil && a.AssistantPhoneDeployment.InputAudio != nil {
			return a.AssistantPhoneDeployment.InputAudio, nil
		}

	case utils.SDK:
		if a := gr.assistant; a != nil && a.AssistantApiDeployment != nil && a.AssistantApiDeployment.InputAudio != nil {
			return a.AssistantApiDeployment.InputAudio, nil
		}

	case utils.WebPlugin:
		if a := gr.assistant; a != nil && a.AssistantWebPluginDeployment != nil && a.AssistantWebPluginDeployment.InputAudio != nil {
			return a.AssistantWebPluginDeployment.InputAudio, nil
		}

	case utils.Debugger:
		if a := gr.assistant; a != nil && a.AssistantDebuggerDeployment != nil && a.AssistantDebuggerDeployment.InputAudio != nil {
			return a.AssistantDebuggerDeployment.InputAudio, nil
		}
	}
	return nil, errors.New("audio is not enabled for the source")
}

func (gr *genericRequestor) GetTextToSpeechTransformer() (*internal_assistant_entity.AssistantDeploymentAudio, error) {
	switch gr.source {
	case utils.PhoneCall:
		if a := gr.assistant; a != nil && a.AssistantPhoneDeployment != nil && a.AssistantPhoneDeployment.OuputAudio != nil {
			return a.AssistantPhoneDeployment.OuputAudio, nil
		}

	case utils.SDK:
		if a := gr.assistant; a != nil && a.AssistantApiDeployment != nil && a.AssistantApiDeployment.OuputAudio != nil {
			return a.AssistantApiDeployment.OuputAudio, nil
		}

	case utils.WebPlugin:
		if a := gr.assistant; a != nil && a.AssistantWebPluginDeployment != nil && a.AssistantWebPluginDeployment.OuputAudio != nil {
			return a.AssistantWebPluginDeployment.OuputAudio, nil
		}

	case utils.Debugger:
		if a := gr.assistant; a != nil && a.AssistantDebuggerDeployment != nil && a.AssistantDebuggerDeployment.OuputAudio != nil {
			return a.AssistantDebuggerDeployment.OuputAudio, nil
		}
	}
	return nil, errors.New("audio is not enabled for the source")
}

func (gr *genericRequestor) GetAssistant(
	ctx context.Context,
	auth types.SimplePrinciple,
	assistantId uint64,
	version string) (*internal_assistant_entity.Assistant, error) {
	versionId := utils.GetVersionDefinition(version)
	assistantOpts := &internal_services.GetAssistantOption{
		InjectTag: false,
		//
		InjectAssistantProvider:      true,
		InjectKnowledgeConfiguration: true,

		// these are very specific for deployment

		InjectTool:          true,
		InjectAnalysis:      true,
		InjectWebhook:       true,
		InjectConversations: false,
	}
	switch gr.source {
	case utils.PhoneCall:
		assistantOpts.InjectPhoneDeployment = true
	case utils.Whatsapp:
		assistantOpts.InjectWhatsappDeployment = true
	case utils.SDK:
		assistantOpts.InjectApiDeployment = true
	case utils.WebPlugin:
		assistantOpts.InjectWebpluginDeployment = true
	case utils.Debugger:
		assistantOpts.InjectDebuggerDeployment = true
	}
	return gr.assistantService.Get(ctx, auth, assistantId, versionId, assistantOpts)
}

/*
 * Auth retrieves the authentication information associated with the debugger.
 *
 * This method returns the SimplePrinciple object that represents the current
 * authentication state of the debugger. The SimplePrinciple typically contains
 * information such as user ID, roles, or any other relevant authentication data.
 *
 * Returns:
 *   - types.SimplePrinciple: The authentication information for the debugger.
 */
func (dm *genericRequestor) Auth() types.SimplePrinciple {
	return dm.auth
}

/*
 * SetAuth sets the authentication information for the debugger.
 *
 * This method allows updating the authentication state of the debugger by
 * providing a new SimplePrinciple object. This is typically used when the
 * authentication state changes, such as after a successful login or when
 * switching users.
 *
 * Parameters:
 *   - auth: types.SimplePrinciple - The new authentication information to set.
 */
func (deb *genericRequestor) SetAuth(auth types.SimplePrinciple) {
	deb.auth = auth
}

/*
 * Metadata Management for Talking Conversations
 * ---------------------------------------------
 * These methods provide functionality to manage metadata associated with
 * a talking conversation. Metadata can be used to store additional
 * information about the conversation that may be useful for processing,
 * analysis, or integration with other systems.
 *
 * GetMetadata(): Retrieves the entire metadata map.
 * AddMetadata(): Adds a single key-value pair to the metadata.
 * SetMetadata(): Replaces the entire metadata map with a new one.
 *
 * Note: Proper use of these methods ensures consistent handling of
 * conversation metadata across the application.
 */
func (tc *genericRequestor) GetMetadata() map[string]interface{} {
	return tc.metadata
}

func (tc *genericRequestor) onSetMetadata(ctx context.Context, auth types.SimplePrinciple, mt map[string]interface{}) {
	modified := make(map[string]interface{})
	for k, v := range mt {
		vl, ok := tc.metadata[k]
		if ok && vl == v {
			continue
		}
		tc.metadata[k] = v
		modified[k] = v
	}
	utils.Go(ctx, func() {
		dbCtx, cancel := context.WithTimeout(context.Background(), dbWriteTimeout)
		defer cancel()
		start := time.Now()
		tc.conversationService.ApplyConversationMetadata(
			dbCtx,
			auth, tc.assistant.Id, tc.assistantConversation.Id, types.NewMetadataList(modified))
		tc.logger.Benchmark("genericRequestor.SetMetadata", time.Since(start))
	})

}

func (tc *genericRequestor) onAddMetadata(ctx context.Context, metadata ...*protos.Metadata) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), dbWriteTimeout)
	defer cancel()
	_, err := tc.conversationService.ApplyConversationMetadata(
		dbCtx,
		tc.auth,
		tc.assistant.Id,
		tc.assistantConversation.Id,
		types.ToMetadatas(metadata),
	)
	if err != nil {
		tc.logger.Errorf("unable to flush metadata for conversation %+v", err)
	}
	return err
}

func (tc *genericRequestor) onAddMetrics(ctx context.Context, metrics ...*protos.Metric) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), dbWriteTimeout)
	defer cancel()
	_, err := tc.conversationService.ApplyConversationMetrics(
		dbCtx,
		tc.auth,
		tc.assistant.Id,
		tc.assistantConversation.Id,
		types.ToMetrics(metrics),
	)
	if err != nil {
		tc.logger.Errorf("unable to flush metrics for conversation %+v", err)
	}
	return err
}

func (deb *genericRequestor) onMessageMetric(ctx context.Context, messageId string, metrics []*protos.Metric) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), dbWriteTimeout)
	defer cancel()
	if _, err := deb.conversationService.ApplyMessageMetrics(dbCtx, deb.Auth(), deb.Conversation().Id, messageId, types.ToMetrics(metrics)); err != nil {
		deb.logger.Errorf("error updating metrics for message: %v", err)
		return err
	}
	return nil
}
