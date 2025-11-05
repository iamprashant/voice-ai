package integration_api

import (
	"context"

	"github.com/gin-gonic/gin"
	config "github.com/rapidaai/config"
	internal_callers "github.com/rapidaai/internal/callers"
	internal_bedrock_callers "github.com/rapidaai/internal/callers/bedrock"
	commons "github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/connectors"
	integration_api "github.com/rapidaai/protos"
)

type bedrockIntegrationApi struct {
	integrationApi
}

type bedrockIntegrationRPCApi struct {
	bedrockIntegrationApi
}

type bedrockIntegrationGRPCApi struct {
	bedrockIntegrationApi
}

// Chat implements lexatic_backend.BedrockServiceServer.

// Generate implements lexatic_backend.BedrockServiceServer.
func (bedrock *bedrockIntegrationGRPCApi) Chat(c context.Context, irRequest *integration_api.ChatRequest) (*integration_api.ChatResponse, error) {
	return bedrock.integrationApi.Chat(c, irRequest, "BEDROCK", internal_bedrock_callers.NewLargeLanguageCaller(bedrock.logger, irRequest.Credential))
}

func NewBedrockRPC(config *config.AppConfig, logger commons.Logger, postgres connectors.PostgresConnector) *bedrockIntegrationRPCApi {
	return &bedrockIntegrationRPCApi{
		bedrockIntegrationApi{
			integrationApi: NewInegrationApi(config, logger, postgres),
			// caller:         caller,
		},
	}
}

func NewBedrockGRPC(config *config.AppConfig, logger commons.Logger, postgres connectors.PostgresConnector) integration_api.BedrockServiceServer {
	return &bedrockIntegrationGRPCApi{
		bedrockIntegrationApi{
			integrationApi: NewInegrationApi(config, logger, postgres),
		},
	}
}

// all the rpc handler
func (oiRPC *bedrockIntegrationRPCApi) Generate(c *gin.Context) {
	oiRPC.logger.Debugf("Generate from rpc with gin context %v", c)
}

func (oiRPC *bedrockIntegrationRPCApi) Chat(c *gin.Context) {
	oiRPC.logger.Debugf("Chat from rpc with gin context %v", c)
}

// Embedding implements lexatic_backend.BedrockServiceServer.
func (bedrock *bedrockIntegrationGRPCApi) Embedding(c context.Context, irRequest *integration_api.EmbeddingRequest) (*integration_api.EmbeddingResponse, error) {
	return bedrock.integrationApi.Embedding(c, irRequest, "BEDROCK", internal_bedrock_callers.NewEmbeddingCaller(bedrock.logger, irRequest.Credential))
}

func (dgGRPC *bedrockIntegrationApi) VerifyCredential(c context.Context, irRequest *integration_api.VerifyCredentialRequest) (*integration_api.VerifyCredentialResponse, error) {
	bedrockCaller := internal_bedrock_callers.NewVerifyCredentialCaller(dgGRPC.logger, irRequest.Credential)
	st, err := bedrockCaller.CredentialVerifier(
		c,
		&internal_callers.CredentialVerifierOptions{},
	)
	if err != nil {
		return &integration_api.VerifyCredentialResponse{
			Code:         401,
			Success:      false,
			ErrorMessage: err.Error(),
		}, nil
	}
	return &integration_api.VerifyCredentialResponse{
		Code:     200,
		Success:  true,
		Response: st,
	}, nil
}
