/*
 *  Copyright (c) 2024. Rapida
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in
 *  all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 *  THE SOFTWARE.
 *
 *  Author: Prashant <prashant@rapida.ai>
 *
 */
package types

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Metadata struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
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
		// Use JSON marshalling for all other types, including maps and complex structures
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
