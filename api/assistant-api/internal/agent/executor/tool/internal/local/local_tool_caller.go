// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_tool_local

import (
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

type toolCaller struct {
	logger      commons.Logger
	toolOptions *internal_assistant_entity.AssistantTool
}

func (executor *toolCaller) Name() string {
	return executor.toolOptions.Name
}

func (executor *toolCaller) Id() uint64 {
	return executor.toolOptions.Id
}

func (executor *toolCaller) ExecutionMethod() string {
	return executor.toolOptions.ExecutionMethod
}

func (executor *toolCaller) Definition() (*protos.FunctionDefinition, error) {
	definition := &protos.FunctionDefinition{
		Name:       executor.toolOptions.Name,
		Parameters: &protos.FunctionParameter{},
	}
	if executor.toolOptions.Description != nil && *executor.toolOptions.Description != "" {
		definition.Description = *executor.toolOptions.Description
	}
	if err := utils.Cast(executor.toolOptions.Fields, definition.Parameters); err != nil {
		return nil, err
	}
	return definition, nil
}
