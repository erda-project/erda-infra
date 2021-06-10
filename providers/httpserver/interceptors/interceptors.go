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
	"runtime"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// CORS .
func CORS() interface{} {
	return middleware.CORS()
}

// Recover .
func Recover(log logs.Logger) interface{} {
	const StackSize = 4 << 10
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, StackSize)
					length := runtime.Stack(stack, true)
					log.Errorf("[PANIC RECOVER] %v %s\n", err, stack[:length])
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}
