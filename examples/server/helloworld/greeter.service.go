package example

import (
	"context"

	"github.com/erda-project/erda-infra/examples/protocol/pb"
)

type greeterService struct {
	p *provider
}

func (s *greeterService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Success: true,
		Data:    "hello " + req.Name,
	}, nil
}
