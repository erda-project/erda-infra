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

package pkgpath

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"

	"github.com/erda-project/erda-infra/tools/gohub/cmd"
	"github.com/spf13/cobra"
)

func FindPkgDir(path, srcDir string) string {
	if path == "" {
		return ""
	}
	// Don't require the source files to be present.
	if abs, err := filepath.Abs(srcDir); err == nil {
		srcDir = abs
	}
	bp, _ := build.Import(path, srcDir, build.FindOnly|build.AllowBinary)
	return bp.Dir
}

func init() {
	cmd.AddCommand(gopkgCmd)
}

var gopkgCmd = &cobra.Command{
	Use:     "pkgpath [package]",
	Aliases: []string{"pkg", "package", "gopkg"},
	Short:   "Print the absolute path of go package",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			cmd.Usage()
			os.Exit(1)
		}
		for _, pkg := range args {
			dir := FindPkgDir(pkg, ".")
			if len(dir) <= 0 {
				fmt.Fprintf(os.Stderr, "not found path of package %q\n", pkg)
				os.Exit(1)
			}
			fmt.Println(dir)
		}
	},
}
