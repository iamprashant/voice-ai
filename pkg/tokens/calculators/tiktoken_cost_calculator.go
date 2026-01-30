// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package token_tiktoken_calculators

import (
	"fmt"
	"strings"

	"github.com/pkoukk/tiktoken-go"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/tokens"
	type_enums "github.com/rapidaai/pkg/types/enums"
	"github.com/rapidaai/protos"
)

type tikTokenCostCalculator struct {
	logger commons.Logger
	model  string
}

func NewTikTokenCostCalculator(
	logger commons.Logger,
	providerModel string) tokens.TokenCalculator {
	return &tikTokenCostCalculator{
		logger: logger,
		model:  providerModel,
	}
}

func (occ *tikTokenCostCalculator) Token(in []*protos.Message, out *protos.Message) []*protos.Metric {
	mt := make([]*protos.Metric, 0)
	ti, to := occ.token(occ.model, in, out)
	mt = append(mt, &protos.Metric{
		Name:        type_enums.INPUT_TOKEN.String(),
		Value:       fmt.Sprintf("%d", ti),
		Description: "Number of input tokens",
	})

	mt = append(mt, &protos.Metric{
		Name:        type_enums.OUTPUT_TOKEN.String(),
		Value:       fmt.Sprintf("%d", to),
		Description: "Number of output tokens",
	})
	// If you want to add the total token count as well
	totalTokens := ti + to
	mt = append(mt, &protos.Metric{
		Name:        type_enums.TOTAL_TOKEN.String(),
		Value:       fmt.Sprintf("%d", totalTokens),
		Description: "Total number of tokens",
	})
	return mt
}

func (occ *tikTokenCostCalculator) token(name string,
	in []*protos.Message, out *protos.Message) (int, int) {
	tkm, err := tiktoken.EncodingForModel(name)
	if err != nil {
		return 0, 0
	}

	var tokensPerMessage = 0
	switch name {
	case "gpt-3.5-turbo-0613",
		"gpt-3.5-turbo-16k-0613",
		"gpt-4-0314",
		"gpt-4-32k-0314",
		"gpt-4-0613",
		"gpt-4-32k-0613":
		tokensPerMessage = 3
	case "gpt-3.5-turbo-0301":
		tokensPerMessage = 4 // every message follows <|start|>{role/name}\n{content}<|end|>\n
	default:
		if strings.Contains(name, "gpt-3.5-turbo") {
			return occ.token("gpt-3.5-turbo-0613", in, out)
		} else if strings.Contains(name, "gpt-4") {
			return occ.token("gpt-4-0613", in, out)
		} else {
			return 0, 0
		}
	}
	inTokenCount := 0
	for _, message := range in {
		switch msg := message.Message.(type) {
		case *protos.Message_User:
			inTokenCount += tokensPerMessage
			inTokenCount += len(tkm.Encode(msg.User.GetContent(), nil, nil))
			inTokenCount += len(tkm.Encode(message.GetRole(), nil, nil))
		case *protos.Message_Assistant:
			inTokenCount += tokensPerMessage
			inTokenCount += len(tkm.Encode(strings.Join(out.GetAssistant().GetContents(), ""), nil, nil))
			inTokenCount += len(tkm.Encode(message.GetRole(), nil, nil))
		case *protos.Message_System:
			inTokenCount += tokensPerMessage
			inTokenCount += len(tkm.Encode(msg.System.GetContent(), nil, nil))
			inTokenCount += len(tkm.Encode(message.GetRole(), nil, nil))
		}
	}
	// every reply is primed with <|start|>assistant<|message|>
	inTokenCount += 3
	outputToken := 0
	outputToken += len(tkm.Encode(strings.Join(out.GetAssistant().GetContents(), ""), nil, nil))
	outputToken += len(tkm.Encode("assistant", nil, nil))

	return inTokenCount, outputToken
}
