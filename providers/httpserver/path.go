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
	libcontext "context"
	"net/http"
	"net/url"
	"strings"

	"github.com/erda-project/erda-infra/pkg/transport/http/runtime"
	"github.com/erda-project/erda-infra/providers/httpserver/server"
)

// PathFormat .
type PathFormat int32

// PathFormat values
const (
	PathFormatEcho       = 0
	PathFormatGoogleAPIs = 1
)

type contextKey int

const (
	varsKey contextKey = iota
)

// WithPathFormat .
func WithPathFormat(format PathFormat) interface{} {
	formater := &pathFormater{typ: format}
	switch format {
	case PathFormatGoogleAPIs:
		formater.format = buildGoogleAPIsPath
		formater.parser = googleAPIsPathParamsInterceptor
	default:
		formater.format = buildEchoPath
	}
	return formater
}

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
	parser func(path string) func(server.HandlerFunc) server.HandlerFunc
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

func googleAPIsPathParamsInterceptor(path string) func(server.HandlerFunc) server.HandlerFunc {
	raw := func(handler server.HandlerFunc) server.HandlerFunc {
		return handler
	}
	matcher, err := runtime.Compile(path)
	if err != nil {
		// panic(fmt.Errorf("path %q error: %s", path, err))
		return func(server.HandlerFunc) server.HandlerFunc {
			return func(ctx server.Context) error {
				return server.NotFoundHandler(ctx)
			}
		}
	}
	if matcher.IsStatic() {
		return raw
	}
	return func(handler server.HandlerFunc) server.HandlerFunc {
		return func(ctx server.Context) error {
			path := ctx.Request().URL.Path
			vars, err := matcher.Match(path)
			if err != nil {
				return server.NotFoundHandler(ctx)
			}
			c := ctx.(*context)
			c.vars = vars
			ctx = c
			ctx.SetRequest(ctx.Request().WithContext(makeCtxWithVars(ctx.Request().Context(), vars)))
			return handler(ctx)
		}
	}
}

func makeCtxWithVars(ctx libcontext.Context, vars map[string]string) libcontext.Context {
	return libcontext.WithValue(ctx, varsKey, vars)
}

// Vars returns the route variables for the current request, if any.
func Vars(r *http.Request) map[string]string {
	if rv := r.Context().Value(varsKey); rv != nil {
		return rv.(map[string]string)
	}
	return nil
}

// Var return the specified variable value and exist from the current request, if any.
func Var(r *http.Request, key string) (string, bool) {
	vars := Vars(r)
	if vars == nil {
		return "", false
	}
	val, ok := vars[key]
	return val, ok
}
