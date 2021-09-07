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
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	forward "github.com/erda-project/erda-infra/providers/remote-forward"
	yamux "github.com/hashicorp/yamux"
)

type (
	// Handshaker .
	Handshaker func(req *forward.RequestHeader, resp *forward.ResponseHeader) error
	// Interface .
	Interface interface {
		AddHandshaker(h Handshaker)
	}
)

var _ (Interface) = (*provider)(nil)

type config struct {
	Addr  string `file:"addr"`
	Token string `file:"token"`
}

type provider struct {
	Cfg        *config
	Log        logs.Logger
	ln         net.Listener
	handshaker []Handshaker
}

func (p *provider) Init(ctx servicehub.Context) error {
	ln, err := net.Listen("tcp", p.Cfg.Addr)
	if err != nil {
		return err
	}
	p.ln = ln
	p.Log.Infof("forward server listen at %s", p.Cfg.Addr)
	return nil
}

func (p *provider) AddHandshaker(h Handshaker) {
	p.handshaker = append(p.handshaker, h)
}

func (p *provider) Start() error {
	ln := p.ln
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}
		go p.handleConn(conn)
	}
}

func (p *provider) Close() error {
	if p.ln != nil {
		ln := p.ln
		p.ln = nil
		err := ln.Close()
		if !errors.Is(err, net.ErrClosed) {
			return err
		}
	}
	return nil
}

func (p *provider) handleConn(conn net.Conn) {
	defer conn.Close()
	req, err := p.handshake(conn)
	if err != nil {
		p.responseError(conn, err)
		return
	}
	if req == nil {
		return
	}
	resp := &forward.ResponseHeader{Values: make(map[string]interface{})}
	for _, h := range p.handshaker {
		err := h(req, resp)
		if err != nil {
			p.responseError(conn, err)
			return
		}
	}

	ln, err := net.Listen("tcp", req.ShadowAddr)
	if err != nil {
		p.responseError(conn, err)
		return
	}
	defer ln.Close()

	session, err := yamux.Client(conn, nil)
	if err != nil {
		p.responseError(conn, err)
		return
	}
	defer session.Close()

	resp.ShadowAddr = ln.Addr().String()
	p.Log.Infof("%q shadow address listen at %s", req.Name, resp.ShadowAddr)
	if err := p.responseOK(conn, resp); err != nil {
		return
	}

	go func() {
		<-session.CloseChan()
		ln.Close()
	}()
	for {
		source, err := ln.Accept()
		if err != nil {
			if !errors.Is(err, net.ErrClosed) {
				p.Log.Errorf("accept error: %s", err)
			}
			return
		}
		go func() {
			defer source.Close()
			target, err := session.Open()
			if err != nil {
				if !errors.Is(err, net.ErrClosed) {
					p.Log.Errorf("failed to open connect in session: %s", err)
				}
				return
			}
			defer target.Close()
			forward.Pipe(p.Log, target, source)
		}()
	}
}

func (p *provider) handshake(conn net.Conn) (header *forward.RequestHeader, err error) {
	defer func() {
		if err != nil && errors.Is(err, net.ErrClosed) {
			header, err = nil, nil
		} else if err != nil {
			err = fmt.Errorf("handshake error: %w", err)
		}
	}()
	err = conn.SetDeadline(time.Now().Add(forward.HandshakeTimeout))
	if err != nil {
		return nil, err
	}
	header, err = forward.DecodeRequestHeader(conn)
	if err != nil {
		return nil, err
	}
	if header.Version != forward.ProtocolVersion {
		return nil, fmt.Errorf("not support version %q", header.Version)
	}
	if header.Token != p.Cfg.Token {
		return nil, fmt.Errorf("invalid token")
	}
	err = conn.SetDeadline(time.Time{})
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (p *provider) responseError(conn net.Conn, err error) error {
	if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
		return err
	}
	err = forward.EncodeResponseHeader(conn, &forward.ResponseHeader{Error: err.Error()})
	if err != nil && !errors.Is(err, net.ErrClosed) && !errors.Is(err, io.EOF) {
		p.Log.Errorf("failed to encode response: %s", err)
	}
	return err
}

func (p *provider) responseOK(conn net.Conn, resp *forward.ResponseHeader) error {
	resp.Error = ""
	err := forward.EncodeResponseHeader(conn, resp)
	if err != nil && !errors.Is(err, net.ErrClosed) && !errors.Is(err, io.EOF) {
		p.Log.Errorf("failed to encode response: %s", err)
	}
	return err
}

func init() {
	servicehub.Register("remote-forward-server", &servicehub.Spec{
		Services:   []string{"remote-forward-server"},
		Types:      []reflect.Type{reflect.TypeOf((*Interface)(nil)).Elem()},
		ConfigFunc: func() interface{} { return &config{} },
		Creator:    func() servicehub.Provider { return &provider{} },
	})
}
