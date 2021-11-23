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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/protobuf/proto-go/cp/pb"
)

func TestObjJSONTransfer(t *testing.T) {
	src := pb.ComponentProtocol{
		Options: &pb.ProtocolOptions{
			SyncIntervalSecond: 0.01,
		},
	}
	var dest cptype.ComponentProtocol

	err := ObjJSONTransfer(&src, &dest)
	assert.NoError(t, err)
	assert.Equal(t, src.Options.SyncIntervalSecond, dest.Options.SyncIntervalSecond)
}
