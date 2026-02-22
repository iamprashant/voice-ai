// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_asterisk_audiosocket

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"

	callcontext "github.com/rapidaai/api/assistant-api/internal/callcontext"
	internal_telephony_base "github.com/rapidaai/api/assistant-api/internal/channel/telephony/internal/base"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

// Streamer implements AudioSocket media streaming over TCP.
type Streamer struct {
	internal_telephony_base.BaseTelephonyStreamer

	conn           net.Conn
	reader         *bufio.Reader
	writer         *bufio.Writer
	writeMu        sync.Mutex
	audioProcessor *AudioProcessor

	// AudioSocket manages its own context for output lifecycle control.
	ctx          context.Context
	cancel       context.CancelFunc
	outputCtx    context.Context
	outputCancel context.CancelFunc

	initialUUID string
	configSent  bool
}

// NewStreamer creates a new AudioSocket streamer.
// initialUUID is the contextId already read from the first UUID frame by the AudioSocket
// engine — when set, the streamer emits ConversationInitialization on the first Recv()
// without waiting for another UUID frame from the wire.
func NewStreamer(
	logger commons.Logger,
	conn net.Conn,
	reader *bufio.Reader,
	writer *bufio.Writer,
	cc *callcontext.CallContext,
	vaultCred *protos.VaultCredential,
) (internal_type.Streamer, error) {
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
		BaseTelephonyStreamer: internal_telephony_base.NewBaseTelephonyStreamer(
			logger, cc, vaultCred,
		),
		conn:           conn,
		reader:         reader,
		writer:         writer,
		audioProcessor: audioProcessor,
		outputCtx:      ctx,
		outputCancel:   cancel,
		ctx:            ctx,
		cancel:         cancel,
		initialUUID:    cc.ContextID,
	}

	audioProcessor.SetInputAudioCallback(as.sendProcessedInputAudio)
	audioProcessor.SetOutputChunkCallback(as.sendAudioChunk)
	go audioProcessor.RunOutputSender(as.outputCtx)
	return as, nil
}

func (as *Streamer) sendProcessedInputAudio(audio []byte) {
	as.WithInputBuffer(func(buf *bytes.Buffer) {
		buf.Write(audio)
	})
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
func (as *Streamer) Recv() (internal_type.Stream, error) {
	if as.conn == nil {
		return nil, io.EOF
	}
	if !as.configSent && as.initialUUID != "" {
		as.configSent = true
		return as.CreateConnectionRequest(), nil
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
				return &protos.ConversationInitialization{
					AssistantConversationId: as.GetConversationId(),
					Assistant:               as.GetAssistantDefinition(),
					StreamMode:              protos.StreamMode_STREAM_MODE_AUDIO,
				}, nil
			}
		case FrameTypeAudio:
			if err := as.audioProcessor.ProcessInputAudio(frame.Payload); err != nil {
				as.Logger.Debug("Failed to process input audio", "error", err.Error())
				continue
			}

			var audioRequest *protos.ConversationUserMessage
			as.WithInputBuffer(func(buf *bytes.Buffer) {
				if buf.Len() > 0 {
					audioRequest = &protos.ConversationUserMessage{
						Message: &protos.ConversationUserMessage_Audio{Audio: buf.Bytes()},
					}
					buf.Reset()
				}
			})
			if audioRequest != nil {
				return audioRequest, nil
			}
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
func (as *Streamer) Send(response internal_type.Stream) error {
	switch data := response.(type) {
	case *protos.ConversationAssistantMessage:
		switch content := data.GetMessage().(type) {
		case *protos.ConversationAssistantMessage_Audio:
			if err := as.audioProcessor.ProcessOutputAudio(content.Audio); err != nil {
				return err
			}
		}
	case *protos.ConversationInterruption:
		if data.GetType() == protos.ConversationInterruption_INTERRUPTION_TYPE_WORD {
			as.audioProcessor.ClearInputBuffer()
			as.audioProcessor.ClearOutputBuffer()
		}
	case *protos.ConversationDirective:
		if data.GetType() == protos.ConversationDirective_END_CONVERSATION {
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
