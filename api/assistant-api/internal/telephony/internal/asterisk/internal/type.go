// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_asterisk

import "encoding/json"

// AsteriskMediaEvent represents media events from Asterisk WebSocket
// Based on chan_websocket protocol from Asterisk
type AsteriskMediaEvent struct {
	// Event type: MEDIA_START, MEDIA_STOP, MEDIA_XON, MEDIA_XOFF, MEDIA_BUFFERING_COMPLETED, etc.
	Event string `json:"event,omitempty"`

	// Command type for JSON mode: START_MEDIA_BUFFERING, STOP_MEDIA_BUFFERING, MARK_MEDIA, HANGUP
	Command string `json:"command,omitempty"`

	// Channel name (populated in MEDIA_START)
	Channel string `json:"channel,omitempty"`

	// Optimal frame size for audio chunks
	OptimalFrameSize int `json:"optimal_frame_size,omitempty"`

	// Correlation ID for tracking
	CorrelationID string `json:"correlation_id,omitempty"`

	// Raw message for non-JSON text messages
	RawMessage string `json:"-"`
}

// ParseAsteriskEvent parses an Asterisk event from string or JSON
func ParseAsteriskEvent(data string) (*AsteriskMediaEvent, error) {
	event := &AsteriskMediaEvent{}

	// Try JSON parsing first
	if err := json.Unmarshal([]byte(data), event); err == nil && (event.Event != "" || event.Command != "") {
		return event, nil
	}

	// Parse legacy text format: "EVENT_TYPE key:value key:value ..."
	event.RawMessage = data
	event.Event = parseEventType(data)

	// Parse key-value pairs
	params := parseKeyValuePairs(data)
	if v, ok := params["channel"]; ok {
		event.Channel = v
	}
	if v, ok := params["optimal_frame_size"]; ok {
		var size int
		if _, err := parseIntFromString(v, &size); err == nil {
			event.OptimalFrameSize = size
		}
	}

	return event, nil
}

// parseEventType extracts the event type from a space-separated message
func parseEventType(data string) string {
	for i, c := range data {
		if c == ' ' {
			return data[:i]
		}
	}
	return data
}

// parseKeyValuePairs parses "key:value" pairs from a space-separated string
func parseKeyValuePairs(data string) map[string]string {
	result := make(map[string]string)
	parts := splitBySpace(data)

	for _, part := range parts[1:] { // Skip the event type
		for i, c := range part {
			if c == ':' {
				key := part[:i]
				value := part[i+1:]
				result[key] = value
				break
			}
		}
	}

	return result
}

// splitBySpace splits a string by spaces
func splitBySpace(data string) []string {
	var result []string
	var current string

	for _, c := range data {
		if c == ' ' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}

	if current != "" {
		result = append(result, current)
	}

	return result
}

// parseIntFromString parses an integer from a string
func parseIntFromString(s string, out *int) (bool, error) {
	var value int
	for _, c := range s {
		if c < '0' || c > '9' {
			return false, nil
		}
		value = value*10 + int(c-'0')
	}
	*out = value
	return true, nil
}

// AsteriskARIEvent represents an ARI event from Asterisk
type AsteriskARIEvent struct {
	Type      string                 `json:"type"`
	Timestamp string                 `json:"timestamp"`
	RequestID string                 `json:"request_id,omitempty"`
	Channel   *AsteriskChannel       `json:"channel,omitempty"`
	Bridge    *AsteriskBridge        `json:"bridge,omitempty"`
	Peer      *AsteriskChannel       `json:"peer,omitempty"`
	Extra     map[string]interface{} `json:"-"`
}

// AsteriskChannel represents a channel in ARI
type AsteriskChannel struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	State       string            `json:"state"`
	Caller      *AsteriskEndpoint `json:"caller,omitempty"`
	Connected   *AsteriskEndpoint `json:"connected,omitempty"`
	Dialplan    *AsteriskDialplan `json:"dialplan,omitempty"`
	ChannelVars map[string]string `json:"channelvars,omitempty"`
}

// AsteriskEndpoint represents caller/connected endpoint info
type AsteriskEndpoint struct {
	Name   string `json:"name"`
	Number string `json:"number"`
}

// AsteriskDialplan represents dialplan context
type AsteriskDialplan struct {
	Context string `json:"context"`
	Exten   string `json:"exten"`
	AppData string `json:"app_data"`
}

// AsteriskBridge represents a bridge in ARI
type AsteriskBridge struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	BridgeType string   `json:"bridge_type"`
	Channels   []string `json:"channels"`
}

// AsteriskRESTRequest represents a REST request over WebSocket
type AsteriskRESTRequest struct {
	Type         string              `json:"type"`
	RequestID    string              `json:"request_id"`
	Method       string              `json:"method"`
	URI          string              `json:"uri"`
	QueryStrings []map[string]string `json:"query_strings,omitempty"`
}

// AsteriskRESTResponse represents a REST response over WebSocket
type AsteriskRESTResponse struct {
	Type         string `json:"type"`
	RequestID    string `json:"request_id"`
	StatusCode   int    `json:"status_code"`
	ReasonPhrase string `json:"reason_phrase"`
	MessageBody  string `json:"message_body,omitempty"`
}

// AsteriskConfig holds configuration for Asterisk connection
// This configuration is extracted from vault credentials
//
// MINIMAL Vault Credential (for outbound calls only):
//
//	{
//	  "ari_host": "asterisk.example.com",  // REQUIRED - Asterisk server hostname
//	  "ari_user": "asterisk",              // REQUIRED - ARI username
//	  "ari_password": "secret"             // REQUIRED - ARI password
//	}
//
// NOTE: For inbound calls, NO vault config is needed!
// Asterisk connects directly to your WebSocket endpoint.
//
// Optional fields (with defaults):
//
//	{
//	  "ari_port": 8088,              // ARI HTTP port (default: 8088)
//	  "ari_scheme": "http",          // http or https (default: http)
//	  "ari_app": "rapida",           // Stasis app name (default: rapida)
//	  "sip_endpoint": "PJSIP",       // SIP tech: PJSIP or SIP (default: PJSIP)
//	  "sip_context": "from-internal" // Dialplan context (default: from-internal)
//	}
//
// Audio codec (ulaw 8kHz) is hardcoded - no configuration needed.
type AsteriskConfig struct {
	// ARI REST API settings (required for outbound calls)
	ARIHost     string `mapstructure:"ari_host"`     // REQUIRED
	ARIUser     string `mapstructure:"ari_user"`     // REQUIRED
	ARIPassword string `mapstructure:"ari_password"` // REQUIRED

	// ARI optional settings (have defaults)
	ARIPort   int    `mapstructure:"ari_port"`   // default: 8088
	ARIScheme string `mapstructure:"ari_scheme"` // default: http
	ARIApp    string `mapstructure:"ari_app"`    // default: rapida

	// SIP settings (have defaults)
	SIPEndpoint string `mapstructure:"sip_endpoint"` // default: PJSIP
	SIPContext  string `mapstructure:"sip_context"`  // default: from-internal
}

// DefaultAsteriskConfig returns default configuration
// Only ARIHost, ARIUser, ARIPassword need to be set from vault
func DefaultAsteriskConfig() *AsteriskConfig {
	return &AsteriskConfig{
		ARIPort:     8088,
		ARIScheme:   "http",
		ARIApp:      "rapida",
		SIPEndpoint: "PJSIP",
		SIPContext:  "from-internal",
	}
}
