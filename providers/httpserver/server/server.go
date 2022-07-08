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
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

type (
	// Server .
	Server interface {
		Use(middleware ...MiddlewareFunc)
		NewRouter() RouterTx
		Router() Router
		Start(addr string) error
		Close() error
	}
	// RouterTx .
	RouterTx interface {
		Router
		Commit()
	}
)

var (
	// NotFoundHandler .
	NotFoundHandler = echo.NotFoundHandler
)

type (
	routerManager interface {
		GetRouter() Router
		NewRouter() RouterTx
	}
	server struct {
		e          *echo.Echo
		reloadable bool
		router     routerManager
		middleware []MiddlewareFunc
	}
)

// New .
func New(reloadable bool, binder echo.Binder, validator echo.Validator) Server {
	s := &server{
		e:          echo.New(),
		reloadable: reloadable,
	}
	s.e.HideBanner, s.e.HidePort = true, true
	s.e.Binder, s.e.Validator = binder, validator
	s.e.Server.Handler, s.e.TLSServer.Handler = s, s
	if reloadable {
		s.router = newReloadableRouterManager(s.e)
	} else {
		s.router = newFixedRouterManager(s.e)
	}
	return s
}

// Start .
func (s *server) Start(addr string) error {
	err := s.startHTTP(addr)
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Close .
func (s *server) Close() error {
	err := s.e.Server.Close()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// ServeHTTP .
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router := s.router.GetRouter()

	// Acquire context
	c := router.NewContext()
	// Release context
	defer router.ReleaseContext(c)
	c.Reset(r, w)

	// Find handler
	router.Find(r.Method, getPath(r), c)
	h := c.Handler()
	if h == nil {
		h = NotFoundHandler
	}
	for i := len(s.middleware) - 1; i >= 0; i-- {
		h = s.middleware[i](h)
	}

	// Execute chain
	if err := h(c); err != nil {
		c.Logger().Errorf("url method: %s, path: %s, matcherPath: %s, ip: %s, header: %v", c.Request().Method, c.Request().URL.Path, c.Path(), c.RealIP(), c.Request().Header)
		s.e.HTTPErrorHandler(err, c)
	}
}

func getPath(r *http.Request) string {
	path := r.URL.RawPath
	if path == "" {
		path = r.URL.Path
	}
	return path
}

// NewRouter .
func (s *server) NewRouter() RouterTx { return s.router.NewRouter() }

// Router .
func (s *server) Router() Router { return s.router.GetRouter() }

// Use .
func (s *server) Use(middleware ...MiddlewareFunc) {
	s.middleware = append(s.middleware, middleware...)
}

// Start starts an HTTP server.
func (s *server) startHTTP(address string) error {
	s.e.Server.Addr = address
	return s.startServer(s.e.Server)
}

// StartServer starts a custom http server.
func (s *server) startServer(svr *http.Server) (err error) {
	if svr.TLSConfig == nil {
		if s.e.Listener == nil {
			s.e.Listener, err = newListener(svr.Addr)
			if err != nil {
				return err
			}
		}
		return svr.Serve(s.e.Listener)
	}
	if s.e.TLSListener == nil {
		l, err := newListener(svr.Addr)
		if err != nil {
			return err
		}
		s.e.TLSListener = tls.NewListener(l, svr.TLSConfig)
	}
	return svr.Serve(s.e.TLSListener)
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func newListener(address string) (*tcpKeepAliveListener, error) {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	return &tcpKeepAliveListener{l.(*net.TCPListener)}, nil
}
