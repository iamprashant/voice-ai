// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_transformer_cartesia

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	internal_transformer "github.com/rapidaai/api/assistant-api/internal/transformer"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

type cartesiaTTS struct {
	*cartesiaOption
	mu         sync.Mutex
	ctx        context.Context
	logger     commons.Logger
	connection *websocket.Conn
	options    *internal_transformer.TextToSpeechInitializeOptions
}

func NewCartesiaTextToSpeech(
	ctx context.Context,
	logger commons.Logger,
	credential *protos.VaultCredential,
	opts *internal_transformer.TextToSpeechInitializeOptions,
) (internal_transformer.TextToSpeechTransformer, error) {
	cartesiaOpts, err := NewCartesiaOption(logger, credential,
		opts.AudioConfig,
		opts.ModelOptions)
	if err != nil {
		logger.Errorf("intializing cartesia failed %+v", err)
		return nil, err
	}

	return &cartesiaTTS{
		cartesiaOption: cartesiaOpts,
		logger:         logger,
		ctx:            ctx,
		options:        opts,
	}, nil
}

func (ct *cartesiaTTS) Initialize() error {
	conn, _, err := websocket.DefaultDialer.Dial(ct.GetTextToSpeechConnectionString(), nil)
	if err != nil {
		return err
	}
	ct.connection = conn
	go ct.textToSpeechCallback(ct.ctx)
	return nil
}

func (cst *cartesiaTTS) textToSpeechCallback(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			cst.logger.Infof("cartesia-tts: context cancelled, stopping response listener")
			return
		default:
			_, msg, err := cst.connection.ReadMessage()
			if err != nil {
				return
			}
			var payload TextToSpeechOuput
			if err := json.Unmarshal(msg, &payload); err != nil {
				cst.logger.Errorf("cartesia-tts: invalid json from cartesia error : %v", err)
				continue
			}
			if payload.Done {
				_ = cst.options.OnComplete(payload.ContextID)
				continue
			}
			if payload.Data == "" {
				continue
			}
			decoded, err := base64.StdEncoding.DecodeString(payload.Data)
			if err != nil {
				cst.logger.Error("cartesia-tts: failed to decode audio payload error: %v", err)
				continue
			}
			_ = cst.options.OnSpeech(payload.ContextID, decoded)
		}
	}
}

// Name returns the name of this transformer.
func (*cartesiaTTS) Name() string {
	return "cartesia-text-to-speech"
}

func (ct *cartesiaTTS) Transform(ctx context.Context, in string, opts *internal_transformer.TextToSpeechOption) error {
	ct.mu.Lock()
	defer ct.mu.Unlock()

	if ct.connection == nil {
		return fmt.Errorf("cartesia-tts: websocket connection is not initialized")
	}
	message := ct.GetTextToSpeechInput(
		in,
		map[string]interface{}{
			"continue":   !opts.IsComplete,
			"context_id": opts.ContextId,
		})

	if err := ct.connection.WriteJSON(message); err != nil {
		return err
	}
	return nil
}

func (ct *cartesiaTTS) Close(ctx context.Context) error {
	if ct.connection != nil {
		_ = ct.connection.Close()
	}
	return nil
}
