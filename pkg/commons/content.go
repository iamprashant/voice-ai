// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package commons

const (
	TEXT_CONTENT        ResponseContentType = "text"
	AUDIO_CONTENT       ResponseContentType = "audio"
	IMAGE_CONTENT       ResponseContentType = "image"
	MULTI_MEDIA_CONTENT ResponseContentType = "multi"
)

const (
	// raw string in byte format
	TEXT_CONTENT_FORMAT_RAW ResponseContentFormat = "raw"

	//
	AUDIO_CONTENT_FORMAT_RAW ResponseContentFormat = "raw"
	AUDIO_CONTENT_FORMAT_URL ResponseContentFormat = "url"

	//
	IMAGE_CONTENT_FORMAT_RAW ResponseContentFormat = "raw"
	IMAGE_CONTENT_FORMAT_URL ResponseContentFormat = "url"

	MULTI_MEDIA_CONTENT_FORMAT_RAW ResponseContentFormat = "raw"
	MULTI_MEDIA_CONTENT_FORMAT_URL ResponseContentFormat = "url"
)

type MessageFinishReason string

const (
	MessageFinishReasonContentFiltered   MessageFinishReason = "content_filter"
	MessageFinishReasonFunctionCall      MessageFinishReason = "function_call"
	MessageFinishReasonStopped           MessageFinishReason = "stop"
	MessageFinishReasonTokenLimitReached MessageFinishReason = "length"
	MessageFinishReasonToolCalls         MessageFinishReason = "tool_calls"
)
