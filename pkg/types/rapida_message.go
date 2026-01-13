// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rapidaai/pkg/commons"
	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

type Message struct {
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`

	// role
	Role string `protobuf:"bytes,1,opt,name=role,proto3" json:"role,omitempty"`

	// contents
	Contents []*Content `protobuf:"bytes,2,rep,name=contents,proto3" json:"contents,omitempty"`

	// all the tool calls
	ToolCalls []*ToolCall `protobuf:"bytes,2,rep,name=toolCalls,proto3" json:"toolCalls,omitempty"`

	// time
	Time time.Time `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`

	// meta
	Meta map[string]interface{} `protobuf:"bytes,1,opt,name=meta,proto3" json:"meta,omitempty"`

	// metrics
	Metrics []*Metric `protobuf:"bytes,1,opt,name=metrics,proto3" json:"metrics,omitempty"`
}

// Assuming you have a types.ToolCall struct like this:
type ToolCall struct {
	Id       *string       `json:"id"`
	Type     *string       `json:"type"`
	Function *FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      *string `json:"name"`
	Arguments *string `json:"arguments"`
}

func (fc *FunctionCall) MergeArguments(newArgs *string) {
	if newArgs == nil {
		return
	}
	if fc.Arguments == nil {
		// Initialize Arguments if it's nil
		fc.Arguments = newArgs
	} else {
		// Append newArgs to Arguments
		*fc.Arguments += *newArgs
	}
}

func (fc *FunctionCall) MergeName(newName *string) {
	if newName == nil {
		return
	}
	if fc.Name == nil {
		fc.Name = newName
	} else {
		// Append newName to Name
		*fc.Name += *newName
	}
}

func (msg *Message) GetLanguage() Language {
	if msg.Meta != nil {
		if langCode, ok := msg.Meta["language_iso_1"].(string); ok {
			return GetLanguageByName(langCode)
		}
	}
	return GetLanguageByName("en")
}

func (m *Message) AddMetadata(key string, value interface{}) {
	if m.Meta == nil {
		m.Meta = make(map[string]interface{})
	}
	m.Meta[key] = value
}

func (m *Message) WithMetadata(x map[string]interface{}) *Message {
	if m.Meta == nil {
		m.Meta = make(map[string]interface{})
	}
	for key, value := range x {
		m.Meta[key] = value
	}
	return m
}

func (msg *Message) GetId() string {
	return msg.Id
}

func (msg *Message) GetContents() []*Content {
	return msg.Contents
}

func (msg *Message) GetToolCalls() []*ToolCall {
	return msg.ToolCalls
}

func (msg *Message) GetRole() string {
	return msg.Role
}

func (msg *Message) GetTime() string {
	return msg.Time.UTC().Format(time.RFC3339)
}

func (msg *Message) ToProto() *protos.Message {
	protoMsg := &protos.Message{}
	err := utils.Cast(msg, protoMsg)
	if err != nil {
		fmt.Printf("error while casting %v", err)
	}
	return protoMsg
}

func (msg *Message) String() string {
	var builder strings.Builder
	if len(msg.GetContents()) == 0 {
		return ""
	}
	for _, c := range msg.GetContents() {
		if commons.ResponseContentType(c.GetContentType()) == commons.TEXT_CONTENT {
			if commons.ResponseContentFormat(c.GetContentFormat()) == commons.TEXT_CONTENT_FORMAT_RAW {
				if commons.ResponseContentFormat(c.GetContentFormat()) == commons.TEXT_CONTENT_FORMAT_RAW {
					if builder.Len() > 0 {
						builder.WriteString(" ") // Add space if there's already content
					}
					builder.Write(c.Content)
				}
			}
		}
	}
	return builder.String()
}

func NewMessage(role string, content ...*Content) *Message {
	return &Message{
		Id:       uuid.NewString(),
		Role:     role,
		Contents: content,
		Time:     time.Now(),
	}
}

func (msg *Message) MergeContent(content ...*Content) *Message {
	msg.Contents = append(msg.Contents, content...)
	return msg
}
