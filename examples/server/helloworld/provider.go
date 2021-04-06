package example

import (
	"context"
	"fmt"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/examples/protocol/pb"
	"github.com/erda-project/erda-infra/pkg/transport"
	"github.com/erda-project/erda-infra/pkg/transport/interceptor"
)

type config struct {
}

type provider struct {
	Cfg      *config
	Log      logs.Logger
	Register transport.Register
}

func (p *provider) Init(ctx servicehub.Context) error {
	// TODO initialize something ...

	greeterService := &greeterService{p}
	pb.RegisterGreeterServiceImp(p.Register, greeterService,
		transport.WithInterceptors(func(h interceptor.Handler) interceptor.Handler {
			fmt.Println("wrap greeterService methods")
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				info := ctx.Value(transport.ServiceInfoContextKey).(transport.ServiceInfo)
				fmt.Printf("before %s/%s\n", info.Service(), info.Method())
				fmt.Println(req)
				out, err := h(ctx, req)
				fmt.Printf("after %s/%s\n", info.Service(), info.Method())
				return out, err
			}
		}),
	)

	userService := &userService{p}
	pb.RegisterUserServiceImp(p.Register, userService,
		transport.WithInterceptors(func(h interceptor.Handler) interceptor.Handler {
			fmt.Println("wrap userService methods")
			return func(ctx context.Context, req interface{}) (interface{}, error) {
				info := ctx.Value(transport.ServiceInfoContextKey).(transport.ServiceInfo)
				fmt.Printf("before %s/%s\n", info.Service(), info.Method())
				fmt.Println(req)
				out, err := h(ctx, req)
				fmt.Printf("after %s/%s\n", info.Service(), info.Method())
				return out, err
			}
		}),
	)

	return nil
}

func init() {
	servicehub.Register("erda.infra.example", &servicehub.Spec{
		Services:     pb.ServiceNames(),
		Types:        pb.Types(),
		Dependencies: []string{"service-register"},
		Description:  "",
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
