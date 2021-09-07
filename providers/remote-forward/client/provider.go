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

package client

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"os"
	"reflect"
	"time"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	forward "github.com/erda-project/erda-infra/providers/remote-forward"
	yamux "github.com/hashicorp/yamux"
	uuid "github.com/satori/go.uuid"
)

// Interface .
type Interface interface {
	RemoteShadowAddr() string
	Values() map[string]interface{}
}

var _ (Interface) = (*provider)(nil)

type config struct {
	RemoteAddr       string `file:"remote_addr"`
	RemoteShadowAddr string `file:"remote_shadow_addr"`
	TargetAddr       string `file:"target_addr"`
	Name             string `file:"name"`
	Token            string `file:"token"`
}

type provider struct {
	Cfg      *config
	Log      logs.Logger
	conn     net.Conn
	response *forward.ResponseHeader
}

func (p *provider) RemoteShadowAddr() string       { return p.response.ShadowAddr }
func (p *provider) Values() map[string]interface{} { return p.response.Values }

func (p *provider) Init(ctx servicehub.Context) error {
	if len(p.Cfg.Name) <= 0 {
		hostname, err := os.Hostname()
		if err == nil {
			p.Cfg.Name = fmt.Sprintf("%s@%s->%s", hex.EncodeToString(uuid.NewV4().Bytes()[8:16]), hostname, p.Cfg.TargetAddr)
		} else {
			p.Cfg.Name = hex.EncodeToString(uuid.NewV4().Bytes())
		}
		p.Log.Infof("forward name is %q", p.Cfg.Name)
	}
	conn, err := net.Dial("tcp", p.Cfg.RemoteAddr)
	if err != nil {
		return fmt.Errorf("failed to connect remote forward server: %s", err)
	}
	p.conn = conn
	p.response, err = p.handshake(conn)
	if err != nil {
		return err
	}
	return nil
}

func (p *provider) Start() error {
	conn := p.conn
	session, err := yamux.Server(conn, nil)
	if err != nil {
		if errors.Is(err, net.ErrClosed) {
			return nil
		}
		return err
	}
	defer session.Close()
	defer conn.Close()
	for {
		conn, err := session.Accept()
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
	if p.conn != nil {
		conn := p.conn
		p.conn = nil
		err := conn.Close()
		if !errors.Is(err, net.ErrClosed) {
			return err
		}
	}
	return nil
}

func (p *provider) handshake(conn net.Conn) (resp *forward.ResponseHeader, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("handshake error: %w", err)
		}
	}()
	err = conn.SetDeadline(time.Now().Add(60 * time.Second))
	if err != nil {
		return nil, err
	}
	err = forward.EncodeRequestHeader(conn, &forward.RequestHeader{
		Version:    forward.ProtocolVersion,
		Name:       p.Cfg.Name,
		Token:      p.Cfg.Token,
		ShadowAddr: p.Cfg.RemoteShadowAddr,
	})
	if err != nil {
		return nil, err
	}
	resp, err = forward.DecodeResponseHeader(conn)
	if err != nil {
		return nil, err
	}
	if len(resp.Error) > 0 {
		return nil, errors.New(resp.Error)
	}
	err = conn.SetDeadline(time.Time{})
	if err != nil {
		return nil, err
	}
	p.Log.Infof("remote shadow addr: %s", resp.ShadowAddr)
	return resp, nil
}

func (p *provider) handleConn(conn net.Conn) {
	defer conn.Close()
	localConn, err := net.Dial("tcp", p.Cfg.TargetAddr)
	if err != nil {
		p.Log.Error(err)
		return
	}
	defer localConn.Close()
	p.Log.Debugf("forward %s -> %s", conn.RemoteAddr(), p.Cfg.TargetAddr)
	forward.Pipe(p.Log, conn, localConn)
}

func init() {
	servicehub.Register("remote-forward-client", &servicehub.Spec{
		Services:   []string{"remote-forward-client"},
		Types:      []reflect.Type{reflect.TypeOf((*Interface)(nil)).Elem()},
		ConfigFunc: func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
