package endpoint_api

import (
	"context"
	"errors"
	"time"

	internal_services "github.com/rapidaai/internal/services"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	endpoint_grpc_api "github.com/rapidaai/protos"
)

func (endpointGRPCApi *endpointGRPCApi) GetEndpoint(ctx context.Context, cepm *endpoint_grpc_api.GetEndpointRequest) (*endpoint_grpc_api.GetEndpointResponse, error) {
	start := time.Now()
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated || !iAuth.HasProject() {
		endpointGRPCApi.logger.Errorf("unauthenticated request for invoke")
		return utils.Error[endpoint_grpc_api.GetEndpointResponse](
			errors.New("unauthenticated request for get endpoint"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)
	}

	ep, err := endpointGRPCApi.endpointService.Get(ctx, iAuth, cepm.GetId(), cepm.EndpointProviderModelId, internal_services.NewDefaultGetEndpointOption())
	if err != nil {
		return utils.Error[endpoint_grpc_api.GetEndpointResponse](
			err,
			"Unable to get the endpoint for given endpoint id.",
		)
	}

	endpointGRPCApi.logger.Benchmark("endpointGRPCApi.GetEndpoint", time.Since(start))
	out := &endpoint_grpc_api.Endpoint{}
	err = utils.Cast(ep, out)
	if err != nil {
		endpointGRPCApi.logger.Errorf("unable to cast endpoint provider model %v", err)
	}

	// analytics := endpointGRPCApi.auditService.GetAggregatedEndpointAnalytics(ctx, ep.Id, 30)
	// out.EndpointAnalytics = analytics
	endpointGRPCApi.logger.Benchmark("endpointGRPCApi.GetEndpoint.EndpointAnalytics", time.Since(start))
	return utils.Success[endpoint_grpc_api.GetEndpointResponse](out)
}

func (endpointGRPCApi *endpointGRPCApi) GetAllEndpoint(ctx context.Context, cepm *endpoint_grpc_api.GetAllEndpointRequest) (*endpoint_grpc_api.GetAllEndpointResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated || !iAuth.HasProject() {
		endpointGRPCApi.logger.Errorf("unauthenticated request for GetAllEndpoint")
		return utils.Error[endpoint_grpc_api.GetAllEndpointResponse](
			errors.New("unauthenticated request for get allendpoint"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)
	}
	cnt, endpoints, err := endpointGRPCApi.endpointService.GetAll(ctx, iAuth,
		cepm.GetCriterias(),
		cepm.GetPaginate())
	if err != nil {
		return utils.Error[endpoint_grpc_api.GetAllEndpointResponse](
			err,
			"Unable to get all the endpoint.",
		)
	}
	out := []*endpoint_grpc_api.Endpoint{}
	err = utils.Cast(endpoints, &out)
	if err != nil {
		endpointGRPCApi.logger.Errorf("unable to cast endpoint provider model %v", err)
	}

	for _, e := range out {
		analytics := endpointGRPCApi.endpointLogService.GetAggregatedEndpointAnalytics(ctx, iAuth, e.Id)
		e.EndpointAnalytics = analytics
	}
	return utils.PaginatedSuccess[endpoint_grpc_api.GetAllEndpointResponse](
		uint32(cnt),
		cepm.GetPaginate().GetPage(),
		out)
}

func (endpointGRPCApi *endpointGRPCApi) GetAllEndpointProviderModel(ctx context.Context, gaep *endpoint_grpc_api.GetAllEndpointProviderModelRequest) (*endpoint_grpc_api.GetAllEndpointProviderModelResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated || !iAuth.HasProject() {
		endpointGRPCApi.logger.Errorf("unauthenticated request for invoke")
		return utils.Error[endpoint_grpc_api.GetAllEndpointProviderModelResponse](
			errors.New("unauthenticated request for getallendpointprovidermodel"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)
	}

	cnt, epms, err := endpointGRPCApi.endpointService.GetAllEndpointProviderModel(ctx,
		iAuth,
		gaep.GetEndpointId(),
		gaep.GetCriterias(),
		gaep.GetPaginate())
	if err != nil {
		return utils.Error[endpoint_grpc_api.GetAllEndpointProviderModelResponse](
			err,
			"Unable to get all the endpoint provider model.",
		)
	}
	out := []*endpoint_grpc_api.EndpointProviderModel{}
	err = utils.Cast(epms, &out)
	if err != nil {
		endpointGRPCApi.logger.Errorf("unable to cast endpoint provider model %v", err)
	}

	return utils.PaginatedSuccess[endpoint_grpc_api.GetAllEndpointProviderModelResponse, []*endpoint_grpc_api.EndpointProviderModel](
		uint32(cnt),
		gaep.GetPaginate().GetPage(),
		out)

}

func (endpointGRPCApi *endpointGRPCApi) GetAllEndpointLog(ctx context.Context, gaep *endpoint_grpc_api.GetAllEndpointLogRequest) (*endpoint_grpc_api.GetAllEndpointLogResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated || !iAuth.HasProject() {
		endpointGRPCApi.logger.Errorf("unauthenticated request for invoke")
		return utils.Error[endpoint_grpc_api.GetAllEndpointLogResponse](
			errors.New("unauthenticated request for getallendpointprovidermodel"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)
	}
	cnt, epms, err := endpointGRPCApi.endpointLogService.GetAllEndpointLog(ctx,
		iAuth,
		gaep.GetEndpointId(),
		gaep.GetCriterias(),
		gaep.GetPaginate())
	if err != nil {
		return utils.Error[endpoint_grpc_api.GetAllEndpointLogResponse](
			err,
			"Unable to get all the endpoint provider model.",
		)
	}
	out := []*endpoint_grpc_api.EndpointLog{}
	err = utils.Cast(epms, &out)
	if err != nil {
		endpointGRPCApi.logger.Errorf("unable to cast endpoint provider model %v", err)
	}

	return utils.PaginatedSuccess[endpoint_grpc_api.GetAllEndpointLogResponse, []*endpoint_grpc_api.EndpointLog](
		uint32(cnt),
		gaep.GetPaginate().GetPage(),
		out)

}

func (endpointGRPCApi *endpointGRPCApi) GetEndpointLog(ctx context.Context, cepm *endpoint_grpc_api.GetEndpointLogRequest) (*endpoint_grpc_api.GetEndpointLogResponse, error) {
	start := time.Now()
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated || !iAuth.HasProject() {
		endpointGRPCApi.logger.Errorf("unauthenticated request for invoke")
		return utils.Error[endpoint_grpc_api.GetEndpointLogResponse](
			errors.New("unauthenticated request for get endpoint"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)
	}

	ep, err := endpointGRPCApi.endpointLogService.GetEndpointLog(ctx,
		iAuth,
		cepm.GetId(),
		cepm.GetEndpointId())
	if err != nil {
		return utils.Error[endpoint_grpc_api.GetEndpointLogResponse](
			err,
			"Unable to get the endpoint log for given id.",
		)
	}

	endpointGRPCApi.logger.Benchmark("endpointGRPCApi.GetEndpoint", time.Since(start))
	out := &endpoint_grpc_api.EndpointLog{}
	err = utils.Cast(ep, out)
	if err != nil {
		endpointGRPCApi.logger.Errorf("unable to cast endpoint provider model %v", err)
	}
	endpointGRPCApi.logger.Benchmark("endpointGRPCApi.GetEndpoint.EndpointAnalytics", time.Since(start))
	return utils.Success[endpoint_grpc_api.GetEndpointLogResponse](out)
}
