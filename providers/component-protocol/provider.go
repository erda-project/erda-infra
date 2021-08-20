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

package componentprotocol

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

// Init .
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
		definition.RegisterDefaultCompProtocols(absPath)
	}
	for key := range definition.DefaultProtocols {
		p.Log.Infof("default protocol registered for scenario: %s", key)
	}

	return nil
}

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
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
