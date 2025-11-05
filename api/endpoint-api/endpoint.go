package endpoint_api

import (
	config "github.com/rapidaai/config"
	internal_services "github.com/rapidaai/internal/services"
	commons "github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
	endpoint_grpc_api "github.com/rapidaai/protos"
)

type endpointApi struct {
	cfg                *config.AppConfig
	logger             commons.Logger
	postgres           connectors.PostgresConnector
	endpointService    internal_services.EndpointService
	endpointLogService internal_services.EndpointLogService
}

type endpointGRPCApi struct {
	endpointApi
}

func NewEndpointGRPCApi(config *config.AppConfig, logger commons.Logger,
	postgres connectors.PostgresConnector,
	redis connectors.RedisConnector,
	opensearch connectors.OpenSearchConnector,
) endpoint_grpc_api.EndpointServiceServer {
	return &endpointGRPCApi{
		endpointApi{
			cfg:                config,
			logger:             logger,
			postgres:           postgres,
			endpointService:    internal_services.NewEndpointService(config, logger, postgres, opensearch),
			endpointLogService: internal_services.NewEndpointLogService(logger, postgres),
		},
	}
}
