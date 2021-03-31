// Author: recallsong
// Email: songruiguo@qq.com

package http

import (
	"encoding/json"
	"net/http"

	"github.com/erda-project/erda-infra/pkg/transport/http/encoding"
)

// DefaultHandleOptions .
func DefaultHandleOptions() *HandleOptions {
	return &HandleOptions{
		Decode: encoding.DecodeRequest,
		Encode: encoding.EncodeResponse,
		Error:  EncodeError,
	}
}

// EncodeError default EncodeErrorFunc implement
func EncodeError(w http.ResponseWriter, r *http.Request, err error) {
	// TODO optimize
	byts, _ := json.Marshal(map[string]interface{}{
		"code": "400",
		"err":  err.Error(),
	})
	w.Write(byts)
}
