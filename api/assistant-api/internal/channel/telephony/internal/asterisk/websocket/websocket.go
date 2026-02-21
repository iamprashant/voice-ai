// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_asterisk_websocket

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	callcontext "github.com/rapidaai/api/assistant-api/internal/callcontext"
	internal_asterisk "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/asterisk/internal"
	internal_telephony_base "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/base"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

// asteriskWebsocketStreamer handles WebSocket communication with Asterisk chan_websocket.
type asteriskWebsocketStreamer struct {
	internal_telephony_base.BaseTelephonyStreamer

	audioProcessor *AudioProcessor
	connection     *websocket.Conn
	channelName    string

	// Output sender state
	outputSenderStarted bool
	outputSenderMu      sync.Mutex
	audioCtx            context.Context
	audioCancel         context.CancelFunc

	// Media buffering state
	mediaBuffering bool
	mediaBufferMu  sync.Mutex
}

// NewAsteriskWebsocketStreamer creates a new Asterisk WebSocket streamer.
func NewAsteriskWebsocketStreamer(
	logger commons.Logger,
	connection *websocket.Conn,
	cc *callcontext.CallContext,
	vaultCred *protos.VaultCredential,
) internal_type.Streamer {
	audioProcessor, err := NewAudioProcessor(logger)
	if err != nil {
		logger.Error("Failed to create audio processor", "error", err)
		return nil
	}

	aws := &asteriskWebsocketStreamer{
		BaseTelephonyStreamer: internal_telephony_base.NewBaseTelephonyStreamer(
			logger, cc, vaultCred,
		),
		audioProcessor: audioProcessor,
		connection:     connection,
	}

	// Set up callbacks
	audioProcessor.SetInputAudioCallback(aws.sendProcessedInputAudio)
	audioProcessor.SetOutputChunkCallback(aws.sendAudioChunk)
	return aws
}

// sendProcessedInputAudio is the callback for processed input audio.
func (aws *asteriskWebsocketStreamer) sendProcessedInputAudio(audio []byte) {
	aws.WithInputBuffer(func(buf *bytes.Buffer) {
		buf.Write(audio)
	})
}

// sendAudioChunk sends an audio chunk to Asterisk
func (aws *asteriskWebsocketStreamer) sendAudioChunk(chunk *AudioChunk) error {
	if aws.connection == nil {
		return nil
	}

	// Send binary audio data directly to Asterisk
	return aws.connection.WriteMessage(websocket.BinaryMessage, chunk.Data)
}

// Context returns the streamer context.
func (aws *asteriskWebsocketStreamer) Context() context.Context {
	return aws.BaseTelephonyStreamer.Context()
}

// Recv receives and processes messages from Asterisk WebSocket
func (aws *asteriskWebsocketStreamer) Recv() (internal_type.Stream, error) {
	if aws.connection == nil {
		return nil, aws.handleError("WebSocket connection is nil", io.EOF)
	}
	messageType, message, err := aws.connection.ReadMessage()
	if err != nil {
		return nil, aws.handleWebSocketError(err)
	}

	switch messageType {
	case websocket.BinaryMessage:
		msg, err := aws.handleAudioData(message)
		if msg == nil {
			return nil, err
		}
		return msg, err
	case websocket.TextMessage:
		event, err := internal_asterisk.ParseAsteriskEvent(string(message))
		if err != nil {
			aws.Logger.Warn("Failed to parse Asterisk event", "error", err.Error(), "message", message)
			return nil, nil
		}
		switch event.Event {
		case "MEDIA_START":
			aws.channelName = event.Channel
			aws.Logger.Info("Asterisk media started", "channel", aws.channelName, "optimal_frame_size", event.OptimalFrameSize)
			if event.OptimalFrameSize > 0 {
				aws.audioProcessor.SetOptimalFrameSize(event.OptimalFrameSize)
			}
			aws.startOutputSender()
			return aws.CreateConnectionRequest(), nil
		case "MEDIA_STOP":
			aws.Logger.Info("Asterisk media stopped")
			aws.stopAudioProcessing()
			aws.Cancel()
			return nil, io.EOF

		case "MEDIA_XON":
			// Resume audio output (flow control)
			aws.audioProcessor.SetXON()

		case "MEDIA_XOFF":
			// Pause audio output (flow control)
			aws.audioProcessor.SetXOFF()
		case "MEDIA_BUFFERING_COMPLETED":
			aws.setMediaBuffering(false)
		default:
			// Handle JSON command responses
			if event.Command != "" {
				aws.Logger.Debug("Received Asterisk command response", "command", event.Command)
			} else if event.RawMessage != "" {
				aws.Logger.Debug("Received unhandled Asterisk message", "message", event.RawMessage)
			}
		}
	case websocket.CloseMessage:
		return nil, io.EOF
	default:
		aws.Logger.Warn("Received unsupported WebSocket message type", "type", messageType)
	}

	return nil, nil
}

// handleAudioData processes incoming binary audio data from Asterisk.
func (aws *asteriskWebsocketStreamer) handleAudioData(audio []byte) (*protos.ConversationUserMessage, error) {
	// Process input audio through audio processor (converts ulaw 8kHz -> linear16 16kHz)
	if err := aws.audioProcessor.ProcessInputAudio(audio); err != nil {
		aws.Logger.Debug("Failed to process input audio", "error", err.Error())
		return nil, nil
	}

	// Check if we have enough buffered audio to send downstream
	var audioRequest *protos.ConversationUserMessage
	aws.WithInputBuffer(func(buf *bytes.Buffer) {
		if buf.Len() > 0 {
			audioRequest = aws.CreateVoiceRequest(buf.Bytes())
			buf.Reset()
		}
	})

	return audioRequest, nil
}

// Send sends output to Asterisk
func (aws *asteriskWebsocketStreamer) Send(response internal_type.Stream) error {
	switch data := response.(type) {
	case *protos.ConversationAssistantMessage:
		switch content := data.Message.(type) {
		case *protos.ConversationAssistantMessage_Audio:
			// Process audio through the audio processor (converts 16kHz -> 8kHz ulaw)
			// The audio will be sent at consistent 20ms intervals by RunOutputSender
			if err := aws.audioProcessor.ProcessOutputAudio(content.Audio); err != nil {
				aws.Logger.Error("Failed to process output audio", "error", err.Error())
				return err
			}
		}

	case *protos.ConversationInterruption:
		if data.Type == protos.ConversationInterruption_INTERRUPTION_TYPE_WORD {
			// Clear both input and output buffers
			aws.audioProcessor.ClearInputBuffer()
			aws.audioProcessor.ClearOutputBuffer()

			// No direct "clear" command in Asterisk media WebSocket,
			// but we can stop buffering if active
			if aws.isMediaBuffering() {
				aws.sendCommand("STOP_MEDIA_BUFFERING")
			}
		}

	case *protos.ConversationDirective:
		if data.GetType() == protos.ConversationDirective_END_CONVERSATION {
			aws.stopAudioProcessing()
			if err := aws.sendCommand("HANGUP"); err != nil {
				aws.Logger.Warn("Failed to send HANGUP via WebSocket, trying ARI API", "error", err)
				if aws.channelName != "" {
					if err := aws.hangupViaARI(); err != nil {
						aws.Logger.Error("Failed to hangup via ARI API", "error", err)
					}
				}
			}
			if err := aws.Cancel(); err != nil {
				aws.Logger.Errorf("Error disconnecting:", err)
			}
		}
	}

	return nil
}

// stopAudioProcessing stops the output sender goroutine
func (aws *asteriskWebsocketStreamer) stopAudioProcessing() {
	aws.outputSenderMu.Lock()
	if aws.audioCancel != nil {
		aws.audioCancel()
		aws.audioCancel = nil
	}
	aws.outputSenderMu.Unlock()
}

// startOutputSender starts the consistent audio output sender
func (aws *asteriskWebsocketStreamer) startOutputSender() {
	aws.outputSenderMu.Lock()
	defer aws.outputSenderMu.Unlock()

	if aws.outputSenderStarted {
		return
	}

	aws.audioCtx, aws.audioCancel = context.WithCancel(aws.BaseTelephonyStreamer.Context())
	aws.outputSenderStarted = true
	go aws.audioProcessor.RunOutputSender(aws.audioCtx)
}

// sendCommand sends a text command to Asterisk
func (aws *asteriskWebsocketStreamer) sendCommand(command string) error {
	if aws.connection == nil {
		return nil
	}
	return aws.connection.WriteMessage(websocket.TextMessage, []byte(command))
}

// setMediaBuffering sets the media buffering state
func (aws *asteriskWebsocketStreamer) setMediaBuffering(buffering bool) {
	aws.mediaBufferMu.Lock()
	aws.mediaBuffering = buffering
	aws.mediaBufferMu.Unlock()
}

// isMediaBuffering returns whether media buffering is active
func (aws *asteriskWebsocketStreamer) isMediaBuffering() bool {
	aws.mediaBufferMu.Lock()
	defer aws.mediaBufferMu.Unlock()
	return aws.mediaBuffering
}

// handleError logs an error and returns it.
func (aws *asteriskWebsocketStreamer) handleError(message string, err error) error {
	aws.Logger.Error(message, "error", err.Error())
	return err
}

// handleWebSocketError handles WebSocket errors
func (aws *asteriskWebsocketStreamer) handleWebSocketError(err error) error {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		aws.Logger.Error("Unexpected websocket close error", "error", err.Error())
	} else {
		aws.Logger.Error("Failed to read message from WebSocket", "error", err.Error())
	}
	aws.Cancel()
	return io.EOF
}

// hangupViaARI hangs up the call using the Asterisk ARI API
// This is a fallback mechanism when the WebSocket HANGUP command fails
func (aws *asteriskWebsocketStreamer) hangupViaARI() error {
	vaultCredential := aws.VaultCredential()
	if vaultCredential == nil {
		return fmt.Errorf("vault credential is nil")
	}

	credMap := vaultCredential.GetValue().AsMap()

	ariURL, _ := credMap["ari_url"].(string)
	ariURL = fmt.Sprintf("%s/ari/channels/%s", ariURL, aws.channelName)
	user, _ := credMap["ari_user"].(string)
	password, _ := credMap["ari_password"].(string)

	req, err := http.NewRequest("DELETE", ariURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(user, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("ARI API returned status: %d", resp.StatusCode)
	}

	aws.Logger.Info("Successfully hung up call via ARI API", "channel", aws.channelName)
	return nil
}

func (tws *asteriskWebsocketStreamer) Cancel() error {
	if tws.connection != nil {
		tws.connection.Close()
		tws.connection = nil
	}
	return nil
}
