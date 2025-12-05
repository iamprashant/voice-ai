// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_transformer_elevenlabs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	internal_transformer "github.com/rapidaai/api/assistant-api/internal/transformer"
	"github.com/rapidaai/pkg/commons"
	protos "github.com/rapidaai/protos"
)

type elevenlabsTTS struct {
	*elevenLabsOption
	ctx        context.Context
	mu         sync.Mutex
	logger     commons.Logger
	connection *websocket.Conn
	options    *internal_transformer.TextToSpeechInitializeOptions
}

func NewElevenlabsTextToSpeech(
	ctx context.Context,
	logger commons.Logger,
	credential *protos.VaultCredential,
	opts *internal_transformer.TextToSpeechInitializeOptions) (internal_transformer.TextToSpeechTransformer, error) {
	eleOpts, err := NewElevenLabsOption(
		logger,
		credential,
		opts.AudioConfig,
		opts.ModelOptions)
	if err != nil {
		logger.Errorf("elevenlabs-tts: intializing elevenlabs failed %+v", err)
		return nil, err
	}

	return &elevenlabsTTS{
		ctx:              ctx,
		options:          opts,
		logger:           logger,
		elevenLabsOption: eleOpts,
	}, nil
}

// Initialize implements internal_transformer.OutputAudioTransformer.
func (ct *elevenlabsTTS) Initialize() error {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	header := http.Header{}
	header.Set("xi-api-key", ct.GetKey())
	conn, resp, err := websocket.DefaultDialer.Dial(ct.GetTextToSpeechConnectionString(), header)
	if err != nil {
		ct.logger.Errorf("elevenlab-tts: error while elevenlabs %s with response %v", err, resp)
		return err
	}

	ct.connection = conn
	go ct.textToSpeechCallback(ct.ctx)
	return nil
}

// Name implements internal_transformer.OutputAudioTransformer.
func (*elevenlabsTTS) Name() string {
	return "elevenlabs-text-to-speech"
}

func (elt *elevenlabsTTS) textToSpeechCallback(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			elt.logger.Infof("elevenlabs-tts: context cancelled, stopping response listener")
			return
		default:
			_, audioChunk, err := elt.connection.ReadMessage()
			if err != nil {
				elt.logger.Errorf("elevenlab-tts: Error reading from TTS WebSocket: %v", err)
			}

			var audioData map[string]interface{}
			if err := json.Unmarshal(audioChunk, &audioData); err != nil {
				elt.logger.Errorf("elevenlab-tts: Error parsing audio chunk: %v", err)
				break
			}

			contextId, ok := audioData["contextId"].(string)
			if !ok {
				continue
			}
			done, ok := audioData["isFinal"].(bool)
			if ok && done {
				elt.options.OnComplete(contextId)
			}
			payload, ok := audioData["audio"].(string)
			if !ok {
				continue
			}
			rawAudioData, err := base64.StdEncoding.DecodeString(payload)
			if err != nil {
				elt.logger.Errorf("elevenlab-tts: Error decoding base64 string: %v", err)
			}
			elt.options.OnSpeech(contextId, rawAudioData)
		}
	}

}

func (t *elevenlabsTTS) Transform(ctx context.Context, in string, opts *internal_transformer.TextToSpeechOption) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.connection == nil {
		return fmt.Errorf("cartesia-stt: websocket connection is not initialized")
	}
	ttsMessage := map[string]interface{}{
		"text":       in,
		"context_id": opts.ContextId,
		"flush":      !opts.IsComplete,
	}

	if err := t.connection.WriteJSON(ttsMessage); err != nil {
		t.logger.Errorf("elevenlab-tts: unable to write json for text to speech: %v", err)
		return err
	}
	return nil
}

func (t *elevenlabsTTS) Close(ctx context.Context) error {
	if t.connection != nil {
		t.connection.Close()
	}
	return nil
}
