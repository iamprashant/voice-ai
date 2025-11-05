package web_client

import (
	"context"
	"errors"
	"time"

	"github.com/rapidaai/config"
	clients "github.com/rapidaai/pkg/clients"
	commons "github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
	provider_api "github.com/rapidaai/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProviderServiceClient interface {
	// GetModel(c context.Context, auth types.SimplePrinciple, modelId uint64) (*provider_api.ProviderModel, error)
	// GetAllProviderModel(c context.Context, modelType string) ([]*provider_api.ProviderModel, error)
	GetAllModelProviders(c context.Context) ([]*provider_api.Provider, error)
}

type providerServiceClient struct {
	clients.InternalClient
	cfg            *config.AppConfig
	logger         commons.Logger
	providerClient provider_api.ProviderServiceClient
}

func NewProviderServiceClientGRPC(config *config.AppConfig, logger commons.Logger, redis connectors.RedisConnector) ProviderServiceClient {
	logger.Debugf("conntecting to provider client with %s", config.ProviderHost)
	conn, err := grpc.NewClient(config.WebHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatalf("Unable to create connection %v to the provider service", err)
	}
	// providerClient :=
	return &providerServiceClient{
		clients.NewInternalClient(config, logger, redis),
		config,
		logger,
		provider_api.NewProviderServiceClient(conn),
	}
}

// // Used by experiment and test. modelType is type of experiment text, chat etc
// func (client *providerServiceClient) GetAllProviderModel(c context.Context, modelType string) ([]*provider_api.ProviderModel, error) {
// 	start := time.Now()
// 	key, value := "endpoint", "complete"
// 	if modelType == "image" {
// 		value = "text-to-image"
// 	}

// 	if modelType == "chat" {
// 		value = "chat-complete"
// 	}

// 	res, err := client.providerClient.GetAllModel(c, &provider_api.GetAllModelRequest{Criterias: []*provider_api.Criteria{{Key: key, Value: value}}})
// 	if err != nil {
// 		client.logger.Errorf("got an error while calling provider service client error here %v", err)
// 		client.logger.Benchmark("providerServiceClient.GetAllProviderModel", time.Since(start))
// 		return nil, err
// 	}
// 	if !res.GetSuccess() {
// 		client.logger.Benchmark("providerServiceClient.GetAllProviderModel", time.Since(start))
// 		return nil, errors.New("illegal request with provider clients")
// 	}
// 	client.logger.Benchmark("providerServiceClient.GetAllProviderModel", time.Since(start))
// 	return res.GetData(), nil
// }

// // GetModel implements internal_clients.ProviderServiceClient.
// func (client *providerServiceClient) GetModel(c context.Context, auth types.SimplePrinciple, modelId uint64) (*provider_api.ProviderModel, error) {
// 	start := time.Now()
// 	// Generate cache key
// 	cacheKey := client.CacheKey(c, "GetModel", fmt.Sprintf("%d", *auth.GetCurrentOrganizationId()), fmt.Sprintf("%d", modelId))

// 	// Retrieve data from cache
// 	cachedValue := client.Retrieve(c, cacheKey)
// 	if cachedValue.HasError() {
// 		client.logger.Errorf("Cache missed for the request: %v", cachedValue.Err)
// 	}

// 	// Initialize data variable
// 	data := &provider_api.ProviderModel{}
// 	err := cachedValue.ResultStruct(data)
// 	if err != nil {
// 		client.logger.Errorf("Failed to parse cached data: %v", err)
// 		res, err := client.providerClient.GetModel(client.WithAuth(c, auth), &provider_api.GetModelRequest{ModelId: modelId})
// 		if err != nil {
// 			client.logger.Errorf("Failed to get model from provider service: %v", err)
// 			return nil, err
// 		}

// 		if res.GetSuccess() {
// 			if res.GetData() != nil {
// 				client.Cache(c, cacheKey, res.GetData())
// 			}
// 			client.logger.Benchmark("providerServiceClient.GetModel", time.Since(start))
// 			return res.GetData(), nil
// 		}

// 	}
// 	// Log benchmarking information
// 	client.logger.Benchmark("providerServiceClient.GetModel", time.Since(start))
// 	return data, nil

// }

func (client *providerServiceClient) GetAllModelProviders(c context.Context) ([]*provider_api.Provider, error) {
	start := time.Now()
	res, err := client.providerClient.GetAllModelProvider(c, &provider_api.GetAllModelProviderRequest{})
	if err != nil {
		client.logger.Errorf("Failed to get all providers from provider service: %v", err)
		return nil, err
	}
	if res.GetSuccess() {
		client.logger.Benchmark("providerServiceClient.GetAllModelProviders", time.Since(start))
		return res.GetData(), nil
	}
	client.logger.Benchmark("providerServiceClient.GetAllModelProviders", time.Since(start))
	return nil, errors.New("unknown error occured while fetching providers")
}
