// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package forward

import (
	"encoding/gob"
	"io"
	"time"
)

const (
	// ProtocolVersion .
	ProtocolVersion = 1
	// HandshakeTimeout .
	HandshakeTimeout = 60 * time.Second
)

func init() {
	gob.Register(&RequestHeader{})
	gob.Register(&ResponseHeader{})
}

// RequestHeader .
type RequestHeader struct {
	Version    uint32
	Name       string
	Token      string
	ShadowAddr string
}

// DecodeRequestHeader .
func DecodeRequestHeader(r io.Reader) (*RequestHeader, error) {
	dec := gob.NewDecoder(r)
	h := &RequestHeader{}
	err := dec.Decode(h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// EncodeRequestHeader .
func EncodeRequestHeader(w io.Writer, h *RequestHeader) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(h)
}

// ResponseHeader .
type ResponseHeader struct {
	ShadowAddr string
	Error      string
	Values     map[string]interface{}
}

// DecodeResponseHeader .
func DecodeResponseHeader(r io.Reader) (*ResponseHeader, error) {
	dec := gob.NewDecoder(r)
	h := &ResponseHeader{}
	err := dec.Decode(h)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// EncodeResponseHeader .
func EncodeResponseHeader(w io.Writer, h *ResponseHeader) error {
	enc := gob.NewEncoder(w)
	return enc.Encode(h)
}
