// Code generated by protoc-gen-go-client. DO NOT EDIT.
// Sources: protocol.proto

package client

import (
	context "context"

	grpc "github.com/erda-project/erda-infra/pkg/transport/grpc"
	pb "github.com/erda-project/erda-infra/providers/component-protocol/protobuf/proto-go/cp/pb"
	grpc1 "google.golang.org/grpc"
)

// Client provide all service clients.
type Client interface {
	// CPService protocol.proto
	CPService() pb.CPServiceClient
}

// New create client
func New(cc grpc.ClientConnInterface) Client {
	return &serviceClients{
		cpservice: pb.NewCPServiceClient(cc),
	}
}

type serviceClients struct {
	cpservice pb.CPServiceClient
}

func (c *serviceClients) CPService() pb.CPServiceClient {
	return c.cpservice
}

type cpserviceWrapper struct {
	client pb.CPServiceClient
	opts   []grpc1.CallOption
}

func (s *cpserviceWrapper) Render(ctx context.Context, req *pb.RenderRequest) (*pb.RenderResponse, error) {
	return s.client.Render(ctx, req, append(grpc.CallOptionFromContext(ctx), s.opts...)...)
}
