// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_tool

import (
	"context"
	"encoding/json"

	internal_type "github.com/rapidaai/api/assistant-api/internal/type"

	"github.com/rapidaai/protos"
)

type ToolCallResult map[string]interface{}

func Result(msg string, success bool) ToolCallResult {
	if success {
		return map[string]interface{}{"data": msg, "status": "SUCCESS"}
	} else {
		return map[string]interface{}{"error": msg, "status": "FAIL"}
	}
}

func JustResult(data map[string]interface{}) ToolCallResult {
	return ToolCallResult(data)
}

func (rt ToolCallResult) Result() string {
	bytes, err := json.Marshal(rt)
	if err != nil {
		return `{"error":"failed to marshal result","success":false,"status":"FAIL"}`
	}

	return string(bytes)
}

// ToolCaller defines the contract for invoking a tool/function that can be
// executed by the agent runtime. Implementations encapsulate tool metadata,
// execution semantics, and request/response handling.
//
// A ToolCaller is responsible for:
//   - Exposing a unique identifier and human-readable name
//   - Providing a function definition consumable by the LLM/runtime
//   - Declaring the execution method (e.g., sync, async, streaming)
//   - Executing the tool call and returning response packets
type ToolCaller interface {
	// Id returns the unique identifier of the tool.
	Id() uint64

	// Name returns the human-readable name of the tool.
	Name() string

	// Definition returns the function definition describing the tool's
	// input parameters and behavior, or an error if the definition
	// cannot be constructed.
	Definition() (*protos.FunctionDefinition, error)

	// ExecutionMethod returns the execution strategy used by the tool
	// (for example, synchronous or asynchronous execution).
	ExecutionMethod() string

	// Call executes the tool with the given arguments and communication
	// context. It returns a slice of Packets representing the tool's
	// response(s) to be consumed by the agent runtime.
	Call(ctx context.Context, messageId string, toolId string, args map[string]interface{}, communication internal_type.Communication) ToolCallResult
}
