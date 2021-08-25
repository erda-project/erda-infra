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
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"

	"github.com/erda-project/erda-infra/providers/httpserver/server"
	"github.com/labstack/echo"
	"github.com/recallsong/go-utils/net/httpx/filesystem"
)

type (
	// Router .
	Router interface {
		GET(path string, handler interface{}, options ...interface{})
		POST(path string, handler interface{}, options ...interface{})
		DELETE(path string, handler interface{}, options ...interface{})
		PUT(path string, handler interface{}, options ...interface{})
		PATCH(path string, handler interface{}, options ...interface{})
		HEAD(path string, handler interface{}, options ...interface{})
		CONNECT(path string, handler interface{}, options ...interface{})
		OPTIONS(path string, handler interface{}, options ...interface{})
		TRACE(path string, handler interface{}, options ...interface{})

		Any(path string, handler interface{}, options ...interface{})
		Static(prefix, root string, options ...interface{})
		File(path, filepath string, options ...interface{})

		Add(method, path string, handler interface{}, options ...interface{}) error
	}
	// RouterManager .
	RouterManager interface {
		NewRouter(opts ...interface{}) RouterTx
		Reloadable() bool
	}
	// RouterTx .
	RouterTx interface {
		Router
		Commit() error
		Rollback()
		Reloadable() bool
	}
)

type (
	routeKey struct {
		method string
		path   string
	}
	route struct {
		method string
		path   string
		group  string
		hide   bool
		desc   string

		handler server.HandlerFunc
	}
	router struct {
		lock         *sync.Mutex
		done         bool
		err          error
		updateRoutes func(map[routeKey]*route)
		tx           server.RouterTx
		pathFormater *pathFormater
		routes       map[routeKey]*route
		group        string
		interceptors []server.MiddlewareFunc
	}
)

func (r *router) Add(method, path string, handler interface{}, options ...interface{}) error {
	pathFormater := r.getPathFormater(options)
	var pathParser server.MiddlewareFunc
	if pathFormater.parser != nil {
		pathParser = pathFormater.parser(path)
	}
	path = pathFormater.format(path)
	method = strings.ToUpper(method)

	key := routeKey{method: method, path: path}
	if rt, ok := r.routes[key]; ok {
		if rt.group != r.group {
			r.err = fmt.Errorf("httpserver routes [%s %s] conflict between groups (%s, %s)",
				key.method, key.path, rt.group, r.group)
		} else {
			r.err = fmt.Errorf("httpserver routes [%s %s] conflict in group %s",
				key.method, key.path, rt.group)
		}
		return r.err
	}
	route := &route{
		method: method,
		path:   path,
		group:  r.group,
	}
	for _, opt := range options {
		processRouteOptions(route, opt)
	}
	r.routes[key] = route

	if handler != nil {
		interceptors := getInterceptors(options)
		route.handler = r.add(method, path, handler, interceptors, pathParser)
	}
	return nil
}

type routeOption func(r *route)

func processRouteOptions(r *route, opt interface{}) {
	if fn, ok := opt.(routeOption); ok {
		fn(r)
	}
}

// WithDescription for Route, description for this route
func WithDescription(desc string) interface{} {
	return routeOption(func(r *route) {
		r.desc = desc
	})
}

// WithHide for Route, not print this route
func WithHide(hide bool) interface{} {
	return routeOption(func(r *route) {
		r.hide = hide
	})
}

// WithInterceptor for Router
func WithInterceptor(fn func(handler func(ctx Context) error) func(ctx Context) error) interface{} {
	return Interceptor(fn)
}

func (r *router) GET(path string, handler interface{}, options ...interface{}) {
	r.Add(http.MethodGet, path, handler, options...)
}

func (r *router) POST(path string, handler interface{}, options ...interface{}) {
	r.Add(http.MethodPost, path, handler, options...)
}

func (r *router) DELETE(path string, handler interface{}, options ...interface{}) {
	r.Add(http.MethodDelete, path, handler, options...)
}

func (r *router) PUT(path string, handler interface{}, options ...interface{}) {
	r.Add(http.MethodPut, path, handler, options...)
}

func (r *router) PATCH(path string, handler interface{}, options ...interface{}) {
	r.Add(http.MethodPatch, path, handler, options...)
}

func (r *router) HEAD(path string, handler interface{}, options ...interface{}) {
	r.Add(http.MethodHead, path, handler, options...)
}

func (r *router) CONNECT(path string, handler interface{}, options ...interface{}) {
	r.Add(http.MethodConnect, path, handler, options...)
}

func (r *router) OPTIONS(path string, handler interface{}, options ...interface{}) {
	r.Add(http.MethodOptions, path, handler, options...)
}

func (r *router) TRACE(path string, handler interface{}, options ...interface{}) {
	r.Add(http.MethodTrace, path, handler, options...)
}

var allMethods = []string{
	http.MethodConnect,
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
	http.MethodTrace,
}

func init() { sort.Strings(allMethods) }

func (r *router) Any(path string, handler interface{}, options ...interface{}) {
	for _, method := range allMethods {
		r.Add(method, path, handler, options...)
	}
}

// WithFileSystem for Static And File
func WithFileSystem(fs http.FileSystem) interface{} {
	return fs
}

type filesystemPath string

// WithFileSystemPath for Static And File
func WithFileSystemPath(root string) interface{} {
	return filesystemPath(root)
}

func (r *router) Static(prefix, root string, options ...interface{}) {
	var fs http.FileSystem
	for _, opt := range options {
		if files, ok := opt.(http.FileSystem); ok {
			fs = files
		} else if path, ok := opt.(filesystemPath); ok {
			root = filepath.Join(string(path), root)
		}
	}
	if root == "" {
		root = "."
	}
	if fs == nil {
		r.addPrefix(prefix, root, func(c echo.Context) error {
			p, err := url.PathUnescape(c.Param("*"))
			if err != nil {
				return err
			}
			name := filepath.Join(root, path.Clean("/"+p)) // "/"+ for security
			return c.File(name)
		}, options...)
	} else {
		fs := filesystem.New(fs).SetRoot(root).SetRoute(prefix)
		handler := fs.Handler
		r.addPrefix(prefix, root, func(c server.Context) error {
			handler.ServeHTTP(c.Response(), c.Request())
			return nil
		}, options...)
	}
}

func (r *router) addPrefix(prefix, root string, h func(c echo.Context) error, options ...interface{}) {
	r.GET(prefix, h, options...)
	if prefix == "/" {
		r.GET(prefix+"*", h, options...)
	} else {
		r.GET(prefix+"/*", h, options...)
	}
}

func (r *router) File(path, file string, options ...interface{}) {
	var fs http.FileSystem
	for _, opt := range options {
		if files, ok := opt.(http.FileSystem); ok {
			fs = files
		} else if path, ok := opt.(filesystemPath); ok {
			file = filepath.Join(string(path), file)
		}
	}
	if fs == nil {
		r.GET(path, func(c echo.Context) error {
			return c.File(file)
		}, options...)
	} else {
		fs := filesystem.New(fs).SetRoot(file).SetRoute(path)
		handler := fs.Handler
		r.GET(path, func(c server.Context) error {
			handler.ServeHTTP(c.Response(), c.Request())
			return nil
		}, options...)
	}
}

// Commit .
func (r *router) Commit() error {
	if r.lock != nil {
		if r.done {
			return fmt.Errorf("routes commited")
		}
		r.done = true
		if r.err != nil {
			r.lock.Unlock()
			return r.err
		}
		r.tx.Commit()
		r.updateRoutes(r.routes)
		r.lock.Unlock()
	}
	return nil
}

// Rollback .
func (r *router) Rollback() {
	if r.lock != nil {
		r.lock.Unlock()
	}
}

// Reloadable .
func (r *router) Reloadable() bool { return r.lock != nil }

type routesSorter []*route

func (s routesSorter) Len() int      { return len(s) }
func (s routesSorter) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s routesSorter) Less(i, j int) bool {
	if s[i].group == s[j].group {
		if s[i].path == s[j].path {
			return s[i].method < s[j].method
		}
		return s[i].path < s[j].path
	}
	return s[i].group < s[j].group
}

func listRoutes(routeMap map[routeKey]*route) []*route {
	routes := make([]*route, 0, len(routeMap))
	for _, route := range routeMap {
		routes = append(routes, route)
	}
	sort.Sort(routesSorter(routes))
	return routes
}

func (p *provider) printRoutes(routes map[routeKey]*route) {
	list := listRoutes(routes)
	var group, path string
	var methods []string
	printRoute := func(group, path string, methods []string) {
		sort.Strings(methods)
		if reflect.DeepEqual(methods, allMethods) {
			p.Log.Infof("%s --> [%s] %-7s %s", p.Cfg.Addr, group, "*", path)
		} else {
			for _, method := range methods {
				p.Log.Infof("%s --> [%s] %-7s %s", p.Cfg.Addr, group, method, path)
			}
		}
	}
	for _, route := range list {
		if route.hide {
			continue
		}
		if methods == nil {
			group, path = route.group, route.path
			methods = []string{route.method}
			continue
		} else if path == route.path && group == route.group {
			methods = append(methods, route.method)
			continue
		}
		printRoute(group, path, methods)
		group, path = route.group, route.path
		methods = []string{route.method}
	}
	if len(methods) > 0 {
		printRoute(group, path, methods)
	}
}

type routerManager struct {
	group string
	opts  []interface{}
	p     *provider
}

func (rm *routerManager) NewRouter(opts ...interface{}) RouterTx {
	args := make([]interface{}, len(rm.opts)+len(opts))
	copy(args, rm.opts)
	copy(args[len(rm.opts):], opts)
	return rm.p.newRouter(rm.group, args...)
}

func (rm *routerManager) Reloadable() bool { return rm.p.Cfg.Reloadable }
