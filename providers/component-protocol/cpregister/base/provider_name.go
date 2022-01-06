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

package base

import (
	"fmt"
	"strings"

	"github.com/erda-project/erda-infra/pkg/strutil"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

const (
	componentProviderNamePrefix = "component-protocol.components."
)

var (
	componentProviderDefaultNamespacePrefix = componentProviderNamePrefix + cptype.DefaultComponentNamespace + "."
)

// MustGetScenarioAndCompNameFromProviderKey .
func MustGetScenarioAndCompNameFromProviderKey(providerKey string) (scenario, compName, instanceName string) {
	scenario, compName, instanceName, err := GetScenarioAndCompNameFromProviderKey(providerKey)
	if err != nil {
		panic(err)
	}

	return scenario, compName, instanceName
}

// GetScenarioAndCompNameFromProviderKey .
func GetScenarioAndCompNameFromProviderKey(providerKey string) (scenario, compName, instanceName string, err error) {
	// validate prefix
	if !strutil.HasPrefixes(providerKey, componentProviderNamePrefix, componentProviderDefaultNamespacePrefix) {
		return "", "", "", fmt.Errorf("invalid prefix")
	}
	// parse as std comp providerKey
	ss := strings.SplitN(providerKey, ".", 4)
	if len(ss) != 4 {
		return "", "", "", fmt.Errorf("not standard provider key: %s", providerKey)
	}
	scenario = ss[2]
	// default namespace doesn't belong to any scenario
	if scenario == cptype.DefaultComponentNamespace {
		scenario = ""
	}
	// split comp and instance name
	compName, instanceName = splitCompAndInstance(ss[3])

	return
}

// splitCompAndInstance split compPartKey to compName and instanceName.
func splitCompAndInstance(compPartKey string) (compName, instanceName string) {
	vv := strings.SplitN(compPartKey, "@", 2)
	if len(vv) == 2 {
		compName = vv[0]
		instanceName = vv[1]
	} else {
		compName = compPartKey
		instanceName = compName
	}
	return
}

// MakeComponentProviderName .
func MakeComponentProviderName(scenario, compType string) string {
	return fmt.Sprintf("%s%s.%s", componentProviderNamePrefix, scenario, compType)
}
