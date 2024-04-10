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

package httpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/recallsong/go-utils/errorx"
	"github.com/recallsong/go-utils/reflectx"

	"github.com/erda-project/erda-infra/providers/httpserver/server"
)

type (
	// Response .
	Response interface {
		Status(Context) int
		ReadCloser(Context) io.ReadCloser
		Error(Context) error
	}
	// ResponseGetter .
	ResponseGetter interface {
		Response(ctx Context) Response
	}
	// Interceptor .
	Interceptor func(handler func(ctx Context) error) func(ctx Context) error
)

func getInterceptors(options []interface{}) []server.MiddlewareFunc {
	var list []server.MiddlewareFunc
	for _, opt := range options {
		var inter Interceptor
		switch val := opt.(type) {
		case Interceptor:
			inter = val
		case func(handler func(ctx Context) error) func(ctx Context) error:
			inter = Interceptor(val)
		case server.MiddlewareFunc:
			list = append(list, val)
		case func(server.HandlerFunc) server.HandlerFunc:
			list = append(list, val)
		default:
			continue
		}
		if inter != nil {
			list = append(list, func(fn server.HandlerFunc) server.HandlerFunc {
				handler := inter(func(ctx Context) error {
					return fn(ctx.(*context))
				})
				return func(ctx server.Context) error {
					return handler(ctx.(*context))
				}
			})
		}
	}
	return list
}

func (r *router) add(method, path string, handler interface{}, inters []server.MiddlewareFunc, outer server.MiddlewareFunc) server.HandlerFunc {
	var echoHandler server.HandlerFunc
	switch fn := handler.(type) {
	case server.HandlerFunc:
		echoHandler = fn
	case func(server.Context) error:
		echoHandler = server.HandlerFunc(fn)
	case func(server.Context):
		echoHandler = server.HandlerFunc(func(ctx server.Context) error {
			fn(ctx)
			return nil
		})
	case http.HandlerFunc:
		echoHandler = server.HandlerFunc(func(ctx server.Context) error {
			fn(ctx.Response(), ctx.Request())
			return nil
		})
	case func(http.ResponseWriter, *http.Request):
		echoHandler = server.HandlerFunc(func(ctx server.Context) error {
			fn(ctx.Response(), ctx.Request())
			return nil
		})
	case func(*http.Request, http.ResponseWriter):
		echoHandler = server.HandlerFunc(func(ctx server.Context) error {
			fn(ctx.Request(), ctx.Response())
			return nil
		})
	case http.Handler:
		echoHandler = server.HandlerFunc(func(ctx server.Context) error {
			fn.ServeHTTP(ctx.Response(), ctx.Request())
			return nil
		})
	default:
		echoHandler = r.handlerWrap(handler)
		if echoHandler == nil {
			panic(fmt.Errorf("%s %s: not support http server handler type: %v", method, path, handler))
		}
	}
	if outer != nil {
		list := make([]server.MiddlewareFunc, 1+len(r.interceptors)+len(inters))
		list[0] = outer
		copy(list[1:], r.interceptors)
		copy(list[1+len(r.interceptors):], inters)
		inters = list
	} else {
		inters = append(r.interceptors[0:len(r.interceptors):len(r.interceptors)], inters...)
	}
	if len(inters) > 0 {
		handler := echoHandler
		for i := len(inters) - 1; i >= 0; i-- {
			handler = inters[i](handler)
		}
		echoHandler = handler
	}
	r.tx.Add(method, path, echoHandler)
	return echoHandler
}

var (
	readerType      = reflect.TypeOf((*io.Reader)(nil)).Elem()
	readCloserType  = reflect.TypeOf((*io.ReadCloser)(nil)).Elem()
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	requestType     = reflect.TypeOf((*http.Request)(nil))
	responseType    = reflect.TypeOf((*http.ResponseWriter)(nil)).Elem()
	echoContextType = reflect.TypeOf((*server.Context)(nil)).Elem()
	contextType     = reflect.TypeOf((*Context)(nil)).Elem()
	interfaceType   = reflect.TypeOf((*interface{})(nil)).Elem()
)

func (r *router) handlerWrap(handler interface{}) server.HandlerFunc {
	typ := reflect.TypeOf(handler)
	if typ.Kind() == reflect.Func {
		val := reflect.ValueOf(handler)
		var argGets []func(ctx server.Context) (interface{}, error)
		argNum := typ.NumIn()
		for i := 0; i < argNum; i++ {
			argTyp := typ.In(i)
			getter := argGetter(argTyp)
			if getter == nil {
				return nil
			}
			argGets = append(argGets, getter)
		}
		retNum := typ.NumOut()
		if retNum > 3 {
			return nil
		}
		var retGet func(values []reflect.Value) (*int, io.ReadCloser, io.Reader, interface{}, error)
		var retIndex [5]*int
		var hasRet bool
		for i := 0; i < retNum; i++ {
			retTyp := typ.Out(i)
			index := i
			if retTyp.Kind() == reflect.Int {
				if retIndex[0] == nil {
					retIndex[0] = &index
					hasRet = true
					continue
				}
			} else if retTyp.AssignableTo(readCloserType) {
				if retIndex[1] == nil {
					retIndex[1] = &index
					hasRet = true
					continue
				}
			} else if retTyp.AssignableTo(readerType) {
				if retIndex[2] == nil {
					retIndex[2] = &index
					hasRet = true
					continue
				}
			} else if retTyp == errorType {
				if retIndex[3] == nil {
					retIndex[3] = &index
					hasRet = true
					continue
				}
			} else if retTyp == interfaceType {
				if retIndex[4] == nil {
					retIndex[4] = &index
					hasRet = true
					continue
				}
			}
			return nil
		}
		if hasRet {
			retGet = func(values []reflect.Value) (status *int, readerCloser io.ReadCloser, reader io.Reader, data interface{}, err error) {
				if retIndex[0] != nil {
					val := int(values[*retIndex[0]].Int())
					status = &val
				}
				if retIndex[1] != nil {
					val := values[*retIndex[1]].Interface()
					readerCloser = val.(io.ReadCloser)
				}
				if retIndex[2] != nil {
					val := values[*retIndex[2]].Interface()
					reader = val.(io.Reader)
				}
				if retIndex[3] != nil {
					val := values[*retIndex[3]].Interface()
					if val != nil {
						err = val.(error)
					}
				}
				if retIndex[4] != nil {
					data = values[*retIndex[4]].Interface()
				}
				return
			}
		}
		return server.HandlerFunc(func(ctx server.Context) error {
			var values []reflect.Value
			for _, getter := range argGets {
				val, err := getter(ctx)
				if err != nil {
					if _, ok := err.(validator.ValidationErrors); ok {
						//TODO: custom error encode
						return ctx.JSON(400, map[string]interface{}{
							"success": false,
							"err": map[string]interface{}{
								"code": "400",
								"msg":  err.Error(),
							},
						})
					}
					if herr, ok := err.(*echo.HTTPError); ok {
						if http.StatusBadRequest <= herr.Code && herr.Code < http.StatusInternalServerError {
							//TODO: custom error encode
							ctx.JSON(400, map[string]interface{}{
								"success": false,
								"err": map[string]interface{}{
									"code": strconv.Itoa(herr.Code),
									"msg":  herr.Message,
								},
							})
						}
					}
					return err
				}
				value := reflect.ValueOf(val)
				values = append(values, value)
			}
			returns := val.Call(values)
			if retGet == nil {
				return nil
			}
			status, readCloser, reader, data, err := retGet(returns)
			if data != nil {
				var resp Response
				context := ctx.(Context)
				switch val := data.(type) {
				case ResponseGetter:
					resp = val.Response(context)
				case Response:
					resp = val
				}
				if resp != nil {
					rc := resp.ReadCloser(context)
					if rc != nil {
						readCloser = rc
					}
					statusCode := resp.Status(context)
					if statusCode > 0 {
						status = &statusCode
					}
					e := resp.Error(context)
					if e != nil {
						err = e
					}
				}
			}
			if status != nil {
				ctx.Response().WriteHeader(*status)
			}
			var errs errorx.Errors
			if err != nil {
				errs = append(errs, err)
			}
			if readCloser != nil {
				defer readCloser.Close()
				_, err = io.Copy(ctx.Response(), readCloser)
				if err != nil {
					errs = append(errs, err)
				}
			} else if reader != nil {
				_, err = io.Copy(ctx.Response(), reader)
				if err != nil {
					errs = append(errs, err)
				}
			} else if data != nil {
				switch val := data.(type) {
				case string:
					_, err = ctx.Response().Write(reflectx.StringToBytes(val))
				case []byte:
					_, err = ctx.Response().Write(val)
				default:
					err = json.NewEncoder(ctx.Response()).Encode(data)
				}
				if err != nil {
					errs = append(errs, err)
				}
			}
			return errs.MaybeUnwrap()
		})
	}
	return nil
}

func argGetter(argTyp reflect.Type) func(ctx server.Context) (interface{}, error) {
	if argTyp == requestType {
		return requestGetter
	} else if argTyp == responseType {
		return responseGetter
	} else if argTyp == contextType || argTyp == echoContextType {
		return contextGetter
	} else {
		kind := argTyp.Kind()
		if kind == reflect.String {
			return requestBodyStirngGetter
		} else if kind == reflect.Slice && argTyp.Elem().Kind() == reflect.Uint8 {
			return requestBodyBytesGetter
		}
		typ := argTyp
		for kind == reflect.Ptr {
			typ = typ.Elem()
			kind = typ.Kind()
		}
		switch kind {
		case reflect.Struct:
			var validate bool
			for i, num := 0, typ.NumField(); i < num; i++ {
				if len(typ.Field(i).Tag.Get("validate")) > 0 {
					validate = true
					break
				}
			}
			return requestDataBind(argTyp, validate)
		case reflect.Map, reflect.Interface:
			return requestDataBind(argTyp, false)
		case reflect.String:
			return requestBodyStirngGetter
		case reflect.Bool,
			reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.Array, reflect.Slice:
			return requestValuesGetter(argTyp)
		default:
			return nil
		}
	}
}

func requestGetter(ctx server.Context) (interface{}, error)  { return ctx.Request(), nil }
func responseGetter(ctx server.Context) (interface{}, error) { return ctx.Response(), nil }
func contextGetter(ctx server.Context) (interface{}, error)  { return ctx, nil }
func requestDataBind(typ reflect.Type, validate bool) func(server.Context) (interface{}, error) {
	return func(ctx server.Context) (data interface{}, err error) {
		outVal := reflect.New(typ)
		if typ.Kind() != reflect.Ptr {
			data = outVal.Interface()
			err = ctx.Bind(data)
		} else {
			eval := outVal.Elem()
			etype := typ.Elem()
			for etype.Kind() == reflect.Ptr {
				v := reflect.New(etype)
				eval.Set(v)
				eval = v.Elem()
				etype = etype.Elem()
			}
			switch etype.Kind() {
			case reflect.Map:
				v := reflect.New(etype)
				v.Elem().Set(reflect.MakeMap(etype))
				eval.Set(v)
			case reflect.Slice:
				v := reflect.New(etype)
				v.Elem().Set(reflect.MakeSlice(etype, 0, 0))
				eval.Set(v)
			default:
				eval.Set(reflect.New(etype))
			}
			data = eval.Interface()
			err = ctx.Bind(data)
		}
		if err != nil {
			return nil, err
		}
		if validate {
			err = ctx.Validate(data)
			if err != nil {
				return nil, err
			}
		}
		return outVal.Elem().Interface(), nil
	}
}
func requestValuesGetter(typ reflect.Type) func(ctx server.Context) (interface{}, error) {
	return func(ctx server.Context) (interface{}, error) {
		out := reflect.New(typ)
		byts, err := io.ReadAll(ctx.Request().Body)
		if err != nil {
			return nil, fmt.Errorf("fail to read body: %s", err)
		}
		ctx.Request().Body = io.NopCloser(bytes.NewBuffer(byts))
		err = json.Unmarshal(byts, out.Interface())
		if err != nil {
			return nil, fmt.Errorf("fail to Unmarshal body: %s", err)
		}
		return out.Elem().Interface(), nil
	}
}
func requestBodyBytesGetter(ctx server.Context) (interface{}, error) {
	byts, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil, fmt.Errorf("fail to read body: %s", err)
	}
	ctx.Request().Body = io.NopCloser(bytes.NewBuffer(byts))
	return byts, nil
}

func requestBodyStirngGetter(ctx server.Context) (interface{}, error) {
	byts, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return "", fmt.Errorf("fail to read body: %s", err)
	}
	return reflectx.BytesToString(byts), nil
}

type structValidator struct {
	validator *validator.Validate
}

// Validate .
func (v *structValidator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
