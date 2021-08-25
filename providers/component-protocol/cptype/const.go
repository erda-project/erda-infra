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

// GlobalInnerCtxKey .
type GlobalInnerCtxKey string

func (k GlobalInnerCtxKey) String() string { return string(k) }

const (
	GlobalInnerKeyCtxSDK   GlobalInnerCtxKey = "_sdk_"
	GlobalInnerKeyUserIDs  GlobalInnerCtxKey = "_userIDs_"
	GlobalInnerKeyIdentity GlobalInnerCtxKey = "_identity_"
	GlobalInnerKeyError    GlobalInnerCtxKey = "_error_"
)

const (
	DefaultRenderingKey     = "__DefaultRendering__"
	InParamsStateBindingKey = "__InParams__"
)

const (
	DefaultComponentNamespace = "__DefaultComponentNamespace__"
)
