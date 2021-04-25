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

package uninstall

import (
	"os"
	"path/filepath"

	"github.com/erda-project/erda-infra/tools/gohub/cmd"
	"github.com/erda-project/erda-infra/tools/gohub/cmd/tools"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

func init() {
	tools.AddCommand(uninstallCmd)
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall tools from {home directory}/." + cmd.Name,
	Run: func(command *cobra.Command, args []string) {
		home := homeDir()
		dir := filepath.Join(home, "."+cmd.Name)
		if cmd.IsFileExist(dir) {
			err := os.RemoveAll(dir)
			cmd.CheckError(err)
		}
	},
}

func homeDir() string {
	home, err := homedir.Dir()
	cmd.CheckError(err)
	return home
}
