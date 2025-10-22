package web_handler

import (
	"context"
	"errors"

	assistant_client "github.com/lexatic/web-backend/pkg/clients/workflow"
	"github.com/lexatic/web-backend/pkg/exceptions"
	"github.com/lexatic/web-backend/pkg/utils"
	web_api "github.com/lexatic/web-backend/protos/lexatic-backend"

	config "github.com/lexatic/web-backend/config"
	commons "github.com/lexatic/web-backend/pkg/commons"
	"github.com/lexatic/web-backend/pkg/connectors"
	"github.com/lexatic/web-backend/pkg/types"
)

type webAssistantApi struct {
	WebApi
	cfg             *config.AppConfig
	logger          commons.Logger
	postgres        connectors.PostgresConnector
	redis           connectors.RedisConnector
	assistantClient assistant_client.AssistantServiceClient
}

type webAssistantGRPCApi struct {
	webAssistantApi
}

func NewAssistantGRPC(config *config.AppConfig, logger commons.Logger, postgres connectors.PostgresConnector, redis connectors.RedisConnector) web_api.AssistantServiceServer {
	return &webAssistantGRPCApi{
		webAssistantApi{
			WebApi:          NewWebApi(config, logger, postgres, redis),
			cfg:             config,
			logger:          logger,
			postgres:        postgres,
			redis:           redis,
			assistantClient: assistant_client.NewAssistantServiceClientGRPC(config, logger, redis),
		},
	}
}

func (assistant *webAssistantGRPCApi) GetAllAssistantConversation(c context.Context, iRequest *web_api.GetAllAssistantConversationRequest) (*web_api.GetAllAssistantConversationResponse, error) {
	assistant.logger.Debugf("GetAllAssistantConversation from grpc with requestPayload %v, %v", iRequest, c)
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(c)
	if !isAuthenticated {
		assistant.logger.Errorf("unauthenticated request for get actvities")
		return exceptions.AuthenticationError[web_api.GetAllAssistantConversationResponse]()
	}

	_page, _assistant, err := assistant.assistantClient.GetAllAssistantConversation(c, iAuth, iRequest.GetAssistantId(), iRequest.GetCriterias(), iRequest.GetPaginate(), nil)
	if err != nil {
		return exceptions.InternalServerError[web_api.GetAllAssistantConversationResponse](err, "Unable to get all the assistant sessions")
	}

	return utils.PaginatedSuccess[web_api.GetAllAssistantConversationResponse, []*web_api.AssistantConversation](
		_page.GetTotalItem(), _page.GetCurrentPage(),
		_assistant)
}

func (assistant *webAssistantGRPCApi) GetAllConversationMessage(c context.Context, iRequest *web_api.GetAllConversationMessageRequest) (*web_api.GetAllConversationMessageResponse, error) {
	assistant.logger.Debugf("GetAllConversationMessage from grpc with requestPayload %v, %v", iRequest, c)
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(c)
	if !isAuthenticated {
		assistant.logger.Errorf("unauthenticated request for get actvities")
		return exceptions.AuthenticationError[web_api.GetAllConversationMessageResponse]()
	}

	_page, _assistant, err := assistant.assistantClient.GetAllConversationMessage(c, iAuth, iRequest.GetAssistantId(), iRequest.GetAssistantConversationId(), iRequest.GetCriterias(), iRequest.GetPaginate(), nil)
	if err != nil {
		return exceptions.InternalServerError[web_api.GetAllConversationMessageResponse](err, "Unable to get all the assistant sessions")
	}

	return utils.PaginatedSuccess[web_api.GetAllConversationMessageResponse, []*web_api.AssistantConversationMessage](
		_page.GetTotalItem(), _page.GetCurrentPage(),
		_assistant)
}

// GetAllAssistantMessage implements web_api.AssistantServiceServer.
func (assistant *webAssistantGRPCApi) GetAllAssistantMessage(c context.Context, iRequest *web_api.GetAllAssistantMessageRequest) (*web_api.GetAllAssistantMessageResponse, error) {
	assistant.logger.Debugf("GetAllAssistantMessage from grpc with requestPayload %v, %v", iRequest, c)
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(c)
	if !isAuthenticated {
		assistant.logger.Errorf("unauthenticated request for get actvities")
		return exceptions.AuthenticationError[web_api.GetAllAssistantMessageResponse]()
	}

	_page, _assistant, err := assistant.assistantClient.GetAllAssistantMessage(c, iAuth, iRequest.GetAssistantId(), iRequest.GetCriterias(), iRequest.GetPaginate(), iRequest.GetOrder(), iRequest.GetSelectors())
	if err != nil {
		return exceptions.InternalServerError[web_api.GetAllAssistantMessageResponse](err, "Unable to get all the assistant messages")
	}

	return utils.PaginatedSuccess[web_api.GetAllAssistantMessageResponse, []*web_api.AssistantConversationMessage](
		_page.GetTotalItem(), _page.GetCurrentPage(),
		_assistant)
}

func (assistant *webAssistantGRPCApi) GetAllMessage(c context.Context, iRequest *web_api.GetAllMessageRequest) (*web_api.GetAllMessageResponse, error) {
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(c)
	if !isAuthenticated {
		assistant.logger.Errorf("unauthenticated request for get actvities")
		return exceptions.AuthenticationError[web_api.GetAllMessageResponse]()
	}

	_page, _assistant, err := assistant.assistantClient.GetAllMessage(c, iAuth, iRequest.GetCriterias(), iRequest.GetPaginate(), iRequest.GetOrder(), iRequest.GetSelectors())
	if err != nil {
		return exceptions.InternalServerError[web_api.GetAllMessageResponse](err, "Unable to get all the assistant messages")
	}

	return utils.PaginatedSuccess[web_api.GetAllMessageResponse, []*web_api.AssistantConversationMessage](
		_page.GetTotalItem(), _page.GetCurrentPage(),
		_assistant)
}

func (assistant *webAssistantGRPCApi) GetAssistant(c context.Context, iRequest *web_api.GetAssistantRequest) (*web_api.GetAssistantResponse, error) {
	assistant.logger.Debugf("GetAssistant from grpc with requestPayload %v, %v", iRequest, c)
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(c)
	if !isAuthenticated {
		assistant.logger.Errorf("unauthenticated request for get actvities")
		return nil, errors.New("unauthenticated request")
	}
	_assistant, err := assistant.assistantClient.GetAssistant(c, iAuth, iRequest)
	if err != nil {
		return utils.Error[web_api.GetAssistantResponse](
			err,
			"Unable to get your assistant, please try again in sometime.")
	}

	if _assistant.GetSuccess() {
		providerModel := _assistant.GetData().GetAssistantProviderModel()
		if providerModel != nil {
			user := assistant.GetUser(c, iAuth, providerModel.GetCreatedBy())
			providerModel.CreatedUser = user
			_assistant.GetData().AssistantProviderModel = providerModel
		}

		agentKit := _assistant.GetData().GetAssistantProviderAgentkit()
		if agentKit != nil {
			user := assistant.GetUser(c, iAuth, agentKit.GetCreatedBy())
			agentKit.CreatedUser = user
			_assistant.GetData().AssistantProviderAgentkit = agentKit
		}

		websocket := _assistant.GetData().GetAssistantProviderWebsocket()
		if websocket != nil {
			user := assistant.GetUser(c, iAuth, websocket.GetCreatedBy())
			websocket.CreatedUser = user
			_assistant.GetData().AssistantProviderWebsocket = websocket
		}

	}
	return _assistant, nil
}

/*
 */

/*
 */
func (assistant *webAssistantGRPCApi) GetAllAssistant(c context.Context, iRequest *web_api.GetAllAssistantRequest) (*web_api.GetAllAssistantResponse, error) {
	assistant.logger.Debugf("GetAllAssistant from grpc with requestPayload %v, %v", iRequest, c)
	iAuth, isAuthenticated := types.GetSimplePrincipleGRPC(c)
	if !isAuthenticated {
		assistant.logger.Errorf("unauthenticated request for get actvities")
		return nil, errors.New("unauthenticated request")
	}

	_page, _assistant, err := assistant.assistantClient.GetAllAssistant(c, iAuth, iRequest.GetCriterias(), iRequest.GetPaginate())
	if err != nil {
		return utils.Error[web_api.GetAllAssistantResponse](
			err,
			"Unable to get your assistant, please try again in sometime.")
	}

	for _, ast := range _assistant {
		providerModel := ast.GetAssistantProviderModel()
		if providerModel != nil {
			user := assistant.GetUser(c, iAuth, providerModel.GetCreatedBy())
			providerModel.CreatedUser = user
			ast.AssistantProviderModel = providerModel
		}

		agentKit := ast.GetAssistantProviderAgentkit()
		if agentKit != nil {
			user := assistant.GetUser(c, iAuth, agentKit.GetCreatedBy())
			agentKit.CreatedUser = user
			ast.AssistantProviderAgentkit = agentKit
		}

		websocket := ast.GetAssistantProviderWebsocket()
		if websocket != nil {
			user := assistant.GetUser(c, iAuth, websocket.GetCreatedBy())
			websocket.CreatedUser = user
			ast.AssistantProviderWebsocket = websocket
		}
	}
	return utils.PaginatedSuccess[web_api.GetAllAssistantResponse, []*web_api.Assistant](
		_page.GetTotalItem(), _page.GetCurrentPage(),
		_assistant)
}

func (assistant *webAssistantGRPCApi) CreateAssistant(c context.Context, iRequest *web_api.CreateAssistantRequest) (*web_api.GetAssistantResponse, error) {
	assistant.logger.Debugf("Create assistant from grpc with requestPayload %v, %v", iRequest, c)
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(c)
	if !isAuthenticated {
		assistant.logger.Errorf("unauthenticated request for get actvities")
		return nil, errors.New("unauthenticated request")
	}
	return assistant.assistantClient.CreateAssistant(c, iAuth, iRequest)
}

func (assistantGRPCApi *webAssistantGRPCApi) GetAllAssistantProvider(ctx context.Context, iRequest *web_api.GetAllAssistantProviderRequest) (*web_api.GetAllAssistantProviderResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request for get actvities")
		return nil, errors.New("unauthenticated request")
	}

	_page, _assistantProviders, err := assistantGRPCApi.assistantClient.GetAllAssistantProvider(ctx, iAuth, iRequest.GetAssistantId(), iRequest.GetCriterias(), iRequest.GetPaginate())
	if err != nil {
		return utils.Error[web_api.GetAllAssistantProviderResponse](
			err,
			"Unable to get your assistant provider models, please try again in sometime.")
	}

	for _, ast := range _assistantProviders {
		if ast.GetAssistantProvider() != nil {
			switch assistantProvider := ast.GetAssistantProvider().(type) {
			case *web_api.GetAllAssistantProviderResponse_AssistantProvider_AssistantProviderAgentkit:
				user := assistantGRPCApi.GetUser(ctx, iAuth, assistantProvider.AssistantProviderAgentkit.GetCreatedBy())
				assistantProvider.AssistantProviderAgentkit.CreatedUser = user
				ast.AssistantProvider = assistantProvider
			case *web_api.GetAllAssistantProviderResponse_AssistantProvider_AssistantProviderModel:
				user := assistantGRPCApi.GetUser(ctx, iAuth, assistantProvider.AssistantProviderModel.GetCreatedBy())
				assistantProvider.AssistantProviderModel.CreatedUser = user
				ast.AssistantProvider = assistantProvider
			case *web_api.GetAllAssistantProviderResponse_AssistantProvider_AssistantProviderWebsocket:
				user := assistantGRPCApi.GetUser(ctx, iAuth, assistantProvider.AssistantProviderWebsocket.GetCreatedBy())
				assistantProvider.AssistantProviderWebsocket.CreatedUser = user
				ast.AssistantProvider = assistantProvider
			}
		}
	}
	return &web_api.GetAllAssistantProviderResponse{
		Code:      200,
		Success:   true,
		Paginated: _page,
		Data:      _assistantProviders,
	}, nil

}

func (assistantGRPCApi *webAssistantGRPCApi) UpdateAssistantVersion(ctx context.Context, iRequest *web_api.UpdateAssistantVersionRequest) (*web_api.GetAssistantResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request for get actvities")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.UpdateAssistantVersion(
		ctx,
		iAuth,
		iRequest.GetAssistantId(), iRequest.GetAssistantProviderId())
}

func (assistantGRPCApi *webAssistantGRPCApi) CreateAssistantProvider(ctx context.Context, iRequest *web_api.CreateAssistantProviderRequest) (*web_api.GetAssistantProviderResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant provider model")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.CreateAssistantProvider(ctx, iAuth, iRequest)
}

// CreateAssistantTag implements web_api.AssistantServiceServer.
func (assistantGRPCApi *webAssistantGRPCApi) CreateAssistantTag(ctx context.Context, iRequest *web_api.CreateAssistantTagRequest) (*web_api.GetAssistantResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.CreateAssistantTag(ctx, iAuth, iRequest)
}

func (assistantGRPCApi *webAssistantGRPCApi) UpdateAssistantDetail(ctx context.Context, iRequest *web_api.UpdateAssistantDetailRequest) (*web_api.GetAssistantResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.UpdateAssistantDetail(ctx, iAuth, iRequest)
}

// CreateAssistantWebhook implements web_api.AssistantServiceServer.
func (assistantGRPCApi *webAssistantGRPCApi) CreateAssistantWebhook(ctx context.Context, iRequest *web_api.CreateAssistantWebhookRequest) (*web_api.GetAssistantWebhookResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.CreateAssistantWebhook(ctx, iAuth, iRequest)

}

func (assistantGRPCApi *webAssistantGRPCApi) UpdateAssistantWebhook(ctx context.Context, iRequest *web_api.UpdateAssistantWebhookRequest) (*web_api.GetAssistantWebhookResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to Update assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.UpdateAssistantWebhook(ctx, iAuth, iRequest)

}

func (assistantGRPCApi *webAssistantGRPCApi) DeleteAssistantWebhook(ctx context.Context, iRequest *web_api.DeleteAssistantWebhookRequest) (*web_api.GetAssistantWebhookResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to Delete assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.DeleteAssistantWebhook(ctx, iAuth, iRequest)

}

// GetAllAssistantWebhook implements web_api.AssistantServiceServer.
func (assistantGRPCApi *webAssistantGRPCApi) GetAllAssistantWebhook(ctx context.Context, iRequest *web_api.GetAllAssistantWebhookRequest) (*web_api.GetAllAssistantWebhookResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}

	page, tls, err := assistantGRPCApi.assistantClient.GetAllAssistantWebhook(ctx, iAuth, iRequest.GetAssistantId(), iRequest.GetCriterias(), iRequest.GetPaginate())
	if err != nil {
		return utils.Error[web_api.GetAllAssistantWebhookResponse](
			err,
			"Unable to get all the webhooks, please try again later.",
		)
	}
	return utils.PaginatedSuccess[web_api.GetAllAssistantWebhookResponse, []*web_api.AssistantWebhook](
		page.GetTotalItem(), page.GetCurrentPage(),
		tls)
}

// GetAllAssistantWebhookLog implements web_api.AssistantServiceServer.
func (assistantGRPCApi *webAssistantGRPCApi) GetAllAssistantWebhookLog(ctx context.Context, iRequest *web_api.GetAllAssistantWebhookLogRequest) (*web_api.GetAllAssistantWebhookLogResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}

	page, tls, err := assistantGRPCApi.assistantClient.GetAllAssistantWebhookLog(ctx, iAuth,
		iRequest.GetProjectId(),
		iRequest.GetCriterias(), iRequest.GetPaginate(), iRequest.GetOrder())
	if err != nil {
		return utils.Error[web_api.GetAllAssistantWebhookLogResponse](
			err,
			"Unable to get all the webhook logs, please try again later.",
		)
	}

	return utils.PaginatedSuccess[web_api.GetAllAssistantWebhookLogResponse, []*web_api.AssistantWebhookLog](
		page.GetTotalItem(), page.GetCurrentPage(),
		tls)

}

func (assistantGRPCApi *webAssistantGRPCApi) GetAssistantWebhookLog(ctx context.Context, iRequest *web_api.GetAssistantWebhookLogRequest) (*web_api.GetAssistantWebhookLogResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}

	return assistantGRPCApi.assistantClient.GetAssistantWebhookLog(ctx, iAuth, iRequest)
}

// GetAssistantWebhook implements web_api.AssistantServiceServer.
func (assistantGRPCApi *webAssistantGRPCApi) GetAssistantWebhook(ctx context.Context, iRequest *web_api.GetAssistantWebhookRequest) (*web_api.GetAssistantWebhookResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}

	return assistantGRPCApi.assistantClient.GetAssistantWebhook(ctx, iAuth, iRequest)
}

// GetAssistantWebhook implements web_api.AssistantServiceServer.
func (assistantGRPCApi *webAssistantGRPCApi) GetAssistantConversation(ctx context.Context, iRequest *web_api.GetAssistantConversationRequest) (*web_api.GetAssistantConversationResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.GetAssistantConversation(ctx, iAuth, iRequest)

}

func (assistantGRPCApi *webAssistantGRPCApi) DeleteAssistant(ctx context.Context, iRequest *web_api.DeleteAssistantRequest) (*web_api.GetAssistantResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.DeleteAssistant(ctx, iAuth, iRequest)

}

// CreateAssistantWebhook implements web_api.AssistantServiceServer.
func (assistantGRPCApi *webAssistantGRPCApi) CreateAssistantKnowledge(ctx context.Context, iRequest *web_api.CreateAssistantKnowledgeRequest) (*web_api.GetAssistantKnowledgeResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.CreateAssistantKnowledge(ctx, iAuth, iRequest)

}

func (assistantGRPCApi *webAssistantGRPCApi) UpdateAssistantKnowledge(ctx context.Context, iRequest *web_api.UpdateAssistantKnowledgeRequest) (*web_api.GetAssistantKnowledgeResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to Update assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.UpdateAssistantKnowledge(ctx, iAuth, iRequest)

}

func (assistantGRPCApi *webAssistantGRPCApi) DeleteAssistantKnowledge(ctx context.Context, iRequest *web_api.DeleteAssistantKnowledgeRequest) (*web_api.GetAssistantKnowledgeResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to Delete assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.DeleteAssistantKnowledge(ctx, iAuth, iRequest)

}

func (assistantGRPCApi *webAssistantGRPCApi) GetAllAssistantKnowledge(ctx context.Context, iRequest *web_api.GetAllAssistantKnowledgeRequest) (*web_api.GetAllAssistantKnowledgeResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}

	page, tls, err := assistantGRPCApi.assistantClient.GetAllAssistantKnowledge(ctx, iAuth, iRequest.GetAssistantId(), iRequest.GetCriterias(), iRequest.GetPaginate())
	if err != nil {
		return utils.Error[web_api.GetAllAssistantKnowledgeResponse](
			err,
			"Unable to get all the assistant knowledge, please try again later.",
		)
	}

	return utils.PaginatedSuccess[web_api.GetAllAssistantKnowledgeResponse, []*web_api.AssistantKnowledge](
		page.GetTotalItem(), page.GetCurrentPage(),
		tls)

}

func (assistantGRPCApi *webAssistantGRPCApi) GetAssistantKnowledge(ctx context.Context, iRequest *web_api.GetAssistantKnowledgeRequest) (*web_api.GetAssistantKnowledgeResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create knowledge tag")
		return nil, errors.New("unauthenticated request")
	}

	return assistantGRPCApi.assistantClient.GetAssistantKnowledge(ctx, iAuth, iRequest)
}

func (assistantGRPCApi *webAssistantGRPCApi) CreateAssistantTool(ctx context.Context, iRequest *web_api.CreateAssistantToolRequest) (*web_api.GetAssistantToolResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.CreateAssistantTool(ctx, iAuth, iRequest)

}

func (assistantGRPCApi *webAssistantGRPCApi) UpdateAssistantTool(ctx context.Context, iRequest *web_api.UpdateAssistantToolRequest) (*web_api.GetAssistantToolResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to Update assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.UpdateAssistantTool(ctx, iAuth, iRequest)

}

func (assistantGRPCApi *webAssistantGRPCApi) DeleteAssistantTool(ctx context.Context, iRequest *web_api.DeleteAssistantToolRequest) (*web_api.GetAssistantToolResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to Delete assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.DeleteAssistantTool(ctx, iAuth, iRequest)

}

func (assistantGRPCApi *webAssistantGRPCApi) GetAllAssistantTool(ctx context.Context, iRequest *web_api.GetAllAssistantToolRequest) (*web_api.GetAllAssistantToolResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}

	page, tls, err := assistantGRPCApi.assistantClient.GetAllAssistantTool(ctx, iAuth, iRequest.GetAssistantId(), iRequest.GetCriterias(), iRequest.GetPaginate())
	if err != nil {
		return utils.Error[web_api.GetAllAssistantToolResponse](
			err,
			"Unable to get all the webhook analysis, please try again later.",
		)
	}

	return utils.PaginatedSuccess[web_api.GetAllAssistantToolResponse, []*web_api.AssistantTool](
		page.GetTotalItem(), page.GetCurrentPage(),
		tls)

}

func (assistantGRPCApi *webAssistantGRPCApi) GetAssistantTool(ctx context.Context, iRequest *web_api.GetAssistantToolRequest) (*web_api.GetAssistantToolResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}

	return assistantGRPCApi.assistantClient.GetAssistantTool(ctx, iAuth, iRequest)
}

func (assistantGRPCApi *webAssistantGRPCApi) CreateAssistantAnalysis(ctx context.Context, iRequest *web_api.CreateAssistantAnalysisRequest) (*web_api.GetAssistantAnalysisResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to create assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.CreateAssistantAnalysis(ctx, iAuth, iRequest)

}

func (assistantGRPCApi *webAssistantGRPCApi) UpdateAssistantAnalysis(ctx context.Context, iRequest *web_api.UpdateAssistantAnalysisRequest) (*web_api.GetAssistantAnalysisResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to Update assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.UpdateAssistantAnalysis(ctx, iAuth, iRequest)

}

func (assistantGRPCApi *webAssistantGRPCApi) DeleteAssistantAnalysis(ctx context.Context, iRequest *web_api.DeleteAssistantAnalysisRequest) (*web_api.GetAssistantAnalysisResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to Delete assistant tag")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.DeleteAssistantAnalysis(ctx, iAuth, iRequest)

}

func (assistantGRPCApi *webAssistantGRPCApi) GetAllAssistantAnalysis(ctx context.Context, iRequest *web_api.GetAllAssistantAnalysisRequest) (*web_api.GetAllAssistantAnalysisResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to GetAllAssistantAnalysis")
		return nil, errors.New("unauthenticated request")
	}

	page, tls, err := assistantGRPCApi.assistantClient.GetAllAssistantAnalysis(ctx, iAuth, iRequest.GetAssistantId(), iRequest.GetCriterias(), iRequest.GetPaginate())
	if err != nil {
		return utils.Error[web_api.GetAllAssistantAnalysisResponse](
			err,
			"Unable to get all the webhook analysis, please try again later.",
		)
	}

	return utils.PaginatedSuccess[web_api.GetAllAssistantAnalysisResponse, []*web_api.AssistantAnalysis](
		page.GetTotalItem(), page.GetCurrentPage(),
		tls)

}

func (assistantGRPCApi *webAssistantGRPCApi) GetAssistantAnalysis(ctx context.Context, iRequest *web_api.GetAssistantAnalysisRequest) (*web_api.GetAssistantAnalysisResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to GetAssistantAnalysis")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.GetAssistantAnalysis(ctx, iAuth, iRequest)
}

func (assistantGRPCApi *webAssistantGRPCApi) GetAssistantToolLog(ctx context.Context, iRequest *web_api.GetAssistantToolLogRequest) (*web_api.GetAssistantToolLogResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to GetAssistantToolLog")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.GetAssistantToolLog(ctx, iAuth, iRequest)
}

func (assistantGRPCApi *webAssistantGRPCApi) GetAllAssistantToolLog(ctx context.Context, iRequest *web_api.GetAllAssistantToolLogRequest) (*web_api.GetAllAssistantToolLogResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to GetAllAssistantToolLog")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.GetAllAssistantToolLog(ctx, iAuth, iRequest)
}

func (assistantGRPCApi *webAssistantGRPCApi) GetAllAssistantTelemetry(ctx context.Context, iRequest *web_api.GetAllAssistantTelemetryRequest) (*web_api.GetAllAssistantTelemetryResponse, error) {
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistantGRPCApi.logger.Errorf("unauthenticated request to GetAllAssistantTelemetry")
		return nil, errors.New("unauthenticated request")
	}
	return assistantGRPCApi.assistantClient.GetAllAssistantTelemetry(ctx, iAuth, iRequest)
}
