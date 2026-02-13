// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package adapter_internal

import (
	"context"
	"fmt"
	"io"

	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/types"
	type_enums "github.com/rapidaai/pkg/types/enums"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

// =============================================================================
// Talk - Main Entry Point
// =============================================================================

// Talk handles the main conversation loop for different streamer types.
// It processes incoming messages and manages the connection lifecycle.
//
// Context management:
//   - The streamer owns its own context (returned by streamer.Context()) which
//     is a child of the gRPC stream context. When the streamer closes for any
//     reason (client disconnect, connection failure, END_CONVERSATION directive)
//     it cancels its context, which causes this loop to exit cleanly.
//   - We use the streamer's context so that all disconnect paths — whether
//     initiated by the client, the server, or the network — converge on a
//     single shutdown signal.
func (t *genericRequestor) Talk(streamerCtx context.Context, auth types.SimplePrinciple) error {
	var initialized bool
	for {
		select {

		// if context is done, it means streamer has closed (either client disconnect, connection failure, or END_CONVERSATION directive)
		case <-streamerCtx.Done():
			t.logger.Infof("talk loop exiting, streamer context done, initialized=%v", initialized)
			if initialized {
				t.Disconnect(context.Background())
			}
			return streamerCtx.Err()

		default:
			req, err := t.streamer.Recv()
			if err != nil {
				// if the error is EOF, it means the streamer closed the connection gracefully (e.g. client disconnect or END_CONVERSATION directive)
				if err == io.EOF {
					t.logger.Infof("talk loop: streamer returned EOF, initialized=%v", initialized)
					if initialized {
						t.Disconnect(context.Background())
					}
					return nil
				}
				continue
			}

			switch payload := req.(type) {
			case *protos.ConversationInitialization:
				t.logger.Infof("talk: received initialization, initialized=%v", initialized)
				if err := t.Connect(streamerCtx, auth, payload); err != nil {
					t.logger.Errorf("unexpected error while connect assistant, might be problem in configuration %+v", err)
					return fmt.Errorf("talking.Connect error: %w", err)
				}
				initialized = true

			case *protos.ConversationConfiguration:
				t.logger.Infof("talk: received configuration, initialized=%v", initialized)
				if initialized {
					switch payload.GetStreamMode() {
					case protos.StreamMode_STREAM_MODE_TEXT:
						utils.Go(streamerCtx, func() {
							t.disconnectSpeechToText(streamerCtx)
						})
						utils.Go(streamerCtx, func() {
							t.disconnectTextToSpeech(streamerCtx)
						})
						t.messaging.SwitchMode(type_enums.TextMode)
					case protos.StreamMode_STREAM_MODE_AUDIO:
						utils.Go(streamerCtx, func() {
							t.logger.Debugf("connecting text to speech")
							t.initializeTextToSpeech(streamerCtx)
						})
						utils.Go(streamerCtx, func() {
							t.logger.Debugf("connecting speech to text")
							t.initializeSpeechToText(streamerCtx)
						})
						t.messaging.SwitchMode(type_enums.AudioMode)
					}
				}

			case *protos.ConversationUserMessage:
				if initialized {
					switch msg := payload.GetMessage().(type) {
					case *protos.ConversationUserMessage_Audio:
						if err := t.OnPacket(streamerCtx, internal_type.UserAudioPacket{Audio: msg.Audio}); err != nil {
							t.logger.Errorf("error processing user audio: %v", err)
						}
					case *protos.ConversationUserMessage_Text:
						if err := t.OnPacket(streamerCtx, internal_type.UserTextPacket{Text: msg.Text}); err != nil {
							t.logger.Errorf("error processing user text: %v", err)
						}
					default:
						t.logger.Errorf("illegal input from the user %+v", msg)
					}
				}

			case *protos.ConversationMetadata:
				if initialized {
					if err := t.OnPacket(streamerCtx,
						internal_type.ConversationMetadataPacket{
							ContextID: payload.GetAssistantConversationId(),
							Metadata:  payload.GetMetadata(),
						}); err != nil {
						t.logger.Errorf("error while accepting metadata: %v", err)
					}
				}

			case *protos.ConversationMetric:
				if initialized {
					if err := t.OnPacket(streamerCtx,
						internal_type.ConversationMetricPacket{
							ContextID: payload.GetAssistantConversationId(),
							Metrics:   payload.GetMetrics(),
						}); err != nil {
						t.logger.Errorf("error while accepting metrics: %v", err)
					}
				}

			case *protos.ConversationDirective:
				// The streamer pushes END_CONVERSATION back into inputCh
				// so the adapter can run cleanup before the context dies.
				if initialized && payload.GetType() == protos.ConversationDirective_END_CONVERSATION {
					t.logger.Infof("talk: received END_CONVERSATION directive, disconnecting")
					t.Disconnect(context.Background())
					return nil
				}
			}
		}
	}
}

// Notify sends notifications to websocket for various events.
func (t *genericRequestor) Notify(ctx context.Context, actionDatas ...internal_type.Stream) error {
	ctx, span, _ := t.Tracer().StartSpan(ctx, utils.AssistantNotifyStage)
	defer span.EndSpan(ctx, utils.AssistantNotifyStage)
	for _, actionData := range actionDatas {
		t.streamer.Send(actionData)
	}
	return nil
}
