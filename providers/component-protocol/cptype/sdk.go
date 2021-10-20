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

package cptype

import (
	"github.com/erda-project/erda-infra/providers/component-protocol/protobuf/proto-go/cp/pb"
	"github.com/erda-project/erda-infra/providers/i18n"
)

// SDK .
type SDK struct {
	Tran     i18n.Translator
	Identity *pb.IdentityInfo
	InParams map[string]interface{}
	Lang     i18n.LanguageCodes
}

// I18n .
func (sdk *SDK) I18n(key string, args ...interface{}) string {
	if len(args) == 0 {
		try := sdk.Tran.Text(sdk.Lang, key)
		if try != key {
			return try
		}
	}
	return sdk.Tran.Sprintf(sdk.Lang, key, args...)
}
