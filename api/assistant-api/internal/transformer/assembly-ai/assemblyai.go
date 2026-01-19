// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_transformer_assemblyai

import (
	"fmt"
	"net/url"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

func (opts *assemblyaiOption) GetEncoding() string {
	switch opts.audioConfig.GetAudioFormat() {
	case protos.AudioConfig_LINEAR16:
		return "pcm_s16le"
	case protos.AudioConfig_MuLaw8:
		return "pcm_mulaw"
	default:
		return "pcm_s16le"
	}
}

type assemblyaiOption struct {
	logger      commons.Logger
	key         string
	mdlOpts     utils.Option
	audioConfig *protos.AudioConfig
}

func NewAssemblyaiOption(
	logger commons.Logger,
	vaultCredential *protos.VaultCredential,
	audioConfig *protos.AudioConfig,
	mdlOpts utils.Option) (*assemblyaiOption, error) {
	cx, ok := vaultCredential.GetValue().AsMap()["key"]
	if !ok {
		return nil, fmt.Errorf("illegal vault config")
	}
	return &assemblyaiOption{
		logger:      logger,
		mdlOpts:     mdlOpts,
		audioConfig: audioConfig,
		key:         cx.(string),
	}, nil
}

func (co *assemblyaiOption) GetKey() string {
	return co.key
}

func (co *assemblyaiOption) GetSpeechToTextConnectionString() string {
	baseURL := "wss://streaming.assemblyai.com/v3/ws"
	params := url.Values{}
	params.Add("sample_rate", fmt.Sprintf("%d", co.audioConfig.SampleRate))
	params.Add("encoding", co.GetEncoding())
	params.Add("format_turns", "true")
	// Check and add language
	if language, err := co.mdlOpts.
		GetString("listen.language"); err == nil {
		params.Add("language", language)
	}

	// Check and add model
	if model, err := co.mdlOpts.
		GetString("listen.model"); err == nil {
		params.Add("model", model)
	}

	return fmt.Sprintf("%s?%s", baseURL, params.Encode())
}
