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

package grpcclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"reflect"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	grpccontext "github.com/erda-project/erda-infra/pkg/trace/inject/context/grpc"
	transgrpc "github.com/erda-project/erda-infra/pkg/transport/grpc"
)

// Interface .
type Interface interface {
	Get() *grpc.ClientConn
	NewConnect(opts ...grpc.DialOption) (*grpc.ClientConn, error)
}

var (
	clientConnType          = reflect.TypeOf((*grpc.ClientConn)(nil))
	clientConnInterfaceType = reflect.TypeOf((*transgrpc.ClientConnInterface)(nil)).Elem()
	interfaceType           = reflect.TypeOf((*Interface)(nil)).Elem()
)

type config struct {
	Addr string `file:"addr" default:":7070" desc:"the server address in the format of host:port"`
	TLS  struct {
		ServerNameOverride string `file:"cert_file" desc:"the server name used to verify the hostname returned by the TLS handshake"`
		CAFile             string `file:"ca_file" desc:"the file containing the CA root cert file"`
		InsecureSkipVerify bool   `file:"insecure_skip_verify" desc:"skip verify"`
	} `file:"tls"`
	Singleton   bool `file:"singleton" default:"true" desc:"one client instance"`
	Block       bool `file:"block" default:"true" desc:"block until the connection is up"`
	TraceEnable bool `file:"trace_enable" default:"true"`
}

type provider struct {
	Cfg  *config
	Log  logs.Logger
	conn *grpc.ClientConn
	opts []grpc.DialOption
}

func (p *provider) Init(ctx servicehub.Context) error {
	var opts []grpc.DialOption
	if len(p.Cfg.TLS.CAFile) > 0 {
		creds, err := credentials.NewClientTLSFromFile(p.Cfg.TLS.CAFile, p.Cfg.TLS.ServerNameOverride)
		if err != nil {
			return fmt.Errorf("fail to create tls credentials %s", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		// distinguish `no tls` or `tls: insecure skip verify`
		notls := true // default no tls, compatible with old config
		if p.Cfg.TLS.InsecureSkipVerify {
			notls = false
		}
		if notls {
			opts = append(opts, grpc.WithInsecure())
		} else {
			insecureSkipVerifyTLS := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})
			opts = append(opts, grpc.WithTransportCredentials(insecureSkipVerifyTLS))
		}
	}
	if p.Cfg.TraceEnable {
		opts = append(opts,
			grpc.WithUnaryInterceptor(grpccontext.UnaryClientInterceptor()),
			grpc.WithStreamInterceptor(grpccontext.StreamClientInterceptor()),
		)
	}
	p.opts = opts
	if p.Cfg.Singleton {
		opts = nil
		if p.Cfg.Block {
			opts = append(opts, grpc.WithBlock())
		}
		conn, err := p.NewConnect(opts...)
		if err != nil {
			return fmt.Errorf("fail to dial: %s", err)
		}
		p.conn = conn
	}
	return nil
}

func (p *provider) Get() *grpc.ClientConn { return p.conn }
func (p *provider) NewConnect(opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.Dial(p.Cfg.Addr, append(opts, p.opts...)...)
}

func (p *provider) Run(ctx context.Context) error {
	if p.Cfg.Singleton {
		select {
		case <-ctx.Done():
			p.conn.Close()
			return nil
		}
	}
	return nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	if ctx.Service() == "grpc-client-connector" || ctx.Type() == interfaceType {
		return p
	}
	return p.conn
}

func init() {
	servicehub.Register("grpc-client", &servicehub.Spec{
		Services: []string{"grpc-client", "grpc-client-connector"},
		Types: []reflect.Type{
			clientConnType,
			clientConnInterfaceType,
			interfaceType,
		},
		Description: "grpc client",
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
