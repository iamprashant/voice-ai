package internal_mistral_callers

import (
	"context"
	"encoding/json"
	"errors"

	internal_callers "github.com/rapidaai/api/integration-api/internal/caller"
	internal_caller_metrics "github.com/rapidaai/api/integration-api/internal/caller/metrics"
	"github.com/rapidaai/pkg/commons"
	integration_api "github.com/rapidaai/protos"
	protos "github.com/rapidaai/protos"
)

type largeLanguageCaller struct {
	Mistral
}

func NewLargeLanguageCaller(logger commons.Logger, credential *integration_api.Credential) internal_callers.LargeLanguageCaller {
	return &largeLanguageCaller{
		Mistral: mistral(logger, credential),
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
		// "model": providerModel,
	}

	msg := make([]map[string]string, 0)
	var lastRole string

	for _, cntn := range allMessages {
		currentRole := cntn.GetRole()
		if currentRole == "user" || currentRole == "system" {
			if lastRole == "user" {
				// Skip this message to ensure alternation
				continue
			}

			if user := cntn.GetUser(); user != nil {
				txt := user.GetContent()
				msg = append(msg, map[string]string{
					"role":    "user",
					"content": txt,
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
				msg = append(msg, map[string]string{
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
	res, err := llc.Call(ctx, "chat/completions", "POST", headers, requestBody)

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
	var resp MistralMessageResponse
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
	contents := make([]string, len(resp.Choices))
	metrics.OnAddMetrics(llc.UsageMetrics(resp.Usage)...)

	for idx, choice := range resp.Choices {
		contents[idx] = choice.Message.Content
	}
	message := &protos.Message{
		Role: "assistant",
		Message: &protos.Message_Assistant{
			Assistant: &protos.AssistantMessage{
				Contents: contents,
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
	llc.logger.Debugf("getting for completion for mistral")
	metrics := internal_caller_metrics.NewMetricBuilder(options.RequestId)
	metrics.OnStart()

	return nil, metrics.OnFailure().Build(), errors.New("illegal implementation")
}
