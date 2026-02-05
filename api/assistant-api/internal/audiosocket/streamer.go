// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package audiosocket

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

// Streamer implements AudioSocket media streaming over TCP.
type Streamer struct {
	logger         commons.Logger
	conn           net.Conn
	reader         *bufio.Reader
	writer         *bufio.Writer
	writeMu        sync.Mutex
	audioProcessor *AudioProcessor
	assistant      *internal_assistant_entity.Assistant
	conversation   *internal_conversation_entity.AssistantConversation
	inputBuffer    *bytes.Buffer
	inputMu        sync.Mutex
	ctx            context.Context
	cancel         context.CancelFunc

	outputCtx    context.Context
	outputCancel context.CancelFunc

	initialUUID string
	configSent  bool
}

// NewStreamer creates a new AudioSocket streamer.
func NewStreamer(
	logger commons.Logger,
	conn net.Conn,
	reader *bufio.Reader,
	writer *bufio.Writer,
	assistant *internal_assistant_entity.Assistant,
	conversation *internal_conversation_entity.AssistantConversation,
	vlt *protos.VaultCredential,
) (internal_type.TelephonyStreamer, error) {
	audioProcessor, err := NewAudioProcessor(logger)
	if err != nil {
		return nil, err
	}

	if reader == nil {
		reader = bufio.NewReader(conn)
	}
	if writer == nil {
		writer = bufio.NewWriter(conn)
	}

	ctx, cancel := context.WithCancel(context.Background())

	as := &Streamer{
		logger:         logger,
		conn:           conn,
		reader:         reader,
		writer:         writer,
		audioProcessor: audioProcessor,
		outputCtx:      ctx,
		outputCancel:   cancel,
		assistant:      assistant,
		conversation:   conversation,
		inputBuffer:    new(bytes.Buffer),
		ctx:            ctx,
		cancel:         cancel,
	}

	audioProcessor.SetInputAudioCallback(as.sendProcessedInputAudio)
	audioProcessor.SetOutputChunkCallback(as.sendAudioChunk)

	go audioProcessor.RunOutputSender(as.outputCtx)

	return as, nil
}

// SetInitialUUID sets the AudioSocket UUID before streaming starts.
func (as *Streamer) SetInitialUUID(uuid string) {
	as.initialUUID = uuid
}

func (as *Streamer) sendProcessedInputAudio(audio []byte) {
	as.inputMu.Lock()
	as.inputBuffer.Write(audio)
	as.inputMu.Unlock()
}

func (as *Streamer) sendAudioChunk(chunk *AudioChunk) error {
	if as.conn == nil {
		return nil
	}
	return as.writeFrame(FrameTypeAudio, chunk.Data)
}

func (as *Streamer) writeFrame(frameType byte, payload []byte) error {
	as.writeMu.Lock()
	defer as.writeMu.Unlock()

	if err := WriteFrame(as.writer, frameType, payload); err != nil {
		return err
	}
	return as.writer.Flush()
}

// Context returns the streamer context.
func (as *Streamer) Context() context.Context {
	return as.ctx
}

// Recv reads AudioSocket frames and returns input for the talker.
func (as *Streamer) Recv() (*protos.AssistantTalkInput, error) {
	if as.conn == nil {
		return nil, io.EOF
	}

	if !as.configSent && as.initialUUID != "" {
		as.configSent = true
		downstreamConfig := as.audioProcessor.GetDownstreamConfig()
		return as.createConnectionRequest(downstreamConfig, downstreamConfig)
	}

	for {
		frame, err := ReadFrame(as.reader)
		if err != nil {
			if err == io.EOF {
				return nil, io.EOF
			}
			return nil, err
		}

		switch frame.Type {
		case FrameTypeUUID:
			as.initialUUID = strings.TrimSpace(string(frame.Payload))
			if !as.configSent {
				as.configSent = true
				downstreamConfig := as.audioProcessor.GetDownstreamConfig()
				return as.createConnectionRequest(downstreamConfig, downstreamConfig)
			}
		case FrameTypeAudio:
			if err := as.audioProcessor.ProcessInputAudio(frame.Payload); err != nil {
				as.logger.Debug("Failed to process input audio", "error", err.Error())
				continue
			}

			as.inputMu.Lock()
			if as.inputBuffer.Len() > 0 {
				audioRequest := as.createVoiceRequest(as.inputBuffer.Bytes())
				as.inputBuffer.Reset()
				as.inputMu.Unlock()
				return audioRequest, nil
			}
			as.inputMu.Unlock()
		case FrameTypeSilence:
			// Silence frame, no action needed
		case FrameTypeHangup:
			return nil, io.EOF
		case FrameTypeError:
			return nil, fmt.Errorf("audiosocket error frame received")
		default:
			// Ignore unknown frame types
		}
	}
}

// Send writes audio/output frames back to Asterisk.
func (as *Streamer) Send(response *protos.AssistantTalkOutput) error {
	switch data := response.GetData().(type) {
	case *protos.AssistantTalkOutput_Assistant:
		switch content := data.Assistant.Message.(type) {
		case *protos.ConversationAssistantMessage_Audio:
			if err := as.audioProcessor.ProcessOutputAudio(content.Audio); err != nil {
				return err
			}
		}
	case *protos.AssistantTalkOutput_Interruption:
		if data.Interruption.Type == protos.ConversationInterruption_INTERRUPTION_TYPE_WORD {
			as.audioProcessor.ClearInputBuffer()
			as.audioProcessor.ClearOutputBuffer()
		}
	case *protos.AssistantTalkOutput_Directive:
		if data.Directive.GetType() == protos.ConversationDirective_END_CONVERSATION {
			_ = as.writeFrame(FrameTypeHangup, nil)
			return as.close()
		}
	}

	return nil
}

func (as *Streamer) close() error {
	if as.outputCancel != nil {
		as.outputCancel()
	}
	if as.cancel != nil {
		as.cancel()
	}
	if as.conn != nil {
		_ = as.conn.Close()
		as.conn = nil
	}
	return nil
}

func (as *Streamer) createVoiceRequest(audioData []byte) *protos.AssistantTalkInput {
	return &protos.AssistantTalkInput{
		Request: &protos.AssistantTalkInput_Message{
			Message: &protos.ConversationUserMessage{
				Message: &protos.ConversationUserMessage_Audio{
					Audio: audioData,
				},
			},
		},
	}
}

func (as *Streamer) createConnectionRequest(in, out *protos.AudioConfig) (*protos.AssistantTalkInput, error) {
	return &protos.AssistantTalkInput{
		Request: &protos.AssistantTalkInput_Configuration{
			Configuration: &protos.ConversationConfiguration{
				AssistantConversationId: as.conversation.Id,
				Assistant: &protos.AssistantDefinition{
					AssistantId: as.assistant.Id,
					Version:     utils.GetVersionString(as.assistant.AssistantProviderId),
				},
				InputConfig:  &protos.StreamConfig{Audio: in},
				OutputConfig: &protos.StreamConfig{Audio: out},
			},
		},
	}, nil
}
