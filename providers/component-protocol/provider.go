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
	"net/http"
	"reflect"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/pkg/transport"
	transhttp "github.com/erda-project/erda-infra/pkg/transport/http"
	"github.com/erda-project/erda-infra/providers/i18n"
	"github.com/erda-project/erda-proto-go/cp/pb"
	jsi "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

type config struct {
}

// +provider
type provider struct {
	Cfg      *config
	Log      logs.Logger
	Register transport.Register

	tran             i18n.Translator
	customContextKVs map[interface{}]interface{}

	protocolService *protocolService
	// internalTran    i18n.Translator `translator:"18n-cp-internal"`
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.customContextKVs = make(map[interface{}]interface{})
	p.protocolService = &protocolService{p: p}
	if p.Register != nil {
		pb.RegisterCPServiceImp(p.Register, p.protocolService, transport.WithHTTPOptions(
			transhttp.WithEncoder(encoder),
			transhttp.WithErrorEncoder(errorEncoder),
		))
	}
	return nil
}

// Error .
type Error interface {
	HTTPStatus() int
}

func encoder(rw http.ResponseWriter, r *http.Request, obj interface{}) error {
	if obj == nil {
		return nil
	}
	renderResp, ok := obj.(*pb.RenderResponse)
	if !ok {
		logrus.Errorf("response obj is not *pb.RenderResponse type")
		return nil
	}
	if renderResp.Protocol != nil && renderResp.Protocol.GlobalState != nil &&
		renderResp.Protocol.GlobalState["_userIDs_"] != nil {
		rw.Header().Set("X-Need-User-Info", "true")
	}
	resp := map[string]interface{}{
		"success": true,
		"data":    renderResp,
		"err": map[string]interface{}{
			"code": "",
			"msg":  "",
			"ctx":  nil,
		},
	}
	data, err := jsi.Marshal(resp)
	if err != nil {
		return err
	}
	if _, err = rw.Write(data); err != nil {
		return err
	}
	rw.Header().Set("Content-Type", "application/json")
	return nil
}

func errorEncoder(rw http.ResponseWriter, request *http.Request, err error) {
	var status int
	if e, ok := err.(Error); ok {
		status = e.HTTPStatus()
	} else {
		status = http.StatusInternalServerError
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)
	byts, _ := jsi.Marshal(map[string]interface{}{
		"success": false,
		"err": map[string]interface{}{
			"code": status,
			"msg":  err.Error(),
		},
	})
	rw.Write(byts)
}

// Provide .
func (p *provider) Provide(ctx servicehub.DependencyContext, args ...interface{}) interface{} {
	return p
}

func init() {
	interfaceType := reflect.TypeOf((*Interface)(nil)).Elem()
	servicehub.Register("component-protocol", &servicehub.Spec{
		Services:             pb.ServiceNames(),
		Types:                append(pb.Types(), interfaceType),
		OptionalDependencies: []string{"service-register"},
		Description:          "",
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}
