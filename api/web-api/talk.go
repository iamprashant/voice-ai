package web_api

import (
	"context"
	"errors"
	"io"

	assistant_client "github.com/lexatic/web-backend/pkg/clients/workflow"
	"github.com/lexatic/web-backend/pkg/utils"
	web_api "github.com/lexatic/web-backend/protos/lexatic-backend"

	config "github.com/lexatic/web-backend/config"
	commons "github.com/lexatic/web-backend/pkg/commons"
	"github.com/lexatic/web-backend/pkg/connectors"
	"github.com/lexatic/web-backend/pkg/types"
)

type webTalkApi struct {
	WebApi
	cfg                         *config.AppConfig
	logger                      commons.Logger
	postgres                    connectors.PostgresConnector
	redis                       connectors.RedisConnector
	assistantConversationClient assistant_client.AssistantConversationServiceClient
}

type webTalkGRPCApi struct {
	webTalkApi
}

// MessageFeedback implements lexatic_backend.TalkServiceServer.
func (*webTalkGRPCApi) MessageFeedback(context.Context, *web_api.MessageFeedbackRequest) (*web_api.MessageFeedbackResponse, error) {
	panic("unimplemented")
}

// AssistantTalk implements lexatic_backend.TalkServiceServer.
func (*webTalkGRPCApi) AssistantTalk(web_api.TalkService_AssistantTalkServer) error {
	panic("unimplemented")
}

func NewTalkGRPC(config *config.AppConfig, logger commons.Logger, postgres connectors.PostgresConnector, redis connectors.RedisConnector) web_api.TalkServiceServer {
	return &webTalkGRPCApi{
		webTalkApi{
			WebApi:                      NewWebApi(config, logger, postgres, redis),
			cfg:                         config,
			logger:                      logger,
			postgres:                    postgres,
			redis:                       redis,
			assistantConversationClient: assistant_client.NewAssistantConversationServiceClientGRPC(config, logger, redis),
		},
	}
}

//
//

// GetAllConversationMessage implements lexatic_backend.AssistantConversationServiceServer.
func (assistant *webTalkGRPCApi) GetAllConversationMessage(ctx context.Context, iRequest *web_api.GetAllConversationMessageRequest) (*web_api.GetAllConversationMessageResponse, error) {
	assistant.logger.Debugf("GetAllConversationMessage started")
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(ctx)
	if !isAuthenticated {
		assistant.logger.Errorf("unauthenticated request for get actvities")
		return nil, errors.New("unauthenticated request")
	}

	_page, _assistant, err := assistant.assistantConversationClient.GetAllConversationMessage(ctx, iAuth, iRequest.GetAssistantId(), iRequest.GetAssistantConversationId(), iRequest.GetCriterias(), iRequest.GetPaginate())
	if err != nil {
		return utils.Error[web_api.GetAllConversationMessageResponse](
			err,
			"Unable to get your assistant, please try again in sometime.")
	}

	return utils.PaginatedSuccess[web_api.GetAllConversationMessageResponse, []*web_api.AssistantConversationMessage](
		_page.GetTotalItem(), _page.GetCurrentPage(),
		_assistant)

}

func (assistant *webTalkGRPCApi) AssistantMessaging(cer *web_api.AssistantMessagingRequest, stream web_api.TalkService_AssistantMessagingServer) error {
	assistant.logger.Debugf("AssistantMessaging started")
	c := stream.Context()
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(c)
	if !isAuthenticated {
		assistant.logger.Errorf("unauthenticated request for get actvities")
		return errors.New("unauthenticated request")
	}
	out, err := assistant.assistantConversationClient.AssistantMessaging(c, iAuth, cer)
	if err != nil {
		return err
	}

	// Channel to handle errors from the upstream stream
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		for {
			st, recvErr := out.Recv()
			if recvErr == io.EOF {
				return // End of upstream stream
			}
			if recvErr != nil {
				errCh <- recvErr
				return
			}
			// Forward message to downstream stream
			if err := stream.Send(st); err != nil {
				errCh <- err
				return
			}
		}
	}()

	// Wait for any errors from the upstream stream or the context cancellation
	select {
	case err := <-errCh:
		return err
	case <-c.Done():
		return c.Err()
	}

}

// GetAllAssistantConversation implements lexatic_backend.AssistantConversationServiceServer.
func (assistant *webTalkGRPCApi) GetAllAssistantConversation(c context.Context, iRequest *web_api.GetAllAssistantConversationRequest) (*web_api.GetAllAssistantConversationResponse, error) {
	assistant.logger.Debugf("AssistantMessaging started")
	iAuth, isAuthenticated := types.GetAuthPrincipleGPRC(c)
	if !isAuthenticated {
		assistant.logger.Errorf("unauthenticated request for get actvities")
		return nil, errors.New("unauthenticated request")
	}

	_page, _assistantConvo, err := assistant.assistantConversationClient.GetAllAssistantConversation(c, iAuth,
		iRequest.GetAssistantId(),
		iRequest.GetCriterias(), iRequest.GetPaginate())
	if err != nil {
		return utils.Error[web_api.GetAllAssistantConversationResponse](
			err,
			"Unable to get your assistant, please try again in sometime.")
	}

	return utils.PaginatedSuccess[web_api.GetAllAssistantConversationResponse, []*web_api.AssistantConversation](
		_page.GetTotalItem(), _page.GetCurrentPage(),
		_assistantConvo)
}
