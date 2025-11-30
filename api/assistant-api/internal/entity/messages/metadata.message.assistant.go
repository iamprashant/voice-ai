package internal_message_gorm

import (
	gorm_model "github.com/rapidaai/pkg/models/gorm"
)

type AssistantConversationMessageMetadata struct {
	gorm_model.Audited
	gorm_model.Mutable
	gorm_model.Metadata
	AssistantConversationId        uint64 `json:"assistantConversationId" gorm:"type:bigint;not null"`
	AssistantConversationMessageId string `json:"messageId" gorm:"type:string;not null"`
}
