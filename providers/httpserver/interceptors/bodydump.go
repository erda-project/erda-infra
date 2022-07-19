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

package interceptors

import (
	"fmt"
	"net/http/httputil"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// BodyDump dump body for http request.
func BodyDump(enable bool, maxBodySizeBytes int) echo.MiddlewareFunc {
	return middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Skipper: func(c echo.Context) bool { return !enable },
		Handler: func(c echo.Context, reqBody []byte, respBody []byte) {
			// request
			reqBase, err := httputil.DumpRequest(c.Request(), false)
			if err == nil {
				fmt.Printf("(%s) Request:\n%s", GetRequestID(c), reqBase)
			}
			if len(reqBody) <= maxBodySizeBytes {
				fmt.Printf("(%s) Request Body:\n%s\n-END-\n", GetRequestID(c), string(reqBody))
			} else {
				fmt.Printf("(%s) Request Body: (Ignored, Body too long)\n", GetRequestID(c))
			}
			// response
			if len(respBody) <= maxBodySizeBytes {
				fmt.Printf("(%s) Response Body:\n%s-END-\n", GetRequestID(c), string(respBody))
			} else {
				fmt.Printf("(%s) Response Body: (Ignored, Body too long)\n", GetRequestID(c))
			}
		},
	})
}
