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

package dependency

import (
	"reflect"

	"github.com/erda-project/erda-infra/base/servicehub"
)

// Interface .
type Interface interface {
	Hello(name string) string
}

type provider struct{}

func (p *provider) Hello(name string) string {
	return "hello " + name
}

func init() {
	servicehub.Register("example-dependency-provider", &servicehub.Spec{
		Services:    []string{"example-dependency"},
		Types:       []reflect.Type{reflect.TypeOf((*Interface)(nil)).Elem()},
		Description: "dependency for example",
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
