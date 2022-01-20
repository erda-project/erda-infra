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

package protocol

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-infra/pkg/strutil"
	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
	"github.com/erda-project/erda-infra/providers/component-protocol/protocol/posthook"
	"github.com/erda-project/erda-infra/providers/component-protocol/utils/cputil"
)

// RunScenarioRender .
func RunScenarioRender(ctx context.Context, req *cptype.ComponentProtocolRequest) error {
	// check debug options
	if err := checkDebugOptions(ctx, req.DebugOptions); err != nil {
		return err
	}

	// get scenario key
	sk, err := getScenarioKey(req.Scenario)
	if err != nil {
		return err
	}

	var useDefaultProtocol bool
	if req.Protocol == nil || req.Event.Component == "" {
		useDefaultProtocol = true
		p, err := getDefaultProtocol(ctx, sk)
		if err != nil {
			return err
		}
		var tmp cptype.ComponentProtocol
		if err := cputil.ObjJSONTransfer(&p, &tmp); err != nil {
			logrus.Errorf("deep copy failed, err: %v", err)
			return err

		}
		req.Protocol = &tmp
	}

	sr, err := getScenarioRenders(sk)
	if err != nil {
		logrus.Errorf("failed to get scenario render, err: %v", err)
		return err
	}

	var compRending []cptype.RendingItem
	if useDefaultProtocol {
		crs, ok := req.Protocol.Rendering[cptype.DefaultRenderingKey]
		if !ok {
			orders, err := calculateDefaultRenderOrderByHierarchy(req.Protocol)
			if err != nil {
				logrus.Errorf("failed to calculate default render order by hierarchy: %v", err)
				return err
			}
			for _, compName := range orders {
				compRending = append(compRending, cptype.RendingItem{Name: compName})
			}
		} else {
			compRending = append(compRending, crs...)
		}

	} else {
		// root is always rendered
		if req.Event.Component != req.Protocol.Hierarchy.Root {
			compRending = append(compRending, cptype.RendingItem{Name: req.Protocol.Hierarchy.Root})
		}
		// 如果是前端触发一个组件操作，则先渲染该组件;
		// 再根据定义的渲染顺序，依次完成其他组件的state注入和渲染;
		compName := req.Event.Component
		compRending = append(compRending, cptype.RendingItem{Name: compName})

		crs, ok := req.Protocol.Rendering[compName]
		if !ok {
			logrus.Infof("empty protocol rending for component: %s, use hierarchy and start from %s", compName, compName)
			orders, err := calculateDefaultRenderOrderByHierarchy(req.Protocol)
			if err != nil {
				logrus.Errorf("failed to calculate default render order by hierarchy for empty rendering: %v", err)
				return err
			}
			subRenderingOrders := getDefaultHierarchyRenderOrderFromCompExclude(orders, compName)
			for _, comp := range subRenderingOrders {
				compRending = append(compRending, cptype.RendingItem{Name: comp})
			}
		} else {
			compRending = append(compRending, crs...)
		}
	}
	compRending = polishComponentRendering(req.DebugOptions, compRending)
	compRending = polishComponentRenderingByInitOp(req.Protocol, req.Event, compRending)
	compRending = polishComponentRenderingByAsyncAtInitOp(req.Protocol, req.Event, compRending)

	if req.Protocol.GlobalState == nil {
		gs := make(cptype.GlobalStateData)
		req.Protocol.GlobalState = &gs
	}

	// clean pre render error
	setGlobalStateKV(req.Protocol, cptype.GlobalInnerKeyError.String(), nil)

	polishProtocol(req.Protocol)

	// if hierarchy.Parallel specified, use new rendering
	if len(req.Protocol.Hierarchy.Parallel) > 0 {
		rootNode, err := parseParallelRendering(req.Protocol, compRending)
		if err != nil {
			logrus.Errorf("failed to parse parallel rendering, err: %v", err)
			return err
		}
		fmt.Println(rootNode.String())
		return renderFromNode(ctx, req, *sr, rootNode)
	}

	for _, v := range compRending {
		if err := renderOneComp(ctx, req, *sr, v); err != nil {
			return err
		}
	}

	posthook.HandleContinueRender(compRending, req.Protocol)
	//posthook.OnlyReturnRenderingComps(compRending, req.Protocol)

	return nil
}

func polishComponentRendering(debugOptions *cptype.ComponentProtocolDebugOptions, compRendering []cptype.RendingItem) []cptype.RendingItem {
	if debugOptions == nil || debugOptions.ComponentKey == "" {
		return compRendering
	}
	var result []cptype.RendingItem
	for _, item := range compRendering {
		if item.Name == debugOptions.ComponentKey {
			result = append(result, item)
			break
		}
	}
	return result
}

func polishComponentRenderingByInitOp(protocol *cptype.ComponentProtocol, event cptype.ComponentEvent, compRendering []cptype.RendingItem) []cptype.RendingItem {
	// judge event
	if event.Component != cptype.InitializeOperation.String() {
		return compRendering
	}
	// judge async comps
	asyncCompsByName := make(map[string]struct{})
	for _, comp := range protocol.Components {
		if comp.Options != nil && comp.Options.AsyncAtInit {
			asyncCompsByName[comp.Name] = struct{}{}
		}
	}
	// skip comp with option: asyncAtInit
	var result []cptype.RendingItem
	for _, item := range compRendering {
		if _, needAsyncAtInit := asyncCompsByName[item.Name]; needAsyncAtInit {
			continue
		}
		result = append(result, item)
	}
	return result
}

func polishComponentRenderingByAsyncAtInitOp(protocol *cptype.ComponentProtocol, event cptype.ComponentEvent, compRendering []cptype.RendingItem) []cptype.RendingItem {
	// judge event
	if event.Component != cptype.AsyncAtInitOperation.String() {
		return compRendering
	}
	asyncCompsByName := make(map[string]struct{})
	if len(event.OperationData) > 0 {
		v, ok := event.OperationData["components"]
		if ok {
			for _, vv := range v.([]interface{}) {
				asyncCompsByName[strutil.String(vv)] = struct{}{}
			}
		}
	}
	if len(asyncCompsByName) == 0 {
		// analyze from protocol
		for _, comp := range protocol.Components {
			if comp.Options != nil && comp.Options.AsyncAtInit {
				asyncCompsByName[comp.Name] = struct{}{}
			}
		}
	}
	// only render async comps
	var result []cptype.RendingItem
	for _, item := range compRendering {
		if _, needAsyncAtInit := asyncCompsByName[item.Name]; needAsyncAtInit {
			result = append(result, item)
		}
	}
	return result
}

func getCompNameAndInstanceName(name string) (compName, instanceName string) {
	ss := strings.SplitN(name, "@", 2)
	if len(ss) == 2 {
		compName = ss[0]
		instanceName = ss[1]
		return
	}
	compName = name
	// use name as instance name
	instanceName = name
	return
}

func calculateDefaultRenderOrderByHierarchy(p *cptype.ComponentProtocol) ([]string, error) {
	allCompSubMap := make(map[string][]string)
	for k, v := range p.Hierarchy.Structure {
		switch subs := v.(type) {
		case []interface{}:
			for i := range subs {
				allCompSubMap[k] = append(allCompSubMap[k], *recursiveGetSubComps(nil, subs[i])...)
			}
		case map[string]interface{}:
			allCompSubMap[k] = append(allCompSubMap[k], *recursiveGetSubComps(nil, subs["left"])...)
			allCompSubMap[k] = append(allCompSubMap[k], *recursiveGetSubComps(nil, subs["right"])...)
			childrenComps := *recursiveGetSubComps(nil, subs["children"])
			footerComps := *recursiveGetSubComps(nil, subs["footer"])
			// recognized structKey: left, right, children, footer
			for structKey, compName := range subs {
				if structKey == "left" || structKey == "right" || structKey == "children" || structKey == "footer" {
					continue
				}
				allCompSubMap[k] = append(allCompSubMap[k], *recursiveGetSubComps(nil, compName)...)
			}
			for _, comp := range childrenComps {
				allCompSubMap[k] = append(allCompSubMap[k], *recursiveGetSubComps(nil, comp)...)
			}
			for _, comp := range footerComps {
				allCompSubMap[k] = append(allCompSubMap[k], *recursiveGetSubComps(nil, comp)...)
			}
		}
		allCompSubMap[k] = strutil.DedupSlice(allCompSubMap[k], true)
	}

	root := p.Hierarchy.Root
	var results []string
	if walkErr := recursiveWalkCompOrder(root, &results, allCompSubMap); walkErr != nil {
		return nil, walkErr
	}
	results = strutil.DedupSlice(results, true)
	return results, nil
}

func recursiveGetSubComps(result *[]string, subs interface{}) *[]string {
	if result == nil {
		result = &[]string{}
	}
	if subs == nil {
		return result
	}
	switch v := subs.(type) {
	case []interface{}:
		for _, vv := range v {
			recursiveGetSubComps(result, vv)
		}
	case map[string]interface{}:
		for _, vv := range v {
			recursiveGetSubComps(result, vv)
		}
	case string:
		*result = append(*result, v)
	case float64:
		*result = append(*result, strutil.String(v))
	default:
		panic(fmt.Errorf("not supported type: %v, subs: %v", reflect.TypeOf(subs), subs))
	}
	return result
}

// recursiveWalkCompOrder
// TODO check cycle visited
func recursiveWalkCompOrder(current string, orders *[]string, allCompSubMap map[string][]string) error {
	*orders = append(*orders, current)
	subs := allCompSubMap[current]
	for _, sub := range subs {
		if err := recursiveWalkCompOrder(sub, orders, allCompSubMap); err != nil {
			return err
		}
	}
	*orders = strutil.DedupSlice(*orders, true)

	return nil
}

func getDefaultHierarchyRenderOrderFromCompExclude(fullOrders []string, startFromCompExclude string) []string {
	fromIdx := -1
	for i, comp := range fullOrders {
		if startFromCompExclude == comp {
			fromIdx = i
			break
		}
	}
	if fromIdx == -1 {
		return []string{}
	}
	return fullOrders[fromIdx+1:]
}

func simplifyComp(comp *cptype.Component) {
	if len(comp.Data) == 0 {
		comp.Data = nil
	}
	if len(comp.State) == 0 {
		comp.State = nil
	}
	if len(comp.Operations) == 0 {
		comp.Operations = nil
	}
	if comp.Options != nil {
		if !comp.Options.Visible &&
			!comp.Options.AsyncAtInit &&
			!comp.Options.FlatExtra &&
			!comp.Options.RemoveExtraAfterFlat {
			comp.Options = nil
		}
	}
	if len(comp.Props) == 0 {
		comp.Props = nil
	}
}

func renderOneComp(ctx context.Context, req *cptype.ComponentProtocolRequest, sr ScenarioRender, v cptype.RendingItem) error {
	// 组件状态渲染
	err := protoCompStateRending(ctx, req.Protocol, v)
	if err != nil {
		logrus.Errorf("protocol component state rending failed, request: %+v, err: %v", v, err)
		return err
	}
	// 获取协议中相关组件
	c, err := getProtoComp(ctx, req.Protocol, v.Name)
	if err != nil {
		logrus.Errorf("get component from protocol failed, scenario: %s, component: %s", req.Scenario.ScenarioKey, req.Event.Component)
		return nil
	}
	// 获取组件渲染函数
	cr, err := getCompRender(ctx, sr, v.Name, c.Type)
	if err != nil {
		logrus.Errorf("get component render failed, scenario: %s, component: %s", req.Scenario.ScenarioKey, req.Event.Component)
		return err
	}
	// 生成组件对应事件，如果不是组件自身事件则为默认事件
	event := eventConvert(v.Name, req.Event)
	// 运行组件渲染函数
	start := time.Now() // 获取当前时间
	_, instanceName := getCompNameAndInstanceName(v.Name)
	c.Name = instanceName
	err = wrapCompRender(cr.RenderC(), req.Protocol.Version).Render(ctx, c, req.Scenario, event, req.Protocol.GlobalState)
	if err != nil {
		logrus.Errorf("render component failed, err: %s, scenario: %+v, component: %s", err.Error(), req.Scenario, cr.CompName)
		return err
	}
	simplifyComp(c)
	elapsed := time.Since(start)
	logrus.Infof("[component render time cost] scenario: %s, component: %s, cost: %s", req.Scenario.ScenarioKey, v.Name, elapsed)
	return nil
}
