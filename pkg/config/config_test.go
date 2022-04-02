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

package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Multi_Non_Bool_Env_Replace(t *testing.T) {
	content := `
example:
	name: "${ES|CK:foo}"
`
	want := `
example:
	name: "ck"
`
	_ = os.Setenv("CK", "ck")
	result := EscapeEnv([]byte(content))

	assert.Equal(t, want, string(result))
	_ = os.Unsetenv("CK")

	want = `
example:
	name: "foo"
`
	result = EscapeEnv([]byte(content))

	assert.Equal(t, want, string(result))

	want = `
example:
	name: "es"
`

	_ = os.Setenv("ES", "es")
	result = EscapeEnv([]byte(content))

	assert.Equal(t, want, string(result))
	_ = os.Unsetenv("ES")

	want = `
example:
	name: "es"
`

	_ = os.Setenv("ES", "es")
	_ = os.Setenv("CK", "ck")
	result = EscapeEnv([]byte(content))

	assert.Equal(t, want, string(result))
	_ = os.Unsetenv("ES")
	_ = os.Unsetenv("CK")

}

func Test_Multi_Bool_Env_Replace(t *testing.T) {
	content := `
example:
	name: "${ES|CK:default}"
`
	want := `
example:
	name: "false"
`
	_ = os.Setenv("CK", "false")
	result := EscapeEnv([]byte(content))

	assert.Equal(t, want, string(result))
	_ = os.Unsetenv("CK")

	want = `
example:
	name: "false"
`

	_ = os.Setenv("ES", "false")
	result = EscapeEnv([]byte(content))

	assert.Equal(t, want, string(result))
	_ = os.Unsetenv("ES")

	want = `
example:
	name: "True"
`
	_ = os.Setenv("ES", "false")
	_ = os.Setenv("CK", "True")
	result = EscapeEnv([]byte(content))

	assert.Equal(t, want, string(result))
	_ = os.Unsetenv("ES")
	_ = os.Unsetenv("CK")

}

func Test_Multi_Mix_Bool_Env_Replace(t *testing.T) {
	content := `
example:
	name: "${ES|CK:default}"
`

	want := `
example:
	name: "false"
`
	_ = os.Setenv("ES", "false")
	_ = os.Setenv("CK", "ddd")
	result := EscapeEnv([]byte(content))

	assert.Equal(t, want, string(result))
	_ = os.Unsetenv("ES")
	_ = os.Unsetenv("CK")

}
