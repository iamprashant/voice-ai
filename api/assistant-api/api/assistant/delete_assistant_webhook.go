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

func (assistantApi *assistantGrpcApi) DeleteAssistantWebhook(ctx context.Context, cer *assistant_api.DeleteAssistantWebhookRequest) (*assistant_api.GetAssistantWebhookResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated {
		assistantApi.logger.Errorf("unauthenticated request for DeleteAssistantWebhookRequest")
		return utils.Error[assistant_api.GetAssistantWebhookResponse](
			errors.New("unauthenticated request for DeleteAssistantWebhookRequest"),
			"Please provider valid service credentials to perfom DeleteAssistantWebhookRequest, read docs @ docs.rapida.ai",
		)
	}
	analysis, err := assistantApi.assistantWebhookService.Delete(ctx,
		iAuth,
		cer.GetId(), cer.GetAssistantId())
	if err != nil {
		return utils.Error[assistant_api.GetAssistantWebhookResponse](
			err,
			"Unable to update assistant analysis, please try again in sometime",
		)
	}
	out := &assistant_api.AssistantWebhook{}
	err = utils.Cast(analysis, out)
	if err != nil {
		assistantApi.logger.Errorf("unable to cast the assistant analysis to the response object")
	}
	return utils.Success[assistant_api.GetAssistantWebhookResponse, *assistant_api.AssistantWebhook](out)

}
