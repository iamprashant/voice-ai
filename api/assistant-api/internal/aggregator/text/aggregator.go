// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_sentence_aggregator

import (
	"context"

	internal_default_aggregator "github.com/rapidaai/api/assistant-api/internal/aggregator/text/internal/default"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
)

type TextAggregatorType string

const (
	TextAggregatorDefault    TextAggregatorType = "default"
	OptionsKeyTextAggregator string             = "speaker.sentence.aggregator"
)

func GetLLMTextAggregator(
	context context.Context,
	logger commons.Logger,
	options utils.Option,
) (internal_type.LLMTextAggregator, error) {
	typ, _ := options.GetString(OptionsKeyTextAggregator)
	switch TextAggregatorType(typ) {
	default:
		return internal_default_aggregator.NewDefaultLLMTextAggregator(context, logger, options)
	}
}
