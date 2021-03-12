package httpendpoints

import (
	"context"
	"net/http"

	"github.com/erda-project/erda-infra/providers/legacy/httpendpoints/i18n"
	"github.com/gorilla/mux"
)

// Endpoint contains URL path and corresponding handler
type Endpoint struct {
	Path           string
	Method         string
	Handler        func(context.Context, *http.Request, map[string]string) (Responser, error)
	WriterHandler  func(context.Context, http.ResponseWriter, *http.Request, map[string]string) error
	ReverseHandler func(context.Context, *http.Request, map[string]string) error
}

// Interface .
type Interface interface {
	RegisterEndpoint(endpoints []Endpoint)
	Router() *mux.Router
}

// Responser is an interface for http response
type Responser interface {
	GetLocaledResp(locale i18n.LocaleResource) HTTPResponse
	GetStatus() int
	GetContent() interface{}
}
