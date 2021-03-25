package errorresp

import (
	"strings"

	"github.com/erda-project/erda-infra/providers/legacy/httpendpoints/i18n"
)

// MetaMessage .
type MetaMessage struct {
	Key     string
	Args    []interface{}
	Default string
}

// APIError .
type APIError struct {
	httpCode           int
	code               string
	msg                string
	localeMetaMessages []MetaMessage
}

// Error .
func (e *APIError) Error() string {
	if e.msg == "" {
		e.Render(i18n.NewNopLocaleResource())
	}
	return e.msg
}

// Code .
func (e *APIError) Code() string {
	return e.code
}

// HTTPCode ...
func (e *APIError) HTTPCode() int {
	return e.httpCode
}

// Option 可选参数
type Option func(*APIError)

// New .
func New(options ...Option) *APIError {
	e := &APIError{}
	for _, op := range options {
		op(e)
	}
	return e
}

// WithMessage .
func WithMessage(msg string) Option {
	return func(a *APIError) {
		a.msg = msg
	}
}

// WithTemplateMessage .
func WithTemplateMessage(template, defaultValue string, args ...interface{}) Option {
	return func(a *APIError) {
		_ = a.appendMeta(template, defaultValue, args...)
	}
}

// WithCode .
func WithCode(httpCode int, code string) Option {
	return func(a *APIError) {
		_ = a.appendCode(httpCode, code)
	}
}

func (e *APIError) appendCode(httpCode int, code string) *APIError {
	e.httpCode = httpCode
	e.code = code
	return e
}

func (e *APIError) appendMsg(template *i18n.Template, args ...interface{}) *APIError {
	msg := template.Render(args...)
	if e.msg == "" {
		e.msg = msg
		return e
	}
	e.msg = strings.Join([]string{e.msg, ": ", msg}, "")
	return e
}

func (e *APIError) appendMeta(key string, defaultContent string, args ...interface{}) *APIError {
	e.localeMetaMessages = append(e.localeMetaMessages, MetaMessage{
		Key:     key,
		Args:    args,
		Default: defaultContent,
	})
	return e
}

func (e *APIError) appendLocaleTemplate(template *i18n.Template, args ...interface{}) *APIError {
	e.localeMetaMessages = append(e.localeMetaMessages, MetaMessage{
		Key:     template.Key(),
		Args:    args,
		Default: template.Content(),
	})
	return e
}

// Render .
func (e *APIError) Render(localeResource i18n.LocaleResource) string {
	for _, metaMessage := range e.localeMetaMessages {
		var template *i18n.Template
		if !localeResource.ExistKey(metaMessage.Key) && metaMessage.Default != "" {
			template = i18n.NewTemplate(metaMessage.Key, metaMessage.Default)
		} else {
			template = localeResource.GetTemplate(metaMessage.Key)
		}
		msg := template.Render(metaMessage.Args...)
		if e.msg == "" {
			e.msg = msg
		} else {
			e.msg = strings.Join([]string{e.msg, ": ", msg}, "")
		}
	}
	return e.msg
}

func (e *APIError) dup() *APIError {
	return &APIError{
		httpCode:           e.httpCode,
		code:               e.code,
		msg:                e.msg,
		localeMetaMessages: e.localeMetaMessages,
	}
}
