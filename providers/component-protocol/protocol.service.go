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

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/protobuf/proto-go/cp/pb"
	"github.com/erda-project/erda-infra/providers/component-protocol/protocol"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

type protocolService struct {
	p *provider
}

// Render .
func (s *protocolService) Render(ctx context.Context, req *pb.RenderRequest) (*pb.RenderResponse, error) {
	s.p.Log.Debugf("request header: %v", cputil.GetAllHeaders(ctx))

	// transfer pb to ComponentProtocolRequest for easy use
	renderReq := &cptype.ComponentProtocolRequest{}
	if err := cputil.ObjJSONTransfer(req, renderReq); err != nil {
		return nil, err
	}

	// make sdk
	sdk := cptype.SDK{
		Scenario: req.Scenario.ScenarioKey,
		Tran:     s.p.tran.Translator(req.Scenario.ScenarioKey),
		Identity: cputil.GetIdentity(ctx),
		InParams: renderReq.InParams,
		Lang:     cputil.Language(ctx),
	}

	// make ctx with sdk
	ctx = context.WithValue(ctx, cptype.GlobalInnerKeyCtxSDK, &sdk)
	for k, v := range s.p.customContextKVs {
		ctx = context.WithValue(ctx, k, v)
	}

	// temp state
	ctx = context.WithValue(ctx, cptype.GlobalInnerKeyStateTemp, make(map[string]interface{}))

	// render concrete scenario
	if err := protocol.RunScenarioRender(ctx, renderReq); err != nil {
		return nil, err
	}

	// make response
	resp, err := s.makeResponse(renderReq)

	return resp, err
}

func (s *protocolService) makeResponse(renderReq *cptype.ComponentProtocolRequest) (*pb.RenderResponse, error) {
	// check if business error exist, not platform error
	businessErr := protocol.GetGlobalStateKV(renderReq.Protocol, cptype.GlobalInnerKeyError.String())
	if businessErr != nil {
		if err, ok := businessErr.(error); ok {
			return nil, err
		}
		s.p.Log.Warnf("business error type is not error, %#v", businessErr)
	}

	// render response
	polishedPbReq := &pb.RenderRequest{}
	if err := cputil.ObjJSONTransfer(renderReq, polishedPbReq); err != nil {
		return nil, err
	}
	pbResp := pb.RenderResponse{
		Scenario: polishedPbReq.Scenario,
		Protocol: polishedPbReq.Protocol,
	}

	return &pbResp, nil
}
