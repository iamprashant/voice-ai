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
	"sync"

	internal_tool "github.com/rapidaai/api/assistant-api/internal/agent/executor/tool/internal"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

// MCPToolCaller implements the ToolCaller interface for MCP server tools.
// It acts as a dynamic proxy that forwards tool calls to the connected MCP server.
type MCPToolCaller struct {
	logger           commons.Logger
	client           *Client
	toolId           uint64
	toolName         string // The name exposed to the LLM (may be prefixed)
	originalToolName string // The original name on the MCP server
	toolDescription  string
	toolDefinition   *protos.FunctionDefinition
	serverURL        string
}

// NewMCPToolCaller creates a new MCP tool caller for a specific tool
func NewMCPToolCaller(
	logger commons.Logger,
	client *Client,
	toolId uint64,
	toolName string,
	toolDefinition *protos.FunctionDefinition,
) internal_tool.ToolCaller {
	return &MCPToolCaller{
		logger:           logger,
		client:           client,
		toolId:           toolId,
		toolName:         toolName,
		originalToolName: toolName, // Same by default
		toolDescription:  toolDefinition.Description,
		toolDefinition:   toolDefinition,
		serverURL:        client.GetServerURL(),
	}
}

// NewMCPToolCallerWithOriginalName creates a new MCP tool caller with a different exposed name
func NewMCPToolCallerWithOriginalName(
	logger commons.Logger,
	client *Client,
	toolId uint64,
	exposedName string,
	originalName string,
	toolDefinition *protos.FunctionDefinition,
) internal_tool.ToolCaller {
	return &MCPToolCaller{
		logger:           logger,
		client:           client,
		toolId:           toolId,
		toolName:         exposedName,
		originalToolName: originalName,
		toolDescription:  toolDefinition.Description,
		toolDefinition:   toolDefinition,
		serverURL:        client.GetServerURL(),
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
	m.logger.Debugf("MCP tool call: %s (original: %s) with args: %s", m.toolName, m.originalToolName, args)

	// Parse the arguments from JSON string to map
	arguments, err := m.parseArguments(args)
	if err != nil {
		m.logger.Errorf("failed to parse arguments for MCP tool %s: %v", m.toolName, err)
		return m.errorPacket(pkt.ContextId(), fmt.Sprintf("failed to parse arguments: %v", err))
	}

	// Execute the tool call via MCP client using the original tool name
	response, err := m.client.Execute(ctx, m.serverURL, m.originalToolName, arguments)
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

// MCPServerConfig holds configuration for an MCP server connection
type MCPServerConfig struct {
	ServerURL  string
	ToolIdBase uint64
	Prefix     string // Optional prefix to avoid tool name conflicts
}

// MCPConnectionResult holds the result of connecting to an MCP server
type MCPConnectionResult struct {
	ServerURL   string
	Definitions []*protos.FunctionDefinition
	Callers     []internal_tool.ToolCaller
	Error       error
}

// MCPToolManager manages multiple MCP server connections and their tools
type MCPToolManager struct {
	logger  commons.Logger
	clients map[string]*Client // serverURL -> client
	mu      sync.RWMutex
}

// NewMCPToolManager creates a new MCP tool manager
func NewMCPToolManager(logger commons.Logger) *MCPToolManager {
	return &MCPToolManager{
		logger:  logger,
		clients: make(map[string]*Client),
	}
}

// ConnectServersAsync connects to multiple MCP servers concurrently and returns results via channel
// This is non-blocking and allows parallel initialization of multiple MCP servers
func (m *MCPToolManager) ConnectServersAsync(ctx context.Context, configs []MCPServerConfig) <-chan MCPConnectionResult {
	results := make(chan MCPConnectionResult, len(configs))

	var wg sync.WaitGroup
	for _, cfg := range configs {
		wg.Add(1)
		go func(config MCPServerConfig) {
			defer wg.Done()

			definitions, callers, err := m.ConnectServer(ctx, config.ServerURL, config.ToolIdBase, config.Prefix)
			results <- MCPConnectionResult{
				ServerURL:   config.ServerURL,
				Definitions: definitions,
				Callers:     callers,
				Error:       err,
			}
		}(cfg)
	}

	// Close results channel when all connections complete
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// ConnectServer connects to an MCP server and returns all available tools as callers
func (m *MCPToolManager) ConnectServer(
	ctx context.Context,
	serverURL string,
	toolIdBase uint64,
	prefix string,
) ([]*protos.FunctionDefinition, []internal_tool.ToolCaller, error) {
	// Check if already connected
	m.mu.RLock()
	client, exists := m.clients[serverURL]
	m.mu.RUnlock()

	if exists {
		return m.getToolsFromClient(ctx, client, toolIdBase, prefix)
	}

	// Create new connection
	client, err := NewClientWithURL(ctx, m.logger, serverURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to MCP server %s: %w", serverURL, err)
	}

	// Store the client
	m.mu.Lock()
	m.clients[serverURL] = client
	m.mu.Unlock()

	return m.getToolsFromClient(ctx, client, toolIdBase, prefix)
}

// getToolsFromClient fetches tools from a connected client and creates callers
func (m *MCPToolManager) getToolsFromClient(
	ctx context.Context,
	client *Client,
	toolIdBase uint64,
	prefix string,
) ([]*protos.FunctionDefinition, []internal_tool.ToolCaller, error) {
	// List available tools from the MCP server
	definitions, err := client.ListTools(ctx, client.GetServerURL())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list tools: %w", err)
	}

	// Create tool callers for each tool
	callers := make([]internal_tool.ToolCaller, 0, len(definitions))
	resultDefs := make([]*protos.FunctionDefinition, 0, len(definitions))

	for i, def := range definitions {
		toolName := def.Name
		originalName := def.Name

		// Apply prefix if configured (to avoid naming conflicts between servers)
		if prefix != "" {
			toolName = prefix + "_" + def.Name
			// Update the definition name for LLM
			def = &protos.FunctionDefinition{
				Name:        toolName,
				Description: def.Description,
				Parameters:  def.Parameters,
			}
		}

		caller := NewMCPToolCallerWithOriginalName(
			m.logger,
			client,
			toolIdBase+uint64(i),
			toolName,
			originalName,
			def,
		)
		callers = append(callers, caller)
		resultDefs = append(resultDefs, def)
		m.logger.Debugf("Created MCP tool caller: %s (original: %s) from server %s", toolName, originalName, client.GetServerURL())
	}

	return resultDefs, callers, nil
}

// GetClient returns the MCP client for a server URL
func (m *MCPToolManager) GetClient(serverURL string) (*Client, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	client, exists := m.clients[serverURL]
	return client, exists
}

// Close closes all MCP connections
func (m *MCPToolManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	var lastErr error
	for url, client := range m.clients {
		if err := client.Close(); err != nil {
			m.logger.Errorf("failed to close MCP client for %s: %v", url, err)
			lastErr = err
		}
	}
	m.clients = make(map[string]*Client)
	return lastErr
}
