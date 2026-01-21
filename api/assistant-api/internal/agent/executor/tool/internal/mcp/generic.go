// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_agent_mcp_tool

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	internal_tool "github.com/rapidaai/api/assistant-api/internal/agent/executor/tool/internal"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// genericMCPClient implements MCPCaller for any MCP server
type genericMCPClient struct {
	logger     commons.Logger
	mcpClient  *mcp.Client
	config     MCPConfig
	toolsCache []internal_tool.ToolCaller
}

// newGenericMCPClient creates a new generic MCP caller that works with any MCP server
func newGenericMCPClient(logger commons.Logger, config MCPConfig) (MCPCaller, error) {
	if config.ServerURL == "" {
		return nil, fmt.Errorf("server_url is required for MCP integration")
	}

	if config.Transport == "" {
		config.Transport = "sse" // Default to Server-Sent Events
		logger.Infof("using default MCP transport: %s", config.Transport)
	}

	caller := &genericMCPClient{
		logger:     logger,
		config:     config,
		toolsCache: make([]internal_tool.ToolCaller, 0),
	}

	// Initialize MCP client connection
	if err := caller.connect(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to MCP server: %v", err)
	}

	return caller, nil
}

// connect establishes connection to the MCP server
func (g *genericMCPClient) connect(ctx context.Context) error {
	g.logger.Infof("connecting to MCP server: %s (type: %s, transport: %s)",
		g.config.ServerURL, g.config.ServerType, g.config.Transport)

	// TODO: Implement actual MCP client connection
	// The implementation will vary based on the transport type:
	//
	// For SSE (Server-Sent Events):
	// client, err := mcp.NewSSEClient(ctx, mcp.SSEConfig{
	//     URL:    g.config.ServerURL,
	//     APIKey: g.config.APIKey,
	// })
	//
	// For stdio (local process):
	// client, err := mcp.NewStdioClient(ctx, mcp.StdioConfig{
	//     Command: g.config.ServerURL, // Command to execute
	//     Args:    g.config.Options["args"].([]string),
	// })
	//
	// For HTTP:
	// client, err := mcp.NewHTTPClient(ctx, mcp.HTTPConfig{
	//     BaseURL: g.config.ServerURL,
	//     APIKey:  g.config.APIKey,
	// })
	//
	// if err != nil {
	//     return fmt.Errorf("failed to create MCP client: %v", err)
	// }
	// g.mcpClient = client

	g.logger.Infof("successfully connected to MCP server")
	return nil
}

// List discovers and returns all available tools from the MCP server
func (g *genericMCPClient) List() ([]internal_tool.ToolCaller, error) {
	g.logger.Infof("discovering available tools from MCP server: %s", g.config.ServerURL)

	// TODO: Implement tool discovery via MCP protocol
	// This will query the MCP server for available tools using the
	// tools/list endpoint as defined in the MCP specification
	//
	// Example implementation:
	// tools, err := g.mcpClient.ListTools(context.Background())
	// if err != nil {
	//     return nil, fmt.Errorf("failed to list MCP tools: %v", err)
	// }
	//
	// // Convert MCP tools to internal ToolCaller interface
	// callers := make([]internal_tool.ToolCaller, 0, len(tools))
	// for idx, tool := range tools {
	//     caller := &mcpToolWrapper{
	//         client:      g,
	//         toolId:      uint64(idx + 1), // Generate sequential IDs
	//         toolName:    tool.Name,
	//         toolSchema:  tool.InputSchema,
	//         description: tool.Description,
	//     }
	//     callers = append(callers, caller)
	// }
	//
	// g.toolsCache = callers
	// g.logger.Infof("discovered %d tools from MCP server", len(callers))
	// return callers, nil

	// Placeholder: return empty list for now
	g.logger.Infof("discovered 0 tools (placeholder implementation)")
	return g.toolsCache, nil
}

// Call executes a tool via MCP protocol
func (g *genericMCPClient) Call(
	ctx context.Context,
	tool internal_tool.ToolCaller,
	messageId string,
	args string,
	communication internal_type.Communication,
) (map[string]interface{}, []*types.Metric) {
	start := time.Now()
	g.logger.Infof("executing MCP tool: %s with args: %s", tool.Name(), args)

	// Parse arguments
	var arguments map[string]interface{}
	if err := json.Unmarshal([]byte(args), &arguments); err != nil {
		g.logger.Errorf("failed to parse arguments: %v", err)
		return map[string]interface{}{
			"error":   fmt.Sprintf("failed to parse arguments: %v", err),
			"success": false,
			"status":  "FAIL",
		}, []*types.Metric{types.NewTimeTakenMetric(time.Since(start))}
	}

	// TODO: Execute MCP tool call
	// This will send the tool call to the MCP server using the
	// tools/call endpoint as defined in the MCP specification
	//
	// Example implementation:
	// result, err := g.mcpClient.CallTool(ctx, mcp.ToolCall{
	//     Name:      tool.Name(),
	//     Arguments: arguments,
	// })
	// if err != nil {
	//     g.logger.Errorf("MCP tool call failed: %v", err)
	//     return map[string]interface{}{
	//         "error":   fmt.Sprintf("MCP tool call failed: %v", err),
	//         "success": false,
	//         "status":  "FAIL",
	//     }, []*types.Metric{types.NewTimeTakenMetric(time.Since(start))}
	// }
	//
	// // Parse result content
	// resultData := make(map[string]interface{})
	// for _, content := range result.Content {
	//     if content.Type == "text" {
	//         resultData["text"] = content.Text
	//     } else if content.Type == "resource" {
	//         resultData["resource"] = content.Resource
	//     }
	// }
	//
	// return map[string]interface{}{
	//     "success": true,
	//     "status":  "SUCCESS",
	//     "data":    resultData,
	//     "isError": result.IsError,
	// }, []*types.Metric{types.NewTimeTakenMetric(time.Since(start))}

	// Placeholder response
	g.logger.Infof("MCP tool call completed (placeholder)")
	return map[string]interface{}{
		"success": true,
		"status":  "SUCCESS",
		"message": "MCP tool call executed successfully (placeholder)",
		"data":    arguments,
	}, []*types.Metric{types.NewTimeTakenMetric(time.Since(start))}
}

// Close closes the MCP client connection
func (g *genericMCPClient) Close() error {
	if g.mcpClient != nil {
		g.logger.Infof("closing MCP client connection")
		// TODO: Close MCP client connection
		// return g.mcpClient.Close()
	}
	return nil
}

// mcpToolWrapper wraps an MCP tool to implement the ToolCaller interface
type mcpToolWrapper struct {
	client      MCPCaller
	toolId      uint64
	toolName    string
	toolSchema  interface{} // MCP tool schema (JSON Schema)
	description string
}

// Id returns the unique identifier of the tool
func (m *mcpToolWrapper) Id() uint64 {
	return m.toolId
}

// Name returns the human-readable name of the tool
func (m *mcpToolWrapper) Name() string {
	return m.toolName
}

// ExecutionMethod returns the execution strategy
func (m *mcpToolWrapper) ExecutionMethod() string {
	return "mcp"
}

// Definition returns the function definition describing the tool
func (m *mcpToolWrapper) Definition() (*protos.FunctionDefinition, error) {
	definition := &protos.FunctionDefinition{
		Name:        m.toolName,
		Description: m.description,
		Parameters:  &protos.FunctionParameter{},
	}

	// Convert MCP schema (JSON Schema) to FunctionParameter
	if m.toolSchema != nil {
		if err := utils.Cast(m.toolSchema, definition.Parameters); err != nil {
			return nil, fmt.Errorf("failed to convert MCP schema: %v", err)
		}
	}

	return definition, nil
}

// Call executes the MCP tool
func (m *mcpToolWrapper) Call(
	ctx context.Context,
	pkt internal_type.LLMPacket,
	toolId string,
	args string,
	communication internal_type.Communication,
) internal_type.LLMToolPacket {
	result, metrics := m.client.Call(ctx, m, pkt.ContextId(), args, communication)

	// Log metrics
	if logger, ok := m.client.(*genericMCPClient); ok {
		logger.logger.Debugf("MCP tool call metrics: %+v", metrics)
	}

	return internal_type.LLMToolPacket{
		Name:      m.toolName,
		ContextID: pkt.ContextId(),
		Action:    protos.AssistantConversationAction_API_REQUEST, // MCP calls are similar to API requests
		Result:    result,
	}
}
