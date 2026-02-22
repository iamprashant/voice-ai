// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_exotel_telephony

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/gorilla/websocket"
	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	callcontext "github.com/rapidaai/api/assistant-api/internal/callcontext"
	internal_telephony_base "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/base"
	internal_exotel "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/exotel/internal"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

type exotelWebsocketStreamer struct {
	internal_telephony_base.BaseTelephonyStreamer

	connection *websocket.Conn
	streamID   string
}

func NewExotelWebsocketStreamer(logger commons.Logger, connection *websocket.Conn, cc *callcontext.CallContext, vaultCred *protos.VaultCredential,
) internal_type.Streamer {
	return &exotelWebsocketStreamer{
		BaseTelephonyStreamer: internal_telephony_base.NewBaseTelephonyStreamer(
			logger, cc, vaultCred,
			internal_telephony_base.WithSourceAudioConfig(internal_audio.NewLinear8khzMonoAudioConfig()),
		),
		streamID:   "",
		connection: connection,
	}
}

func (exotel *exotelWebsocketStreamer) Recv() (internal_type.Stream, error) {
	if exotel.connection == nil {
		return nil, io.EOF
	}

	_, message, err := exotel.connection.ReadMessage()
	if err != nil {
		exotel.connection.Close()
		exotel.connection = nil
		return nil, io.EOF
	}

	var mediaEvent internal_exotel.ExotelMediaEvent
	if err := json.Unmarshal(message, &mediaEvent); err != nil {
		exotel.Logger.Error("Failed to unmarshal Exotel media event", "error", err.Error())
		return nil, nil
	}

	switch mediaEvent.Event {
	case "connected":
		return exotel.CreateConnectionRequest(), nil
	case "start":
		exotel.handleStartEvent(mediaEvent)
		return nil, nil
	case "media":
		msg, err := exotel.handleMediaEvent(mediaEvent)
		if msg == nil {
			return nil, err
		}
		return msg, err
	case "dtmf":
		return nil, nil
	case "stop":
		exotel.Cancel()
		return nil, io.EOF
	default:
		exotel.Logger.Warn("Unhandled Exotel event", "event", mediaEvent.Event)
		return nil, nil
	}
}

func (exotel *exotelWebsocketStreamer) Send(response internal_type.Stream) error {
	switch data := response.(type) {
	case *protos.ConversationAssistantMessage:
		switch content := data.Message.(type) {
		case *protos.ConversationAssistantMessage_Audio:
			var sendErr error
			exotel.WithOutputBuffer(func(buf *bytes.Buffer) {
				buf.Write(content.Audio)
				for buf.Len() >= exotel.OutputFrameSize() && exotel.streamID != "" {
					chunk := buf.Next(exotel.OutputFrameSize())
					if err := exotel.sendingExotelMessage("media", map[string]interface{}{
						"payload": exotel.Encoder().EncodeToString(chunk),
					}); err != nil {
						exotel.Logger.Error("Failed to send audio chunk", "error", err.Error())
						sendErr = err
						return
					}
				}
				// Flush remaining audio when response is marked complete
				if data.GetCompleted() && buf.Len() > 0 {
					remainingChunk := buf.Bytes()
					if err := exotel.sendingExotelMessage("media", map[string]interface{}{
						"payload": exotel.Encoder().EncodeToString(remainingChunk),
					}); err != nil {
						exotel.Logger.Errorf("Failed to send final audio chunk", "error", err.Error())
						sendErr = err
						return
					}
					buf.Reset()
				}
			})
			return sendErr
		}
	case *protos.ConversationInterruption:
		// interrupt on word given by stt
		if data.Type == protos.ConversationInterruption_INTERRUPTION_TYPE_WORD {
			exotel.ResetOutputBuffer()
			if err := exotel.sendingExotelMessage("clear", nil); err != nil {
				exotel.Logger.Errorf("Error sending clear command:", err)
			}
		}
	case *protos.ConversationDirective:
		if data.GetType() == protos.ConversationDirective_END_CONVERSATION {
			if err := exotel.Cancel(); err != nil {
				exotel.Logger.Errorf("Error disconnecting command:", err)
			}
		}
	}
	return nil
}

// start event contains streamSid to be used for subsequent media messages
func (exotel *exotelWebsocketStreamer) handleStartEvent(mediaEvent internal_exotel.ExotelMediaEvent) {
	exotel.streamID = mediaEvent.StreamSid
}

func (exotel *exotelWebsocketStreamer) handleMediaEvent(mediaEvent internal_exotel.ExotelMediaEvent) (*protos.ConversationUserMessage, error) {
	payloadBytes, err := exotel.Encoder().DecodeString(mediaEvent.Media.Payload)
	if err != nil {
		exotel.Logger.Warn("Failed to decode media payload", "error", err.Error())
		return nil, nil
	}

	var audioRequest *protos.ConversationUserMessage
	exotel.WithInputBuffer(func(buf *bytes.Buffer) {
		buf.Write(payloadBytes)
		if buf.Len() >= exotel.InputBufferThreshold() {
			audioRequest = exotel.CreateVoiceRequest(buf.Bytes())
			buf.Reset()
		}
	})
	return audioRequest, nil
}

func (exotel *exotelWebsocketStreamer) sendingExotelMessage(eventType string, mediaData map[string]interface{}) error {
	if exotel.connection == nil || exotel.streamID == "" {
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
	if err := exotel.connection.WriteMessage(websocket.TextMessage, exotelMessageJSON); err != nil {
		return exotel.handleError("Failed to send message to Exotel", err)
	}
	return nil
}

func (exo *exotelWebsocketStreamer) handleError(message string, err error) error {
	exo.Logger.Error(message, "error", err.Error())
	return err
}

func (tws *exotelWebsocketStreamer) Cancel() error {
	tws.connection.Close()
	tws.connection = nil
	return nil
}
