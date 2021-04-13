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
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/erda-project/erda-infra/tools/gohub/cmd"
	"github.com/erda-project/erda-infra/tools/gohub/cmd/pkgpath"
	"github.com/mitchellh/go-homedir"
)

// Download .
func Download(override, verbose bool) {
	dir := ensureToolsDir()

	// download protoc
	if !cmd.IsFileExist(filepath.Join(dir, "protoc")) || override {
		file := downloadProtoc(dir, verbose)
		err := unzip(file, func(f *zip.File) (string, bool) {
			if f.Name == "bin/protoc" {
				return filepath.Join(dir, "protoc"), true
			}
			return "", false
		})
		cmd.CheckError(err)
		err = os.Remove(file)
		cmd.CheckError(err)
	}

	plugins := []string{
		"protoc-gen-go-grpc",
		"protoc-gen-go-client",
		"protoc-gen-go-http",
		"protoc-gen-go-form",
		"protoc-gen-go-register",
		"protoc-gen-go-provider",
	}

	var checkPlugins bool
	for _, plugin := range plugins {
		if !cmd.IsFileExist(filepath.Join(dir, plugin)) {
			checkPlugins = true
		}
	}

	if checkPlugins || override {
		// find package path
		pkgPath := pkgpath.FindPkgDir(cmd.PackagePath, ".")
		if len(pkgPath) <= 0 {
			command := exec.Command("go", "get", cmd.PackagePath)
			err := command.Run()
			cmd.CheckError(err)
		}
		pkgPath = pkgpath.FindPkgDir(cmd.PackagePath, ".")
		if len(pkgPath) <= 0 {
			cmd.CheckError(fmt.Errorf("not found package %q", pkgPath))
		}
		fmt.Printf("tools go package path: %s\n", pkgPath)

		// build protoc plugins
		for _, plugin := range plugins {
			if !cmd.IsFileExist(filepath.Join(dir, plugin)) || override {
				fmt.Printf("building %s ...\n", plugin)
				command := exec.Command("go", "build", "-o", filepath.Join(dir, plugin))
				command.Dir = filepath.Join(pkgPath, "tools", "protoc", plugin)
				err := command.Run()
				cmd.CheckError(err)
				fmt.Printf("build %s successfully !\n", plugin)
			}
		}
	}
	paths := os.Getenv("PATH")
	os.Setenv("PATH", dir+string(os.PathListSeparator)+paths)
}

func ensureToolsDir() string {
	home := homeDir()
	dir := filepath.Join(home, "."+cmd.Name)
	stat, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
			cmd.CheckError(err)
			return dir
		}
		cmd.CheckError(err)
	}
	if !stat.IsDir() {
		cmd.CheckError(fmt.Errorf("%s file already exist, it not a directory.", dir))
	}
	return dir
}

func downloadProtoc(dir string, verbose bool) string {
	var url string
	switch {
	case runtime.GOOS == "darwin" && runtime.GOARCH == "amd64":
		url = "https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-osx-x86_64.zip"
	case runtime.GOOS == "linux" && runtime.GOARCH == "amd64":
		url = "https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip"
	default:
		cmd.CheckError(fmt.Errorf("not support %s-%s environment", runtime.GOARCH, runtime.GOOS))
	}
	idx := strings.LastIndex(url, "/")
	if idx <= 0 && !strings.HasSuffix(url, ".zip") {
		cmd.CheckError(fmt.Errorf("invaid url %q", url))
	}
	filename := url[idx+1:]
	path := filepath.Join(dir, filename)
	if verbose {
		fmt.Printf("downloading %s to %s ...\n", url, path)
	} else {
		fmt.Printf("downloading %s ...\n", url)
	}
	res, err := http.Get(url)
	cmd.CheckError(err)
	f, err := os.Create(path)
	cmd.CheckError(err)
	io.Copy(f, res.Body)
	fmt.Printf("download %s successfully !\n", path)
	return path
}

func unzip(zipFile string, filter func(*zip.File) (string, bool)) error {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, f := range reader.File {
		fpath, ok := filter(f)
		if !ok {
			continue
		}
		if f.FileInfo().IsDir() {
			err = os.MkdirAll(fpath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func homeDir() string {
	home, err := homedir.Dir()
	cmd.CheckError(err)
	return home
}
