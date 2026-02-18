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
// Shutdown relies on Recv() returning an error (EOF or context-cancelled)
// or a ConversationDisconnection message. All streamer implementations
// guarantee one of these when the connection ends.
func (t *genericRequestor) Talk(_ context.Context, auth types.SimplePrinciple) error {
	var initialized bool
	totalTime := time.Now()
	for {
		req, err := t.streamer.Recv()
		if err != nil {
			if initialized {
				t.OnPacket(
					context.Background(),
					internal_type.ConversationMetricPacket{
						ContextID: t.Conversation().Id,
						Metrics: []*protos.Metric{{
							Name:        type_enums.STATUS.String(),
							Value:       "completed",
							Description: "Status of current conversation",
						}},
					},
					internal_type.ConversationMetricPacket{
						ContextID: t.Conversation().Id,
						Metrics: []*protos.Metric{{
							Name:        type_enums.TIME_TAKEN.String(),
							Value:       fmt.Sprintf("%d", time.Since(totalTime)),
							Description: "Time taken to complete the conversation from the first message received to the end of the conversation.",
						}},
					},
				)
				t.Disconnect(context.Background())
			}
			return nil
		}

		switch payload := req.(type) {
		case *protos.ConversationInitialization:
			t.logger.Infof("talk: received initialization, initialized=%v", initialized)
			if err := t.Connect(t.streamer.Context(), auth, payload); err != nil {
				t.logger.Errorf("unexpected error while connect assistant, might be problem in configuration %+v", err)
				return fmt.Errorf("talking.Connect error: %w", err)
			}
			initialized = true

		case *protos.ConversationConfiguration:
			if initialized {
				utils.Go(t.streamer.Context(), func() {
					t.disconnectSpeechToText(t.streamer.Context())
				})
				utils.Go(t.streamer.Context(), func() {
					t.disconnectTextToSpeech(t.streamer.Context())
				})
				switch payload.StreamMode {
				case protos.StreamMode_STREAM_MODE_TEXT:
					t.messaging.SwitchMode(type_enums.TextMode)
				case protos.StreamMode_STREAM_MODE_AUDIO:
					utils.Go(t.streamer.Context(), func() {
						t.initializeTextToSpeech(t.streamer.Context())
					})
					utils.Go(t.streamer.Context(), func() {
						t.initializeSpeechToText(t.streamer.Context())
					})
					t.messaging.SwitchMode(type_enums.AudioMode)
				}
			}

		case *protos.ConversationUserMessage:
			if initialized {
				switch msg := payload.GetMessage().(type) {
				case *protos.ConversationUserMessage_Audio:
					if err := t.OnPacket(t.streamer.Context(), internal_type.UserAudioPacket{Audio: msg.Audio}); err != nil {
						t.logger.Errorf("error processing user audio: %v", err)
					}
				case *protos.ConversationUserMessage_Text:
					if err := t.OnPacket(t.streamer.Context(), internal_type.UserTextPacket{Text: msg.Text}); err != nil {
						t.logger.Errorf("error processing user text: %v", err)
					}
				default:
					t.logger.Errorf("illegal input from the user %+v", msg)
				}
			}

		case *protos.ConversationMetadata:
			if initialized {
				if err := t.OnPacket(t.streamer.Context(),
					internal_type.ConversationMetadataPacket{
						ContextID: payload.GetAssistantConversationId(),
						Metadata:  payload.GetMetadata(),
					}); err != nil {
					t.logger.Errorf("error while accepting metadata: %v", err)
				}
			}

		case *protos.ConversationMetric:
			if initialized {
				if err := t.OnPacket(t.streamer.Context(),
					internal_type.ConversationMetricPacket{
						ContextID: payload.GetAssistantConversationId(),
						Metrics:   payload.GetMetrics(),
					}); err != nil {
					t.logger.Errorf("error while accepting metrics: %v", err)
				}
			}

		case *protos.ConversationDisconnection:
			if initialized {
				t.OnPacket(context.Background(),
					internal_type.ConversationMetadataPacket{
						ContextID: t.Conversation().Id,
						Metadata: []*protos.Metadata{{
							Key:   "disconnect_reason",
							Value: payload.GetType().String(),
						}},
					},
				)
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
