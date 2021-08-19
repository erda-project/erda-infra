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

package cptype

const (
	// 协议定义的操作
	// 用户通过URL初次访问
	InitializeOperation OperationKey = "__Initialize__"
	// 用于替换掉DefaultOperation，前端触发某组件，在协议Rending中定义了关联渲染组件，传递的事件是RendingOperation
	RenderingOperation OperationKey = "__Rendering__"
)

// 组件化协议定义
type ComponentProtocol struct {
	Version     string                   `json:"version" yaml:"version"`
	Scenario    string                   `json:"scenario" yaml:"scenario"`
	GlobalState *GlobalStateData         `json:"state" yaml:"state"`
	Hierarchy   Hierarchy                `json:"hierarchy" yaml:"hierarchy"`
	Components  map[string]*Component    `json:"components" yaml:"components"`
	Rendering   map[string][]RendingItem `json:"rendering" yaml:"rendering"`
}

type GlobalStateData map[string]interface{}

// Hierarchy只是前端关心，只读，且有些字结构是字典有些是列表，后端不需要关心这部分
type Hierarchy struct {
	Version string `json:"version" yaml:"version"`
	Root    string `json:"root" yaml:"root"`
	// structure的结构可能是list、map
	Structure map[string]interface{} `json:"structure" yaml:"structure"`
}

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

type ComponentData map[string]interface{}

type ComponentOps map[string]interface{}

type Operation struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Reload   bool   `json:"reload"`
	FillMeta string `json:"fillMeta"`
}

type RendingItem struct {
	Name  string         `json:"name" yaml:"name"`
	State []RendingState `json:"state" yaml:"state"`
}

type RendingState struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

//type ComponentRenderCtx ComponentProtocolRequest

// request
type ComponentProtocolRequest struct {
	Scenario ComponentProtocolScenario `json:"scenario"`
	Event    ComponentEvent            `json:"event"`
	InParams map[string]interface{}    `json:"inParams"`
	// 初次请求为空，事件出发后，把包含状态的protocol传到后端
	Protocol *ComponentProtocol `json:"protocol"`

	// DebugOptions debug 选项
	DebugOptions *ComponentProtocolDebugOptions `json:"debugOptions,omitempty"`
}

type ComponentProtocolScenario struct {
	// 场景类型, 没有则为空
	ScenarioType string `json:"scenarioType" query:"scenarioType"`
	// 场景Key
	ScenarioKey string `json:"scenarioKey" query:"scenarioKey"`
}

type ComponentEvent struct {
	Component     string                 `json:"component"`
	Operation     OperationKey           `json:"operation"`
	OperationData map[string]interface{} `json:"operationData"`
}

type OperationKey string

func (o OperationKey) String() string {
	return string(o)
}

type ComponentProtocolParams interface{}

type ComponentProtocolDebugOptions struct {
	ComponentKey string `json:"componentKey"`
}
