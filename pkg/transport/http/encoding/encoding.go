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

package encoding

import (
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/erda-project/erda-infra/pkg/urlenc"
)

type notSupportMediaTypeErr struct {
	text string
}

func (e notSupportMediaTypeErr) HTTPStatus() int { return http.StatusNotAcceptable }
func (e notSupportMediaTypeErr) Error() string   { return e.text }

// DecodeRequest .
func DecodeRequest(r *http.Request, out interface{}) error {
	if out == nil {
		return nil
	}
	if bytesPtr, ok := out.(*[]byte); ok {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		*bytesPtr = body
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
		if r.ContentLength <= 0 {
			return nil
		}
		if msg, ok := out.(proto.Message); ok {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return err
			}
			if len(body) <= 0 {
				return nil
			}
			return proto.Unmarshal(body, msg)
		}
	case "application/x-www-form-urlencoded", "multipart/form-data":
		if un, ok := out.(urlenc.URLValuesUnmarshaler); ok {
			err := r.ParseForm()
			if err != nil {
				return err
			}
			return un.UnmarshalURLValues("", r.Form)
		}
	default:
		if r.ContentLength <= 0 {
			return nil
		}
		if mtype == "application/json" || (strings.HasPrefix(mtype, "application/vnd.") && strings.HasSuffix(mtype, "+json")) {
			if um, ok := out.(json.Unmarshaler); ok {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					return err
				}
				return um.UnmarshalJSON(body)
			} else if msg, ok := out.(proto.Message); ok {
				body, err := io.ReadAll(r.Body)
				if err != nil {
					return err
				}
				if len(body) <= 0 {
					return nil
				}
				return protojson.Unmarshal(body, msg)
			}
			return json.NewDecoder(r.Body).Decode(out)
		}
	}
	return notSupportMediaTypeErr{text: fmt.Sprintf("not support media type: %s", mtype)}
}

// EncodeResponse .
func EncodeResponse(w http.ResponseWriter, r *http.Request, out interface{}) error {
	if out == nil {
		return nil
	}
	accept := r.Header.Get("Accept")
	var acceptAny bool
	if len(accept) > 0 {
		// TODO select MediaType of max q
		for _, item := range strings.Split(accept, ",") {
			mtype, _, err := mime.ParseMediaType(item)
			if err != nil {
				return err
			}
			if mtype == "*/*" {
				acceptAny = true
				continue
			}
			ok, err := encodeResponse(mtype, w, r, out)
			if ok {
				if err != nil {
					return err
				}
				return nil
			}
		}
	} else {
		_, err := encodeResponse("application/json", w, r, out)
		return err
	}
	if acceptAny {
		contentType := r.Header.Get("Content-Type")
		if len(contentType) > 0 {
			mtype, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				return err
			}
			ok, err := encodeResponse(mtype, w, r, out)
			if ok {
				if err != nil {
					return err
				}
				return nil
			}
		}
		_, err := encodeResponse("application/json", w, r, out)
		return err
	}
	return notSupportMediaTypeErr{text: fmt.Sprintf("not support media type: %s", accept)}
}

func encodeResponse(mtype string, w http.ResponseWriter, r *http.Request, out interface{}) (bool, error) {
	switch mtype {
	case "application/protobuf", "application/x-protobuf":
		if msg, ok := out.(proto.Message); ok {
			byts, err := proto.Marshal(msg)
			if err != nil {
				return false, err
			}
			w.Header().Set("Content-Type", "application/protobuf")
			_, err = w.Write(byts)
			return true, err
		}
	case "application/x-www-form-urlencoded", "multipart/form-data":
		if m, ok := out.(urlenc.URLValuesMarshaler); ok {
			vals := make(url.Values)
			w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
			return true, m.MarshalURLValues("", vals)
		}
	default:
		if mtype == "application/json" || (strings.HasPrefix(mtype, "application/vnd.") && strings.HasSuffix(mtype, "+json")) {
			if msg, ok := out.(proto.Message); ok {
				byts, err := protojson.Marshal(msg)
				if err != nil {
					return false, err
				}
				w.Header().Set("Content-Type", "application/json")
				_, err = w.Write(byts)
				return true, err
			}
			w.Header().Set("Content-Type", "application/json")
			return true, json.NewEncoder(w).Encode(out)
		}
	}
	return false, nil
}
