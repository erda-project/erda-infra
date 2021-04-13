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
	"github.com/erda-project/erda-infra/tools/gohub/cmd/pkgpath"
	"github.com/erda-project/erda-infra/tools/gohub/cmd/tools/install"
	"github.com/spf13/cobra"
)

func init() {
	messageCmd.Flags().String("msg_out", ".", "output directory of Message files")
	protoCmd.AddCommand(messageCmd)

	protocolCmd.Flags().Bool("grpc", true, "support expose gRPC APIs")
	protocolCmd.Flags().Bool("http", true, "support expose HTTP APIs")
	protocolCmd.Flags().String("client_out", "./client", "output directory of gRPC Client files")
	protocolCmd.Flags().String("msg_out", "./pb", "output directory of Message files")
	protocolCmd.Flags().String("service_out", "./pb", "output directory of Service files")
	protoCmd.AddCommand(protocolCmd)

	implementCmd.Flags().String("imp_out", ".", "output directory of implementation")
	implementCmd.Flags().Bool("grpc", true, "implements gRPC APIs")
	implementCmd.Flags().Bool("http", true, "implements HTTP APIs")
	protoCmd.AddCommand(implementCmd)

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

func execProtoc(files, dirs []string, params ...string) {
	for _, d := range dirs {
		params = append(params, fmt.Sprintf("-I=%s", d))
	}
	pkgPath := pkgpath.FindPkgDir(cmd.PackagePath, ".")
	if len(pkgPath) >= 0 {
		params = append(params, fmt.Sprintf("-I=%s", filepath.Join(pkgPath, "/tools/protoc/include/")))
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
	execProtoc(files, dirs,
		fmt.Sprintf("--go_out=%s", output), "--go_opt=paths=source_relative",
	)
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
		srvDir := ensureOutputDir(command, "service_out", "Service")
		execProtoc(files, dirs,
			fmt.Sprintf("--go-register_out=%s", srvDir), "--go-register_opt=paths=source_relative",
			"--go-register_opt=grpc="+strconv.FormatBool(createGRPC), "--go-register_opt=http="+strconv.FormatBool(createHTTP),
		)
	}
}

func createGRPCService(command *cobra.Command, args, files, dirs []string) {
	clientDir := ensureOutputDir(command, "client_out", "Client")
	grpcDir := ensureOutputDir(command, "service_out", "Service")
	execProtoc(files, dirs,
		fmt.Sprintf("--go-grpc_out=%s", grpcDir), "--go-grpc_opt=paths=source_relative",
		fmt.Sprintf("--go-client_out=%s", clientDir), "--go-client_opt=paths=source_relative",
	)
}

func createHTTPService(command *cobra.Command, args, files, dirs []string) {
	msgDir := ensureOutputDir(command, "msg_out", "Message")
	httpDir := ensureOutputDir(command, "service_out", "Service")
	execProtoc(files, dirs,
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
	execProtoc(files, dirs,
		fmt.Sprintf("--go-provider_out=%s", impDir), "--go-provider_opt=paths=source_relative",
		"--go-provider_opt=grpc="+strconv.FormatBool(createGRPC), "--go-provider_opt=http="+strconv.FormatBool(createHTTP),
	)
}
