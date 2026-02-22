// Copyright (c) 2023-2025 RapidaAI
// Author: Prashant Srivastav <prashant@rapida.ai>
//
// Licensed under GPL-2.0 with Rapida Additional Terms.
// See LICENSE.md or contact sales@rapida.ai for commercial usage.
package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/rapidaai/pkg/utils"
	"github.com/rapidaai/protos"
)

type Metadata struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (d *Metadata) Type() string {
	return "metadata"
}
func (d *Metadata) SetValue(src interface{}) error {
	switch v := src.(type) {
	case string:
		d.Value = v
	case []byte:
		d.Value = string(v)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		d.Value = fmt.Sprintf("%d", v)
	case float32, float64:
		d.Value = fmt.Sprintf("%f", v)
	case bool:
		d.Value = strconv.FormatBool(v)
	case nil:
		d.Value = ""
	default:
		// Use JSON marshaling for all other types, including maps and complex structures
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		d.Value = string(jsonBytes)
	}
	return nil
}

func NewMetadata(k string, v interface{}) *Metadata {
	md := &Metadata{
		Key: k,
	}
	md.SetValue(v)
	return md
}

func NewMetadataList(data map[string]interface{}) []*Metadata {
	var metadataList []*Metadata
	for key, value := range data {
		metadataList = append(metadataList, NewMetadata(key, value))
	}
	return metadataList
}

func ToMetadata(mtr *protos.Metadata) *Metadata {
	out := &Metadata{}
	err := utils.Cast(mtr, out)
	if err != nil {
		return nil
	}
	return out
}

func ToMetadatas(mtr []*protos.Metadata) []*Metadata {
	out := make([]*Metadata, len(mtr))
	for idx, k := range mtr {
		out[idx] = ToMetadata(k)
	}
	return out
}
