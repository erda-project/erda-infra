package example

import (
	context "context"

	pb "github.com/erda-project/erda-infra/examples/protocol/pb"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type userService struct {
	p *provider
}

func (s *userService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	// TODO .
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (s *userService) UpdateUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UpdateUserResponse, error) {
	// TODO .
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
