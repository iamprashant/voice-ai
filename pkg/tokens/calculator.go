// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package tokens

import (
	"github.com/rapidaai/pkg/types"
)

type TokenCalculator interface {
	Token(in []*types.Message, out *types.Message) []*types.Metric
}
