package endpoint_api

import (
	config "github.com/rapidaai/api/endpoint-api/config"
	internal_services "github.com/rapidaai/api/endpoint-api/internal/service"
	internal_endpoint_service "github.com/rapidaai/api/endpoint-api/internal/service/endpoint"
	internal_log_service "github.com/rapidaai/api/endpoint-api/internal/service/log"
	commons "github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
	endpoint_grpc_api "github.com/rapidaai/protos"
)

type endpointApi struct {
	cfg                *config.EndpointConfig
	logger             commons.Logger
	postgres           connectors.PostgresConnector
	endpointService    internal_services.EndpointService
	endpointLogService internal_services.EndpointLogService
}

type endpointGRPCApi struct {
	endpointApi
}

func NewEndpointGRPCApi(config *config.EndpointConfig, logger commons.Logger,
	postgres connectors.PostgresConnector,
	redis connectors.RedisConnector,
	opensearch connectors.OpenSearchConnector,
) endpoint_grpc_api.EndpointServiceServer {
	return &endpointGRPCApi{
		endpointApi{
			cfg:                config,
			logger:             logger,
			postgres:           postgres,
			endpointService:    internal_endpoint_service.NewEndpointService(config, logger, postgres, opensearch),
			endpointLogService: internal_log_service.NewEndpointLogService(logger, postgres),
		},
	}
}
