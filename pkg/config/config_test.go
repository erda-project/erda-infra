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
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
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

func Test_polishBuffer(t *testing.T) {
	const envLocale = "LOCALE"
	content := `
i18n:
  LOCALE: "${LOCALE:zh-CN}"
  num: 1
`

	// init buf
	buf := bytes.NewBufferString(content)
	// polish buf
	assert.NoError(t, polishBuffer(buf))
	// unmarshal
	cfg := map[string]interface{}{}
	err := yaml.Unmarshal(buf.Bytes(), &cfg)
	assert.NoError(t, err)
	i18nCfg, ok := cfg["i18n"]
	assert.True(t, ok)
	c, ok := i18nCfg.(map[string]interface{})
	assert.True(t, ok)
	assert.True(t, c[envLocale] == "zh-CN")
	assert.True(t, c["num"] == 1)

	// set env
	assert.NoError(t, os.Setenv(envLocale, "en-US"))
	defer func() { _ = os.Unsetenv(envLocale) }()
	// init buf
	buf = bytes.NewBufferString(content)
	// polish buf
	assert.NoError(t, polishBuffer(buf))
	// unmarshal
	cfg = map[string]interface{}{}
	err = yaml.Unmarshal(buf.Bytes(), &cfg)
	assert.NoError(t, err)
	i18nCfg, ok = cfg["i18n"]
	assert.True(t, ok)
	c, ok = i18nCfg.(map[string]interface{})
	assert.True(t, ok)
	assert.True(t, c[envLocale] == "en-US")
	assert.True(t, c["num"] == 1)
}

func TestUnmarshalToMap(t *testing.T) {
	const envLocale = "LOCALE"
	content := `
i18n:
  LOCALE: "${LOCALE:zh-CN}"
  num: 1
`

	// init buf
	buf := bytes.NewBufferString(content)
	cfg := make(map[string]interface{})
	// parse conf
	assert.NoError(t, UnmarshalToMap(buf, "yaml", cfg))
	i18nCfg, ok := cfg["i18n"]
	assert.True(t, ok)
	c, ok := i18nCfg.(map[string]interface{})
	assert.True(t, ok)
	assert.True(t, c[envLocale] == "zh-CN")
	assert.True(t, c["num"] == 1)

	// set env
	assert.NoError(t, os.Setenv(envLocale, "en-US"))
	defer func() { _ = os.Unsetenv(envLocale) }()
	// init buf
	buf = bytes.NewBufferString(content)
	cfg = make(map[string]interface{})
	// parse conf
	assert.NoError(t, UnmarshalToMap(buf, "yaml", cfg))
	i18nCfg, ok = cfg["i18n"]
	assert.True(t, ok)
	c, ok = i18nCfg.(map[string]interface{})
	assert.True(t, ok)
	assert.True(t, c[envLocale] == "en-US")
	assert.True(t, c["num"] == 1)
}
