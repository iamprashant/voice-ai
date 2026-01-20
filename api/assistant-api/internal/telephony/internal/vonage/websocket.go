// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_vonage_telephony

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gorilla/websocket"
	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_streamers "github.com/rapidaai/api/assistant-api/internal/streamers"
	internal_telephony_base "github.com/rapidaai/api/assistant-api/internal/telephony/internal/base"
	"github.com/rapidaai/pkg/commons"
	protos "github.com/rapidaai/protos"
	"github.com/vonage/vonage-go-sdk"
)

type vonageWebsocketStreamer struct {
	streamer internal_telephony_base.BaseTelephonyStreamer
	logger   commons.Logger
}

func NewVonageWebsocketStreamer(logger commons.Logger, connection *websocket.Conn, assistant *internal_assistant_entity.Assistant, conversation *internal_conversation_entity.AssistantConversation, vlt *protos.VaultCredential) internal_streamers.Streamer {
	return &vonageWebsocketStreamer{
		logger:   logger,
		streamer: internal_telephony_base.NewBaseTelephonyStreamer(logger, connection, assistant, conversation, vlt),
	}
}

func (vng *vonageWebsocketStreamer) Context() context.Context {
	return vng.streamer.Context()
}

func (vng *vonageWebsocketStreamer) Recv() (*protos.AssistantMessagingRequest, error) {
	if vng.streamer.Connection() == nil {
		return nil, vng.handleError("WebSocket connection is nil", io.EOF)
	}
	messageType, message, err := vng.streamer.Connection().ReadMessage()
	if err != nil {
		return nil, vng.handleWebSocketError(err)
	}
	switch messageType {
	case websocket.TextMessage:
		var textEvent map[string]interface{}
		if err := json.Unmarshal(message, &textEvent); err != nil {
			vng.logger.Error("Failed to unmarshal text event", "error", err.Error())
			return nil, err
		}
		switch textEvent["event"] {
		case "websocket:connected":
			return vng.streamer.CreateConnectionRequest(internal_audio.NewLinear16khzMonoAudioConfig(), internal_audio.NewLinear16khzMonoAudioConfig())

		case "stop":
			return nil, io.EOF

		default:
			vng.logger.Debugf("Unhandled event type: %s", textEvent["event"])
		}

	case websocket.BinaryMessage:
		return vng.handleMediaEvent(message)
	default:
		vng.logger.Warn("Unhandled message type", "type", messageType)
		return nil, nil
	}
	return nil, nil
}

func (vng *vonageWebsocketStreamer) Send(response *protos.AssistantMessagingResponse) error {
	if vng.streamer.Connection() == nil {
		return nil
	}
	switch data := response.GetData().(type) {
	case *protos.AssistantMessagingResponse_Assistant:
		switch content := data.Assistant.Message.(type) {
		case *protos.AssistantConversationAssistantMessage_Audio:
			//	1ms 32  10ms 320byte @ 16000Hz, 16-bit mono PCM = 640 bytes
			// Each message needs to be a 20ms sample of audio.
			// At 8kHz the message should be 320 bytes.
			// At 16kHz the message should be 640 bytes.
			bufferSizeThreshold := 32 * 20
			audioData := content.Audio.GetContent()

			// Use vng.audioBuffer to handle pending data across calls
			vng.streamer.LockOutputAudioBuffer()
			defer vng.streamer.UnlockOutputAudioBuffer()

			// Append incoming audio data to the buffer
			vng.streamer.OutputBuffer().Write(audioData)
			// Process and send chunks of 640 bytes
			for vng.streamer.OutputBuffer().Len() >= bufferSizeThreshold {
				chunk := vng.streamer.OutputBuffer().Next(bufferSizeThreshold) // Get and remove the next 640 bytes
				if err := vng.streamer.Connection().WriteMessage(websocket.BinaryMessage, chunk); err != nil {
					vng.logger.Error("Failed to send audio chunk", "error", err.Error())
					return err
				}
			}

			// If response is marked as completed, flush any remaining audio in the buffer
			if data.Assistant.GetCompleted() && vng.streamer.OutputBuffer().Len() > 0 {
				remainingChunk := vng.streamer.OutputBuffer().Bytes()
				if err := vng.streamer.Connection().WriteMessage(websocket.BinaryMessage, remainingChunk); err != nil {
					vng.logger.Errorf("Failed to send final audio chunk", "error", err.Error())
					return err
				}
				vng.streamer.OutputBuffer().Reset() // Clear the buffer after flushing
			}
		}
	case *protos.AssistantMessagingResponse_Interruption:
		if data.Interruption.Type == protos.AssistantConversationInterruption_INTERRUPTION_TYPE_WORD {
			vng.streamer.LockOutputAudioBuffer()
			vng.streamer.OutputBuffer().Reset()
			vng.streamer.UnlockOutputAudioBuffer()

			// Clear the buffer after flushing
			if err := vng.streamer.Connection().WriteMessage(websocket.TextMessage, []byte(`{"action":"clear"}`)); err != nil {
				vng.logger.Errorf("Error sending clear command:", err)
			}
		}
	case *protos.AssistantMessagingResponse_Action:
		if data.Action.GetAction() == protos.AssistantConversationAction_END_CONVERSATION {
			if vng.streamer.GetUuid() != "" {
				cAuth, err := vng.Auth(vng.streamer.VaultCredential())
				if err != nil {
					vng.logger.Errorf("Error creating Twilio client:", err)
					if err := vng.streamer.Cancel(); err != nil {
						vng.logger.Errorf("Error disconnecting command:", err)
					}
					return nil
				}

				if _, _, err := vonage.NewVoiceClient(cAuth).Hangup(vng.streamer.GetUuid()); err != nil {
					vng.logger.Errorf("Error ending Twilio call:", err)
					if err := vng.streamer.Cancel(); err != nil {
						vng.logger.Errorf("Error disconnecting command:", err)
					}
					return nil
				}
			}
			if err := vng.streamer.Cancel(); err != nil {
				vng.logger.Errorf("Error disconnecting command:", err)
			}
		} else {
			if err := vng.streamer.Cancel(); err != nil {
				vng.logger.Errorf("Error disconnecting command:", err)
			}
		}
	}
	return nil
}

func (vng *vonageWebsocketStreamer) handleMediaEvent(message []byte) (*protos.AssistantMessagingRequest, error) {
	vng.streamer.LockInputAudioBuffer()
	defer vng.streamer.UnlockInputAudioBuffer()

	vng.streamer.InputBuffer().Write(message)
	const bufferSizeThreshold = 32 * 60

	if vng.streamer.InputBuffer().Len() >= bufferSizeThreshold {
		audioRequest := vng.streamer.CreateVoiceRequest(vng.streamer.InputBuffer().Bytes())
		vng.streamer.InputBuffer().Reset()
		return audioRequest, nil
	}
	return nil, nil
}

func (vng *vonageWebsocketStreamer) handleError(message string, err error) error {
	vng.logger.Error(message, "error", err.Error())
	return err
}
func (vng *vonageWebsocketStreamer) handleWebSocketError(err error) error {
	vng.streamer.Cancel()
	return io.EOF
}

func (tpc *vonageWebsocketStreamer) Auth(vaultCredential *protos.VaultCredential) (vonage.Auth, error) {
	privateKey, ok := vaultCredential.GetValue().AsMap()["private_key"]
	if !ok {
		return nil, fmt.Errorf("illegal vault config privateKey is not found")
	}
	applicationId, ok := vaultCredential.GetValue().AsMap()["application_id"]
	if !ok {
		return nil, fmt.Errorf("illegal vault config application_id is not found")
	}
	clientAuth, err := vonage.CreateAuthFromAppPrivateKey(applicationId.(string), []byte(privateKey.(string)))
	if err != nil {
		return nil, err
	}
	return clientAuth, nil
}
