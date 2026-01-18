// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_adapter_generic

import (
	"context"
	"errors"
	"time"

	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_gorm "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_services "github.com/rapidaai/api/assistant-api/internal/services"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
)

func (dm *GenericRequestor) Assistant() *internal_assistant_entity.Assistant {
	return dm.assistant
}

func (gr *GenericRequestor) Conversation() *internal_conversation_gorm.AssistantConversation {
	return gr.assistantConversation
}

func (gr *GenericRequestor) GetSpeechToTextTransformer() (
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

func (gr *GenericRequestor) GetTextToSpeechTransformer() (*internal_assistant_entity.AssistantDeploymentAudio, error) {
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

func (gr *GenericRequestor) GetAssistant(
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
	return gr.assistantService.Get(gr.ctx, auth, assistantId, versionId, assistantOpts)
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
func (dm *GenericRequestor) Auth() types.SimplePrinciple {
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
func (deb *GenericRequestor) SetAuth(auth types.SimplePrinciple) {
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
func (tc *GenericRequestor) GetMetadata() map[string]interface{} {
	return tc.metadata
}

func (tc *GenericRequestor) onSetMetadata(auth types.SimplePrinciple, mt map[string]interface{}) {
	modified := make(map[string]interface{})
	for k, v := range mt {
		vl, ok := tc.metadata[k]
		if ok && vl == v {
			continue
		}
		tc.metadata[k] = v
		modified[k] = v
	}
	utils.Go(tc.Context(), func() {
		start := time.Now()
		tc.conversationService.ApplyConversationMetadata(
			tc.Context(),
			auth, tc.assistant.Id, tc.assistantConversation.Id, types.NewMetadataList(modified))
		tc.logger.Benchmark("GenericRequestor.SetMetadata", time.Since(start))
	})

}

// for the conversation metrics
// for adding another metrics
// -----------------------------------------------------------------------------
// Metrics Management
// -----------------------------------------------------------------------------
//
// The following methods are responsible for managing metrics associated with
// the GenericRequestor. Metrics provide valuable insights into the
// conversation's performance, usage, and other relevant statistics.
//
// GetMetrics retrieves the current set of metrics associated with this
// conversation. It returns a slice of Metric pointers, allowing the caller
// to access and analyze various aspects of the conversation's performance.
//
// AddMetrics allows for the addition of new metrics to the conversation.
// This method can be used to update or extend the existing set of metrics
// with new data points or measurements.
//
// These methods play a crucial role in monitoring and analyzing the behavior
// and performance of the GenericRequestor, enabling data-driven
// improvements and optimizations.
//
// -----------------------------------------------------------------------------

func (tc *GenericRequestor) onAddMetrics(auth types.SimplePrinciple, metrics ...*types.Metric) {
	utils.Go(tc.ctx, func() {
		start := time.Now()
		_, err := tc.conversationService.ApplyConversationMetrics(
			tc.ctx,
			auth,
			tc.assistant.Id,
			tc.assistantConversation.Id,
			metrics,
		)
		tc.logger.Benchmark("GenericRequestor.AddMetrics", time.Since(start))
		if err != nil {
			tc.logger.Errorf("unable to flush metrics for conversation %+v", err)
		}
	})
}

func (deb *GenericRequestor) onMessageMetric(ctx context.Context, messageId string, metrics []*types.Metric) error {
	if _, err := deb.
		conversationService.ApplyMessageMetrics(ctx, deb.Auth(), deb.Conversation().Id, messageId, metrics); err != nil {
		deb.logger.Errorf("error updating metrics for message: %v", err)
		return err
	}
	return nil
}
