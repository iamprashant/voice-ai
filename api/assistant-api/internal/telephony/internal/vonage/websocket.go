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
	"sync"

	"github.com/gorilla/websocket"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_telephony_base "github.com/rapidaai/api/assistant-api/internal/telephony/internal/base"
	internal_vonage "github.com/rapidaai/api/assistant-api/internal/telephony/internal/vonage/internal"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	protos "github.com/rapidaai/protos"
	"github.com/vonage/vonage-go-sdk"
)

type vonageWebsocketStreamer struct {
	streamer       internal_telephony_base.BaseTelephonyStreamer
	logger         commons.Logger
	audioProcessor *internal_vonage.AudioProcessor

	// Output sender state
	outputSenderStarted bool
	outputSenderMu      sync.Mutex
	audioCtx            context.Context
	audioCancel         context.CancelFunc
}

func NewVonageWebsocketStreamer(logger commons.Logger, connection *websocket.Conn, assistant *internal_assistant_entity.Assistant, conversation *internal_conversation_entity.AssistantConversation, vlt *protos.VaultCredential) internal_type.TelephonyStreamer {
	audioProcessor, err := internal_vonage.NewAudioProcessor(logger)
	if err != nil {
		logger.Error("Failed to create audio processor", "error", err)
		return nil
	}

	vng := &vonageWebsocketStreamer{
		logger:         logger,
		streamer:       internal_telephony_base.NewBaseTelephonyStreamer(logger, connection, assistant, conversation, vlt),
		audioProcessor: audioProcessor,
	}

	// Set up callbacks
	audioProcessor.SetInputAudioCallback(vng.sendProcessedInputAudio)
	audioProcessor.SetOutputChunkCallback(vng.sendAudioChunk)

	return vng
}

// sendProcessedInputAudio is the callback for processed input audio
func (vng *vonageWebsocketStreamer) sendProcessedInputAudio(audio []byte) {
	// This will be called when enough audio has been buffered
	// The audio is already in 16kHz linear16 format
	vng.streamer.LockInputAudioBuffer()
	vng.streamer.InputBuffer().Write(audio)
	vng.streamer.UnlockInputAudioBuffer()
}

// sendAudioChunk sends an audio chunk to Vonage
func (vng *vonageWebsocketStreamer) sendAudioChunk(chunk *internal_vonage.AudioChunk) error {
	if vng.streamer.Connection() == nil {
		return nil
	}
	return vng.streamer.Connection().WriteMessage(websocket.BinaryMessage, chunk.Data)
}

// stopAudioProcessing stops the output sender goroutine
func (vng *vonageWebsocketStreamer) stopAudioProcessing() {
	vng.outputSenderMu.Lock()
	if vng.audioCancel != nil {
		vng.audioCancel()
		vng.audioCancel = nil
	}
	vng.outputSenderMu.Unlock()
}

// startOutputSender starts the consistent audio output sender
func (vng *vonageWebsocketStreamer) startOutputSender() {
	vng.outputSenderMu.Lock()
	defer vng.outputSenderMu.Unlock()

	if vng.outputSenderStarted {
		return
	}

	vng.audioCtx, vng.audioCancel = context.WithCancel(vng.streamer.Context())
	vng.outputSenderStarted = true
	go vng.audioProcessor.RunOutputSender(vng.audioCtx)
}

func (vng *vonageWebsocketStreamer) Context() context.Context {
	return vng.streamer.Context()
}

func (vng *vonageWebsocketStreamer) Recv() (*protos.AssistantTalkInput, error) {
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
			// Start the consistent output sender when connected
			vng.startOutputSender()
			// Return downstream config (16kHz linear16) for STT/TTS
			downstreamConfig := vng.audioProcessor.GetDownstreamConfig()
			return vng.streamer.CreateConnectionRequest(downstreamConfig, downstreamConfig)

		case "stop":
			vng.stopAudioProcessing()
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

func (vng *vonageWebsocketStreamer) Send(response *protos.AssistantTalkOutput) error {
	if vng.streamer.Connection() == nil {
		return nil
	}
	switch data := response.GetData().(type) {
	case *protos.AssistantTalkOutput_Assistant:
		switch content := data.Assistant.Message.(type) {
		case *protos.ConversationAssistantMessage_Audio:
			// Process audio through the audio processor
			// The audio will be sent at consistent 20ms intervals by RunOutputSender
			if err := vng.audioProcessor.ProcessOutputAudio(content.Audio); err != nil {
				vng.logger.Error("Failed to process output audio", "error", err.Error())
				return err
			}
		}
	case *protos.AssistantTalkOutput_Interruption:
		if data.Interruption.Type == protos.ConversationInterruption_INTERRUPTION_TYPE_WORD {
			// Clear both input and output buffers
			vng.audioProcessor.ClearInputBuffer()
			vng.audioProcessor.ClearOutputBuffer()

			// Send clear command to Vonage
			if err := vng.streamer.Connection().WriteMessage(websocket.TextMessage, []byte(`{"action":"clear"}`)); err != nil {
				vng.logger.Errorf("Error sending clear command:", err)
			}
		}
	case *protos.AssistantTalkOutput_Directive:
		if data.Directive.GetType() == protos.ConversationDirective_END_CONVERSATION {
			vng.stopAudioProcessing()
			if vng.streamer.GetUuid() != "" {
				cAuth, err := vng.Auth(vng.streamer.VaultCredential())
				if err != nil {
					vng.logger.Errorf("Error creating Vonage client:", err)
					if err := vng.streamer.Cancel(); err != nil {
						vng.logger.Errorf("Error disconnecting command:", err)
					}
					return nil
				}

				if _, _, err := vonage.NewVoiceClient(cAuth).Hangup(vng.streamer.GetUuid()); err != nil {
					vng.logger.Errorf("Error ending Vonage call:", err)
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

func (vng *vonageWebsocketStreamer) handleMediaEvent(message []byte) (*protos.AssistantTalkInput, error) {
	// Process input audio through audio processor
	// Vonage sends linear16 16kHz, which matches downstream format
	if err := vng.audioProcessor.ProcessInputAudio(message); err != nil {
		vng.logger.Debug("Failed to process input audio", "error", err.Error())
		return nil, nil
	}

	// Check if we have enough buffered audio to send downstream
	vng.streamer.LockInputAudioBuffer()
	defer vng.streamer.UnlockInputAudioBuffer()

	if vng.streamer.InputBuffer().Len() > 0 {
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
