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

package example

import (
	"github.com/erda-project/erda-infra/base/servicehub"
)

// Interface .
type Interface interface {
	Hello(name string) string
	Add(a, b int) int
}

var _ Interface = (*provider)(nil) // check interface implemented

type provider struct{}

func (p *provider) Hello(name string) string {
	return "hello " + name
}

func (p *provider) Add(a, b int) int {
	return a + b
}

func (p *provider) sub(a, b int) int {
	return a - b
}

func init() {
	servicehub.Register("example-provider", &servicehub.Spec{
		Services:    []string{"example"},
		Description: "example",
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
