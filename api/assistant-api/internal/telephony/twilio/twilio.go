// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_twilio_telephony

import (
	"fmt"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
	"github.com/twilio/twilio-go"
)

type twl struct {
	logger commons.Logger
}

func NewTwilio(logger commons.Logger) twl {
	return twl{
		logger: logger,
	}
}

func (tpc twl) Client(vaultCredential *protos.VaultCredential) (*twilio.RestClient, error) {
	clientParams, err := tpc.ClientParam(vaultCredential)
	if err != nil {
		return nil, err
	}
	return twilio.NewRestClientWithParams(*clientParams), nil
}

func (tpc twl) ClientParam(vaultCredential *protos.VaultCredential) (*twilio.ClientParams, error) {
	accountSid, ok := vaultCredential.GetValue().AsMap()["account_sid"]
	if !ok {
		return nil, fmt.Errorf("illegal vault config accountSid is not found")
	}
	authToken, ok := vaultCredential.GetValue().AsMap()["account_token"]
	if !ok {
		return nil, fmt.Errorf("illegal vault config account_token not found")
	}
	return &twilio.ClientParams{
		Username: accountSid.(string),
		Password: authToken.(string),
	}, nil
}
