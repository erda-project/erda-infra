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
	"reflect"
	"testing"

	"github.com/erda-project/erda-infra/pkg/transport/http"
	"google.golang.org/grpc"
)

type MockRouter struct {

}

func (m *MockRouter) Add(method, path string, handler http.HandlerFunc) {
}

func (m *MockRouter) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
}

func TestInit(t *testing.T) {
	p := provider{Register: &MockRouter{}}
	if err := p.Init(nil); err != nil {
		t.Error(err)
	}
}

func TestConvertErrToResp(t *testing.T) {


	obj := cpErrResponse{
		Code: 500,
		Err:  "Failed",
	}
	resp, err := convertErrToResp(obj)
	if err != nil {
		t.Error(err)
	}
	success, ok := resp["success"].(bool)
	if !ok {
		t.Errorf("test failed, expected type of success field is bool, got %v", reflect.TypeOf(resp["success"].(bool)))
	}
	if success {
		t.Error("test failed, expect values of success field is false, got true")
	}

	respErr, ok := resp["err"].(map[string]interface{})
	if !ok {
		t.Errorf("test failed, expected type of err field is map[string]interface{}, got %v", reflect.TypeOf(resp["err"]))
	}

	code, ok := respErr["code"].(string)
	if !ok {
		t.Errorf("test failed, expected type of code field is string, got %v", reflect.TypeOf(respErr["code"]))
	}
	expectedCode := "Proxy Error: 500"
	if code != expectedCode {
		t.Errorf("test failed, expected values of code is %s, got %s", expectedCode, code)
	}

	msg, ok := respErr["msg"].(string)
	if !ok {
		t.Errorf("test failed, expected type of msg field is string, got %v", reflect.TypeOf(respErr["msg"]))
	}
	expectedMsg := "Failed"
	if msg != expectedMsg {
		t.Errorf("test failed, expected values of msg field is %s, got %s", expectedMsg, msg)
	}
}
