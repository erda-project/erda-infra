// Copyright 2021 Terminus
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
	return nil, status.Errorf(codes.Unimplemented, "method GetUser not implemented")
}
func (s *userService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	// TODO .
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
