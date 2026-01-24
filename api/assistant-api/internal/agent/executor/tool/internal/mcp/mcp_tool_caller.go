// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_tool_mcp

import (
	"context"
	"encoding/json"
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
	pkt internal_type.LLMPacket,
	toolId string,
	args string,
	communication internal_type.Communication,
) internal_type.LLMToolPacket {
	m.logger.Debugf("MCP tool call: %s with args: %s", m.toolName, args)

	// Parse the arguments from JSON string to map
	arguments, err := m.parseArguments(args)
	if err != nil {
		m.logger.Errorf("failed to parse arguments for MCP tool %s: %v", m.toolName, err)
		return m.errorPacket(pkt.ContextId(), fmt.Sprintf("failed to parse arguments: %v", err))
	}

	// Execute the tool call via MCP client
	response, err := m.client.Execute(ctx, m.toolName, arguments)
	if err != nil {
		m.logger.Errorf("MCP tool execution failed for %s: %v", m.toolName, err)
		return m.errorPacket(pkt.ContextId(), fmt.Sprintf("tool execution failed: %v", err))
	}

	// Convert response to LLMToolPacket
	return internal_type.LLMToolPacket{
		Name:      m.toolName,
		ContextID: pkt.ContextId(),
		Action:    protos.AssistantConversationAction_MCP_TOOL_CALL,
		Result:    response.ToMap(),
	}
}

// parseArguments converts a JSON string to a map of arguments
func (m *MCPToolCaller) parseArguments(args string) (map[string]any, error) {
	if args == "" || args == "{}" {
		return make(map[string]any), nil
	}

	var arguments map[string]any
	if err := json.Unmarshal([]byte(args), &arguments); err != nil {
		// Try to wrap as a single value if it's not valid JSON object
		return map[string]any{"input": args}, nil
	}
	return arguments, nil
}

// errorPacket creates an error response packet
func (m *MCPToolCaller) errorPacket(contextId, errorMsg string) internal_type.LLMToolPacket {
	return internal_type.LLMToolPacket{
		Name:      m.toolName,
		ContextID: contextId,
		Action:    protos.AssistantConversationAction_MCP_TOOL_CALL,
		Result: map[string]interface{}{
			"success": false,
			"status":  "FAIL",
			"error":   errorMsg,
		},
	}
}
