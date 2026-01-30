// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_tool_mcp

import (
	"context"
	"fmt"

	internal_tool "github.com/rapidaai/api/assistant-api/internal/agent/executor/tool/internal"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

// MCPToolCaller implements the ToolCaller interface for MCP server tools.
// It forwards tool calls to the connected MCP server.
type MCPToolCaller struct {
	logger         commons.Logger
	client         *Client
	toolId         uint64
	toolName       string
	toolDefinition *protos.FunctionDefinition
}

// NewMCPToolCaller creates a new MCP tool caller for a specific tool
func NewMCPToolCaller(logger commons.Logger, client *Client, toolId uint64, toolName string, toolDefinition *protos.FunctionDefinition,
) internal_tool.ToolCaller {
	return &MCPToolCaller{
		logger:         logger,
		client:         client,
		toolId:         toolId,
		toolName:       toolName,
		toolDefinition: toolDefinition,
	}
}

// Id returns the unique identifier of the tool
func (m *MCPToolCaller) Id() uint64 {
	return m.toolId
}

// Name returns the human-readable name of the tool
func (m *MCPToolCaller) Name() string {
	return m.toolName
}

// Definition returns the function definition describing the tool's input parameters
func (m *MCPToolCaller) Definition() (*protos.FunctionDefinition, error) {
	if m.toolDefinition == nil {
		return nil, fmt.Errorf("tool definition not available for %s", m.toolName)
	}
	return m.toolDefinition, nil
}

// ExecutionMethod returns the execution strategy used by the tool
func (m *MCPToolCaller) ExecutionMethod() string {
	return "mcp"
}

// Call executes the MCP tool with the given arguments and returns the response
func (m *MCPToolCaller) Call(
	ctx context.Context,
	contextID,
	toolId string,
	args map[string]interface{},
	communication internal_type.Communication,
) internal_tool.ToolCallResult {
	response, err := m.client.Execute(ctx, m.toolName, args)
	if err != nil {
		m.logger.Errorf("MCP tool execution failed for %s: %v", m.toolName, err)
		return m.errorPacket(contextID, fmt.Sprintf("tool execution failed: %v", err))
	}
	if response.Error != "" {
		return m.errorPacket(contextID, response.Error)
	}
	return internal_tool.JustResult(map[string]interface{}{
		"status": "SUCCESS",
		"result": response.Result,
	})
}

// errorPacket creates an error response packet
func (m *MCPToolCaller) errorPacket(contextId, errorMsg string) internal_tool.ToolCallResult {
	return internal_tool.JustResult(map[string]interface{}{
		"status": "FAIL",
		"error":  errorMsg,
	})
}
