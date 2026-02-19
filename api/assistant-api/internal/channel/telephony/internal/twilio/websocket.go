// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_twilio_telephony

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gorilla/websocket"
	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	internal_telephony_base "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/base"
	internal_twilio "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/twilio/internal"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type twilioWebsocketStreamer struct {
	internal_telephony_base.BaseTelephonyStreamer

	streamID   string
	connection *websocket.Conn
}

func NewTwilioWebsocketStreamer(logger commons.Logger, connection *websocket.Conn, assistant *internal_assistant_entity.Assistant, conversation *internal_conversation_entity.AssistantConversation, vlt *protos.VaultCredential) internal_type.Streamer {
	return &twilioWebsocketStreamer{
		BaseTelephonyStreamer: internal_telephony_base.NewBaseTelephonyStreamer(
			logger, assistant, conversation, vlt,
			internal_telephony_base.WithSourceAudioConfig(internal_audio.NewMulaw8khzMonoAudioConfig()),
		),
		streamID:   "",
		connection: connection,
	}
}

func (tws *twilioWebsocketStreamer) Recv() (internal_type.Stream, error) {
	if tws.connection == nil {
		return nil, tws.handleError("WebSocket connection is nil", io.EOF)
	}
	_, message, err := tws.connection.ReadMessage()
	if err != nil {
		return nil, tws.handleWebSocketError(err)
	}

	var mediaEvent internal_twilio.TwilioMediaEvent
	if err := json.Unmarshal(message, &mediaEvent); err != nil {
		tws.Logger.Error("Failed to unmarshal Twilio media event", "error", err.Error())
		return nil, nil
	}
	switch mediaEvent.Event {
	case "connected":
		return tws.CreateConnectionRequest(), nil
	case "start":
		tws.handleStartEvent(mediaEvent)
		return nil, nil
	case "media":
		return tws.handleMediaEvent(mediaEvent)
	case "stop":
		tws.Logger.Info("Twilio stream stopped")
		tws.connection.Close()
		tws.connection = nil
		return nil, io.EOF
	default:
		tws.Logger.Warn("Unhandled Twilio event", "event", mediaEvent.Event)
		return nil, nil
	}
}

func (tws *twilioWebsocketStreamer) Send(response internal_type.Stream) error {
	switch data := response.(type) {
	case *protos.ConversationAssistantMessage:
		switch content := data.Message.(type) {
		case *protos.ConversationAssistantMessage_Audio:
			audioData := content.Audio

			var sendErr error
			tws.WithOutputBuffer(func(buf *bytes.Buffer) {
				buf.Write(audioData)
				for buf.Len() >= tws.OutputFrameSize() && tws.streamID != "" {
					chunk := buf.Next(tws.OutputFrameSize())
					if err := tws.sendTwilioMessage("media", map[string]interface{}{
						"payload": tws.Encoder().EncodeToString(chunk),
					}); err != nil {
						tws.Logger.Error("Failed to send audio chunk", "error", err.Error())
						sendErr = err
						return
					}
				}
				// Flush remaining audio when response is marked complete
				if data.GetCompleted() && buf.Len() > 0 {
					remainingChunk := buf.Bytes()
					if err := tws.sendTwilioMessage("media", map[string]interface{}{
						"payload": tws.Encoder().EncodeToString(remainingChunk),
					}); err != nil {
						tws.Logger.Errorf("Failed to send final audio chunk", "error", err.Error())
						sendErr = err
						return
					}
					buf.Reset()
				}
			})
			return sendErr
		}
	case *protos.ConversationInterruption:
		if data.Type == protos.ConversationInterruption_INTERRUPTION_TYPE_WORD {
			tws.ResetOutputBuffer()
			if err := tws.sendTwilioMessage("clear", nil); err != nil {
				tws.Logger.Errorf("Error sending clear command:", err)
			}
		}
	case *protos.ConversationDirective:
		if data.GetType() == protos.ConversationDirective_END_CONVERSATION {
			if tws.GetConversationUuid() != "" {
				client, err := tws.client(tws.VaultCredential())
				if err != nil {
					tws.Logger.Errorf("Error creating Twilio client:", err)
					if err := tws.Cancel(); err != nil {
						tws.Logger.Errorf("Error disconnecting command:", err)
					}
					return nil
				}
				params := &openapi.UpdateCallParams{}
				params.SetStatus("completed")
				if _, err := client.Api.UpdateCall(tws.GetConversationUuid(), params); err != nil {
					tws.Logger.Errorf("Error ending Twilio call:", err)
					if err := tws.Cancel(); err != nil {
						tws.Logger.Errorf("Error disconnecting command:", err)
					}
					return nil
				}
			}
			if err := tws.Cancel(); err != nil {
				tws.Logger.Errorf("Error disconnecting command:", err)
			}
		}
	}
	return nil
}

// start event contains streamSid to be used for subsequent media messages
func (tws *twilioWebsocketStreamer) handleStartEvent(mediaEvent internal_twilio.TwilioMediaEvent) {
	tws.streamID = mediaEvent.StreamSid
}

func (tws *twilioWebsocketStreamer) GetConversationUuid() string {
	v, err := tws.GetAssistatntConversation().GetMetadatas().GetString("telephony.uuid")
	if err != nil {
		return ""
	}
	return v
}

func (tws *twilioWebsocketStreamer) Cancel() error {
	tws.connection.Close()
	tws.connection = nil
	return nil
}

func (tws *twilioWebsocketStreamer) handleMediaEvent(mediaEvent internal_twilio.TwilioMediaEvent) (*protos.ConversationUserMessage, error) {
	payloadBytes, err := tws.Encoder().DecodeString(mediaEvent.Media.Payload)
	if err != nil {
		tws.Logger.Warn("Failed to decode media payload", "error", err.Error())
		return nil, nil
	}

	var audioRequest *protos.ConversationUserMessage
	tws.WithInputBuffer(func(buf *bytes.Buffer) {
		buf.Write(payloadBytes)
		if buf.Len() >= tws.InputBufferThreshold() {
			audioRequest = tws.CreateVoiceRequest(buf.Bytes())
			buf.Reset()
		}
	})
	return audioRequest, nil
}

func (tws *twilioWebsocketStreamer) sendTwilioMessage(
	eventType string,
	mediaData map[string]interface{}) error {
	if tws.connection == nil || tws.streamID == "" {
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

	if err := tws.connection.WriteMessage(websocket.TextMessage, twilioMessageJSON); err != nil {
		return tws.handleError("Failed to send message to Twilio", err)
	}

	return nil
}

func (tws *twilioWebsocketStreamer) handleError(message string, err error) error {
	tws.Logger.Error(message, "error", err.Error())
	return err
}

func (tws *twilioWebsocketStreamer) handleWebSocketError(err error) error {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		tws.Logger.Error("Unexpected websocket close error", "error", err.Error())
	} else {
		tws.Logger.Error("Failed to read message from WebSocket", "error", err.Error())
	}
	tws.Cancel()
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
