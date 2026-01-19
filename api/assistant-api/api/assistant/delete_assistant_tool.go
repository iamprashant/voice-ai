// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package assistant_api

import (
	"context"
	"errors"

	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	assistant_api "github.com/rapidaai/protos"
)

func (assistantApi *assistantGrpcApi) DeleteAssistantTool(ctx context.Context, cer *assistant_api.DeleteAssistantToolRequest) (*assistant_api.GetAssistantToolResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated {
		assistantApi.logger.Errorf("unauthenticated request for DeleteAssistantToolRequest")
		return utils.Error[assistant_api.GetAssistantToolResponse](
			errors.New("unauthenticated request for DeleteAssistantToolRequest"),
			"Please provider valid service credentials to perfom DeleteAssistantToolRequest, read docs @ docs.rapida.ai",
		)
	}
	analysis, err := assistantApi.assistantToolService.Delete(ctx,
		iAuth,
		cer.GetId(),
		cer.GetAssistantId())
	if err != nil {
		return utils.Error[assistant_api.GetAssistantToolResponse](
			err,
			"Unable to update assistant analysis, please try again in sometime",
		)
	}
	out := &assistant_api.AssistantTool{}
	err = utils.Cast(analysis, out)
	if err != nil {
		assistantApi.logger.Errorf("unable to cast the assistant analysis to the response object")
	}
	return utils.Success[assistant_api.GetAssistantToolResponse, *assistant_api.AssistantTool](out)

}
