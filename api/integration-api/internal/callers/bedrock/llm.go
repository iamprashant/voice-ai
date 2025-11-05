package internal_bedrock_callers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	bedrock_types "github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"

	internal_callers "github.com/rapidaai/internal/callers"
	internal_caller_metrics "github.com/rapidaai/internal/callers/metrics"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	integration_api "github.com/rapidaai/protos"
	lexatic_backend "github.com/rapidaai/protos"
)

type largeLanguageCaller struct {
	Bedrock
}

// StreamChatCompletion implements internal_callers.LargeLanguageCaller.
func (*largeLanguageCaller) StreamChatCompletion(
	ctx context.Context,
	// providerModel string,
	allMessages []*lexatic_backend.Message,
	options *internal_callers.ChatCompletionOptions,
	onStream func(types.Message) error,
	onMetrics func(*types.Message, types.Metrics) error,
	onError func(err error),
) error {
	panic("unimplemented")
}

// GetChatCompletion implements internal_callers.LargeLanguageCaller.
func (llc *largeLanguageCaller) GetChatCompletion(ctx context.Context,
	// providerModel string,
	allMessages []*lexatic_backend.Message, options *internal_callers.ChatCompletionOptions) (*types.Message, types.Metrics, error) {
	llc.logger.Debugf("chat complition request started for openai")
	//
	// Working with chat complition with vision
	//
	metrics := internal_caller_metrics.NewMetricBuilder(options.RequestId)
	metrics.OnStart()

	client, err := llc.GetClient()
	if err != nil {
		llc.logger.Errorf("chat complition unable to get client for openai %v", err)
		return nil, metrics.OnFailure().Build(), err
	}

	msg := make([]bedrock_types.Message, len(allMessages))
	systemMsg := make([]bedrock_types.SystemContentBlock, 0)
	idx := 0
	for _, cntn := range allMessages {
		if len(cntn.GetContents()) == 0 {
			// there might be problem in initiator
			continue
		}

		allMessageContent := make([]bedrock_types.ContentBlock, len(cntn.GetContents()))
		for idy, ct := range cntn.GetContents() {
			if ct.ContentType == commons.TEXT_CONTENT.String() {
				content := string(ct.GetContent())
				allMessageContent[idy] = &bedrock_types.ContentBlockMemberText{
					Value: content,
				}
			}
			if ct.ContentType == commons.IMAGE_CONTENT.String() {
				if ct.GetContentFormat() == commons.AUDIO_CONTENT_FORMAT_RAW.String() {
					allMessageContent[idy] = &bedrock_types.ContentBlockMemberImage{
						Value: bedrock_types.ImageBlock{
							Source: &bedrock_types.ImageSourceMemberBytes{
								Value: ct.GetContent(),
							},
						},
					}
				}
			}
			if cntn.GetRole() == "user" || cntn.GetRole() == "system" {
				msg[idx] = bedrock_types.Message{Content: allMessageContent, Role: "user"}
			} else {
				msg[idx] = bedrock_types.Message{Content: allMessageContent, Role: "assistant"}
			}
		}
		idx++
	}

	llmRequest := bedrockruntime.ConverseInput{
		Messages: msg[:idx],
		// ModelId:  to.Ptr(providerModel),
		System: systemMsg,
	}
	options.AIOptions.PreHook(utils.ToJson(llmRequest))
	// single minute timeout and cancellable by the client as context will get cancel
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	resp, err := client.Converse(ctx, &llmRequest)

	if err != nil {
		options.AIOptions.PostHook(
			map[string]interface{}{
				"error":  err,
				"result": resp,
			}, metrics.OnFailure().Build())
		llc.logger.Errorf("complition response from bedrock %v", err)
		return nil, metrics.OnFailure().Build(), err
	}

	metrics.OnSuccess()
	output := make([]*types.Content, 0)
	// var union types.ConverseOutput
	switch v := resp.Output.(type) {
	case *bedrock_types.ConverseOutputMemberMessage:
		for _, choice := range v.Value.Content {
			switch x := choice.(type) {
			case *bedrock_types.ContentBlockMemberImage:
				switch s := x.Value.Source.(type) {
				case *bedrock_types.ImageSourceMemberBytes:
					output = append(output, &types.Content{
						ContentType:   commons.IMAGE_CONTENT.String(),
						ContentFormat: commons.IMAGE_CONTENT_FORMAT_RAW.String(),
						Content:       s.Value,
					})
				}
			case *bedrock_types.ContentBlockMemberText:
				output = append(output, &types.Content{
					ContentType:   commons.TEXT_CONTENT.String(),
					ContentFormat: commons.TEXT_CONTENT_FORMAT_RAW.String(),
					Content:       []byte(x.Value),
				})

			}
		}

	}

	options.AIOptions.PostHook(map[string]interface{}{
		"result": resp,
	}, metrics.Build())
	return &types.Message{
		Contents: output,
	}, metrics.Build(), nil

}

// GetCompletion implements internal_callers.LargeLanguageCaller.
func (llc *largeLanguageCaller) GetCompletion(ctx context.Context,
	// providerModel string,
	prompts []string, options *internal_callers.CompletionOptions) ([]*types.Content, types.Metrics, error) {
	llc.logger.Debugf("chat complition request started for openai")
	//
	// Working with chat complition with vision
	//
	metrics := internal_caller_metrics.NewMetricBuilder(options.RequestId)
	metrics.OnStart()

	client, err := llc.GetClient()
	if err != nil {
		llc.logger.Errorf("chat complition unable to get client for openai %v", err)
		return nil, metrics.OnFailure().Build(), err
	}

	// Create a map[string]interface{} as the request body
	requestBody := map[string]interface{}{
		"prompt": "Translate the following text to French: 'Hello, world!'",
	}

	// Marshal the map into JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		llc.logger.Errorf("failed to marshal request body, %v", err)
	}

	llmRequest := bedrockruntime.InvokeModelInput{
		Body: jsonBody,
		// ModelId: to.Ptr(providerModel),
		Accept: aws.String("application/json"),
	}
	options.AIOptions.PreHook(utils.ToJson(llmRequest))
	// single minute timeout and cancellable by the client as context will get cancel
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	resp, err := client.InvokeModel(ctx, &llmRequest, nil)

	if err != nil {
		options.AIOptions.PostHook(
			map[string]interface{}{
				"error":  err,
				"result": resp,
			}, metrics.OnFailure().Build())
		llc.logger.Errorf("complition response from bedrock %v", err)
		return nil, metrics.OnFailure().Build(), err
	}

	metrics.OnSuccess()
	var response AWSCompletionResponse
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		options.AIOptions.PostHook(
			map[string]interface{}{
				"error":  err,
				"result": resp,
			}, metrics.OnFailure().Build())
		llc.logger.Errorf("complition to parse api response from bedrock %v", err)
		return nil, metrics.OnFailure().Build(), err
	}

	output := make([]*types.Content, len(response.Completions))

	for idx, choice := range response.Completions {
		if choice.Data != nil {
			output[idx] = &types.Content{
				ContentType:   commons.TEXT_CONTENT.String(),
				ContentFormat: commons.TEXT_CONTENT_FORMAT_RAW.String(),
				Content:       []byte(choice.Data.Text),
			}
		}
	}
	options.AIOptions.PostHook(map[string]interface{}{
		"result": resp,
	}, metrics.Build())
	return output, metrics.Build(), nil

}

func NewLargeLanguageCaller(logger commons.Logger, credential *integration_api.Credential) internal_callers.LargeLanguageCaller {
	logger.Debugf("creating large language instance with proxy")
	return &largeLanguageCaller{
		Bedrock: bedrock(logger, credential),
	}
}
