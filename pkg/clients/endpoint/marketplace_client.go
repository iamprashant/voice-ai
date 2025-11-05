package endpoint_client

import (
	"context"

	"github.com/rapidaai/config"
	clients "github.com/rapidaai/pkg/clients"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
	"github.com/rapidaai/pkg/types"
	endpoint_api "github.com/rapidaai/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type MarketplaceServiceClient interface {
	GetAllDeployment(c context.Context, auth types.SimplePrinciple, criterias []*endpoint_api.Criteria, paginate *endpoint_api.Paginate) (*endpoint_api.Paginated, []*endpoint_api.SearchableDeployment, error)
}

type marketplaceServiceClient struct {
	clients.InternalClient
	cfg               *config.AppConfig
	logger            commons.Logger
	marketplaceClient endpoint_api.MarketplaceServiceClient
}

func NewMarketplaceServiceClientGRPC(config *config.AppConfig, logger commons.Logger, redis connectors.RedisConnector) MarketplaceServiceClient {
	conn, err := grpc.NewClient(config.EndpointHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Errorf("Unable to create connection %v", err)
	}
	return &marketplaceServiceClient{
		InternalClient:    clients.NewInternalClient(config, logger, redis),
		cfg:               config,
		logger:            logger,
		marketplaceClient: endpoint_api.NewMarketplaceServiceClient(conn),
	}
}

func (client *marketplaceServiceClient) GetAllDeployment(c context.Context, auth types.SimplePrinciple, criterias []*endpoint_api.Criteria, paginate *endpoint_api.Paginate) (*endpoint_api.Paginated, []*endpoint_api.SearchableDeployment, error) {
	res, err := client.marketplaceClient.GetAllDeployment(client.WithAuth(c, auth), &endpoint_api.GetAllDeploymentRequest{
		Paginate:  paginate,
		Criterias: criterias,
	})
	if err != nil {
		client.logger.Errorf("error while calling to get all endpoint %v", err)
		return nil, nil, err
	}
	if !res.GetSuccess() {
		client.logger.Errorf("error while calling to get all endpoint %v", err)
		return nil, nil, err
	}
	return res.GetPaginated(), res.GetData(), nil
}
