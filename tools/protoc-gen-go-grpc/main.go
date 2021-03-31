// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"flag"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

var requireUnimplemented *bool

func main() {
	var flags flag.FlagSet
	requireUnimplemented = flags.Bool("requireUnimplementedServers", false, "unset to match legacy behavior")

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if f.Generate {
				generateFile(gen, f)
			}
		}
		return nil
	})
}
