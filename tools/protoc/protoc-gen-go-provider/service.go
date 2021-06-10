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
	"strconv"
	"strings"
	"unicode"

	"google.golang.org/protobuf/compiler/protogen"
)

const (
	contextPackage = protogen.GoImportPath("context")
	statusPackage  = protogen.GoImportPath("google.golang.org/grpc/status")
	codesPackage   = protogen.GoImportPath("google.golang.org/grpc/codes")
	testingPackage = protogen.GoImportPath("testing")
	reflectPackage = protogen.GoImportPath("reflect")
)

func genServices(gen *protogen.Plugin, files []*protogen.File, root *protogen.File) error {
	for _, file := range files {
		for _, ser := range file.Services {
			// service file
			filename := strings.ToLower(strings.Join(splitCase(ser.GoName), ".")) + ".go"
			g := gen.NewGeneratedFile(filename, protogen.GoImportPath(root.Desc.Package().Name()))
			g.P("package ", root.Desc.Package().Name())
			g.P()
			typeName := lowerCaptain(ser.GoName)
			g.P("type ", typeName, " struct {")
			g.P("	p *provider")
			g.P("}")
			g.P()
			for _, m := range ser.Methods {
				g.P("func (s *", typeName, ") ", m.GoName, "(ctx ", contextPackage.Ident("Context"), ",req *", m.Input.GoIdent, ") (*", m.Output.GoIdent, ", error) {")
				g.P("	// TODO .")
				g.P("	return nil, ", statusPackage.Ident("Errorf"), "(", codesPackage.Ident("Unimplemented"), ", \"method ", m.GoName, " not implemented\")")
				g.P("}")
			}

			// test file
			filename = strings.ToLower(strings.Join(splitCase(ser.GoName), ".")) + "_test.go"
			g = gen.NewGeneratedFile(filename, protogen.GoImportPath(root.Desc.Package().Name()))
			g.P("package ", root.Desc.Package().Name())
			g.P()
			service := strings.TrimRight(string(file.Desc.Package()), ".") + "." + ser.GoName
			for _, m := range ser.Methods {
				g.P("func Test_", typeName, "_", m.GoName, "(t *", testingPackage.Ident("T"), ") {")
				g.P("	type args struct {")
				g.P("		ctx ", contextPackage.Ident("Context"))
				g.P("		req *", m.Input.GoIdent)
				g.P("	}")
				g.P("	tests := []struct {")
				g.P("		name     string")
				g.P("		service  string")
				g.P("		config   string")
				g.P("		args     args")
				g.P("		wantResp *", m.Output.GoIdent)
				g.P("		wantErr  bool")
				g.P("	}{")
				g.P("		// TODO: Add test cases.")
				g.P("		{")
				g.P(`			"case 1",`)
				g.P(`			`, strconv.Quote(service), `,`)
				g.P("			`")
				g.P(string(root.Desc.Package()), ":")
				g.P("`,")
				g.P(`			args{`)
				g.P(`				`, contextPackage.Ident("TODO"), `(),`)
				g.P(`				&`, m.Input.GoIdent, `{`)
				g.P(`					// TODO: setup fields`)
				g.P(`				},`)
				g.P(`			},`)
				g.P(`			&`, m.Output.GoIdent, `{`)
				g.P(`				// TODO: setup fields.`)
				g.P(`			},`)
				g.P(`			false,`)
				g.P("		},")
				g.P("	}")
				g.P("	for _, tt := range tests {")
				g.P("		t.Run(tt.name, func(t *testing.T) {")
				g.P("			hub := ", servicehubPackage.Ident("New"), "()")
				g.P("			events := hub.Events()")
				g.P("			go func() {")
				g.P("				hub.RunWithOptions(&", servicehubPackage.Ident("RunOptions"), "{Content: tt.config})")
				g.P("			}()")
				g.P("			err := <-events.Started()")
				g.P("			if err != nil {")
				g.P("				t.Error(err)")
				g.P("				return")
				g.P("			}")
				if *genGRPC {
					g.P("			srv := hub.Service(tt.service).(", root.GoImportPath.Ident(ser.GoName+"Server"), ")")
				} else if *genHTTP {
					g.P("			srv := hub.Service(tt.service).(", root.GoImportPath.Ident(ser.GoName+"Handler"), ")")
				} else {
					g.P("			srv := hub.Service(tt.service).(*", lowerCaptain(ser.GoName), ")")
				}
				g.P("			got, err := srv.", m.GoName, "(tt.args.ctx, tt.args.req)")
				g.P("			if (err != nil) != tt.wantErr {")
				g.P("				t.Errorf(\"", typeName, ".", m.GoName, "() error = %v, wantErr %v\", err, tt.wantErr)")
				g.P("				return")
				g.P("			}")
				g.P("			if !", reflectPackage.Ident("DeepEqual"), "(got, tt.wantResp) {")
				g.P("				t.Errorf(\"", typeName, ".", m.GoName, "() = %v, want %v\", got, tt.wantResp)")
				g.P("			}")
				g.P("		})")
				g.P("	}")
				g.P("}")
				g.P()
			}
		}
	}
	return nil
}

func splitCase(name string) (list []string) {
	if len(name) <= 0 {
		return nil
	}
	chars := []rune(name)
	pre, start, idx, num := chars[0], 0, 1, len(chars)
	for ; idx < num; idx++ {
		if unicode.IsUpper(chars[idx]) != unicode.IsUpper(pre) {
			pre = chars[idx]
			if idx-start == 1 {
				continue
			}
			if unicode.IsLower(chars[idx]) {
				list = append(list, string(chars[start:idx-1]))
				start = idx - 1
				continue
			}
			list = append(list, string(chars[start:idx]))
			start = idx
		}
	}
	if start < num {
		return append(list, string(chars[start:]))
	}
	return list
}
