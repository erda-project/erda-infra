// Copyright 2021 Terminus
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
