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
	"context"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
)

// SDK return cp sdk for easy use.
func SDK(ctx context.Context) (sdk *cptype.SDK) {
	v := ctx.Value(cptype.GlobalInnerKeyCtxSDK)
	if v == nil {
		return
	}
	vv, ok := v.(*cptype.SDK)
	if !ok {
		return
	}
	return vv
}

// GetInParamByKey return cp inParam by key for easy use.
func GetInParamByKey(ctx context.Context, key string) interface{} {
	sdk := SDK(ctx)
	if sdk == nil {
		return nil
	}
	if sdk.InParams == nil {
		return nil
	}
	return sdk.InParams[key]
}
