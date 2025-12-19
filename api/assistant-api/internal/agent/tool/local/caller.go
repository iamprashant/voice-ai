// Copyright (c) Rapida
// Author: Prashant <prashant@rapida.ai>
//
// Licensed under the Rapida internal use license.
// This file is part of Rapida's proprietary software.
// Unauthorized copying, modification, or redistribution is strictly prohibited.
package internal_agent_local_tool

import (
	"context"

	internal_adapter_requests "github.com/rapidaai/api/assistant-api/internal/adapters/requests"
	"github.com/rapidaai/pkg/types"
	protos "github.com/rapidaai/protos"
)

type ToolCaller interface {
	// tool call id
	Id() uint64

	//
	Name() string

	//
	Definition() (*protos.FunctionDefinition, error)

	//
	ExecutionMethod() string

	//
	Call(
		ctx context.Context,
		messageId string,
		args string,
		communication internal_adapter_requests.Communication,
	) (map[string]interface{}, []*types.Metric)
}
