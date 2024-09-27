package web_routers

import (
	webApi "github.com/lexatic/web-backend/api/web-api"
	"github.com/lexatic/web-backend/config"
	"github.com/lexatic/web-backend/pkg/commons"
	"github.com/lexatic/web-backend/pkg/connectors"
	workflow_api "github.com/lexatic/web-backend/protos/lexatic-backend"
	"google.golang.org/grpc"
)

func DocumentApiRoute(
	Cfg *config.AppConfig,
	S *grpc.Server,
	Logger commons.Logger,
	Postgres connectors.PostgresConnector,
	Redis connectors.RedisConnector,
) {
	workflow_api.RegisterDocumentServiceServer(S,
		webApi.NewDocumentGRPCApi(Cfg,
			Logger,
			Postgres,
			Redis,
		))
}
