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

	"github.com/gorilla/websocket"
	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_streamers "github.com/rapidaai/api/assistant-api/internal/streamers"
	internal_telephony_base "github.com/rapidaai/api/assistant-api/internal/telephony/internal/base"
	internal_twilio "github.com/rapidaai/api/assistant-api/internal/telephony/internal/twilio/internal"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type twilioWebsocketStreamer struct {
	streamID string
	streamer internal_telephony_base.BaseTelephonyStreamer
	logger   commons.Logger
}

func NewTwilioWebsocketStreamer(logger commons.Logger, connection *websocket.Conn, assistant *internal_assistant_entity.Assistant, conversation *internal_conversation_entity.AssistantConversation, vlt *protos.VaultCredential) internal_streamers.Streamer {
	return &twilioWebsocketStreamer{
		logger:   logger,
		streamID: "",
		streamer: internal_telephony_base.NewBaseTelephonyStreamer(logger, connection, assistant, conversation, vlt),
	}
}

func (tws *twilioWebsocketStreamer) Context() context.Context {
	return tws.streamer.Context()
}

func (tws *twilioWebsocketStreamer) Recv() (*protos.AssistantMessagingRequest, error) {
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
		return tws.streamer.CreateConnectionRequest(internal_audio.NewMulaw8khzMonoAudioConfig(), internal_audio.NewMulaw8khzMonoAudioConfig())
	case "start":
		tws.handleStartEvent(mediaEvent)
		return nil, nil
	case "media":
		return tws.handleMediaEvent(mediaEvent)
	case "stop":
		tws.logger.Info("Twilio stream stopped")
		tws.streamer.Cancel()
		return nil, io.EOF
	default:
		tws.logger.Warn("Unhandled Twilio event", "event", mediaEvent.Event)
		return nil, nil
	}
}

func (tws *twilioWebsocketStreamer) Send(response *protos.AssistantMessagingResponse) error {
	switch data := response.GetData().(type) {
	case *protos.AssistantMessagingResponse_Assistant:
		switch content := data.Assistant.Message.(type) {
		case *protos.AssistantConversationAssistantMessage_Audio:
			//1ms 32  10ms 320byte @ 16000Hz, 16-bit mono PCM = 640 bytes
			// Each message needs to be a 20ms sample of audio.
			// At 8kHz the message should be 320 bytes.
			// At 16kHz the message should be 640 bytes.
			bufferSizeThreshold := 8 * 20
			audioData := content.Audio.GetContent()

			// Use vng.audioBuffer to handle pending data across calls
			tws.streamer.LockOutputAudioBuffer()
			defer tws.streamer.UnlockOutputAudioBuffer()

			// Append incoming audio data to the buffer
			tws.streamer.OutputBuffer().Write(audioData)
			// Process and send chunks of 640 bytes
			for tws.streamer.OutputBuffer().Len() >= bufferSizeThreshold && tws.streamID != "" {
				chunk := tws.streamer.OutputBuffer().Next(bufferSizeThreshold) // Get and remove the next 640 bytes
				if err := tws.sendTwilioMessage("media", map[string]interface{}{
					"payload": tws.streamer.Encoder().EncodeToString(chunk),
				}); err != nil {
					tws.logger.Error("Failed to send audio chunk", "error", err.Error())
					return err
				}
			}

			// If response is marked as completed, flush any remaining audio in the buffer
			if data.Assistant.GetCompleted() && tws.streamer.OutputBuffer().Len() > 0 {
				remainingChunk := tws.streamer.OutputBuffer().Bytes()
				if err := tws.sendTwilioMessage("media", map[string]interface{}{
					"payload": tws.streamer.Encoder().EncodeToString(remainingChunk),
				}); err != nil {
					tws.logger.Errorf("Failed to send final audio chunk", "error", err.Error())
					return err
				}
				tws.streamer.OutputBuffer().Reset() // Clear the buffer after flushing
			}
		}
	case *protos.AssistantMessagingResponse_Interruption:
		if data.Interruption.Type == protos.AssistantConversationInterruption_INTERRUPTION_TYPE_WORD {
			tws.streamer.LockOutputAudioBuffer()
			tws.streamer.OutputBuffer().Reset() // Clear the buffer after flushing
			tws.streamer.UnlockOutputAudioBuffer()

			if err := tws.sendTwilioMessage("clear", nil); err != nil {
				tws.logger.Errorf("Error sending clear command:", err)
			}
		}
	case *protos.AssistantMessagingResponse_Action:
		if data.Action.GetAction() == protos.AssistantConversationAction_END_CONVERSATION {
			if tws.streamer.GetUuid() != "" {
				//
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
}

func (tws *twilioWebsocketStreamer) handleMediaEvent(mediaEvent internal_twilio.TwilioMediaEvent) (*protos.AssistantMessagingRequest, error) {
	payloadBytes, err := tws.streamer.Encoder().DecodeString(mediaEvent.Media.Payload)
	if err != nil {
		tws.logger.Warn("Failed to decode media payload", "error", err.Error())
		return nil, nil
	}

	tws.streamer.LockInputAudioBuffer()
	defer tws.streamer.UnlockInputAudioBuffer()

	// 1ms 8 bytes @ 8kHz µ-law mono 60ms of audio as silero can't process smaller chunk for mulaw
	tws.streamer.InputBuffer().Write(payloadBytes)
	const bufferSizeThreshold = 8 * 60

	if tws.streamer.InputBuffer().Len() >= bufferSizeThreshold {
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
