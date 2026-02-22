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
func (executor *toolExecutor) initializeLocalTool(ctx context.Context, logger commons.Logger, toolOpts *internal_assistant_entity.AssistantTool, communication internal_type.Communication) (internal_tool.ToolCaller, error) {
	switch toolOpts.ExecutionMethod {
	case "knowledge_retrieval":
		return internal_tool_local.NewKnowledgeRetrievalToolCaller(ctx, logger, toolOpts, communication)
	case "api_request":
		return internal_tool_local.NewApiRequestToolCaller(ctx, logger, toolOpts, communication)
	case "endpoint_request":
		return internal_tool_local.NewEndpointToolCaller(ctx, logger, toolOpts, communication)
	case "end_of_conversation":
		return internal_tool_local.NewEndOfConversationCaller(ctx, logger, toolOpts, communication)
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
			caller, err := executor.initializeLocalTool(ctx, executor.logger, tool, communication)
			if err != nil {
				executor.logger.Errorf("Failed to initialize local tool %s: %v", tool.Name, err)
				continue
			}

			def, err := caller.Definition()
			if err != nil {
				executor.logger.Errorf("Failed to get definition for tool %s: %v", tool.Name, err)
				continue
			}

			tracer.AddAttributes(ctx, internal_adapter_telemetry.KV{K: caller.Name(), V: internal_adapter_telemetry.StringValue(caller.ExecutionMethod())})
			executor.registerTool(caller, def)
		}

	}
}

// Initialize sets up all tools (local + MCP) for the assistant
func (executor *toolExecutor) Initialize(ctx context.Context, communication internal_type.Communication) error {
	ctx, span, _ := communication.Tracer().StartSpan(ctx, utils.AssistantToolConnectStage)
	defer span.EndSpan(ctx, utils.AssistantToolConnectStage)
	executor.initializeTools(ctx, communication.Assistant().AssistantTools, communication, span)
	return nil
}

func (executor *toolExecutor) GetFunctionDefinitions() []*protos.FunctionDefinition {
	return executor.availableToolFunctions
}

func (executor *toolExecutor) execute(ctx context.Context, contextID string, call *protos.ToolCall, communication internal_type.Communication) *protos.ToolMessage_Tool {
	ctx, span, _ := communication.Tracer().StartSpan(ctx, utils.AssistantToolExecuteStage, internal_adapter_telemetry.MessageKV(contextID))
	defer span.EndSpan(ctx, utils.AssistantToolExecuteStage)

	start := time.Now()
	// metrics := make([]*types.Metric, 0)
	funC, ok := executor.getTool(call.GetFunction().GetName())
	if !ok {
		return &protos.ToolMessage_Tool{Name: call.GetFunction().GetName(), Id: call.Id, Content: "unable to find tool: " + call.GetFunction().GetName()}
	}
	span.AddAttributes(ctx,
		internal_adapter_telemetry.KV{K: "function", V: internal_adapter_telemetry.StringValue(call.GetFunction().GetName())},
		internal_adapter_telemetry.KV{K: "argument", V: internal_adapter_telemetry.StringValue(call.GetFunction().GetArguments())})

	arguments := executor.parseArgument(call.GetFunction().GetArguments())
	// on packge
	communication.OnPacket(ctx, internal_type.LLMToolCallPacket{ToolID: call.GetId(), Name: call.GetFunction().GetName(), Arguments: arguments})
	// output
	output := funC.Call(ctx, contextID, call.GetId(), arguments, communication)

	communication.OnPacket(ctx, internal_type.LLMToolResultPacket{ToolID: call.GetId(), Name: call.GetFunction().GetName(), Result: output})

	// log tool execution
	// executor.tool(contextID, map[string]interface{}{
	// 	"id":        call.Id,
	// 	"name":      call.GetFunction().GetName(),
	// 	"arguments": arguments,
	// }, output, append(metrics, types.NewTimeTakenMetric(time.Since(start))), communication)

	// log execution record
	executor.log(ctx, funC, communication, contextID, type_enums.RECORD_COMPLETE, int64(time.Since(start)),
		map[string]interface{}{
			"id":        call.Id,
			"name":      call.GetFunction().GetName(),
			"arguments": arguments,
		}, output)

	return &protos.ToolMessage_Tool{Name: call.GetFunction().GetName(), Id: call.Id, Content: output.Result()}
}

func (executor *toolExecutor) ExecuteAll(ctx context.Context, contextID string, calls []*protos.ToolCall, communication internal_type.Communication) *protos.Message {
	if len(calls) == 0 {
		return nil
	}
	// Use mutex-protected slices for concurrent writes
	result := make([]*protos.ToolMessage_Tool, 0, len(calls))
	var wg sync.WaitGroup
	for _, xt := range calls {
		xtCopy := xt
		wg.Add(1)
		utils.Go(context.Background(), func() {
			defer wg.Done()
			result = append(result, executor.execute(ctx, contextID, xtCopy, communication))
		})
	}
	wg.Wait()
	return &protos.Message{Role: "tool", Message: &protos.Message_Tool{Tool: &protos.ToolMessage{Tools: result}}}
}

func (executor *toolExecutor) parseArgument(arguments string) map[string]interface{} {
	var argMap map[string]interface{}
	err := json.Unmarshal([]byte(arguments), &argMap)
	if err != nil {
		return map[string]interface{}{"raw": arguments}
	} else {
		return argMap
	}
}

// Close releases all resources held by the tool executor
func (executor *toolExecutor) Close(ctx context.Context) error {
	for _, client := range executor.mcpClients {
		if err := client.Close(ctx); err != nil {
			executor.logger.Errorf("failed to close MCP client: %v", err)
		}
	}
	return nil
}

func (executor *toolExecutor) log(ctx context.Context, toolCaller internal_tool.ToolCaller, communication internal_type.Communication, assistantConversationMessageId string, recordStatus type_enums.RecordState, timeTaken int64, in, out map[string]interface{}) {
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
		communication.CreateToolLog(ctx, toolCaller.Id(), assistantConversationMessageId, toolCaller.Name(), toolCaller.ExecutionMethod(), recordStatus, timeTaken, i, o)
	})
}
