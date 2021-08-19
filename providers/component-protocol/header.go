// Copyright (c) 2021 Terminus, Inc.
//
// This program is free software: you can use, redistribute, and/or modify
// it under the terms of the GNU Affero General Public License, version 3
// or later ("AGPL"), as published by the Free Software Foundation.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package component_protocol

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/erda-project/erda-infra/pkg/transport"
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
