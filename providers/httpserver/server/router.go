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

package server

import (
	"sync"
	"sync/atomic"

	"github.com/labstack/echo"
)

type (
	// Context .
	Context = echo.Context
	// HandlerFunc .
	HandlerFunc = echo.HandlerFunc
	// MiddlewareFunc .
	MiddlewareFunc = echo.MiddlewareFunc

	// Router .
	Router interface {
		Add(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc)
		Find(method, path string, c Context)
		NewContext() Context
		ReleaseContext(c Context)
	}
)

type router struct {
	*echo.Router
	e    *echo.Echo
	pool sync.Pool
}

func newRouter(e *echo.Echo, binder echo.Binder, validator echo.Validator) Router {
	if e == nil {
		e = echo.New()
		e.Binder, e.Validator = binder, validator
	}
	r := &router{
		Router: e.Router(),
		e:      e,
	}
	r.pool.New = func() interface{} {
		return r.e.NewContext(nil, nil)
	}
	return r
}

// Add registers a new route for an HTTP method and path with matching handler
// in the router with optional route-level middleware.
func (rt *router) Add(method, path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	h := handler
	// Chain middleware
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	rt.Router.Add(method, path, h)
}

// NewContext .
func (rt *router) NewContext() Context {
	return rt.pool.Get().(Context)
}

func (rt *router) ReleaseContext(c Context) {
	rt.pool.Put(c)
}

type fixedRouterManager struct{ Router }

func (r *fixedRouterManager) GetRouter() Router   { return r.Router }
func (r *fixedRouterManager) NewRouter() RouterTx { return r }
func (r *fixedRouterManager) Commit()             {}

func newFixedRouterManager(e *echo.Echo) routerManager {
	return &fixedRouterManager{
		Router: newRouter(nil, e.Binder, e.Validator),
	}
}

type reloadableRouterManager struct {
	v         atomic.Value
	binder    echo.Binder
	validator echo.Validator
}

func (r *reloadableRouterManager) GetRouter() Router { return r.v.Load().(Router) }
func (r *reloadableRouterManager) NewRouter() RouterTx {
	return &reloadableRouterTx{
		Router: newRouter(nil, r.binder, r.validator),
		rr:     r,
	}
}

type reloadableRouterTx struct {
	Router
	rr *reloadableRouterManager
}

func (r *reloadableRouterTx) Commit() { r.rr.v.Store(r.Router) }

func newReloadableRouterManager(e *echo.Echo) routerManager {
	r := &reloadableRouterManager{
		binder:    e.Binder,
		validator: e.Validator,
	}
	r.NewRouter().Commit()
	return r
}
