// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_message_gorm

import (
	gorm_model "github.com/rapidaai/pkg/models/gorm"
)

type AssistantConversationMessage struct {
	gorm_model.Audited
	gorm_model.Mutable
	MessageId                string                                  `json:"messageId" gorm:"type:string;not null"`
	AssistantConversationId  uint64                                  `json:"assistantConversationId" gorm:"type:bigint;not null"`
	AssistantId              uint64                                  `json:"assistantId" gorm:"type:bigint;not null"`
	AssistantProviderModelId uint64                                  `json:"assistantProviderModelId" gorm:"type:bigint;not null"`
	Source                   string                                  `json:"source" gorm:"type:string;size:50;not null"`
	Role                     string                                  `json:"role" gorm:"type:string;size:50;not null"`
	Body                     string                                  `json:"body"`
	Metadatas                []*AssistantConversationMessageMetadata `json:"metadata" gorm:"foreignKey:AssistantConversationMessageId;references:MessageId"`
	Metrics                  []*AssistantConversationMessageMetric   `json:"metrics" gorm:"foreignKey:AssistantConversationMessageId;references:MessageId"`
}
