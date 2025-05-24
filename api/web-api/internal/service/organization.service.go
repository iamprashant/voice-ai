package internal_service

import (
	"context"

	internal_entity "github.com/lexatic/web-backend/api/web-api/internal/entity"
	"github.com/lexatic/web-backend/pkg/types"
)

type OrganizationService interface {
	Create(ctx context.Context, auth types.Principle, name string, size string, industry string) (*internal_entity.Organization, error)
	Get(ctx context.Context, organizationId uint64) (*internal_entity.Organization, error)
	Update(ctx context.Context, auth types.Principle, organizationId uint64, name *string, industry *string, email *string) (*internal_entity.Organization, error)
}
