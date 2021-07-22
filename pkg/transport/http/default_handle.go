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

package http

import (
	"encoding/json"
	"net/http"

	"github.com/erda-project/erda-infra/pkg/transport/http/encoding"
)

// DefaultHandleOptions .
func DefaultHandleOptions() *HandleOptions {
	return &HandleOptions{
		Decode: encoding.DecodeRequest,
		Encode: encoding.EncodeResponse,
		Error:  EncodeError,
	}
}

// Error .
type Error interface {
	HTTPStatus() int
}

// EncodeError default EncodeErrorFunc implement
func EncodeError(w http.ResponseWriter, r *http.Request, err error) {
	var status int
	if e, ok := err.(Error); ok {
		status = e.HTTPStatus()
	} else {
		status = http.StatusInternalServerError
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	byts, _ := json.Marshal(map[string]interface{}{
		"code": status,
		"err":  err.Error(),
	})
	w.Write(byts)
}
