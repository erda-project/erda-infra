// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"fmt"
	"strconv"
	"strings"

	protocutils "github.com/erda-project/erda-infra/pkg/protoc-utils"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
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
	g.P("// Register", s.ServiceType, "Handler register ", s.ServiceType, "Handler to ", transportPackage.Ident("Router"), ".")
	g.P("func Register", s.ServiceType, "Handler(r ", transportPackage.Ident("Router"), ", srv ", s.ServiceType, "Handler, opts ...", transportPackage.Ident("HandleOption"), ") {")
	g.P("	h := ", transportPackage.Ident("DefaultHandleOptions"), "()")
	g.P("	for _, op := range opts {")
	g.P("		op(h)")
	g.P("	}")
	g.P("	type ConvertFunc func(", httpPackage.Ident("ResponseWriter"), ", *", httpPackage.Ident("Request"), ") (interface{}, error)")
	g.P("	encodeFunc := func (fn ConvertFunc) ", transportPackage.Ident("HandlerFunc"), " {")
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
		g.P("	", s.convertHandlerFuncName(m), " := func(fn func(", contextPackage.Ident("Context"), ", *", m.Request, ") (*", m.Response, ", error)) ConvertFunc {")
		g.P("	handler := func(ctx ", contextPackage.Ident("Context"), ", req interface{}) (interface{}, error) {")
		g.P("		return fn(ctx, req.(*", m.Request, "))")
		g.P("	}")
		g.P("	if h.Interceptor != nil {")
		g.P("		handler = h.Interceptor(handler)")
		g.P("	}")
		g.P("	return func(w ", httpPackage.Ident("ResponseWriter"), ", r *", httpPackage.Ident("Request"), ") (interface{}, error) {")
		g.P("		var in ", m.Request)
		if len(m.ReqBody) > 0 {
			path, err := protocutils.GetFieldPath(m.ReqBody, m.Meta.Input.Fields)
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
					var field *protogen.Field
					for _, fd := range m.Meta.Input.Fields {
						if string(fd.Desc.Name()) == names[0] {
							field = fd
							break
						}
					}
					if field == nil {
						return fmt.Errorf("field %q not exist", name)
					}
					g.P("if vals := params[", strconv.Quote(key), "]; len(vals) > 0 {")
					err := genQueryString(g, "in", names, field)
					if err != nil {
						return fmt.Errorf("service %q, method %q : %s", s.ServiceType, m.Name, err)
					}
					g.P("}")
				}
			}
		}
		g.P("		out, err := handler(r.Context(), &in)")
		g.P("		if err != nil {")
		g.P("			return out, err")
		g.P("		}")
		if len(m.RespBody) > 0 {
			g.P("	if out != nil {")
			g.P("		resp := out.(*", m.Response, ")")
			path, err := protocutils.GetFieldPath(m.RespBody, m.Meta.Output.Fields)
			if err != nil {
				return fmt.Errorf("service %q, method %q : %s", s.ServiceType, m.Name, err)
			}
			g.P("		out = resp.", path)
			g.P("	}")
		}
		g.P("		return out, nil")
		g.P("	}")
		g.P("}")
		g.P()
	}
	g.P()
	for _, m := range s.Methods {
		g.P("	r.Add(", strconv.Quote(m.Method), ", ", strconv.Quote(m.Path), ", encodeFunc(", s.convertHandlerFuncName(m), "(srv.", m.Name, ")))")
	}
	g.P("}")
	g.P()
	return nil
}

func (s *serviceDesc) convertHandlerFuncName(m *methodDesc) string {
	return "convert_" + m.Name + "_to_HandlerFunc"
}

func genQueryString(g *protogen.GeneratedFile, prefix string, names []string, field *protogen.Field) error {
	switch len(names) {
	case 0:
		return nil
	case 1:
		return genQueryStringValue(g, prefix+"."+field.GoName, field)
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
			return genQueryString(g, name, names[1:], fd)
		}
	}
	return fmt.Errorf("not found field %s", names[1])
}

func genQueryStringValue(g *protogen.GeneratedFile, path string, field *protogen.Field) error {
	switch field.Desc.Kind() {
	case protoreflect.BoolKind:
		if field.Desc.IsList() {
			g.P("list := make([]bool, 0, len(vals))")
			g.P("for _, text := range vals {")
			g.P("	val, err := ", strconvPackage.Ident("ParseBool"), "(text)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, val)")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", strconvPackage.Ident("ParseBool"), "(vals[0])")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = val")
		}
	case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind:
		if field.Desc.IsList() {
			g.P("list := make([]int32, 0, len(vals))")
			g.P("for _, text := range vals {")
			g.P("	val, err := ", strconvPackage.Ident("ParseInt"), "(text, 10, 32)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, int32(val))")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", strconvPackage.Ident("ParseInt"), "(vals[0], 10, 32)")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = int32(val)")
		}
	case protoreflect.Uint32Kind, protoreflect.Fixed32Kind:
		if field.Desc.IsList() {
			g.P("list := make([]uint32, 0, len(vals))")
			g.P("for _, text := range vals {")
			g.P("	val, err := ", strconvPackage.Ident("ParseUint"), "(text, 10, 32)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, uint32(val))")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", strconvPackage.Ident("ParseUint"), "(vals[0], 10, 32)")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = uint32(val)")
		}
	case protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
		if field.Desc.IsList() {
			g.P("list := make([]int64, 0, len(vals))")
			g.P("for _, text := range vals {")
			g.P("	val, err := ", strconvPackage.Ident("ParseInt"), "(text, 10, 64)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, val)")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", strconvPackage.Ident("ParseInt"), "(vals[0], 10, 64)")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = val")
		}
	case protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
		if field.Desc.IsList() {
			g.P("list := make([]uint64, 0, len(vals))")
			g.P("for _, text := range vals {")
			g.P("	val, err := ", strconvPackage.Ident("ParseUint"), "(text, 10, 64)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, val)")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", strconvPackage.Ident("ParseUint"), "(vals[0], 10, 64)")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = val")
		}
	case protoreflect.FloatKind:
		if field.Desc.IsList() {
			g.P("list := make([]float32, 0, len(vals))")
			g.P("for _, text := range vals {")
			g.P("	val, err := ", strconvPackage.Ident("ParseFloat"), "(text, 32)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, float32(val))")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", strconvPackage.Ident("ParseFloat"), "(vals[0], 32)")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = float32(val)")
		}
	case protoreflect.DoubleKind:
		if field.Desc.IsList() {
			g.P("list := make([]float64, 0, len(vals))")
			g.P("for _, text := range vals {")
			g.P("	val, err := ", strconvPackage.Ident("ParseFloat"), "(text, 64)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, val)")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", strconvPackage.Ident("ParseFloat"), "(vals[0], 64)")
			g.P("if err != nil {")
			g.P("	return nil, err")
			g.P("}")
			g.P(path, " = val")
		}
	case protoreflect.StringKind:
		if field.Desc.IsList() {
			g.P(path, " = vals")
		} else {
			g.P(path, " = vals[0]")
		}
	case protoreflect.BytesKind:
		if field.Desc.IsList() {
			g.P("list := make([][]byte, 0, len(vals))")
			g.P("for _, text := range vals {")
			g.P("	val, err := ", base64Package.Ident("StdEncoding.DecodeString"), "(text)")
			g.P("	if err != nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P("	list = append(list, val)")
			g.P("}")
			g.P(path, " = list")
		} else {
			g.P("val, err := ", base64Package.Ident("StdEncoding.DecodeString"), "(vals[0])")
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
