package web_handler

import (
	"context"

	internal_service "github.com/rapidaai/api/web-api/internal/service"
	internal_organization_service "github.com/rapidaai/api/web-api/internal/service/organization"
	internal_provider_service "github.com/rapidaai/api/web-api/internal/service/provider"
	internal_user_service "github.com/rapidaai/api/web-api/internal/service/user"
	config "github.com/rapidaai/config"
	commons "github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	web_api "github.com/rapidaai/protos"
)

type WebApi struct {
	cfg             *config.AppConfig
	logger          commons.Logger
	postgres        connectors.PostgresConnector
	redis           connectors.RedisConnector
	userService     internal_service.UserService
	providerService internal_service.ProviderService
	orgService      internal_service.OrganizationService
}

func NewWebApi(cfg *config.AppConfig, logger commons.Logger, postgres connectors.PostgresConnector, redis connectors.RedisConnector) WebApi {
	return WebApi{
		cfg: cfg, logger: logger, postgres: postgres, redis: redis,
		userService:     internal_user_service.NewUserService(logger, postgres),
		orgService:      internal_organization_service.NewOrganizationService(logger, postgres),
		providerService: internal_provider_service.NewProviderService(logger, postgres),
	}
}

func (w *WebApi) GetUser(c context.Context, auth types.SimplePrinciple, userId uint64) *web_api.User {
	usr, err := w.userService.GetUser(c, userId)
	if err != nil {
		w.logger.Errorf("unable to get user form the database %+v", err)
		return nil
	}
	ot := &web_api.User{}
	err = utils.Cast(usr, ot)
	if err != nil {
		w.logger.Errorf("unable to cast project to proto object %v", err)
	}
	return ot
}

func (w *WebApi) GetOrganization(c context.Context, auth types.SimplePrinciple, orgId uint64) *web_api.Organization {
	org, err := w.orgService.Get(c, orgId)
	if err != nil {
		w.logger.Errorf("unable to get organization form the database %+v", err)
		return nil
	}
	ot := &web_api.Organization{}
	err = utils.Cast(org, ot)
	if err != nil {
		w.logger.Errorf("unable to cast project to proto object %v", err)
	}
	return ot
}
