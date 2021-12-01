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

package etcdclientv3

import (
	_ "unsafe"

	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"

	grpccontext "github.com/erda-project/erda-infra/pkg/trace/inject/context/grpc"
	"github.com/erda-project/erda-infra/pkg/trace/inject/hook"
)

//go:linkname newClient github.com/coreos/etcd/clientv3.newClient
//go:noinline
func newClient(cfg *clientv3.Config) (*clientv3.Client, error)

//go:noinline
func originalNewClient(cfg *clientv3.Config) (*clientv3.Client, error) {
	return newClient(cfg)
}

//go:noinline
func wrappedNewClient(cfg *clientv3.Config) (*clientv3.Client, error) {
	cfg.DialOptions = append([]grpc.DialOption{
		grpc.WithUnaryInterceptor(grpccontext.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(grpccontext.StreamClientInterceptor()),
	}, cfg.DialOptions...)
	return originalNewClient(cfg)
}

func init() {
	hook.Hook(newClient, wrappedNewClient, originalNewClient)
}
