// Copyright (c) Rapida
// Author: Prashant <prashant@rapida.ai>
//
// Licensed under the Rapida internal use license.
// This file is part of Rapida's proprietary software.
// Unauthorized copying, modification, or redistribution is strictly prohibited.
package internal_adapter_request_customizers

import (
	internal_adapter_requests "github.com/rapidaai/api/assistant-api/internal/adapters/requests"
	"github.com/rapidaai/pkg/utils"
	protos "github.com/rapidaai/protos"
)

type customizer struct {
	args     map[string]interface{}
	options  map[string]interface{}
	metadata map[string]interface{}
}

func NewRequestBaseCustomizer(req *protos.AssistantConversationConfiguration) (internal_adapter_requests.Customization, error) {
	arg, err := utils.AnyMapToInterfaceMap(req.GetArgs())
	if err != nil {
		return nil, err
	}
	opts, err := utils.AnyMapToInterfaceMap(req.GetOptions())
	if err != nil {
		return nil, err
	}
	mtd, err := utils.AnyMapToInterfaceMap(req.GetMetadata())
	if err != nil {
		return nil, err
	}
	return &customizer{
		metadata: mtd,
		options:  opts,
		args:     arg,
	}, nil

}

func (ctmzr *customizer) GetMetadata() map[string]interface{} {
	return ctmzr.metadata
}

func (ctmzr *customizer) GetOptions() map[string]interface{} {
	return ctmzr.options
}

func (ctmzr *customizer) GetArgs() map[string]interface{} {
	return ctmzr.args
}
