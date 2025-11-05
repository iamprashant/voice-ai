package web_api

import (
	"context"

	config "github.com/rapidaai/api/web-api/config"
	internal_service "github.com/rapidaai/api/web-api/internal/service"
	internal_provider_service "github.com/rapidaai/api/web-api/internal/service/provider"
	commons "github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
	"github.com/rapidaai/pkg/utils"
	protos "github.com/rapidaai/protos"
)

type webProviderApi struct {
	WebApi
	cfg             *config.WebAppConfig
	logger          commons.Logger
	postgres        connectors.PostgresConnector
	redis           connectors.RedisConnector
	providerService internal_service.ProviderService
}

type webProviderRPCApi struct {
	webProviderApi
}

type webProviderGRPCApi struct {
	webProviderApi
}

func NewProviderGRPC(config *config.WebAppConfig, logger commons.Logger, postgres connectors.PostgresConnector, redis connectors.RedisConnector) protos.ProviderServiceServer {
	return &webProviderGRPCApi{
		webProviderApi{
			WebApi:          NewWebApi(config, logger, postgres, redis),
			cfg:             config,
			logger:          logger,
			postgres:        postgres,
			redis:           redis,
			providerService: internal_provider_service.NewProviderService(logger, postgres),
		},
	}
}

// GetAllProvider implements lexatic_backend.ProviderServiceServer.
func (w *webProviderGRPCApi) GetAllModelProvider(ctx context.Context, gat *protos.GetAllModelProviderRequest) (*protos.GetAllModelProviderResponse, error) {
	providers, err := w.providerService.GetAllModelProvider(ctx, gat.GetCriterias())
	if err != nil {
		return utils.Error[protos.GetAllModelProviderResponse](
			err,
			"Unable to get tool providers, please try again in sometime.")
	}
	prds := []*protos.Provider{}
	err = utils.Cast(providers, &prds)
	if err != nil {
		w.logger.Errorf("unable to cast tool provider to proto object %v", err)
	}
	return utils.PaginatedSuccess[protos.GetAllModelProviderResponse, []*protos.Provider](
		uint32(len(prds)), 0,
		prds)
}

func (w *webProviderGRPCApi) GetAllToolProvider(ctx context.Context, gat *protos.GetAllToolProviderRequest) (*protos.GetAllToolProviderResponse, error) {
	providers, err := w.providerService.GetAllToolProvider(ctx, gat.GetCriterias())
	if err != nil {
		return utils.Error[protos.GetAllToolProviderResponse](
			err,
			"Unable to get tool providers, please try again in sometime.")
	}
	prds := []*protos.ToolProvider{}
	err = utils.Cast(providers, &prds)
	if err != nil {
		w.logger.Errorf("unable to cast tool provider to proto object %v", err)
	}
	return utils.PaginatedSuccess[protos.GetAllToolProviderResponse, []*protos.ToolProvider](
		uint32(len(prds)), 1,
		prds)
}
