// Copyright (c) Rapida
// Author: Prashant <prashant@rapida.ai>
//
// Licensed under the Rapida internal use license.
// This file is part of Rapida's proprietary software and is not open source.
// Unauthorized copying, modification, or redistribution is strictly prohibited.

package internal_transformer_aws

import (
	"context"

	internal_transformer "github.com/rapidaai/api/assistant-api/internal/transformer"
	"github.com/rapidaai/config"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

type awsSpeechToTextTransformer struct {
	cfg    *config.AppConfig
	logger commons.Logger
}

func NewAWSTextToSpeech(ctx context.Context, logger commons.Logger,
	vaultCredential *protos.VaultCredential,
	opts *internal_transformer.TextToSpeechInitializeOptions) (internal_transformer.TextToSpeechTransformer, error) {
	return nil, nil
	// return &awsSpeechToTextTransformer{
	// 	cfg:    cfg,
	// 	logger: logger,
	// }
}

// a implements internal_transformer.SpeechToTextTransformer.
func (*awsSpeechToTextTransformer) Name() string {
	panic("unimplemented")
}

// trf implements internal_transformer.SpeechToTextTransformer.
func (*awsSpeechToTextTransformer) Transform(a []byte) (string, error) {
	panic("unimplemented")
}
