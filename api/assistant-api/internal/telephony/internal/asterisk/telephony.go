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
	"github.com/gorilla/websocket"
	"github.com/rapidaai/api/assistant-api/config"
	internal_assistant_entity "github.com/rapidaai/api/assistant-api/internal/entity/assistants"
	internal_conversation_entity "github.com/rapidaai/api/assistant-api/internal/entity/conversations"
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

// Streamer creates a new WebSocket streamer for Asterisk chan_websocket
func (apt *asteriskTelephony) Streamer(
	c *gin.Context,
	connection *websocket.Conn,
	assistant *internal_assistant_entity.Assistant,
	conversation *internal_conversation_entity.AssistantConversation,
	vlt *protos.VaultCredential,
) internal_type.TelephonyStreamer {
	return NewAsteriskWebsocketStreamer(apt.logger, connection, assistant, conversation, vlt)
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

	// Get Asterisk ARI configuration from vault
	ariConfig, err := apt.getARIConfig(vaultCredential)
	if err != nil {
		return append(mtds,
			types.NewMetadata("telephony.error", fmt.Sprintf("configuration error: %s", err.Error())),
			types.NewMetric("STATUS", "FAILED", utils.Ptr("Status of telephony api")),
		), err
	}

	// Build ARI REST endpoint URL using scheme from vault config
	ariURL := fmt.Sprintf("%s://%s:%d/ari/channels", ariConfig.ARIScheme, ariConfig.ARIHost, ariConfig.ARIPort)

	// Build query parameters for channel origination
	// Uses SIP endpoint type and context from vault configuration
	params := url.Values{}
	params.Set("endpoint", fmt.Sprintf("%s/%s", ariConfig.SIPEndpoint, toPhone))
	params.Set("extension", toPhone)
	params.Set("context", ariConfig.SIPContext)
	params.Set("priority", "1")
	params.Set("callerId", fromPhone)
	params.Set("app", ariConfig.ARIApp)
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

	// Set authentication
	req.SetBasicAuth(ariConfig.ARIUser, ariConfig.ARIPassword)

	// Execute request
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
// This returns the WebSocket connection URL for Asterisk to connect to
func (apt *asteriskTelephony) InboundCall(
	c *gin.Context,
	auth types.SimplePrinciple,
	assistantId uint64,
	clientNumber string,
	assistantConversationId uint64,
) error {
	// Return WebSocket connection URL for Asterisk
	wsPath := internal_type.GetAnswerPath("asterisk", auth, assistantId, assistantConversationId, clientNumber)

	response := map[string]interface{}{
		"websocket_url": fmt.Sprintf("wss://%s/%s", apt.appCfg.PublicAssistantHost, wsPath),
		"protocol":      "media",
		"assistant_id":  assistantId,
		"client_number": clientNumber,
	}

	c.JSON(http.StatusOK, response)
	return nil
}

// getARIConfig extracts ARI configuration from vault credential
//
// MINIMAL vault credential (3 required fields):
//
//	{
//	  "ari_host": "asterisk.example.com",  // REQUIRED
//	  "ari_user": "asterisk",              // REQUIRED
//	  "ari_password": "secret"             // REQUIRED
//	}
//
// All other fields have sensible defaults:
//   - ari_port: 8088
//   - ari_scheme: http
//   - ari_app: rapida
//   - sip_endpoint: PJSIP
//   - sip_context: from-internal
func (apt *asteriskTelephony) getARIConfig(vaultCredential *protos.VaultCredential) (*ARIConfig, error) {
	if vaultCredential == nil {
		return nil, fmt.Errorf("vault credential is nil")
	}

	credMap := vaultCredential.GetValue().AsMap()

	// Start with defaults for optional fields
	config := &ARIConfig{
		ARIPort:     8088,            // default ARI port
		ARIScheme:   "http",          // default scheme
		ARIApp:      "rapida",        // default app name
		SIPEndpoint: "PJSIP",         // default SIP technology
		SIPContext:  "from-internal", // default dialplan context
	}

	// Parse REQUIRED fields
	if host, ok := credMap["ari_host"]; ok && host != nil {
		config.ARIHost = fmt.Sprintf("%v", host)
	}
	if user, ok := credMap["ari_user"]; ok && user != nil {
		config.ARIUser = fmt.Sprintf("%v", user)
	}
	if password, ok := credMap["ari_password"]; ok && password != nil {
		config.ARIPassword = fmt.Sprintf("%v", password)
	}

	// Parse ARI port
	if port, ok := credMap["ari_port"]; ok && port != nil {
		switch v := port.(type) {
		case float64:
			config.ARIPort = int(v)
		case int:
			config.ARIPort = v
		case string:
			if p, err := parsePort(v); err == nil {
				config.ARIPort = p
			}
		}
	}

	// Parse ARI scheme (http/https)
	if scheme, ok := credMap["ari_scheme"]; ok && scheme != nil {
		config.ARIScheme = fmt.Sprintf("%v", scheme)
	}

	// Parse ARI app name
	if app, ok := credMap["ari_app"]; ok && app != nil {
		config.ARIApp = fmt.Sprintf("%v", app)
	}

	// Parse SIP endpoint type
	if endpoint, ok := credMap["sip_endpoint"]; ok && endpoint != nil {
		config.SIPEndpoint = fmt.Sprintf("%v", endpoint)
	}

	// Parse SIP context
	if context, ok := credMap["sip_context"]; ok && context != nil {
		config.SIPContext = fmt.Sprintf("%v", context)
	}

	// Validate required fields (only 3 needed!)
	if config.ARIHost == "" || config.ARIHost == "localhost" {
		return nil, fmt.Errorf("ari_host is required in vault credential")
	}
	if config.ARIUser == "" {
		return nil, fmt.Errorf("ari_user is required in vault credential")
	}
	if config.ARIPassword == "" {
		return nil, fmt.Errorf("ari_password is required in vault credential")
	}

	return config, nil
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

// ARIConfig holds ARI configuration extracted from vault
type ARIConfig struct {
	ARIHost     string // Asterisk ARI host
	ARIPort     int    // Asterisk ARI port
	ARIScheme   string // http or https
	ARIApp      string // Stasis application name
	ARIUser     string // ARI username
	ARIPassword string // ARI password
	SIPEndpoint string // SIP endpoint type (PJSIP, SIP)
	SIPContext  string // Dialplan context
}
