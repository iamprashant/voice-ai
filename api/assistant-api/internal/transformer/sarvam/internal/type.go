// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.

package sarvam_internal

type SpeechToTextTranscriptionData struct {
	RequestID  string `json:"request_id"`
	Transcript string `json:"transcript"`
	Metrics    struct {
		AudioDuration     float64 `json:"audio_duration"`
		ProcessingLatency float64 `json:"processing_latency"`
	} `json:"metrics"`
	Timestamps         interface{} `json:"timestamps,omitempty"`
	DiarizedTranscript interface{} `json:"diarized_transcript,omitempty"`
	LanguageCode       string      `json:"language_code,omitempty"`
}
type ErrorData struct {
	Error string `json:"error"` // The error message
	Code  string `json:"code"`  // The error code
}

type EventsData struct {
	EventType  string  `json:"event_type,omitempty"`  // Optional: Type of event
	Timestamp  string  `json:"timestamp,omitempty"`   // Optional: Timestamp of the event
	SignalType string  `json:"signal_type,omitempty"` // Optional: Voice Activity Detection (VAD) signal type, e.g., "START_SPEECH", "END_SPEECH"
	OccurredAt float64 `json:"occurred_at,omitempty"` // Optional: Epoch timestamp when the event occurred
}
