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

package version

import (
	"fmt"
	"os"
)

var (
	// Version .
	Version string
	// BuildTime .
	BuildTime string
	// GoVersion .
	GoVersion string
	// CommitID git commot id
	CommitID string
	// DockerImage docker image url
	DockerImage string
)

// String 返回版本信息
func String() string {
	return fmt.Sprintf("Version: %s\nBuildTime: %s\nGoVersion: %s\nCommitID: %s\nDockerImage: %s\n",
		Version, BuildTime, GoVersion, CommitID, DockerImage)
}

// Print print version information
func Print() {
	fmt.Print(String())
}

// PrintIfCommand .
func PrintIfCommand() {
	if len(os.Args) == 2 && os.Args[1] == "version" {
		Print()
		os.Exit(0)
	}
}
