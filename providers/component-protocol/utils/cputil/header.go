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

package cputil

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/erda-project/erda-infra/pkg/transport"
	commonpb "github.com/erda-project/erda-proto-go/common/pb"
)

// GetOrgID .
func GetOrgID(ctx context.Context) string {
	return GetHeader(ctx, "org-id")
}

// GetUserID .
func GetUserID(ctx context.Context) string {
	return GetHeader(ctx, "user-id")
}

// GetHeader .
func GetHeader(ctx context.Context, key string) string {
	header := transport.ContextHeader(ctx)
	if header != nil {
		for _, v := range header.Get(key) {
			if len(v) > 0 {
				return v
			}
		}
	}
	return ""
}

func GetAllHeaders(ctx context.Context) metadata.MD {
	return transport.ContextHeader(ctx)
}

func GetIdentity(ctx context.Context) *commonpb.IdentityInfo {
	return &commonpb.IdentityInfo{
		UserID: GetUserID(ctx),
		OrgID:  GetOrgID(ctx),
	}
}
