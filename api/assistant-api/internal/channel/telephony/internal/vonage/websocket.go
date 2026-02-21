// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_vonage_telephony

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/gorilla/websocket"
	callcontext "github.com/rapidaai/api/assistant-api/internal/callcontext"
	internal_telephony_base "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/base"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	protos "github.com/rapidaai/protos"
	"github.com/vonage/vonage-go-sdk"
)

type vonageWebsocketStreamer struct {
	internal_telephony_base.BaseTelephonyStreamer

	connection *websocket.Conn
}

// NewVonageWebsocketStreamer creates a Vonage WebSocket streamer.
// Vonage sends linear16 16kHz — same as the internal Rapida format, so no
// resampling is needed (nil source audio config defaults to linear16 16kHz).
func NewVonageWebsocketStreamer(logger commons.Logger, connection *websocket.Conn, cc *callcontext.CallContext, vaultCred *protos.VaultCredential) internal_type.Streamer {
	return &vonageWebsocketStreamer{
		BaseTelephonyStreamer: internal_telephony_base.NewBaseTelephonyStreamer(
			logger, cc, vaultCred,
		),
		connection: connection,
	}
}

func (vng *vonageWebsocketStreamer) Recv() (internal_type.Stream, error) {
	if vng.connection == nil {
		return nil, vng.handleError("WebSocket connection is nil", io.EOF)
	}
	messageType, message, err := vng.connection.ReadMessage()
	if err != nil {
		return nil, vng.handleWebSocketError(err)
	}
	switch messageType {
	case websocket.TextMessage:
		var textEvent map[string]interface{}
		if err := json.Unmarshal(message, &textEvent); err != nil {
			vng.Logger.Error("Failed to unmarshal text event", "error", err.Error())
			return nil, err
		}
		switch textEvent["event"] {
		case "websocket:connected":
			return vng.CreateConnectionRequest(), nil

		case "stop":
			return nil, io.EOF

		default:
			vng.Logger.Debugf("Unhandled event type: %s", textEvent["event"])
		}

	case websocket.BinaryMessage:
		msg, err := vng.handleMediaEvent(message)
		if msg == nil {
			return nil, err
		}
		return msg, err
	default:
		vng.Logger.Warn("Unhandled message type", "type", messageType)
		return nil, nil
	}
	return nil, nil
}

func (vng *vonageWebsocketStreamer) Send(response internal_type.Stream) error {
	if vng.connection == nil {
		return nil
	}
	switch data := response.(type) {
	case *protos.ConversationAssistantMessage:
		switch content := data.Message.(type) {
		case *protos.ConversationAssistantMessage_Audio:
			audioData := content.Audio

			var sendErr error
			vng.WithOutputBuffer(func(buf *bytes.Buffer) {
				buf.Write(audioData)
				for buf.Len() >= vng.OutputFrameSize() {
					chunk := buf.Next(vng.OutputFrameSize())
					if err := vng.connection.WriteMessage(websocket.BinaryMessage, chunk); err != nil {
						vng.Logger.Error("Failed to send audio chunk", "error", err.Error())
						sendErr = err
						return
					}
				}
				// Flush remaining audio when response is marked complete
				if data.GetCompleted() && buf.Len() > 0 {
					remainingChunk := buf.Bytes()
					if err := vng.connection.WriteMessage(websocket.BinaryMessage, remainingChunk); err != nil {
						vng.Logger.Errorf("Failed to send final audio chunk", "error", err.Error())
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
			vng.ResetOutputBuffer()
			if err := vng.connection.WriteMessage(websocket.TextMessage, []byte(`{"action":"clear"}`)); err != nil {
				vng.Logger.Errorf("Error sending clear command:", err)
			}
		}
	case *protos.ConversationDirective:
		if data.GetType() == protos.ConversationDirective_END_CONVERSATION {
			if vng.GetConversationUuid() != "" {
				cAuth, err := vng.Auth(vng.VaultCredential())
				if err != nil {
					vng.Logger.Errorf("Error creating Vonage client:", err)
					if err := vng.Cancel(); err != nil {
						vng.Logger.Errorf("Error disconnecting command:", err)
					}
					return nil
				}

				if _, _, err := vonage.NewVoiceClient(cAuth).Hangup(vng.GetConversationUuid()); err != nil {
					vng.Logger.Errorf("Error ending Vonage call:", err)
					if err := vng.Cancel(); err != nil {
						vng.Logger.Errorf("Error disconnecting command:", err)
					}
					return nil
				}
			}
			if err := vng.Cancel(); err != nil {
				vng.Logger.Errorf("Error disconnecting command:", err)
			}
		} else {
			if err := vng.Cancel(); err != nil {
				vng.Logger.Errorf("Error disconnecting command:", err)
			}
		}
	}
	return nil
}

func (vng *vonageWebsocketStreamer) handleMediaEvent(message []byte) (*protos.ConversationUserMessage, error) {
	var audioRequest *protos.ConversationUserMessage
	vng.WithInputBuffer(func(buf *bytes.Buffer) {
		buf.Write(message)
		if buf.Len() >= vng.InputBufferThreshold() {
			audioRequest = vng.CreateVoiceRequest(buf.Bytes())
			buf.Reset()
		}
	})
	return audioRequest, nil
}

func (vng *vonageWebsocketStreamer) handleError(message string, err error) error {
	vng.Logger.Error(message, "error", err.Error())
	return err
}

func (vng *vonageWebsocketStreamer) handleWebSocketError(err error) error {
	vng.Cancel()
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

func (tws *vonageWebsocketStreamer) GetConversationUuid() string {
	return tws.ChannelUUID
}

func (tws *vonageWebsocketStreamer) Cancel() error {
	tws.connection.Close()
	tws.connection = nil
	return nil
}
