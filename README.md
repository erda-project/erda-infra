# Erda Infra

[![License](https://img.shields.io/badge/license-Apache%202-4EB1BA.svg)](https://www.apache.org/licenses/LICENSE-2.0.html)
[![codecov](https://codecov.io/gh/erda-project/erda-infra/branch/develop/graph/badge.svg?token=SVROJLY8UK)](https://codecov.io/gh/erda-project/erda-infra)

Translations: [English](README.md) | [简体中文](README_zh.md)

Erda Infra is a lightweight microservices framework implements by golang, which offers many useful modules and tools to help you quickly build a module-driven applications.

Many Go projects are built using Erda Infra including:
* [Erda](https://github.com/erda-project/erda)

## Features
* modular design to drive the implementation of the application system, and each module is pluggable.
* each module is configurable，supports setting defaults、reading from files (YAML, HCL,、JSON、TOML、.env file)、environment、flags.
* manage the lifecycle of the module, includes initialization, startup, and shutdown.
* manage dependencies between modules.
* support Dependency Injection between modules.
* offers many commonly modules that can be used directly.
* support define APIs and models in protobuf file to expose both gRPC and HTTP APIs.
* offers tools to help you quickly build a module.
* etc.

## Concept
* Service, represents a function.
* Provider, service provider, equivalent to module, provide some services. It can also depend on other services, witch provide by other provider.
* ProviderDefine, describe a provider, includes provider's name, constructor function of provider, services, etc. Register by it *servicehub.RegisterProvider* function.
* Hub, is a container for all providers, and manage the life cycle of all loaded providers.

A configuration is used to determine whether all registered Providers are loaded, and the Hub initializes, starts, and closes the loaded Providers.

## Define Provider
Define a provider by implementing the *servicehub.ProviderDefine* interface, and register it through the *servicehub.RegisterProvider* function.

But, it is simpler to describe a provider through *servicehub.Spec* and register it through the *servicehub.Register* function.

[Examples](./base/servicehub/examples)

## Quick Start

```sh
➜ # create service interface
➜ ROOT_PATH=$(pwd)
➜ ${ROOT_PATH}/tools/protoc.sh protocol "examples/protocol/*.proto"
➜ 
➜ # create module 
➜ mkdir -p examples/server/helloworld
➜ cd examples/server/helloworld
➜ ${ROOT_PATH}/tools/protoc.sh init "${ROOT_PATH}/examples/protocol/*.proto"
➜ 
➜ # implement the service interface in examples/server/helloworld directory
➜ 
➜ cd ${ROOT_PATH}/examples/server
➜ 
➜ # create main.go, like examples/server/main.go
➜ # create server.yaml, like examples/server/server.yaml
➜ 
➜ go run main.go
```
![example](./examples/example.jpg)

[Hello World](./examples) \( [Server](./examples/server) | [Client](./examples/client) \)

## Useful Providers
Many available providers have been packaged in this project, it can be found in the [providers/](./providers) directory.

Under each module, there is an examples directory, which contains examples of the use of the module.

* elasticsearch, provide elasticsearch client APIs, and it easy to write batch data.
* etcd, provide etcd client APIs.
* etcd-mutex, distributed lock implemented by etcd.
* grpcserver, start a grpc server.
* grpcclient, provide gRPC client, and manage client's connection.
* health, provide health check API, and can register some health check function into this provider.
* httpserver, provide an HTTP server, support any form of handle function, interceptor, parameter binding, parameter verification, etc.
* i18n, provide internationalization support, manage i18n files, support templates.
* kafka, provide kafka sdk, and easy to produce and consume messages.
* kubernetes, provide kubernetes client APIs.
* mysql, provide mysql client APIs.
* pprof, expose pprof HTTP APIs By httpserver.
* redis, provide redis client APIs.
* zk-master-election, provide interface about master-slave election, it is implemented by zookeeper.
* zookeeper, provide zookeeper client.
* serviceregister, use it to register services to expose gRPC and HTTP APIs.

# 工具
protoc-gen-go-* tools depends on protobuf compiler，see [protobuf](https://github.com/protocolbuffers/protobuf) to install protoc。

You can also use the following tools through a Docker container.

```sh
docker run --rm -ti -v $(pwd):/go \
    registry.cn-hangzhou.aliyuncs.com/dice/erda-tools:1.0 protoc.sh usage
```

* protoc-gen-go-grpc, according to *.proto file, provide gRPC server and client support
* protoc-gen-go-http, according to the *.proto file, provide HTTP server support for the defined service.
* protoc-gen-go-form, according to the *.proto file, provide HTTP form codec support for the defined message.
* protoc-gen-go-client, according to the *.proto file, generate clients and provider for defined service.
* protoc-gen-go-register, provide some functions to help a provider register services.
* protoc-gen-go-provider, according to the *.proto file, generate a provider template that implements the service interface to facilitate rapid module development.
* protoc.sh, wrap the protoc-gen-go-* series of tools for the development.

## License
Erda is under the Apache 2.0 license. See the [LICENSE](/LICENSE) file for details.