package example

import (
	logs "github.com/erda-project/erda-infra/base/logs"
	servicehub "github.com/erda-project/erda-infra/base/servicehub"
	pb "github.com/erda-project/erda-infra/examples/protocol/pb"
	serviceregister "github.com/erda-project/erda-infra/providers/serviceregister"
)

type config struct {
}

type provider struct {
	Cfg      *config
	Log      logs.Logger
	Register serviceregister.Interface
}

func (p *provider) Init(ctx servicehub.Context) error {
	// TODO initialize something ...

	greeterService := &greeterService{}
	userService := &userService{}

	pb.RegisterServices(p.Register, p.Register,
		greeterService,
		userService,
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
