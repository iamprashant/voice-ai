// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package assistant_talk_api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"

	internal_adapter "github.com/rapidaai/api/assistant-api/internal/adapters"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
	internal_services "github.com/rapidaai/api/assistant-api/internal/services"
	telephony "github.com/rapidaai/api/assistant-api/internal/telephony"
	"github.com/rapidaai/pkg/types"
	type_enums "github.com/rapidaai/pkg/types/enums"
	"github.com/rapidaai/pkg/utils"
	protos "github.com/rapidaai/protos"
)

func (cApi *ConversationApi) UnviersalCallback(c *gin.Context) {
	body, err := c.GetRawData() // Extract raw request body
	if err != nil {
		cApi.logger.Errorf("failed to read event body with error %+v", err)
	}
	cApi.logger.Debugf("event from exotel | body: %s", string(body))
}

func (cApi *ConversationApi) Callback(c *gin.Context) {
	iAuth, isAuthenticated := types.GetAuthPrinciple(c)
	if !isAuthenticated {
		cApi.logger.Debugf("illegal unable to authenticate")
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthenticated request"})
		return
	}
	assistantIdStr := c.Param("assistantId")
	assistantId, err := strconv.ParseUint(assistantIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assistant ID"})
		return
	}

	conversationIdStr := c.Param("conversationId")
	conversationId, err := strconv.ParseUint(conversationIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
		return
	}

	tlp := c.Param("telephony")
	_telephony, err := telephony.GetTelephony(
		telephony.Telephony(tlp),
		cApi.cfg,
		cApi.logger,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid telephony"})
		return
	}

	mtr, evts, err := _telephony.StatusCallback(c, iAuth, assistantId, conversationId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event to process"})
		return
	}

	if _, err := cApi.assistantConversationService.ApplyConversationTelephonyEvent(c, iAuth, tlp, assistantId, conversationId, evts); err != nil {
		c.Status(http.StatusOK)
		return
	}

	if _, err := cApi.assistantConversationService.ApplyConversationMetrics(c, iAuth, assistantId, conversationId, mtr); err != nil {
		c.Status(http.StatusOK)
		return
	}
	c.Status(http.StatusCreated)
	return
}

// CallReciever handles incoming calls for the given assistant.
// @Router /v1/call/:assistantId [post]
// @Summary Recieve call for given assistant
// @Produce json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
func (cApi *ConversationApi) CallReciever(c *gin.Context) {
	iAuth, isAuthenticated := types.GetAuthPrinciple(c)
	if !isAuthenticated {
		cApi.logger.Debugf("illegal unable to authenticate")
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthenticated request"})
		return
	}

	tlp := c.Param("telephony")
	_telephony, err := telephony.GetTelephony(telephony.Telephony(tlp), cApi.cfg, cApi.logger)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Telephony is not connected"})
		return
	}

	clientNumber, assistantIdStr, err := _telephony.AcceptCall(c)
	if err != nil {
		cApi.logger.Errorf(err.Error())
		return
	}

	assistantId, err := strconv.ParseUint(*assistantIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assistant ID"})
		return
	}

	assistant, err := cApi.assistantService.Get(c, iAuth, assistantId, utils.GetVersionDefinition("latest"), &internal_services.GetAssistantOption{InjectPhoneDeployment: true})
	if err != nil {
		cApi.logger.Debugf("illegal unable to find assistant %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to initiate talker"})
		return
	}

	conversation, err := cApi.assistantConversationService.CreateConversation(c, iAuth, internal_adapter.Identifier(utils.PhoneCall, c, iAuth, *clientNumber), assistant.Id, assistant.AssistantProviderId, type_enums.DIRECTION_INBOUND, utils.PhoneCall)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to initiate talker"})
		return
	}

	//
	if _, err := cApi.assistantConversationService.ApplyConversationMetrics(c, iAuth, assistantId, conversation.Id, []*types.Metric{types.NewStatusMetric(type_enums.RECORD_CONNECTED)}); err != nil {
		cApi.logger.Errorf("error while applying metrics %v", err)
	}

	if err := _telephony.IncomingCall(c, iAuth, assistant.Id, *clientNumber, conversation.Id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to initiate talker"})
		return
	}
	return

}

func (cApi *ConversationApi) CallTalker(c *gin.Context) {
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool { return true }}
	websocketConnection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to upgrade connection"})
		return
	}
	auth, isAuthenticated := types.GetAuthPrinciple(c)
	if !isAuthenticated {
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthenticated request"})
		return
	}

	// Extract the client source from the stream context
	assistantId, err := strconv.ParseUint(c.Param("assistantId"), 10, 64)
	if err != nil {
		cApi.logger.Errorf("Invalid assistantId: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assistantId"})
		return
	}

	conversationId, err := strconv.ParseUint(c.Param("conversationId"), 10, 64)
	if err != nil {
		cApi.logger.Errorf("Invalid conversationId: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversationId"})
		return
	}

	identifier := c.Param("identifier")
	tlp := c.Param("telephony")
	_telephony, err := telephony.GetTelephony(telephony.Telephony(tlp), cApi.cfg, cApi.logger)
	if err != nil {
		cApi.assistantConversationService.ApplyConversationMetrics(c, auth, assistantId, conversationId, []*types.Metric{{Name: type_enums.STATUS.String(), Value: type_enums.RECORD_FAILED.String(), Description: "Invalid telephony"}})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid telephony"})
		return
	}

	var (
		wg                    errgroup.Group
		vltC                  *protos.VaultCredential
		assistantConversation *internal_conversation_entity.AssistantConversation
		assistant             *internal_assistant_entity.Assistant
	)
	wg.Go(func() error {
		assistant, err = cApi.assistantService.Get(c, auth, assistantId, nil, &internal_services.GetAssistantOption{InjectPhoneDeployment: true})
		if err != nil {
			return err
		}

		if !assistant.IsPhoneDeploymentEnable() {
			return err
		}
		credentialID, err := assistant.AssistantPhoneDeployment.GetOptions().GetUint64("rapida.credential_id")
		if err != nil {
			return err
		}

		cApi.logger.Debugf("got aith %+v", auth)
		vltC, err = cApi.vaultClient.GetCredential(c, auth, credentialID)
		if err != nil {
			return err
		}
		return nil
	})

	// get converstation to make sure if it exists and valid
	wg.Go(func() error {
		assistantConversation, err = cApi.assistantConversationService.Get(c, auth, assistantId, conversationId, internal_services.NewDefaultGetConversationOption())
		if err != nil {
			return err
		}
		return nil
	})

	//
	if err := wg.Wait(); err != nil {
		cApi.assistantConversationService.ApplyConversationMetrics(c, auth, assistantId, conversationId, []*types.Metric{{Name: type_enums.STATUS.String(), Value: type_enums.RECORD_FAILED.String(), Description: "Invalid phone deployment credential"}})
		cApi.logger.Errorf("error while recieving call %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid phone deployment"})
		return
	}
	talker, err := internal_adapter.GetTalker(
		utils.PhoneCall, c, cApi.cfg, cApi.logger, cApi.postgres, cApi.opensearch, cApi.redis, cApi.storage,
		_telephony.Streamer(c, websocketConnection, assistant, assistantConversation, vltC))
	if err != nil {
		cApi.logger.Errorf("error while recieving call %v", err)
		cApi.assistantConversationService.ApplyConversationMetrics(c, auth, assistantId, conversationId, []*types.Metric{{Name: type_enums.STATUS.String(), Value: type_enums.RECORD_FAILED.String(), Description: "Internal server error"}})
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid talker"})
		return
	}

	if err := talker.Talk(c, auth, internal_adapter.Identifier(utils.PhoneCall, c, auth, identifier)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid talk"})
	}
}
