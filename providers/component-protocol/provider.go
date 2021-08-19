package component_protocol

import (
	"reflect"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/pkg/transport"
	"github.com/erda-project/erda-infra/providers/i18n"
	"github.com/erda-project/erda-proto-go/cp/pb"
)

type config struct {
}

// +provider
type provider struct {
	Cfg      *config
	Log      logs.Logger
	Register transport.Register
	Tran     i18n.Translator

	protocolService *protocolService
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.protocolService = &protocolService{p: p}
	if p.Register != nil {
		pb.RegisterCPServiceImp(p.Register, p.protocolService)
	}
	return nil
}

func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	//switch {
	//case ctx.Service() == "erda.cp.CPService" || ctx.Type() == pb.CPServiceServerType() || ctx.Type() == pb.CPServiceHandlerType():
	//	return p.protocolService
	//}
	return p
}

func init() {
	interfaceType := reflect.TypeOf((*Interface)(nil)).Elem()
	servicehub.Register("erda.cp", &servicehub.Spec{
		Services:             pb.ServiceNames(),
		Types:                append(pb.Types(), interfaceType),
		OptionalDependencies: []string{"service-register"},
		Description:          "",
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
