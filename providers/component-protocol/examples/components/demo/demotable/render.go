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

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

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
