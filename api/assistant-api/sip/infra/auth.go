// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package sip_infra

import (
	"fmt"
	"strings"
)

// SipURICredentials holds parsed credentials from a SIP URI.
//
// URI format: sip:{assistantID}:{apiKey}@aws.ap-south-east-01.rapida.ai
//
// The assistantID and apiKey together form a valid authentication token
// for all SIP requests — not just INVITE.
type SipURICredentials struct {
	AssistantID string
	APIKey      string
}

// ParseCredentialsFromURI extracts assistant ID and API key from a SIP URI.
// Exported so middleware and tests can call it without a Server instance.
//
// URI format: sip:{assistantID}:{apiKey}@aws.ap-south-east-01.rapida.ai
//
// The userinfo part ({assistantID}:{apiKey}) is the authentication credential.
// The host part is the SIP server routing domain and is ignored for auth purposes.
func ParseCredentialsFromURI(uri string) (*SipURICredentials, error) {
	// Remove sip:/sips: prefix
	uri = strings.TrimPrefix(uri, "sip:")
	uri = strings.TrimPrefix(uri, "sips:")
	// Split user@host
	parts := strings.SplitN(uri, "@", 2)
	if len(parts) == 0 || parts[0] == "" {
		return nil, fmt.Errorf("invalid SIP URI: missing userinfo")
	}
	userinfo := parts[0]

	// Parse {assistantID}:{apiKey} from userinfo
	colonIdx := strings.Index(userinfo, ":")
	if colonIdx <= 0 {
		return nil, fmt.Errorf("invalid SIP URI: expected {assistantID}:{apiKey}@host, got %q", uri)
	}

	assistantID := userinfo[:colonIdx]
	apiKey := userinfo[colonIdx+1:]

	if apiKey == "" {
		return nil, fmt.Errorf("invalid SIP URI: missing API key in %q", uri)
	}

	// Strip legacy suffixes from assistant ID if present
	if idx := strings.Index(assistantID, ".rapid-sip"); idx > 0 {
		assistantID = assistantID[:idx]
	} else if idx := strings.Index(assistantID, ".rapida"); idx > 0 {
		assistantID = assistantID[:idx]
	}

	return &SipURICredentials{
		AssistantID: assistantID,
		APIKey:      apiKey,
	}, nil
}

// CredentialMiddleware is the first middleware in the SIP authentication chain.
// It parses the SIP URI from the To header to extract the assistantID and apiKey,
// then sets them on the SIPRequestContext for downstream middlewares.
//
// URI format: sip:{assistantID}:{apiKey}@aws.ap-south-east-01.rapida.ai
//
// This middleware runs for ALL SIP requests (INVITE, BYE, REGISTER, OPTIONS, etc.)
// to ensure every request carries valid authentication credentials.
func CredentialMiddleware(ctx *SIPRequestContext, next func() (*InviteResult, error)) (*InviteResult, error) {
	// Parse credentials from To URI — the To header contains the target identity
	// which embeds the authentication: sip:{assistantID}:{apiKey}@host
	creds, err := ParseCredentialsFromURI(ctx.ToURI)
	if err != nil {
		return Reject(400, fmt.Sprintf("Invalid SIP URI: %v", err)), nil
	}

	ctx.AssistantID = creds.AssistantID
	ctx.APIKey = creds.APIKey

	return next()
}
