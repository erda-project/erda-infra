// Author: recallsong
// Email: songruiguo@qq.com

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

type define struct{}

func (d *define) Services() []string {
	return []string{"grpc-server"}
}
func (d *define) Types() []reflect.Type {
	return []reflect.Type{
		reflect.TypeOf((*grpc.Server)(nil)),
		reflect.TypeOf((*Interface)(nil)).Elem(),
	}
}
func (d *define) Description() string { return "grpc server" }
func (d *define) Config() interface{} { return &config{} }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		p := &provider{}
		return p
	}
}

// config .
type config struct {
	Addr string `file:"addr" default:":7800" desc:"http address to listen"`
	TLS  struct {
		CertFile string `file:"cert_file"`
		KeyFile  string `file:"key_file"`
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
	servicehub.RegisterProvider("grpc-server", &define{})
}
