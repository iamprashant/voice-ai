// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_tool_mcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

// ToolResponse represents the response from an MCP tool execution
type ToolResponse struct {
	Success bool           `json:"success"`
	Result  any            `json:"result,omitempty"`
	Error   string         `json:"error,omitempty"`
	Data    map[string]any `json:"data,omitempty"`
	Meta    map[string]any `json:"meta,omitempty"`
}

// NewToolResponse creates a new ToolResponse with the given success status
func NewToolResponse(success bool) *ToolResponse {
	return &ToolResponse{
		Success: success,
		Data:    make(map[string]any),
		Meta:    make(map[string]any),
	}
}

// WithResult sets the result of the tool response
func (r *ToolResponse) WithResult(result any) *ToolResponse {
	r.Result = result
	return r
}

// WithError sets the error message of the tool response
func (r *ToolResponse) WithError(err string) *ToolResponse {
	r.Error = err
	r.Success = false
	return r
}

// WithData adds additional data to the tool response
func (r *ToolResponse) WithData(key string, value any) *ToolResponse {
	r.Data[key] = value
	return r
}

// WithMeta adds metadata to the tool response
func (r *ToolResponse) WithMeta(key string, value any) *ToolResponse {
	r.Meta[key] = value
	return r
}

// ToMap converts the ToolResponse to a map for use in LLMToolPacket
func (r *ToolResponse) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"success": r.Success,
	}
	if r.Result != nil {
		result["result"] = r.Result
	}
	if r.Error != "" {
		result["error"] = r.Error
	}
	if len(r.Data) > 0 {
		result["data"] = r.Data
	}
	if len(r.Meta) > 0 {
		result["meta"] = r.Meta
	}
	if r.Success {
		result["status"] = "SUCCESS"
	} else {
		result["status"] = "FAIL"
	}
	return result
}

// Client manages MCP server connections and tool execution
type Client struct {
	logger    commons.Logger
	session   *mcp.ClientSession
	serverURL string
	tools     map[string]*protos.FunctionDefinition // toolName -> definition
	mu        sync.RWMutex
}

// NewClient creates a new MCP client
func NewClient(ctx context.Context, logger commons.Logger, clientOption utils.Option) (*Client, error) {
	serverUrl, err := clientOption.GetString("mcp.server_url")
	if err != nil {
		return nil, fmt.Errorf("mcp.server_url is required: %w", err)
	}

	client := mcp.NewClient(&mcp.Implementation{
		Name:    "rapida-voice-ai",
		Version: "1.0.0",
	}, nil)

	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: serverUrl}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MCP server: %w", err)
	}
	return &Client{
		logger:    logger,
		session:   session,
		serverURL: serverUrl,
		tools:     make(map[string]*protos.FunctionDefinition),
	}, nil
}

// NewClientWithURL creates a new MCP client with direct server URL
func NewClientWithURL(ctx context.Context, logger commons.Logger, serverURL string) (*Client, error) {
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "rapida-voice-ai",
		Version: "1.0.0",
	}, nil)

	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: serverURL}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MCP server at %s: %w", serverURL, err)
	}
	return &Client{
		logger:    logger,
		session:   session,
		serverURL: serverURL,
		tools:     make(map[string]*protos.FunctionDefinition),
	}, nil
}

// GetServerURL returns the server URL for this client
func (c *Client) GetServerURL() string {
	return c.serverURL
}

// GetToolDefinition returns the cached tool definition by name
func (c *Client) GetToolDefinition(toolName string) (*protos.FunctionDefinition, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	def, ok := c.tools[toolName]
	return def, ok
}

// ListTools returns all tool definitions from the MCP server and caches them
func (c *Client) ListTools(ctx context.Context, serverURL string) ([]*protos.FunctionDefinition, error) {
	result, err := c.session.ListTools(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	definitions := make([]*protos.FunctionDefinition, 0, len(result.Tools))
	for _, tool := range result.Tools {
		c.logger.Debugf("got all tools %+v", *tool)
		def := convertTool(tool)
		definitions = append(definitions, def)
		// Cache the tool definition for later lookup
		c.tools[def.Name] = def
	}

	c.logger.Infof("Found %d tools on MCP server %s", len(definitions), c.serverURL)
	return definitions, nil
}

// ListToolNames returns all tool names available on the MCP server
func (c *Client) ListToolNames(ctx context.Context) ([]string, error) {
	result, err := c.session.ListTools(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list tools: %w", err)
	}

	names := make([]string, 0, len(result.Tools))
	for _, tool := range result.Tools {
		names = append(names, tool.Name)
	}
	return names, nil
}

// Execute calls an MCP tool and returns the response
func (c *Client) Execute(ctx context.Context, serverURL, toolName string, args map[string]any) (*ToolResponse, error) {
	result, err := c.session.CallTool(ctx, &mcp.CallToolParams{
		Name:      toolName,
		Arguments: args,
	})
	if err != nil {
		return NewToolResponse(false).WithError(err.Error()), nil
	}

	return convertResult(result), nil
}

// Close closes all active sessions
func (c *Client) Close() error {
	c.session.Close()
	return nil
}

// convertTool converts MCP Tool to FunctionDefinition
func convertTool(tool *mcp.Tool) *protos.FunctionDefinition {
	params := convertSchema(tool.InputSchema)

	// Ensure properties map is never nil - protobuf omitempty will skip empty maps,
	// but we need {"type": "object", "properties": {}} for valid JSON schema.
	if len(params.Properties) == 0 {
		// Initialize with a non-nil map to ensure proper serialization
		params.Properties = make(map[string]*protos.FunctionParameterProperty)
	}

	return &protos.FunctionDefinition{
		Name:        tool.Name,
		Description: tool.Description,
		Parameters:  params,
	}
}

// convertSchema converts MCP inputSchema to FunctionParameter
// Ensures valid JSON schema format even when no properties exist
func convertSchema(schema any) *protos.FunctionParameter {
	// Always initialize with empty properties map to ensure valid JSON schema
	// Output: {"type": "object", "properties": {}} instead of {"type": "object"}
	params := &protos.FunctionParameter{
		Type:       "object",
		Properties: make(map[string]*protos.FunctionParameterProperty),
		Required:   make([]string, 0),
	}

	if schema == nil {
		return params
	}

	// Convert to map via JSON
	data, err := json.Marshal(schema)
	if err != nil {
		return params
	}

	var schemaMap map[string]any
	if err := json.Unmarshal(data, &schemaMap); err != nil {
		return params
	}

	if t, ok := schemaMap["type"].(string); ok {
		params.Type = t
	}

	if required, ok := schemaMap["required"].([]any); ok {
		for _, r := range required {
			if s, ok := r.(string); ok {
				params.Required = append(params.Required, s)
			}
		}
	}

	if props, ok := schemaMap["properties"].(map[string]any); ok {
		for name, prop := range props {
			propMap, ok := prop.(map[string]any)
			if !ok {
				continue
			}
			p := &protos.FunctionParameterProperty{}
			if t, ok := propMap["type"].(string); ok {
				p.Type = t
			}
			if d, ok := propMap["description"].(string); ok {
				p.Description = d
			}
			if e, ok := propMap["enum"].([]any); ok {
				for _, v := range e {
					if s, ok := v.(string); ok {
						p.Enum = append(p.Enum, s)
					}
				}
			}
			params.Properties[name] = p
		}
	}

	return params
}

// convertResult converts MCP CallToolResult to ToolResponse
func convertResult(result *mcp.CallToolResult) *ToolResponse {
	resp := NewToolResponse(!result.IsError)

	var texts []string
	for _, content := range result.Content {
		switch ct := content.(type) {
		case *mcp.TextContent:
			texts = append(texts, ct.Text)
		case *mcp.ImageContent:
			resp.WithData("image", map[string]string{
				"mimeType": ct.MIMEType,
				"data":     base64.StdEncoding.EncodeToString(ct.Data),
			})
		case *mcp.EmbeddedResource:
			resp.WithData("resource", map[string]any{
				"uri":      ct.Resource.URI,
				"mimeType": ct.Resource.MIMEType,
			})
		}
	}

	if len(texts) == 1 {
		resp.WithResult(texts[0])
	} else if len(texts) > 1 {
		resp.WithResult(texts)
	}

	if result.IsError && len(texts) > 0 {
		resp.WithError(texts[0])
	}

	return resp
}
