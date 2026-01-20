// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package internal_exotel

type ExotelMediaEvent struct {
	Event     string       `json:"event"`
	StreamSid string       `json:"stream_sid"`
	Media     *ExotelMedia `json:"media,omitempty"`
}

type ExotelMedia struct {
	Payload string `json:"payload"`
}

type MakeCallResponse struct {
	Call struct {
		Sid              string  `json:"Sid"`
		Status           string  `json:"Status"`
		RecordingUrl     string  `json:"RecordingUrl"`
		ConversationUuid *string `json:"ParentCallSid"` // Use pointers for nullable fields
	} `json:"Call"`
}
