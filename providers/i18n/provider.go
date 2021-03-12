// Author: recallsong
// Email: songruiguo@qq.com

package i18n

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	cfg "github.com/recallsong/go-utils/config"
	"github.com/recallsong/go-utils/reflectx"
)

// Translator .
type Translator interface {
	Get(lang LanguageCodes, key, def string) string
	Text(lang LanguageCodes, key string) string
	Sprintf(lang LanguageCodes, key string, args ...interface{}) string
}

// I18n .
type I18n interface {
	Get(namespace string, lang LanguageCodes, key, def string) string
	Text(namespace string, lang LanguageCodes, key string) string
	Sprintf(namespace string, lang LanguageCodes, key string, args ...interface{}) string
	Translator(namespace string) Translator
}

// NopTranslator .
type NopTranslator struct{}

// Text .
func (t *NopTranslator) Text(lang LanguageCodes, key string) string { return key }

// Sprintf .
func (t *NopTranslator) Sprintf(lang LanguageCodes, key string, args ...interface{}) string {
	return fmt.Sprintf(key, args...)
}

type define struct{}

func (d *define) Service() []string   { return []string{"i18n"} }
func (d *define) Description() string { return "i18n" }
func (d *define) Config() interface{} { return &config{} }
func (d *define) Creator() servicehub.Creator {
	return func() servicehub.Provider {
		return &provider{
			common: make(map[string]map[string]string),
			dic:    make(map[string]map[string]map[string]string),
		}
	}
}

type config struct {
	Files  []string `file:"files"`
	Common []string `file:"common"`
}

type provider struct {
	C      *config
	L      logs.Logger
	common map[string]map[string]string
	dic    map[string]map[string]map[string]string
}

func (p *provider) Init(ctx servicehub.Context) error {
	log := ctx.Logger()
	for _, file := range p.C.Common {
		f, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("fail to load i18n file: %s", err)
		}
		if f.IsDir() {
			err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				if skipFile(path, info, err) {
					return nil
				}
				return p.loadToDic(file, p.common)
			})
			if err != nil {
				return err
			}
			continue
		}
		err = p.loadToDic(file, p.common)
		if err != nil {
			return err
		}
	}
	for _, file := range p.C.Files {
		f, err := os.Stat(file)
		if err != nil {
			return fmt.Errorf("fail to load i18n file: %s", err)
		}
		if f.IsDir() {
			err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				if skipFile(path, info, err) {
					return nil
				}
				return p.loadI18nFile(path)
			})
			if err != nil {
				return err
			}
			continue
		}
		err = p.loadI18nFile(file)
		if err != nil {
			return err
		}
	}
	log.Infof("load i18n files: %v, %v", p.C.Common, p.C.Files)
	return nil
}

func skipFile(path string, info os.FileInfo, err error) bool {
	if err != nil || info == nil || info.IsDir() {
		return true
	}
	if strings.HasPrefix(filepath.Base(path), ".") {
		return true
	}
	return false
}

func (p *provider) loadI18nFile(file string) error {
	base := filepath.Base(file)
	name := base[0 : len(base)-len(filepath.Ext(base))]
	dic := p.dic[name]
	if dic == nil {
		dic = make(map[string]map[string]string)
		p.dic[name] = dic
	}
	err := p.loadToDic(file, dic)
	if err != nil {
		return err
	}
	return nil
}

func (p *provider) loadToDic(file string, dic map[string]map[string]string) error {
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

func (p *provider) Text(namespace string, lang LanguageCodes, key string) string {
	return p.Translator(namespace).Text(lang, key)
}

func (p *provider) Sprintf(namespace string, lang LanguageCodes, key string, args ...interface{}) string {
	return p.Translator(namespace).Sprintf(lang, key, args...)
}

func (p *provider) Get(namespace string, lang LanguageCodes, key, def string) string {
	return p.Translator(namespace).Get(lang, key, def)
}

func (p *provider) Translator(namespace string) Translator {
	return &translator{
		common: p.common,
		dic:    p.dic[namespace],
	}
}

type translator struct {
	common map[string]map[string]string
	dic    map[string]map[string]string
}

func (t *translator) Text(lang LanguageCodes, key string) string {
	text := t.getText(lang, key)
	if len(text) > 0 {
		return text
	}
	return key
}

func (t *translator) Sprintf(lang LanguageCodes, key string, args ...interface{}) string {
	return fmt.Sprintf(t.escape(lang, key), args...)
}

func (t *translator) Get(lang LanguageCodes, key, def string) string {
	text := t.getText(lang, key)
	if len(text) > 0 {
		return text
	}
	return def
}

func (t *translator) getText(langs LanguageCodes, key string) string {
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
		text := t.common[lang.Code]
		if text != nil {
			if value, ok := text[key]; ok {
				return value
			}
		}
		text = t.common[lang.RestrictedCode()]
		if text != nil {
			if value, ok := text[key]; ok {
				return value
			}
		}
	}
	return ""
}

var regExp = regexp.MustCompile(`\$\{(\w+)(:[^}]*)?\}|\$\{([^}]*)?\}`)

func (t *translator) escape(lang LanguageCodes, text string) string {
	contents := reflectx.StringToBytes(text)
	params := regExp.FindAllSubmatch(contents, -1)
	for _, param := range params {
		if len(param) != 4 {
			continue
		}
		var key, defval []byte
		if len(param[1]) > 0 {
			key = param[1]
			if len(param[2]) > 0 {
				defval = param[2][1:]
			}
		} else if len(param[3]) > 0 {
			key = param[3]
		} else {
			continue
		}
		k := reflectx.BytesToString(key)
		val := t.getText(lang, k)
		if len(val) <= 0 {
			if len(param[3]) > 0 {
				val = k
			} else {
				val = strings.Trim(reflectx.BytesToString(defval), `"`)
			}
		}
		contents = bytes.Replace(contents, param[0], reflectx.StringToBytes(val), 1)
	}
	return reflectx.BytesToString(contents)
}

func init() {
	servicehub.RegisterProvider("i18n", &define{})
}
