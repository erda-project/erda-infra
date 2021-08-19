package httpserver

import (
	"net/http"

	"github.com/erda-project/erda-infra/providers/httpserver/server"
)

// Context handler context.
type Context interface {
	SetAttribute(key string, val interface{})
	Attribute(key string) interface{}
	Attributes() map[string]interface{}
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
	Param(name string) string
	ParamNames() []string
}

var _ server.Context = (*context)(nil)

type context struct {
	server.Context
	data map[string]interface{}
	vars map[string]string
}

func (c *context) SetAttribute(key string, val interface{}) {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	c.data[key] = val
}

func (c *context) Attribute(key string) interface{} {
	if c.data == nil {
		return nil
	}
	return c.data[key]
}

func (c *context) Attributes() map[string]interface{} {
	if c.data == nil {
		c.data = make(map[string]interface{})
	}
	return c.data
}

func (c *context) ResponseWriter() http.ResponseWriter {
	return c.Context.Response()
}

func (c *context) Bind(i interface{}) error {
	return c.Echo().Binder.Bind(i, c)
}
