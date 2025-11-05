package endpoint_api

import (
	"context"
	"errors"

	internal_gorm "github.com/rapidaai/internal/gorm"
	internal_services "github.com/rapidaai/internal/services"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	endpoint_grpc_api "github.com/rapidaai/protos"
	"google.golang.org/protobuf/encoding/protojson"
)

func (endpointGRPCApi *endpointGRPCApi) CreateEndpoint(ctx context.Context, cer *endpoint_grpc_api.CreateEndpointRequest) (*endpoint_grpc_api.CreateEndpointResponse, error) {
	endpointGRPCApi.logger.Debugf("Create endpoint request %v, %v", cer, ctx)
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated {
		endpointGRPCApi.logger.Errorf("unauthenticated request for invoke")
		return utils.Error[endpoint_grpc_api.CreateEndpointResponse](
			errors.New("unauthenticated request for invoke"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)
	}
	endpoint, err := endpointGRPCApi.endpointService.CreateEndpoint(
		ctx,
		iAuth,
		cer.GetEndpointAttribute().GetName(),
		&cer.GetEndpointAttribute().Description,
		&cer.GetEndpointAttribute().Visibility,
		&cer.GetEndpointAttribute().Source,
		&cer.GetEndpointAttribute().SourceIdentifier)
	if err != nil {
		return utils.Error[endpoint_grpc_api.CreateEndpointResponse](
			err,
			"Unable to create endpoint, please try again later",
		)
	}
	epModel, err := endpointGRPCApi.createEndpointProviderModel(ctx, iAuth, endpoint, cer.GetEndpointProviderModelAttribute())
	if err != nil {
		return utils.Error[endpoint_grpc_api.CreateEndpointResponse](
			err,
			"Unable to create endpoint provider model, please try again later",
		)
	}

	_, err = endpointGRPCApi.endpointService.AttachProviderModelToEndpoint(ctx, iAuth, epModel.Id, endpoint.Id)
	if err != nil {
		return utils.Error[endpoint_grpc_api.CreateEndpointResponse](
			err,
			"Unable to attach endpoint provider model, please try again later",
		)
	}

	if cer.GetRetryConfiguration() != nil {
		_, err = endpointGRPCApi.endpointService.ConfigureEndpointRetry(ctx,
			iAuth,
			endpoint.Id,
			internal_gorm.Retry(cer.GetRetryConfiguration().GetRetryType()),
			cer.GetRetryConfiguration().GetMaxAttempts(),
			cer.GetRetryConfiguration().GetDelaySeconds(),
			cer.GetRetryConfiguration().GetExponentialBackoff(),
			cer.GetRetryConfiguration().GetRetryables(),
		)
		if err != nil {
			return utils.Error[endpoint_grpc_api.CreateEndpointResponse](
				errors.New("unauthenticated request for CreateEndpointProviderModel"),
				"Unable to configure endpoint retry, please try again later",
			)

		}
	}
	if cer.GetCacheConfiguration() != nil {
		_, err = endpointGRPCApi.endpointService.ConfigureEndpointCaching(ctx,
			iAuth,
			endpoint.Id,
			internal_gorm.Cache(cer.GetCacheConfiguration().GetCacheType()),
			cer.GetCacheConfiguration().GetExpiryInterval(),
			cer.GetCacheConfiguration().GetMatchThreshold())
		if err != nil {
			return utils.Error[endpoint_grpc_api.CreateEndpointResponse](
				errors.New("unauthenticated request for CreateEndpointProviderModel"),
				"Unable to configure endpoint retry, please try again later",
			)

		}
	}
	_, err = endpointGRPCApi.endpointService.CreateOrUpdateEndpointTag(ctx, iAuth, endpoint.Id, cer.GetTags())
	if err != nil {
		return utils.Error[endpoint_grpc_api.CreateEndpointResponse](
			err,
			"Unable to create endpoint tags, please try again.",
		)
	}

	endpoint.EndpointProviderModel = epModel
	out := &endpoint_grpc_api.Endpoint{}
	err = utils.Cast(endpoint, out)
	if err != nil {
		endpointGRPCApi.logger.Errorf("unable to cast the endpoint provider model to the response object")
	}

	// calling to index the endpoint
	// endpointGRPCApi.endpointService.IndexEndpoint(ctx, iAuth, endpoint.Id)
	return utils.Success[endpoint_grpc_api.CreateEndpointResponse, *endpoint_grpc_api.Endpoint](out)
}

func (endpointGRPCApi *endpointGRPCApi) CreateEndpointProviderModel(ctx context.Context, iRequest *endpoint_grpc_api.CreateEndpointProviderModelRequest) (*endpoint_grpc_api.CreateEndpointProviderModelResponse, error) {
	endpointGRPCApi.logger.Debugf("create Endpoint provider model request %v %v", iRequest, ctx)

	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated {
		endpointGRPCApi.logger.Errorf("unauthenticated request for invoke")
		return utils.Error[endpoint_grpc_api.CreateEndpointProviderModelResponse](
			errors.New("unauthenticated request for CreateEndpointProviderModel"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)

	}

	ep, err := endpointGRPCApi.endpointService.Get(ctx, iAuth, iRequest.GetEndpointId(), nil, internal_services.NewDefaultGetEndpointOption())
	if err != nil {
		return utils.Error[endpoint_grpc_api.CreateEndpointProviderModelResponse](
			err,
			"Unable to create endpoint version, please try again later",
		)
	}

	epm, err := endpointGRPCApi.createEndpointProviderModel(ctx, iAuth, ep, iRequest.GetEndpointProviderModelAttribute())
	if err != nil {
		return utils.Error[endpoint_grpc_api.CreateEndpointProviderModelResponse](
			err,
			"Unable to create endpoint version, please try again later",
		)
	}
	out := &endpoint_grpc_api.EndpointProviderModel{}
	err = utils.Cast(epm, out)
	if err != nil {
		endpointGRPCApi.logger.Errorf("unable to cast the endpoint provider model to the response object")
	}
	return utils.Success[endpoint_grpc_api.CreateEndpointProviderModelResponse, *endpoint_grpc_api.EndpointProviderModel](out)
}

func (endpointGRPCApi *endpointGRPCApi) createEndpointProviderModel(ctx context.Context,
	iAuth types.SimplePrinciple,
	endpoint *internal_gorm.Endpoint,
	ea *endpoint_grpc_api.EndpointProviderModelAttribute,
) (*internal_gorm.EndpointProviderModel, error) {
	endpointGRPCApi.logger.Debugf("creating endpoint provider model with provider model attributes %+v", ea)
	prompt := protojson.Format(ea.GetChatCompletePrompt())
	epm, err := endpointGRPCApi.endpointService.CreateEndpointProviderModel(ctx,
		iAuth,
		endpoint.Id,
		ea.GetDescription(),
		ea.GetModelProviderId(),
		ea.GetModelProviderName(),
		prompt,
		ea.GetEndpointModelOptions(),
	)
	if err != nil {
		endpointGRPCApi.logger.Errorf("unable to create endpoint provider model wuth error %+v", err)
		return nil, err
	}
	return epm, nil
}

func (endpointGRPCApi *endpointGRPCApi) CreateEndpointCacheConfiguration(ctx context.Context, eRequest *endpoint_grpc_api.CreateEndpointCacheConfigurationRequest) (*endpoint_grpc_api.CreateEndpointCacheConfigurationResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated {
		endpointGRPCApi.logger.Errorf("unauthenticated request for invoke")
		return utils.Error[endpoint_grpc_api.CreateEndpointCacheConfigurationResponse](
			errors.New("unauthenticated request for CreateEndpointProviderModel"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)

	}

	cec, err := endpointGRPCApi.endpointService.ConfigureEndpointCaching(ctx,
		iAuth,
		eRequest.GetEndpointId(),
		internal_gorm.Cache(eRequest.GetData().GetCacheType()),
		eRequest.GetData().GetExpiryInterval(),
		eRequest.GetData().GetMatchThreshold())
	if err != nil {
		return utils.Error[endpoint_grpc_api.CreateEndpointCacheConfigurationResponse](
			err,
			"Unable to configure endpoint caching, please try again later",
		)

	}
	out := &endpoint_grpc_api.EndpointCacheConfiguration{}
	err = utils.Cast(cec, out)
	if err != nil {
		endpointGRPCApi.logger.Errorf("unable to cast the endpoint cache configuration to the response object")
	}
	return utils.Success[endpoint_grpc_api.CreateEndpointCacheConfigurationResponse, *endpoint_grpc_api.EndpointCacheConfiguration](out)

}
func (endpointGRPCApi *endpointGRPCApi) CreateEndpointRetryConfiguration(ctx context.Context, eRequest *endpoint_grpc_api.CreateEndpointRetryConfigurationRequest) (*endpoint_grpc_api.CreateEndpointRetryConfigurationResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated {
		endpointGRPCApi.logger.Errorf("unauthenticated request for invoke")
		return utils.Error[endpoint_grpc_api.CreateEndpointRetryConfigurationResponse](
			errors.New("unauthenticated request for CreateEndpointProviderModel"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)

	}
	erc, err := endpointGRPCApi.endpointService.ConfigureEndpointRetry(ctx,
		iAuth,
		eRequest.GetEndpointId(),
		internal_gorm.Retry(eRequest.GetData().GetRetryType()),
		eRequest.GetData().GetMaxAttempts(),
		eRequest.GetData().GetDelaySeconds(),
		eRequest.GetData().GetExponentialBackoff(),
		eRequest.GetData().GetRetryables(),
	)
	if err != nil {
		return utils.Error[endpoint_grpc_api.CreateEndpointRetryConfigurationResponse](
			err,
			"Unable to configure endpoint retry, please try again later",
		)

	}
	out := &endpoint_grpc_api.EndpointRetryConfiguration{}
	err = utils.Cast(erc, out)
	if err != nil {
		endpointGRPCApi.logger.Errorf("unable to cast the endpoint retry configuration to the response object")
	}

	return utils.Success[endpoint_grpc_api.CreateEndpointRetryConfigurationResponse, *endpoint_grpc_api.EndpointRetryConfiguration](out)
}
func (endpointGRPCApi *endpointGRPCApi) CreateEndpointTag(ctx context.Context, eRequest *endpoint_grpc_api.CreateEndpointTagRequest) (*endpoint_grpc_api.GetEndpointResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated {
		endpointGRPCApi.logger.Errorf("unauthenticated request for invoke")
		return utils.Error[endpoint_grpc_api.GetEndpointResponse](
			errors.New("unauthenticated request for CreateEndpointProviderModel"),
			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
		)
	}
	_, err := endpointGRPCApi.endpointService.CreateOrUpdateEndpointTag(ctx, iAuth, eRequest.GetEndpointId(), eRequest.GetTags())
	if err != nil {
		return utils.Error[endpoint_grpc_api.GetEndpointResponse](
			err,
			"Unable to create endpoint tags for endpoint",
		)
	}
	// // calling to index the endpoint
	// endpointGRPCApi.endpointService.IndexEndpoint(ctx, iAuth, eRequest.GetEndpointId())
	ep, err := endpointGRPCApi.endpointService.Get(ctx, iAuth, eRequest.GetEndpointId(), nil, internal_services.NewDefaultGetEndpointOption())
	if err != nil {
		return utils.Error[endpoint_grpc_api.GetEndpointResponse](
			err,
			"Unable to get the endpoint for given endpoint id.",
		)
	}

	out := &endpoint_grpc_api.Endpoint{}
	err = utils.Cast(ep, out)
	if err != nil {
		endpointGRPCApi.logger.Errorf("unable to cast endpoint %v", err)
	}

	return utils.Success[endpoint_grpc_api.GetEndpointResponse](out)
}
