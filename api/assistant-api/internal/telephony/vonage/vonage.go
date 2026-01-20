// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_vonage_telephony

import (
	"fmt"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/protos"
	vng "github.com/vonage/vonage-go-sdk"
)

type vg struct {
	logger commons.Logger
}

func NewVonage(logger commons.Logger) vg {
	return vg{
		logger: logger,
	}
}

func (vt vg) Auth(vaultCredential *protos.VaultCredential) (vng.Auth, error) {
	privateKey, ok := vaultCredential.GetValue().AsMap()["private_key"]
	if !ok {
		return nil, fmt.Errorf("illegal vault config privateKey is not found")
	}
	applicationId, ok := vaultCredential.GetValue().AsMap()["application_id"]
	if !ok {
		return nil, fmt.Errorf("illegal vault config application_id is not found")
	}
	clientAuth, err := vng.CreateAuthFromAppPrivateKey(applicationId.(string), []byte(privateKey.(string)))
	if err != nil {
		return nil, err
	}
	return clientAuth, nil
}
