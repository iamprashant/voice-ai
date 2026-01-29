package internal_voyageai_callers

import (
	"context"
	"encoding/json"

	internal_callers "github.com/rapidaai/api/integration-api/internal/caller"
	internal_caller_metrics "github.com/rapidaai/api/integration-api/internal/caller/metrics"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
	integration_api "github.com/rapidaai/protos"
)

type rerankingCaller struct {
	Voyageai
}

func NewRerankingCaller(logger commons.Logger, credential *integration_api.Credential) internal_callers.RerankingCaller {
	return &rerankingCaller{
		Voyageai: voyageai(logger, credential),
	}
}

func (rr *rerankingCaller) GetReranking(ctx context.Context,
	query string,
	content map[int32]string,
	options *internal_callers.RerankerOptions,
) ([]*integration_api.Reranking, []*protos.Metric, error) {
	metrics := internal_caller_metrics.NewMetricBuilder(options.RequestId)
	metrics.OnStart()

	// preseving the order

	input := make([]string, len(content))
	for k, v := range content {
		input[k] = v
	}

	request := map[string]interface{}{
		"query":     query,
		"documents": input,
		// "model":     providerModel,
	}

	headers := map[string]string{}
	options.PreHook(request)
	res, err := rr.Call(ctx, "rerank", "POST", headers, request)

	//
	if err != nil {
		rr.logger.Errorf("getting error for chat completion %v", err)
		options.PostHook(map[string]interface{}{
			"result": res,
			"error":  err,
		}, metrics.OnFailure().Build())
		return nil, metrics.Build(), err
	}
	metrics.OnSuccess()

	var resp VoyageaiRerankingResponse
	if err := json.Unmarshal([]byte(*res), &resp); err != nil {
		rr.logger.Errorf("error while parsing reranking response %v", err)
		options.PostHook(map[string]interface{}{
			"result": res,
			"error":  err,
		}, metrics.Build())
		return nil, metrics.Build(), err
	}

	output := make([]*integration_api.Reranking, len(resp.Data))

	for _, rerankedData := range resp.Data {
		// preserve the index of the chunk
		output[rerankedData.Index] = &integration_api.Reranking{
			Index:          int32(rerankedData.Index),
			Content:        content[rerankedData.Index],
			RelevanceScore: rerankedData.RelevanceScore,
		}
	}
	options.PostHook(map[string]interface{}{
		"result": res,
		"error":  err,
	}, metrics.Build())
	return output, metrics.Build(), nil
}
