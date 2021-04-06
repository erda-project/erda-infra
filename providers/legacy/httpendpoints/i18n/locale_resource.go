// Copyright 2021 Terminus
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
	"net/http"

	i18nprovider "github.com/erda-project/erda-infra/providers/i18n"
)

// LocaleResource .
type LocaleResource interface {
	// Name() string
	ExistKey(key string) bool
	Get(key string, defaults ...string) string
	GetTemplate(key string) *Template
}

type localeResource struct {
	// name string
	t     i18nprovider.Translator
	langs i18nprovider.LanguageCodes
}

// func (lr *localeResource) Name() string { return lr.name }

func (lr *localeResource) ExistKey(key string) bool {
	text := lr.t.Get(lr.langs, key, "")
	return len(text) > 0
}

func (lr *localeResource) Get(key string, defaults ...string) string {
	if len(defaults) > 0 {
		return lr.t.Get(lr.langs, key, defaults[0])
	}
	return lr.t.Text(lr.langs, key)
}

func (lr *localeResource) GetTemplate(key string) *Template {
	content := lr.t.Text(lr.langs, key)
	return NewTemplate(key, content)
}

// WrapLocaleResource .
func WrapLocaleResource(t i18nprovider.Translator, langs i18nprovider.LanguageCodes) LocaleResource {
	return &localeResource{
		t:     t,
		langs: langs,
	}
}

type nopLocaleResource struct{}

func (lr *nopLocaleResource) ExistKey(key string) bool { return false }
func (lr *nopLocaleResource) Get(key string, defaults ...string) string {
	if len(defaults) > 0 {
		return defaults[0]
	}
	return key
}
func (lr *nopLocaleResource) GetTemplate(key string) *Template {
	return NewTemplate(key, key)
}

// NewNopLocaleResource .
func NewNopLocaleResource() LocaleResource {
	return &nopLocaleResource{}
}

// Language .
func Language(r *http.Request) i18nprovider.LanguageCodes {
	lang := r.Header.Get("Lang")
	if len(lang) <= 0 {
		lang = r.Header.Get("Accept-Language")
	}
	langs, _ := i18nprovider.ParseLanguageCode(lang)
	return langs
}
