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
	"bytes"
	"io"
	"math"
	"testing"
)

func TestRequestHeader(t *testing.T) {
	tests := []struct {
		name    string
		h       *RequestHeader
		tail    string
		want    *RequestHeader
		wantErr bool
	}{
		{
			h: &RequestHeader{
				Version:    ProtocolVersion,
				Name:       "test-name",
				Token:      "test-token",
				ShadowAddr: "test-addr",
			},
			tail: "test",
			want: &RequestHeader{
				Version:    ProtocolVersion,
				Name:       "test-name",
				Token:      "test-token",
				ShadowAddr: "test-addr",
			},
		},
		{
			h: &RequestHeader{
				Version:    math.MaxUint32,
				Name:       "test2-name",
				Token:      "test2-token",
				ShadowAddr: "test2-addr",
			},
			tail: "test",
			want: &RequestHeader{
				Version:    math.MaxUint32,
				Name:       "test2-name",
				Token:      "test2-token",
				ShadowAddr: "test2-addr",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if err := EncodeRequestHeader(w, tt.h); (err != nil) != tt.wantErr {
				t.Errorf("EncodeRequestHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if (err != nil) && tt.wantErr {
				return
			}
			w.WriteString(tt.tail)

			r := bytes.NewReader(w.Bytes())
			h, err := DecodeRequestHeader(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeRequestHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if (err != nil) && tt.wantErr {
				return
			}
			if *h != *tt.h {
				t.Errorf("EncodeRequestHeader() != DecodeRequestHeader()")
				return
			}
			byts, _ := io.ReadAll(r)
			if tt.tail != string(byts) {
				t.Errorf("EncodeRequestHeader() tail bytes != DecodeRequestHeader() tail bytes")
				return
			}

			if *h != *tt.want {
				t.Errorf("DecodeRequestHeader() got %v, want %v", *h, *tt.want)
				return
			}
		})
	}
}
