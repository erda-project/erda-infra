package component_protocol

import (
	"os"
	"path/filepath"
	"reflect"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/pkg/transport"
	"github.com/erda-project/erda-infra/providers/component-protocol/definition"
	"github.com/erda-project/erda-infra/providers/i18n"
	"github.com/erda-project/erda-proto-go/cp/pb"
)

type config struct {
	DefaultProtocolYamlScanBasePaths []string `file:"default_protocol_yaml_scan_base_paths" env:"DEFAULT_PROTOCOL_YAML_SCAN_BASE_PATHS"`
}

// +provider
type provider struct {
	Cfg      *config
	Log      logs.Logger
	Register transport.Register

	Tran             i18n.Translator
	CustomContextKVs map[interface{}]interface{}

	protocolService *protocolService
}

func (p *provider) Init(ctx servicehub.Context) error {
	p.CustomContextKVs = make(map[interface{}]interface{})
	p.protocolService = &protocolService{p: p}
	if p.Register != nil {
		pb.RegisterCPServiceImp(p.Register, p.protocolService)
	}

	// register default protocol yaml files
	for _, basePath := range p.Cfg.DefaultProtocolYamlScanBasePaths {
		pwd, _ := os.Getwd()
		absPath := filepath.Join(pwd, basePath)
		definition.InitDefaultCompProtocols(absPath)
	}
	for key := range definition.DefaultProtocols {
		p.Log.Infof("default protocol registered for scenario: %s", key)
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
