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

// InitializeOperation .
const (
	// 协议定义的操作
	// 用户通过URL初次访问
	InitializeOperation OperationKey = "__Initialize__"
	// 用于替换掉DefaultOperation，前端触发某组件，在协议Rending中定义了关联渲染组件，传递的事件是RendingOperation
	RenderingOperation OperationKey = "__Rendering__"
)

// ComponentProtocol is protocol definition.
type ComponentProtocol struct {
	Version     string                   `json:"version" yaml:"version"`
	Scenario    string                   `json:"scenario" yaml:"scenario"`
	GlobalState *GlobalStateData         `json:"state" yaml:"state"`
	Hierarchy   Hierarchy                `json:"hierarchy" yaml:"hierarchy"`
	Components  map[string]*Component    `json:"components" yaml:"components"`
	Rendering   map[string][]RendingItem `json:"rendering" yaml:"rendering"`
}

// GlobalStateData .
type GlobalStateData map[string]interface{}

// Hierarchy represents components' hierarchy.
type Hierarchy struct {
	Version string `json:"version" yaml:"version"`
	Root    string `json:"root" yaml:"root"`
	// structure的结构可能是list、map
	Structure map[string]interface{} `json:"structure" yaml:"structure"`
}

// Component defines a component.
type Component struct {
	Version string `json:"version,omitempty" yaml:"version,omitempty"`
	// 组件类型
	Type string `json:"type,omitempty" yaml:"type,omitempty"`
	// 组件名字
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
	// table 动态字段
	Props interface{} `json:"props,omitempty" yaml:"props,omitempty"`
	// 组件业务数据
	Data ComponentData `json:"data,omitempty" yaml:"data,omitempty"`
	// 前端组件状态
	State map[string]interface{} `json:"state,omitempty" yaml:"state,omitempty"`
	// 组件相关操作（前端定义）
	Operations ComponentOps `json:"operations,omitempty" yaml:"operations,omitempty"`
}

// ComponentData .
type ComponentData map[string]interface{}

// ComponentOps .
type ComponentOps map[string]interface{}

// Operation .
type Operation struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Reload   bool   `json:"reload"`
	FillMeta string `json:"fillMeta"`
}

// RendingItem .
type RendingItem struct {
	Name  string         `json:"name" yaml:"name"`
	State []RendingState `json:"state" yaml:"state"`
}

// RendingState .
type RendingState struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

// ComponentProtocolRequest .
type ComponentProtocolRequest struct {
	Scenario Scenario               `json:"scenario"`
	Event    ComponentEvent         `json:"event"`
	InParams map[string]interface{} `json:"inParams"`
	// 初次请求为空，事件出发后，把包含状态的protocol传到后端
	Protocol *ComponentProtocol `json:"protocol"`

	// DebugOptions debug 选项
	DebugOptions *ComponentProtocolDebugOptions `json:"debugOptions,omitempty"`
}

// Scenario .
type Scenario struct {
	// 场景类型, 没有则为空
	ScenarioType string `json:"scenarioType" query:"scenarioType"`
	// 场景Key
	ScenarioKey string `json:"scenarioKey" query:"scenarioKey"`
}

// ComponentEvent .
type ComponentEvent struct {
	Component     string                 `json:"component"`
	Operation     OperationKey           `json:"operation"`
	OperationData map[string]interface{} `json:"operationData"`
}

// OperationKey .
type OperationKey string

// String .
func (o OperationKey) String() string {
	return string(o)
}

// ComponentProtocolParams .
type ComponentProtocolParams interface{}

// ComponentProtocolDebugOptions .
type ComponentProtocolDebugOptions struct {
	ComponentKey string `json:"componentKey"`
}
