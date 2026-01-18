// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_sentence_assembler

import (
	"context"

	internal_default_assembler "github.com/rapidaai/api/assistant-api/internal/assembler/text/internal/default"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
)

type TextAssemblerType string

const (
	TextAssemblerDefault    TextAssemblerType = "default"
	OptionsKeyTextAssembler string            = "speaker.sentence.assembler"
)

func GetLLMTextAssembler(
	context context.Context,
	logger commons.Logger,
	options utils.Option,
) (internal_type.LLMTextAssembler, error) {
	typ, _ := options.GetString(OptionsKeyTextAssembler)
	switch TextAssemblerType(typ) {
	default:
		return internal_default_assembler.NewDefaultLLMTextAssembler(context, logger, options)
	}
}
