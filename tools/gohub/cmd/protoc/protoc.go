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

package protoc

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/erda-project/erda-infra/tools/gohub/cmd"
	"github.com/erda-project/erda-infra/tools/gohub/cmd/tools/install"
	"github.com/spf13/cobra"
)

func init() {
	messageCmd.Flags().String("msg_out", ".", "output directory of Message files")
	messageCmd.Flags().Bool("validate", false, "generate Validate function")
	messageCmd.Flags().Bool("json", true, "generate JSON Marshal and Unmarshal")
	messageCmd.Flags().StringSlice("json_opt", nil, "options for JSON Marshal and Unmarshal")
	messageCmd.Flags().StringSlice("include", nil, "include directory")
	messageCmd.Flags().Bool("gogofast", false, "use gogofast")
	protoCmd.AddCommand(messageCmd)

	protocolCmd.Flags().Bool("grpc", true, "support expose gRPC APIs")
	protocolCmd.Flags().Bool("http", true, "support expose HTTP APIs")
	protocolCmd.Flags().Bool("validate", false, "generate Validate function")
	protocolCmd.Flags().Bool("json", true, "generate JSON function")
	protocolCmd.Flags().StringSlice("json_opt", nil, "options for JSON Marshal and Unmarshal")
	protocolCmd.Flags().String("client_out", "./client", "output directory of gRPC Client files")
	protocolCmd.Flags().String("msg_out", "./pb", "output directory of Message files")
	protocolCmd.Flags().String("service_out", "./pb", "output directory of Service files")
	protocolCmd.Flags().StringSlice("include", nil, "include directory")
	protocolCmd.Flags().Bool("gogofast", false, "use gogofast")
	protoCmd.AddCommand(protocolCmd)

	implementCmd.Flags().String("imp_out", ".", "output directory of implementation")
	implementCmd.Flags().Bool("grpc", true, "implements gRPC APIs")
	implementCmd.Flags().Bool("http", true, "implements HTTP APIs")
	implementCmd.Flags().StringSlice("include", nil, "include directory")
	protoCmd.AddCommand(implementCmd)

	pluginCmd.Flags().StringSlice("opt", nil, "options for protoc plugin")
	pluginCmd.Flags().String("out", ".", "output directory of plugin")
	pluginCmd.Flags().StringSlice("include", nil, "include directory")
	protoCmd.AddCommand(pluginCmd)

	cmd.AddCommand(protoCmd)
}

var protoCmd = &cobra.Command{
	Use:     "protoc",
	Aliases: []string{"proto"},
	Short:   "ProtoBuf compiler tools",
}

var messageCmd = &cobra.Command{
	Use:     "message",
	Aliases: []string{"msg"},
	Short:   "Compile message only",
	Run: func(command *cobra.Command, args []string) {
		install.Download(false, cmd.Verbose())
		files := protoFiles(args)
		dirs := protoDirs(files)
		createMessage(command, args, files, dirs)
		fmt.Println("build successfully !")

	},
}

var protocolCmd = &cobra.Command{
	Use:   "protocol",
	Short: "Compile message, grpc, http codes",
	Run: func(command *cobra.Command, args []string) {
		install.Download(false, cmd.Verbose())
		files := protoFiles(args)
		dirs := protoDirs(files)
		createMessage(command, args, files, dirs)
		createService(command, args, files, dirs)
		fmt.Println("build successfully !")
	},
}

var implementCmd = &cobra.Command{
	Use:     "implement",
	Aliases: []string{"imp", "impl", "implements"},
	Short:   "Create a provider to implements protocol",
	Run: func(command *cobra.Command, args []string) {
		install.Download(false, cmd.Verbose())
		files := protoFiles(args)
		dirs := protoDirs(files)
		createImplementTemp(command, args, files, dirs)
		fmt.Println("build successfully !")
	},
}

var pluginCmd = &cobra.Command{
	Use:   "exec [plugin]",
	Short: "exec",
	Run: func(command *cobra.Command, args []string) {
		install.Download(false, cmd.Verbose())
		if len(args) < 1 {
			command.Usage()
			os.Exit(1)
		}
		plugin := args[0]
		out, err := command.Flags().GetString("out")
		cmd.CheckError(err)
		params := []string{
			fmt.Sprintf("--%s_out=%s", plugin, out), fmt.Sprintf("--%s_opt=paths=source_relative", plugin),
		}
		opts, err := command.Flags().GetStringSlice("opt")
		cmd.CheckError(err)
		if len(opts) > 0 {
			for _, op := range opts {
				params = append(params, fmt.Sprintf("--%s_opt=%s", plugin, op))
			}
		}
		files := protoFiles(args[1:])
		dirs := protoDirs(files)
		includes, _ := command.Flags().GetStringSlice("include")
		execProtoc(files, dirs, includes, params...)
	},
}

func protoFiles(args []string) []string {
	var files []string
	for i := len(args) - 1; i >= 0; i-- {
		if strings.HasPrefix(args[i], "-") {
			files = args[i+1:]
			break
		}
	}
	if files == nil {
		files = args
	}
	if len(files) <= 0 {
		cmd.CheckError(fmt.Errorf("no *.proto files specified"))
	}
	return files
}

func protoDirs(files []string) []string {
	dirset := make(map[string]struct{})
	for _, file := range files {
		dirset[filepath.Dir(file)] = struct{}{}
	}
	var dirs []string
	for d := range dirset {
		dirs = append(dirs, d)
	}
	sort.Strings(dirs)
	return dirs
}

func ensureOutputDir(command *cobra.Command, key, typ string) string {
	output, err := command.Flags().GetString(key)
	cmd.CheckError(err)
	if len(output) <= 0 {
		cmd.CheckError(fmt.Errorf("missing %s output directives", typ))
	}
	if !cmd.IsFileExist(output) {
		if !cmd.Confirm(fmt.Sprintf("%s not exit, create it?", output)) {
			cmd.CheckError(fmt.Errorf("%s not exit", output))
		}
		err = os.MkdirAll(output, os.ModePerm)
		cmd.CheckError(err)
	}
	return output
}

func execProtoc(files, dirs, include []string, params ...string) {
	for _, d := range dirs {
		params = append(params, fmt.Sprintf("-I=%s", d))
	}
	for _, d := range include {
		params = append(params, fmt.Sprintf("-I=%s", d))
	}
	includes := install.IncludeDirs()
	for _, include := range includes {
		if len(include) > 0 {
			params = append(params, fmt.Sprintf("-I=%s", include))
		}
	}

	params = append(params, fmt.Sprintf("-I=%s", "/usr/local/include/"))
	params = append(params, files...)
	fmt.Println("protoc", strings.Join(params, " "))
	proc := exec.Command("protoc", params...)
	proc.Stderr = os.Stderr
	proc.Stdout = os.Stdout
	proc.Stdin = os.Stdin
	err := proc.Run()
	cmd.CheckError(err)
}

func createMessage(command *cobra.Command, args, files, dirs []string) {
	output := ensureOutputDir(command, "msg_out", "Message")
	includes, _ := command.Flags().GetStringSlice("include")

	gogo, err := command.Flags().GetBool("gogofast")
	cmd.CheckError(err)
	if gogo {
		execProtoc(files, dirs, includes, fmt.Sprintf("--gogofast_out=%s", output))
	} else {
		execProtoc(files, dirs, includes,
			fmt.Sprintf("--go_out=%s", output), "--go_opt=paths=source_relative",
		)
	}

	valid, err := command.Flags().GetBool("validate")
	cmd.CheckError(err)
	if valid {
		execProtoc(files, dirs, includes,
			fmt.Sprintf("--govalidators_out=%s", output), "--govalidators_opt=paths=source_relative",
		)
	}

	json, err := command.Flags().GetBool("json")
	cmd.CheckError(err)
	if json {
		jsonOpts, _ := command.Flags().GetStringSlice("json_opt")
		params := []string{
			fmt.Sprintf("--go-json_out=%s", output), "--go-json_opt=paths=source_relative",
		}
		for _, opt := range jsonOpts {
			params = append(params, fmt.Sprintf("--go-json_opt=%s", opt))
		}
		execProtoc(files, dirs, includes, params...)
	}
}

func createService(command *cobra.Command, args, files, dirs []string) {
	createGRPC, err := command.Flags().GetBool("grpc")
	cmd.CheckError(err)
	if createGRPC {
		createGRPCService(command, args, files, dirs)
	}

	createHTTP, err := command.Flags().GetBool("http")
	cmd.CheckError(err)
	if createHTTP {
		createHTTPService(command, args, files, dirs)
	}

	if createGRPC || createHTTP {
		includes, _ := command.Flags().GetStringSlice("include")
		srvDir := ensureOutputDir(command, "service_out", "Service")
		execProtoc(files, dirs, includes,
			fmt.Sprintf("--go-register_out=%s", srvDir), "--go-register_opt=paths=source_relative",
			"--go-register_opt=grpc="+strconv.FormatBool(createGRPC), "--go-register_opt=http="+strconv.FormatBool(createHTTP),
		)
	}
}

func createGRPCService(command *cobra.Command, args, files, dirs []string) {
	clientDir := ensureOutputDir(command, "client_out", "Client")
	grpcDir := ensureOutputDir(command, "service_out", "Service")
	includes, _ := command.Flags().GetStringSlice("include")
	execProtoc(files, dirs, includes,
		fmt.Sprintf("--go-grpc_out=%s", grpcDir), "--go-grpc_opt=paths=source_relative",
		fmt.Sprintf("--go-client_out=%s", clientDir), "--go-client_opt=paths=source_relative",
	)
}

func createHTTPService(command *cobra.Command, args, files, dirs []string) {
	msgDir := ensureOutputDir(command, "msg_out", "Message")
	httpDir := ensureOutputDir(command, "service_out", "Service")
	includes, _ := command.Flags().GetStringSlice("include")
	execProtoc(files, dirs, includes,
		fmt.Sprintf("--go-http_out=%s", httpDir), "--go-http_opt=paths=source_relative",
		fmt.Sprintf("--go-form_out=%s", msgDir), "--go-form_opt=paths=source_relative",
	)
}

func createImplementTemp(command *cobra.Command, args, files, dirs []string) {
	createGRPC, err := command.Flags().GetBool("grpc")
	cmd.CheckError(err)
	createHTTP, err := command.Flags().GetBool("http")
	cmd.CheckError(err)
	impDir := ensureOutputDir(command, "imp_out", "Implementation")
	includes, _ := command.Flags().GetStringSlice("include")
	execProtoc(files, dirs, includes,
		fmt.Sprintf("--go-provider_out=%s", impDir), "--go-provider_opt=paths=source_relative",
		"--go-provider_opt=grpc="+strconv.FormatBool(createGRPC), "--go-provider_opt=http="+strconv.FormatBool(createHTTP),
	)
}
