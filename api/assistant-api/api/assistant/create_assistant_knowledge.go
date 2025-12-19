// Copyright (c) Rapida
// Author: Prashant <prashant@rapida.ai>
//
// Licensed under the Rapida internal use license.
// This file is part of Rapida's proprietary software.
// Unauthorized copying, modification, or redistribution is strictly prohibited.
package assistant_api

import (
	"context"
	"errors"

	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	"github.com/rapidaai/pkg/exceptions"
	gorm_types "github.com/rapidaai/pkg/models/gorm/types"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	assistant_api "github.com/rapidaai/protos"
)

// CreateAssistantKnowledgeConfiguration implements assistant_api.AssistantServiceServer.
func (assistantApi *assistantGrpcApi) CreateAssistantKnowledge(ctx context.Context, cepm *assistant_api.CreateAssistantKnowledgeRequest) (*assistant_api.GetAssistantKnowledgeResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated || !iAuth.HasProject() {
		assistantApi.logger.Errorf("unauthenticated request for invoke")
		return utils.Error[assistant_api.GetAssistantKnowledgeResponse](
			errors.New("unauthenticated request for get assistant"),
			"Please provider valid service credentials to perform CreateAssistantKnowledge, read docs @ docs.rapida.ai",
		)
	}
	aK, err := assistantApi.assistantKnowledgeService.Create(
		ctx,
		iAuth,
		cepm.GetAssistantId(),
		cepm.GetKnowledgeId(),
		gorm_types.RetrievalMethod(cepm.GetRetrievalMethod()),
		cepm.GetRerankerEnable(),
		cepm.GetScoreThreshold(),
		cepm.GetTopK(),
		&cepm.RerankerModelProviderId,
		&cepm.RerankerModelProviderName,
		cepm.GetAssistantKnowledgeRerankerOptions(),
	)
	if err != nil {
		return exceptions.BadRequestError[assistant_api.GetAssistantKnowledgeResponse](err.Error())
	}

	out := &assistant_api.AssistantKnowledge{}
	err = utils.Cast(aK, out)
	if err != nil {
		assistantApi.logger.Errorf("unable to cast assistant knowledge %v", err)
	}
	return utils.Success[assistant_api.GetAssistantKnowledgeResponse, *assistant_api.AssistantKnowledge](out)
}

func (assistantApi *assistantGrpcApi) createAssistantKnowledge(
	ctx context.Context,
	iAuth types.SimplePrinciple,
	assistantId uint64, cepm *assistant_api.CreateAssistantKnowledgeRequest) (*internal_assistant_entity.AssistantKnowledge, error) {
	return assistantApi.assistantKnowledgeService.Create(
		ctx,
		iAuth,
		assistantId,
		cepm.GetKnowledgeId(),
		gorm_types.RetrievalMethod(cepm.GetRetrievalMethod()),
		cepm.GetRerankerEnable(),
		cepm.GetScoreThreshold(),
		cepm.GetTopK(),
		&cepm.RerankerModelProviderId,
		&cepm.RerankerModelProviderName,
		cepm.GetAssistantKnowledgeRerankerOptions(),
	)

}
