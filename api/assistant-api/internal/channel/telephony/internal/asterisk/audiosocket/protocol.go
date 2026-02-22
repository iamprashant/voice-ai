// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_asterisk_audiosocket

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
)

// AudioSocket frame types (Asterisk res_audiosocket)
const (
	FrameTypeHangup  byte = 0x00
	FrameTypeUUID    byte = 0x01
	FrameTypeSilence byte = 0x02
	FrameTypeAudio   byte = 0x10
	FrameTypeError   byte = 0xFF
)

const maxFrameSize = 65535

// Frame represents a single AudioSocket frame.
type Frame struct {
	Type    byte
	Payload []byte
}

// ReadFrame reads a single frame from the AudioSocket stream.
// Frame format: 1-byte type + 2-byte big-endian length + payload.
func ReadFrame(r *bufio.Reader) (*Frame, error) {
	frameType, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	lenBuf := make([]byte, 2)
	if _, err := io.ReadFull(r, lenBuf); err != nil {
		return nil, err
	}

	payloadLen := int(binary.BigEndian.Uint16(lenBuf))
	if payloadLen < 0 || payloadLen > maxFrameSize {
		return nil, fmt.Errorf("invalid frame length: %d", payloadLen)
	}

	payload := make([]byte, payloadLen)
	if payloadLen > 0 {
		if _, err := io.ReadFull(r, payload); err != nil {
			return nil, err
		}
	}

	return &Frame{Type: frameType, Payload: payload}, nil
}

// WriteFrame writes a single frame to the AudioSocket stream.
// Frame format: 1-byte type + 2-byte big-endian length + payload.
func WriteFrame(w io.Writer, frameType byte, payload []byte) error {
	if len(payload) > maxFrameSize {
		return fmt.Errorf("frame payload too large: %d", len(payload))
	}

	header := []byte{frameType, 0x00, 0x00}
	binary.BigEndian.PutUint16(header[1:], uint16(len(payload)))

	if _, err := w.Write(header); err != nil {
		return err
	}
	if len(payload) > 0 {
		if _, err := w.Write(payload); err != nil {
			return err
		}
	}
	return nil
}
