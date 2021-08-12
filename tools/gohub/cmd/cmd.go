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

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// common value
const (
	Name        = "gohub"
	PackagePath = "github.com/erda-project/erda-infra"
)

// RootCmd .
var RootCmd = &cobra.Command{
	Use:   Name,
	Short: Name,
	Long: `The ` + Name + ` is CLI tools for erda-infra. 
Complete documentation is available at https://` + PackagePath + ` .`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
		os.Exit(1)
	},
}

var verbose = RootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose information")

// Verbose .
func Verbose() bool { return *verbose }

// AddCommand .
func AddCommand(cmd *cobra.Command) {
	RootCmd.AddCommand(cmd)
}

// CheckError .
func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// Confirm .
func Confirm(msg string) bool {
	var yes string
	for {
		fmt.Printf("%s [yes/no]: ", msg)
		fmt.Scan(&yes)
		yes = strings.ToLower(yes)
		switch yes {
		case "yes", "y":
			return true
		case "no", "n":
			return false
		}
	}
}

// IsFileExist .
func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		CheckError(err)
	}
	return true
}

// Execute .
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
