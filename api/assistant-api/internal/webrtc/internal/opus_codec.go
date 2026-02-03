// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package webrtc_internal

import (
	"encoding/binary"
	"fmt"

	"gopkg.in/hraban/opus.v2"
)

const opusFrameSamples = 960 // 20ms at 48kHz

// OpusCodec handles Opus audio encoding/decoding for WebRTC (48kHz mono)
type OpusCodec struct {
	encoder *opus.Encoder
	decoder *opus.Decoder
}

// NewOpusCodec creates a new Opus codec optimized for voice
func NewOpusCodec() (*OpusCodec, error) {
	enc, err := opus.NewEncoder(OpusSampleRate, 1, opus.AppVoIP)
	if err != nil {
		return nil, fmt.Errorf("failed to create Opus encoder: %w", err)
	}

	enc.SetBitrate(24000)
	enc.SetComplexity(5)
	enc.SetInBandFEC(true)
	enc.SetPacketLossPerc(10)

	dec, err := opus.NewDecoder(OpusSampleRate, 1)
	if err != nil {
		return nil, fmt.Errorf("failed to create Opus decoder: %w", err)
	}

	return &OpusCodec{encoder: enc, decoder: dec}, nil
}

// Encode encodes PCM16 bytes (48kHz mono, little-endian) to Opus
func (c *OpusCodec) Encode(pcm []byte) ([]byte, error) {
	if len(pcm) == 0 {
		return nil, nil
	}

	numSamples := len(pcm) / 2
	samples := make([]int16, numSamples)
	for i := 0; i < numSamples; i++ {
		samples[i] = int16(binary.LittleEndian.Uint16(pcm[i*2 : i*2+2]))
	}

	output := make([]byte, 1000)
	n, err := c.encoder.Encode(samples, output)
	if err != nil {
		return nil, fmt.Errorf("Opus encode failed: %w", err)
	}

	return output[:n], nil
}

// Decode decodes Opus to PCM16 bytes (48kHz mono, little-endian)
func (c *OpusCodec) Decode(encoded []byte) ([]byte, error) {
	if len(encoded) == 0 {
		return nil, nil
	}

	samples := make([]int16, opusFrameSamples)
	n, err := c.decoder.Decode(encoded, samples)
	if err != nil {
		return nil, fmt.Errorf("Opus decode failed: %w", err)
	}

	pcm := make([]byte, n*2)
	for i := 0; i < n; i++ {
		binary.LittleEndian.PutUint16(pcm[i*2:i*2+2], uint16(samples[i]))
	}

	return pcm, nil
}
