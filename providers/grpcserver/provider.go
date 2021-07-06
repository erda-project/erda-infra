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

package grpcserver

import (
	"fmt"
	"net"
	"reflect"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Interface .
type Interface interface {
	RegisterService(sd *grpc.ServiceDesc, ss interface{})
}

// config .
type config struct {
	Addr string `file:"addr" default:":7070" desc:"grpc address to listen"`
	TLS  struct {
		CertFile string `file:"cert_file" desc:"the TLS cert file"`
		KeyFile  string `file:"key_file" desc:"the TLS key file"`
	} `file:"tls"`
}

type provider struct {
	Cfg    *config
	Log    logs.Logger
	listen net.Listener
	server *grpc.Server
}

func (p *provider) Init(ctx servicehub.Context) error {
	lis, err := net.Listen("tcp", p.Cfg.Addr)
	if err != nil {
		return err
	}
	p.listen = lis

	var opts []grpc.ServerOption
	if len(p.Cfg.TLS.CertFile) > 0 || len(p.Cfg.TLS.KeyFile) > 0 {
		creds, err := credentials.NewServerTLSFromFile(p.Cfg.TLS.CertFile, p.Cfg.TLS.KeyFile)
		if err != nil {
			return fmt.Errorf("fail to generate credentials %v", err)
		}
		opts = append(opts, grpc.Creds(creds))
	}
	p.server = grpc.NewServer(opts...)
	return nil
}

func (p *provider) Start() error {
	p.Log.Infof("starting grpc server at %s", p.Cfg.Addr)
	return p.server.Serve(p.listen)
}

func (p *provider) Close() error {
	p.server.Stop()
	return nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	return p.server
}

func init() {
	servicehub.Register("grpc-server", &servicehub.Spec{
		Services: []string{"grpc-server"},
		Types: []reflect.Type{
			reflect.TypeOf((*grpc.Server)(nil)),
			reflect.TypeOf((*Interface)(nil)).Elem(),
		},
		Description: "grpc server",
		ConfigFunc:  func() interface{} { return &config{} },
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
