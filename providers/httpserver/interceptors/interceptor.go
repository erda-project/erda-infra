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

	"github.com/erda-project/erda-infra/base/logs"
)

const (
	defaultEnableFlag = "__debug__"
)

type Option struct {
	EnableFetchFuncs []EnableFetchFunc
	Log              logs.Logger
}

func NewOption(funcs []EnableFetchFunc, log logs.Logger) Option {
	return Option{
		EnableFetchFuncs: funcs,
		Log:              log,
	}
}

type EnableFetchFunc func(c echo.Context) bool

var defaultEnableFetchFunc EnableFetchFunc = func(c echo.Context) bool {
	if c.Request().URL.Query().Has(defaultEnableFlag) {
		return true
	}
	if _, ok := c.Request().Header[defaultEnableFlag]; ok {
		return true
	}
	return false
}

// judgeAnyEnable judge enable if any func executed return true
func judgeAnyEnable(c echo.Context, enableFetchFuncs []EnableFetchFunc) bool {
	for _, f := range enableFetchFuncs {
		if f(c) {
			return true
		}
	}
	return defaultEnableFetchFunc(c)
}
