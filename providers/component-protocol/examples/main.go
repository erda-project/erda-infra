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

package main

import (
	"embed"
	"os"

	"github.com/erda-project/erda-infra/base/logs"
	"github.com/erda-project/erda-infra/base/servicehub"
	"github.com/erda-project/erda-infra/providers/component-protocol"
	_ "github.com/erda-project/erda-infra/providers/component-protocol"
	"github.com/erda-project/erda-infra/providers/component-protocol/protocol"
	"github.com/erda-project/erda-infra/providers/i18n"
	_ "github.com/erda-project/erda-infra/providers/serviceregister"

	// import all scenarios
	_ "github.com/erda-project/erda-infra/providers/component-protocol/examples/scenarios"
)

//go:embed scenarios
var scenarioFS embed.FS

type config struct {
}

// +provider
type provider struct {
	Cfg *config
	Log logs.Logger

	Protocol componentprotocol.Interface
	Tran     i18n.Translator `translator:"dic"` // match dic.yaml
}

// Init .
func (p *provider) Init(ctx servicehub.Context) error {
	p.Log.Info("init demo")
	p.Protocol.SetI18nTran(p.Tran)          // use custom i18n translator
	p.Protocol.WithContextValue("k1", "v1") // test custom context kv

	// register protocols
	protocol.MustRegisterProtocolsFromFS(scenarioFS)

	return nil
}

func init() {
	servicehub.Register("demo", &servicehub.Spec{
		Services: []string{
			"demo-service",
		},
		Description: "here is description of demo",
		ConfigFunc: func() interface{} {
			return &config{}
		},
		Creator: func() servicehub.Provider {
			return &provider{}
		},
	})
}

func main() {
	hub := servicehub.New()
	hub.Run("demo", "", os.Args...)
}

// COMMAND: curl -s -XPOST localhost:8080/api/component-protocol/actions/render -H "Content-Type: application/json" \
//               -H "lang: zh" -d '{"scenario":{"scenarioType":"demo","scenarioKey":"demo"}}' | jq
// HTTP RESPONSE OUTPUT:
// {
//  "scenario": {
//    "scenarioKey": "demo"
//  },
//  "protocol": {
//    "scenario": "demo",
//    "hierarchy": {
//      "root": "demoTable",
//      "structure": {
//        "demoTable": []
//      }
//    },
//    "components": {
//      "demoTable": {
//        "type": "Table",
//        "name": "demoTable",
//        "props": {
//          "columns": [
//            {
//              "dataIndex": "sn",
//              "title": "编号"
//            },
//            {
//              "dataIndex": "name",
//              "title": "名字"
//            },
//            {
//              "dataIndex": "helloMsg",
//              "title": "欢迎消息"
//            }
//          ],
//          "pageSizeOptions": [
//            "10",
//            "20",
//            "1000"
//          ],
//          "rowKey": "sn"
//        },
//        "data": {
//          "list": [
//            {
//              "helloMsg": "你好: 张三 (666)",
//              "name": "张三",
//              "sn": "1"
//            },
//            {
//              "helloMsg": "你好 克里托斯",
//              "name": "克里托斯",
//              "sn": "2"
//            }
//          ]
//        }
//      }
//    },
//    "rendering": {
//      "__DefaultRendering__": [
//        {
//          "name": "demoTable",
//          "state": []
//        }
//      ]
//    }
//  }
//}

// COMMAND: curl -s -XPOST localhost:8080/api/component-protocol/actions/render -H "Content-Type: application/json" \
//               -H "lang: en" -d '{"scenario":{"scenarioType":"demo","scenarioKey":"demo"}}' | jq
// HTTP RESPONSE OUTPUT:
// {
//  "scenario": {
//    "scenarioKey": "demo"
//  },
//  "protocol": {
//    "scenario": "demo",
//    "hierarchy": {
//      "root": "demoTable",
//      "structure": {
//        "demoTable": []
//      }
//    },
//    "components": {
//      "demoTable": {
//        "type": "Table",
//        "name": "demoTable",
//        "props": {
//          "columns": [
//            {
//              "dataIndex": "sn",
//              "title": "SN"
//            },
//            {
//              "dataIndex": "name",
//              "title": "Name"
//            },
//            {
//              "dataIndex": "helloMsg",
//              "title": "Hello Message"
//            }
//          ],
//          "pageSizeOptions": [
//            "10",
//            "20",
//            "1000"
//          ],
//          "rowKey": "sn"
//        },
//        "data": {
//          "list": [
//            {
//              "helloMsg": "hello my friend: zhangsan (666)",
//              "name": "zhangsan",
//              "sn": "1"
//            },
//            {
//              "helloMsg": "hello my friend kratos",
//              "name": "kratos",
//              "sn": "2"
//            }
//          ]
//        }
//      }
//    },
//    "rendering": {
//      "__DefaultRendering__": [
//        {
//          "name": "demoTable",
//          "state": []
//        }
//      ]
//    }
//  }
//}
