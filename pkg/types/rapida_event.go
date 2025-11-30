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

type Event struct {
	EventType string                 `protobuf:"bytes,1,opt,name=eventType,proto3" json:"eventType,omitempty"`
	Payload   map[string]interface{} `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (d *Event) SetValue(src interface{}) error {
	switch v := src.(type) {
	case string:
		d.Payload = map[string]interface{}{"value": v}
	case []byte:
		d.Payload = map[string]interface{}{"value": string(v)}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		d.Payload = map[string]interface{}{"value": fmt.Sprintf("%d", v)}
	case float32, float64:
		d.Payload = map[string]interface{}{"value": fmt.Sprintf("%f", v)}
	case bool:
		d.Payload = map[string]interface{}{"value": strconv.FormatBool(v)}
	case nil:
		d.Payload = map[string]interface{}{}
	default:
		// Use JSON marshalling for all other types
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal Payload: %w", err)
		}
		var result map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			return fmt.Errorf("failed to unmarshal into InterfaceMap: %w", err)
		}
		d.Payload = result
	}
	return nil
}

func NewEvent(k string, v interface{}) *Event {
	md := &Event{
		EventType: k,
	}
	md.SetValue(v)
	return md
}
