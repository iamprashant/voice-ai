// Rapida – Open Source Voice AI Orchestration Platform
// Copyright (C) 2023-2025 Prashant Srivastav <prashant@rapida.ai>
// Licensed under a modified GPL-2.0. See the LICENSE file for details.
package internal_services

import (
	"context"

	internal_gorm "github.com/rapidaai/api/integration-api/internal/entity"
	type_enums "github.com/rapidaai/pkg/types/enums"
	protos "github.com/rapidaai/protos"
)

type AuditService interface {
	Get(ctx context.Context, organizationId uint64, projectId uint64, auditId uint64) (*internal_gorm.ExternalAudit, error)
	GetAll(ctx context.Context, organizationId uint64, projectId uint64, paginate *protos.Paginate, opts []*protos.Criteria) (int64, []*internal_gorm.ExternalAudit, error)
	Create(ctx context.Context, requestId, organizationId, projectId, credentialId uint64, intName string, assetPrefix string, metrics []*protos.Metric, status type_enums.RecordState) (*internal_gorm.ExternalAudit, error)
	CreateMetadata(c context.Context, auditId uint64, metadata map[string]string) ([]*internal_gorm.ExternalAuditMetadata, error)
	UpdateMetadata(c context.Context, auditId uint64, metadata map[string]string) ([]*internal_gorm.ExternalAuditMetadata, error)
}
