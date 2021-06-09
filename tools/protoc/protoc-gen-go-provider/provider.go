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
	"sort"
	"strconv"
	"strings"
	"unicode"

	"google.golang.org/protobuf/compiler/protogen"
)

const (
	logPackage        = protogen.GoImportPath("github.com/erda-project/erda-infra/base/logs")
	servicehubPackage = protogen.GoImportPath("github.com/erda-project/erda-infra/base/servicehub")
	transportPackage  = protogen.GoImportPath("github.com/erda-project/erda-infra/pkg/transport")
	transhttpPackage  = protogen.GoImportPath("github.com/erda-project/erda-infra/pkg/transport/http")
	transgrpcPackage  = protogen.GoImportPath("github.com/erda-project/erda-infra/pkg/transport/grpc")
)

func generateFiles(gen *protogen.Plugin, files []*protogen.File) error {
	sort.Slice(files, func(i, j int) bool {
		return files[i].Desc.Name() < files[j].Desc.Name()
	})
	var file *protogen.File
	var count int
	for _, f := range files {
		if len(f.Services) <= 0 {
			continue
		}
		count += len(f.Services)
		if file == nil {
			file = f
		}
		if f.GoImportPath != file.GoImportPath {
			return fmt.Errorf("package path conflict between %s and %s", file.GoImportPath, f.GoImportPath)
		}
		if f.Desc.Package() != file.Desc.Package() {
			return fmt.Errorf("package path conflict between %s and %s", file.Desc.Package(), f.Desc.Package())
		}
	}
	if count <= 0 {
		return fmt.Errorf("not found service in all proto files")
	}
	err := genProvider(gen, files, file)
	if err != nil {
		return err
	}
	return genServices(gen, files, file)
}

func genProvider(gen *protogen.Plugin, files []*protogen.File, root *protogen.File) error {
	const filename = "provider.go"
	g := gen.NewGeneratedFile(filename, protogen.GoImportPath(root.Desc.Package().Name()))
	g.P("package ", root.Desc.Package().Name())
	g.P()
	g.P("type config struct {")
	g.P("}")
	g.P()
	g.P("// +provider")
	g.P("type provider struct {")
	g.P("	Cfg    *config")
	g.P("	Log    ", logPackage.Ident("Logger"))
	if *genHTTP && *genGRPC {
		g.P("	Register   ", transportPackage.Ident("Register"))
	} else if *genHTTP {
		g.P("	Register   ", transhttpPackage.Ident("Router"))
	} else if *genGRPC {
		g.P("	Register   ", transgrpcPackage.Ident("ServiceRegistrar"))
	}
	for _, file := range files {
		for _, ser := range file.Services {
			g.P(lowerCaptain(ser.GoName), " *", lowerCaptain(ser.GoName))
		}
	}
	g.P("}")
	g.P()
	g.P("func (p *provider) Init(ctx ", servicehubPackage.Ident("Context"), ") error {")
	g.P("	// TODO initialize something ...")
	g.P()
	for _, file := range files {
		for _, ser := range file.Services {
			g.P("p.", lowerCaptain(ser.GoName), " = &", lowerCaptain(ser.GoName), "{p}")
		}
	}
	g.P("	if p.Register != nil {")
	for i, file := range files {
		for _, ser := range file.Services {
			if *genHTTP && *genGRPC {
				g.P(root.GoImportPath.Ident("Register"+ser.GoName+"Imp"), "(p.Register, p.", lowerCaptain(ser.GoName), ")")
			} else if *genGRPC {
				g.P(root.GoImportPath.Ident("Register"+ser.GoName+"Server"), "(p.Register, p.", lowerCaptain(ser.GoName), ")")
			} else if *genHTTP {
				g.P(root.GoImportPath.Ident("Register"+ser.GoName+"Handler"), "(p.Register, p.", lowerCaptain(ser.GoName), ")")
			}
			if i < len(files)-1 {
				g.P()
			}
		}
	}
	g.P("	}")
	g.P("	return nil")
	g.P("}")
	g.P()
	g.P("func (p *provider) Provide(ctx ", servicehubPackage.Ident("DependencyContext"), ", args ...interface{}) interface{} {")
	g.P("	switch {")
	for _, file := range files {
		for _, ser := range file.Services {
			service := strings.TrimRight(string(file.Desc.Package()), ".") + "." + ser.GoName
			if *genHTTP && *genGRPC {
				g.P("	case ctx.Service() == ", strconv.Quote(service),
					" || ctx.Type() == ", root.GoImportPath.Ident(ser.GoName+"ServerType"), "() || ctx.Type() == ",
					root.GoImportPath.Ident(ser.GoName+"HandlerType"), "():")
			} else if *genGRPC {
				g.P("	case ctx.Service() == ", strconv.Quote(service),
					" || ctx.Type() == ", root.GoImportPath.Ident(ser.GoName+"ServerType"), "():")
			} else if *genHTTP {
				g.P("	case ctx.Service() == ", strconv.Quote(service),
					" || ctx.Type() == ", root.GoImportPath.Ident(ser.GoName+"HandlerType"), "():")
			} else {
				g.P("	case ctx.Service() == ", strconv.Quote(service), ":")
			}
			g.P("		return p.", lowerCaptain(ser.GoName))
		}
	}
	g.P("	}")
	g.P("	return p")
	g.P("}")
	g.P()
	g.P("func init() {")
	g.P("	", servicehubPackage.Ident("Register"), "(", strconv.Quote(string(root.Desc.Package())), ", &", servicehubPackage.Ident("Spec"), "{")
	g.P("		Services: ", root.GoImportPath.Ident("ServiceNames"), "(),")
	g.P("		Types: ", root.GoImportPath.Ident("Types"), "(),")
	g.P("		OptionalDependencies: []string{\"service-register\"},")
	g.P("		Description: \"\",")
	g.P("		ConfigFunc: func() interface{} {")
	g.P("			return &config{}")
	g.P("		},")
	g.P("		Creator: func() ", servicehubPackage.Ident("Provider"), " {")
	g.P("			return &provider{}")
	g.P("		},")
	g.P("	})")
	g.P("}")
	return nil
}

func lowerCaptain(name string) string {
	if len(name) <= 0 {
		return name
	}
	chars := []rune(name)
	pre := chars[0]
	if unicode.IsLower(pre) {
		return name
	}
	for i, c := range chars {
		if unicode.IsUpper(c) != unicode.IsUpper(pre) {
			break
		}
		chars[i] = unicode.ToLower(c)
	}
	return string(chars)
}
