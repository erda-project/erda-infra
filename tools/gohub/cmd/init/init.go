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

package init

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/erda-project/erda-infra/tools/gohub/cmd"
	"github.com/spf13/cobra"
)

func init() {
	initCmd.Flags().StringP("template", "t", "simple", "full|simple")
	initCmd.Flags().StringP("out", "o", ".", "output directory")
	cmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a provider with name",
	Run: func(command *cobra.Command, args []string) {
		var provider string
		if len(args) <= 0 {
			fmt.Printf("Input Service Provider Name: ")
			fmt.Scan(&provider)
		} else if len(args) == 1 {
			provider = args[0]
		} else {
			command.Usage()
			os.Exit(1)
		}
		out, err := command.Flags().GetString("out")
		cmd.CheckError(err)
		if !cmd.IsFileExist(out) {
			err := os.MkdirAll(out, os.ModePerm)
			cmd.CheckError(err)
		}

		tname, err := command.Flags().GetString("template")
		cmd.CheckError(err)
		t, ok := templates[tname]
		if !ok {
			cmd.CheckError(fmt.Errorf("not found template %q", tname))
		}
		tmpl := template.New("provider").Funcs(template.FuncMap{
			"quote":  strconv.Quote,
			"printf": fmt.Sprintf,
		})
		ctx := &tempContext{
			Package:  getPackageName(),
			Provider: provider,
		}

		createFile(tmpl, ctx, filepath.Join(out, "provider.go"), t.Content)
		createFile(tmpl, ctx, filepath.Join(out, "provider_test.go"), t.TestContext)
	},
}

func createFile(tmpl *template.Template, ctx *tempContext, filename, content string) {
	tmpl, err := tmpl.Parse(content)
	cmd.CheckError(err)
	_, err = os.Stat(filename)
	if err == nil {
		cmd.CheckError(fmt.Errorf("%s file already exist", filename))
	}
	file, err := os.Create(filename)
	cmd.CheckError(err)
	defer file.Close()
	err = tmpl.Execute(file, ctx)
	cmd.CheckError(err)
}

func getPackageName() string {
	wd, err := os.Getwd()
	cmd.CheckError(err)
	base := filepath.Base(wd)
	idx := strings.LastIndex(base, "-")
	if idx >= 0 {
		base = base[idx+1:]
	}
	sb := strings.Builder{}
	for _, c := range base {
		if unicode.IsLetter(c) {
			sb.WriteRune(c)
		}
	}
	return sb.String()
}
