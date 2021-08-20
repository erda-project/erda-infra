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

package component_protocol

import (
	"context"

	"github.com/erda-project/erda-infra/providers/i18n"
	"github.com/erda-project/erda-proto-go/cp/pb"
)

type Interface interface {
	Render(ctx context.Context, req *pb.RenderRequest) (*pb.RenderResponse, error)
	SetI18nTran(tran i18n.Translator)
	WithContextValue(key, value interface{})
}

func (p *provider) Render(ctx context.Context, req *pb.RenderRequest) (*pb.RenderResponse, error) {
	return p.protocolService.Render(ctx, req)
}
func (p *provider) SetI18nTran(tran i18n.Translator)        { p.Tran = tran }
func (p *provider) WithContextValue(key, value interface{}) { p.CustomContextKVs[key] = value }
