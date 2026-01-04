// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package commons

import (
	"testing"
)

func TestResponseContentTypeConstants(t *testing.T) {
	if TEXT_CONTENT != "text" {
		t.Errorf("TEXT_CONTENT = %v, want %v", TEXT_CONTENT, "text")
	}
	if AUDIO_CONTENT != "audio" {
		t.Errorf("AUDIO_CONTENT = %v, want %v", AUDIO_CONTENT, "audio")
	}
	if IMAGE_CONTENT != "image" {
		t.Errorf("IMAGE_CONTENT = %v, want %v", IMAGE_CONTENT, "image")
	}
	if MULTI_MEDIA_CONTENT != "multi" {
		t.Errorf("MULTI_MEDIA_CONTENT = %v, want %v", MULTI_MEDIA_CONTENT, "multi")
	}
}

func TestResponseContentFormatConstants(t *testing.T) {
	if TEXT_CONTENT_FORMAT_RAW != "raw" {
		t.Errorf("TEXT_CONTENT_FORMAT_RAW = %v, want %v", TEXT_CONTENT_FORMAT_RAW, "raw")
	}
	if AUDIO_CONTENT_FORMAT_RAW != "raw" {
		t.Errorf("AUDIO_CONTENT_FORMAT_RAW = %v, want %v", AUDIO_CONTENT_FORMAT_RAW, "raw")
	}
	if AUDIO_CONTENT_FORMAT_URL != "url" {
		t.Errorf("AUDIO_CONTENT_FORMAT_URL = %v, want %v", AUDIO_CONTENT_FORMAT_URL, "url")
	}
	if IMAGE_CONTENT_FORMAT_RAW != "raw" {
		t.Errorf("IMAGE_CONTENT_FORMAT_RAW = %v, want %v", IMAGE_CONTENT_FORMAT_RAW, "raw")
	}
	if IMAGE_CONTENT_FORMAT_URL != "url" {
		t.Errorf("IMAGE_CONTENT_FORMAT_URL = %v, want %v", IMAGE_CONTENT_FORMAT_URL, "url")
	}
	if MULTI_MEDIA_CONTENT_FORMAT_RAW != "raw" {
		t.Errorf("MULTI_MEDIA_CONTENT_FORMAT_RAW = %v, want %v", MULTI_MEDIA_CONTENT_FORMAT_RAW, "raw")
	}
	if MULTI_MEDIA_CONTENT_FORMAT_URL != "url" {
		t.Errorf("MULTI_MEDIA_CONTENT_FORMAT_URL = %v, want %v", MULTI_MEDIA_CONTENT_FORMAT_URL, "url")
	}
}

func TestMessageFinishReasonConstants(t *testing.T) {
	if MessageFinishReasonContentFiltered != "content_filter" {
		t.Errorf("MessageFinishReasonContentFiltered = %v, want %v", MessageFinishReasonContentFiltered, "content_filter")
	}
	if MessageFinishReasonFunctionCall != "function_call" {
		t.Errorf("MessageFinishReasonFunctionCall = %v, want %v", MessageFinishReasonFunctionCall, "function_call")
	}
	if MessageFinishReasonStopped != "stop" {
		t.Errorf("MessageFinishReasonStopped = %v, want %v", MessageFinishReasonStopped, "stop")
	}
	if MessageFinishReasonTokenLimitReached != "length" {
		t.Errorf("MessageFinishReasonTokenLimitReached = %v, want %v", MessageFinishReasonTokenLimitReached, "length")
	}
	if MessageFinishReasonToolCalls != "tool_calls" {
		t.Errorf("MessageFinishReasonToolCalls = %v, want %v", MessageFinishReasonToolCalls, "tool_calls")
	}
}