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

// TODO switch to servicehub listener mechanism later, now this mechanism is missing.

package translator

import (
	"bytes"
	_ "embed" // embed
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/recallsong/go-utils/reflectx"
	"github.com/sirupsen/logrus"

	cfg "github.com/erda-project/erda-infra/pkg/config"
	"github.com/erda-project/erda-infra/providers/i18n"
)

// InternalI18nConfigs contains all protocl internal i18n configs.
//
//go:embed i18n-cp-internal.yaml
var InternalI18nConfigs string

// Tran is a translator.
type Tran struct {
	dic map[string]map[string]string
}

// NewInternalTranslator .
func NewInternalTranslator() *Tran {
	// make embed content as a temp file
	f, _ := ioutil.TempFile(os.TempDir(), "*.yaml")
	defer func() {
		if err := f.Close(); err != nil {
			logrus.Errorf("failed to close i18n config file, err: %v", err)
		}
	}()
	if _, err := f.WriteString(InternalI18nConfigs); err != nil {
		panic(fmt.Errorf("failed to write i18n config content to temp file, err: %v", err))
	}
	fPath, err := filepath.Abs(f.Name())
	if err != nil {
		panic(fmt.Errorf("failed to get abs path of temp i18n config file, err: %v", err))
	}
	defer func() {
		if err := os.RemoveAll(fPath); err != nil {
			logrus.Errorf("failed to remove temp i18n config file, err: %v", err)
		}
	}()

	// load to dic
	dic := make(map[string]map[string]string)
	if err := loadToDic(fPath, dic); err != nil {
		panic(err)
	}
	return &Tran{dic: dic}
}

// Get .
func (t *Tran) Get(lang i18n.LanguageCodes, key, def string) string {
	text := t.getText(lang, key)
	if len(text) > 0 {
		return text
	}
	return def
}

// Text .
func (t *Tran) Text(lang i18n.LanguageCodes, key string) string {
	text := t.getText(lang, key)
	if len(text) > 0 {
		return text
	}
	return key
}

// Sprintf .
func (t *Tran) Sprintf(lang i18n.LanguageCodes, key string, args ...interface{}) string {
	return fmt.Sprintf(t.escape(lang, key), args...)
}

func (t *Tran) getText(langs i18n.LanguageCodes, key string) string {
	key = strings.ToLower(key)
	for _, lang := range langs {
		if t.dic != nil {
			text := t.dic[lang.Code]
			if text != nil {
				if value, ok := text[key]; ok {
					return value
				}
			}
			text = t.dic[lang.RestrictedCode()]
			if text != nil {
				if value, ok := text[key]; ok {
					return value
				}
			}
		}
	}
	return ""
}

var regExp = regexp.MustCompile(`\$\{([^:}]*)(:[^}]*)?\}`)

func (t *Tran) escape(lang i18n.LanguageCodes, text string) string {
	contents := reflectx.StringToBytes(text)
	params := regExp.FindAllSubmatch(contents, -1)
	for _, param := range params {
		if len(param) != 3 {
			continue
		}
		var key, defval []byte = param[1], nil
		if len(param[2]) > 0 {
			defval = param[2][1:]
		}
		k := reflectx.BytesToString(key)
		val := t.getText(lang, k)
		if len(val) <= 0 {
			val = strings.Trim(reflectx.BytesToString(defval), `"`)
		}
		contents = bytes.Replace(contents, param[0], reflectx.StringToBytes(val), 1)
	}
	return reflectx.BytesToString(contents)
}

func loadToDic(file string, dic map[string]map[string]string) error {
	m := make(map[string]interface{})
	err := cfg.LoadToMap(file, m)
	if err != nil {
		return fmt.Errorf("fail to load i18n file: %s", err)
	}
	for lang, v := range m {
		text := dic[lang]
		if text == nil {
			text = make(map[string]string)
			dic[lang] = text
		}
		switch m := v.(type) {
		case map[string]string:
			for k, v := range m {
				text[strings.ToLower(k)] = fmt.Sprint(v)
			}
		case map[string]interface{}:
			for k, v := range m {
				text[strings.ToLower(k)] = fmt.Sprint(v)
			}
		case map[interface{}]interface{}:
			for k, v := range m {
				text[strings.ToLower(fmt.Sprint(k))] = fmt.Sprint(v)
			}
		default:
			return fmt.Errorf("invalid i18n file format: %s", file)
		}
	}
	return nil
}
