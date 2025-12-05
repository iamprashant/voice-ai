// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_transformer_resemble

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	internal_transformer "github.com/rapidaai/api/assistant-api/internal/transformer"
	"github.com/rapidaai/pkg/commons"
	protos "github.com/rapidaai/protos"
)

type resembleTTS struct {
	*resembleOption
	ctx       context.Context
	mu        sync.Mutex
	contextId string

	logger     commons.Logger
	connection *websocket.Conn
	options    *internal_transformer.TextToSpeechInitializeOptions
}

func NewResembleTextToSpeech(
	ctx context.Context,
	logger commons.Logger,
	credential *protos.VaultCredential,
	options *internal_transformer.TextToSpeechInitializeOptions,
) (internal_transformer.TextToSpeechTransformer, error) {
	rsmblOpts, err := NewResembleOption(
		logger,
		credential,
		options.AudioConfig,
		options.ModelOptions)
	if err != nil {
		logger.Errorf("resemble-tts: intializing resembleai failed %+v", err)
		return nil, err
	}
	return &resembleTTS{
		resembleOption: rsmblOpts,
		ctx:            ctx,
		logger:         logger,
		options:        options,
	}, nil
}

// Initialize implements internal_transformer.OutputAudioTransformer.
func (rt *resembleTTS) Initialize() error {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	headers := map[string][]string{
		"Authorization": {"Bearer " + rt.GetKey()},
	}
	conn, _, err := websocket.DefaultDialer.Dial("wss://websocket.cluster.resemble.ai/stream", headers)
	if err != nil {
		rt.logger.Errorf("resemble-tts: unable to connect to websocket err: %v", err)
		return err
	}
	rt.connection = conn
	go rt.textToSpeechCallback(rt.ctx)
	return nil
}

// Name implements internal_transformer.OutputAudioTransformer.
func (*resembleTTS) Name() string {
	return "resemble-text-to-speech"
}

func (rt *resembleTTS) textToSpeechCallback(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			rt.logger.Infof("sarvam-tts: context cancelled, stopping response listener")
			return
		default:
			_, audioChunk, err := rt.connection.ReadMessage()
			if err != nil {
				rt.logger.Errorf("resemble-tts: error reading from Resemble WebSocket: %v", err)
				continue
			}

			var audioData map[string]interface{}
			if err := json.Unmarshal(audioChunk, &audioData); err != nil {
				rt.logger.Errorf("resemble-tts: error parsing audio chunk: %v", err)
				continue
			}
			if audioData["type"] == "audio_end" {
				break
			}
			if audioData["type"] == "audio" {
				payload, ok := audioData["audio_content"].(string)
				if !ok {
					continue
				}
				rawAudioData, err := base64.StdEncoding.DecodeString(payload)
				if err != nil {
					log.Fatalf("Error decoding base64 string: %v", err)
				}
				rt.options.OnSpeech(rt.contextId, rawAudioData)
			}
		}
	}

}

func (rt *resembleTTS) Transform(ctx context.Context, in string, opts *internal_transformer.TextToSpeechOption) error {
	rt.mu.Lock()
	defer rt.mu.Unlock()

	rt.contextId = opts.ContextId
	if err := rt.connection.WriteJSON(rt.GetTextToSpeechRequest(opts.ContextId, in)); err != nil {
		rt.logger.Errorf("resemble-tts: error while writing request to websocket %v", err)
		return err
	}
	return nil
}

func (rt *resembleTTS) Close(ctx context.Context) error {
	rt.connection.Close()
	return nil
}
