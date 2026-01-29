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

	"github.com/gorilla/websocket"
	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	streamers "github.com/rapidaai/api/assistant-api/internal/streamers"
	internal_telephony_base "github.com/rapidaai/api/assistant-api/internal/telephony/internal/base"
	internal_exotel "github.com/rapidaai/api/assistant-api/internal/telephony/internal/exotel/internal"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

type exotelWebsocketStreamer struct {
	streamer internal_telephony_base.BaseTelephonyStreamer
	logger   commons.Logger
	streamID string
}

func NewExotelWebsocketStreamer(logger commons.Logger, connection *websocket.Conn, assistant *internal_assistant_entity.Assistant, conversation *internal_conversation_entity.AssistantConversation, vlt *protos.VaultCredential,
) streamers.Streamer {
	return &exotelWebsocketStreamer{
		logger:   logger,
		streamID: "",
		streamer: internal_telephony_base.NewBaseTelephonyStreamer(logger, connection, assistant, conversation, vlt),
	}
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
		return exotel.streamer.CreateConnectionRequest(internal_audio.NewLinear8khzMonoAudioConfig(), internal_audio.NewLinear8khzMonoAudioConfig())
	case "start":
		exotel.handleStartEvent(mediaEvent)
		return nil, nil
	case "media":
		return exotel.handleMediaEvent(mediaEvent)
	case "dtmf":
		return nil, nil
	case "stop":
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
			//1ms 32  10ms 320byte @ 16000Hz, 16-bit mono PCM = 640 bytes
			// Each message needs to be a 20ms sample of audio.
			// At 8kHz the message should be 320 bytes.
			// At 16kHz the message should be 640 bytes.
			bufferSizeThreshold := 32 * 20
			audioData := content.Audio

			// Use vng.audioBuffer to handle pending data across calls
			exotel.streamer.LockOutputAudioBuffer()
			defer exotel.streamer.UnlockOutputAudioBuffer()

			// Append incoming audio data to the buffer
			exotel.streamer.OutputBuffer().Write(audioData)
			// Process and send chunks of 640 bytes
			for exotel.streamer.OutputBuffer().Len() >= bufferSizeThreshold && exotel.streamID != "" {
				chunk := exotel.streamer.OutputBuffer().Next(bufferSizeThreshold) // Get and remove the next 640 bytes
				if err := exotel.sendingExotelMessage("media", map[string]interface{}{
					"payload": exotel.streamer.Encoder().EncodeToString(chunk),
				}); err != nil {
					exotel.logger.Error("Failed to send audio chunk", "error", err.Error())
					return err
				}
			}

			// If response is marked as completed, flush any remaining audio in the buffer
			if data.Assistant.GetCompleted() && exotel.streamer.OutputBuffer().Len() > 0 {
				remainingChunk := exotel.streamer.OutputBuffer().Bytes()
				if err := exotel.sendingExotelMessage("media", map[string]interface{}{
					"payload": exotel.streamer.Encoder().EncodeToString(remainingChunk),
				}); err != nil {
					exotel.logger.Errorf("Failed to send final audio chunk", "error", err.Error())
					return err
				}
				exotel.streamer.OutputBuffer().Reset() // Clear the buffer after flushing
			}
		}
	case *protos.AssistantTalkOutput_Interruption:
		// interrupt on word given by stt
		if data.Interruption.Type == protos.ConversationInterruption_INTERRUPTION_TYPE_WORD {
			exotel.streamer.LockInputAudioBuffer()
			exotel.streamer.OutputBuffer().Reset()
			exotel.streamer.UnlockInputAudioBuffer()

			//
			if err := exotel.sendingExotelMessage("clear", nil); err != nil {
				exotel.logger.Errorf("Error sending clear command:", err)
			}
		}
	case *protos.AssistantTalkOutput_Directive:
		if data.Directive.GetType() == protos.ConversationDirective_END_CONVERSATION {
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
}

// when exotel is connected then connect the assistant

func (exotel *exotelWebsocketStreamer) handleMediaEvent(mediaEvent internal_exotel.ExotelMediaEvent) (*protos.AssistantTalkInput, error) {
	payloadBytes, err := exotel.streamer.Encoder().DecodeString(mediaEvent.Media.Payload)
	if err != nil {
		exotel.logger.Warn("Failed to decode media payload", "error", err.Error())
		return nil, nil
	}

	exotel.streamer.LockInputAudioBuffer()
	defer exotel.streamer.UnlockInputAudioBuffer()

	// 1ms 8 bytes @ 8kHz µ-law mono 60ms of audio as silero can't process smaller chunk for mulaw
	exotel.streamer.InputBuffer().Write(payloadBytes)
	const bufferSizeThreshold = 32 * 60

	if exotel.streamer.InputBuffer().Len() >= bufferSizeThreshold {
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
