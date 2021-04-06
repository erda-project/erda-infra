// Author: recallsong
// Email: songruiguo@qq.com

package grpcclient

import (
	"context"
	"fmt"
	"reflect"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	transgrpc "github.com/erda-project/erda-infra/pkg/transport/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Interface .
type Interface interface {
	Get() *grpc.ClientConn
	NewConnect(opts ...grpc.DialOption) *grpc.ClientConn
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
	} `file:"tls"`
	Singleton bool `file:"singleton" default:"true" desc:"one client instance"`
	Block     bool `file:"block" default:"true" desc:"block until the connection is up"`
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
		opts = append(opts, grpc.WithInsecure())
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
