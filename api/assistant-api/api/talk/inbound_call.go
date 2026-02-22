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

	internal_adapter "github.com/rapidaai/api/assistant-api/internal/adapters"
	telephony "github.com/rapidaai/api/assistant-api/internal/channel/telephony"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
)

func (cApi *ConversationApi) UnviersalCallback(c *gin.Context) {
	body, err := c.GetRawData() // Extract raw request body
	if err != nil {
		cApi.logger.Errorf("failed to read event body with error %+v", err)
	}
	cApi.logger.Debugf("event body: %s", string(body))
}

// CallbackByContext handles status callback webhooks using a contextId stored in Postgres.
// The contextId resolves to the full call context (auth, assistant, conversation, provider).
// The context is NOT deleted — callbacks fire multiple times during a call.
// Route: GET/POST /:telephony/ctx/:contextId/event
func (cApi *ConversationApi) CallbackByContext(c *gin.Context) {
	contextID := c.Param("contextId")
	if contextID == "" {
		cApi.logger.Errorf("missing contextId in CallbackByContext")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing contextId"})
		return
	}

	if err := cApi.inboundDispatcher.HandleStatusCallbackByContext(c, contextID); err != nil {
		cApi.logger.Errorf("status callback failed for context %s: %v", contextID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event to process"})
		return
	}
	c.Status(http.StatusCreated)
}

// CallReciever handles incoming calls for the given assistant.
// The telephony provider sends a webhook when an inbound call arrives.
// This handler creates a conversation, saves a CallContext to Postgres, and returns
// provider-specific instructions (TwiML, NCCO, contextId) to answer the call.
// @Router /v1/talk/:telephony/call/:assistantId [get]
// @Summary Receive call for given assistant
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

	assistantID := c.Param("assistantId")
	if assistantID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assistant ID"})
		return
	}
	assistantId, err := strconv.ParseUint(assistantID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid assistant ID"})
		return
	}

	if _, err := cApi.inboundDispatcher.HandleReceiveCall(c, c.Param("telephony"), iAuth, assistantId); err != nil {
		cApi.logger.Errorf("failed to handle inbound call: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to initiate talker"})
		return
	}
}

// CallTalkerByContext handles WebSocket connections using a contextId stored in Postgres.
// The contextId was returned by CallReciever (inbound) or CreatePhoneCall (outbound).
// All auth, assistant, conversation, and provider info is resolved from the call context.
// Route: GET /:telephony/ctx/:contextId
func (cApi *ConversationApi) CallTalkerByContext(c *gin.Context) {
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool { return true }}
	websocketConnection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to upgrade connection"})
		return
	}

	contextID := c.Param("contextId")
	if contextID == "" {
		cApi.logger.Errorf("missing contextId in CallTalkerByContext")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing contextId"})
		return
	}

	cc, vaultCred, err := cApi.inboundDispatcher.ResolveCallSessionByContext(c, contextID)
	if err != nil {
		cApi.logger.Errorf("error resolving session for context %s: %v", contextID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid call context"})
		return
	}

	streamer, err := telephony.Telephony(cc.Provider).NewStreamer(cApi.logger, cc, vaultCred, telephony.StreamerOption{
		WebSocketConn: websocketConnection,
	})
	if err != nil {
		cApi.logger.Errorf("error creating streamer for context %s: %v", contextID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid telephony streamer"})
		return
	}

	talker, err := internal_adapter.GetTalker(utils.PhoneCall, c, cApi.cfg, cApi.logger, cApi.postgres, cApi.opensearch, cApi.redis, cApi.storage, streamer)
	if err != nil {
		cApi.logger.Errorf("error creating talker for context %s: %v", contextID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid talker"})
		return
	}

	if err := talker.Talk(c, cc.ToAuth()); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid talk"})
	}

	// Mark call context as completed now that the call has ended
	cApi.inboundDispatcher.CompleteCallSession(c, contextID)
}
