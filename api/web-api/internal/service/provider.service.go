package internal_service

import (
	"context"

	internal_entity "github.com/lexatic/web-backend/api/web-api/internal/entity"
	provider_api "github.com/lexatic/web-backend/protos/lexatic-backend"
)

type ProviderService interface {
	GetModel(c context.Context, modelId uint64) (*internal_entity.ProviderModel, error)
	GetAllModel(c context.Context, criterias []*provider_api.Criteria) ([]*internal_entity.ProviderModel, error)

	GetAllModelProvider(c context.Context, criterias []*provider_api.Criteria) ([]*internal_entity.Provider, error)
	GetAllToolProvider(ctx context.Context, criterias []*provider_api.Criteria) ([]*internal_entity.ToolProvider, error)
}
