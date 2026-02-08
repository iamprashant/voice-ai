// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_sip_telephony

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rapidaai/api/assistant-api/config"
	internal_type "github.com/rapidaai/api/assistant-api/internal/type"
	sip_infra "github.com/rapidaai/api/assistant-api/sip/infra"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

// sipTelephony implements the Telephony interface for native SIP
type sipTelephony struct {
	appCfg       *config.AssistantConfig
	logger       commons.Logger
	sharedServer *sip_infra.Server // Shared SIP server for outbound calls (injected from SIPManager)
}

// NewSIPTelephony creates a new SIP telephony provider.
// sipServer is the shared SIP server instance from SIPManager used for outbound calls.
func NewSIPTelephony(cfg *config.AssistantConfig, logger commons.Logger, sipServer *sip_infra.Server) (internal_type.Telephony, error) {
	return &sipTelephony{
		appCfg:       cfg,
		logger:       logger,
		sharedServer: sipServer,
	}, nil
}

// parseConfig parses SIP provider credentials from vault, then overlays
// platform operational settings (port, transport, RTP range) from app config.
// Twilio/providers only give: sip_uri, sip_username, sip_password, sip_realm, sip_domain
// Our platform provides: port, transport, rtp_port_range from app config
func (t *sipTelephony) parseConfig(vaultCredential *protos.VaultCredential) (*sip_infra.Config, error) {
	if vaultCredential == nil || vaultCredential.GetValue() == nil {
		return nil, fmt.Errorf("vault credential is required")
	}

	credMap := vaultCredential.GetValue().AsMap()
	cfg := &sip_infra.Config{}

	// --- Provider credentials (from vault / Twilio) ---

	// Parse sip_uri to extract server and port (e.g. "sip:192.168.1.5:5060")
	if sipURI, ok := credMap["sip_uri"].(string); ok && sipURI != "" {
		uri := strings.TrimPrefix(strings.TrimPrefix(sipURI, "sips:"), "sip:")
		host, portStr, err := net.SplitHostPort(uri)
		if err != nil {
			// No port in URI, treat entire string as host
			cfg.Server = uri
		} else {
			cfg.Server = host
			if p, err := strconv.Atoi(portStr); err == nil {
				cfg.Port = p
			}
		}
	}

	// Explicit sip_server overrides sip_uri
	if server, ok := credMap["sip_server"].(string); ok && server != "" {
		cfg.Server = server
	}
	if username, ok := credMap["sip_username"].(string); ok {
		cfg.Username = username
	}
	if password, ok := credMap["sip_password"].(string); ok {
		cfg.Password = password
	}
	if realm, ok := credMap["sip_realm"].(string); ok {
		cfg.Realm = realm
	}
	if domain, ok := credMap["sip_domain"].(string); ok {
		cfg.Domain = domain
	}

	// --- Platform operational settings (from app config) ---
	if t.appCfg.SIPConfig != nil {
		cfg.ApplyOperationalDefaults(
			t.appCfg.SIPConfig.Port,
			sip_infra.Transport(t.appCfg.SIPConfig.Transport),
			t.appCfg.SIPConfig.RTPPortRangeStart,
			t.appCfg.SIPConfig.RTPPortRangeEnd,
		)
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// StatusCallback handles status callbacks from SIP events
func (t *sipTelephony) StatusCallback(
	c *gin.Context,
	auth types.SimplePrinciple,
	assistantId uint64,
	assistantConversationId uint64,
) ([]types.Telemetry, error) {
	body, err := c.GetRawData()
	if err != nil {
		t.logger.Error("Failed to read SIP status callback body", "error", err)
		return nil, fmt.Errorf("failed to read request body")
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.logger.Error("Failed to parse SIP status callback", "error", err)
		return nil, fmt.Errorf("failed to parse request body")
	}

	// Extract event type
	eventType, _ := payload["event"].(string)
	callID, _ := payload["call_id"].(string)

	t.logger.Debug("SIP status callback received",
		"event", eventType,
		"call_id", callID,
		"assistant_id", assistantId,
		"conversation_id", assistantConversationId)

	return []types.Telemetry{
		types.NewMetric("STATUS", eventType, utils.Ptr("SIP event status")),
		types.NewEvent(eventType, payload),
	}, nil
}

// CatchAllStatusCallback handles catch-all status callbacks
func (t *sipTelephony) CatchAllStatusCallback(ctx *gin.Context) ([]types.Telemetry, error) {
	return nil, nil
}

// OutboundCall initiates an outbound SIP call
func (t *sipTelephony) OutboundCall(
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
		types.NewMetadata("telephony.provider", "sip"),
	}

	cfg, err := t.parseConfig(vaultCredential)
	if err != nil {
		return append(mtds,
			types.NewMetadata("telephony.error", fmt.Sprintf("config error: %s", err.Error())),
			types.NewMetric("STATUS", "FAILED", utils.Ptr("Status of telephony api")),
		), err
	}

	// Validate shared server is available and running
	if t.sharedServer == nil {
		return append(mtds,
			types.NewMetadata("telephony.error", "SIP server not initialized"),
			types.NewMetric("STATUS", "FAILED", utils.Ptr("Status of telephony api")),
		), fmt.Errorf("shared SIP server not available")
	}
	if !t.sharedServer.IsRunning() {
		return append(mtds,
			types.NewMetadata("telephony.error", "SIP server not running"),
			types.NewMetric("STATUS", "FAILED", utils.Ptr("Status of telephony api")),
		), fmt.Errorf("shared SIP server is not running")
	}

	// Initiate outbound call via the shared SIP server.
	// Pass metadata upfront so it is set on the session BEFORE the
	// handleOutboundDialog goroutine starts. On fast LANs the 200 OK
	// can arrive before MakeCall returns, causing handleOutboundAnswered
	// to fail with "outbound session missing assistant_id metadata".
	callMetadata := map[string]interface{}{
		"assistant_id":    assistantId,
		"conversation_id": assistantConversationId,
		"to_phone":        toPhone,
		"auth":            auth,
		"sip_config":      cfg,
	}
	session, err := t.sharedServer.MakeCall(context.Background(), cfg, toPhone, fromPhone, callMetadata)
	if err != nil {
		return append(mtds,
			types.NewMetadata("telephony.error", fmt.Sprintf("call error: %s", err.Error())),
			types.NewMetric("STATUS", "FAILED", utils.Ptr("Status of telephony api")),
		), err
	}

	t.logger.Info("SIP outbound call initiated",
		"to", toPhone,
		"from", fromPhone,
		"call_id", session.GetCallID(),
		"assistant_id", assistantId,
		"conversation_id", assistantConversationId)

	return append(mtds,
		types.NewMetadata("telephony.status", "initiated"),
		types.NewMetadata("telephony.call_id", session.GetCallID()),
		types.NewEvent("initiated", map[string]interface{}{
			"to":              toPhone,
			"from":            fromPhone,
			"call_id":         session.GetCallID(),
			"assistant_id":    assistantId,
			"conversation_id": assistantConversationId,
		}),
		types.NewMetric("STATUS", "SUCCESS", utils.Ptr("Status of telephony api")),
	), nil
}

// InboundCall handles incoming SIP calls
func (t *sipTelephony) InboundCall(
	c *gin.Context,
	auth types.SimplePrinciple,
	assistantId uint64,
	clientNumber string,
	assistantConversationId uint64,
) error {
	// For native SIP, inbound calls are handled directly by the SIP server
	// This endpoint just returns a confirmation
	c.JSON(http.StatusOK, gin.H{
		"status":          "ready",
		"assistant_id":    assistantId,
		"conversation_id": assistantConversationId,
		"client_number":   clientNumber,
		"message":         "SIP inbound call ready - connect via SIP signaling",
	})
	return nil
}

// ReceiveCall processes incoming call webhook data
func (t *sipTelephony) ReceiveCall(c *gin.Context) (*string, []types.Telemetry, error) {
	queryParams := make(map[string]string)
	telemetry := []types.Telemetry{}

	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			queryParams[key] = values[0]
		}
	}

	// Extract caller information
	clientNumber, ok := queryParams["from"]
	if !ok || clientNumber == "" {
		// Try alternative parameter names
		clientNumber, ok = queryParams["caller"]
		if !ok || clientNumber == "" {
			return nil, telemetry, fmt.Errorf("missing caller information")
		}
	}

	if callID, ok := queryParams["call_id"]; ok && callID != "" {
		telemetry = append(telemetry, types.NewMetadata("telephony.uuid", callID))
	}

	return utils.Ptr(clientNumber), append(telemetry,
		types.NewEvent("webhook", queryParams),
		types.NewMetric("STATUS", "SUCCESS", utils.Ptr("Status of telephony api")),
	), nil
}
