// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package integration_client_builders

import (
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"

	gorm_types "github.com/rapidaai/pkg/models/gorm/types"
	"github.com/rapidaai/protos"
)

type InputChatBuilder interface {
	Credential(i uint64, dp *structpb.Struct) *protos.Credential
	Chat(
		requestId string,
		credential *protos.Credential,
		modelOpts map[string]*anypb.Any,
		tools []*protos.FunctionDefinition,
		additionalData map[string]string,
		conversations ...*protos.Message,
	) *protos.ChatRequest

	Message(
		templates []*gorm_types.PromptTemplate,
		arguments map[string]interface{},
	) []*protos.Message

	Options(
		opts map[string]interface{},
		options map[string]*anypb.Any,
	) map[string]*anypb.Any

	Arguments(
		variables []*gorm_types.PromptVariable,
		arguments map[string]*anypb.Any,
	) map[string]interface{}

	PromptArguments(
		variables []*gorm_types.PromptVariable,
	) map[string]interface{}
}

type InputEmbeddingBuilder interface {
	Credential(i uint64, dp *structpb.Struct) *protos.Credential
	Embedding(
		credential *protos.Credential,
		modelOpts map[string]*anypb.Any,
		additionalData map[string]string,
		contents map[int32]string,
	) *protos.EmbeddingRequest
	Options(
		opts map[string]interface{},
		options map[string]*anypb.Any) map[string]*anypb.Any
}

type InputRerankingBuilder interface {
	Credential(i uint64, dp *structpb.Struct) *protos.Credential
	Reranking(
		credential *protos.Credential,
		modelOpts map[string]*anypb.Any,
		additionalData map[string]string,
		contents map[int32]string,
	) *protos.RerankingRequest
	Options(
		opts map[string]interface{},
		options map[string]*anypb.Any) map[string]*anypb.Any
}
