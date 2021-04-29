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

package servicehub

import (
	"reflect"
	"testing"
)

type testSpecDefine struct {
	services     []string
	types        []reflect.Type
	dependencies []string
	summary      string
	description  string
}

func (d *testSpecDefine) Services() []string         { return d.services }
func (d *testSpecDefine) Types() []reflect.Type      { return d.types }
func (d *testSpecDefine) Summary() string            { return d.summary }
func (d *testSpecDefine) Description() string        { return d.description }
func (d *testSpecDefine) Dependencies(*Hub) []string { return d.dependencies }

func Test_specDefine_Services(t *testing.T) {
	tests := []struct {
		name string
		spec *Spec
		want []string
	}{
		{
			name: "empty",
			spec: &Spec{},
			want: nil,
		},
		{
			name: "direct",
			spec: &Spec{Services: []string{"test-service"}},
			want: []string{"test-service"},
		},
		{
			name: "override",
			spec: &Spec{
				Define: &testSpecDefine{services: []string{"test-service-2"}},
			},
			want: []string{"test-service-2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &specDefine{s: tt.spec}
			if got := d.Services(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("specDefine.Services() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_specDefine_Types(t *testing.T) {
	testType := reflect.TypeOf((*testSpecDefine)(nil))
	tests := []struct {
		name string
		spec *Spec
		want []reflect.Type
	}{
		{
			name: "empty",
			spec: &Spec{},
			want: nil,
		},
		{
			name: "direct",
			spec: &Spec{Types: []reflect.Type{testType}},
			want: []reflect.Type{testType},
		},
		{
			name: "override",
			spec: &Spec{
				Define: &testSpecDefine{types: []reflect.Type{testType}},
			},
			want: []reflect.Type{testType},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &specDefine{s: tt.spec}
			if got := d.Types(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("specDefine.Types() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_specDefine_Dependencies(t *testing.T) {
	hub := &Hub{servicesMap: make(map[string][]*providerContext)}
	hub.servicesMap["test-service"] = []*providerContext{nil}
	hub.servicesMap["test-service-1"] = []*providerContext{nil}
	tests := []struct {
		name string
		spec *Spec
		hub  *Hub
		want []string
	}{
		{
			name: "empty",
			spec: &Spec{},
			hub:  hub,
			want: nil,
		},
		{
			name: "dependencies",
			spec: &Spec{Dependencies: []string{"test-service"}},
			hub:  hub,
			want: []string{"test-service"},
		},
		{
			name: "optional dependencies",
			spec: &Spec{
				OptionalDependencies: []string{"test-service", "test-service-3"},
			},
			hub:  hub,
			want: []string{"test-service"},
		},
		{
			name: "merge dependencies",
			spec: &Spec{
				Dependencies:         []string{"test-service"},
				OptionalDependencies: []string{"test-service-1", "test-service-2"},
			},
			hub:  hub,
			want: []string{"test-service", "test-service-1"},
		},
		{
			name: "override",
			spec: &Spec{
				Define: &testSpecDefine{dependencies: []string{"test-service-override"}},
			},
			hub:  hub,
			want: []string{"test-service-override"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &specDefine{s: tt.spec}
			if got := d.Dependencies(tt.hub); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("specDefine.Dependencies() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_specDefine_Summary(t *testing.T) {
	tests := []struct {
		name string
		spec *Spec
		want string
	}{
		{
			name: "empty",
			spec: &Spec{},
			want: "",
		},
		{
			name: "direct",
			spec: &Spec{
				Summary: "test-summary",
			},
			want: "test-summary",
		},
		{
			name: "override",
			spec: &Spec{
				Define: &testSpecDefine{summary: "test-summary-override"},
			},
			want: "test-summary-override",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &specDefine{
				s: tt.spec,
			}
			if got := d.Summary(); got != tt.want {
				t.Errorf("specDefine.Summary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_specDefine_Description(t *testing.T) {
	tests := []struct {
		name string
		spec *Spec
		want string
	}{
		{
			name: "empty",
			spec: &Spec{},
			want: "",
		},
		{
			name: "direct",
			spec: &Spec{
				Description: "test-description",
			},
			want: "test-description",
		},
		{
			name: "override",
			spec: &Spec{
				Define: &testSpecDefine{description: "test-description-override"},
			},
			want: "test-description-override",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &specDefine{
				s: tt.spec,
			}
			if got := d.Description(); got != tt.want {
				t.Errorf("specDefine.Description() = %v, want %v", got, tt.want)
			}
		})
	}
}
