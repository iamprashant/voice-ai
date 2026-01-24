// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_agent_executor_tool

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	internal_agent_executor "github.com/rapidaai/api/assistant-api/internal/agent/executor"
	internal_tool "github.com/rapidaai/api/assistant-api/internal/agent/executor/tool/internal"
	internal_tool_local "github.com/rapidaai/api/assistant-api/internal/agent/executor/tool/internal/local"
	internal_tool_mcp "github.com/rapidaai/api/assistant-api/internal/agent/executor/tool/internal/mcp"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_adapter_telemetry "github.com/rapidaai/api/assistant-api/internal/telemetry"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"

	"github.com/rapidaai/protos"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/types"
	type_enums "github.com/rapidaai/pkg/types/enums"
	"github.com/rapidaai/pkg/utils"
)

type toolExecutor struct {
	logger                 commons.Logger
	tools                  map[string]internal_tool.ToolCaller
	availableToolFunctions []*protos.FunctionDefinition
	mcpClients             []*internal_tool_mcp.Client
}

func NewToolExecutor(logger commons.Logger) internal_agent_executor.ToolExecutor {
	return &toolExecutor{
		logger:                 logger,
		mcpClients:             make([]*internal_tool_mcp.Client, 0),
		tools:                  make(map[string]internal_tool.ToolCaller),
		availableToolFunctions: make([]*protos.FunctionDefinition, 0),
	}
}

// registerTool safely registers a tool caller and its definition
func (executor *toolExecutor) registerTool(caller internal_tool.ToolCaller, def *protos.FunctionDefinition) {
	executor.tools[caller.Name()] = caller
	executor.availableToolFunctions = append(executor.availableToolFunctions, def)
}

// getTool safely retrieves a tool caller by name
func (executor *toolExecutor) getTool(name string) (internal_tool.ToolCaller, bool) {
	caller, ok := executor.tools[name]
	return caller, ok
}

// initializeLocalTool creates a tool caller for local execution methods
func (executor *toolExecutor) initializeLocalTool(logger commons.Logger, toolOpts *internal_assistant_entity.AssistantTool, communication internal_type.Communication) (internal_tool.ToolCaller, error) {
	switch toolOpts.ExecutionMethod {
	case "knowledge_retrieval":
		return internal_tool_local.NewKnowledgeRetrievalToolCaller(logger, toolOpts, communication)
	case "api_request":
		return internal_tool_local.NewApiRequestToolCaller(logger, toolOpts, communication)
	case "endpoint":
		return internal_tool_local.NewEndpointToolCaller(logger, toolOpts, communication)
	case "put_on_hold":
		return internal_tool_local.NewPutOnHoldToolCaller(logger, toolOpts, communication)
	case "end_of_conversation":
		return internal_tool_local.NewEndOfConversationCaller(logger, toolOpts, communication)
	default:
		return nil, errors.New("illegal tool action provided")
	}
}

// initializeLocalTools initializes all local tools synchronously
func (executor *toolExecutor) initializeTools(ctx context.Context, tools []*internal_assistant_entity.AssistantTool, communication internal_type.Communication, tracer internal_adapter_telemetry.Tracer[utils.RapidaStage]) {
	for _, tool := range tools {
		switch tool.ExecutionMethod {
		case "mcp":
			client, err := internal_tool_mcp.NewClient(ctx, executor.logger, tool.GetOptions())
			if err != nil {
				continue
			}
			executor.mcpClients = append(executor.mcpClients, client)
			definitions, err := client.ListTools(ctx)
			if err != nil {
				continue
			}
			for i, def := range definitions {
				caller := internal_tool_mcp.NewMCPToolCaller(executor.logger, client, tool.Id+uint64(i), def.Name, def)
				tracer.AddAttributes(ctx, internal_adapter_telemetry.KV{K: caller.Name(), V: internal_adapter_telemetry.StringValue(caller.ExecutionMethod())})
				executor.registerTool(caller, def)
			}
		default:
			caller, err := executor.initializeLocalTool(executor.logger, tool, communication)
			if err != nil {
				executor.logger.Errorf("Failed to initialize local tool %s: %v", tool.Name, err)
				continue
			}

			def, err := caller.Definition()
			if err != nil {
				executor.logger.Errorf("Failed to get definition for tool %s: %v", tool.Name, err)
				continue
			}

			tracer.AddAttributes(communication.Context(), internal_adapter_telemetry.KV{K: caller.Name(), V: internal_adapter_telemetry.StringValue(caller.ExecutionMethod())})
			executor.registerTool(caller, def)
		}

	}
}

// Initialize sets up all tools (local + MCP) for the assistant
func (executor *toolExecutor) Initialize(ctx context.Context, communication internal_type.Communication) error {
	ctx, span, _ := communication.Tracer().StartSpan(ctx, utils.AssistantToolConnectStage)
	defer span.EndSpan(ctx, utils.AssistantToolConnectStage)

	start := time.Now()
	tools := communication.Assistant().AssistantTools

	// Initialize local tools
	executor.initializeTools(ctx, tools, communication, span)
	executor.logger.Benchmark("ToolExecutor.Init", time.Since(start))
	executor.logger.Infof("Initialized %d tools total", len(executor.tools))

	return nil
}

func (executor *toolExecutor) GetFunctionDefinitions() []*protos.FunctionDefinition {
	return executor.availableToolFunctions
}

func (executor *toolExecutor) tool(messageId string, in, out map[string]interface{}, metrics []*types.Metric, communication internal_type.Communication) error {
	utils.Go(communication.Context(), func() {
		communication.CreateConversationToolLog(messageId, in, out, metrics)
	})
	return nil
}

func (executor *toolExecutor) execute(ctx context.Context, message internal_type.LLMPacket, call *protos.ToolCall, communication internal_type.Communication) internal_type.LLMToolPacket {
	ctx, span, _ := communication.Tracer().StartSpan(ctx, utils.AssistantToolExecuteStage, internal_adapter_telemetry.MessageKV(message.ContextId()))
	defer span.EndSpan(ctx, utils.AssistantToolExecuteStage)

	start := time.Now()
	metrics := make([]*types.Metric, 0)

	funC, ok := executor.getTool(call.GetFunction().GetName())
	if !ok {
		return internal_type.LLMToolPacket{ContextID: message.ContextId(),
			Action: protos.AssistantConversationAction_ACTION_UNSPECIFIED, Result: map[string]interface{}{
				"error":   "unable to find tool: " + call.GetFunction().GetName(),
				"success": false,
				"status":  "FAIL",
			}}
	}

	// should return multiple things
	span.AddAttributes(ctx, internal_adapter_telemetry.KV{K: "function", V: internal_adapter_telemetry.StringValue(call.GetFunction().GetName())}, internal_adapter_telemetry.KV{K: "argument", V: internal_adapter_telemetry.StringValue(call.GetFunction().GetArguments())})

	output := funC.Call(ctx, message, call.GetId(), call.GetFunction().GetArguments(), communication)
	metrics = append(metrics, types.NewTimeTakenMetric(time.Since(start)))

	//
	executor.tool(message.ContextId(), map[string]interface{}{
		"id":        call.Id,
		"name":      call.GetFunction().GetName(),
		"arguments": call.GetFunction().GetArguments(),
	}, output.Result, metrics, communication)

	executor.Log(ctx, funC, communication, message.ContextId(), type_enums.RECORD_COMPLETE, int64(time.Since(start)),
		map[string]interface{}{
			"id":        call.Id,
			"name":      call.GetFunction().GetName(),
			"arguments": call.GetFunction().GetArguments(),
		}, output.Result)

	return output
}

func (executor *toolExecutor) ExecuteAll(ctx context.Context, message internal_type.LLMPacket, calls []*protos.ToolCall, communication internal_type.Communication) ([]internal_type.Packet, []*types.Content) {
	if len(calls) == 0 {
		return nil, nil
	}

	// Use mutex-protected slices for concurrent writes
	contents := make([]internal_type.Packet, 0, len(calls))
	result := make([]*types.Content, 0, len(calls))

	var wg sync.WaitGroup
	for _, xt := range calls {
		xtCopy := xt
		wg.Add(1)
		utils.Go(context.Background(), func() {
			defer wg.Done()

			cntn := executor.execute(ctx, message, xtCopy, communication)

			bt, err := json.Marshal(cntn.Result)
			content := &types.Content{
				ContentType:   xtCopy.GetId(),
				ContentFormat: commons.TEXT_CONTENT_FORMAT_RAW.String(),
				Name:          xtCopy.GetFunction().GetName(),
			}
			if err != nil {
				content.Content = []byte("unable to parse the response.")
			} else {
				content.Content = bt
			}

			// Thread-safe append
			contents = append(contents, cntn)
			result = append(result, content)
		})
	}
	wg.Wait()

	return contents, result
}

// Close releases all resources held by the tool executor
func (executor *toolExecutor) Close() error {
	for _, client := range executor.mcpClients {
		if err := client.Close(); err != nil {
			executor.logger.Errorf("failed to close MCP client: %v", err)
		}
	}
	return nil
}

func (executor *toolExecutor) Log(ctx context.Context, toolCaller internal_tool.ToolCaller, communication internal_type.Communication, assistantConversationMessageId string, recordStatus type_enums.RecordState, timeTaken int64, in, out map[string]interface{}) {
	utils.Go(ctx, func() {
		// Include function definition in log
		if def, err := toolCaller.Definition(); err == nil {
			in["llm_definition"] = map[string]interface{}{
				"name":        def.Name,
				"description": def.Description,
				"parameters": map[string]interface{}{
					"type":       def.Parameters.Type,
					"required":   def.Parameters.Required,
					"properties": def.Parameters.Properties,
				},
			}
		}
		i, _ := json.Marshal(in)
		o, _ := json.Marshal(out)
		communication.CreateToolLog(toolCaller.Id(), assistantConversationMessageId, toolCaller.Name(), toolCaller.ExecutionMethod(), recordStatus, timeTaken, i, o)
	})
}
