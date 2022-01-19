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

package protocol

import (
	"context"
	"regexp"

	"github.com/erda-project/erda-infra/providers/component-protocol/protocol/translator"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
	pi18n "github.com/erda-project/erda-infra/providers/i18n"
)

var tran pi18n.Translator

func init() {
	tran = translator.NewInternalTranslator()
}

func i18n(ctx context.Context, key string, args ...interface{}) string {
	if len(args) == 0 {
		try := tran.Text(cputil.Language(ctx), key)
		if try != key {
			return try
		}
	}
	return tran.Sprintf(cputil.Language(ctx), key, args...)
}

var CpRe = regexp.MustCompile(`\${{[ ]{1}([^{}\s]+)[ ]{1}}}`)

const I18n = "i18n"
