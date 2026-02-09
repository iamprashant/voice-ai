// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package adapter_internal

import (
	"context"
	"fmt"
	"time"

	internal_adapter_telemetry "github.com/rapidaai/api/assistant-api/internal/telemetry"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/types"
	type_enums "github.com/rapidaai/pkg/types/enums"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

// tracerSpan is a type alias for the tracer interface used as span.
type tracerSpan = internal_adapter_telemetry.Tracer[utils.RapidaStage]

// =============================================================================
// Talk - Main Entry Point
// =============================================================================

// Talk handles the main conversation loop for different streamer types.
// It processes incoming messages and manages the connection lifecycle.
func (t *genericRequestor) Talk(ctx context.Context, auth types.SimplePrinciple) error {
	t.StartedAt = time.Now()
	var initialized bool
	for {
		select {
		case <-ctx.Done():
			if initialized {
				t.Disconnect()
			}
			return ctx.Err()
		default:
			req, err := t.streamer.Recv()
			if err != nil {
				continue
			}

			//
			switch payload := req.(type) {
			case *protos.ConversationInitialization:
				if err := t.Connect(ctx, auth, payload); err != nil {
					t.logger.Errorf("unexpected error while connect assistant, might be problem in configuration %+v", err)
					return fmt.Errorf("talking.Connect error: %w", err)
				}
				initialized = true
			case *protos.ConversationConfiguration:
				if initialized {
					switch payload.GetStreamMode() {
					case protos.StreamMode_STREAM_MODE_TEXT:
						utils.Go(ctx, func() {
							t.disconnectSpeechToText(ctx)
						})
						utils.Go(ctx, func() {
							t.disconnectTextToSpeech(ctx)
						})
						t.messaging.SwitchMode(type_enums.TextMode)
					case protos.StreamMode_STREAM_MODE_AUDIO:
						utils.Go(ctx, func() {
							t.initializeTextToSpeech(ctx)
						})
						utils.Go(ctx, func() {
							t.initializeSpeechToText(ctx)
						})
						t.messaging.SwitchMode(type_enums.AudioMode)
					}
				}
			case *protos.ConversationUserMessage:
				if initialized {
					if err := t.Input(payload); err != nil {
						t.logger.Errorf("error while accepting input %v", err)
					}
				}
			case *protos.ConversationMetadata:
				// might be used for future enhancements
			case *protos.ConversationMerics:
				// Handle metrics if needed
			case *protos.ConversationDisconnection:
				// Handle disconnection if needed
			}
		}
	}
}

// // Notify sends notifications to websocket for various events.
func (t *genericRequestor) Notify(ctx context.Context, actionDatas ...internal_type.Stream) error {
	ctx, span, _ := t.Tracer().StartSpan(ctx, utils.AssistantNotifyStage)
	defer span.EndSpan(ctx, utils.AssistantNotifyStage)
	for _, actionData := range actionDatas {
		t.streamer.Send(actionData)
	}
	return nil
}
