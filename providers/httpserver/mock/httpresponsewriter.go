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

package mock

import (
	"bytes"
	"net/http"
)

// HTTPResponseWriter .
type HTTPResponseWriter struct {
	Status     int
	MockHeader http.Header
	Bytes      *bytes.Buffer
}

// NewHTTPResponseWriter .
func NewHTTPResponseWriter() *HTTPResponseWriter {
	return &HTTPResponseWriter{
		MockHeader: make(http.Header),
		Bytes:      new(bytes.Buffer),
	}
}

// Header .
func (rw *HTTPResponseWriter) Header() http.Header {
	return rw.MockHeader
}

// Write .
func (rw *HTTPResponseWriter) Write(byts []byte) (int, error) {
	if rw.Bytes == nil {
		rw.Bytes = new(bytes.Buffer)
	}
	rw.Bytes.Write(byts)
	return len(byts), nil
}

// WriteHeader .
func (rw *HTTPResponseWriter) WriteHeader(statusCode int) {
	rw.Status = statusCode
}
