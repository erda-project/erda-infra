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

package demotable

import (
	"context"
	"fmt"
	"reflect"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/protocol"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

type Interface interface {
	Render(ctx context.Context, c *cptype.Component, scenario cptype.Scenario, event cptype.ComponentEvent, gs *cptype.GlobalStateData) error
}

type config struct {
	Scenario string
	Name     string
}

type provider struct {
	Cfg *config
}

type column struct {
	DataIndex string `json:"dataIndex,omitempty"`
	Title     string `json:"title,omitempty"`
}
type tableLine struct {
	SN       string `json:"sn,omitempty"`
	Name     string `json:"name,omitempty"`
	HelloMsg string `json:"helloMsg,omitempty"`
}

// Render .
func (p *provider) Render(ctx context.Context, c *cptype.Component, scenario cptype.Scenario, event cptype.ComponentEvent, gs *cptype.GlobalStateData) error {
	sdk := cputil.SDK(ctx)
	tran := cputil.SDK(ctx).Tran
	c.Props = map[string]interface{}{
		"columns": []column{
			{DataIndex: "sn", Title: cputil.I18n(ctx, "column_serialNumber")},
			{DataIndex: "name", Title: tran.Text(cputil.Language(ctx), "column_name")},
			{DataIndex: "helloMsg", Title: sdk.I18n("column_helloMsg")},
		},
		"rowKey":          "sn",
		"pageSizeOptions": []string{"10", "20", "1000"},
	}
	c.Data = map[string]interface{}{
		"list": []tableLine{
			{SN: "1", Name: sdk.I18n("nameOfSN1"), HelloMsg: sdk.I18n("${helloMsg}: ${nameOfSN1} (%d)", 666)},
			{SN: "2", Name: sdk.I18n("nameOfSN2"), HelloMsg: sdk.I18n("${helloMsg} ${nameOfSN2}")},
		},
	}

	// print custom kv in context
	fmt.Println(ctx.Value("k1"))
	// print i18n
	fmt.Println(sdk.I18n("nameOfSN1"))

	return nil
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	compName := "demoTable"
	if ctx.Label() != "" {
		compName = ctx.Label()
	}
	protocol.MustRegisterComponent(&protocol.CompRenderSpec{
		Scenario: p.Cfg.Scenario,
		CompName: compName,
		RenderC:  func() protocol.CompRender { return &provider{} },
	})
	return nil
}

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	return p
}

func init() {
	interfaceType := reflect.TypeOf((*Interface)(nil)).Elem()
	servicehub.Register("erda.cp.components.table.demo", &servicehub.Spec{
		Types: []reflect.Type{interfaceType},
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}

// OUTPUT with header: Lang: zh
// v1
// 张三
// {demo demo} demoTable  执行完成耗时： 58.299µs
//
// OUTPUT with header: Lang: jp
// jpjpjpname
// {demo demo} demoTable  执行完成耗时： 105.248µs
