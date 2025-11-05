package integration_api

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"time"

// 	config "github.com/rapidaai/config"
// 	callers "github.com/rapidaai/internal/callers"
// 	commons "github.com/rapidaai/pkg/commons"
// 	"github.com/rapidaai/pkg/connectors"
// 	"github.com/rapidaai/pkg/types"
// 	"github.com/rapidaai/pkg/utils"
// 	integration_api "github.com/rapidaai/protos"
// )

// type togetherAiIntegrationApi struct {
// 	integrationApi
// 	caller callers.Caller
// }

// type togetherAiIntegrationRPCApi struct {
// 	togetherAiIntegrationApi
// }

// type togetherAiIntegrationGRPCApi struct {
// 	togetherAiIntegrationApi
// }

// func NewTogetherAiRPCApi(config *config.AppConfig, logger commons.Logger, caller callers.Caller, postgres connectors.PostgresConnector) *togetherAiIntegrationRPCApi {
// 	return &togetherAiIntegrationRPCApi{
// 		togetherAiIntegrationApi{
// 			integrationApi: NewInegrationApi(config, logger, postgres),
// 			caller:         caller,
// 		},
// 	}
// }

// func NewTogetherAiGRPC(config *config.AppConfig, logger commons.Logger, caller callers.Caller, postgres connectors.PostgresConnector) integration_api.TogetherAiServiceServer {
// 	return &togetherAiIntegrationGRPCApi{
// 		togetherAiIntegrationApi{
// 			integrationApi: NewInegrationApi(config, logger, postgres),
// 			caller:         caller,
// 		},
// 	}
// }

// func (oiGRPC *togetherAiIntegrationGRPCApi) VerifyCredential(context.Context, *integration_api.VerifyCredentialRequest) (*integration_api.VerifyCredentialResponse, error) {
// 	return &integration_api.VerifyCredentialResponse{
// 		Code:    200,
// 		Success: true,
// 	}, nil
// }

// func (oiGRPC *togetherAiIntegrationGRPCApi) GenerateTextToImage(c context.Context, irRequest *integration_api.GenerateTextToImageRequest) (*integration_api.GenerateTextToImageResponse, error) {
// 	oiGRPC.logger.Debugf("request for image generate togetherAi with request %+v", irRequest)
// 	iAuth, isAuthenticated := types.GetClaimPrincipleGRPC[*types.ServiceScope](c)
// 	if !isAuthenticated || !iAuth.HasProject() {
// 		oiGRPC.logger.Errorf("unauthenticated request for invoke")
// 		return utils.Error[integration_api.GenerateTextToImageResponse](
// 			errors.New("unauthenticated request for text to image"),
// 			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
// 		)
// 	}

// 	headers := map[string]string{
// 		"Authorization": fmt.Sprintf("Bearer %s", irRequest.Credential.Value),
// 	}
// 	requestBody := map[string]interface{}{
// 		"prompt": irRequest.GetPrompt(),
// 		"model":
// 	}

// 	// construct parameter

// 	requestBody = oiGRPC.ConstructParameter(irRequest.GetModelParameters(), requestBody)
// 	adt := oiGRPC.PreHook(c, iAuth, irRequest.GetAdditionalData(), irRequest.GetCredential().GetId(), "TOGETHER_AI", requestBody)
// 	// only calculate api call
// 	start := time.Now()
// 	res, err := oiGRPC.caller.Call(c, "/v1/completions", "POST", headers, requestBody)
// 	timeTaken := int64(time.Since(start))

// 	if err == nil {
// 		oiGRPC.PostHook(c, iAuth, irRequest.GetAdditionalData(), irRequest.GetCredential().GetId(), adt, 200, timeTaken, res)
// 		return &integration_api.GenerateTextToImageResponse{
// 			Code:      200,
// 			Success:   true,
// 			Response:  res,
// 			RequestId: adt,
// 			TimeTaken: timeTaken,
// 		}, nil
// 	}

// 	oiGRPC.logger.Debugf("Exception occurred while calling completions %v", err)
// 	ex, ok := err.(callers.TogetherAiError)
// 	if ok {
// 		// can be used as defer
// 		errMessage := ex.Error()
// 		oiGRPC.PostHook(c, iAuth, irRequest.GetAdditionalData(), irRequest.GetCredential().GetId(), adt, int64(ex.Err.StatusCode), timeTaken, &errMessage)
// 		return &integration_api.GenerateTextToImageResponse{
// 			Code:    500,
// 			Success: false,
// 			Error: &integration_api.Error{
// 				ErrorCode:    uint64(ex.Err.StatusCode),
// 				ErrorMessage: errMessage,
// 				HumanMessage: ex.Err.Message,
// 			},
// 			RequestId: adt,
// 			TimeTaken: timeTaken,
// 		}, nil
// 	}

// 	return utils.Error[integration_api.GenerateTextToImageResponse](errors.New("illegal token while processing request"), "Illegal request, please try again")
// }

// func (oiGRPC *togetherAiIntegrationGRPCApi) Generate(c context.Context, irRequest *integration_api.GenerateRequest) (*integration_api.GenerateResponse, error) {
// 	oiGRPC.logger.Debugf("request for generate togetherAi with request %+v", irRequest)
// 	iAuth, isAuthenticated := types.GetClaimPrincipleGRPC[*types.ServiceScope](c)
// 	if !isAuthenticated || !iAuth.HasProject() {
// 		oiGRPC.logger.Errorf("unauthenticated request for invoke")
// 		return utils.Error[integration_api.GenerateResponse](
// 			errors.New("unauthenticated request for generate"),
// 			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
// 		)
// 	}

// 	headers := map[string]string{
// 		"Authorization": fmt.Sprintf("Bearer %s", irRequest.GetCredential().GetValue()),
// 	}

// 	requestBody := map[string]interface{}{
// 		"prompt": irRequest.GetPrompt(),
// 		"model":
// 	}
// 	requestBody = oiGRPC.ConstructParameter(irRequest.GetModelParameters(), requestBody)

// 	adt := oiGRPC.PreHook(c, iAuth, irRequest.GetAdditionalData(), irRequest.GetCredential().GetId(), "TOGETHER_AI", requestBody)
// 	// only calculate api call

// 	start := time.Now()
// 	res, err := oiGRPC.caller.Call(c, "/v1/completions", "POST", headers, requestBody)
// 	timeTaken := int64(time.Since(start))

// 	if err == nil {
// 		oiGRPC.PostHook(c, iAuth, irRequest.GetAdditionalData(), irRequest.GetCredential().GetId(), adt, 200, timeTaken, res)
// 		return &integration_api.GenerateResponse{
// 			Code:      200,
// 			Success:   true,
// 			Response:  res,
// 			RequestId: adt,
// 			TimeTaken: timeTaken,
// 		}, nil
// 	}

// 	oiGRPC.logger.Debugf("Exception occurred while calling completions %v", err)
// 	ex, ok := err.(callers.TogetherAiError)
// 	if ok {
// 		// can be used as defer
// 		errMessage := ex.Error()
// 		oiGRPC.PostHook(c, iAuth, irRequest.GetAdditionalData(), irRequest.GetCredential().GetId(), adt, 400, timeTaken, &errMessage)
// 		return &integration_api.GenerateResponse{
// 			Code:    500,
// 			Success: false,
// 			Error: &integration_api.Error{
// 				ErrorCode:    uint64(ex.Err.StatusCode),
// 				ErrorMessage: errMessage,
// 				HumanMessage: ex.Err.Message,
// 			},
// 			RequestId: adt,
// 			TimeTaken: timeTaken,
// 		}, nil
// 	}

// 	return utils.Error[integration_api.GenerateResponse](errors.New("illegal token while processing request"), "Illegal request, please try again")
// }

// func (oiGRPC *togetherAiIntegrationGRPCApi) Chat(c context.Context, irRequest *integration_api.ChatRequest) (*integration_api.ChatResponse, error) {
// 	oiGRPC.logger.Debugf("request for chat togetherAi with request %+v", irRequest)

// 	iAuth, isAuthenticated := types.GetClaimPrincipleGRPC[*types.ServiceScope](c)
// 	if !isAuthenticated || !iAuth.HasProject() {
// 		oiGRPC.logger.Errorf("unauthenticated request for invoke")
// 		return utils.Error[integration_api.ChatResponse](
// 			errors.New("unauthenticated request for chat"),
// 			"Please provider valid service credentials to perfom invoke, read docs @ docs.rapida.ai",
// 		)
// 	}

// 	headers := map[string]string{
// 		"Authorization": fmt.Sprintf("Bearer %s", irRequest.GetCredential().GetValue()),
// 	}

// 	requestBody := map[string]interface{}{
// 		"model":
// 		"messages": irRequest.Conversations,
// 	}
// 	requestBody = oiGRPC.ConstructParameter(irRequest.GetModelParameters(), requestBody)

// 	adt := oiGRPC.PreHook(c, iAuth, irRequest.GetAdditionalData(), irRequest.GetCredential().GetId(), "TOGETHER_AI", requestBody)
// 	// only calculate api call
// 	start := time.Now()
// 	res, err := oiGRPC.caller.Call(c, "/v1/chat/completions", "POST", headers, requestBody)
// 	timeTaken := int64(time.Since(start))

// 	if err == nil {
// 		oiGRPC.PostHook(c, iAuth, irRequest.GetAdditionalData(), irRequest.GetCredential().GetId(), adt, 200, timeTaken, res)
// 		return &integration_api.ChatResponse{
// 			Code:      200,
// 			Success:   true,
// 			Response:  res,
// 			RequestId: adt,
// 			TimeTaken: timeTaken,
// 		}, nil
// 	}

// 	oiGRPC.logger.Debugf("Exception occurred while calling completions %v", err)
// 	ex, ok := err.(callers.TogetherAiError)
// 	if ok {
// 		// can be used as defer
// 		errMessage := ex.Error()
// 		oiGRPC.PostHook(c, iAuth, irRequest.GetAdditionalData(), irRequest.GetCredential().GetId(), adt, 400, timeTaken, &errMessage)
// 		return &integration_api.ChatResponse{
// 			Code:    500,
// 			Success: false,
// 			Error: &integration_api.Error{
// 				ErrorCode:    uint64(ex.Err.StatusCode),
// 				ErrorMessage: errMessage,
// 				HumanMessage: ex.Err.Message,
// 			},
// 			RequestId: adt,
// 			TimeTaken: timeTaken,
// 		}, nil
// 	}

// 	return utils.Error[integration_api.ChatResponse](errors.New("illegal token while processing request"), "Illegal request, please try again")
// }
