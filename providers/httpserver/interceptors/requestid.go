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
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// InjectRequestID inject request id for http request.
func InjectRequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			rid := req.Header.Get(echo.HeaderXRequestID)
			if rid == "" {
				rid = uuid.NewString()
				// set rid to request headers for proxy use, otherwise the forwarded target httpserver will inject a new X-Request-ID
				req.Header.Set(echo.HeaderXRequestID, rid)
			}
			res.Header().Set(echo.HeaderXRequestID, rid)

			return next(c)
		}
	}
}

// GetRequestID get request id from context.
func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}
