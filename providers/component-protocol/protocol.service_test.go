package component_protocol

import (
	"context"
	"reflect"
	"testing"

	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-proto-go/cp/pb"
)

func Test_protocolService_Render(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.RenderRequest
	}
	tests := []struct {
		name     string
		service  string
		config   string
		args     args
		wantResp *pb.RenderResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
		{
			"case 1",
			"erda.component_protocol.ProtocolService",
			`
erda.component_protocol:
`,
			args{
				context.TODO(),
				&pb.RenderRequest{
					// TODO: setup fields
				},
			},
			&pb.RenderResponse{
				// TODO: setup fields.
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hub := servicehub.New()
			events := hub.Events()
			go func() {
				hub.RunWithOptions(&servicehub.RunOptions{Content: tt.config})
			}()
			err := <-events.Started()
			if err != nil {
				t.Error(err)
				return
			}
			srv := hub.Service(tt.service).(pb.CPServiceServer)
			got, err := srv.Render(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("protocolService.Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantResp) {
				t.Errorf("protocolService.Render() = %v, want %v", got, tt.wantResp)
			}
		})
	}
}
