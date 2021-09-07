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

package encoding

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockReadCloser struct {
}

func (r mockReadCloser) Read(p []byte) (n int, err error) {
	p[0] = 't'
	p[1] = 'e'
	p[2] = 's'
	p[3] = 't'
	return 4, io.EOF
}

func (r mockReadCloser) Close() error {
	return nil
}

func makeMockRequest() *http.Request {
	r := &http.Request{}
	r.Body = mockReadCloser{}
	return r
}

func TestDecodeRequest(t *testing.T) {
	var output []byte
	type args struct {
		r   *http.Request
		out interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"case1",
			args{
				makeMockRequest(),
				&output,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DecodeRequest(tt.args.r, tt.args.out); (err != nil) != tt.wantErr {
				t.Errorf("DecodeRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	assert.Equal(t, []byte("test"), output)
}
