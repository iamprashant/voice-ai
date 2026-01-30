// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"reflect"
	"testing"
)

func TestContent_GetContent(t *testing.T) {
	content := []byte("test content")
	c := &Content{Content: content}
	if got := c.GetContent(); !reflect.DeepEqual(got, content) {
		t.Errorf("GetContent() = %v, want %v", got, content)
	}
}

func TestContent_GetString(t *testing.T) {
	content := []byte("test string")
	c := &Content{Content: content}
	if got := c.GetString(); got != "test string" {
		t.Errorf("GetString() = %v, want %v", got, "test string")
	}
}

func TestContent_GetContentFormat(t *testing.T) {
	format := "raw"
	c := &Content{ContentFormat: format}
	if got := c.GetContentFormat(); got != format {
		t.Errorf("GetContentFormat() = %v, want %v", got, format)
	}
}

func TestContent_GetContentType(t *testing.T) {
	contentType := "text"
	c := &Content{ContentType: contentType}
	if got := c.GetContentType(); got != contentType {
		t.Errorf("GetContentType() = %v, want %v", got, contentType)
	}
}

// func TestContent_ToProto(t *testing.T) {
// 	c := &Content{
// 		Name:          "name",
// 		ContentType:   "type",
// 		ContentFormat: "format",
// 		Content:       []byte("content"),
// 		Meta:          map[string]interface{}{"key": "value"},
// 	}
// 	proto := c.ToProto()
// 	if proto.Name != "name" {
// 		t.Errorf("ToProto() Name = %v, want %v", proto.Name, "name")
// 	}
// 	// Assuming Cast works, just check it's not nil
// 	if proto == nil {
// 		t.Errorf("ToProto() returned nil")
// 	}
// }

// func TestContents_ToProto(t *testing.T) {
// 	contents := Contents{
// 		&Content{Name: "1"},
// 		&Content{Name: "2"},
// 	}
// 	protos := contents.ToProto()
// 	if len(protos) != 2 {
// 		t.Errorf("ToProto() length = %v, want %v", len(protos), 2)
// 	}
// 	if protos[0].Name != "1" {
// 		t.Errorf("ToProto()[0].Name = %v, want %v", protos[0].Name, "1")
// 	}
// }
