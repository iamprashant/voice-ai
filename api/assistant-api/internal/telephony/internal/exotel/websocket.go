// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_exotel_telephony

import (
	"context"
	"encoding/json"
	"io"
	"sync"

	"github.com/gorilla/websocket"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"

	internal_telephony_base "github.com/rapidaai/api/assistant-api/internal/telephony/internal/base"
	internal_exotel "github.com/rapidaai/api/assistant-api/internal/telephony/internal/exotel/internal"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

type exotelWebsocketStreamer struct {
	streamer       internal_telephony_base.BaseTelephonyStreamer
	logger         commons.Logger
	streamID       string
	audioProcessor *internal_exotel.AudioProcessor

	// Output sender state
	outputSenderStarted bool
	outputSenderMu      sync.Mutex
	audioCtx            context.Context
	audioCancel         context.CancelFunc
}

func NewExotelWebsocketStreamer(logger commons.Logger, connection *websocket.Conn, assistant *internal_assistant_entity.Assistant, conversation *internal_conversation_entity.AssistantConversation, vlt *protos.VaultCredential,
) internal_type.TelephonyStreamer {
	audioProcessor, err := internal_exotel.NewAudioProcessor(logger)
	if err != nil {
		logger.Error("Failed to create audio processor", "error", err)
		return nil
	}

	exo := &exotelWebsocketStreamer{
		logger:         logger,
		streamID:       "",
		streamer:       internal_telephony_base.NewBaseTelephonyStreamer(logger, connection, assistant, conversation, vlt),
		audioProcessor: audioProcessor,
	}

	// Set up callbacks
	audioProcessor.SetInputAudioCallback(exo.sendProcessedInputAudio)
	audioProcessor.SetOutputChunkCallback(exo.sendAudioChunk)

	return exo
}

// sendProcessedInputAudio is the callback for processed input audio
func (exotel *exotelWebsocketStreamer) sendProcessedInputAudio(audio []byte) {
	// This will be called when enough audio has been buffered
	// The audio is already converted to 16kHz linear16
	exotel.streamer.LockInputAudioBuffer()
	exotel.streamer.InputBuffer().Write(audio)
	exotel.streamer.UnlockInputAudioBuffer()
}

// sendAudioChunk sends an audio chunk to Exotel
func (exotel *exotelWebsocketStreamer) sendAudioChunk(chunk *internal_exotel.AudioChunk) error {
	if exotel.streamID == "" {
		return nil
	}
	return exotel.sendingExotelMessage("media", map[string]interface{}{
		"payload": exotel.streamer.Encoder().EncodeToString(chunk.Data),
	})
}

// stopAudioProcessing stops the output sender goroutine
func (exotel *exotelWebsocketStreamer) stopAudioProcessing() {
	exotel.outputSenderMu.Lock()
	if exotel.audioCancel != nil {
		exotel.audioCancel()
		exotel.audioCancel = nil
	}
	exotel.outputSenderMu.Unlock()
}

// startOutputSender starts the consistent audio output sender
func (exotel *exotelWebsocketStreamer) startOutputSender() {
	exotel.outputSenderMu.Lock()
	defer exotel.outputSenderMu.Unlock()

	if exotel.outputSenderStarted {
		return
	}

	exotel.audioCtx, exotel.audioCancel = context.WithCancel(exotel.streamer.Context())
	exotel.outputSenderStarted = true
	go exotel.audioProcessor.RunOutputSender(exotel.audioCtx)
}

func (exotel *exotelWebsocketStreamer) Context() context.Context {
	return exotel.streamer.Context()
}

func (exotel *exotelWebsocketStreamer) Recv() (*protos.AssistantTalkInput, error) {
	if exotel.streamer.Connection() == nil {
		return nil, io.EOF
	}

	_, message, err := exotel.streamer.Connection().ReadMessage()
	if err != nil {
		exotel.stopAudioProcessing()
		exotel.streamer.Cancel()
		return nil, io.EOF
	}

	var mediaEvent internal_exotel.ExotelMediaEvent
	if err := json.Unmarshal(message, &mediaEvent); err != nil {
		exotel.logger.Error("Failed to unmarshal Exotel media event", "error", err.Error())
		return nil, nil
	}

	switch mediaEvent.Event {
	case "connected":
		// Return downstream config (16kHz linear16) for STT/TTS
		downstreamConfig := exotel.audioProcessor.GetDownstreamConfig()
		return exotel.streamer.CreateConnectionRequest(downstreamConfig, downstreamConfig)
	case "start":
		exotel.handleStartEvent(mediaEvent)
		return nil, nil
	case "media":
		return exotel.handleMediaEvent(mediaEvent)
	case "dtmf":
		return nil, nil
	case "stop":
		exotel.stopAudioProcessing()
		exotel.streamer.Cancel()
		return nil, io.EOF
	default:
		exotel.logger.Warn("Unhandled Exotel event", "event", mediaEvent.Event)
		return nil, nil
	}
}

func (exotel *exotelWebsocketStreamer) Send(response *protos.AssistantTalkOutput) error {
	switch data := response.GetData().(type) {
	case *protos.AssistantTalkOutput_Assistant:
		switch content := data.Assistant.Message.(type) {
		case *protos.ConversationAssistantMessage_Audio:
			// Process audio through the audio processor (converts 16kHz -> 8kHz linear16)
			// The audio will be sent at consistent 20ms intervals by RunOutputSender
			if err := exotel.audioProcessor.ProcessOutputAudio(content.Audio); err != nil {
				exotel.logger.Error("Failed to process output audio", "error", err.Error())
				return err
			}
		}
	case *protos.AssistantTalkOutput_Interruption:
		// interrupt on word given by stt
		if data.Interruption.Type == protos.ConversationInterruption_INTERRUPTION_TYPE_WORD {
			// Clear both input and output buffers
			exotel.audioProcessor.ClearInputBuffer()
			exotel.audioProcessor.ClearOutputBuffer()

			if err := exotel.sendingExotelMessage("clear", nil); err != nil {
				exotel.logger.Errorf("Error sending clear command:", err)
			}
		}
	case *protos.AssistantTalkOutput_Directive:
		if data.Directive.GetType() == protos.ConversationDirective_END_CONVERSATION {
			exotel.stopAudioProcessing()
			if err := exotel.streamer.Connection().Close(); err != nil {
				// terminate the conversation as end tool call is triggered
				exotel.logger.Errorf("Error disconnecting command:", err)
			}
		}
	}
	return nil
}

// start event contains streamSid to be used for subsequent media messages
func (exotel *exotelWebsocketStreamer) handleStartEvent(mediaEvent internal_exotel.ExotelMediaEvent) {
	exotel.streamID = mediaEvent.StreamSid
	// Start the consistent output sender when stream starts
	exotel.startOutputSender()
}

func (exotel *exotelWebsocketStreamer) handleMediaEvent(mediaEvent internal_exotel.ExotelMediaEvent) (*protos.AssistantTalkInput, error) {
	payloadBytes, err := exotel.streamer.Encoder().DecodeString(mediaEvent.Media.Payload)
	if err != nil {
		exotel.logger.Warn("Failed to decode media payload", "error", err.Error())
		return nil, nil
	}

	// Process input audio through audio processor (converts linear16 8kHz -> linear16 16kHz)
	if err := exotel.audioProcessor.ProcessInputAudio(payloadBytes); err != nil {
		exotel.logger.Debug("Failed to process input audio", "error", err.Error())
		return nil, nil
	}

	// Check if we have enough buffered audio to send downstream
	exotel.streamer.LockInputAudioBuffer()
	defer exotel.streamer.UnlockInputAudioBuffer()

	if exotel.streamer.InputBuffer().Len() > 0 {
		audioRequest := exotel.streamer.CreateVoiceRequest(exotel.streamer.InputBuffer().Bytes())
		exotel.streamer.InputBuffer().Reset()
		return audioRequest, nil
	}

	return nil, nil
}

func (exotel *exotelWebsocketStreamer) sendingExotelMessage(eventType string, mediaData map[string]interface{}) error {
	if exotel.streamer.Connection() == nil || exotel.streamID == "" {
		return nil
	}
	message := map[string]interface{}{
		"event":     eventType,
		"streamSid": exotel.streamID,
	}
	if mediaData != nil {
		message["media"] = mediaData
	}
	exotelMessageJSON, err := json.Marshal(message)
	if err != nil {
		return exotel.handleError("Failed to marshal Exotel message", err)
	}
	if err := exotel.streamer.Connection().WriteMessage(websocket.TextMessage, exotelMessageJSON); err != nil {
		return exotel.handleError("Failed to send message to Exotel", err)
	}
	return nil
}

func (exo *exotelWebsocketStreamer) handleError(message string, err error) error {
	exo.logger.Error(message, "error", err.Error())
	return err
}
