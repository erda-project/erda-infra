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
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"

	"github.com/erda-project/erda-infra/tools/gohub/cmd"
	"github.com/erda-project/erda-infra/tools/gohub/cmd/pkgpath"
	"github.com/erda-project/erda-infra/tools/gohub/cmd/version"
)

// IncludeDirs .
func IncludeDirs() []string {
	home := homeDir()
	repo := filepath.Base(cmd.PackagePath)
	return []string{
		filepath.Join(home, "."+cmd.Name, repo, "tools/protoc/include"),
		filepath.Join(home, "."+cmd.Name),
	}
}

func getVersion() string {
	home := homeDir()
	file := filepath.Join(home, "."+cmd.Name, ".version")
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return ""
		}
		cmd.CheckError(err)
	}
	return strings.TrimSpace(string(bytes))
}

func updateVersion() {
	home := homeDir()
	file := filepath.Join(home, "."+cmd.Name, ".version")
	err := ioutil.WriteFile(file, []byte(version.Version), os.ModePerm)
	cmd.CheckError(err)
}

// Download .
func Download(override, verbose bool) {
	dir := ensureToolsDir()

	// check version
	if getVersion() != version.Version {
		override = true
	}

	// download protoc
	if !cmd.IsFileExist(filepath.Join(dir, "protoc")) || (!*localInstall && override) {
		file := downloadProtoc(dir, verbose)
		err := unzip(file, func(f *zip.File) (string, bool) {
			if f.Name == "bin/protoc" {
				return filepath.Join(dir, "protoc"), true
			}
			if strings.HasPrefix(f.Name, "include") {
				return filepath.Join(dir, f.Name), true
			}
			return "", false
		})
		cmd.CheckError(err)
		err = os.Remove(file)
		cmd.CheckError(err)
	}

	// install plugins
	for _, p := range []struct {
		Name string
		URL  string
		Path string
	}{
		{
			Name: "protoc-gen-go",
			URL:  "https://github.com/golang/protobuf",
			Path: "protoc-gen-go",
		},
		{
			Name: "protoc-gen-validate",
			URL:  "https://github.com/envoyproxy/protoc-gen-validate",
			Path: "protoc-gen-validate",
		},
	} {
		if !cmd.IsFileExist(filepath.Join(dir, p.Name)) || (!*localInstall && override) {
			u, err := url.Parse(p.URL)
			cmd.CheckError(err)
			host, _, err := net.SplitHostPort(u.Host)
			if err != nil {
				host = u.Host
			}
			repodir := filepath.Join(dir, host, u.Path)
			tmpdir := repodir + ".tmp"
			// create tmpdir
			err = os.RemoveAll(tmpdir)
			cmd.CheckError(err)
			err = os.MkdirAll(tmpdir, os.ModePerm)
			cmd.CheckError(err)
			// clone
			runCommand(dir, nil, "git", "clone", "--depth", "1", wrapGhProxy(p.URL), tmpdir)
			// rename
			err = os.RemoveAll(repodir)
			cmd.CheckError(err)
			err = os.Rename(tmpdir, repodir)
			cmd.CheckError(err)
			// build
			fmt.Printf("building %s ...\n", p.Name)
			buildDir := repodir
			goPath := os.Getenv("GOPATH")
			if p.Name == "protoc-gen-validate" {
				setEnvPathWithDir(dir)
				runCommand(buildDir, []string{fmt.Sprintf("GOBIN=%s", dir)}, "make", "build")
				err = os.RemoveAll(filepath.Join(dir, "include"))
				cmd.CheckError(err)
			} else {
				if len(p.Path) > 0 {
					buildDir = filepath.Join(repodir, p.Path)
				}
				src := filepath.Join(dir, p.Name)
				runCommand(buildDir, nil, "go", "build", "-o", src)
				runCommand(buildDir, nil, "cp", "-f", src, filepath.Join(goPath, "bin"))
			}
			fmt.Printf("build %s successfully !\n", p.Name)
		}
	}

	plugins := []string{
		"protoc-gen-go-grpc",
		"protoc-gen-go-client",
		"protoc-gen-go-http",
		"protoc-gen-go-form",
		"protoc-gen-go-json",
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
		repo := filepath.Base(cmd.PackagePath)
		repodir := filepath.Join(dir, repo)
		var pkgPath string
		if *localInstall {
			pkgPath = pkgpath.FindPkgDir(cmd.PackagePath+"/tools", ".")
			if len(pkgPath) <= 0 {
				cmd.CheckError(fmt.Errorf("not found package %q", cmd.PackagePath))
			}
			fmt.Printf("tools go package path: %s\n", pkgPath)
			destDir := filepath.Join(repodir, "tools/protoc/")
			err := os.MkdirAll(destDir, os.ModePerm)
			cmd.CheckError(err)
			copyDir(filepath.Join(pkgPath, "protoc/include"), destDir)
		} else {
			tmpdir := filepath.Join(dir, repo+".tmp")
			err := os.RemoveAll(tmpdir)
			cmd.CheckError(err)
			runCommand(dir, nil, "git", "clone", "--depth", "1", "https://"+cmd.PackagePath, tmpdir)
			err = os.RemoveAll(repodir)
			cmd.CheckError(err)
			err = os.Rename(tmpdir, repodir)
			cmd.CheckError(err)
			pkgPath = filepath.Join(repodir, "tools")
		}
		// build protoc plugins
		for _, plugin := range plugins {
			if !cmd.IsFileExist(filepath.Join(dir, plugin)) || override {
				fmt.Printf("building %s ...\n", plugin)
				buildDir := filepath.Join(pkgPath, "protoc", plugin)
				runCommand(buildDir, nil, "go", "build", "-o", filepath.Join(dir, plugin))
				fmt.Printf("build %s successfully !\n", plugin)
			}
		}
	}

	if override {
		updateVersion()
	}

	setEnvPathWithDir(dir)
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
		cmd.CheckError(fmt.Errorf("%s file already exist, it not a directory", dir))
	}
	return dir
}

func runCommand(wd string, envs []string, exe string, params ...string) {
	command := exec.Command(exe, params...)
	if cmd.Verbose() {
		fmt.Fprintf(os.Stdout, "\n>>> run command: %s\n", command.String())
	}
	command.Dir = wd
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Env = append(os.Environ(), getGoProxyEnv())
	command.Env = append(command.Env, envs...)
	err := command.Run()
	cmd.CheckError(err)
}

func downloadProtoc(dir string, verbose bool) string {
	var url string
	switch {
	case runtime.GOOS == "darwin" && runtime.GOARCH == "arm64":
		url = "https://erda-project.oss-cn-hangzhou.aliyuncs.com/protobuf/protoc-3.15.8-osx-x86_64.zip"
	case runtime.GOOS == "darwin" && runtime.GOARCH == "amd64":
		url = "https://erda-project.oss-cn-hangzhou.aliyuncs.com/protobuf/protoc-3.15.8-osx-x86_64.zip"
	case runtime.GOOS == "linux" && runtime.GOARCH == "amd64":
		url = "https://erda-project.oss-cn-hangzhou.aliyuncs.com/protobuf/protoc-3.15.8-linux-x86_64.zip"
	case runtime.GOOS == "linux" && runtime.GOARCH == "arm64":
		url = "https://erda-project.oss-cn-hangzhou.aliyuncs.com/protobuf/protoc-3.15.8-linux-aarch_64.zip"
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

func copyDir(src string, dest string) {
	formatPath := func(s string) string {
		switch runtime.GOOS {
		case "windows":
			return strings.Replace(s, "/", "\\", -1)
		case "darwin", "linux":
			return strings.Replace(s, "\\", "/", -1)
		default:
			return s
		}
	}
	src, dest = formatPath(src), formatPath(dest)
	switch runtime.GOOS {
	case "windows":
		runCommand("", nil, "xcopy", src, dest, "/I", "/E")
	case "darwin", "linux":
		runCommand("", nil, "cp", "-R", src, dest)
	}
}

func joinPathList(list ...string) {
	sep := string(os.PathListSeparator)
	paths := os.Getenv("PATH")
	set := make(map[string]bool)
	for _, p := range strings.Split(paths, sep) {
		set[filepath.Clean(p)] = true
	}
	for i := len(list) - 1; i >= 0; i-- {
		p := filepath.Clean(list[i])
		if !set[p] {
			set[p] = true
			paths = p + sep + paths
		}
	}
	os.Setenv("PATH", paths)
}

func setEnvPathWithDir(dir string) {
	paths := []string{dir}
	goPath := os.Getenv("GOPATH")
	if len(goPath) > 0 {
		for _, p := range strings.Split(goPath, string(os.PathListSeparator)) {
			paths = append(paths, filepath.Join(p, "bin"))
		}
	}
	joinPathList(paths...)
}
