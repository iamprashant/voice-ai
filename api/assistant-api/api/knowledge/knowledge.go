package knowledge_api

import (
	internal_services "github.com/rapidaai/api/internal/services"
	internal_knowledge_service "github.com/rapidaai/api/internal/services/knowledge"
	"github.com/rapidaai/config"
	document_client "github.com/rapidaai/pkg/clients/document"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
	storage_files "github.com/rapidaai/pkg/storages/file-storage"
	knowledge_api "github.com/rapidaai/protos"
)

type knowledgeApi struct {
	cfg                      *config.AppConfig
	logger                   commons.Logger
	postgres                 connectors.PostgresConnector
	redis                    connectors.RedisConnector
	knowledgeService         internal_services.KnowledgeService
	indexerServiceClient     document_client.IndexerServiceClient
	knowledgeDocumentService internal_services.KnowledgeDocumentService
}

type knowledgeGrpcApi struct {
	knowledgeApi
}

func NewKnowledgeGRPCApi(config *config.AppConfig, logger commons.Logger,
	postgres connectors.PostgresConnector,
	redis connectors.RedisConnector,
	opensearch connectors.OpenSearchConnector,
) knowledge_api.KnowledgeServiceServer {
	return &knowledgeGrpcApi{
		knowledgeApi{
			cfg:                      config,
			logger:                   logger,
			postgres:                 postgres,
			redis:                    redis,
			knowledgeService:         internal_knowledge_service.NewKnowledgeService(config, logger, postgres, storage_files.NewStorage(config.AssetStoreConfig, logger)),
			knowledgeDocumentService: internal_knowledge_service.NewKnowledgeDocumentService(config, logger, postgres, opensearch),
			indexerServiceClient:     document_client.NewIndexerServiceClient(config, logger, redis),
		},
	}
}
