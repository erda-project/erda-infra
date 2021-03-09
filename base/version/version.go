// Author: recallsong
// Email: songruiguo@qq.com

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
