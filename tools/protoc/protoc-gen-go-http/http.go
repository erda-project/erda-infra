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
	"net/url"
	"sort"
	"strings"

	"github.com/erda-project/erda-infra/pkg/transport/http/httprule"
	"github.com/erda-project/erda-infra/pkg/transport/http/runtime"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	transhttpPackage = protogen.GoImportPath("github.com/erda-project/erda-infra/pkg/transport/http")
)

func generateFile(gen *protogen.Plugin, file *protogen.File) (*protogen.GeneratedFile, error) {
	// canGen := false
	// for _, srv := range file.Services {
	// 	for _, method := range srv.Methods {
	// 		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
	// 			continue
	// 		}
	// 		canGen = true
	// 	}
	// }
	// if !canGen {
	// 	return nil, nil
	// }
	filename := file.GeneratedFilenamePrefix + ".http.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by ", genName, ". DO NOT EDIT.")
	g.P("// Source: ", file.Desc.Path())
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	g.P("// This is a compile-time assertion to ensure that this generated file")
	g.P("// is compatible with the ", transhttpPackage, " package it is being compiled against.")
	g.P("const _ = ", transhttpPackage.Ident("SupportPackageIsVersion1"))

	for _, service := range file.Services {
		err := genService(gen, file, g, service)
		if err != nil {
			return g, err
		}
	}
	return g, nil
}

func genService(gen *protogen.Plugin, file *protogen.File, g *protogen.GeneratedFile, service *protogen.Service) error {
	g.P()
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P(deprecationComment)
	}

	sd := &serviceDesc{
		ServiceType: service.GoName,
		ServiceName: string(service.Desc.FullName()),
		Metadata:    file.Desc.Path(),
	}
	for _, method := range service.Methods {
		if method.Desc.IsStreamingClient() || method.Desc.IsStreamingServer() {
			continue
		}
		rule, ok := proto.GetExtension(method.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
		if rule != nil && ok {
			for _, bind := range rule.AdditionalBindings {
				m, err := buildHTTPRule(g, service, method, bind)
				if err != nil {
					return err
				}
				sd.Methods = append(sd.Methods, m)
			}
			m, err := buildHTTPRule(g, service, method, rule)
			if err != nil {
				return err
			}
			sd.Methods = append(sd.Methods, m)
		} else if *genAll {
			path := fmt.Sprintf("/%s/%s", service.Desc.FullName(), method.Desc.Name())
			m, err := buildMethodDesc(g, method, "POST", path)
			if err != nil {
				return err
			}
			sd.Methods = append(sd.Methods, m)
		}
	}
	return sd.execute(g)
}

func buildHTTPRule(g *protogen.GeneratedFile, service *protogen.Service, m *protogen.Method, rule *annotations.HttpRule) (*methodDesc, error) {
	var path, method string
	switch pattern := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		path = pattern.Get
		method = "GET"
	case *annotations.HttpRule_Put:
		path = pattern.Put
		method = "PUT"
	case *annotations.HttpRule_Post:
		path = pattern.Post
		method = "POST"
	case *annotations.HttpRule_Delete:
		path = pattern.Delete
		method = "DELETE"
	case *annotations.HttpRule_Patch:
		path = pattern.Patch
		method = "PATCH"
	case *annotations.HttpRule_Custom:
		path = pattern.Custom.Path
		method = pattern.Custom.Kind
	}
	if len(path) <= 0 {
		path = fmt.Sprintf("/%s/%s", service.Desc.FullName(), m.Desc.Name())
	}
	if len(method) <= 0 {
		method = "POST"
	}

	md, err := buildMethodDesc(g, m, method, path)
	if err != nil {
		return nil, err
	}
	reqbody, respBody := strings.TrimSpace(rule.Body), strings.TrimSpace(rule.ResponseBody)
	if reqbody == "*" {
		reqbody = ""
	}
	md.ReqBody = reqbody
	if respBody == "*" {
		respBody = ""
	}
	md.RespBody = respBody
	return md, nil
}

func buildMethodDesc(g *protogen.GeneratedFile, m *protogen.Method, method, path string) (*methodDesc, error) {
	queryString, idx := "", strings.Index(path, "?")
	var queryParams map[string][]string
	var queryParamKeys []string
	if idx >= 0 {
		queryString = path[idx+1:]
		if len(queryString) > 0 {
			values := make(map[string][]string)
			params, err := url.ParseQuery(queryString)
			if err != nil {
				return nil, err
			}
			for key, vals := range params {
				for _, val := range vals {
					if strings.HasPrefix(val, "{") && strings.HasSuffix(val, "}") {
						val = strings.TrimRight(strings.TrimLeft(val, "{"), "}")
						values[key] = append(values[key], val)
					}
				}
			}
			queryParams = values
			for key := range queryParams {
				queryParamKeys = append(queryParamKeys, key)
			}
			sort.Strings(queryParamKeys)
		}
		path = path[0:idx]
	}
	path = "/" + strings.TrimLeft(strings.TrimSpace(path), "/")
	var pathParams []string
	if len(path) > 1 {
		compiler, err := httprule.Parse(path)
		if err != nil {
			return nil, fmt.Errorf("invalid path %q : %s", path, err)
		}
		tp := compiler.Compile()
		pathParams = tp.Fields
		_, err = runtime.NewPattern(httprule.SupportPackageIsVersion1, tp.OpCodes, tp.Pool, tp.Verb)
		if err != nil {
			return nil, fmt.Errorf("path %q NewPattern return error: %s", path, err)
		}
	}
	return &methodDesc{
		Name:           m.GoName,
		Comment:        m.Comments.Leading.String(),
		Path:           path,
		Method:         method,
		QueryParams:    queryParams,
		QueryParamKeys: queryParamKeys,
		PathParams:     pathParams,
		Request:        g.QualifiedGoIdent(m.Input.GoIdent),
		Response:       g.QualifiedGoIdent(m.Output.GoIdent),
		Meta:           m,
	}, nil
}

const deprecationComment = "// Deprecated: Do not use."
