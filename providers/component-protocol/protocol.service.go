package component_protocol

import (
	"context"
	"encoding/json"

	"github.com/erda-project/erda-infra/providers/component-protocol/definition"
	"github.com/erda-project/erda-infra/providers/component-protocol/definition/cptype"
	commonpb "github.com/erda-project/erda-proto-go/common/pb"
	"github.com/erda-project/erda-proto-go/cp/pb"
)

type protocolService struct {
	p      *provider
}

func (s *protocolService) Render(ctx context.Context, req *pb.RenderRequest) (*pb.RenderResponse, error) {
	s.p.Log.Debugf("request header: %v", GetAllHeaders(ctx))

	// transfer pb to ComponentProtocolRequest for easy use
	renderReq := &cptype.ComponentProtocolRequest{}
	if err := objTransfer(req, renderReq); err != nil {
		return nil, err
	}

	// TODO set bundle in component-protocol-erda
	// TODO set i18n
	// TODO get locale
	identity := GetIdentity(ctx)

	ctxBdl := definition.ContextBundle{
		Tran:     s.p.Tran,
		Identity: identity,
		InParams: renderReq.InParams,
	}
	ctx = context.WithValue(ctx, definition.GlobalInnerKeyCtxBundle.String(), ctxBdl)

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

func GetIdentity(ctx context.Context) *commonpb.IdentityInfo {
	return &commonpb.IdentityInfo{
		UserID: GetHeader(ctx, "User-ID"),
		OrgID:  GetHeader(ctx, "Org-ID"),
	}
}

// objTransfer transfer from src to dst using json.
func objTransfer(src interface{}, dst interface{}) error {
	b, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dst)
}
