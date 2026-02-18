// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_recorder

import (
	"bytes"
	"context"
	"encoding/binary"
	"testing"

	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
)

func newTestLogger(t *testing.T) commons.Logger {
	t.Helper()
	logger, err := commons.NewApplicationLogger(
		commons.Name("test-recorder"),
		commons.Path(t.TempDir()),
		commons.Level("debug"),
	)
	if err != nil {
		t.Fatalf("failed to create test logger: %v", err)
	}
	return logger
}

func newTestRecorder(t *testing.T) *audioRecorder {
	t.Helper()
	rec, err := NewDefaultAudioRecorder(newTestLogger(t))
	if err != nil {
		t.Fatalf("failed to create recorder: %v", err)
	}
	return rec.(*audioRecorder)
}

func pcm(val byte, length int) []byte {
	buf := make([]byte, length)
	for i := range buf {
		buf[i] = val
	}
	return buf
}

func wavPCMData(wav []byte) []byte { return wav[44:] }

func TestRecordUserAudio(t *testing.T) {
	rec := newTestRecorder(t)
	data := pcm(0x01, 320)
	rec.Record(context.Background(), internal_type.UserAudioPacket{Audio: data})

	if len(rec.chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(rec.chunks))
	}
	if rec.chunks[0].Track != trackUser {
		t.Errorf("expected trackUser")
	}
	if !bytes.Equal(rec.chunks[0].Data, data) {
		t.Errorf("data mismatch")
	}
}

func TestRecordSystemAudio(t *testing.T) {
	rec := newTestRecorder(t)
	rec.Record(context.Background(), internal_type.TextToSpeechAudioPacket{ContextID: "c1", AudioChunk: pcm(0x02, 640)})

	if len(rec.chunks) != 1 || rec.chunks[0].Track != trackSystem {
		t.Errorf("expected 1 system chunk")
	}
}

func TestRecordEmptyDataIsIgnored(t *testing.T) {
	rec := newTestRecorder(t)
	ctx := context.Background()
	rec.Record(ctx, internal_type.UserAudioPacket{Audio: nil})
	rec.Record(ctx, internal_type.UserAudioPacket{Audio: []byte{}})
	rec.Record(ctx, internal_type.TextToSpeechAudioPacket{ContextID: "c", AudioChunk: nil})

	if len(rec.chunks) != 0 {
		t.Fatalf("expected 0 chunks, got %d", len(rec.chunks))
	}
}

func TestTTSBurstChunksPreserveOrder(t *testing.T) {
	rec := newTestRecorder(t)
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		rec.Record(ctx, internal_type.TextToSpeechAudioPacket{
			ContextID:  "c1",
			AudioChunk: pcm(byte(i+1), 320),
		})
	}
	if len(rec.chunks) != 5 {
		t.Fatalf("expected 5 chunks, got %d", len(rec.chunks))
	}
	for i, c := range rec.chunks {
		if c.Data[0] != byte(i+1) {
			t.Errorf("chunk %d: expected first byte %d, got %d", i, i+1, c.Data[0])
		}
		if c.Track != trackSystem {
			t.Errorf("chunk %d: expected trackSystem", i)
		}
	}
}

func TestPushCopiesData(t *testing.T) {
	rec := newTestRecorder(t)
	data := pcm(0xFF, 100)
	rec.Record(context.Background(), internal_type.UserAudioPacket{Audio: data})
	data[0] = 0x00
	if rec.chunks[0].Data[0] != 0xFF {
		t.Error("push must copy data")
	}
}

func TestPersistEmptyReturnsError(t *testing.T) {
	rec := newTestRecorder(t)
	if _, _, err := rec.Persist(); err == nil {
		t.Fatal("expected error for empty recorder")
	}
}

func TestPersistProducesValidWAV(t *testing.T) {
	rec := newTestRecorder(t)
	ctx := context.Background()
	rec.Record(ctx, internal_type.UserAudioPacket{Audio: pcm(0x01, 3200)})
	rec.Record(ctx, internal_type.TextToSpeechAudioPacket{ContextID: "c1", AudioChunk: pcm(0x02, 6400)})

	userWAV, systemWAV, err := rec.Persist()
	if err != nil {
		t.Fatalf("Persist error: %v", err)
	}
	for name, wav := range map[string][]byte{"user": userWAV, "system": systemWAV} {
		if len(wav) < 44 {
			t.Fatalf("%s WAV too short", name)
		}
		if string(wav[0:4]) != "RIFF" || string(wav[8:12]) != "WAVE" {
			t.Errorf("%s WAV missing RIFF/WAVE header", name)
		}
		if sr := binary.LittleEndian.Uint32(wav[24:28]); sr != audioConfig.SampleRate {
			t.Errorf("%s sample rate: got %d", name, sr)
		}
	}
	// Both tracks must have same length
	if len(wavPCMData(userWAV)) != len(wavPCMData(systemWAV)) {
		t.Error("user and system WAV lengths differ")
	}
	// Total PCM = user chunk + system chunk
	if got := len(wavPCMData(userWAV)); got != 3200+6400 {
		t.Errorf("expected %d PCM bytes, got %d", 3200+6400, got)
	}
}

func TestPersistSilenceFilling(t *testing.T) {
	rec := newTestRecorder(t)
	ctx := context.Background()
	rec.Record(ctx, internal_type.UserAudioPacket{Audio: pcm(0x11, 100)})
	rec.Record(ctx, internal_type.TextToSpeechAudioPacket{ContextID: "c1", AudioChunk: pcm(0x22, 200)})

	userWAV, systemWAV, _ := rec.Persist()
	userPCM := wavPCMData(userWAV)
	systemPCM := wavPCMData(systemWAV)

	// User track: 100 bytes audio, 200 bytes silence
	for i := 0; i < 100; i++ {
		if userPCM[i] != 0x11 {
			t.Errorf("user byte %d: expected 0x11, got 0x%02x", i, userPCM[i])
			break
		}
	}
	for i := 100; i < 300; i++ {
		if userPCM[i] != 0x00 {
			t.Errorf("user byte %d: expected silence, got 0x%02x", i, userPCM[i])
			break
		}
	}
	// System track: 100 bytes silence, 200 bytes audio
	for i := 0; i < 100; i++ {
		if systemPCM[i] != 0x00 {
			t.Errorf("system byte %d: expected silence, got 0x%02x", i, systemPCM[i])
			break
		}
	}
	for i := 100; i < 300; i++ {
		if systemPCM[i] != 0x22 {
			t.Errorf("system byte %d: expected 0x22, got 0x%02x", i, systemPCM[i])
			break
		}
	}
}
