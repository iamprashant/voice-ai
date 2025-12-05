// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_transformer_assemblyai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	internal_transformer "github.com/rapidaai/api/assistant-api/internal/transformer"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

type assemblyaiSTT struct {
	*assemblyaiOption
	ctx        context.Context
	mu         sync.Mutex
	logger     commons.Logger
	connection *websocket.Conn
	options    *internal_transformer.SpeechToTextInitializeOptions
}

func NewAssemblyaiSpeechToText(
	ctx context.Context,
	logger commons.Logger,
	credential *protos.VaultCredential,
	iOption *internal_transformer.SpeechToTextInitializeOptions,
) (internal_transformer.SpeechToTextTransformer, error) {
	ayOptions, err := NewAssemblyaiOption(
		logger,
		credential,
		iOption.AudioConfig,
		iOption.ModelOptions,
	)
	if err != nil {
		logger.Errorf("assembly-ai-stt: key from credential failed %v", err)
		return nil, err
	}
	return &assemblyaiSTT{
		ctx:              ctx,
		logger:           logger,
		options:          iOption,
		assemblyaiOption: ayOptions,
	}, nil
}

func (aai *assemblyaiSTT) Name() string {
	return "assemblyai-speech-to-text"
}

func (aai *assemblyaiSTT) Initialize() error {
	headers := http.Header{}
	headers.Set("Authorization", aai.GetKey())
	dialer := websocket.Dialer{
		Proxy:            nil,              // Skip proxy for direct connection
		HandshakeTimeout: 10 * time.Second, // Reduced handshake timeout for quick failover
	}
	conenction, _, err := dialer.Dial(aai.GetSpeechToTextConnectionString(), headers)
	if err != nil {
		aai.logger.Errorf("assembly-ai-stt: key from credential failed %v", err)
		return fmt.Errorf("failed to connect to assemblyai websocket: %w", err)
	}
	aai.connection = conenction
	go aai.speechToTextCallback()
	return nil
}

func (aai *assemblyaiSTT) speechToTextCallback() {
	for {
		select {
		case <-aai.ctx.Done():
			aai.logger.Info("assembly-ai-stt: read goroutine exiting due to context cancellation")
			return
		default:
			_, msg, err := aai.connection.ReadMessage()
			if err != nil {
				aai.logger.Error("assembly-ai-stt: read error: ", err)
				return
			}
			var transcript TranscriptMessage
			if err := json.Unmarshal(msg, &transcript); err != nil {
				continue
			}
			switch transcript.Type {
			case "Turn":
				confidence := 0.0
				for _, v := range transcript.Words {
					confidence += v.Confidence
				}
				averageConfidence := confidence / float64(len(transcript.Words))
				if transcript.EndOfTurn {
					aai.options.OnTranscript(
						transcript.Transcript,
						averageConfidence,
						"en",
						true,
					)
				} else {
					aai.options.OnTranscript(
						transcript.Transcript,
						averageConfidence,
						"en",
						false,
					)
				}
			case "Begin":
			default:
			}
		}
	}
}
func (aai *assemblyaiSTT) Transform(ctx context.Context, in []byte, opts *internal_transformer.SpeechToTextOption) error {
	if aai.connection == nil {
		return fmt.Errorf("assembly-ai-stt: websocket connection is not initialized")
	}
	if err := aai.connection.WriteMessage(websocket.BinaryMessage, in); err != nil {
		return fmt.Errorf("error sending audio: %w", err)
	}
	return nil
}

func (aai *assemblyaiSTT) Close(ctx context.Context) error {
	if aai.connection != nil {
		return aai.connection.Close()
	}
	return nil
}
