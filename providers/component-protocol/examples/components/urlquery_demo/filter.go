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

package urlquery_demo

import (
	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-infra/providers/component-protocol/components/filter"
	"github.com/erda-project/erda-infra/providers/component-protocol/components/filter/impl"
	"github.com/erda-project/erda-infra/providers/component-protocol/cpregister"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

type Filter struct {
	impl.DefaultFilter
}

func init() {
	cpregister.RegisterComponent("urlquery-demo", "filter", func() cptype.IComponent { return &Filter{} })
}

type UrlQueryStruct struct {
	Name string `json:"name,omitempty"`
	Age  int    `json:"age,omitempty"`
}

func (f *Filter) RegisterInitializeOp() (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) cptype.IStdStructuredPtr {
		// get url query
		var urlQuery UrlQueryStruct
		cputil.MustGetUrlQuery(sdk, &urlQuery)
		logrus.Infof("urlQuery: %+v", urlQuery)
		// set url query
		cputil.SetUrlQuery(sdk, &UrlQueryStruct{Name: "bob", Age: 20})
		return nil
	}
}

func (f *Filter) RegisterFilterOp(opData filter.OpFilter) (opFunc cptype.OperationFunc) {
	return nil
}

func (f *Filter) RegisterFilterItemSaveOp(opData filter.OpFilterItemSave) (opFunc cptype.OperationFunc) {
	return nil
}

func (f *Filter) RegisterFilterItemDeleteOp(opData filter.OpFilterItemDelete) (opFunc cptype.OperationFunc) {
	return nil
}
