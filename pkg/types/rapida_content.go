// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

type Content struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// audio, image, video, text etc
	ContentType string `protobuf:"bytes,2,opt,name=contentType,proto3" json:"contentType,omitempty"`
	// from raw string to url all can be
	ContentFormat string `protobuf:"bytes,3,opt,name=contentFormat,proto3" json:"contentFormat,omitempty"`
	// actual content
	Content []byte `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
	// added meta data incase you want to add something which is not supported in
	Meta map[string]interface{} `protobuf:"bytes,5,opt,name=meta,proto3" json:"meta,omitempty"`
}

func (c *Content) GetContent() []byte {
	return c.Content
}

func (c *Content) GetString() string {
	return string(c.Content)
}

func (c *Content) GetContentFormat() string {
	return string(c.ContentFormat)
}

func (c *Content) GetContentType() string {
	return string(c.ContentType)
}

// func (c *Content) ToProto() *protos.Content {
// 	protoC := &protos.Content{}
// 	utils.Cast(c, protoC)
// 	return protoC
// }

// type Contents []*Content

// func (m Contents) ToProto() []*protos.Content {
// 	out := make([]*protos.Content, len(m))
// 	for idx, k := range m {
// 		out[idx] = k.ToProto()
// 	}
// 	return out
// }
