package httpendpoints

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/erda-project/erda-infra/providers/legacy/httpendpoints/i18n"
	"github.com/erda-project/erda-infra/providers/legacy/httpendpoints/ierror"
	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
)

const (
	// ContentTypeJSON Content Type
	ContentTypeJSON = "application/json"
	// ResponseWriter context value key
	ResponseWriter = "responseWriter"
	// Base64EncodedRequestBody .
	Base64EncodedRequestBody = "base64-encoded-request-body"
	// TraceID .
	TraceID = "dice-trace-id"
)

// RegisterEndpoints match URL path to corresponding handler
func (p *provider) RegisterEndpoints(endpoints []Endpoint) {
	for _, ep := range endpoints {
		if ep.WriterHandler != nil {
			p.router.Path(ep.Path).Methods(ep.Method).HandlerFunc(p.internalWriterHandler(ep.WriterHandler))
		} else if ep.ReverseHandler != nil {
			p.router.Path(ep.Path).Methods(ep.Method).Handler(p.internalReverseHandler(ep.ReverseHandler))
		} else {
			p.router.Path(ep.Path).Methods(ep.Method).HandlerFunc(p.internal(ep.Handler))
		}
		p.L.Infof("Added endpoint: %s %s", ep.Method, ep.Path)
	}
}

func (p *provider) internal(handler func(context.Context, *http.Request, map[string]string) (Responser, error)) http.HandlerFunc {
	pctx := context.Background()
	pctx = injectTraceID(pctx)

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		p.L.Debugf("start %s %s", r.Method, r.URL.String())

		ctx, cancel := context.WithCancel(pctx)
		defer func() {
			cancel()
			p.L.Debugf("finished handle request %s %s (took %v)", r.Method, r.URL.String(), time.Since(start))
		}()
		ctx = context.WithValue(ctx, ResponseWriter, w)

		handleRequest(r)

		langs := i18n.Language(r)
		locale := i18n.WrapLocaleResource(p.t, langs)
		response, err := handler(ctx, r, mux.Vars(r))
		if err == nil {
			response = response.GetLocaledResp(locale)
		}
		if err != nil {
			apiError, isApiError := err.(ierror.IAPIError)
			if isApiError {
				response = HTTPResponse{
					Status: apiError.HttpCode(),
					Content: Resp{
						Success: false,
						Err: ErrorResponse{
							Code: apiError.Code(),
							Msg:  apiError.Render(locale),
						},
					},
				}
			} else {
				p.L.Errorf("failed to handle request: %s (%v)", r.URL.String(), err)

				statusCode := http.StatusInternalServerError
				if response != nil {
					statusCode = response.GetStatus()
				}
				w.WriteHeader(statusCode)
				io.WriteString(w, err.Error())
				return
			}
		}

		w.Header().Set("Content-Type", ContentTypeJSON)
		w.WriteHeader(response.GetStatus())

		encoder := json.NewEncoder(w)
		vals := r.URL.Query()
		pretty, ok := vals["pretty"]
		if ok && strings.Compare(pretty[0], "true") == 0 {
			encoder.SetIndent("", "    ")
		}

		if err := encoder.Encode(response.GetContent()); err != nil {
			p.L.Errorf("failed to send response: %s (%v)", r.URL.String(), err)
			return
		}
	}
}

func (p *provider) internalWriterHandler(handler func(context.Context, http.ResponseWriter, *http.Request, map[string]string) error) http.HandlerFunc {
	pctx := context.Background()
	pctx = injectTraceID(pctx)

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		p.L.Debugf("start %s %s", r.Method, r.URL.String())

		ctx, cancel := context.WithCancel(pctx)
		defer func() {
			cancel()
			p.L.Debugf("finished handle request %s %s (took %v)", r.Method, r.URL.String(), time.Since(start))
		}()

		handleRequest(r)

		err := handler(ctx, w, r, mux.Vars(r))
		if err != nil {
			p.L.Errorf("failed to handle request: %s (%v)", r.URL.String(), err)

			statusCode := http.StatusInternalServerError
			w.WriteHeader(statusCode)
			io.WriteString(w, err.Error())
		}
	}
}

// internalReverseHandler .
func (p *provider) internalReverseHandler(handler func(context.Context, *http.Request, map[string]string) error) http.Handler {
	pctx := context.Background()
	pctx = injectTraceID(pctx)

	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			start := time.Now()
			p.L.Debugf("start %s %s", r.Method, r.URL.String())

			ctx, cancel := context.WithCancel(pctx)
			defer func() {
				cancel()
				p.L.Debugf("finished handle request %s %s (took %v)", r.Method, r.URL.String(), time.Since(start))
			}()

			handleRequest(r)

			err := handler(ctx, r, mux.Vars(r))
			if err != nil {
				p.L.Errorf("failed to handle request: %s (%v)", r.URL.String(), err)
				return
			}
		},
		FlushInterval: -1,
	}
}

func handleRequest(r *http.Request) {
	// base64 decode request body if declared in header
	if strings.EqualFold(r.Header.Get(Base64EncodedRequestBody), "true") {
		r.Body = ioutil.NopCloser(base64.NewDecoder(base64.StdEncoding, r.Body))
	}
}

func injectTraceID(ctx context.Context) context.Context {
	id, _ := uuid.NewV4()
	return context.WithValue(ctx, TraceID, id.String())
}
