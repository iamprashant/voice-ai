// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_asterisk_telephony

import (
	"bytes"
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

const asteriskProvider = "asterisk"

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

// StatusCallback handles ARI status callback events from Asterisk.
// Asterisk ARI sends JSON payloads with a "type" field indicating the event type.
func (apt *asteriskTelephony) StatusCallback(
	c *gin.Context,
	auth types.SimplePrinciple,
	assistantId uint64,
	assistantConversationId uint64,
) (*internal_type.StatusInfo, error) {
	var eventDetails map[string]interface{}
	if err := c.ShouldBindJSON(&eventDetails); err != nil {
		apt.logger.Errorf("failed to parse ARI event body: %+v", err)
		return nil, fmt.Errorf("failed to parse ARI event body: %w", err)
	}

	eventType := "unknown"
	if v, ok := eventDetails["type"]; ok {
		eventType = fmt.Sprintf("%v", v)
	}

	return &internal_type.StatusInfo{Event: eventType, Payload: eventDetails}, nil
}

// CatchAllStatusCallback handles catch-all status callbacks
func (apt *asteriskTelephony) CatchAllStatusCallback(ctx *gin.Context) (*internal_type.StatusInfo, error) {
	return nil, nil
}

// ReceiveCall handles incoming call webhooks from Asterisk.
// The caller number is passed as the `from` query parameter by the Asterisk dialplan.
// Channel ID (if present) is captured as ChannelUUID.
func (apt *asteriskTelephony) ReceiveCall(c *gin.Context) (*internal_type.CallInfo, error) {
	clientNumber := c.Query("from")
	if clientNumber == "" {
		clientNumber = c.Query("callerid")
	}
	if clientNumber == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing caller information — provide 'from' query parameter"})
		return nil, fmt.Errorf("missing caller information in query params")
	}

	info := &internal_type.CallInfo{
		CallerNumber: clientNumber,
		Provider:     asteriskProvider,
		Status:       "SUCCESS",
		StatusInfo:   internal_type.StatusInfo{Event: "webhook", Payload: map[string]string{"from": clientNumber}},
	}
	if channelID := c.Query("channel_id"); channelID != "" {
		info.ChannelUUID = channelID
	}
	return info, nil
}

// OutboundCall initiates an outbound call via Asterisk ARI.
//
// ARI POST /ari/channels requires:
//   - endpoint: Technology/Resource (e.g. "PJSIP/+15551234567")
//   - app: Stasis application name (default "rapida")
//   - callerId: caller ID string
//   - context/extension/priority: dialplan entry point (optional, for non-Stasis routing)
//   - variables: JSON object of channel variables (passed in request body)
//
// The RAPIDA_CONTEXT_ID channel variable is critical — the Asterisk dialplan uses
// it as the AudioSocket UUID so the AudioSocket server can resolve the call context.
//
// Vault credential must contain:
//   - ari_url: base URL (e.g. "http://asterisk:8088")
//   - ari_user / ari_password: ARI authentication
//   - Optional: endpoint_technology (default "PJSIP"), trunk
//
// Deployment options may contain:
//   - context, extension: dialplan entry point for the outbound call
//   - caller_id: override callerId
//   - app: Stasis application name (default "rapida")
//   - endpoint_technology: SIP technology (default "PJSIP")
//   - trunk: SIP trunk name (e.g. "mytrunk")
func (apt *asteriskTelephony) OutboundCall(
	auth types.SimplePrinciple,
	toPhone string,
	fromPhone string,
	assistantId, assistantConversationId uint64,
	vaultCredential *protos.VaultCredential,
	opts utils.Option,
) (*internal_type.CallInfo, error) {
	info := &internal_type.CallInfo{Provider: asteriskProvider}

	if vaultCredential == nil {
		info.Status = "FAILED"
		info.ErrorMessage = "Missing vault credential for Asterisk ARI"
		return info, fmt.Errorf("missing vault credential for Asterisk ARI")
	}

	credMap := vaultCredential.GetValue().AsMap()
	ariBaseURL, _ := credMap["ari_url"].(string)
	if ariBaseURL == "" {
		info.Status = "FAILED"
		info.ErrorMessage = "Missing ari_url in vault credential"
		return info, fmt.Errorf("missing ari_url in vault credential")
	}

	// Build the endpoint in Technology/Resource format.
	// Default technology is PJSIP; can be overridden via vault credential or deployment opts.
	endpointTech := "PJSIP"
	if tech, ok := credMap["endpoint_technology"].(string); ok && tech != "" {
		endpointTech = tech
	}
	if tech, err := opts.GetString("endpoint_technology"); err == nil && tech != "" {
		endpointTech = tech
	}

	// Build endpoint: Technology/trunk/number or Technology/number
	endpoint := fmt.Sprintf("%s/%s", endpointTech, toPhone)
	if trunk, ok := credMap["trunk"].(string); ok && trunk != "" {
		endpoint = fmt.Sprintf("%s/%s/%s", endpointTech, trunk, toPhone)
	}
	if trunk, err := opts.GetString("trunk"); err == nil && trunk != "" {
		endpoint = fmt.Sprintf("%s/%s/%s", endpointTech, trunk, toPhone)
	}

	// Resolve caller ID
	callerId := fromPhone
	if callerIdVal, err := opts.GetString("caller_id"); err == nil && callerIdVal != "" {
		callerId = callerIdVal
	}

	// Build query parameters for ARI POST /ari/channels
	params := url.Values{}
	params.Set("endpoint", endpoint)
	params.Set("callerId", callerId)

	// Stasis app name — required by ARI to route the channel into a Stasis application.
	// If no context/extension is set, ARI requires "app" to handle the channel.
	// If context/extension IS set, ARI creates the channel in the dialplan instead.
	appName := "rapida"
	if appVal, err := opts.GetString("app"); err == nil && appVal != "" {
		appName = appVal
	}

	// Dialplan context/extension — if set, ARI originates the channel into the
	// given dialplan context (which should route to AudioSocket with RAPIDA_CONTEXT_ID).
	// If NOT set, ARI originates into the Stasis app.
	hasDialplan := false
	if ctxVal, err := opts.GetString("context"); err == nil && ctxVal != "" {
		params.Set("context", ctxVal)
		params.Set("priority", "1")
		hasDialplan = true
		if extVal, err := opts.GetString("extension"); err == nil && extVal != "" {
			params.Set("extension", extVal)
		} else {
			params.Set("extension", "s")
		}
	}

	if !hasDialplan {
		// No dialplan context — use Stasis app mode
		params.Set("app", appName)
		params.Set("appArgs", fmt.Sprintf("incoming,assistant_id=%d,conversation_id=%d", assistantId, assistantConversationId))
	}

	// Build channel variables as a JSON body.
	// ARI expects variables as {"variables": {"KEY": "VALUE"}} in the request body.
	// The RAPIDA_CONTEXT_ID variable is essential — the Asterisk dialplan uses it
	// as the AudioSocket UUID so the AudioSocket server can resolve the call context.
	channelVars := map[string]string{}
	if contextID, err := opts.GetString("rapida.context_id"); err == nil && contextID != "" {
		channelVars["RAPIDA_CONTEXT_ID"] = contextID
	}

	var bodyBytes []byte
	if len(channelVars) > 0 {
		bodyMap := map[string]interface{}{
			"variables": channelVars,
		}
		var err error
		bodyBytes, err = json.Marshal(bodyMap)
		if err != nil {
			info.Status = "FAILED"
			info.ErrorMessage = fmt.Sprintf("failed to marshal channel variables: %s", err.Error())
			return info, err
		}
	}

	// Create HTTP request
	ariURL := fmt.Sprintf("%s/ari/channels?%s", ariBaseURL, params.Encode())
	var req *http.Request
	var err error
	if bodyBytes != nil {
		req, err = http.NewRequest("POST", ariURL, bytes.NewReader(bodyBytes))
		if err == nil {
			req.Header.Set("Content-Type", "application/json")
		}
	} else {
		req, err = http.NewRequest("POST", ariURL, nil)
	}
	if err != nil {
		info.Status = "FAILED"
		info.ErrorMessage = fmt.Sprintf("request creation error: %s", err.Error())
		return info, err
	}

	user, _ := credMap["ari_user"].(string)
	password, _ := credMap["ari_password"].(string)
	req.SetBasicAuth(user, password)

	apt.logger.Infof("ARI outbound call: endpoint=%s, callerId=%s, url=%s", endpoint, callerId, ariURL)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		info.Status = "FAILED"
		info.ErrorMessage = fmt.Sprintf("ARI request error: %s", err.Error())
		return info, err
	}
	defer resp.Body.Close()

	// Parse response body (for both success and error cases)
	var ariResp map[string]interface{}
	if decodeErr := json.NewDecoder(resp.Body).Decode(&ariResp); decodeErr != nil {
		apt.logger.Warnf("failed to decode ARI response: %v", decodeErr)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		errMsg := fmt.Sprintf("ARI returned status %d", resp.StatusCode)
		if msg, ok := ariResp["message"]; ok {
			errMsg = fmt.Sprintf("ARI returned status %d: %v", resp.StatusCode, msg)
		}
		info.Status = "FAILED"
		info.ErrorMessage = errMsg
		apt.logger.Errorf("ARI outbound call failed: %s, response: %+v", errMsg, ariResp)
		return info, fmt.Errorf(errMsg)
	}

	if id, ok := ariResp["id"]; ok {
		info.ChannelUUID = fmt.Sprintf("%v", id)
	}

	info.Status = "SUCCESS"
	info.StatusInfo = internal_type.StatusInfo{Event: "channel_created", Payload: ariResp}
	apt.logger.Infof("ARI outbound call succeeded: channelId=%s, endpoint=%s", info.ChannelUUID, endpoint)
	return info, nil
}

// InboundCall handles inbound call setup for Asterisk.
// Returns the contextId as plain text — Asterisk dialplan uses this as the AudioSocket UUID
// so the AudioSocket server can resolve the full call context from Postgres.
//
// For AudioSocket: same = n,AudioSocket(${CURL(https://api.rapida.ai/v1/talk/asterisk/call/${ASSISTANT_ID})},host:port)
// For chan_websocket: the contextId is used in the WS URL path: wss://host/v1/talk/asterisk/ctx/${contextId}
func (apt *asteriskTelephony) InboundCall(
	c *gin.Context,
	auth types.SimplePrinciple,
	assistantId uint64,
	clientNumber string,
	assistantConversationId uint64,
) error {
	// contextId was set by CallReciever after saving the call context to Postgres.
	// Return it as plain text — Asterisk dialplan uses this as the AudioSocket UUID
	// or as part of the WebSocket URL path for chan_websocket.
	contextID, exists := c.Get("contextId")
	if !exists || contextID == "" {
		return fmt.Errorf("missing contextId — CallReciever must save call context before InboundCall")
	}
	c.String(http.StatusOK, fmt.Sprintf("%v", contextID))
	return nil
}
