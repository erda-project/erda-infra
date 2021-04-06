// Author: recallsong
// Email: songruiguo@qq.com

package main

import (
	"fmt"
	"go/build"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		return
	}
	dir := FindPkgDir(os.Args[1], ".")
	fmt.Println(dir)
}

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
