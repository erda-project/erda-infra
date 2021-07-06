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

package jsonpb

import (
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/runtime/protoimpl"
)

// AnyResolver takes a type URL, present in an Any message,
// and resolves it into an instance of the associated message.
type AnyResolver interface {
	Resolve(typeURL string) (proto.Message, error)
}

type anyResolver struct{ AnyResolver }

func (r anyResolver) FindMessageByName(message protoreflect.FullName) (protoreflect.MessageType, error) {
	return r.FindMessageByURL(string(message))
}

func (r anyResolver) FindMessageByURL(url string) (protoreflect.MessageType, error) {
	m, err := r.Resolve(url)
	if err != nil {
		return nil, err
	}
	return protoimpl.X.MessageTypeOf(m), nil
}

func (r anyResolver) FindExtensionByName(field protoreflect.FullName) (protoreflect.ExtensionType, error) {
	return protoregistry.GlobalTypes.FindExtensionByName(field)
}

func (r anyResolver) FindExtensionByNumber(message protoreflect.FullName, field protoreflect.FieldNumber) (protoreflect.ExtensionType, error) {
	return protoregistry.GlobalTypes.FindExtensionByNumber(message, field)
}

func wellKnownType(s protoreflect.FullName) string {
	if s.Parent() == "google.protobuf" {
		switch s.Name() {
		case "Empty", "Any",
			"BoolValue", "BytesValue", "StringValue",
			"Int32Value", "UInt32Value", "FloatValue",
			"Int64Value", "UInt64Value", "DoubleValue",
			"Duration", "Timestamp",
			"NullValue", "Struct", "Value", "ListValue":
			return string(s.Name())
		}
	}
	return ""
}

func isMessageSet(md protoreflect.MessageDescriptor) bool {
	ms, ok := md.(interface{ IsMessageSet() bool })
	return ok && ms.IsMessageSet()
}
