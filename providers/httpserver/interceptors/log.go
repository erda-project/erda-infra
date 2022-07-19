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
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/erda-project/erda-infra/base/logs"
)

// SimpleRecord record begin and end for http request.
func SimpleRecord(log logs.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			log.Infof("(%s) begin handle request: %s\n", GetRequestID(c), c.Request().URL)
			defer log.Infof("(%s) end handle request: %s\n", GetRequestID(c), c.Request().URL)
			return next(c)
		}
	}
}

// DetailLog print detail log for http request.
// Like: {"time":"2022-07-19T16:03:44.525493+08:00","id":"eSRLySVRiUAXs0VRu0AC0ETteIRKtAHg","remote_ip":"127.0.0.1","host":"localhost:9529","method":"GET","uri":"/test","user_agent":"curl/7.79.1","status":401,"error":"","latency":273532041,"latency_human":"273.532041ms","bytes_in":0,"bytes_out":25}
func DetailLog(enable bool) echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool { return !enable },
	})
}
