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

package componentprotocol

import (
	"context"

	"github.com/erda-project/erda-infra/providers/component-protocol/protobuf/proto-go/cp/pb"
	"github.com/erda-project/erda-infra/providers/i18n"
)

// Interface wrap CPService and other logic.
type Interface interface {
	Render(ctx context.Context, req *pb.RenderRequest) (*pb.RenderResponse, error)
	SetI18nTran(tran i18n.I18n)
	WithContextValue(key, value interface{})
}

// Render .
func (p *provider) Render(ctx context.Context, req *pb.RenderRequest) (*pb.RenderResponse, error) {
	return p.protocolService.Render(ctx, req)
}

// SetI18nTran .
func (p *provider) SetI18nTran(tran i18n.I18n) { p.tran = tran }

// WithContextValue .
func (p *provider) WithContextValue(key, value interface{}) { p.customContextKVs[key] = value }
