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

package list

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

func TestItemExtra(t *testing.T) {
	item := Item{
		ID: "ID",
		Extra: cptype.Extra{
			Extra: map[string]interface{}{
				"extKey1": "extVal1",
				"extKey2": "extVal2",
			},
		},
	}
	b, _ := json.Marshal(item)
	// output: {"id":"ID","extra":{"extKey1":"extVal1","extKey2":"extVal2"}}"
	logrus.Infof("item bytes: %s", string(b))
}

func TestModel(t *testing.T) {
	data := Data{Total: 0}
	compData := cptype.ComponentData{"total": 45}
	cputil.MustObjJSONTransfer(&data, &compData)
	fmt.Printf("%#v", compData)
	_ = data
}
