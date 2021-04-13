# Erda Infra

[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
[![codecov](https://codecov.io/gh/erda-project/erda-infra/branch/develop/graph/badge.svg?token=SVROJLY8UK)](https://codecov.io/gh/erda-project/erda-infra)

翻译: [English](README.md) | [简体中文](README_zh.md)

Erda Infra 一套轻量级 Go 微服务框架，包含大量现成的模块和工具，能够快速构建起以模块化驱动的应用程序。

一些 Go 项目基于该框架进行构建:
* [Erda](https://github.com/erda-project/erda)

## 特性
* 以模块化设计方式来驱动应用系统实现，支持模块可插拔
* 统一配置读取，支持默认值、支持从文件、环境变量、命令行参数读取
* 统一模块的初始化、启动、关闭
* 统一管理模块间的依赖关系
* 支持模块间的依赖注入
* 包含大量现成的微模块
* 支持统一 gRPC 和 HTTP 接口设计、以及拦截器
* 提供快速构建模块的脚本
* 等等

## 概念 
* Service，服务，表示某个具体的功能
* Provider，服务的提供者，提供0个或多个 Service，也可以依赖0个或多个其他 Service，被依赖的 Service 由其他 Provider 提供
* ProviderDefine，提供 Provider 相关的元信息，比如：提供 Provider 的构造函数。通过 *servicehub.RegisterProvider* 来注册 Provider
* Hub，是所有 Provider 的容器，管理所有已加载的 Provider 的生命周期

所有已注册的 Provider 通过 一份配置来确定是否 加载，由 Hub 对已加载的 Provider 的进行初始化、启动、关闭等。

![servicehub](./docs/servicehub.jpg)

## Provider 定义
通过实现 *servicehub.ProviderDefine* 接口来定义一个模块，并 通过 *servicehub.RegisterProvider* 函数进行注册。

但更简单的是通过 *servicehub.Spec* 来描述一个模块，并 通过 *servicehub.Register* 函数进行注册。

[例子](./base/servicehub/examples)

## Quick Start
### 快速创建一个模块
第一步，创建模块
```sh
➜ gohub init -o helloworld
Input Service Provider Name: helloworld
➜ # 以上命令创建了一个模块的模版代码，文件如下：
➜ tree helloworld
helloworld
├── provider.go
└── provider_test.go
```

第二步，创建 main.go
```go
package main

import (
	"github.com/erda-project/erda-infra/base/servicehub"
	_ "./helloworld" // your package import path
)

func main() {
	servicehub.Run(&servicehub.RunOptions{
		Content: `
helloworld:
`,
	})
}
```

第三步，运行程序
```sh
➜ go run main.go
INFO[2021-04-13 13:17:36.416] message: hi                                   module=helloworld
INFO[2021-04-13 13:17:36.416] provider helloworld initialized              
INFO[2021-04-13 13:17:36.416] signals to quit: [hangup interrupt terminated quit] 
INFO[2021-04-13 13:17:36.426] provider helloworld running ...              
INFO[2021-04-13 13:17:39.429] do something...                               module=helloworld
```
[Hello World](./examples/example) \( [helloworld/](./examples/example/helloworld) | [main.go](./examples/example/main.go) \)

### 创建 HTTP/gRPC 服务
这些服务既可以被远程调用，也可以被本地模块调用。

第一步，在 *.proto 文件中定义协议 (消息结构 和 接口)
```protobuf
syntax = "proto3";

package erda.infra.example;
import "google/api/annotations.proto";
option go_package = "github.com/erda-project/erda-infra/examples/service/protocol/pb";

// the greeting service definition.
service GreeterService {
  // say hello
  rpc SayHello (HelloRequest) returns (HelloResponse)  {
    option (google.api.http) = {
      get: "/api/greeter/{name}",
    };
  }
}

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  bool success = 1;
  string data = 2;
}
```

第二步，编译生成接口 和 客户端代码
```sh
➜ gohub protoc protocol *.proto 
➜ tree 
.
├── client
│   ├── client.go
│   └── provider.go
├── greeter.proto
└── pb
    ├── greeter.form.pb.go
    ├── greeter.http.pb.go
    ├── greeter.pb.go
    ├── greeter_grpc.pb.go
    └── register.services.pb.go
```

第三步，实现协议接口
```sh
➜ gohub protoc imp *.proto --imp_out=../server/helloworld
➜ tree ../server/helloworld
../server/helloworld
├── greeter.service.go
├── greeter.service_test.go
└── provider.go
```

第四步，创建 main.go 启动程序

*main.go*
```
package main

import (
	"os"

	"github.com/erda-project/erda-infra/base/servicehub"

	// import all providers
	_ "github.com/erda-project/erda-infra/examples/service/server/helloworld"
	_ "github.com/erda-project/erda-infra/providers"
)

func main() {
	hub := servicehub.New()
	hub.Run("server", "server.yaml", os.Args...)
}
```

*server.yaml*
```yaml
# optional
http-server:
    addr: ":8080"
grpc-server:
    addr: ":7070"
service-register:
# expose services and interface
erda.infra.example:
```

[Service](./examples/service) \( [Protocol](./examples/service/protocol) | [Implementation](./examples/service/server/helloworld) | [Server](./examples/service/server) | [Caller](./examples/service/caller) | [Client](./examples/service/client)  \)


## 微模块
该项目中已经封装了许多可用的模块，在 [providers/](./providers) 目录下可以找到。

每一个模块下面，都有一个 examples 目录，包含了该模块的使用例子。

* elasticsearch，对 elasticsearch 客户端的封装，更方便的进行批量数据的写入
* etcd，对 etcd 客户端的封装
* etcd-mutex，利用 etcd 实现的分布式锁
* grpcserver，启动一个 gRPC server
* grpcclient，统一管理 gRPC 客户端
* health，通过 httpserver 注册一个健康检查的接口
* httpserver，提供一个 HTTP server, 支持任意形式的处理函数、拦截器、参数绑定、参数校验等
* i18n，提供了国际化的支持，可以统一管理国际化文件、支持模版
* kafka，提供了访问 kafka 相关的能力，更方便地去批量消费和推送消息
* kubernetes，对 kubernetes 客户端的封装
* mysql，对 mysql 客户端的封装
* pprof，通过 httpserver 注册一些 pprof 相关的接口
* redis，对 redis 客户端的封装
* zk-master-election，通过 zookeeper 实现主从选举
* zookeeper，对 zookeeper 客户端的封装
* cassandra，对 Cassandra 客户端的封装
* serviceregister，封装提供统一注册 gRPC 和 HTTP 接口的能力

# 工具
*gohub* 是一个能够帮助您快速构建模块的命令行工具，可以通过如下方式安装： 
```sh
git clone https://github.com/erda-project/erda-infra 
cd erda-infra/tools/gohub
go install .
```

也可以通过 Docker 容器来使用以下工具:
```sh
➜ docker run --rm -ti -v $(pwd):/go \
    registry.cn-hangzhou.aliyuncs.com/dice/erda-tools:1.0 gohub                                                                
Usage:
  gohub [flags]
  gohub [command]

Available Commands:
  help        Help about any command
  init        Initialize a provider with name
  pkgpath     Print the absolute path of go package
  protoc      ProtoBuf compiler tools
  tools       Tools
  version     Print the version number

Flags:
  -h, --help   help for gohub

Use "gohub [command] --help" for more information about a command.
```

## License
Erda Infra is under the Apache 2.0 license. See the [LICENSE](/LICENSE) file for details.