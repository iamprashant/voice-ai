// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package assistant_talk_api

import (
	"context"
	"fmt"

	callcontext "github.com/rapidaai/api/assistant-api/internal/callcontext"
	internal_services "github.com/rapidaai/api/assistant-api/internal/services"
	"github.com/rapidaai/pkg/types"
	type_enums "github.com/rapidaai/pkg/types/enums"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

// InitiateAssistantTalk implements protos.TalkServiceServer.
func (cApi *ConversationGrpcApi) CreatePhoneCall(ctx context.Context, ir *protos.CreatePhoneCallRequest) (*protos.CreatePhoneCallResponse, error) {
	auth, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated {
		cApi.logger.Errorf("unable to resolve the authentication object, please check the parameter for authentication")
		return utils.AuthenticateError[protos.CreatePhoneCallResponse]()
	}

	toNumber := ir.GetToNumber()
	if utils.IsEmpty(toNumber) {
		return utils.ErrorWithCode[protos.CreatePhoneCallResponse](200, fmt.Errorf("missing to_phone parameter"), "Please provide the required to_phone parameter.")
	}

	mtd, err := utils.AnyMapToInterfaceMap(ir.GetMetadata())
	if err != nil {
		return utils.ErrorWithCode[protos.CreatePhoneCallResponse](200, err, "Illegal metadata for initialize request, please check and try again.")
	}

	args, err := utils.AnyMapToInterfaceMap(ir.GetArgs())
	if err != nil {
		return utils.ErrorWithCode[protos.CreatePhoneCallResponse](200, err, "Illegal options for initialize request, please check and try again.")
	}

	opts, err := utils.AnyMapToInterfaceMap(ir.GetOptions())
	if err != nil {
		return utils.ErrorWithCode[protos.CreatePhoneCallResponse](200, err, "Illegal arguments for initialize request, please check and try again.")
	}

	assistant, err := cApi.assistantService.Get(ctx, auth, ir.GetAssistant().GetAssistantId(), utils.GetVersionDefinition(ir.GetAssistant().GetVersion()), &internal_services.GetAssistantOption{InjectPhoneDeployment: true})
	if err != nil {
		cApi.logger.Debugf("illegal unable to find assistant %v", err)
		return utils.ErrorWithCode[protos.CreatePhoneCallResponse](200, err, "Invalid assistant id, please check and try again.")
	}

	if !assistant.IsPhoneDeploymentEnable() {
		cApi.logger.Debugf("illegal deployment for phone %v", err)
		return utils.ErrorWithCode[protos.CreatePhoneCallResponse](200, err, "Phone deployment not enabled or incomplete, please check rapida console and update the deployment")
	}

	// creating conversation
	conversation, err := cApi.assistantConversationService.CreateConversation(ctx, auth, toNumber, assistant.Id, assistant.AssistantProviderId, type_enums.DIRECTION_OUTBOUND, utils.PhoneCall)
	if err != nil {
		cApi.logger.Errorf("unable to create conversation %+v", err)
		return utils.ErrorWithCode[protos.CreatePhoneCallResponse](200, err, "Unable to create conversation session, please check and try again.")
	}
	o, err := cApi.assistantConversationService.ApplyConversationOption(ctx, auth, assistant.Id, conversation.Id, opts)
	if err != nil {
		cApi.logger.Debugf("unable to create options %v", err)
		return utils.ErrorWithCode[protos.CreatePhoneCallResponse](200, err, "Unable to create conversation options, please check and try again.")
	}
	conversation.Options = o
	// updating arguments
	arguments, err := cApi.assistantConversationService.ApplyConversationArgument(ctx, auth, assistant.Id, conversation.Id, args)
	if err != nil {
		cApi.logger.Debugf("unable to create argument %v", err)
		return utils.ErrorWithCode[protos.CreatePhoneCallResponse](200, err, "Unable to create conversation arguments, please check and try again.")
	}
	conversation.Arguments = arguments

	// Resolve from phone number
	fromPhone := ir.GetFromNumber()
	if utils.IsEmpty(fromPhone) {
		fromNumber, err := assistant.AssistantPhoneDeployment.GetOptions().GetString("phone")
		if err != nil {
			cApi.assistantConversationService.ApplyConversationMetrics(ctx, auth, assistant.Id, conversation.Id, []*types.Metric{types.NewStatusMetric(type_enums.RECORD_FAILED)})
			return utils.ErrorWithCode[protos.CreatePhoneCallResponse](200, fmt.Errorf("failed to get phone number"), "Unable to retrieve the default phone number.")
		}
		fromPhone = fromNumber
	}

	// Apply metadata early (before async dispatch)
	if len(mtd) > 0 {
		mtdas, err := cApi.assistantConversationService.ApplyConversationMetadata(ctx, auth, assistant.Id, conversation.Id, types.NewMetadataList(mtd))
		if err != nil {
			cApi.logger.Errorf("failed to apply conversation metadata: %v", err)
		} else {
			conversation.Metadatas = mtdas
		}
	}

	// Store call context in Postgres for async worker resolution
	cc := &callcontext.CallContext{
		AssistantID:         assistant.Id,
		ConversationID:      conversation.Id,
		AssistantProviderId: assistant.AssistantProviderId,
		AuthToken:           auth.GetCurrentToken(),
		AuthType:            auth.Type(),
		Direction:           "outbound",
		CallerNumber:        toNumber,
		CalleeNumber:        toNumber,
		FromNumber:          fromPhone,
		Provider:            assistant.AssistantPhoneDeployment.TelephonyProvider,
		Status:              "queued",
	}
	if auth.GetCurrentProjectId() != nil {
		cc.ProjectID = *auth.GetCurrentProjectId()
	}
	if auth.GetCurrentOrganizationId() != nil {
		cc.OrganizationID = *auth.GetCurrentOrganizationId()
	}
	contextID, err := cApi.callContextStore.Save(ctx, cc)
	if err != nil {
		cApi.logger.Errorf("failed to save call context for outbound call: %v", err)
		cApi.assistantConversationService.ApplyConversationMetrics(ctx, auth, assistant.Id, conversation.Id, []*types.Metric{types.NewStatusMetric(type_enums.RECORD_FAILED)})
		return utils.ErrorWithCode[protos.CreatePhoneCallResponse](500, err, "Failed to create call context")
	}

	// Dispatch outbound call asynchronously — the goroutine resolves the call context,
	// fetches vault credentials, and initiates the telephony call without blocking the gRPC response.
	err = cApi.outboundDispatcher.Dispatch(context.Background(), contextID)
	if err != nil {
		cApi.logger.Errorf("failed to dispatch outbound call: %v", err)
		cApi.assistantConversationService.ApplyConversationMetrics(ctx, auth, assistant.Id, conversation.Id, []*types.Metric{types.NewStatusMetric(type_enums.RECORD_FAILED)})
		return utils.ErrorWithCode[protos.CreatePhoneCallResponse](500, err, "Failed to dispatch outbound call")
	}

	cApi.logger.Infof("outbound call dispatched: contextId=%s, provider=%s, assistant=%d, conversation=%d",
		contextID, cc.Provider, assistant.Id, conversation.Id)

	// Apply initial metadata for the queued call
	cApi.assistantConversationService.ApplyConversationMetadata(ctx, auth, assistant.Id, conversation.Id, []*types.Metadata{
		types.NewMetadata("telephony.contextId", contextID),
		types.NewMetadata("telephony.toPhone", toNumber),
		types.NewMetadata("telephony.fromPhone", fromPhone),
		types.NewMetadata("telephony.provider", cc.Provider),
	})

	// Return immediately — the worker will handle the actual telephony call
	out := &protos.AssistantConversation{}
	err = utils.Cast(conversation, out)
	if err != nil {
		cApi.logger.Errorf("unable to cast assistant conversation %v", err)
	}
	return utils.Success[protos.CreatePhoneCallResponse, *protos.AssistantConversation](out)
}

// InitiateBulkAssistantTalk implements protos.TalkServiceServer.
func (cApi *ConversationGrpcApi) CreateBulkPhoneCall(ctx context.Context, ir *protos.CreateBulkPhoneCallRequest) (*protos.CreateBulkPhoneCallResponse, error) {
	_, isAuthenticated := types.GetSimplePrincipleGRPC(ctx)
	if !isAuthenticated {
		cApi.logger.Errorf("unable to resolve the authentication object, please check the parameter for authentication")
		return utils.AuthenticateError[protos.CreateBulkPhoneCallResponse]()
	}

	out := make([]*protos.AssistantConversation, 0)
	for _, v := range ir.GetPhoneCalls() {
		resp, err := cApi.CreatePhoneCall(ctx, v)
		if err != nil {
			cApi.logger.Errorf("error while making call %+v", err)
		}
		if resp.GetData() != nil {
			out = append(out, resp.GetData())
		}
	}
	return utils.Success[protos.CreateBulkPhoneCallResponse, []*protos.AssistantConversation](out)
}
