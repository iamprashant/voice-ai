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
	"sync"

	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

// GetAllAssistantProviderModel implements protos.AssistantServiceServer.
func (assistantApi *assistantGrpcApi) GetAllAssistantProvider(ctx context.Context, gaep *protos.GetAllAssistantProviderRequest) (*protos.GetAllAssistantProviderResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated || !iAuth.HasProject() {
		assistantApi.logger.Errorf("unauthenticated request for invoke")
		return utils.Error[protos.GetAllAssistantProviderResponse](
			errors.New("unauthenticated request for GetAllAssistantProvider"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)
	}
	combinedProviders := make([]*protos.GetAllAssistantProviderResponse_AssistantProvider, 0)
	var wg sync.WaitGroup
	wg.Add(1)
	// provider models
	utils.Go(ctx, func() {
		defer wg.Done()
		_, pModels, err := assistantApi.
			assistantService.
			GetAllAssistantProviderModel(ctx,
				iAuth,
				gaep.GetAssistantId(),
				gaep.GetCriterias(),
				gaep.GetPaginate())
		if err != nil {
			return
		}
		out := []*protos.AssistantProviderModel{}
		err = utils.Cast(pModels, &out)
		if err != nil {
			assistantApi.logger.Errorf("unable to cast assistant provider model %v", err)
		}
		for _, v := range out {
			combinedProviders = append(combinedProviders, &protos.GetAllAssistantProviderResponse_AssistantProvider{
				AssistantProvider: &protos.GetAllAssistantProviderResponse_AssistantProvider_AssistantProviderModel{
					AssistantProviderModel: v,
				},
			})
		}

	})

	wg.Add(1)
	// agentKit
	utils.Go(ctx, func() {
		defer wg.Done()
		_, pModels, err := assistantApi.
			assistantService.
			GetAllAssistantProviderAgentkit(ctx,
				iAuth,
				gaep.GetAssistantId(),
				gaep.GetCriterias(),
				gaep.GetPaginate())
		if err != nil {
			return
		}
		out := []*protos.AssistantProviderAgentkit{}
		err = utils.Cast(pModels, &out)
		if err != nil {
			assistantApi.logger.Errorf("unable to cast assistant provider agentkit %v", err)
		}
		for _, v := range out {
			combinedProviders = append(combinedProviders, &protos.GetAllAssistantProviderResponse_AssistantProvider{
				AssistantProvider: &protos.GetAllAssistantProviderResponse_AssistantProvider_AssistantProviderAgentkit{
					AssistantProviderAgentkit: v,
				},
			})
		}

	})

	wg.Add(1)
	utils.Go(ctx, func() {
		defer wg.Done()
		_, pModels, err := assistantApi.
			assistantService.
			GetAllAssistantProviderWebsocket(ctx,
				iAuth,
				gaep.GetAssistantId(),
				gaep.GetCriterias(),
				gaep.GetPaginate())
		if err != nil {
		}
		out := []*protos.AssistantProviderWebsocket{}
		err = utils.Cast(pModels, &out)
		if err != nil {
			assistantApi.logger.Errorf("unable to cast assistant provider websocket %v", err)
		}
		for _, v := range out {
			combinedProviders = append(combinedProviders, &protos.GetAllAssistantProviderResponse_AssistantProvider{
				AssistantProvider: &protos.GetAllAssistantProviderResponse_AssistantProvider_AssistantProviderWebsocket{
					AssistantProviderWebsocket: v,
				},
			})
		}

	})

	wg.Wait()
	return &protos.GetAllAssistantProviderResponse{
		Code:    200,
		Success: true,
		Paginated: &protos.Paginated{
			CurrentPage: gaep.GetPaginate().GetPage(),
			TotalItem:   uint32(len(combinedProviders)),
		},
		Data: combinedProviders,
	}, nil

}
