// Author: recallsong
// Email: songruiguo@qq.com

package register

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/pkg/transport"
	transgrpc "github.com/erda-project/erda-infra/pkg/transport/grpc"
	transhttp "github.com/erda-project/erda-infra/pkg/transport/http"
	"github.com/erda-project/erda-infra/providers/grpcserver"
	"github.com/erda-project/erda-infra/providers/httpserver"
	"google.golang.org/grpc"
)

// Interface .
type Interface = transport.Register

type provider struct {
	Log    logs.Logger
	router httpserver.Router
	grpc   grpcserver.Interface
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.router, _ = ctx.Service("http-server").(httpserver.Router)
	p.grpc, _ = ctx.Service("grpc-server").(grpcserver.Interface)
	if p.router == nil && p.grpc == nil {
		return fmt.Errorf("not found http-server of grpc-server")
	}
	return nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	return &service{
		name:   ctx.Caller(),
		router: p.router,
		grpc:   p.grpc,
	}
}

var _ Interface = (*service)(nil)

type service struct {
	name   string
	router httpserver.Router
	grpc   grpcserver.Interface
}

func (s *service) Add(method, path string, handler transhttp.HandlerFunc) {
	if s.router != nil {
		path = buildPath(path)
		s.router.Add(method, path, handler)
	}
}

func (s *service) RegisterService(sd *grpc.ServiceDesc, impl interface{}) {
	if s.grpc != nil {
		s.grpc.RegisterService(sd, impl)
	}
}

// buildPath convert googleapis path to echo path
func buildPath(path string) string {
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
						sb.WriteString(":" + name + strings.ReplaceAll(string(chars[begin:i]), ":", "%3A")) // replace ":" to %3A
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

func init() {
	servicehub.Register("service-register", &servicehub.Spec{
		Services: []string{"service-register"},
		Types: []reflect.Type{
			reflect.TypeOf((*Interface)(nil)).Elem(),
			reflect.TypeOf((*transgrpc.ServiceRegistrar)(nil)).Elem(),
			reflect.TypeOf((*transhttp.Router)(nil)).Elem(),
		},
		OptionalDependencies: []string{"grpc-server", "http-server"},
		Description:          "provide grpc and http server",
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
