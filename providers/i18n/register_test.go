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

package i18n

import (
	"embed"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed examples/i18n
var i18nFS embed.FS

func TestRegisterFilesFromFS(t *testing.T) {
	p := &provider{
		common: make(map[string]map[string]string),
		dic:    make(map[string]map[string]map[string]string),
	}
	_ = p.RegisterFilesFromFS("examples/i18n", i18nFS)

	assert.Equal(t, "名字来自common", p.common["zh"]["common name"])
	assert.Equal(t, "名字来自file", p.dic["hello"]["zh"]["file name"])
}
