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

package componentprotocol

import (
	"bytes"
	"errors"
	"net/http"
	"strconv"
	"testing"

	"google.golang.org/protobuf/types/known/structpb"

	"github.com/erda-project/erda-infra/providers/component-protocol/protobuf/proto-go/cp/pb"
)

type MockResponseWriter struct {
	Head http.Header
	Data bytes.Buffer
}

func (m *MockResponseWriter) Header() http.Header {
	return m.Head
}

func (m *MockResponseWriter) Write(data []byte) (int, error) {
	return m.Data.Write(data)
}

func (m *MockResponseWriter) WriteHeader(statusCode int) {
	m.Head["test"] = []string{
		strconv.FormatInt(int64(statusCode), 10),
	}
	return
}

type mockStr struct {
}

func TestEncoder(t *testing.T) {
	rw := &MockResponseWriter{
		Head: map[string][]string{},
		Data: bytes.Buffer{},
	}
	err := encoder(rw, nil, nil)
	if err != nil {
		t.Error(err)
	}

	err = encoder(rw, nil, mockStr{})
	if err != nil {
		t.Error(err)
	}

	listValue, _ := structpb.NewList([]interface{}{"testUserID"})
	resp := &pb.RenderResponse{
		Scenario: &pb.Scenario{
			ScenarioKey:  "testScenario",
			ScenarioType: "testScenario",
		},
		Protocol: &pb.ComponentProtocol{
			Version:  "0.1",
			Scenario: "testScenario",
			GlobalState: map[string]*structpb.Value{
				"_userIDs_": structpb.NewListValue(listValue),
			},
			Hierarchy: &pb.Hierarchy{
				Root: "page",
			},
			Components: map[string]*pb.Component{
				"testComponent": {
					Type: "container",
					Name: "test",
				},
			},
		},
	}
	if err = encoder(rw, nil, resp); err != nil {
		t.Error(err)
	}
}

type MockError struct {
}

func (e *MockError) HTTPStatus() int {
	return 500
}

func (e *MockError) Error() string {
	return ""
}

func TestErrorEncoder(t *testing.T) {
	rw := &MockResponseWriter{
		Head: map[string][]string{},
		Data: bytes.Buffer{},
	}
	errorEncoder(rw, nil, &MockError{})
	if rw.Head["test"][0] != "500" {
		t.Errorf("test failed, expect head is 500, got %s", rw.Head["test"][0])
	}
	errorEncoder(rw, nil, errors.New("test"))
}
