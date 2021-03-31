// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"flag"
	"fmt"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	version = "v1.0.0"
	genName = "protoc-gen-go-http"
)

var (
	showVersion = flag.Bool("version", false, "print the version and exit")
	genAll      = flag.Bool("genall", false, "generate all service function")
)

func main() {
	flag.Parse()
	if *showVersion {
		fmt.Printf("%s %v\n", genName, version)
		return
	}

	var flags flag.FlagSet
	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(p *protogen.Plugin) error {
		p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range p.Files {
			if f.Generate {
				if _, err := generateFile(p, f); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
