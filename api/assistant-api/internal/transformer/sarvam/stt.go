// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_transformer_sarvam

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/dvonthenen/websocket"
	internal_transformer "github.com/rapidaai/api/assistant-api/internal/transformer"
	sarvam_internal "github.com/rapidaai/api/assistant-api/internal/transformer/sarvam/internal"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

type sarvamSpeechToText struct {
	*sarvamOption
	// context management
	ctx       context.Context
	ctxCancel context.CancelFunc

	mu     sync.Mutex
	logger commons.Logger

	connection         *websocket.Conn
	transformerOptions *internal_transformer.SpeechToTextInitializeOptions
}

// Name implements internal_transformer.SpeechToTextTransformer.
func (*sarvamSpeechToText) Name() string {
	return "sarvam-speech-to-text"
}

func NewSarvamSpeechToText(
	ctx context.Context,
	logger commons.Logger,
	credential *protos.VaultCredential,
	opts *internal_transformer.SpeechToTextInitializeOptions) (internal_transformer.SpeechToTextTransformer, error) {
	sarvamOpts, err := NewSarvamOption(logger, credential, opts.AudioConfig, opts.ModelOptions)
	if err != nil {
		logger.Errorf("sarvam-stt: intializing sarvam failed %+v", err)
		return nil, err
	}

	ct, ctxCancel := context.WithCancel(ctx)
	return &sarvamSpeechToText{
		ctx:                ct,
		ctxCancel:          ctxCancel,
		logger:             logger,
		sarvamOption:       sarvamOpts,
		transformerOptions: opts,
	}, nil
}

func (cst *sarvamSpeechToText) speechToTextCallback(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			cst.logger.Infof("sarvam-stt: context cancelled, stopping response listener")
			return
		default:
			_, msg, err := cst.connection.ReadMessage()
			if err != nil {
				cst.logger.Error("sarvam-stt: error reading from WebSocket: ", err)
				return
			}
			var response struct {
				Type string          `json:"type"`
				Data json.RawMessage `json:"data"`
			}
			if err := json.Unmarshal(msg, &response); err != nil {
				cst.logger.Errorf("sarvam-stt: failed to unmarshal response: %v", err)
				continue
			}

			switch response.Type {
			case "data":
				var transcriptionData sarvam_internal.SpeechToTextTranscriptionData
				if err := json.Unmarshal(response.Data, &transcriptionData); err == nil {
					cst.logger.Debugf("sarvam-stt: transcription received: %+v", transcriptionData)
					if cst.transformerOptions.OnTranscript != nil {
						cst.transformerOptions.OnTranscript(
							transcriptionData.Transcript,
							0.9,
							transcriptionData.LanguageCode,
							true,
						)
					}
				}
			case "error":
				var errorData sarvam_internal.ErrorData
				if err := json.Unmarshal(response.Data, &errorData); err == nil {
					cst.logger.Errorf("sarvam-stt: error from server: %v", errorData)
				}
			case "events":
				cst.logger.Infof("sarvam-stt: event received: %s", string(response.Data))
			default:
				cst.logger.Warnf("sarvam-stt: unknown response type: %s", response.Type)
			}
		}
	}
}

func (cst *sarvamSpeechToText) Initialize() error {
	cst.mu.Lock()
	defer cst.mu.Unlock()

	headers := make(map[string][]string)
	headers["Api-Subscription-Key"] = []string{cst.GetKey()}
	conn, _, err := websocket.DefaultDialer.Dial(cst.speechToTextUrl(), headers)
	if err != nil {
		return fmt.Errorf("sarvam-stt: failed to connect to Sarvam WebSocket: %w", err)
	}

	cst.connection = conn
	go cst.speechToTextCallback(cst.ctx) // Start processing responses asynchronously
	return nil
}
func (cst *sarvamSpeechToText) Transform(ctx context.Context, in []byte, opts *internal_transformer.SpeechToTextOption) error {
	cst.mu.Lock()
	defer cst.mu.Unlock()

	if cst.connection == nil {
		return fmt.Errorf("sarvam-stt: websocket connection is not initialized")
	}

	in, err := cst.speechToTextMessage(in, opts)
	if err != nil {
		return fmt.Errorf("sarvam-stt: unable to encode byte to base64")
	}
	if err := cst.connection.WriteMessage(
		websocket.TextMessage, in); err != nil {
		return fmt.Errorf("failed to send audio data: %w", err)
	}

	return nil
}

func (cst *sarvamSpeechToText) Close(ctx context.Context) error {
	cst.ctxCancel()
	if cst.connection != nil {
		err := cst.connection.Close()
		if err != nil {
			return fmt.Errorf("error closing WebSocket connection: %w", err)
		}
		cst.logger.Info("sarvam-stt: sarvam websocket connection closed")
	}
	return nil
}
