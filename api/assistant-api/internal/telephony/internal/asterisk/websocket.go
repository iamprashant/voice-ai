// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_asterisk_telephony

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_asterisk "github.com/rapidaai/api/assistant-api/internal/telephony/internal/asterisk/internal"
	internal_telephony_base "github.com/rapidaai/api/assistant-api/internal/telephony/internal/base"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

// asteriskWebsocketStreamer handles WebSocket communication with Asterisk chan_websocket
type asteriskWebsocketStreamer struct {
	channelName    string
	streamer       internal_telephony_base.BaseTelephonyStreamer
	logger         commons.Logger
	audioProcessor *internal_asterisk.AudioProcessor

	// Output sender state
	outputSenderStarted bool
	outputSenderMu      sync.Mutex
	audioCtx            context.Context
	audioCancel         context.CancelFunc

	// Media buffering state
	mediaBuffering bool
	mediaBufferMu  sync.Mutex
}

// NewAsteriskWebsocketStreamer creates a new Asterisk WebSocket streamer
func NewAsteriskWebsocketStreamer(
	logger commons.Logger,
	connection *websocket.Conn,
	assistant *internal_assistant_entity.Assistant,
	conversation *internal_conversation_entity.AssistantConversation,
	vlt *protos.VaultCredential,
) internal_type.TelephonyStreamer {
	audioProcessor, err := internal_asterisk.NewAudioProcessor(logger)
	if err != nil {
		logger.Error("Failed to create audio processor", "error", err)
		return nil
	}

	aws := &asteriskWebsocketStreamer{
		logger:         logger,
		channelName:    "",
		streamer:       internal_telephony_base.NewBaseTelephonyStreamer(logger, connection, assistant, conversation, vlt),
		audioProcessor: audioProcessor,
	}

	// Set up callbacks
	audioProcessor.SetInputAudioCallback(aws.sendProcessedInputAudio)
	audioProcessor.SetOutputChunkCallback(aws.sendAudioChunk)

	return aws
}

// sendProcessedInputAudio is the callback for processed input audio
func (aws *asteriskWebsocketStreamer) sendProcessedInputAudio(audio []byte) {
	aws.streamer.LockInputAudioBuffer()
	aws.streamer.InputBuffer().Write(audio)
	aws.streamer.UnlockInputAudioBuffer()
}

// sendAudioChunk sends an audio chunk to Asterisk
func (aws *asteriskWebsocketStreamer) sendAudioChunk(chunk *internal_asterisk.AudioChunk) error {
	if aws.streamer.Connection() == nil {
		return nil
	}

	// Send binary audio data directly to Asterisk
	return aws.streamer.Connection().WriteMessage(websocket.BinaryMessage, chunk.Data)
}

// Context returns the streamer context
func (aws *asteriskWebsocketStreamer) Context() context.Context {
	return aws.streamer.Context()
}

// Recv receives and processes messages from Asterisk WebSocket
func (aws *asteriskWebsocketStreamer) Recv() (*protos.AssistantTalkInput, error) {
	if aws.streamer.Connection() == nil {
		return nil, aws.handleError("WebSocket connection is nil", io.EOF)
	}

	messageType, message, err := aws.streamer.Connection().ReadMessage()
	if err != nil {
		return nil, aws.handleWebSocketError(err)
	}

	// Handle binary audio data
	if messageType == websocket.BinaryMessage {
		return aws.handleAudioData(message)
	}

	// Handle text messages (events/commands)
	if messageType == websocket.TextMessage {
		return aws.handleTextMessage(string(message))
	}

	return nil, nil
}

// handleAudioData processes incoming binary audio data from Asterisk
func (aws *asteriskWebsocketStreamer) handleAudioData(audio []byte) (*protos.AssistantTalkInput, error) {
	// Process input audio through audio processor (converts ulaw 8kHz -> linear16 16kHz)
	if err := aws.audioProcessor.ProcessInputAudio(audio); err != nil {
		aws.logger.Debug("Failed to process input audio", "error", err.Error())
		return nil, nil
	}

	// Check if we have enough buffered audio to send downstream
	aws.streamer.LockInputAudioBuffer()
	defer aws.streamer.UnlockInputAudioBuffer()

	if aws.streamer.InputBuffer().Len() > 0 {
		audioRequest := aws.streamer.CreateVoiceRequest(aws.streamer.InputBuffer().Bytes())
		aws.streamer.InputBuffer().Reset()
		return audioRequest, nil
	}

	return nil, nil
}

// handleTextMessage processes text messages (events/commands) from Asterisk
func (aws *asteriskWebsocketStreamer) handleTextMessage(message string) (*protos.AssistantTalkInput, error) {
	event, err := internal_asterisk.ParseAsteriskEvent(message)
	if err != nil {
		aws.logger.Warn("Failed to parse Asterisk event", "error", err.Error(), "message", message)
		return nil, nil
	}

	switch event.Event {
	case "MEDIA_START":
		return aws.handleMediaStart(event)

	case "MEDIA_STOP":
		aws.logger.Info("Asterisk media stopped")
		aws.stopAudioProcessing()
		aws.streamer.Cancel()
		return nil, io.EOF

	case "MEDIA_XON":
		// Resume audio output (flow control)
		aws.audioProcessor.SetXON()
		return nil, nil

	case "MEDIA_XOFF":
		// Pause audio output (flow control)
		aws.audioProcessor.SetXOFF()
		return nil, nil

	case "MEDIA_BUFFERING_COMPLETED":
		aws.setMediaBuffering(false)
		return nil, nil

	default:
		// Handle JSON command responses
		if event.Command != "" {
			aws.logger.Debug("Received Asterisk command response", "command", event.Command)
		} else if event.RawMessage != "" {
			aws.logger.Debug("Received unhandled Asterisk message", "message", event.RawMessage)
		}
		return nil, nil
	}
}

// handleMediaStart handles the MEDIA_START event from Asterisk
func (aws *asteriskWebsocketStreamer) handleMediaStart(event *internal_asterisk.AsteriskMediaEvent) (*protos.AssistantTalkInput, error) {
	aws.channelName = event.Channel
	aws.logger.Info("Asterisk media started", "channel", aws.channelName, "optimal_frame_size", event.OptimalFrameSize)

	// Set optimal frame size if provided
	if event.OptimalFrameSize > 0 {
		aws.audioProcessor.SetOptimalFrameSize(event.OptimalFrameSize)
	}

	// Start the consistent output sender when stream starts
	aws.startOutputSender()

	// Return downstream config (16kHz linear16) for STT/TTS
	downstreamConfig := aws.audioProcessor.GetDownstreamConfig()
	return aws.streamer.CreateConnectionRequest(downstreamConfig, downstreamConfig)
}

// Send sends output to Asterisk
func (aws *asteriskWebsocketStreamer) Send(response *protos.AssistantTalkOutput) error {
	switch data := response.GetData().(type) {
	case *protos.AssistantTalkOutput_Assistant:
		switch content := data.Assistant.Message.(type) {
		case *protos.ConversationAssistantMessage_Audio:
			// Process audio through the audio processor (converts 16kHz -> 8kHz ulaw)
			// The audio will be sent at consistent 20ms intervals by RunOutputSender
			if err := aws.audioProcessor.ProcessOutputAudio(content.Audio); err != nil {
				aws.logger.Error("Failed to process output audio", "error", err.Error())
				return err
			}
		}

	case *protos.AssistantTalkOutput_Interruption:
		if data.Interruption.Type == protos.ConversationInterruption_INTERRUPTION_TYPE_WORD {
			// Clear both input and output buffers
			aws.audioProcessor.ClearInputBuffer()
			aws.audioProcessor.ClearOutputBuffer()

			// No direct "clear" command in Asterisk media WebSocket,
			// but we can stop buffering if active
			if aws.isMediaBuffering() {
				aws.sendCommand("STOP_MEDIA_BUFFERING")
			}
		}

	case *protos.AssistantTalkOutput_Directive:
		if data.Directive.GetType() == protos.ConversationDirective_END_CONVERSATION {
			aws.stopAudioProcessing()

			// Try to hangup via WebSocket command first
			if err := aws.sendCommand("HANGUP"); err != nil {
				aws.logger.Warn("Failed to send HANGUP via WebSocket, trying ARI API", "error", err)
				// Fallback to ARI API if channel name is available
				if aws.channelName != "" {
					if err := aws.hangupViaARI(); err != nil {
						aws.logger.Error("Failed to hangup via ARI API", "error", err)
					}
				}
			}

			if err := aws.streamer.Cancel(); err != nil {
				aws.logger.Errorf("Error disconnecting:", err)
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

	aws.audioCtx, aws.audioCancel = context.WithCancel(aws.streamer.Context())
	aws.outputSenderStarted = true
	go aws.audioProcessor.RunOutputSender(aws.audioCtx)
}

// sendCommand sends a text command to Asterisk
func (aws *asteriskWebsocketStreamer) sendCommand(command string) error {
	if aws.streamer.Connection() == nil {
		return nil
	}
	return aws.streamer.Connection().WriteMessage(websocket.TextMessage, []byte(command))
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

// handleError logs an error and returns it
func (aws *asteriskWebsocketStreamer) handleError(message string, err error) error {
	aws.logger.Error(message, "error", err.Error())
	return err
}

// handleWebSocketError handles WebSocket errors
func (aws *asteriskWebsocketStreamer) handleWebSocketError(err error) error {
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		aws.logger.Error("Unexpected websocket close error", "error", err.Error())
	} else {
		aws.logger.Error("Failed to read message from WebSocket", "error", err.Error())
	}
	aws.streamer.Cancel()
	return io.EOF
}

// hangupViaARI hangs up the call using the Asterisk ARI API
// This is a fallback mechanism when the WebSocket HANGUP command fails
func (aws *asteriskWebsocketStreamer) hangupViaARI() error {
	vaultCredential := aws.streamer.VaultCredential()
	if vaultCredential == nil {
		return fmt.Errorf("vault credential is nil")
	}

	ariConfig, err := aws.getARIConfig(vaultCredential)
	if err != nil {
		return fmt.Errorf("failed to get ARI config: %w", err)
	}

	// Build ARI DELETE channels endpoint URL
	ariURL := fmt.Sprintf("%s://%s:%d/ari/channels/%s",
		ariConfig.ARIScheme, ariConfig.ARIHost, ariConfig.ARIPort, aws.channelName)

	req, err := http.NewRequest("DELETE", ariURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.SetBasicAuth(ariConfig.ARIUser, ariConfig.ARIPassword)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("ARI API returned status: %d", resp.StatusCode)
	}

	aws.logger.Info("Successfully hung up call via ARI API", "channel", aws.channelName)
	return nil
}

// getARIConfig extracts ARI configuration from vault credential
func (aws *asteriskWebsocketStreamer) getARIConfig(vaultCredential *protos.VaultCredential) (*ARIConfig, error) {
	if vaultCredential == nil {
		return nil, fmt.Errorf("vault credential is nil")
	}

	credMap := vaultCredential.GetValue().AsMap()

	config := &ARIConfig{
		ARIHost:   "localhost",
		ARIPort:   8088,
		ARIScheme: "http",
	}

	if host, ok := credMap["ari_host"]; ok && host != nil {
		config.ARIHost = fmt.Sprintf("%v", host)
	}
	if port, ok := credMap["ari_port"]; ok && port != nil {
		switch v := port.(type) {
		case float64:
			config.ARIPort = int(v)
		case int:
			config.ARIPort = v
		}
	}
	if scheme, ok := credMap["ari_scheme"]; ok && scheme != nil {
		config.ARIScheme = fmt.Sprintf("%v", scheme)
	}
	if user, ok := credMap["ari_user"]; ok && user != nil {
		config.ARIUser = fmt.Sprintf("%v", user)
	}
	if password, ok := credMap["ari_password"]; ok && password != nil {
		config.ARIPassword = fmt.Sprintf("%v", password)
	}

	if config.ARIUser == "" || config.ARIPassword == "" {
		return nil, fmt.Errorf("ARI credentials (ari_user, ari_password) are required")
	}

	return config, nil
}
