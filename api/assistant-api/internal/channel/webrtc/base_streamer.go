// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package channel_webrtc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sync"

	webrtc_internal "github.com/rapidaai/api/assistant-api/internal/channel/webrtc/internal"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ============================================================================
// baseStreamer — transport-agnostic channel & buffer management
// ============================================================================

// baseStreamer owns the input/output channels and audio buffers that every
// concrete streamer (WebRTC, telephony, …) needs. It handles:
//
//   - inputCh / outputCh: unified, ordered message channels
//   - inputAudioBuffer / outputAudioBuffer: PCM accumulation with thresholds
//   - flushAudioCh: interrupt signalling for the output writer
//   - pushInput / pushOutput: non-blocking channel sends
//   - clearInputBuffer / clearOutputBuffer: buffer + channel draining
//   - pushDisconnection: idempotent disconnect signalling
//   - Recv / Context: Streamer interface helpers
//
// The concrete streamer embeds baseStreamer and only implements
// transport-specific logic (WebRTC track I/O, gRPC dispatch, Opus encoding, etc.).
type baseStreamer struct {
	mu sync.Mutex

	// Core components
	logger commons.Logger

	// Lifecycle
	ctx    context.Context
	cancel context.CancelFunc

	// Disconnect tracking — true once pushDisconnection has run.
	closed bool

	// inputCh: all downstream-bound messages (gRPC + decoded audio) funnelled here.
	// recv (non-blocking) -> inputCh -> loop (Recv) -> downstream service
	inputCh              chan internal_type.Stream
	inputAudioBuffer     *bytes.Buffer
	inputAudioBufferLock sync.Mutex

	// outputCh: all upstream-bound messages funnelled here to preserve ordering.
	// send (non-blocking) -> outputCh -> loop (runOutputWriter) -> upstream service
	outputCh              chan internal_type.Stream
	outputAudioBuffer     *bytes.Buffer
	outputAudioBufferLock sync.Mutex

	// flushAudioCh signals the output writer to discard its pending audio queue
	// (used on interruption to silence stale frames immediately).
	flushAudioCh chan struct{}
}

// newBaseStreamer initialises a baseStreamer with channels and buffers.
// The streamer owns its own context (derived from context.Background) so that
// cleanup is never short-circuited by the caller's context being cancelled first.
func newBaseStreamer(logger commons.Logger) baseStreamer {
	ctx, cancel := context.WithCancel(context.Background())
	return baseStreamer{
		logger:            logger,
		ctx:               ctx,
		cancel:            cancel,
		inputCh:           make(chan internal_type.Stream, webrtc_internal.InputChannelSize),
		outputCh:          make(chan internal_type.Stream, webrtc_internal.OutputChannelSize),
		inputAudioBuffer:  new(bytes.Buffer),
		outputAudioBuffer: new(bytes.Buffer),
		flushAudioCh:      make(chan struct{}, 1),
	}
}

// ============================================================================
// Input buffer helpers
// ============================================================================

// bufferAndSendInput accumulates resampled audio and sends it to inputCh
// when the buffer reaches the threshold.
func (s *baseStreamer) bufferAndSendInput(audio []byte) {
	s.inputAudioBufferLock.Lock()
	s.inputAudioBuffer.Write(audio)

	if s.inputAudioBuffer.Len() < webrtc_internal.InputBufferThreshold {
		s.inputAudioBufferLock.Unlock()
		return
	}

	audioData := make([]byte, s.inputAudioBuffer.Len())
	s.inputAudioBuffer.Read(audioData)
	s.inputAudioBufferLock.Unlock()

	s.pushInput(&protos.ConversationUserMessage{
		Message: &protos.ConversationUserMessage_Audio{Audio: audioData},
		Time:    timestamppb.Now(),
	})
}

// clearInputBuffer resets the input PCM buffer and drains the input channel.
func (s *baseStreamer) clearInputBuffer() {
	s.inputAudioBufferLock.Lock()
	s.inputAudioBuffer.Reset()
	s.inputAudioBufferLock.Unlock()
	for {
		select {
		case <-s.inputCh:
		default:
			return
		}
	}
}

// ============================================================================
// Output buffer helpers
// ============================================================================

// bufferAndSendOutput accumulates resampled 48kHz PCM and flushes consistent
// 20ms frames into outputCh as ConversationAssistantMessage_Audio messages.
// Opus encoding happens later in the concrete streamer's output writer.
//
// audio received -> push to outputAudioBuffer -> check size ->
// push ConversationAssistantMessage_Audio -> outputCh
func (s *baseStreamer) bufferAndSendOutput(audio48kHz []byte) {
	s.outputAudioBufferLock.Lock()
	s.outputAudioBuffer.Write(audio48kHz)

	if s.outputAudioBuffer.Len() < webrtc_internal.OutputBufferThreshold {
		s.outputAudioBufferLock.Unlock()
		return
	}

	// Flush as many complete 20ms PCM frames as possible.
	for s.outputAudioBuffer.Len() >= webrtc_internal.OpusFrameBytes {
		frame := make([]byte, webrtc_internal.OpusFrameBytes)
		s.outputAudioBuffer.Read(frame)
		s.outputAudioBufferLock.Unlock()

		// Push raw PCM frame; encoding is deferred to the output writer.
		s.pushOutput(&protos.ConversationAssistantMessage{
			Message: &protos.ConversationAssistantMessage_Audio{Audio: frame},
			Time:    timestamppb.Now(),
		})

		s.outputAudioBufferLock.Lock()
	}
	s.outputAudioBufferLock.Unlock()
}

// clearOutputBuffer resets the output PCM buffer, signals the output writer
// to flush its pending audio queue, and drains the output channel.
func (s *baseStreamer) clearOutputBuffer() {
	// 1. Reset the PCM accumulation buffer so no new frames are encoded.
	s.outputAudioBufferLock.Lock()
	s.outputAudioBuffer.Reset()
	s.outputAudioBufferLock.Unlock()

	// 2. Signal the output writer to flush its local pending audio queue first,
	//    before draining outputCh, to prevent the writer from dequeuing a message
	//    between drain and signal.
	select {
	case s.flushAudioCh <- struct{}{}:
	default:
	}

	// 3. Drain the output channel (pending audio + gRPC messages).
	for {
		select {
		case <-s.outputCh:
		default:
			return
		}
	}
}

// ============================================================================
// Channel push helpers
// ============================================================================

// pushInput sends a message to the unified input channel (non-blocking).
// Safe to call after inputCh is closed — the send is guarded by the closed flag.
func (s *baseStreamer) pushInput(msg internal_type.Stream) {
	select {
	case s.inputCh <- msg:
	default:
		s.logger.Warnw("Input channel full, dropping message", "type", fmt.Sprintf("%T", msg))
	}
}

// pushOutput sends a message to the unified output channel (non-blocking).
func (s *baseStreamer) pushOutput(msg internal_type.Stream) {
	select {
	case s.outputCh <- msg:
	default:
		s.logger.Warnw("Output channel full, dropping message", "type", fmt.Sprintf("%T", msg))
	}
}

// ============================================================================
// Disconnect helpers
// ============================================================================

// pushDisconnection pushes a ConversationDisconnection into inputCh.
// It is idempotent — safe to call from multiple goroutines or multiple times.
// FIFO ordering guarantees the Talk loop processes any preceding metrics before
// the disconnection signal.
func (s *baseStreamer) pushDisconnection(reason protos.ConversationDisconnection_DisconnectionType) {
	s.mu.Lock()
	alreadyClosed := s.closed
	s.closed = true
	s.mu.Unlock()
	if alreadyClosed {
		return
	}

	s.pushInput(&protos.ConversationDisconnection{
		Type: reason,
		Time: timestamppb.Now(),
	})
}

// ============================================================================
// Streamer interface helpers (embedded by concrete streamers)
// ============================================================================

// Context returns the streamer-scoped context.
func (s *baseStreamer) Context() context.Context {
	return s.ctx
}

// Recv reads the next downstream-bound message from the unified input channel.
// Both gRPC messages and decoded WebRTC audio are fed into the same channel
// by background goroutines. Shutdown is signalled by a ConversationDisconnection
// message through inputCh, which the Talk loop handles to trigger Disconnect().
// No context select here — Close() pushes ConversationDisconnection first and
// cancels the context afterwards, so a competing select could skip the message.
func (s *baseStreamer) Recv() (internal_type.Stream, error) {
	select {
	case msg, ok := <-s.inputCh:
		if !ok {
			return nil, io.EOF
		}
		return msg, nil
	case <-s.ctx.Done():
		return nil, io.EOF
	}
}
