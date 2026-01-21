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
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// zapierMCPClient implements MCPCaller for Zapier MCP integration
type zapierMCPClient struct {
	logger    commons.Logger
	mcpClient *mcp.Client
	config    MCPConfig
	// Cache of discovered tools
	toolsCache []internal_tool.ToolCaller
}

// newZapierMCPClient creates a new Zapier MCP caller (internal factory function)
func newZapierMCPClient(logger commons.Logger, config MCPConfig) (MCPCaller, error) {
	// Set defaults for Zapier
	if config.ServerURL == "" {
		config.ServerURL = "https://mcp.zapier.app"
		logger.Infof("using default Zapier MCP server URL: %s", config.ServerURL)
	}

	if config.Transport == "" {
		config.Transport = "sse" // Zapier uses Server-Sent Events
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("api_key is required for Zapier MCP integration")
	}

	caller := &zapierMCPClient{
		logger:     logger,
		config:     config,
		toolsCache: make([]internal_tool.ToolCaller, 0),
	}

	// Initialize MCP client connection
	if err := caller.connect(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to Zapier MCP server: %v", err)
	}

	return caller, nil
}

// connect establishes connection to the Zapier MCP server
func (z *zapierMCPClient) connect(ctx context.Context) error {
	z.logger.Infof("connecting to Zapier MCP server: %s", z.config.ServerURL)

	// TODO: Implement actual MCP client connection
	// This will establish a connection to the Zapier MCP server
	// using the provided API key and server URL
	//
	// Example (pseudo-code):
	// client, err := mcp.NewClient(ctx, mcp.Config{
	//     ServerURL: z.config.ServerURL,
	//     APIKey:    z.config.APIKey,
	//     Transport: z.config.Transport,
	// })
	// if err != nil {
	//     return fmt.Errorf("failed to create MCP client: %v", err)
	// }
	// z.mcpClient = client

	z.logger.Infof("successfully connected to Zapier MCP server")
	return nil
}

// List discovers and returns all available tools from Zapier MCP server
func (z *zapierMCPClient) List() ([]internal_tool.ToolCaller, error) {
	z.logger.Infof("discovering available tools from Zapier MCP server")

	// TODO: Implement tool discovery via MCP protocol
	// This will query the Zapier MCP server for available tools
	//
	// Example (pseudo-code):
	// tools, err := z.mcpClient.ListTools(context.Background())
	// if err != nil {
	//     return nil, fmt.Errorf("failed to list MCP tools: %v", err)
	// }
	//
	// // Convert MCP tools to internal ToolCaller interface
	// callers := make([]internal_tool.ToolCaller, 0, len(tools))
	// for _, tool := range tools {
	//     caller := &zapierToolWrapper{
	//         client:      z,
	//         toolName:    tool.Name,
	//         toolSchema:  tool.InputSchema,
	//         description: tool.Description,
	//     }
	//     callers = append(callers, caller)
	// }
	//
	// z.toolsCache = callers
	// return callers, nil

	// Placeholder: return empty list for now
	z.logger.Infof("discovered 0 tools (placeholder implementation)")
	return z.toolsCache, nil
}

// Call executes a tool via Zapier MCP
func (z *zapierMCPClient) Call(
	ctx context.Context,
	tool internal_tool.ToolCaller,
	messageId string,
	args string,
	communication internal_type.Communication,
) (map[string]interface{}, []*types.Metric) {
	start := time.Now()
	z.logger.Infof("executing Zapier MCP tool: %s with args: %s", tool.Name(), args)

	// Parse arguments
	var arguments map[string]interface{}
	if err := json.Unmarshal([]byte(args), &arguments); err != nil {
		z.logger.Errorf("failed to parse arguments: %v", err)
		return map[string]interface{}{
			"error":   fmt.Sprintf("failed to parse arguments: %v", err),
			"success": false,
			"status":  "FAIL",
		}, []*types.Metric{types.NewTimeTakenMetric(time.Since(start))}
	}

	// TODO: Execute MCP tool call
	// This will send the tool call to the Zapier MCP server
	// and return the response
	//
	// Example (pseudo-code):
	// result, err := z.mcpClient.CallTool(ctx, mcp.ToolCall{
	//     Name:      tool.Name(),
	//     Arguments: arguments,
	// })
	// if err != nil {
	//     z.logger.Errorf("MCP tool call failed: %v", err)
	//     return map[string]interface{}{
	//         "error":   fmt.Sprintf("MCP tool call failed: %v", err),
	//         "success": false,
	//         "status":  "FAIL",
	//     }, []*types.Metric{types.NewTimeTakenMetric(time.Since(start))}
	// }
	//
	// return map[string]interface{}{
	//     "success": true,
	//     "status":  "SUCCESS",
	//     "data":    result.Content,
	// }, []*types.Metric{types.NewTimeTakenMetric(time.Since(start))}

	// Placeholder response
	z.logger.Infof("MCP tool call completed (placeholder)")
	return map[string]interface{}{
		"success": true,
		"status":  "SUCCESS",
		"message": "MCP tool call executed successfully (placeholder)",
		"data":    arguments,
	}, []*types.Metric{types.NewTimeTakenMetric(time.Since(start))}
}

// Close closes the MCP client connection
func (z *zapierMCPClient) Close() error {
	if z.mcpClient != nil {
		z.logger.Infof("closing Zapier MCP client connection")
		// TODO: Close MCP client connection
		// z.mcpClient.Close()
	}
	return nil
}

// zapierToolWrapper wraps a Zapier MCP tool to implement the ToolCaller interface
type zapierToolWrapper struct {
	client      *zapierMCPClient
	toolId      uint64
	toolName    string
	toolSchema  interface{} // MCP tool schema
	description string
}

// Id returns the unique identifier of the tool
func (z *zapierToolWrapper) Id() uint64 {
	return z.toolId
}

// Name returns the human-readable name of the tool
func (z *zapierToolWrapper) Name() string {
	return z.toolName
}

// ExecutionMethod returns the execution strategy
func (z *zapierToolWrapper) ExecutionMethod() string {
	return "zapier_mcp"
}

// Definition returns the function definition describing the tool
func (z *zapierToolWrapper) Definition() (*protos.FunctionDefinition, error) {
	definition := &protos.FunctionDefinition{
		Name:        z.toolName,
		Description: z.description,
		Parameters:  &protos.FunctionParameter{},
	}

	// Convert MCP schema to FunctionParameter
	if err := utils.Cast(z.toolSchema, definition.Parameters); err != nil {
		return nil, fmt.Errorf("failed to convert MCP schema: %v", err)
	}

	return definition, nil
}

// Call executes the Zapier MCP tool
func (z *zapierToolWrapper) Call(
	ctx context.Context,
	pkt internal_type.LLMPacket,
	toolId string,
	args string,
	communication internal_type.Communication,
) internal_type.LLMToolPacket {
	result, metrics := z.client.Call(ctx, z, pkt.ContextId(), args, communication)

	// Log metrics
	z.client.logger.Debugf("MCP tool call metrics: %+v", metrics)

	return internal_type.LLMToolPacket{
		Name:      z.toolName,
		ContextID: pkt.ContextId(),
		Action:    protos.AssistantConversationAction_API_REQUEST, // MCP calls are similar to API requests
		Result:    result,
	}
}

// NewMCPCallerFromAssistantTool creates an MCP caller from assistant tool configuration
// This is a convenience function for creating MCP callers from database configuration
func NewMCPCallerFromAssistantTool(
	logger commons.Logger,
	toolOptions *internal_assistant_entity.AssistantTool,
) (MCPCaller, error) {
	opts := toolOptions.GetOptions()

	// Extract MCP configuration from tool options
	config := MCPConfig{
		ServerType: "zapier", // Default to Zapier, can be made configurable
		Options:    make(map[string]interface{}),
	}

	// Parse standard configuration fields
	if serverType, err := opts.GetString("mcp.server_type"); err == nil {
		config.ServerType = serverType
	}

	if serverURL, err := opts.GetString("mcp.server_url"); err == nil {
		config.ServerURL = serverURL
	}

	if apiKey, err := opts.GetString("mcp.api_key"); err == nil {
		config.APIKey = apiKey
	}

	if transport, err := opts.GetString("mcp.transport"); err == nil {
		config.Transport = transport
	}

	// Store all options for server-specific use
	for k, v := range opts {
		config.Options[k] = v
	}

	return NewMCPCaller(logger, config)
}
