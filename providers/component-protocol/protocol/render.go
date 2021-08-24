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
	"time"

	"github.com/sirupsen/logrus"

	"github.com/erda-project/erda-infra/providers/component-protocol/cptype"
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
		p, err := getDefaultProtocol(sk)
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
		logrus.Errorf("faield to get scenario render, err: %v", err)
		return err
	}

	var compRending []cptype.RendingItem
	if useDefaultProtocol {
		// 如果是加载默认协议，则渲染所有组件，不存在组件及状态依赖
		crs, ok := req.Protocol.Rendering[cptype.DefaultRenderingKey]
		if !ok {
			// 如果不存在默认DefaultRenderingKey，则随机从map中获取各组件key，渲染所有组件，无state绑定
			for k := range *sr {
				if _, ok := req.Protocol.Components[k]; !ok {
					continue
				}
				ri := cptype.RendingItem{
					Name: k,
				}
				compRending = append(compRending, ri)
			}
		} else {
			// 如果存在默认DefaultRenderingKey，则获取默认定义的Rending
			compRending = append(compRending, crs...)
		}

	} else {
		// 如果是前端触发一个组件操作，则先渲染该组件;
		// 再根据定义的渲染顺序，依次完成其他组件的state注入和渲染;
		compName := req.Event.Component
		compRending = append(compRending, cptype.RendingItem{Name: compName})

		crs, ok := req.Protocol.Rendering[compName]
		if !ok {
			logrus.Infof("empty protocol rending for component:%s", compName)
		} else {
			compRending = append(compRending, crs...)
		}
	}
	compRending = polishComponentRendering(req.DebugOptions, compRending)

	if req.Protocol.GlobalState == nil {
		gs := make(cptype.GlobalStateData)
		req.Protocol.GlobalState = &gs
	}

	// clean pre render error
	setGlobalStateKV(req.Protocol, cptype.GlobalInnerKeyError.String(), nil)

	polishProtocol(req.Protocol)

	for _, v := range compRending {
		// 组件状态渲染
		err := protoCompStateRending(ctx, req.Protocol, v)
		if err != nil {
			logrus.Errorf("protocol component state rending failed, request:%+v, err: %v", v, err)
			return err
		}
		// 获取协议中相关组件
		c, err := getProtoComp(ctx, req.Protocol, v.Name)
		if err != nil {
			logrus.Errorf("get component from protocol failed, scenario:%s, component:%s", sk, req.Event.Component)
			return nil
		}
		// 获取组件渲染函数
		cr, err := getCompRender(ctx, *sr, v.Name, c.Type)
		if err != nil {
			logrus.Errorf("get component render failed, scenario:%s, component:%s", sk, req.Event.Component)
			return err
		}
		// 生成组件对应事件，如果不是组件自身事件则为默认事件
		event := eventConvert(v.Name, req.Event)
		// 运行组件渲染函数
		start := time.Now() // 获取当前时间
		err = wrapCompRender(cr.RenderC(), req.Protocol.Version).Render(ctx, c, req.Scenario, event, req.Protocol.GlobalState)
		if err != nil {
			logrus.Errorf("render component failed,err: %s, scenario:%+v, component:%s", err.Error(), req.Scenario, cr.CompName)
			return err
		}
		elapsed := time.Since(start)
		logrus.Infof("[component render time cost] scenario: %s, component: %s, cost: %s", req.Scenario.ScenarioKey, cr.CompName, elapsed)
	}
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

