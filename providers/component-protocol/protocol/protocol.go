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

package protocol

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/erda-project/erda-infra/pkg/strutil"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

// defaultProtocols contains all default protocols.
// map key: scenarioKey
// map value: default protocol
var defaultProtocols = make(map[string]cptype.ComponentProtocol)
var defaultProtocolsRaw = make(map[string]string)

// RegisterDefaultProtocols register protocol contents.
func RegisterDefaultProtocols(protocolYAMLs ...[]byte) {
	for idx, protocolYAML := range protocolYAMLs {
		var p cptype.ComponentProtocol
		if err := yaml.Unmarshal(protocolYAML, &p); err != nil {
			panic(fmt.Errorf("failed to parse protocol yaml, index: %d, err: %v", idx, err))
		}
		if p.Scenario == "" {
			continue
		}
		defaultProtocols[p.Scenario] = p
		if CpPlaceHolderRe.Match(protocolYAML) {
			defaultProtocolsRaw[p.Scenario] = string(protocolYAML)
		}
		logrus.Infof("default protocol registered for scenario: %s", p.Scenario)
	}
}

// RegisterDefaultProtocolsFromBasePath register default component protocols under base path.
// default path: libs/erda-configs/permission
func RegisterDefaultProtocolsFromBasePath(basePath string) {
	var err error
	defer func() {
		if err != nil {
			logrus.Errorf("failed to register default component protocol, err: %v", err)
			panic(err)
		}
	}()
	rd, err := ioutil.ReadDir(basePath)
	if err != nil {
		return
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := basePath + "/" + fi.Name()
			RegisterDefaultProtocolsFromBasePath(fullDir)
		} else {
			if fi.Name() != "protocol.yml" && fi.Name() != "protocol.yaml" {
				continue
			}
			fullName := basePath + "/" + fi.Name()
			yamlFile, er := ioutil.ReadFile(fullName)
			if er != nil {
				err = er
				return
			}
			var p cptype.ComponentProtocol
			if er := yaml.Unmarshal(yamlFile, &p); er != nil {
				err = er
				return
			}
			defaultProtocols[p.Scenario] = p
			logrus.Infof("default protocol registered for scenario: %s", p.Scenario)
		}
	}
}

// getDefaultProtocol get default protocol by scenario.
func getDefaultProtocol(ctx context.Context, scenario string) (cptype.ComponentProtocol, error) {
	rawYamlStr, ok := defaultProtocolsRaw[scenario]
	if !ok {
		// protocol not have cp placeholder
		p, ok := defaultProtocols[scenario]
		if !ok {
			return cptype.ComponentProtocol{}, fmt.Errorf(i18n(ctx, "${default.protocol.not.exist}, ${scenario}: %s", scenario))
		}
		return p, nil
	}
	lang := cputil.Language(ctx)
	tran := ctx.Value(cptype.GlobalInnerKeyCtxSDK).(*cptype.SDK).Tran
	replaced := strutil.ReplaceAllStringSubmatchFunc(CpPlaceHolderRe, rawYamlStr, func(v []string) string {
		if len(v) == 2 && strings.HasPrefix(v[1], I18n+".") {
			key := strings.TrimPrefix(v[1], I18n+".")
			if len(key) > 0 {
				return tran.Text(lang, key)
			}
		}
		return v[0]
	})
	var p cptype.ComponentProtocol
	if err := yaml.Unmarshal([]byte(replaced), &p); err != nil {
		return cptype.ComponentProtocol{}, fmt.Errorf("failed to parse protocol yaml i18n, err: %v", err)
	}
	return p, nil
}

// getProtoComp .
func getProtoComp(ctx context.Context, p *cptype.ComponentProtocol, compName string) (c *cptype.Component, err error) {
	if p.Components == nil {
		err = fmt.Errorf("empty protocol components")
		return
	}

	c, ok := p.Components[compName]
	if !ok {
		defaultProtocol, err := getDefaultProtocol(ctx, p.Scenario)
		if err != nil {
			return c, err
		}
		c, ok = defaultProtocol.Components[compName]
		if !ok {
			err = fmt.Errorf("empty component [%s] in default protocol", compName)
		}
		p.Components[compName] = c
		return c, err
	}
	return
}

// getProtoCompStateValue .
func getProtoCompStateValue(ctx context.Context, p *cptype.ComponentProtocol, compName, sk string) (interface{}, error) {
	c, err := getProtoComp(ctx, p, compName)
	if err != nil {
		return nil, err
	}
	v, err := getCompStateKV(c, sk)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// polishProtocol .
func polishProtocol(req *cptype.ComponentProtocol) {
	if req == nil {
		return
	}
	// polish component name
	for name, component := range req.Components {
		component.Name = name
	}
}
