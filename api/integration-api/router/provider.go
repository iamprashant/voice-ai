// Rapida – Open Source Voice AI Orchestration Platform
// Copyright (C) 2023-2025 Prashant Srivastav <prashant@rapida.ai>
// Licensed under a modified GPL-2.0. See the LICENSE file for details.
package integration_routers

import (
	integrationApi "github.com/rapidaai/api/integration-api/api"
	"github.com/rapidaai/api/integration-api/config"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
	"github.com/rapidaai/protos"
	"google.golang.org/grpc"
)

// all the provider routes
func ProviderApiRoute(Cfg *config.IntegrationConfig, S *grpc.Server, Logger commons.Logger, Postgres connectors.PostgresConnector) {
	protos.RegisterCohereServiceServer(S, integrationApi.NewCohereGRPC(Cfg, Logger, Postgres))
	protos.RegisterOpenAiServiceServer(S, integrationApi.NewOpenAiGRPC(Cfg, Logger, Postgres))
	protos.RegisterGeminiServiceServer(S, integrationApi.NewGeminiGRPC(Cfg, Logger, Postgres))
	protos.RegisterAzureServiceServer(S, integrationApi.NewAzureGRPC(Cfg, Logger, Postgres))
	protos.RegisterAnthropicServiceServer(S, integrationApi.NewAnthropicGRPC(Cfg, Logger, Postgres))
	protos.RegisterVoyageAiServiceServer(S, integrationApi.NewVoyageAiGRPC(Cfg, Logger, Postgres))
	protos.RegisterHuggingfaceServiceServer(S, integrationApi.NewHuggingfaceGRPC(Cfg, Logger, Postgres))
	protos.RegisterMistralServiceServer(S, integrationApi.NewMistralGRPC(Cfg, Logger, Postgres))
	protos.RegisterReplicateServiceServer(S, integrationApi.NewReplicateGRPC(Cfg, Logger, Postgres))
	protos.RegisterVertexAiServiceServer(S, integrationApi.NewVertexaiGRPC(Cfg, Logger, Postgres))
}

// audit logging api route
func AuditLoggingApiRoute(
	Cfg *config.IntegrationConfig,
	S *grpc.Server,
	Logger commons.Logger,
	Postgres connectors.PostgresConnector,
) {
	protos.RegisterAuditLoggingServiceServer(S, integrationApi.NewAuditLoggingGRPC(Cfg, Logger, Postgres))
}
