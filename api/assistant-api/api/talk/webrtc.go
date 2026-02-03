// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package assistant_talk_api

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var webrtcUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// WebRTCConnect handles WebRTC connections for voice conversations
// The WebSocket is used ONLY for signaling (SDP/ICE exchange)
// Audio flows through native WebRTC media tracks (SRTP), NOT WebSocket
//
// @Router /v1/talk/webrtc/:assistantId [get]
// @Summary Connect to assistant via WebRTC
// @Description Establishes a WebRTC connection for real-time voice conversation
// @Param assistantId path uint64 true "Assistant ID"
// @Produce json
// @Success 101 "Switching Protocols"
// @Failure 400 {object} gin.H
// @Failure 401 {object} gin.H
// func (cApi *ConversationApi) WebRTCConnect(c *gin.Context) {
// 	// Upgrade to WebSocket for signaling only
// 	signalingConn, err := webrtcUpgrader.Upgrade(c.Writer, c.Request, nil)
// 	if err != nil {
// 		cApi.logger.Errorf("WebSocket upgrade failed: %v", err)
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to upgrade to WebSocket"})
// 		return
// 	}

// 	// Helper to send error over WebSocket and close
// 	sendErrorAndClose := func(errMsg string) {
// 		errData, _ := json.Marshal(map[string]interface{}{
// 			"type":  "error",
// 			"error": errMsg,
// 		})
// 		signalingConn.WriteMessage(websocket.TextMessage, errData)
// 		signalingConn.Close()
// 	}

// 	auth, isAuthenticated := types.GetAuthPrinciple(c)
// 	if !isAuthenticated {
// 		cApi.logger.Error("WebRTC: Unauthenticated request")
// 		sendErrorAndClose("Unauthenticated request - missing or invalid authorization")
// 		return
// 	}

// 	// Create identifier for the conversation
// 	identifier := internal_adapter.Identifier(utils.Debugger, c, auth, "")
// 	// Create WebRTC streamer
// 	// Audio flows through WebRTC peer connection, WebSocket only for signaling
// 	streamer, err := internal_webrtc.NewStreamer(c.Request.Context(),
// 		cApi.logger,
// 		signalingConn,
// 	)
// 	if err != nil {
// 		sendErrorAndClose(fmt.Sprintf("Failed to create WebRTC streamer: %v", err))
// 		return
// 	}

// 	// Create talker with WebRTC source
// 	talker, err := internal_adapter.GetTalker(
// 		utils.Debugger,
// 		c,
// 		cApi.cfg,
// 		cApi.logger,
// 		cApi.postgres,
// 		cApi.opensearch,
// 		cApi.redis,
// 		cApi.storage,
// 		streamer,
// 	)
// 	if err != nil {
// 		if closeable, ok := streamer.(io.Closer); ok {
// 			closeable.Close()
// 		}
// 		sendErrorAndClose(fmt.Sprintf("Failed to setup conversation: %v", err))
// 		return
// 	}
// 	if err := talker.Talk(c, auth, identifier); err != nil {
// 		cApi.logger.Errorf("WebRTC conversation error: %v", err)
// 	}

// }
