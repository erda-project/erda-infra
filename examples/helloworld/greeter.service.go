package example

import (
	context "context"

	pb "github.com/erda-project/erda-infra/examples/protocol/pb"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type greeterService struct {
	p *provider
}

func (s *greeterService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	// TODO .
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
