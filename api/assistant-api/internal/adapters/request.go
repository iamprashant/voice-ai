// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_adapter

import (
	"context"

	"github.com/rapidaai/api/assistant-api/config"

	adapter_internal "github.com/rapidaai/api/assistant-api/internal/adapters/internal"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
	"github.com/rapidaai/pkg/storages"
	"github.com/rapidaai/pkg/utils"
)

func GetTalker(source utils.RapidaSource, ctx context.Context, cfg *config.AssistantConfig, logger commons.Logger, postgres connectors.PostgresConnector, opensearch connectors.OpenSearchConnector, redis connectors.RedisConnector, storage storages.Storage, streamer internal_type.Streamer,
) (internal_type.Talking, error) {
	return adapter_internal.NewGenericRequestor(ctx, cfg, logger, source, postgres, opensearch, redis, storage, streamer), nil
}
