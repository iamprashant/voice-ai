// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"strings"
	"time"

	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

func OnlyStringContent(Contents []*Content) string {
	var builder strings.Builder
	if len(Contents) == 0 {
		return ""
	}
	for _, c := range Contents {
		if commons.ResponseContentType(c.GetContentType()) == commons.TEXT_CONTENT {
			if commons.ResponseContentFormat(c.GetContentFormat()) == commons.TEXT_CONTENT_FORMAT_RAW {
				builder.Write(c.Content)
			}
		}
	}
	return builder.String()
}

func ContentString(c *protos.Content) string {
	var builder strings.Builder
	if commons.ResponseContentType(c.GetContentType()) == commons.TEXT_CONTENT {
		if commons.ResponseContentFormat(c.GetContentFormat()) == commons.TEXT_CONTENT_FORMAT_RAW {
			builder.Write(c.Content)
		}
	}
	return builder.String()
}

func OnlyStringProtoContent(Contents []*protos.Content) string {
	var builder strings.Builder
	if len(Contents) == 0 {
		return ""
	}
	for _, c := range Contents {
		if commons.ResponseContentType(c.GetContentType()) == commons.TEXT_CONTENT {
			if commons.ResponseContentFormat(c.GetContentFormat()) == commons.TEXT_CONTENT_FORMAT_RAW {
				builder.Write(c.Content)
			}
		}
	}
	return builder.String()
}

func ContainsAudioContent(Contents []*protos.Content) bool {
	for _, c := range Contents {
		if commons.ResponseContentType(c.GetContentType()) == commons.AUDIO_CONTENT {
			return true
		}
	}
	return false
}

func ToMessage(msg *protos.Message) *Message {
	out := &Message{}
	err := utils.Cast(msg, out)
	if err != nil {
		return nil
	}
	out.Time = time.Now()
	return out
}

func ToMessages(msgs []*protos.Message) []*Message {
	out := make([]*Message, 0, len(msgs))
	for _, msg := range msgs {
		if convertedMsg := ToMessage(msg); convertedMsg != nil {
			out = append(out, convertedMsg)
		}
	}
	return out
}
