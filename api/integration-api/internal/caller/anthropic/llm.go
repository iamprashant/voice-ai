package internal_anthropic_callers

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"

	internal_callers "github.com/rapidaai/api/integration-api/internal/caller"
	internal_caller_metrics "github.com/rapidaai/api/integration-api/internal/caller/metrics"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	protos "github.com/rapidaai/protos"
)

type largeLanguageCaller struct {
	Anthropic
}

func NewLargeLanguageCaller(logger commons.Logger, credential *protos.Credential) internal_callers.LargeLanguageCaller {
	return &largeLanguageCaller{
		Anthropic: anthropicAI(logger, credential),
	}
}

func (llc *largeLanguageCaller) BuildHistory(allMessages []*protos.Message) ([]anthropic.TextBlockParam, []anthropic.MessageParam) {
	messages := make([]anthropic.MessageParam, 0)
	systemPrompt := make([]anthropic.TextBlockParam, 0)
	for _, msg := range allMessages {
		switch msg.GetRole() {
		case "assistant":
			mConnect := make([]anthropic.ContentBlockParamUnion, 0)
			for _, c := range msg.GetAssistant().GetContents() {
				mConnect = append(mConnect, anthropic.ContentBlockParamUnion{
					OfText: &anthropic.TextBlockParam{
						Text: string(c),
					},
				})
			}
			for _, tc := range msg.GetAssistant().GetToolCalls() {
				var input map[string]interface{}
				if err := json.Unmarshal([]byte(tc.GetFunction().GetArguments()), &input); err != nil {
					llc.logger.Warnf("Invalid JSON in tool call arguments: %v", err)
					continue
				}
				mConnect = append(mConnect, anthropic.ContentBlockParamUnion{
					OfToolUse: &anthropic.ToolUseBlockParam{
						ID:    tc.GetId(),
						Name:  tc.GetFunction().GetName(),
						Input: input,
					},
				})
			}
			if len(mConnect) > 0 {
				messages = append(messages, anthropic.MessageParam{
					Role:    anthropic.MessageParamRoleAssistant,
					Content: mConnect,
				})
			}
		case "user":
			if u := msg.GetUser(); u != nil && strings.TrimSpace(u.GetContent()) != "" {
				messages = append(messages, anthropic.MessageParam{
					Role: anthropic.MessageParamRoleUser,
					Content: []anthropic.ContentBlockParamUnion{{
						OfText: &anthropic.TextBlockParam{
							Text: string(u.GetContent()),
						},
					},
					}})
			}
		case "tool":
			tContent := make([]anthropic.ContentBlockParamUnion, 0)
			for _, c := range msg.GetTool().GetTools() {
				tContent = append(tContent, anthropic.ContentBlockParamUnion{
					OfToolResult: &anthropic.ToolResultBlockParam{
						ToolUseID: c.GetId(),
						Content: []anthropic.ToolResultBlockParamContentUnion{{
							OfText: &anthropic.TextBlockParam{
								Text: string(c.GetContent()),
							},
						}},
					},
				})
			}
			if len(tContent) > 0 {
				messages = append(messages, anthropic.MessageParam{
					Role:    anthropic.MessageParamRoleUser,
					Content: tContent,
				})
			}
		case "system":
			if c := msg.GetSystem(); c != nil {
				systemPrompt = append(systemPrompt, anthropic.TextBlockParam{
					Text: string(c.GetContent()),
				})
			}
		}
	}
	return systemPrompt, messages
}

func (llc *largeLanguageCaller) StreamChatCompletion(
	ctx context.Context,
	allMessages []*protos.Message,
	options *internal_callers.ChatCompletionOptions,
	onStream func(string, *protos.Message) error,
	onMetrics func(string, *protos.Message, []*protos.Metric) error,
	onError func(string, error),
) error {
	metrics := internal_caller_metrics.NewMetricBuilder(options.RequestId)
	metrics.OnStart()

	instruction, messages := llc.BuildHistory(allMessages)
	params := llc.GetMessageNewParams(options)
	params.Messages = messages
	params.System = instruction

	client, err := llc.GetClient()
	if err != nil {
		llc.logger.Errorf("chat completion unable to get client for anthropic: %v", err)
		onError(options.Request.GetRequestId(), err)
		options.PostHook(map[string]interface{}{
			"error": err,
		}, metrics.OnFailure().Build())
		return err
	}
	options.PreHook(utils.ToJson(params))
	stream := client.Messages.NewStreaming(ctx, params)
	message := anthropic.Message{}
	if stream.Err() != nil {
		options.PostHook(map[string]interface{}{
			"result": utils.ToJson(message),
			"error":  stream.Err(),
		}, metrics.Build())
		onError(options.Request.GetRequestId(), stream.Err())
		return stream.Err()
	}

	completeMessage := &protos.AssistantMessage{}
	var currentToolCall *protos.ToolCall
	var currentContent string
	var textTokenBuffer []string // Buffer to hold text tokens temporarily
	isToolCall := false
	hasToolCalls := false // Flag to track if response contains tool calls
	for stream.Next() {
		event := stream.Current()
		err := message.Accumulate(event)
		if err != nil {
			onError(options.Request.GetRequestId(), err)
			continue
		}

		switch event := event.AsAny().(type) {
		case anthropic.ContentBlockStartEvent:
			llc.logger.Debugf("ContentBlockStartEvent %+v", event.JSON)
			switch event.ContentBlock.Type {
			case "tool_use":
				isToolCall = true
				hasToolCalls = true
				currentToolCall = &protos.ToolCall{
					Id:   event.ContentBlock.ID,
					Type: "function",
					Function: &protos.FunctionCall{
						Name: event.ContentBlock.Name,
					},
				}
			case "text":
				currentContent = ""
			}

		case anthropic.ContentBlockDeltaEvent:
			switch event.Delta.Type {
			case "text_delta":
				content := event.Delta.Text
				if content != "" && !isToolCall {
					currentContent += content
					// Buffer the token instead of streaming immediately
					textTokenBuffer = append(textTokenBuffer, content)
				}
			case "input_json_delta":
				if currentToolCall != nil {
					currentToolCall.Function.Arguments += event.Delta.PartialJSON
				}
			}

		case anthropic.ContentBlockStopEvent:
			if currentToolCall != nil {
				completeMessage.ToolCalls = append(completeMessage.ToolCalls, currentToolCall)
			}
			if currentContent != "" {
				completeMessage.Contents = append(completeMessage.Contents, currentContent)
				currentContent = ""
			}
			isToolCall = false

		case anthropic.MessageStopEvent:
			metrics.OnAddMetrics(llc.UsageMetrics(message.Usage)...)
			options.PostHook(map[string]interface{}{
				"result": utils.ToJson(message),
			}, metrics.Build())

			finalMsg := &protos.Message{
				Role: "assistant",
				Message: &protos.Message_Assistant{
					Assistant: completeMessage,
				},
			}
			// Stream text tokens only if no tool calls in response
			if !hasToolCalls {
				for _, token := range textTokenBuffer {
					if token != "" {
						tokenMsg := &protos.Message{
							Role: "assistant",
							Message: &protos.Message_Assistant{
								Assistant: &protos.AssistantMessage{
									Contents: []string{token},
								},
							},
						}
						if err := onStream(options.Request.GetRequestId(), tokenMsg); err != nil {
							llc.logger.Warnf("error streaming token: %v", err)
						}
					}
				}
			}
			// Send metrics with complete message
			onMetrics(options.Request.GetRequestId(), finalMsg, metrics.Build())
			return nil
		}
	}
	if stream.Err() != nil {
		llc.logger.Errorf("Stream error: %v", stream.Err())
		onError(options.Request.GetRequestId(), stream.Err())
		return stream.Err()
	}
	return nil
}
func (llc *largeLanguageCaller) GetMessageNewParams(opts *internal_callers.ChatCompletionOptions) anthropic.MessageNewParams {
	options := anthropic.MessageNewParams{}
	if len(opts.ToolDefinitions) > 0 {
		fns := make([]anthropic.ToolUnionParam, len(opts.ToolDefinitions))
		for idx, tl := range opts.ToolDefinitions {
			switch tl.Type {
			case "tool":
			case "function":
				fn := tl.Function
				if fn != nil {
					funcDef := &anthropic.ToolParam{
						Name: fn.Name,
					}
					if fn.Description != "" {
						funcDef.Description = anthropic.String(fn.Description)
					}
					if fn.Parameters != nil {
						funcDef.InputSchema = anthropic.ToolInputSchemaParam{
							Properties: fn.Parameters.Properties,
							Required:   fn.Parameters.Required,
						}
					}
					fns[idx] = anthropic.ToolUnionParam{
						OfTool: funcDef,
					}
				}
			}
		}
		options.Tools = fns
	}

	for key, value := range opts.ModelParameter {
		switch key {
		case "model.name":
			if mn, err := utils.AnyToString(value); err == nil {
				options.Model = anthropic.Model(mn)
			}
		case "model.max_tokens":
			if mct, err := utils.AnyToInt64(value); err == nil {
				options.MaxTokens = mct
			}

		case "model.thinking":
			if format, err := utils.AnyToJSON(value); err == nil {
				if enabled, ok := format["enabled"].(bool); ok && enabled {
					if budgetTokens, ok := format["budget_tokens"].(float64); ok {
						options.Thinking = anthropic.ThinkingConfigParamOfEnabled(int64(budgetTokens))
					}
				}
			}

		case "model.stop":
			if stopStr, err := utils.AnyToString(value); err == nil {
				options.StopSequences = strings.Split(stopStr, ",")
			}
		case "model.temperature":
			if temp, err := utils.AnyToFloat64(value); err == nil {
				options.Temperature = anthropic.Float(temp)
			}
		case "model.top_k":
			if topk, err := utils.AnyToInt64(value); err == nil {
				options.TopK = anthropic.Int(topk)
			}
		case "model.top_p":
			if topP, err := utils.AnyToFloat64(value); err == nil {
				options.TopP = anthropic.Float(topP)
			}
		}
	}
	return options
}

func (llc *largeLanguageCaller) GetChatCompletion(
	ctx context.Context,
	allMessages []*protos.Message,
	options *internal_callers.ChatCompletionOptions,
) (*protos.Message, []*protos.Metric, error) {
	metrics := internal_caller_metrics.NewMetricBuilder(options.RequestId)
	metrics.OnStart()

	client, err := llc.GetClient()
	if err != nil {
		options.PostHook(map[string]interface{}{
			"error": err,
		}, metrics.OnFailure().Build())
		return nil, metrics.OnFailure().Build(), err
	}

	instruction, messages := llc.BuildHistory(allMessages)
	params := llc.GetMessageNewParams(options)
	params.Messages = messages
	params.System = instruction

	options.PreHook(utils.ToJson(params))
	resp, err := client.Messages.New(ctx, params)
	if err != nil {
		options.PostHook(map[string]interface{}{
			"error":  err,
			"result": resp,
		}, metrics.OnFailure().Build())
		return nil, metrics.Build(), err
	}

	protoMessage := llc.convertAnthropicMessageToProto(*resp)
	metrics.OnAddMetrics(llc.UsageMetrics(resp.Usage)...)
	options.PostHook(map[string]interface{}{
		"result": resp,
		"error":  err,
	}, metrics.OnSuccess().Build())
	return protoMessage, metrics.Build(), nil
}

func (llc *largeLanguageCaller) convertAnthropicMessageToProto(message anthropic.Message) *protos.Message {
	contents := make([]string, 0)
	toolCalls := make([]*protos.ToolCall, 0)

	for _, content := range message.Content {
		switch c := content.AsAny().(type) {
		case anthropic.TextBlock:
			contents = append(contents, c.Text)
		case anthropic.ToolUseBlock:
			toolCalls = append(toolCalls, &protos.ToolCall{
				Id:   c.ID,
				Type: "function",
				Function: &protos.FunctionCall{
					Name:      c.Name,
					Arguments: string(c.JSON.Input.Raw()),
				},
			})
		}
	}

	return &protos.Message{
		Role: "assistant",
		Message: &protos.Message_Assistant{
			Assistant: &protos.AssistantMessage{
				Contents:  contents,
				ToolCalls: toolCalls,
			},
		},
	}
}
