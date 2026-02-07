// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_asterisk_telephony

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/rapidaai/api/assistant-api/config"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

// asteriskTelephony implements the Telephony interface for Asterisk
type asteriskTelephony struct {
	appCfg *config.AssistantConfig
	logger commons.Logger
}

// NewAsteriskTelephony creates a new Asterisk telephony provider
func NewAsteriskTelephony(config *config.AssistantConfig, logger commons.Logger) (internal_type.Telephony, error) {
	return &asteriskTelephony{
		appCfg: config,
		logger: logger,
	}, nil
}

// StatusCallback handles status callback events from Asterisk
func (apt *asteriskTelephony) StatusCallback(
	c *gin.Context,
	auth types.SimplePrinciple,
	assistantId uint64,
	assistantConversationId uint64,
) ([]types.Telemetry, error) {
	body, err := c.GetRawData()
	if err != nil {
		apt.logger.Errorf("Failed to read event body: %+v", err)
		return nil, fmt.Errorf("failed to read request body")
	}

	// Try to parse as JSON first
	var eventDetails map[string]interface{}
	if err := json.Unmarshal(body, &eventDetails); err != nil {
		// Fall back to form-encoded data
		values, err := url.ParseQuery(string(body))
		if err != nil {
			apt.logger.Errorf("Failed to parse body: %+v", err)
			return nil, fmt.Errorf("failed to parse request body")
		}
		eventDetails = make(map[string]interface{})
		for key, value := range values {
			if len(value) > 0 {
				eventDetails[key] = value[0]
			} else {
				eventDetails[key] = nil
			}
		}
	}

	// Extract event type from various possible fields
	eventType := "unknown"
	if v, ok := eventDetails["type"]; ok {
		eventType = fmt.Sprintf("%v", v)
	} else if v, ok := eventDetails["event"]; ok {
		eventType = fmt.Sprintf("%v", v)
	} else if v, ok := eventDetails["Event"]; ok {
		eventType = fmt.Sprintf("%v", v)
	}

	return []types.Telemetry{
		types.NewMetric("STATUS", eventType, utils.Ptr("Status of conversation")),
		types.NewEvent(eventType, eventDetails),
	}, nil
}

// CatchAllStatusCallback handles catch-all status callbacks
func (apt *asteriskTelephony) CatchAllStatusCallback(ctx *gin.Context) ([]types.Telemetry, error) {
	return nil, nil
}

// ReceiveCall handles incoming call webhooks from Asterisk
func (apt *asteriskTelephony) ReceiveCall(c *gin.Context) (*string, []types.Telemetry, error) {
	queryParams := make(map[string]string)
	telemetry := []types.Telemetry{}

	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}

	// Try to get caller info from query params or body
	var clientNumber string

	// Check query parameters first
	if from, ok := queryParams["from"]; ok && from != "" {
		clientNumber = from
	} else if callerID, ok := queryParams["callerid"]; ok && callerID != "" {
		clientNumber = callerID
	} else if caller, ok := queryParams["caller"]; ok && caller != "" {
		clientNumber = caller
	}

	// Try to parse body if no caller found in query
	if clientNumber == "" {
		body, err := c.GetRawData()
		if err == nil && len(body) > 0 {
			var bodyData map[string]interface{}
			if json.Unmarshal(body, &bodyData) == nil {
				if from, ok := bodyData["from"]; ok {
					clientNumber = fmt.Sprintf("%v", from)
				} else if callerID, ok := bodyData["callerid"]; ok {
					clientNumber = fmt.Sprintf("%v", callerID)
				} else if caller, ok := bodyData["caller"]; ok {
					clientNumber = fmt.Sprintf("%v", caller)
				}
			}
		}
	}

	if clientNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing caller information"})
		return nil, telemetry, fmt.Errorf("missing caller information")
	}

	// Extract channel ID if available
	if channelID, ok := queryParams["channel_id"]; ok && channelID != "" {
		telemetry = append(telemetry, types.NewMetadata("telephony.uuid", channelID))
	} else if channelID, ok := queryParams["channelid"]; ok && channelID != "" {
		telemetry = append(telemetry, types.NewMetadata("telephony.uuid", channelID))
	}

	return utils.Ptr(clientNumber), append(telemetry,
		types.NewEvent("webhook", queryParams),
		types.NewMetric("STATUS", "SUCCESS", utils.Ptr("Status of telephony api")),
	), nil
}

// OutboundCall initiates an outbound call via Asterisk ARI
func (apt *asteriskTelephony) OutboundCall(
	auth types.SimplePrinciple,
	toPhone string,
	fromPhone string,
	assistantId, assistantConversationId uint64,
	vaultCredential *protos.VaultCredential,
	opts utils.Option,
) ([]types.Telemetry, error) {
	mtds := []types.Telemetry{
		types.NewMetadata("telephony.toPhone", toPhone),
		types.NewMetadata("telephony.fromPhone", fromPhone),
		types.NewMetadata("telephony.provider", "asterisk"),
	}

	if vaultCredential == nil {
		return append(mtds,
			types.NewMetadata("telephony.error", "Missing vault credential for Asterisk ARI"),
			types.NewMetric("STATUS", "FAILED", utils.Ptr("Status of telephony api")),
		), fmt.Errorf("missing vault credential for Asterisk ARI")
	}
	credMap := vaultCredential.GetValue().AsMap()
	ariURL, _ := credMap["ari_url"].(string)
	ariURL = fmt.Sprintf("%s/ari/channels", ariURL)
	params := url.Values{}
	params.Set("endpoint", fmt.Sprintf("%s", toPhone))

	if ctxVal, err := opts.GetString("context"); err == nil && ctxVal != "" {
		params.Set("context", ctxVal)
	}

	if extVal, err := opts.GetString("extension"); err == nil && extVal != "" {
		params.Set("extension", extVal)
	}
	callerId := fromPhone
	if callerIdVal, err := opts.GetString("caller_id"); err == nil && callerIdVal != "" {
		callerId = callerIdVal
	}

	params.Set("priority", "1")
	params.Set("callerId", callerId)
	params.Set("appArgs", fmt.Sprintf("incoming,assistant_id=%d,conversation_id=%d", assistantId, assistantConversationId))

	// Add channel variables for AudioSocket routing
	params.Add("variables", fmt.Sprintf("RAPIDA_ASSISTANT_ID=%d", assistantId))
	params.Add("variables", fmt.Sprintf("RAPIDA_CONVERSATION_ID=%d", assistantConversationId))
	if auth.GetCurrentProjectId() != nil {
		params.Add("variables", fmt.Sprintf("RAPIDA_PROJECT_ID=%d", *auth.GetCurrentProjectId()))
	}
	if auth.GetCurrentOrganizationId() != nil {
		params.Add("variables", fmt.Sprintf("RAPIDA_ORG_ID=%d", *auth.GetCurrentOrganizationId()))
	}

	// Create HTTP request
	reqURL := fmt.Sprintf("%s?%s", ariURL, params.Encode())
	req, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		return append(mtds,
			types.NewMetadata("telephony.error", fmt.Sprintf("request creation error: %s", err.Error())),
			types.NewMetric("STATUS", "FAILED", utils.Ptr("Status of telephony api")),
		), err
	}

	user, _ := credMap["ari_user"].(string)
	password, _ := credMap["ari_password"].(string)
	// Set authentication
	req.SetBasicAuth(user, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return append(mtds,
			types.NewMetadata("telephony.error", fmt.Sprintf("API error: %s", err.Error())),
			types.NewMetric("STATUS", "FAILED", utils.Ptr("Status of telephony api")),
		), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return append(mtds,
			types.NewMetadata("telephony.error", fmt.Sprintf("API returned status: %d", resp.StatusCode)),
			types.NewMetric("STATUS", "FAILED", utils.Ptr("Status of telephony api")),
		), fmt.Errorf("ARI API returned status: %d", resp.StatusCode)
	}

	// Parse response
	var ariResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&ariResp); err != nil {
		apt.logger.Warn("Failed to decode ARI response", "error", err)
	}

	channelID := ""
	if id, ok := ariResp["id"]; ok {
		channelID = fmt.Sprintf("%v", id)
	}

	return append(mtds,
		types.NewMetadata("telephony.uuid", channelID),
		types.NewEvent("channel_created", ariResp),
		types.NewMetric("STATUS", "SUCCESS", utils.Ptr("Status of telephony api")),
	), nil
}

// InboundCall handles inbound call setup for Asterisk
// This returns the WebSocket connection URL as plain text for Asterisk to connect to
// Asterisk AudioSocket or chan_websocket expects a plain WebSocket URL
func (apt *asteriskTelephony) InboundCall(
	c *gin.Context,
	auth types.SimplePrinciple,
	assistantId uint64,
	clientNumber string,
	assistantConversationId uint64,
) error {
	// Build WebSocket connection URL for Asterisk
	wsPath := internal_type.GetAnswerPrefix(auth, assistantId, assistantConversationId, clientNumber)
	// wsURL := fmt.Sprintf("wss://%s/%s", apt.appCfg.PublicAssistantHost, wsPath)

	// Return plain text WebSocket URL for Asterisk dialplan to use
	// This can be used directly in Asterisk dialplan with AudioSocket or chan_websocket:
	// same = n,AudioSocket(${CURL(https://api.rapida.ai/v1/talk/asterisk/call/${ASSISTANT_ID})})
	c.String(http.StatusOK, wsPath)
	return nil
}

// parsePort parses a port number from string
func parsePort(s string) (int, error) {
	var port int
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, fmt.Errorf("invalid port: %s", s)
		}
		port = port*10 + int(c-'0')
	}
	return port, nil
}
