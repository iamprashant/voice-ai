// Rapida – Open Source Voice AI Orchestration Platform
// Copyright (C) 2023-2025 Prashant Srivastav <prashant@rapida.ai>
// Licensed under a modified GPL-2.0. See the LICENSE file for details.
package internal_callers

import (
	"github.com/rapidaai/pkg/types"
	protos "github.com/rapidaai/protos"
)

type EmbeddingOptions struct {
	AIOptions
}

func NewEmbeddingOptions(
	requestId uint64,
	irRequest *protos.EmbeddingRequest,
	preHook func(rst map[string]interface{}),
	postHook func(rst map[string]interface{}, metrics types.Metrics),
) *EmbeddingOptions {
	cc := &EmbeddingOptions{
		AIOptions: AIOptions{
			RequestId:      requestId,
			PreHook:        preHook,
			PostHook:       postHook,
			ModelParameter: irRequest.GetModelParameters(),
		},
	}
	return cc
}
