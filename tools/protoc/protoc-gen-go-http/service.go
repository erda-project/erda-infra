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

package main

import (
	"fmt"
	"strconv"
	"strings"

	protocutils "github.com/erda-project/erda-infra/tools/pkg/protoc-utils"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	transportPackage = protogen.GoImportPath("github.com/erda-project/erda-infra/pkg/transport")
	contextPackage   = protogen.GoImportPath("context")
	httpPackage      = protogen.GoImportPath("net/http")
	urlencPackage    = protogen.GoImportPath("github.com/erda-project/erda-infra/pkg/urlenc")
	base64Package    = protogen.GoImportPath("encoding/base64")
	strconvPackage   = protogen.GoImportPath("strconv")
	httprulePackage  = protogen.GoImportPath("github.com/erda-project/erda-infra/pkg/transport/http/httprule")
	runtimePackage   = protogen.GoImportPath("github.com/erda-project/erda-infra/pkg/transport/http/runtime")
	fmtPackage       = protogen.GoImportPath("fmt")
	stringsPackage   = protogen.GoImportPath("strings")
)

type serviceDesc struct {
	ServiceType string
	ServiceName string
	Metadata    string
	Methods     []*methodDesc
}

type methodDesc struct {
	Name        string
	Comment     string
	Path        string
	Method      string
	QueryParams map[string][]string
	PathParams  []string
	Request     string
	Response    string
	ReqBody     string
	RespBody    string
	Meta        *protogen.Method
}

func (s *serviceDesc) execute(g *protogen.GeneratedFile) error {
	g.P("// ", s.ServiceType, "Handler is the server API for ", s.ServiceType, " service.")
	g.P("type ", s.ServiceType, "Handler interface {")
	for _, m := range s.Methods {
		if len(m.Comment) > 0 {
			g.P("	", strings.TrimSpace(m.Comment))
		}
		g.P("// ", m.Method, " ", m.Path)
		if m.Meta.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated() {
			g.P(deprecationComment)
		}
		g.P("	", m.Name, "(", contextPackage.Ident("Context"), ", *", m.Request, ") (*", m.Response, ", error)")
	}
	g.P("}")
	g.P()
	g.P("// Register", s.ServiceType, "Handler register ", s.ServiceType, "Handler to ", transhttpPackage.Ident("Router"), ".")
	g.P("func Register", s.ServiceType, "Handler(r ", transhttpPackage.Ident("Router"), ", srv ", s.ServiceType, "Handler, opts ...", transhttpPackage.Ident("HandleOption"), ") {")
	if len(s.Methods) > 0 {
		g.P("	h := ", transhttpPackage.Ident("DefaultHandleOptions"), "()")
		g.P("	for _, op := range opts {")
		g.P("		op(h)")
		g.P("	}")
		g.P("	encodeFunc := func (fn func(", httpPackage.Ident("ResponseWriter"), ", *", httpPackage.Ident("Request"), ") (interface{}, error)", ") ", transhttpPackage.Ident("HandlerFunc"), " {")
		g.P("		return func(w ", httpPackage.Ident("ResponseWriter"), ", r *", httpPackage.Ident("Request"), ") {")
		g.P("			out, err := fn(w, r)")
		g.P("			if err != nil {")
		g.P("				h.Error(w, r, err)")
		g.P("				return")
		g.P("			}")
		g.P("			if err := h.Encode(w, r, out); err != nil {")
		g.P("				h.Error(w, r, err)")
		g.P("			}")
		g.P("		}")
		g.P("	}")
		g.P()
		for _, m := range s.Methods {
			routeFunc := "add_" + m.Name
			g.P("	", routeFunc, " := func(method, path string, fn func(", contextPackage.Ident("Context"), ", *", m.Request, ") (*", m.Response, ", error)) {")
			g.P("	handler := func(ctx ", contextPackage.Ident("Context"), ", req interface{}) (interface{}, error) {")
			g.P("		return fn(ctx, req.(*", m.Request, "))")
			g.P("	}")
			infoVar := fmt.Sprintf("%s_info", m.Name)
			g.P("	var ", infoVar, " ", transportPackage.Ident("ServiceInfo"))
			g.P("	if h.Interceptor != nil {")
			g.P("		", infoVar, " = ", transportPackage.Ident("NewServiceInfo"), "(", strconv.Quote(string(m.Meta.Parent.Desc.FullName())), ",", strconv.Quote(m.Name), ", srv)")
			g.P("		handler = h.Interceptor(handler)")
			g.P("	}")
			if len(m.PathParams) > 0 {
				g.P("	compiler, _ := ", httprulePackage.Ident("Parse"), "(path)")
				g.P("	temp := compiler.Compile()")
				g.P("	pattern, _ := ", runtimePackage.Ident("NewPattern"), "(", httprulePackage.Ident("SupportPackageIsVersion1"), ", temp.OpCodes, temp.Pool, temp.Verb)")
			}
			g.P("	r.Add(method, path, encodeFunc(")
			g.P("	func(w ", httpPackage.Ident("ResponseWriter"), ", r *", httpPackage.Ident("Request"), ") (interface{}, error) {")
			g.P("		var in ", m.Request)
			if len(m.ReqBody) > 0 {
				path, _, err := protocutils.GetFieldPath(m.ReqBody, m.Meta.Input.Fields)
				if err != nil {
					return fmt.Errorf("service %q, method %q : %s", s.ServiceType, m.Name, err)
				}
				g.P("	if err := h.Decode(r, &in.", path, "); err != nil {")
			} else {
				g.P("	if err := h.Decode(r, &in); err != nil {")
			}
			g.P("			return nil, err")
			g.P("		}")
			g.P("		var input interface{} = &in")
			g.P("		if u, ok := (input).(", urlencPackage.Ident("URLValuesUnmarshaler"), "); ok {")
			g.P("			if err := u.UnmarshalURLValues(\"\", r.URL.Query()); err != nil {")
			g.P("				return nil, err")
			g.P("			}")
			g.P("		}")
			if len(m.QueryParams) > 0 {
				g.P("params := r.URL.Query()")
				for key, fields := range m.QueryParams {
					for _, name := range fields {
						names := strings.Split(name, ".")
						field, err := getField(names[0], m.Meta.Input.Fields)
						if err != nil {
							return err
						}
						g.P("if vals := params[", strconv.Quote(key), "]; len(vals) > 0 {")
						err = genVarValue(g, "in", names, field, "vals", "vals[0]")
						if err != nil {
							return fmt.Errorf("service %q, method %q : %s", s.ServiceType, m.Name, err)
						}
						g.P("}")
					}
				}
			}
			if len(m.PathParams) > 0 {
				g.P("	path := r.URL.Path")
				g.P("	if len(path) > 0 {")
				g.P("		components := ", stringsPackage.Ident("Split"), `(path[1:], "/")`)
				g.P("		last := len(components) - 1")
				g.P("		var verb string")
				g.P("		if idx := ", stringsPackage.Ident("LastIndex"), `(components[last], ":"); idx >= 0 {`)
				g.P("			c := components[last]")
				g.P("			components[last], verb = c[:idx], c[idx+1:]")
				g.P("		}")
				g.P("		vars, err := pattern.Match(components, verb)")
				g.P("		if err != nil {")
				g.P("			return nil, err")
				g.P("		}")
				g.P("		for k, val := range vars {")
				g.P("			switch k {")
				for _, name := range m.PathParams {
					_, field, err := protocutils.GetFieldPath(name, m.Meta.Input.Fields)
					if err != nil {
						return fmt.Errorf("service %q, method %q : %s", s.ServiceType, m.Name, err)
					}
					g.P("			case ", strconv.Quote(name), ":")
					if field.Desc.IsList() {
						g.P("		vals := ", stringsPackage.Ident("Split"), `(val, ",")`)
					}
					names := strings.Split(name, ".")
					field, err = getField(names[0], m.Meta.Input.Fields)
					if err != nil {
						return err
					}
					err = genVarValue(g, "in", names, field, "vals", "val")
					if err != nil {
						return err
					}
				}
				g.P("			}")
				g.P("		}")
				g.P("	}")
			}
			g.P("		ctx := ", transhttpPackage.Ident("WithRequest"), "(r.Context(), r)")
			g.P("		ctx = ", transportPackage.Ident("WithHTTPHeaderForServer"), "(ctx, r.Header)")
			g.P("		if h.Interceptor != nil {")
			g.P("			ctx = ", contextPackage.Ident("WithValue"), "(ctx, ", transportPackage.Ident("ServiceInfoContextKey"), ", ", infoVar, ")")
			g.P("		}")
			g.P("		out, err := handler(ctx, &in)")
			g.P("		if err != nil {")
			g.P("			return out, err")
			g.P("		}")
			if len(m.RespBody) > 0 {
				g.P("	if out != nil {")
				g.P("		resp := out.(*", m.Response, ")")
				path, _, err := protocutils.GetFieldPath(m.RespBody, m.Meta.Output.Fields)
				if err != nil {
					return fmt.Errorf("service %q, method %q : %s", s.ServiceType, m.Name, err)
				}
				g.P("		out = resp.", path)
				g.P("	}")
			}
			g.P("		return out, nil")
			g.P("	}),")
			g.P(")")
			g.P("}")
			g.P()
		}
		g.P()
		for _, m := range s.Methods {
			routeFunc := "add_" + m.Name
			g.P("	", routeFunc, "(", strconv.Quote(m.Method), ", ", strconv.Quote(m.Path), ", srv.", m.Name, ")")
		}
	}
	g.P("}")
	g.P()
	return nil
}

func getField(name string, fields []*protogen.Field) (field *protogen.Field, err error) {
	for _, fd := range fields {
		if string(fd.Desc.Name()) == name {
			return fd, nil
		}
	}
	return nil, fmt.Errorf("field %q not exist", name)
}

func genVarValue(g *protogen.GeneratedFile, prefix string, names []string, field *protogen.Field, listname, varname string) error {
	switch len(names) {
	case 0:
		return nil
	case 1:
		return genSetVarValue(g, prefix+"."+field.GoName, field, listname, varname)
	}
	if field.Message == nil {
		return fmt.Errorf("%s is not message type", field.Desc.Name())
	}
	name := prefix + "." + field.GoName
	g.P("if ", name, " == nil {")
	g.P("	", name, " = &", field.Message.GoIdent, "{}")
	g.P("}")
	for _, fd := range field.Message.Fields {
		if string(fd.Desc.Name()) == names[1] {
			return genVarValue(g, name, names[1:], fd, listname, varname)
		}
	}
	return fmt.Errorf("not found field %s", names[1])
}

func genSetVarValue(g *protogen.GeneratedFile, path string, field *protogen.Field, listname, varname string) error {
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		if field.Desc.IsList() {
			g.P("list := make([]bool, 0, len(", listname, "))")
			g.P("for _, text := range ", listname, " {")
			g.P("	val, err := ", strconvPackage.Ident("ParseBool"), "(text)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, val)")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", strconvPackage.Ident("ParseBool"), "(", varname, ")")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = val")
		}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		if field.Desc.IsList() {
			g.P("list := make([]int32, 0, len(", listname, "))")
			g.P("for _, text := range ", listname, " {")
			g.P("	val, err := ", strconvPackage.Ident("ParseInt"), "(text, 10, 32)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, int32(val))")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", strconvPackage.Ident("ParseInt"), "(", varname, ", 10, 32)")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = int32(val)")
		}
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		if field.Desc.IsList() {
			g.P("list := make([]uint32, 0, len(", listname, "))")
			g.P("for _, text := range ", listname, " {")
			g.P("	val, err := ", strconvPackage.Ident("ParseUint"), "(text, 10, 32)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, uint32(val))")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", strconvPackage.Ident("ParseUint"), "(", varname, ", 10, 32)")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = uint32(val)")
		}
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		if field.Desc.IsList() {
			g.P("list := make([]int64, 0, len(", listname, "))")
			g.P("for _, text := range ", listname, " {")
			g.P("	val, err := ", strconvPackage.Ident("ParseInt"), "(text, 10, 64)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, val)")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", strconvPackage.Ident("ParseInt"), "(", varname, ", 10, 64)")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = val")
		}
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		if field.Desc.IsList() {
			g.P("list := make([]uint64, 0, len(", listname, "))")
			g.P("for _, text := range ", listname, " {")
			g.P("	val, err := ", strconvPackage.Ident("ParseUint"), "(text, 10, 64)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, val)")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", strconvPackage.Ident("ParseUint"), "(", varname, ", 10, 64)")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = val")
		}
	case protoreflect.FloatKind:
		if field.Desc.IsList() {
			g.P("list := make([]float32, 0, len(", listname, "))")
			g.P("for _, text := range ", listname, " {")
			g.P("	val, err := ", strconvPackage.Ident("ParseFloat"), "(text, 32)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, float32(val))")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", strconvPackage.Ident("ParseFloat"), "(", varname, ", 32)")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = float32(val)")
		}
	case protoreflect.DoubleKind:
		if field.Desc.IsList() {
			g.P("list := make([]float64, 0, len(", listname, "))")
			g.P("for _, text := range ", listname, " {")
			g.P("	val, err := ", strconvPackage.Ident("ParseFloat"), "(text, 64)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, val)")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", strconvPackage.Ident("ParseFloat"), "(", varname, ", 64)")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = val")
		}
	case protoreflect.StringKind:
		if field.Desc.IsList() {
			g.P(path, " = ", listname)
		} else {
			g.P(path, " = ", varname)
		}
	case protoreflect.BytesKind:
		if field.Desc.IsList() {
			g.P("list := make([][]byte, 0, len(", listname, "))")
			g.P("for _, text := range ", listname, " {")
			g.P("	val, err := ", base64Package.Ident("StdEncoding.DecodeString"), "(text)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, val)")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", base64Package.Ident("StdEncoding.DecodeString"), "(", varname, ")")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = val")
		}
	default:
		return fmt.Errorf("not support type %q for query string", field.Desc.Kind())
	}
	return nil
}
