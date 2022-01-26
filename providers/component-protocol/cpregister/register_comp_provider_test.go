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

package cpregister

import (
	"fmt"
	"testing"

	"github.com/mohae/deepcopy"
	"github.com/stretchr/testify/assert"

	"github.com/erda-project/erda-infra/providers/component-protocol/components/topn/impl"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/i18n"
)

type MockProvider struct {
	// provider
	Tran        i18n.Translator      `i18n:"hello"`
	InParams    MockServiceInParams  `json:"inParams,omitempty"`
	InParamsPtr *MockServiceInParams `json:"inParamsPtr,omitempty"`

	// component
	impl.DefaultTop
}

func (m MockProvider) RegisterInitializeOp() (opFunc cptype.OperationFunc) {
	return func(sdk *cptype.SDK) cptype.IStdStructuredPtr { return nil }
}

type MockServiceInParams struct {
	InParams *MockModel
}

type MockModel struct {
	StartTime uint64 `json:"startTime,omitempty"`
	EndTime   uint64 `json:"endTime,omitempty"`
	TenantID  string `json:"tenantID,omitempty"`
}

func Test_copyProvider(t *testing.T) {
	originProvider := &MockProvider{InParamsPtr: &MockServiceInParams{}, Tran: &i18n.NopTranslator{}}
	copiedProvider := copyProvider(originProvider).(*MockProvider)
	fmt.Printf("[provider] origin: %p, copied: %p\n", originProvider, copiedProvider)
	fmt.Printf("[tran] origin: %p, copied: %p\n", originProvider.Tran, copiedProvider.Tran)
	fmt.Printf("[InParams] origin: %p, copied: %p\n", &originProvider.InParams, &copiedProvider.InParams)
	fmt.Printf("[InParamsPtr] origin: %p, copied: %p\n", originProvider.InParamsPtr, copiedProvider.InParamsPtr)
	fmt.Printf("[DefaultTop] origin: %p, copied: %p\n", &originProvider.DefaultTop, &copiedProvider.DefaultTop)
	assert.Equal(t, fmt.Sprintf("%p", originProvider.InParamsPtr), fmt.Sprintf("%p", copiedProvider.InParamsPtr))

	fmt.Println()

	deepCopiedProvider := deepcopy.Copy(originProvider).(*MockProvider)
	fmt.Printf("[provider] origin: %p, copied: %p\n", originProvider, deepCopiedProvider)
	fmt.Printf("[tran] origin: %p, copied: %p\n", originProvider.Tran, deepCopiedProvider.Tran)
	fmt.Printf("[InParams] origin: %p, copied: %p\n", &originProvider.InParams, &deepCopiedProvider.InParams)
	fmt.Printf("[InParamsPtr] origin: %p, copied: %p\n", originProvider.InParamsPtr, deepCopiedProvider.InParamsPtr)
	fmt.Printf("[DefaultTop] origin: %p, copied: %p\n", &originProvider.DefaultTop, &deepCopiedProvider.DefaultTop)
	assert.NotEqual(t, fmt.Sprintf("%p", originProvider.InParamsPtr), fmt.Sprintf("%p", deepCopiedProvider.InParamsPtr))
}
