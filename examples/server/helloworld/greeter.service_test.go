package example

import (
	context "context"
	reflect "reflect"
	testing "testing"

	pb "github.com/erda-project/erda-infra/examples/protocol/pb"
)

func Test_greeterService_SayHello(t *testing.T) {
	type fields struct {
		p *provider
	}
	type args struct {
		ctx context.Context
		req *pb.HelloRequest
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantResp *pb.HelloResponse
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &greeterService{
				p: tt.fields.p,
			}
			gotResp, err := s.SayHello(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("greeterService.SayHello() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResp, tt.wantResp) {
				t.Errorf("greeterService.SayHello() = %v, want %v", gotResp, tt.wantResp)
			}
		})
	}
}
