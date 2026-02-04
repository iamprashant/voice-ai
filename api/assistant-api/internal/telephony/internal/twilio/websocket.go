// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_twilio_telephony

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
	internal_twilio "github.com/rapidaai/api/assistant-api/internal/telephony/internal/twilio/internal"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type twilioWebsocketStreamer struct {
	streamID       string
	streamer       internal_telephony_base.BaseTelephonyStreamer
	logger         commons.Logger
	audioProcessor *internal_twilio.AudioProcessor

	// Output sender state
	outputSenderStarted bool
	outputSenderMu      sync.Mutex
	audioCtx            context.Context
	audioCancel         context.CancelFunc
}

func NewTwilioWebsocketStreamer(logger commons.Logger, connection *websocket.Conn, assistant *internal_assistant_entity.Assistant, conversation *internal_conversation_entity.AssistantConversation, vlt *protos.VaultCredential) internal_type.TelephonyStreamer {
	audioProcessor, err := internal_twilio.NewAudioProcessor(logger)
	if err != nil {
		logger.Error("Failed to create audio processor", "error", err)
		return nil
	}

	tws := &twilioWebsocketStreamer{
		logger:         logger,
		streamID:       "",
		streamer:       internal_telephony_base.NewBaseTelephonyStreamer(logger, connection, assistant, conversation, vlt),
		audioProcessor: audioProcessor,
	}

	// Set up callbacks
	audioProcessor.SetInputAudioCallback(tws.sendProcessedInputAudio)
	audioProcessor.SetOutputChunkCallback(tws.sendAudioChunk)

	return tws
}

// sendProcessedInputAudio is the callback for processed input audio
func (tws *twilioWebsocketStreamer) sendProcessedInputAudio(audio []byte) {
	// This will be called when enough audio has been buffered
	// The audio is already converted to 16kHz linear16
	tws.streamer.LockInputAudioBuffer()
	tws.streamer.InputBuffer().Write(audio)
	tws.streamer.UnlockInputAudioBuffer()
}

// sendAudioChunk sends an audio chunk to Twilio
func (tws *twilioWebsocketStreamer) sendAudioChunk(chunk *internal_twilio.AudioChunk) error {
	if tws.streamID == "" {
		return nil
	}
	return tws.sendTwilioMessage("media", map[string]interface{}{
		"payload": tws.streamer.Encoder().EncodeToString(chunk.Data),
	})
}

func (tws *twilioWebsocketStreamer) Context() context.Context {
	return tws.streamer.Context()
}

func (tws *twilioWebsocketStreamer) Recv() (*protos.AssistantTalkInput, error) {
	if tws.streamer.Connection() == nil {
		return nil, tws.handleError("WebSocket connection is nil", io.EOF)
	}
	_, message, err := tws.streamer.Connection().ReadMessage()
	if err != nil {
		return nil, tws.handleWebSocketError(err)
	}

	var mediaEvent internal_twilio.TwilioMediaEvent
	if err := json.Unmarshal(message, &mediaEvent); err != nil {
		tws.logger.Error("Failed to unmarshal Twilio media event", "error", err.Error())
		return nil, nil
	}
	switch mediaEvent.Event {
	case "connected":
		// Return downstream config (16kHz linear16) for STT/TTS
		downstreamConfig := tws.audioProcessor.GetDownstreamConfig()
		return tws.streamer.CreateConnectionRequest(downstreamConfig, downstreamConfig)
	case "start":
		tws.handleStartEvent(mediaEvent)
		return nil, nil
	case "media":
		return tws.handleMediaEvent(mediaEvent)
	case "stop":
		tws.logger.Info("Twilio stream stopped")
		tws.stopAudioProcessing()
		tws.streamer.Cancel()
		return nil, io.EOF
	default:
		tws.logger.Warn("Unhandled Twilio event", "event", mediaEvent.Event)
		return nil, nil
	}
}

// stopAudioProcessing stops the output sender goroutine
func (tws *twilioWebsocketStreamer) stopAudioProcessing() {
	tws.outputSenderMu.Lock()
	if tws.audioCancel != nil {
		tws.audioCancel()
		tws.audioCancel = nil
	}
	tws.outputSenderMu.Unlock()
}

// startOutputSender starts the consistent audio output sender
func (tws *twilioWebsocketStreamer) startOutputSender() {
	tws.outputSenderMu.Lock()
	defer tws.outputSenderMu.Unlock()

	if tws.outputSenderStarted {
		return
	}

	tws.audioCtx, tws.audioCancel = context.WithCancel(tws.streamer.Context())
	tws.outputSenderStarted = true
	go tws.audioProcessor.RunOutputSender(tws.audioCtx)
}

func (tws *twilioWebsocketStreamer) Send(response *protos.AssistantTalkOutput) error {
	switch data := response.GetData().(type) {
	case *protos.AssistantTalkOutput_Assistant:
		switch content := data.Assistant.Message.(type) {
		case *protos.ConversationAssistantMessage_Audio:
			// Process audio through the audio processor (converts 16kHz -> 8kHz mulaw)
			// The audio will be sent at consistent 20ms intervals by RunOutputSender
			if err := tws.audioProcessor.ProcessOutputAudio(content.Audio); err != nil {
				tws.logger.Error("Failed to process output audio", "error", err.Error())
				return err
			}
		}
	case *protos.AssistantTalkOutput_Interruption:
		if data.Interruption.Type == protos.ConversationInterruption_INTERRUPTION_TYPE_WORD {
			// Clear both input and output buffers
			tws.audioProcessor.ClearInputBuffer()
			tws.audioProcessor.ClearOutputBuffer()

			if err := tws.sendTwilioMessage("clear", nil); err != nil {
				tws.logger.Errorf("Error sending clear command:", err)
			}
		}
	case *protos.AssistantTalkOutput_Directive:
		if data.Directive.GetType() == protos.ConversationDirective_END_CONVERSATION {
			tws.stopAudioProcessing()
			if tws.streamer.GetUuid() != "" {
				client, err := tws.client(tws.streamer.VaultCredential())
				if err != nil {
					tws.logger.Errorf("Error creating Twilio client:", err)
					if err := tws.streamer.Cancel(); err != nil {
						tws.logger.Errorf("Error disconnecting command:", err)
					}
					return nil
				}
				// Set parameters to change status
				params := &openapi.UpdateCallParams{}
				params.SetStatus("completed")
				if _, err := client.Api.UpdateCall(tws.streamer.GetUuid(), params); err != nil {
					tws.logger.Errorf("Error ending Twilio call:", err)
					if err := tws.streamer.Cancel(); err != nil {
						tws.logger.Errorf("Error disconnecting command:", err)
					}
					return nil
				}
			}
			if err := tws.streamer.Cancel(); err != nil {
				tws.logger.Errorf("Error disconnecting command:", err)
			}
		}
	}
	return nil
}

// start event contains streamSid to be used for subsequent media messages
func (tws *twilioWebsocketStreamer) handleStartEvent(mediaEvent internal_twilio.TwilioMediaEvent) {
	tws.streamID = mediaEvent.StreamSid
	// Start the consistent output sender when stream starts
	tws.startOutputSender()
}

func (tws *twilioWebsocketStreamer) handleMediaEvent(mediaEvent internal_twilio.TwilioMediaEvent) (*protos.AssistantTalkInput, error) {
	payloadBytes, err := tws.streamer.Encoder().DecodeString(mediaEvent.Media.Payload)
	if err != nil {
		tws.logger.Warn("Failed to decode media payload", "error", err.Error())
		return nil, nil
	}

	// Process input audio through audio processor (converts mulaw 8kHz -> linear16 16kHz)
	if err := tws.audioProcessor.ProcessInputAudio(payloadBytes); err != nil {
		tws.logger.Debug("Failed to process input audio", "error", err.Error())
		return nil, nil
	}

	// Check if we have enough buffered audio to send downstream
	tws.streamer.LockInputAudioBuffer()
	defer tws.streamer.UnlockInputAudioBuffer()

	if tws.streamer.InputBuffer().Len() > 0 {
		audioRequest := tws.streamer.CreateVoiceRequest(tws.streamer.InputBuffer().Bytes())
		tws.streamer.InputBuffer().Reset()
		return audioRequest, nil
	}

	return nil, nil
}

func (tws *twilioWebsocketStreamer) sendTwilioMessage(
	eventType string,
	mediaData map[string]interface{}) error {
	if tws.streamer.Connection() == nil || tws.streamID == "" {
		return nil
	}
	message := map[string]interface{}{
		"event":     eventType,
		"streamSid": tws.streamID,
	}
	if mediaData != nil {
		message["media"] = mediaData
	}

	twilioMessageJSON, err := json.Marshal(message)
	if err != nil {
		return tws.handleError("Failed to marshal Twilio message", err)
	}

	if err := tws.streamer.Connection().WriteMessage(websocket.TextMessage, twilioMessageJSON); err != nil {
		return tws.handleError("Failed to send message to Twilio", err)
	}

	return nil
}

func (tws *twilioWebsocketStreamer) handleError(message string, err error) error {
	tws.logger.Error(message, "error", err.Error())
	return err
}

func (tws *twilioWebsocketStreamer) handleWebSocketError(err error) error {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		tws.logger.Error("Unexpected websocket close error", "error", err.Error())
	} else {
		tws.logger.Error("Failed to read message from WebSocket", "error", err.Error())
	}
	tws.streamer.Cancel()
	return io.EOF
}

func (tpc *twilioWebsocketStreamer) client(vaultCredential *protos.VaultCredential) (*twilio.RestClient, error) {
	clientParams, err := tpc.clientParam(vaultCredential)
	if err != nil {
		return nil, err
	}
	return twilio.NewRestClientWithParams(*clientParams), nil
}

func (tpc *twilioWebsocketStreamer) clientParam(vaultCredential *protos.VaultCredential) (*twilio.ClientParams, error) {
	accountSid, ok := vaultCredential.GetValue().AsMap()["account_sid"]
	if !ok {
		return nil, fmt.Errorf("illegal vault config accountSid is not found")
	}
	authToken, ok := vaultCredential.GetValue().AsMap()["account_token"]
	if !ok {
		return nil, fmt.Errorf("illegal vault config account_token not found")
	}
	return &twilio.ClientParams{
		Username: accountSid.(string),
		Password: authToken.(string),
	}, nil
}
