// Copyright (c) Rapida
// Author: Prashant <prashant@rapida.ai>
//
// Licensed under the Rapida internal use license.
// This file is part of Rapida's proprietary software.
// Unauthorized copying, modification, or redistribution is strictly prohibited.
package assistant_api

import (
	"context"

	"github.com/rapidaai/pkg/exceptions"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	protos "github.com/rapidaai/protos"
)

func (assistantApi *assistantGrpcApi) UpdateAssistantTool(ctx context.Context, cawr *protos.UpdateAssistantToolRequest) (*protos.GetAssistantToolResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated || !iAuth.HasProject() {
		assistantApi.logger.Errorf("unauthenticated request for UpdateAssistantTool")
		return exceptions.AuthenticationError[protos.GetAssistantToolResponse]()
	}

	wl, err := assistantApi.assistantToolService.Update(
		ctx,
		iAuth,
		cawr.GetId(),
		cawr.GetAssistantId(),
		cawr.GetName(),
		cawr.GetDescription(),
		cawr.GetFields().AsMap(),
		cawr.GetExecutionMethod(),
		cawr.GetExecutionOptions())
	if err != nil {
		return exceptions.BadRequestError[protos.GetAssistantToolResponse](err.Error())
	}
	aAnalysis := &protos.AssistantTool{}
	err = utils.Cast(wl, aAnalysis)
	if err != nil {
		assistantApi.logger.Errorf("unable to cast the assistant tool to the response object")
	}
	return utils.Success[protos.GetAssistantToolResponse, *protos.AssistantTool](aAnalysis)
}
