package workflow_routers

import (
	"github.com/gin-gonic/gin"
	healthCheckApi "github.com/rapidaai/api/health-check-api"
	"github.com/rapidaai/config"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
)

func HealthCheckRoutes(cfg *config.AppConfig, engine *gin.Engine, logger commons.Logger, postgres connectors.PostgresConnector) {
	logger.Info("Internal HealthCheckRoutes and Connectors added to engine.")
	apiv1 := engine.Group("")
	hcApi := healthCheckApi.New(cfg, logger, postgres)
	{
		apiv1.GET("/readiness/", hcApi.Readiness)
		apiv1.GET("/healthz/", hcApi.Healthz)
	}
}
