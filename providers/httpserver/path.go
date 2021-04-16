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

package httpserver

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/erda-project/erda-infra/pkg/transport/http/httprule"
	"github.com/erda-project/erda-infra/pkg/transport/http/runtime"
	"github.com/labstack/echo"
)

func (r *router) getPathFormater(options []interface{}) *pathFormater {
	pformater := r.pathFormater
	for _, arg := range options {
		if f, ok := arg.(*pathFormater); ok {
			pformater = f
		}
	}
	if pformater == nil {
		pformater = newPathFormater()
	}
	return pformater
}

type pathFormater struct {
	typ    PathFormat
	format func(string) string
	parser func(path string) func(echo.HandlerFunc) echo.HandlerFunc
}

func newPathFormater() *pathFormater {
	return &pathFormater{
		typ:    PathFormatEcho,
		format: buildEchoPath,
	}
}

func buildEchoPath(p string) string { return p }

// convert googleapis path to echo path
func buildGoogleAPIsPath(path string) string {
	// skip query string
	idx := strings.Index(path, "?")
	if idx >= 0 {
		path = path[0:idx]
	}

	sb := &strings.Builder{}
	chars := []rune(path)
	start, i, l := 0, 0, len(chars)
	for ; i < l; i++ {
		c := chars[i]
		switch c {
		case '{':
			sb.WriteString(string(chars[start:i]))
			start = i
			var hasEq bool
			var name string
			begin := i
			i++ // skip '{'
		loop:
			for ; i < l; i++ {
				c = chars[i]
				switch c {
				case '}':
					begin++ // skip '{' or '='
					if len(chars[begin:i]) <= 0 || len(chars[begin:i]) == 1 && chars[begin] == '*' {
						sb.WriteString("*")
					} else if hasEq {
						pattern := string(chars[begin:i])
						if !strings.HasPrefix(pattern, "/") {
							idx := strings.Index(pattern, "/")
							if idx >= 0 {
								pattern = pattern[idx:]
							} else {
								pattern = ""
							}
						}
						sb.WriteString(":" + name + strings.ReplaceAll(pattern, ":", "%3A")) // replace ":" to %3A
					} else {
						sb.WriteString(":" + name + url.PathEscape(string(chars[begin:i])))
					}
					start = i + 1 // skip '}'
					break loop
				case '=':
					name = url.PathEscape(string(chars[begin+1 : i]))
					hasEq = true
					begin = i
				}
			}
		}
	}
	if start < l {
		sb.WriteString(strings.ReplaceAll(string(chars[start:]), ":", "%3A")) // replace ":" to %3A
	}
	return sb.String()
}

func googleAPIsPathParamsInterceptor(path string) func(echo.HandlerFunc) echo.HandlerFunc {
	raw := func(handler echo.HandlerFunc) echo.HandlerFunc {
		return handler
	}
	if path == "/" {
		return raw
	}
	compiler, err := httprule.Parse(path)
	if err != nil {
		panic(fmt.Errorf("invalid path format: %s", err))
	}
	temp := compiler.Compile()
	if len(temp.Fields) <= 0 {
		return raw
	}
	pattern, err := runtime.NewPattern(httprule.SupportPackageIsVersion1, temp.OpCodes, temp.Pool, temp.Verb)
	if err != nil {
		panic(fmt.Errorf("fail to create path pattern: %s", err))
	}
	return func(handler echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			path := ctx.Request().URL.Path
			if len(path) > 0 {
				components := strings.Split(path[1:], "/")
				last := len(components) - 1
				var verb string
				if idx := strings.LastIndex(components[last], ":"); idx >= 0 {
					c := components[last]
					components[last], verb = c[:idx], c[idx+1:]
				}
				vars, err := pattern.Match(components, verb)
				if err != nil {
					return echo.NotFoundHandler(ctx)
				}
				c := ctx.(*context)
				c.vars = vars
				ctx = c
			}
			return handler(ctx)
		}
	}
}
