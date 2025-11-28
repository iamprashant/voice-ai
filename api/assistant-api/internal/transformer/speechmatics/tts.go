// Copyright (c) Rapida
// Author: Prashant <prashant@rapida.ai>
//
// Licensed under the Rapida internal use license.
// This file is part of Rapida's proprietary software and is not open source.
// Unauthorized copying, modification, or redistribution is strictly prohibited.

package internal_transformer_speechmatics

import (
	"context"

	internal_transformer "github.com/rapidaai/api/assistant-api/internal/transformer"
	"github.com/rapidaai/pkg/commons"
	protos "github.com/rapidaai/protos"
)

func NewSpeechmaticsTextToSpeech(
	ctx context.Context,
	logger commons.Logger,
	credential *protos.VaultCredential,
	onSpeech func([]byte) error) (internal_transformer.TextToSpeechTransformer, error) {
	return nil, nil
}
