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
	"context"
	"encoding/json"

	"github.com/erda-project/erda-infra/providers/component-protocol/definition"
	"github.com/erda-project/erda-infra/providers/component-protocol/definition/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
	"github.com/erda-project/erda-proto-go/cp/pb"
)

type protocolService struct {
	p *provider
}

// Render .
func (s *protocolService) Render(ctx context.Context, req *pb.RenderRequest) (*pb.RenderResponse, error) {
	s.p.Log.Debugf("request header: %v", cputil.GetAllHeaders(ctx))

	// transfer pb to ComponentProtocolRequest for easy use
	renderReq := &cptype.ComponentProtocolRequest{}
	if err := objTransfer(req, renderReq); err != nil {
		return nil, err
	}

	ctxBdl := definition.SDK{
		Tran:     s.p.Tran,
		Identity: cputil.GetIdentity(ctx),
		InParams: renderReq.InParams,
		Lang:     cputil.Language(ctx),
	}
	ctx = context.WithValue(ctx, definition.GlobalInnerKeyCtxSDK, &ctxBdl)
	for k, v := range s.p.CustomContextKVs {
		ctx = context.WithValue(ctx, k, v)
	}

	// render concrete scenario
	if err := definition.RunScenarioRender(ctx, renderReq); err != nil {
		return nil, err
	}

	// make response
	resp, err := s.makeResponse(renderReq)

	return resp, err
}

func (s *protocolService) makeResponse(renderReq *cptype.ComponentProtocolRequest) (*pb.RenderResponse, error) {
	// check if business error exist, not platform error
	businessErr := definition.GetGlobalStateKV(renderReq.Protocol, definition.GlobalInnerKeyError.String())
	if businessErr != nil {
		if err, ok := businessErr.(error); ok {
			return nil, err
		}
		s.p.Log.Warnf("business error type is not error, %#v", businessErr)
	}

	// render response
	polishedPbReq := &pb.RenderRequest{}
	if err := objTransfer(renderReq, polishedPbReq); err != nil {
		return nil, err
	}
	pbResp := pb.RenderResponse{
		Scenario: polishedPbReq.Scenario,
		Protocol: polishedPbReq.Protocol,
	}

	return &pbResp, nil
}

// objTransfer transfer from src to dst using json.
func objTransfer(src interface{}, dst interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}
