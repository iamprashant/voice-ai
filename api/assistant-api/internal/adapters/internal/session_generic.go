// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

// Package internal_adapter_generic provides the generic adapter implementation
// for managing voice assistant sessions. It handles the complete lifecycle of
// assistant conversations including connection, disconnection, audio streaming,
// and state management.
package adapter_internal

import (
	"context"
	"sync"
	"time"

	internal_audio_recorder "github.com/rapidaai/api/assistant-api/internal/audio/recorder"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_telemetry "github.com/rapidaai/api/assistant-api/internal/telemetry"
	"github.com/rapidaai/pkg/types"
	type_enums "github.com/rapidaai/pkg/types/enums"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// =============================================================================
// Constants
// =============================================================================

const (
	// clientInfoMetadataKey is the metadata key used to store client information.
	clientInfoMetadataKey = "talk.client_information"

	// dbWriteTimeout is the maximum duration allowed for database write operations
	// (inserts, updates, metric flushes). Uses a background context so that writes
	// are not cancelled by the caller's context lifecycle.
	dbWriteTimeout = 1 * time.Second
)

// =============================================================================
// Session Lifecycle Management
// =============================================================================

// Disconnect gracefully terminates an active assistant conversation session.
//
// This method orchestrates the complete disconnection lifecycle by:
//   - Closing all active listeners (speech-to-text transformers)
//   - Closing all active speakers (text-to-speech transformers)
//   - Flushing final conversation metrics (duration, status)
//   - Persisting audio recordings to storage
//   - Exporting telemetry data for analytics
//   - Cleaning up the assistant executor
//   - Stopping any active idle timeout timers
//
// The method executes resource cleanup operations concurrently for optimal
// performance while ensuring all operations complete before returning.
//
// Thread Safety: This method is safe to call from any goroutine, but should
// only be called once per session to avoid duplicate cleanup operations.
func (r *genericRequestor) Disconnect(ctx context.Context) {
	ctx, span, _ := r.Tracer().StartSpan(ctx, utils.AssistantDisconnectStage)
	startTime := time.Now()

	// Phase 1: Close all session resources concurrently
	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	// Close speech-to-text listener
	utils.Go(ctx, func() {
		defer waitGroup.Done()
		if err := r.disconnectSpeechToText(ctx); err != nil {
			r.logger.Tracef(ctx, "failed to close input transformer: %+v", err)
		}

		if err := r.disconnectEndOfSpeech(ctx); err != nil {
			r.logger.Tracef(ctx, "failed to close end of speech: %+v", err)
		}

	})

	// Close text-to-speech speaker
	utils.Go(ctx, func() {
		defer waitGroup.Done()
		if err := r.disconnectTextToSpeech(ctx); err != nil {
			r.logger.Tracef(ctx, "failed to close output transformer: %+v", err)
		}

		if err := r.disconnectTextAggregator(); err != nil {
			r.logger.Tracef(ctx, "failed to close text aggregator: %+v", err)
		}
	})
	waitGroup.Wait()

	// Phase 2: Trigger end-of-conversation hooks
	r.OnEndConversation(ctx)

	// Phase 3: Persist audio recording asynchronously
	r.persistRecording(ctx)

	// Phase 4: Complete the tracing span
	span.EndSpan(ctx, utils.AssistantDisconnectStage)

	// Phase 5: Export telemetry and cleanup
	r.exportTelemetry(ctx)

	// Phase 6: Close assistant executor and stop timers
	r.closeExecutor(ctx)
	r.stopTimers()
	r.logger.Benchmark("session.Disconnect", time.Since(startTime))
}

// Connect establishes a new assistant session or resumes an existing one.
//
// This method serves as the primary entry point for initiating assistant
// conversations. Based on the provided configuration, it either creates
// a new session or resumes an existing conversation.
//
// Parameters:
//   - ctx: Context for cancellation and deadline propagation
//   - auth: Authentication principal containing user/organization credentials
//   - identifier: Unique session identifier (e.g., phone number, client ID)
//   - config: Conversation configuration including assistant details and audio settings
//
// Returns:
//   - error: nil on success, or an error describing the failure reason
//
// The method performs the following steps:
//  1. Validates and creates the request customizer
//  2. Sets authentication context
//  3. Retrieves the assistant configuration
//  4. Routes to either new session creation or session resumption
//
// Example:
//
//	err := requestor.Connect(ctx, auth, "user-123", &protos.ConversationConfiguration{
//	    Assistant: &protos.AssistantDefinition{AssistantId: 1},
//	})
func (r *genericRequestor) Connect(
	ctx context.Context,
	auth types.SimplePrinciple,
	config *protos.ConversationInitialization,
) error {
	ctx, span, _ := r.Tracer().StartSpan(ctx, utils.AssistantConnectStage)
	defer span.EndSpan(ctx, utils.AssistantConnectStage)

	// Set authentication context
	r.SetAuth(auth)

	// Retrieve assistant configuration
	assistant, err := r.GetAssistant(ctx, auth, config.Assistant.AssistantId, config.Assistant.Version)
	if err != nil {
		r.logger.Errorf("failed to retrieve assistant configuration: %+v", err)
		return err
	}

	// Route to appropriate session handler based on conversation ID presence
	if conversationID := config.GetAssistantConversationId(); conversationID > 0 {
		span.AddAttributes(ctx, internal_telemetry.KV{K: "conversation_initiation", V: internal_telemetry.StringValue("resume")}, internal_telemetry.KV{K: "conversation_id", V: internal_telemetry.IntValue(conversationID)})
		return r.resumeSession(ctx, config, assistant)
	}

	span.AddAttributes(ctx, internal_telemetry.KV{K: "conversation_initiation", V: internal_telemetry.StringValue("new")})
	return r.createSession(ctx, config, assistant)
}

// persistRecording saves the audio recording asynchronously.
//
// This method runs in a background goroutine to avoid blocking the
// disconnect flow. Any errors are logged but do not affect the
// disconnection process.
func (r *genericRequestor) persistRecording(ctx context.Context) {
	if r.recorder != nil {
		utils.Go(ctx, func() {
			userAudio, systemAudio, err := r.recorder.Persist()
			if err != nil {
				r.logger.Tracef(ctx, "failed to persist audio recording: %+v", err)
				return
			}
			if err = r.CreateConversationRecording(ctx, userAudio, systemAudio); err != nil {
				r.logger.Tracef(ctx, "failed to create conversation recording record: %+v", err)
			}
		})
	}

}

// exportTelemetry exports conversation telemetry data for analytics and monitoring.
func (r *genericRequestor) exportTelemetry(ctx context.Context) {
	exportOptions := &internal_telemetry.VoiceAgentExportOption{
		AssistantId:              r.assistant.Id,
		AssistantProviderModelId: r.assistant.AssistantProviderId,
		AssistantConversationId:  r.assistantConversation.Id,
	}

	if err := r.tracer.Export(ctx, r.auth, exportOptions); err != nil {
		r.logger.Errorf("failed to export telemetry data: %v", err)
	}
}

// closeExecutor shuts down the assistant executor and releases its resources.
func (r *genericRequestor) closeExecutor(ctx context.Context) {
	if err := r.assistantExecutor.Close(ctx); err != nil {
		r.logger.Errorf("failed to close assistant executor: %v", err)
	}
}

// stopTimers stops all active timers (idle timeout and max session duration).
func (r *genericRequestor) stopTimers() {
	if r.idleTimeoutTimer != nil {
		r.idleTimeoutTimer.Stop()
	}
	if r.maxSessionTimer != nil {
		r.maxSessionTimer.Stop()
	}
}

// =============================================================================
// Connect Helpers
// =============================================================================

// resumeSession delegates to OnResumeSession with extracted configuration values.
func (r *genericRequestor) resumeSession(
	ctx context.Context,
	config *protos.ConversationInitialization,
	assistant *internal_assistant_entity.Assistant,
) error {
	ctx, span, _ := r.Tracer().StartSpan(ctx, utils.AssistantResumeConverstaionStage)
	defer span.EndSpan(ctx, utils.AssistantResumeConverstaionStage)

	// Resume existing conversation
	conversation, err := r.ResumeConversation(ctx, assistant, config)
	if err != nil {
		r.logger.Errorf("failed to resume conversation: %+v", err)
		return err
	}
	// Initialize critical components concurrently
	errGroup, _ := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		if err := r.assistantExecutor.Initialize(ctx, r, config); err != nil {
			r.logger.Tracef(ctx, "failed to initialize executor: %+v", err)
			return err
		}
		return nil
	})

	errGroup.Go(func() error {
		r.notifyConfiguration(ctx, config, conversation, assistant)
		return nil
	})

	errGroup.Go(func() error {
		if err := r.initializeTextAggregator(ctx); err != nil {
			r.logger.Errorf("unable to initialize sentence assembler with error %v", err)
		}
		return nil
	})

	// tts is important and speed up the first message
	errGroup.Go(func() error {
		switch config.StreamMode {
		case protos.StreamMode_STREAM_MODE_TEXT:
			r.messaging.SwitchMode(type_enums.TextMode)
		case protos.StreamMode_STREAM_MODE_AUDIO:
			r.initializeTextToSpeech(ctx)
			r.initializeBehavior(ctx)
			r.messaging.SwitchMode(type_enums.AudioMode)
		}
		return nil
	})

	utils.Go(ctx, func() {
		switch config.StreamMode {
		case protos.StreamMode_STREAM_MODE_AUDIO:
			r.initializeSpeechToText(ctx)
		}
	})

	// Start non-critical background tasks (not a new session)

	// Initialize audio recorder when both input and output are configured
	utils.Go(ctx, func() {
		rc, err := internal_audio_recorder.GetRecorder(r.logger)
		if err != nil {
			r.logger.Tracef(ctx, "failed to initialize audio recorder: %+v", err)
			return
		}
		r.recorder = rc
	})

	// Establish speech-to-text listener connection
	utils.Go(ctx, func() {
		if err := r.initializeEndOfSpeech(ctx); err != nil {
			r.logger.Tracef(ctx, "failed to initialize input: %+v", err)
		}
	})

	// Update conversation status metric
	utils.Go(ctx, func() {
		r.onAddMetrics(ctx, &protos.Metric{
			Name:        type_enums.STATUS.String(),
			Value:       type_enums.RECORD_IN_PROGRESS.String(),
			Description: "Conversation is currently in progress",
		})
	})

	// Store client information from gRPC context
	utils.Go(ctx, func() {
		r.storeClientInformation(ctx)
	})

	// Trigger begin conversation hooks for new sessions only
	utils.Go(ctx, func() {
		if err := r.OnResumeConversation(ctx); err != nil {
			r.logger.Errorf("failed to execute resume conversation hooks: %v", err)
		}
	})

	return errGroup.Wait()
}

// createSession delegates to OnCreateSession with extracted configuration values.
func (r *genericRequestor) createSession(
	ctx context.Context,
	config *protos.ConversationInitialization,
	assistant *internal_assistant_entity.Assistant,
) error {
	ctx, span, _ := r.Tracer().StartSpan(ctx, utils.AssistantCreateConversationStage)
	defer span.EndSpan(ctx, utils.AssistantCreateConversationStage)

	//
	conversation, err := r.BeginConversation(
		ctx,
		assistant,
		type_enums.DIRECTION_INBOUND,
		config,
	)
	if err != nil {
		r.logger.Errorf("failed to begin conversation: %+v", err)
		return err
	}

	// Initialize critical components concurrently
	errGroup, _ := errgroup.WithContext(ctx)

	// blocking initialization of assistant executor to ensure it's ready before processing any input or output
	errGroup.Go(func() error {
		if err := r.assistantExecutor.Initialize(ctx, r, config); err != nil {
			r.logger.Tracef(ctx, "failed to initialize executor: %+v", err)
			return err
		}
		return nil
	})

	// blocking initialization of assistant executor to ensure it's ready before processing any input or output
	errGroup.Go(func() error {
		r.notifyConfiguration(ctx, config, conversation, assistant)
		return nil
	})

	errGroup.Go(func() error {
		switch config.StreamMode {
		case protos.StreamMode_STREAM_MODE_TEXT:
			r.messaging.SwitchMode(type_enums.TextMode)
		case protos.StreamMode_STREAM_MODE_AUDIO:
			r.initializeTextToSpeech(ctx)
			r.initializeBehavior(ctx)
			r.messaging.SwitchMode(type_enums.AudioMode)
		}
		return nil
	})

	// blocking initialization of assistant executor to ensure it's ready before processing any input or output
	errGroup.Go(func() error {
		if err := r.initializeTextAggregator(ctx); err != nil {
			r.logger.Errorf("unable to initialize sentence assembler with error %v", err)
		}
		return nil
	})

	utils.Go(ctx, func() {
		switch config.StreamMode {
		case protos.StreamMode_STREAM_MODE_AUDIO:
			r.initializeSpeechToText(ctx)
		}
	})

	// Start non-critical background tasks
	// Initialize audio recorder when both input and output are configured
	utils.Go(ctx, func() {
		rc, err := internal_audio_recorder.GetRecorder(r.logger)
		if err != nil {
			r.logger.Tracef(ctx, "failed to initialize audio recorder: %+v", err)
			return
		}
		r.recorder = rc
	})

	utils.Go(ctx, func() {
		if err := r.initializeEndOfSpeech(ctx); err != nil {
			r.logger.Tracef(ctx, "failed to initialize input: %+v", err)
		}

	})

	// Update conversation status metric
	utils.Go(ctx, func() {
		r.onAddMetrics(ctx, &protos.Metric{
			Name:        type_enums.STATUS.String(),
			Value:       type_enums.RECORD_IN_PROGRESS.String(),
			Description: "Conversation is currently in progress",
		})
	})

	// Store client information from gRPC context
	utils.Go(ctx, func() {
		r.storeClientInformation(ctx)
	})

	// Trigger begin conversation hooks for new sessions only
	utils.Go(ctx, func() {
		if err := r.OnBeginConversation(ctx); err != nil {
			r.logger.Errorf("failed to execute begin conversation hooks: %+v", err)
		}
	})
	return errGroup.Wait()
}

// notifyConfiguration sends the initial conversation configuration to the client.
//
// This notification includes the conversation ID and assistant details,
// allowing the client to track the session.
func (r *genericRequestor) notifyConfiguration(
	ctx context.Context,
	config *protos.ConversationInitialization,
	conversation *internal_conversation_entity.AssistantConversation,
	assistant *internal_assistant_entity.Assistant,
) {
	if err := r.Notify(ctx, &protos.ConversationInitialization{
		AssistantConversationId: conversation.Id,
		Assistant: &protos.AssistantDefinition{
			AssistantId: assistant.Id,
			Version:     utils.GetVersionString(assistant.AssistantProviderId),
		},
		Args:         config.GetArgs(),
		Metadata:     config.GetOptions(),
		Options:      config.GetMetadata(),
		StreamMode:   config.GetStreamMode(),
		UserIdentity: config.GetUserIdentity(),
		Time:         timestamppb.Now(),
	}); err != nil {
		r.logger.Errorf("failed to send configuration notification: %v", err)
	}
}

// storeClientInformation extracts client metadata from the gRPC context
// and persists it as conversation metadata for analytics purposes.
func (r *genericRequestor) storeClientInformation(ctx context.Context) {
	clientInfo := types.GetClientInfoFromGrpcContext(ctx)
	if clientInfo == nil {
		return
	}
	clientJSON, err := clientInfo.ToJson()
	if err != nil {
		r.logger.Tracef(ctx, "failed to serialize client information: %+v", err)
		return
	}
	r.onSetMetadata(ctx, r.Auth(), map[string]interface{}{
		clientInfoMetadataKey: clientJSON,
	})
}
