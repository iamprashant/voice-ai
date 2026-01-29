// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_type

import (
	"github.com/rapidaai/pkg/types"
	"github.com/rapidaai/protos"
)

// Packet represents a generic request packet handled by the adapter layer.
// Concrete packet types (e.g., FlushPacket, InterruptionPacket) are used to
// signal specific control actions within a given context.
type Packet interface {
	ContextId() string
}

// Wrapper for message packet
type MessagePacket interface {
	Packet
	Role() string
	Content() string
}

type AudioPacket interface {
	Packet
	Content() []byte
}

// InterruptionPacket represents a request to interrupt ongoing processing
// within a specific context.

// =============================================================================
// LLM Packets
// =============================================================================

type InterruptionSource string

const (
	InterruptionSourceWord InterruptionSource = "word"
	InterruptionSourceVad  InterruptionSource = "vad"
)

type InterruptionPacket struct {
	// ContextID identifies the context to be interrupted.
	ContextID string

	// Source indicates the origin of the interruption.
	Source InterruptionSource

	// start of interruption
	StartAt float64

	// end of interruption
	EndAt float64
}

// ContextId returns the identifier of the context associated with this interruption request.
func (f InterruptionPacket) ContextId() string {
	return f.ContextID
}

// =============================================================================
// LLM Packets
// =============================================================================

// MetricPacket represents a request to send metrics within a specific context.
type MetricPacket struct {
	// ContextID identifies the context to be flushed.
	ContextID string

	// Metrics holds the list of metrics to be sent within the specified context.
	Metrics []*types.Metric
}

func (f MetricPacket) ContextId() string {
	return f.ContextID
}

// =============================================================================
// Directive Packets
// =============================================================================

type DirectivePacket struct {
	// ContextID identifies the context to be flushed.
	ContextID string

	// Directive
	Directive protos.ConversationDirective_DirectiveType

	// arguments for directive
	Arguments map[string]interface{}
}

func (f DirectivePacket) ContextId() string {
	return f.ContextID
}

// =============================================================================
// LLM Packets
// =============================================================================

type LLMPacket interface {
	Packet
	ContextId() string
}

type LLMErrorPacket struct {
	// contextID identifies the context to be flushed.
	ContextID string

	// error
	Error error

	//

}

func (f LLMErrorPacket) ContextId() string {
	return f.ContextID
}

// LLMResponseDeltaPacket represents a streaming text delta from the LLM.
// This packet is emitted during streaming responses, containing partial text chunks.
type LLMResponseDeltaPacket struct {
	// ContextID identifies the context for this response.
	ContextID string

	// Text contains the partial text content of this delta.
	Text string
}

func (f LLMResponseDeltaPacket) ContextId() string {
	return f.ContextID
}

// LLMResponseDonePacket signals the completion of an LLM response stream.
// This packet is emitted when the LLM has finished generating its response.
type LLMResponseDonePacket struct {
	// ContextID identifies the context for this response.
	ContextID string

	// Text contains the final aggregated text (optional, may be empty for streaming).
	Text string
}

func (f LLMResponseDonePacket) Content() string {
	return f.Text
}

func (f LLMResponseDonePacket) Role() string {
	return "assistant"
}

func (f LLMResponseDonePacket) ContextId() string {
	return f.ContextID
}

// =============================================================================
// LLM Tool Call Packets
// =============================================================================

type LLMToolPacket interface {
	ToolId() string
}

type LLMToolCallPacket struct {
	// id of tool which user has configured
	ToolID string

	// name of tool which user has configured
	Name string

	// contextID identifies the context to be flushed.
	ContextID string

	// arguments for tool call
	Arguments map[string]interface{}
}

func (f LLMToolCallPacket) ContextId() string {
	return f.ContextID
}

func (f LLMToolCallPacket) ToolId() string {
	return f.ToolID
}

type LLMToolResultPacket struct {
	// id of tool which user has configured
	ToolID string

	// name of tool which user has configured
	Name string

	// contextID identifies the context to be flushed.
	ContextID string

	// result for tool call
	Result map[string]interface{}
}

func (f LLMToolResultPacket) ToolId() string {
	return f.ToolID
}

func (f LLMToolResultPacket) ContextId() string {
	return f.ContextID
}

// =============================================================================
// LLM Packets end
// =============================================================================

type StaticPacket struct {
	// contextID identifies the context to be flushed.
	ContextID string

	// message
	Text string
}

func (f StaticPacket) ContextId() string {
	return f.ContextID
}

func (f StaticPacket) Content() string {
	return f.Text
}

func (f StaticPacket) Role() string {
	return "rapida"
}

// =============================================================================
// LLM Packets end
// =============================================================================

type TextToSpeechAudioPacket struct {

	// contextID identifies the context to be flushed.
	ContextID string

	// audio chunk
	AudioChunk []byte
}

func (f TextToSpeechAudioPacket) ContextId() string {
	return f.ContextID
}

type TextToSpeechEndPacket struct {
	// contextID identifies the context to be flushed.
	ContextID string
}

func (f TextToSpeechEndPacket) ContextId() string {
	return f.ContextID
}

// =============================================================================
// User Packet
// =============================================================================

type UserTextPacket struct {
	// contextID identifies the context to be flushed.
	ContextID string

	// text
	Text string
}

func (f UserTextPacket) ContextId() string {
	return f.ContextID
}

func (f UserTextPacket) Content() string {
	return f.Text
}

func (f UserTextPacket) Role() string {
	return "user"
}

type UserAudioPacket struct {
	// contextID identifies the context to be flushed.
	ContextID string

	Audio []byte

	NoiseReduced bool
}

func (f UserAudioPacket) ContextId() string {
	return f.ContextID
}

func (f UserAudioPacket) Content() []byte {
	return f.Audio
}

func (f UserAudioPacket) Role() string {
	return "user"
}

// =============================================================================
// End of speech Packet
// =============================================================================

type EndOfSpeechPacket struct {
	// contextID identifies the context to be flushed.
	ContextID string

	Speech string
}

func (f EndOfSpeechPacket) ContextId() string {
	return f.ContextID
}

type SpeechToTextPacket struct {
	ContextID string

	// script
	Script string

	// confidence
	Confidence float64

	// language
	Language string

	// interim
	Interim bool
}

func (f SpeechToTextPacket) ContextId() string {
	return f.ContextID
}

//

// KnowledgeRetrieveOption contains options for knowledge retrieval operations
type KnowledgeRetrieveOption struct {
	EmbeddingProviderCredential *protos.VaultCredential
	RetrievalMethod             string
	TopK                        uint32
	ScoreThreshold              float32
}

type KnowledgeContextResult struct {
	ID         string                 `json:"id"`
	DocumentID string                 `json:"document_id"`
	Metadata   map[string]interface{} `json:"metadata"`
	Content    string                 `json:"content"`
	Score      float64                `json:"score"`
}
