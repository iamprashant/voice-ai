// Rapida – Open Source Voice AI Orchestration Platform
// Copyright (C) 2023-2025 Prashant Srivastav <prashant@rapida.ai>
// Licensed under a modified GPL-2.0. See the LICENSE file for details.
package internal_gemini_callers

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/genai"

	internal_callers "github.com/rapidaai/api/integration-api/internal/caller"
	"github.com/rapidaai/pkg/commons"
	type_enums "github.com/rapidaai/pkg/types/enums"
	"github.com/rapidaai/protos"
	integration_api "github.com/rapidaai/protos"
)

type Google struct {
	logger     commons.Logger
	credential internal_callers.CredentialResolver
}

var (
	API_KEY = "key"
)

func google(logger commons.Logger, credential *integration_api.Credential) Google {
	return Google{
		logger: logger,
		credential: func() map[string]interface{} {
			return credential.GetValue().AsMap()
		},
	}
}
func (goog *Google) GetClient() (*genai.Client, error) {
	// need to replace with current request context
	ctx := context.Background()
	credentials := goog.credential()
	cx, ok := credentials[API_KEY]
	if !ok {
		goog.logger.Errorf("Unable to get client for user")
		return nil, errors.New("unable to resolve the credential")
	}
	return genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: cx.(string),
	})
}
func (goog *Google) UsageMetrics(usages *genai.GenerateContentResponseUsageMetadata) []*protos.Metric {
	metrics := make([]*protos.Metric, 0)
	if usages != nil {
		metrics = append(metrics, &protos.Metric{
			Name:        type_enums.INPUT_TOKEN.String(),
			Value:       fmt.Sprintf("%d", usages.PromptTokenCount),
			Description: "Input tokens (including cached content)",
		})

		if usages.CachedContentTokenCount > 0 {
			metrics = append(metrics, &protos.Metric{
				Name:        "CACHED_CONTENT_TOKEN",
				Value:       fmt.Sprintf("%d", usages.CachedContentTokenCount),
				Description: "Cached content tokens",
			})
		}

		metrics = append(metrics, &protos.Metric{
			Name:        type_enums.OUTPUT_TOKEN.String(),
			Value:       fmt.Sprintf("%d", usages.CandidatesTokenCount),
			Description: "Output tokens",
		})

		if usages.ToolUsePromptTokenCount > 0 {
			metrics = append(metrics, &protos.Metric{
				Name:        "TOOL_USE_PROMPT_TOKEN",
				Value:       fmt.Sprintf("%d", usages.ToolUsePromptTokenCount),
				Description: "Tool-use prompt tokens",
			})
		}

		if usages.ThoughtsTokenCount > 0 {
			metrics = append(metrics, &protos.Metric{
				Name:        "THOUGHTS_TOKEN",
				Value:       fmt.Sprintf("%d", usages.ThoughtsTokenCount),
				Description: "Thoughts tokens for thinking models",
			})
		}

		metrics = append(metrics, &protos.Metric{
			Name:        type_enums.TOTAL_TOKEN.String(),
			Value:       fmt.Sprintf("%d", usages.TotalTokenCount),
			Description: "Total tokens (prompt, response, and tool-use)",
		})
	}
	return metrics
}
