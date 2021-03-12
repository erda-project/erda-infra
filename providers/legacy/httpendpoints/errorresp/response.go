package errorresp

import (
	"encoding/json"
	"net/http"

	"github.com/erda-project/erda-infra/providers/legacy/httpendpoints"
	"github.com/erda-project/erda-infra/providers/legacy/httpendpoints/i18n"
)

// ToResp 根据 APIError 转为一个 http error response.
func (e *APIError) ToResp() httpendpoints.Responser {
	return &httpendpoints.HTTPResponse{
		Error:  e,
		Status: e.httpCode,
		Content: httpendpoints.Resp{
			Success: false,
			Err: httpendpoints.ErrorResponse{
				Code: e.code,
				Msg:  e.msg,
			},
		},
	}
}

// ErrResp 根据 error 转为一个 http error response.
func ErrResp(e error) (httpendpoints.Responser, error) {
	switch t := e.(type) {
	case *APIError:
		return e.(*APIError).ToResp(), nil
	default:
		_ = t
		return New().InternalError(e).ToResp(), nil
	}
}

// Write 将错误写入 http.ResponseWriter
func (e *APIError) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(e.httpCode)
	return json.NewEncoder(w).Encode(httpendpoints.Resp{
		Success: false,
		Err: httpendpoints.ErrorResponse{
			Code: e.code,
			Msg:  e.Render(i18n.NewNopLocaleResource()),
		},
	})
}

// ErrWrite 根据 error 写入标准错误格式
func ErrWrite(e error, w http.ResponseWriter) error {
	switch t := e.(type) {
	case *APIError:
		return e.(*APIError).Write(w)
	default:
		_ = t
		return New().InternalError(e).Write(w)
	}
}
