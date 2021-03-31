// Author: recallsong
// Email: songruiguo@qq.com

package encoding

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/erda-project/erda-infra/pkg/urlenc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// DecodeRequest .
func DecodeRequest(r *http.Request, out interface{}) error {
	if out == nil {
		return nil
	}
	contentType := r.Header.Get("Content-Type")
	if len(contentType) <= 0 {
		return nil
	}
	mtype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return err
	}
	switch mtype {
	case "application/protobuf", "application/x-protobuf":
		if msg, ok := out.(proto.Message); ok {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return err
			}
			return proto.Unmarshal(body, msg)
		}
	case "application/json":
		if msg, ok := out.(proto.Message); ok {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				return err
			}
			return protojson.Unmarshal(body, msg)
		}
		return json.NewDecoder(r.Body).Decode(out)
	case "application/x-www-form-urlencoded", "multipart/form-data":
		if un, ok := out.(urlenc.URLValuesUnmarshaler); ok {
			err := r.ParseForm()
			if err != nil {
				return err
			}
			return un.UnmarshalURLValues("", r.Form)
		}
	}
	return fmt.Errorf("not support Unmarshal type %s with %s", reflect.TypeOf(out).Name(), mtype)
}

// EncodeResponse .
func EncodeResponse(w http.ResponseWriter, r *http.Request, out interface{}) error {
	accept := r.Header.Get("Accept")
	if len(accept) > 0 {
		// TODO select MediaType of max q
		for _, item := range strings.Split(accept, ",") {
			mtype, _, err := mime.ParseMediaType(item)
			if err != nil {
				return err
			}
			if mtype == "*/*" || mtype == "" {
				continue
			}
			accept = mtype
			break
		}
	}
	if len(accept) <= 0 {
		contentType := r.Header.Get("Content-Type")
		if len(contentType) <= 0 {
			return nil
		}
		mtype, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			return err
		}
		accept = mtype
	}
	switch accept {
	case "application/protobuf", "application/x-protobuf":
		if msg, ok := out.(proto.Message); ok {
			byts, err := proto.Marshal(msg)
			if err != nil {
				return err
			}
			_, err = w.Write(byts)
			return err
		}
	case "application/json":
		if msg, ok := out.(proto.Message); ok {
			byts, err := protojson.Marshal(msg)
			if err != nil {
				return err
			}
			_, err = w.Write(byts)
			return err
		}
		return json.NewEncoder(w).Encode(out)
	case "application/x-www-form-urlencoded", "multipart/form-data":
		if m, ok := out.(urlenc.URLValuesMarshaler); ok {
			vals := make(url.Values)
			return m.MarshalURLValues("", vals)
		}
	}
	return fmt.Errorf("not support Marshal type %s with %s", reflect.TypeOf(out).Name(), accept)
}
