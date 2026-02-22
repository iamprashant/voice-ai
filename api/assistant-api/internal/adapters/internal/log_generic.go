// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package adapter_internal

import (
	"context"

	type_enums "github.com/rapidaai/pkg/types/enums"
	protos "github.com/rapidaai/protos"
)

const ConversationPageHistory uint32 = 50

func (kr *genericRequestor) CreateKnowledgeLog(ctx context.Context, knowledgeId uint64, retrievalMethod string,
	topK uint32,
	scoreThreshold float32,
	documentCount int,
	timeTaken int64,
	additionalData map[string]string,
	status type_enums.RecordState,
	request, response []byte) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), dbWriteTimeout)
	defer cancel()
	_, err := kr.knowledgeService.CreateLog(dbCtx, kr.Auth(), knowledgeId, retrievalMethod, topK, scoreThreshold, documentCount, timeTaken, additionalData, status, request, response)
	return err
}

func (cr *genericRequestor) CreateWebhookLog(
	ctx context.Context,
	webhookID uint64, httpUrl, httpMethod, event string,
	responseStatus int64,
	timeTaken int64,
	retryCount uint32,
	status type_enums.RecordState,
	request, response []byte) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), dbWriteTimeout)
	defer cancel()
	_, err := cr.webhookService.CreateLog(dbCtx, cr.auth, webhookID, cr.assistant.Id, cr.assistantConversation.Id, httpUrl, httpMethod, event, responseStatus, timeTaken, retryCount, status, request, response)
	return err
}

func (cr *genericRequestor) GetConversationLogs(ctx context.Context) []*protos.Message {
	messages := make([]*protos.Message, 0)
	cnt, conversations, err := cr.
		conversationService.
		GetAllMessageActions(
			ctx,
			cr.auth,
			cr.assistantConversation.Id,
			[]*protos.Criteria{
				{
					Key:   "action_type",
					Value: "llm-call",
					Logic: "=",
				},
			},
			&protos.Paginate{
				Page:     1,
				PageSize: ConversationPageHistory,
			},
			&protos.Ordering{
				Column: "created_date",
				Order:  "asc",
			},
		)

	if cnt == 0 || err != nil {
		return messages
	}

	for _, x := range conversations {
		if x.Status == type_enums.RECORD_SUCCESS || x.Status == type_enums.RECORD_ACTIVE {
			messages = append(messages, x.RequestMessage())
			messages = append(messages, x.ResponseMessage())
		}
	}
	return messages
}

func (cr *genericRequestor) CreateConversationMessageLog(ctx context.Context, messageid string, in, out *protos.Message, metrics []*protos.Metric) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), dbWriteTimeout)
	defer cancel()
	cr.conversationService.CreateLLMAction(
		dbCtx,
		cr.Auth(),
		cr.assistant.Id,
		cr.assistantConversation.Id,
		messageid,
		in, out, metrics)
	return nil
}

func (cr *genericRequestor) CreateConversationToolLog(
	ctx context.Context,
	messageid string, in, out map[string]interface{}, metrics []*protos.Metric) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), dbWriteTimeout)
	defer cancel()
	cr.conversationService.CreateToolAction(
		dbCtx,
		cr.Auth(),
		cr.assistant.Id,
		cr.assistantConversation.Id,
		messageid,
		in, out, metrics)
	return nil
}

func (cr *genericRequestor) CreateToolLog(
	ctx context.Context,
	toolId uint64,
	messageId string,
	toolName string,
	executionMethod string,
	status type_enums.RecordState,
	timeTaken int64,
	request, response []byte) error {
	dbCtx, cancel := context.WithTimeout(context.Background(), dbWriteTimeout)
	defer cancel()
	_, err := cr.assistantToolService.CreateLog(
		dbCtx, cr.Auth(), cr.assistant.Id,
		cr.assistantConversation.Id, toolId, messageId, toolName, timeTaken, executionMethod,
		status, request, response,
	)
	return err
}
