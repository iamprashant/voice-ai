// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_transformer_sarvam

import (
	"fmt"

	internal_audio "github.com/rapidaai/api/assistant-api/internal/audio"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

const (
	SARVAM_URL = "wss://api.sarvam.ai/text-to-speech/ws"
	MODEL      = "bulbul:v2"
	VOICE      = "anushka"
)

type sarvamOption struct {
	logger      commons.Logger
	audioConfig *internal_audio.AudioConfig
	modelOpts   utils.Option
	key         string
}

func NewSarvamOption(logger commons.Logger,
	vaultCredential *protos.VaultCredential,
	audioConfig *internal_audio.AudioConfig, option utils.Option) (*sarvamOption, error) {

	cx, ok := vaultCredential.GetValue().AsMap()["key"]
	if !ok {
		return nil, fmt.Errorf("sarvam: illegal vault config")
	}
	return &sarvamOption{
		logger:      logger,
		audioConfig: audioConfig,
		modelOpts:   option,
		key:         cx.(string),
	}, nil
}

func (ro *sarvamOption) GetKey() string {
	return ro.key
}

func (ro *sarvamOption) GetEncoding() string {
	switch ro.audioConfig.Format {

	case internal_audio.Linear16:
		return "linear16"
	case internal_audio.MuLaw8:
		return "mulaw"
	default:
		return "linear16"
	}
}

func (ro *sarvamOption) GetTextToSpeechRequest(contextId, text string) map[string]interface{} {
	return map[string]interface{}{
		"request_id":      contextId,
		"data":            text,
		"binary_response": true,
		"precision":       ro.GetEncoding(),
		"sample_rate":     ro.audioConfig.GetSampleRate(),
	}

}
