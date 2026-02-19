// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package internal_transformer_resemble

import (
	"fmt"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

const (
	RESEMBLE_URL = "f.cluster.resemble.ai/stream"
	VOICE_ID     = "1dcf0222"
)

type resembleOption struct {
	logger    commons.Logger
	modelOpts utils.Option
	key       string
	projectId string
}

func NewResembleOption(logger commons.Logger,
	vaultCredential *protos.VaultCredential, option utils.Option) (*resembleOption, error) {

	credentialsMap := vaultCredential.GetValue().AsMap()
	cx, ok := credentialsMap["key"]
	if !ok {
		return nil, fmt.Errorf("resemble: illegal vault config")
	}

	prj, ok := credentialsMap["project_id"]
	if !ok {
		return nil, fmt.Errorf("resemble: illegal vault config")
	}
	return &resembleOption{
		logger:    logger,
		modelOpts: option,
		key:       cx.(string),
		projectId: prj.(string),
	}, nil
}

func (ro *resembleOption) GetKey() string {
	return ro.key
}

func (ro *resembleOption) GetProject() string {
	return ro.projectId
}

func (ro *resembleOption) GetEncoding() string {
	return "PCM_16"
}

func (ro *resembleOption) GetTextToSpeechRequest(contextId, text string) map[string]interface{} {
	return map[string]interface{}{
		"voice_uuid":      VOICE_ID,
		"request_id":      contextId,
		"project_uuid":    ro.GetProject(),
		"data":            text,
		"binary_response": true,
		"precision":       ro.GetEncoding(),
		"sample_rate":     16000,
	}

}
