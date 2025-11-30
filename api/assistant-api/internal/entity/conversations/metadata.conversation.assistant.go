package internal_conversation_gorm

import gorm "github.com/rapidaai/pkg/models/gorm"

type AssistantConversationMetadata struct {
	gorm.Audited
	gorm.Mutable
	gorm.Metadata
	AssistantId             uint64 `json:"assistantId" gorm:"type:bigint;not null"`
	AssistantConversationId uint64 `json:"assistantConversationId" gorm:"type:bigint;not null"`
}
