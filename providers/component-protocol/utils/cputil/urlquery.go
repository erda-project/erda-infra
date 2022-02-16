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

package cputil

import (
	"encoding/base64"
	"fmt"

	"github.com/erda-project/erda-infra/pkg/strutil"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

const (
	compURLQuerySuffix = "__urlQuery"
)

// MakeCompURLQueryKey make url query key for component.
func MakeCompURLQueryKey(compName string) string { return compName + compURLQuerySuffix }

// SetURLQuery set data to url query.
func SetURLQuery(sdk *cptype.SDK, data interface{}) {
	b, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Errorf("failed to set url query, err: %v", err))
	}
	urlQueryStr := base64.URLEncoding.EncodeToString(b)
	// set into comp options
	if sdk.Comp.Options == nil {
		sdk.Comp.Options = &cptype.ComponentOptions{}
	}
	sdk.Comp.Options.URLQuery = urlQueryStr
}

// GetURLQuery get component's url query and parse to `resultStructPtr`.
func GetURLQuery(sdk *cptype.SDK, resultStructPtr interface{}) error {
	if sdk.InParams == nil {
		return nil
	}
	if resultStructPtr == nil {
		return fmt.Errorf("result receiver pointer can't be nil")
	}
	encodedURLQuery := strutil.String(sdk.InParams[MakeCompURLQueryKey(sdk.Comp.Name)])
	jsonEncodedURLQuery, err := base64.URLEncoding.DecodeString(encodedURLQuery)
	if err != nil {
		return fmt.Errorf("failed to get url query from inParams, err: %v", err)
	}
	err = json.Unmarshal(jsonEncodedURLQuery, resultStructPtr)
	if err != nil {
		return fmt.Errorf("failed to json unmarshal json encoded url query, err: %v", err)
	}
	return nil
}

// MustGetURLQuery must GetURLQuery.
func MustGetURLQuery(sdk *cptype.SDK, resultStructPtr interface{}) {
	err := GetURLQuery(sdk, resultStructPtr)
	if err != nil {
		panic(err)
	}
}
