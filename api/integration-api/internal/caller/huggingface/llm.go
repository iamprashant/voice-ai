package internal_huggingface_callers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	internal_callers "github.com/rapidaai/api/integration-api/internal/caller"
	internal_caller_metrics "github.com/rapidaai/api/integration-api/internal/caller/metrics"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
)

type largeLanguageCaller struct {
	Huggingface
}

func NewLargeLanguageCaller(logger commons.Logger, credential *protos.Credential) internal_callers.LargeLanguageCaller {
	return &largeLanguageCaller{
		Huggingface: huggingface(logger,
			DEFAULT_URL, credential),
	}
}

// StreamChatCompletion implements internal_callers.LargeLanguageCaller.
func (*largeLanguageCaller) StreamChatCompletion(
	ctx context.Context,
	// providerModel string,
	allMessages []*protos.Message,
	options *internal_callers.ChatCompletionOptions,
	onStream func(string, *protos.Message) error,
	onMetrics func(string, *protos.Message, []*protos.Metric) error,
	onError func(string, error),
) error {
	panic("unimplemented")
}

func (llc *largeLanguageCaller) GetChatCompletion(
	ctx context.Context,
	// providerModel string,
	allMessages []*protos.Message,
	options *internal_callers.ChatCompletionOptions,
) (*protos.Message, []*protos.Metric, error) {
	llc.logger.Debugf("getting chat completion from google llc")
	//
	// Working with chat completion with vision
	//
	metrics := internal_caller_metrics.NewMetricBuilder(options.RequestId)
	metrics.OnStart()
	requestBody := map[string]interface{}{
		// "model":      providerModel,
		"max_tokens": 1024,
	}

	msg := make([]map[string]interface{}, 0)
	var lastRole string

	for _, cntn := range allMessages {
		currentRole := cntn.GetRole()
		if currentRole == "user" || currentRole == "system" {
			if lastRole == "user" {
				// Skip this message to ensure alternation
				continue
			}
			if user := cntn.GetUser(); user != nil {
				contents := make([]map[string]interface{}, 0)
				contents = append(contents, map[string]interface{}{
					"type": "text",
					"text": user.GetContent(),
				})
				msg = append(msg, map[string]interface{}{
					"role":    "user",
					"content": contents,
				})
				lastRole = "user"
			}
		}

		if currentRole == "assistant" {
			if lastRole == "assistant" {
				// Skip this message to ensure alternation
				continue
			}
			if assistant := cntn.GetAssistant(); assistant != nil && len(assistant.GetContents()) > 0 {
				txt := assistant.GetContents()[0]
				msg = append(msg, map[string]interface{}{
					"role":    "assistant",
					"content": txt,
				})
				lastRole = "assistant"
			}
		}
	}

	requestBody["messages"] = msg
	headers := map[string]string{}
	options.PreHook(requestBody)
	res, err := llc.Call(ctx, fmt.Sprintf("models/%s", ""),
		"POST", headers, requestBody)

	//
	if err != nil {
		llc.logger.Errorf("getting error for chat completion %v", err)
		options.PostHook(map[string]interface{}{
			"result": res,
			"error":  err,
		}, metrics.OnFailure().Build())
		return nil, metrics.Build(), err
	}
	metrics.OnSuccess()
	var resp HuggingfaceInferenceResponse
	if err := json.Unmarshal([]byte(*res), &resp); err != nil {
		llc.logger.Errorf("error while parsing chat completion response %v", err)
		options.PostHook(map[string]interface{}{
			"result": res,
			"error":  err,
		}, metrics.Build())
		return nil, metrics.Build(), err
	}

	//
	//
	message := &protos.Message{
		Role: "assistant",
		Message: &protos.Message_Assistant{
			Assistant: &protos.AssistantMessage{
				Contents: []string{resp.Text},
			},
		},
	}
	options.PostHook(map[string]interface{}{
		"result": res,
	}, metrics.Build())
	return message, metrics.Build(), nil
}

func (llc *largeLanguageCaller) GetCompletion(
	ctx context.Context,
	providerModel string,
	prompts []string,
	options *internal_callers.CompletionOptions,
) ([]string, []*protos.Metric, error) {
	//
	// Working with chat completion with vision
	//
	llc.logger.Debugf("getting for completion for huggingface")
	metrics := internal_caller_metrics.NewMetricBuilder(options.RequestId)
	metrics.OnStart()

	var buffer bytes.Buffer
	buffer.WriteString("\n\nHuman: %s\n\nAssistant:")
	requestBody := map[string]interface{}{
		"prompt":               fmt.Sprintf(buffer.String(), prompts[0]),
		"model":                providerModel,
		"max_tokens_to_sample": 100,
	}

	headers := map[string]string{}
	options.PreHook(requestBody)
	res, err := llc.Call(ctx, fmt.Sprintf("models/%s", providerModel), "POST", headers, requestBody)
	//
	//
	if err != nil {
		llc.logger.Errorf("getting error for completion %v", err)
		options.PostHook(map[string]interface{}{
			"result": res,
			"error":  err,
		}, metrics.OnFailure().Build())
		return nil, metrics.Build(), err
	}
	metrics.OnSuccess()
	var resp HuggingfaceInferenceResponse
	if err := json.Unmarshal([]byte(*res), &resp); err != nil {
		llc.logger.Errorf("error while parsing completion response %v", err)
		options.PostHook(map[string]interface{}{
			"result": res,
			"error":  err,
		}, metrics.Build())
		return nil, metrics.Build(), err
	}

	options.PostHook(map[string]interface{}{
		"result": res,
	}, metrics.Build())
	return []string{resp.Text}, metrics.Build(), nil
}
