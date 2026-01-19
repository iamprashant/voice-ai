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

func (assistantApi *assistantGrpcApi) DeleteAssistant(ctx context.Context, cer *assistant_api.DeleteAssistantRequest) (*assistant_api.GetAssistantResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated {
		assistantApi.logger.Errorf("unauthenticated request for UpdateAssistantDetail")
		return utils.Error[assistant_api.GetAssistantResponse](
			errors.New("unauthenticated request for UpdateAssistantDetail"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)
	}
	assistant, err := assistantApi.assistantService.DeleteAssistant(ctx,
		iAuth,
		cer.GetId())
	if err != nil {
		return utils.Error[assistant_api.GetAssistantResponse](
			err,
			"Unable to update assistant, please try again in sometime",
		)
	}
	out := &assistant_api.Assistant{}
	err = utils.Cast(assistant, out)
	if err != nil {
		assistantApi.logger.Errorf("unable to cast the assistant provider model to the response object")
	}
	return utils.Success[assistant_api.GetAssistantResponse, *assistant_api.Assistant](out)

}
