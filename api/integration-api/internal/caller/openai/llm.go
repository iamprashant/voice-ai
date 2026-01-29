package internal_openai_callers

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"

	internal_callers "github.com/rapidaai/api/integration-api/internal/caller"
	internal_caller_metrics "github.com/rapidaai/api/integration-api/internal/caller/metrics"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	protos "github.com/rapidaai/protos"
)

type largeLanguageCaller struct {
	OpenAI
}

func NewLargeLanguageCaller(logger commons.Logger, credential *protos.Credential) internal_callers.LargeLanguageCaller {
	return &largeLanguageCaller{
		OpenAI: openAI(logger, credential),
	}
}

func (llc *largeLanguageCaller) getChatCompletionOptions(
	opts *internal_callers.ChatCompletionOptions,
) openai.ChatCompletionNewParams {
	options := openai.ChatCompletionNewParams{}
	if len(opts.ToolDefinitions) > 0 {
		fns := make([]openai.ChatCompletionToolParam, len(opts.ToolDefinitions))
		for idx, tl := range opts.ToolDefinitions {
			switch tl.Type {
			case "tool":
			case "function":
				fn := tl.Function
				if fn != nil {
					funcDef := openai.FunctionDefinitionParam{
						Name: fn.Name,
					}
					if fn.Description != "" {
						funcDef.Description = openai.String(fn.Description)
					}
					// Always set parameters with valid JSON schema format
					if fn.Parameters != nil {
						funcDef.Parameters = fn.Parameters.ToMap()
					} else {
						// Default empty parameters with properties field for valid schema
						funcDef.Parameters = map[string]interface{}{
							"type":       "object",
							"properties": map[string]interface{}{},
						}
					}
					fns[idx] = openai.ChatCompletionToolParam{
						Function: funcDef,
					}
				}
			}
		}
		options.Tools = fns
	}

	for key, value := range opts.ModelParameter {
		switch key {
		case "model.name":
			if modelName, err := utils.AnyToString(value); err == nil {
				options.Model = modelName
			}
		case "model.user":
			if user, err := utils.AnyToString(value); err == nil {
				options.User = openai.String(user)
			}
		case "model.reasoning_effort":
			if re, err := utils.AnyToString(value); err == nil {
				options.ReasoningEffort = shared.ReasoningEffort(re)
			}
		case "model.seed":
			if seed, err := utils.AnyToInt64(value); err == nil {
				options.Seed = openai.Int(seed)
			}
		case "model.service_tier":
			if st, err := utils.AnyToString(value); err == nil {
				options.ServiceTier = openai.ChatCompletionNewParamsServiceTier(st)
			}
		case "model.top_logprobs":
			if tl, err := utils.AnyToInt64(value); err == nil {
				options.TopLogprobs = openai.Int(tl)
			}
		case "model.metadata":
			format, _ := utils.AnyToString(value)
			var mtd map[string]string
			if err := json.Unmarshal([]byte(format), &mtd); err == nil {
				options.Metadata = shared.Metadata(mtd)
			}
		case "model.frequency_penalty":
			if fp, err := utils.AnyToFloat64(value); err == nil {
				options.FrequencyPenalty = openai.Float(fp)
			}
		case "model.temperature":
			if temp, err := utils.AnyToFloat64(value); err == nil {
				options.Temperature = openai.Float(temp)
			}
		case "model.top_p":
			if topP, err := utils.AnyToFloat64(value); err == nil {
				options.TopP = openai.Float(topP)
			}
		case "model.presence_penalty":
			if pp, err := utils.AnyToFloat64(value); err == nil {
				options.PresencePenalty = openai.Float(pp)
			}
		case "model.max_completion_tokens":
			if maxTokens, err := utils.AnyToInt64(value); err == nil {
				options.MaxTokens = openai.Int(maxTokens)
			}
		case "model.stop":
			if stopStr, err := utils.AnyToString(value); err == nil {
				for _, stopper := range strings.Split(stopStr, ",") {
					if strings.TrimSpace(stopper) != "" {
						options.Stop.OfStringArray = append(options.Stop.OfStringArray, stopper)
					}
				}
			}
		case "model.tool_choice":
			if choice, err := utils.AnyToString(value); err == nil {
				switch choice {
				case "auto":
					options.ToolChoice = openai.ChatCompletionToolChoiceOptionUnionParam{
						OfAuto: openai.String("auto"),
					}
				case "required":
					options.ToolChoice = openai.ChatCompletionToolChoiceOptionUnionParam{
						OfAuto: openai.String("required"),
					}
				case "none":
					options.ToolChoice = openai.ChatCompletionToolChoiceOptionUnionParam{
						OfAuto: openai.String("none"),
					}
				default:
					options.ToolChoice = openai.ChatCompletionToolChoiceOptionUnionParam{
						OfAuto: openai.String("none"),
					}
				}
			}
		case "model.response_format":
			if format, err := utils.AnyToJSON(value); err == nil {
				switch format["type"].(string) {
				case "json_object":
					options.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
						OfJSONObject: &openai.ResponseFormatJSONObjectParam{},
					}
				case "text":
					options.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{}
				case "json_schema":
					if schemaData, ok := format["json_schema"].(map[string]interface{}); ok {
						jsonSchemaParam := shared.ResponseFormatJSONSchemaJSONSchemaParam{}
						jsonData, err := json.Marshal(schemaData)
						if err == nil {
							json.Unmarshal(jsonData, &jsonSchemaParam)
						}
						options.ResponseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
							OfJSONSchema: &shared.ResponseFormatJSONSchemaParam{
								JSONSchema: jsonSchemaParam,
							},
						}
					}
				}
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
		llc.logger.Errorf("chat completion unable to get client for openai %v", err)
		return nil, metrics.OnFailure().Build(), err
	}

	// message and options
	llmRequest := llc.getChatCompletionOptions(options)
	llmRequest.Messages = llc.BuildHistory(allMessages)

	// prehook
	options.PreHook(utils.ToJson(llmRequest))

	// chat complitions
	resp, err := client.Chat.Completions.New(ctx, llmRequest)
	if err != nil {
		llc.logger.Errorf("chat completion failed to get response from openai %v", err)
		options.PostHook(map[string]interface{}{
			"error":  err,
			"result": resp,
		}, metrics.OnFailure().Build())
		return nil, metrics.OnFailure().Build(), err
	}

	assistantMsg := &protos.AssistantMessage{
		Contents:  make([]string, 0),
		ToolCalls: make([]*protos.ToolCall, 0),
	}
	metrics.OnSuccess()

	for _, choice := range resp.Choices {
		switch choice.FinishReason {
		case "length", "content_filter":
		case "stop":
			assistantMsg.Contents = append(assistantMsg.Contents, choice.Message.Content)
		case "function_call", "tool_calls":
			if choice.Message.ToolCalls != nil {
				for _, tool := range choice.Message.ToolCalls {
					if tool.Type == "function" {
						assistantMsg.ToolCalls = append(assistantMsg.ToolCalls, &protos.ToolCall{
							Id:   tool.ID,
							Type: string(tool.Type),
							Function: &protos.FunctionCall{
								Name:      tool.Function.Name,
								Arguments: tool.Function.Arguments,
							},
						})
					}
				}
			}
		}
	}

	options.PostHook(map[string]interface{}{
		"result": resp,
	}, metrics.Build())
	return &protos.Message{
		Role: "assistant",
		Message: &protos.Message_Assistant{
			Assistant: assistantMsg,
		},
	}, metrics.Build(), nil
}

func (llc *largeLanguageCaller) StreamChatCompletion(
	ctx context.Context,
	allMessages []*protos.Message,
	options *internal_callers.ChatCompletionOptions,
	onStream func(string, *protos.Message) error,
	onMetrics func(string, *protos.Message, []*protos.Metric) error,
	onError func(string, error),
) error {
	start := time.Now()
	metrics := internal_caller_metrics.NewMetricBuilder(options.RequestId)
	metrics.OnStart()

	client, err := llc.GetClient()
	if err != nil {
		llc.logger.Errorf("chat completion unable to get client for openai: %v", err)
		onError(options.Request.GetRequestId(), err)
		options.PostHook(map[string]interface{}{
			"error": err,
		}, metrics.OnFailure().Build())
		return err
	}

	completionsOptions := llc.getChatCompletionOptions(options)
	completionsOptions.Messages = llc.BuildHistory(allMessages)
	options.PreHook(utils.ToJson(completionsOptions))
	llc.logger.Benchmark("Openai.llm.GetChatCompletion.llmRequestPrepare", time.Since(start))

	// Get streaming response
	resp := client.Chat.Completions.NewStreaming(ctx, completionsOptions)
	if resp.Err() != nil {
		llc.logger.Errorf("Failed to get chat completions stream: %v", resp.Err())
		options.PostHook(map[string]interface{}{
			"result": utils.ToJson(resp),
			"error":  resp.Err(),
		}, metrics.Build())
		onError(options.Request.GetRequestId(), resp.Err())
		return resp.Err()
	}
	defer resp.Close()
	assistantMsg := &protos.AssistantMessage{
		Contents:  make([]string, 0),
		ToolCalls: make([]*protos.ToolCall, 0),
	}
	contentBuffer := make([]string, 0) // Buffer for accumulating content per choice

	accumulate := openai.ChatCompletionAccumulator{}
	for resp.Next() {
		chatCompletions := resp.Current()
		accumulate.AddChunk(chatCompletions)

		if _, ok := accumulate.JustFinishedContent(); ok {
			metrics.OnAddMetrics(llc.GetComplitionUsages(accumulate.Usage)...)
			metrics.OnSuccess()
			options.PostHook(map[string]interface{}{
				"result": utils.ToJson(accumulate),
			}, metrics.Build())

			// Finalize content from buffer
			assistantMsg.Contents = contentBuffer
			protoMsg := &protos.Message{
				Role: "assistant",
				Message: &protos.Message_Assistant{
					Assistant: assistantMsg,
				},
			}
			// Stream if no tool calls
			if len(assistantMsg.ToolCalls) == 0 {
				if err := onStream(options.Request.GetRequestId(), protoMsg); err != nil {
					llc.logger.Warnf("error streaming complete message: %v", err)
				}
			}
			onMetrics(options.Request.GetRequestId(), protoMsg, metrics.Build())
			return nil
		}

		if tool, ok := accumulate.JustFinishedToolCall(); ok {
			assistantMsg.ToolCalls = append(assistantMsg.ToolCalls, &protos.ToolCall{
				Id: tool.ID,
				Function: &protos.FunctionCall{
					Name:      tool.Name,
					Arguments: tool.Arguments,
				},
			})
			metrics.OnAddMetrics(llc.GetComplitionUsages(accumulate.Usage)...)
			options.PostHook(map[string]interface{}{
				"result": utils.ToJson(accumulate),
			}, metrics.Build())

			// Don't stream if tool calls are present
			assistantMsg.Contents = contentBuffer
			protoMsg := &protos.Message{
				Role: "assistant",
				Message: &protos.Message_Assistant{
					Assistant: assistantMsg,
				},
			}
			onMetrics(options.Request.GetRequestId(), protoMsg, metrics.Build())
			return nil
		}

		for i, choice := range chatCompletions.Choices {
			content := choice.Delta.Content
			if content != "" {
				// Accumulate content per choice index
				if len(contentBuffer) <= i {
					contentBuffer = append(contentBuffer, content)
				} else {
					contentBuffer[i] += content
				}
			}
		}
	}

	return nil
}

func (llc *largeLanguageCaller) BuildHistory(allMessages []*protos.Message) []openai.ChatCompletionMessageParamUnion {
	msg := make([]openai.ChatCompletionMessageParamUnion, 0)
	for _, cntn := range allMessages {
		switch cntn.GetRole() {
		case ChatRoleUser:
			if user := cntn.GetUser(); user != nil {
				var messageContent []openai.ChatCompletionContentPartUnionParam
				messageContent = append(messageContent, openai.ChatCompletionContentPartUnionParam{
					OfText: &openai.ChatCompletionContentPartTextParam{
						Text: user.GetContent(),
					},
				})
				msg = append(msg, openai.UserMessage(messageContent))
			}
		case ChatRoleAssistant:
			if assistant := cntn.GetAssistant(); assistant != nil {
				txtContent := strings.Join(assistant.GetContents(), "")
				toolCalls := assistant.GetToolCalls()
				assistantMessage := openai.ChatCompletionAssistantMessageParam{}
				if len(txtContent) > 0 || len(toolCalls) > 0 {
					if len(txtContent) > 0 {
						assistantMessage.Content = openai.ChatCompletionAssistantMessageParamContentUnion{
							OfString: openai.String(txtContent),
						}
					}
					if len(toolCalls) > 0 {
						fctCall := make([]openai.ChatCompletionMessageToolCallParam, 0)
						for _, ttc := range toolCalls {
							fctCall = append(fctCall, openai.ChatCompletionMessageToolCallParam{
								ID: ttc.GetId(),
								Function: openai.ChatCompletionMessageToolCallFunctionParam{
									Name:      ttc.GetFunction().GetName(),
									Arguments: ttc.GetFunction().GetArguments(),
								},
							})
						}
						assistantMessage.ToolCalls = fctCall
					}
					msg = append(msg, openai.ChatCompletionMessageParamUnion{
						OfAssistant: &assistantMessage,
					})
				}
			}

		case ChatRoleSystem:
			if system := cntn.GetSystem(); system != nil {
				txtContent := system.GetContent()
				if len(txtContent) > 0 {
					msg = append(msg, openai.SystemMessage(txtContent))
				}
			}

		case ChatRoleTool:
			if tool := cntn.GetTool(); tool != nil {
				for _, t := range tool.GetTools() {
					msg = append(msg, openai.ToolMessage(t.GetContent(), t.GetId()))
				}
			}
		}
	}
	return msg
}
