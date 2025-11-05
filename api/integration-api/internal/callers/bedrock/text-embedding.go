package internal_bedrock_callers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	internal_callers "github.com/rapidaai/internal/callers"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/types"
	type_enums "github.com/rapidaai/pkg/types/enums"
	"github.com/rapidaai/pkg/utils"
	integration_api "github.com/rapidaai/protos"
)

type embeddingCaller struct {
	Bedrock
}

type BedrockEmbedding struct {
	Embedding []float32
}

func NewEmbeddingCaller(logger commons.Logger, credential *integration_api.Credential) internal_callers.EmbeddingCaller {
	return &embeddingCaller{
		bedrock(logger, credential),
	}
}

// GetText2Speech implements internal_callers.Text2SpeechCaller.
func (brc *embeddingCaller) GetEmbedding(ctx context.Context,
	// providerModel string,
	content map[int32]string,
	options *internal_callers.EmbeddingOptions) ([]*integration_api.Embedding, types.Metrics, error) {
	//
	// Working with chat complition with vision
	//
	start := time.Now()
	timeMetric := &types.Metric{
		Name:        type_enums.TIME_TAKEN.String(),
		Value:       fmt.Sprintf("%d", int64(time.Since(start))),
		Description: "Time taken to serve the llm request",
	}
	client, err := brc.GetClient()
	if err != nil {
		return nil, types.Metrics{timeMetric}, err
	}

	// single minute timeout and cancellable by the client as context will get cancel
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	output := make([]*integration_api.Embedding, len(content))

	for idx, st := range content {
		payload := map[string]string{
			"InputText": st,
		}
		payloadBytes, err := json.Marshal(payload)
		if err != nil {
			log.Fatal(err)
		}
		invokeOut, err := client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
			Body: payloadBytes,
			// ModelId:     aws.String(providerModel),
			ContentType: aws.String("application/json"),
		})

		var resp BedrockEmbedding
		err = json.Unmarshal(invokeOut.Body, &resp)
		if err != nil {
			log.Fatal("failed to unmarshal", err)
		}
		// preserve the order
		output[idx] = &integration_api.Embedding{
			Index:     idx,
			Embedding: utils.EmbeddingToFloat64(resp.Embedding),
			Base64:    utils.EmbeddingToBase64(utils.EmbeddingToFloat64(resp.Embedding)),
		}

	}
	return output, types.Metrics{timeMetric}, nil

}
