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

package cptype

import (
	"context"
	"strconv"
	"sync"

	"github.com/erda-project/erda-infra/pkg/strutil"
	"github.com/erda-project/erda-infra/providers/component-protocol/protobuf/proto-go/cp/pb"
	"github.com/erda-project/erda-infra/providers/i18n"
)

// SDK .
type SDK struct {
	Ctx         context.Context
	Scenario    string
	Tran        i18n.Translator
	Identity    *pb.IdentityInfo
	InParams    InParams
	Lang        i18n.LanguageCodes
	GlobalState *GlobalStateData

	// ONLY FOR STD COMPONENT USE
	Event       ComponentEvent
	CompOpFuncs map[OperationKey]OperationFunc
	Comp        *Component

	// for parallel use, it's request level
	StdStructuredPtr IStdStructuredPtr

	Lock       sync.Mutex
	OriginLock *sync.Mutex
}

// Clone only return general-part of sdk to avoid concurrency issue.
func (sdk *SDK) Clone() *SDK {
	sdk.Lock.Lock()
	defer sdk.Lock.Unlock()

	clonedSDK := SDK{
		Ctx:      sdk.Ctx,
		Scenario: sdk.Scenario,
		Tran:     sdk.Tran,
		Identity: sdk.Identity,
		Lang:     sdk.Lang,

		OriginLock: &sdk.Lock,
	}
	// inParams
	clonedInParams := make(InParams)
	for k, v := range sdk.InParams {
		clonedInParams[k] = v
	}
	clonedSDK.InParams = clonedInParams
	// globalStates
	clonedGlobalStates := make(GlobalStateData)
	for k, v := range *sdk.GlobalState {
		clonedGlobalStates[k] = v
	}
	clonedSDK.GlobalState = &clonedGlobalStates

	return &clonedSDK
}

// I18n .
func (sdk *SDK) I18n(key string, args ...interface{}) string {
	if len(args) == 0 {
		try := sdk.Tran.Text(sdk.Lang, key)
		if try != key {
			return try
		}
	}
	return sdk.Tran.Sprintf(sdk.Lang, key, args...)
}

// String .
func (p InParams) String(key string) string {
	if p == nil {
		return ""
	}
	return strutil.String(p[key])
}

// Int64 .
func (p InParams) Int64(key string) int64 {
	if p == nil {
		return 0
	}
	i, ok := p[key]
	if !ok {
		return 0
	}
	switch v := i.(type) {
	case string:
		res, _ := strconv.ParseInt(v, 10, 64)
		return res
	default:
		res, _ := i.(int64)
		return res
	}
}

// Uint64 .
func (p InParams) Uint64(key string) uint64 {
	return uint64(p.Int64(key))
}

// RegisterOperation .
func (sdk *SDK) RegisterOperation(opKey OperationKey, opFunc OperationFunc) {
	sdk.CompOpFuncs[opKey] = opFunc
}

// SetUserIDs .
func (sdk *SDK) SetUserIDs(userIDs []string) {
	sdk.Lock.Lock()
	defer sdk.Lock.Unlock()
	(*sdk.GlobalState)[GlobalInnerKeyUserIDs.String()] = strutil.DedupSlice(userIDs, true)
}
