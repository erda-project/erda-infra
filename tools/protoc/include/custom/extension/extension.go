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

package extension

import (
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
)

func MethodShouldBeGrpcSkipped(method *protogen.Method) bool {
	pureHTTP := proto.GetExtension(method.Desc.Options(), E_Http).(*HttpMethodOption).GetPure()
	return pureHTTP
}

func GetServiceGrpcMethods(service *protogen.Service) (grpcMethods []*protogen.Method) {
	return GetGrpcMethods(service.Methods)
}

func GetGrpcMethods(methods []*protogen.Method) (grpcMethods []*protogen.Method) {
	for _, method := range methods {
		if !MethodShouldBeGrpcSkipped(method) {
			grpcMethods = append(grpcMethods, method)
		}
	}
	return
}
