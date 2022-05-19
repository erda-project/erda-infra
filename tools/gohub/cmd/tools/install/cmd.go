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

package install

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/erda-project/erda-infra/tools/gohub/cmd"
	"github.com/erda-project/erda-infra/tools/gohub/cmd/tools"
)

var localInstall *bool
var ghProxy *string
var goProxy *string

func init() {
	localInstall = installCmd.Flags().Bool("local", false, "find local package")
	ghProxy = installCmd.Flags().String("ghproxy", "", "github proxy")
	goProxy = installCmd.Flags().String("goproxy", "", "go proxy")
	tools.AddCommand(installCmd)
}

func wrapGhProxy(url string) string {
	if ghProxy == nil || len(*ghProxy) == 0 {
		return url
	}
	*ghProxy = strings.TrimRight(*ghProxy, "/") + "/"
	return *ghProxy + url
}

func getGoProxyEnv() string {
	if goProxy == nil {
		*goProxy = ""
	}
	return "GOPROXY=" + *goProxy
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install all tools to {home directory}/." + cmd.Name,
	Run: func(command *cobra.Command, args []string) {
		Download(true, cmd.Verbose())
	},
}
