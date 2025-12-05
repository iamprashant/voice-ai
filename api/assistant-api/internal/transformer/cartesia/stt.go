// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_transformer_cartesia

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	internal_transformer "github.com/rapidaai/api/assistant-api/internal/transformer"
	"github.com/rapidaai/pkg/commons"
	protos "github.com/rapidaai/protos"
)

type cartesiaSpeechToText struct {
	*cartesiaOption
	mu                 sync.Mutex
	logger             commons.Logger
	ctx                context.Context
	connection         *websocket.Conn
	transformerOptions *internal_transformer.SpeechToTextInitializeOptions
}

// Name implements internal_transformer.SpeechToTextTransformer.
func (*cartesiaSpeechToText) Name() string {
	return "cartesia-speech-to-text"
}

func NewCartesiaSpeechToText(ctx context.Context,
	logger commons.Logger,
	credential *protos.VaultCredential,
	transformerOptions *internal_transformer.SpeechToTextInitializeOptions,
) (internal_transformer.SpeechToTextTransformer, error) {
	cartesiaOpts, err := NewCartesiaOption(logger,
		credential,
		transformerOptions.AudioConfig,
		transformerOptions.ModelOptions)
	if err != nil {
		logger.Errorf("cartesia-stt: intializing cartesia failed %+v", err)
		return nil, err
	}

	return &cartesiaSpeechToText{
		ctx:                ctx,
		logger:             logger,
		cartesiaOption:     cartesiaOpts,
		transformerOptions: transformerOptions,
	}, nil
}

// textToSpeechCallback processes streaming responses asynchronously.
func (cst *cartesiaSpeechToText) speechToTextCallback(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			cst.logger.Infof("cartesia-tts: context cancelled, stopping response listener")
			return
		default:
			_, msg, err := cst.connection.ReadMessage()
			if err != nil {
				cst.logger.Error("cartesia-tts: error reading from Cartesia WebSocket: ", err)
				return
			}
			var resp SpeechToTextOutput
			if err := json.Unmarshal(msg, &resp); err == nil && resp.Text != "" {
				cst.logger.Debug("cartesia-tts: received transcription: %+v", resp)
				if cst.transformerOptions.OnTranscript != nil {
					cst.transformerOptions.OnTranscript(
						resp.Text,
						0.9,
						resp.Language,
						resp.IsFinal,
					)
				}
			}
		}
	}
}

func (cst *cartesiaSpeechToText) Initialize() error {
	cst.mu.Lock()
	defer cst.mu.Unlock()

	conn, _, err := websocket.DefaultDialer.Dial(cst.GetSpeechToTextConnectionString(), nil)
	if err != nil {
		return fmt.Errorf("cartesia-stt: failed to connect to Cartesia WebSocket: %w", err)
	}
	cst.connection = conn
	go cst.speechToTextCallback(cst.ctx)
	return nil
}

func (cst *cartesiaSpeechToText) Transform(ctx context.Context, in []byte, opts *internal_transformer.SpeechToTextOption) error {
	cst.mu.Lock()
	defer cst.mu.Unlock()

	if cst.connection == nil {
		return fmt.Errorf("cartesia-stt: websocket connection is not initialized")
	}
	if err := cst.connection.WriteMessage(
		websocket.BinaryMessage, in); err != nil {
		return fmt.Errorf("failed to send audio data: %w", err)
	}

	return nil
}

func (cst *cartesiaSpeechToText) Close(ctx context.Context) error {
	if cst.connection != nil {
		err := cst.connection.Close()
		if err != nil {
			return fmt.Errorf("error closing WebSocket connection: %w", err)
		}
		cst.logger.Info("cartesia-stt: cartesia websocket connection closed")
	}
	return nil
}
