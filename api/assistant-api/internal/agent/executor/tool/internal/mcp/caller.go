// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_agent_mcp_tool

import (
	"context"

	internal_tool "github.com/rapidaai/api/assistant-api/internal/agent/executor/tool/internal"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/types"
)

// MCPCaller defines the interface for MCP (Model Context Protocol) server integrations.
// This interface enables dynamic discovery and execution of tools from any MCP server
// (e.g., Zapier, Notion, GitHub, etc.).
//
// An MCPCaller is responsible for:
//   - Connecting to an MCP server
//   - Discovering available tools from the server
//   - Executing tool calls via the MCP protocol
//   - Managing the lifecycle of the MCP connection
type MCPCaller interface {
	// List discovers and returns all available tools from the MCP server.
	// Each tool is wrapped as a ToolCaller that can be executed by the agent runtime.
	List() ([]internal_tool.ToolCaller, error)

	// Call executes a tool via the MCP protocol.
	// Returns the result and metrics for the tool execution.
	Call(ctx context.Context, tool internal_tool.ToolCaller, messageId string, args string, communication internal_type.Communication) (map[string]interface{}, []*types.Metric)

	// Close closes the MCP connection and releases resources.
	Close() error
}

// MCPConfig holds configuration for connecting to an MCP server
type MCPConfig struct {
	// ServerType identifies the MCP server type (e.g., "zapier", "notion", "github")
	ServerType string

	// ServerURL is the endpoint URL for the MCP server
	ServerURL string

	// APIKey is the authentication key for the MCP server
	APIKey string

	// Transport is the protocol transport ("sse", "stdio", "http")
	Transport string

	// Additional server-specific configuration options
	Options map[string]interface{}
}

// NewMCPCaller creates a new MCPCaller based on the server type
func NewMCPCaller(logger commons.Logger, config MCPConfig) (MCPCaller, error) {
	switch config.ServerType {
	case "zapier":
		return newZapierMCPClient(logger, config)
	// Add more MCP server types here:
	// case "notion":
	//     return newNotionMCPClient(logger, config)
	// case "github":
	//     return newGitHubMCPClient(logger, config)
	default:
		// For unknown types, use generic MCP client
		return newGenericMCPClient(logger, config)
	}
}
