// Copyright (c) Rapida
// Author: Prashant <prashant@rapida.ai>
//
// Licensed under the Rapida internal use license.
// This file is part of Rapida's proprietary software.
// Unauthorized copying, modification, or redistribution is strictly prohibited.
package internal_adapter_requests

import (
	"context"

	"github.com/rapidaai/pkg/types"
	protos "github.com/rapidaai/protos"
)

type KnowledgeRetriveOption struct {
	EmbeddingProviderCredential *protos.VaultCredential
	RetrievalMethod             string
	TopK                        uint32
	ScoreThreshold              float32
}

type KnowledgeContextResult struct {
	ID         string                 `json:"id"`
	DocumentID string                 `json:"document_id"`
	Metadata   map[string]interface{} `json:"metadata"`
	Content    string                 `json:"content"`
	Score      float64                `json:"score"`
}

type Talking interface {
	Communication
	Connect(
		ctx context.Context,
		auth types.SimplePrinciple,
		identifier string,
		connectionConfig *protos.AssistantConversationConfiguration) error
	Talk(
		ctx context.Context,
		auth types.SimplePrinciple,
		identifier string,
	) error
}
