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

package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/magiconair/properties"
	"github.com/mitchellh/mapstructure"
	"github.com/pelletier/go-toml"
	"github.com/recallsong/go-utils/reflectx"
	"gopkg.in/ini.v1"
	"gopkg.in/yaml.v2"
)

// TrimBOM .
func TrimBOM(f []byte) []byte {
	return bytes.TrimPrefix(f, []byte("\xef\xbb\xbf"))
}

var (
	envVarRe      = regexp.MustCompile(`\$([\w\|]+)|\$\{([\w\|]+)(:[^}]*)?\}`)
	envVarEscaper = strings.NewReplacer(
		`"`, `\"`,
		`\`, `\\`,
	)
)

// EscapeEnv .
func EscapeEnv(contents []byte) []byte {
	params := envVarRe.FindAllSubmatch(contents, -1)
	for _, param := range params {
		if len(param) != 4 {
			continue
		}
		var key, defval []byte
		if len(param[1]) > 0 {
			key = param[1]
		} else if len(param[2]) > 0 {
			key = param[2]
		} else {
			continue
		}
		if len(param[3]) > 0 {
			defval = param[3][1:]
		}
		envKey := strings.TrimPrefix(reflectx.BytesToString(key), "$")
		val, ok := os.LookupEnv(envKey)
		if !ok && strings.Contains(envKey, "|") {
			val, ok = lookupEnvWithBooleanExpression(envKey)
		}
		if !ok {
			if len(param[1]) > 0 {
				continue
			}
			val = string(defval)
		}
		val = envVarEscaper.Replace(val)
		contents = bytes.Replace(contents, param[0], reflectx.StringToBytes(val), 1)
	}
	return contents
}

func lookupEnvWithBooleanExpression(envKey string) (string, bool) {
	keys := strings.Split(envKey, "|")
	if len(keys) == 1 {
		// skip parse boolean value if doesn't contains "|"
		return os.LookupEnv(envKey)
	}

	var vals []string
	var allBool = true
	var firstTrueVal string

	for _, key := range keys {
		val, ok := os.LookupEnv(key)
		if !ok {
			continue
		}

		vals = append(vals, val)
		b, err := strconv.ParseBool(val)
		if err != nil {
			allBool = false
			break
		}
		if b && len(firstTrueVal) == 0 {
			firstTrueVal = val
		}
	}

	if len(vals) == 0 {
		return "", false
	}

	if allBool && len(firstTrueVal) > 0 {
		return firstTrueVal, true
	}
	return vals[0], true
}

// ParseError denotes failing to parse configuration file.
type ParseError struct {
	err error
}

// Error returns the formatted configuration error.
func (pe ParseError) Error() string {
	return fmt.Sprintf("While parsing config: %s", pe.err.Error())
}

// UnmarshalToMap .
func UnmarshalToMap(in io.Reader, typ string, c map[string]interface{}) (err error) {
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(in)
	if err != nil {
		return err
	}
	err = polishBuffer(buf)
	if err != nil {
		return err
	}
	switch strings.ToLower(typ) {
	case "yaml", "yml":
		if err = yaml.Unmarshal(buf.Bytes(), &c); err != nil {
			return ParseError{err}
		}
	case "json":
		if err = json.Unmarshal(buf.Bytes(), &c); err != nil {
			return ParseError{err}
		}
	case "hcl":
		obj, err := hcl.Parse(reflectx.BytesToString(buf.Bytes()))
		if err != nil {
			return ParseError{err}
		}
		if err = hcl.DecodeObject(&c, obj); err != nil {
			return ParseError{err}
		}
	case "toml":
		tree, err := toml.LoadReader(buf)
		if err != nil {
			return ParseError{err}
		}
		tmap := tree.ToMap()
		for k, v := range tmap {
			c[k] = v
		}
	case "properties", "props", "prop":
		props := properties.NewProperties()
		var err error
		if err = props.Load(buf.Bytes(), properties.UTF8); err != nil {
			return ParseError{err}
		}
		for _, key := range props.Keys() {
			value, _ := props.Get(key)
			// recursively build nested maps
			path := strings.Split(key, ".")
			lastKey := strings.ToLower(path[len(path)-1])
			deepestMap := deepSearch(c, path[0:len(path)-1])
			// set innermost value
			deepestMap[lastKey] = value
		}
	case "ini":
		cfg := ini.Empty()
		err = cfg.Append(buf.Bytes())
		if err != nil {
			return ParseError{err}
		}
		sections := cfg.Sections()
		for i := 0; i < len(sections); i++ {
			section := sections[i]
			keys := section.Keys()
			for j := 0; j < len(keys); j++ {
				key := keys[j]
				value := cfg.Section(section.Name()).Key(key.Name()).String()
				c[section.Name()+"."+key.Name()] = value
			}
		}
	}
	toStringKeyMap(c)
	return nil
}

func polishBuffer(buf *bytes.Buffer) error {
	byts := buf.Bytes()
	byts = TrimBOM(byts)
	byts = EscapeEnv(byts)
	buf.Reset()
	_, err := buf.Write(byts)
	return err
}

func toStringKeyMap(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v := range x {
			m[fmt.Sprint(k)] = toStringKeyMap(v)
		}
		return m
	case map[string]interface{}:
		for k, v := range x {
			x[k] = toStringKeyMap(v)
		}
	case []interface{}:
		for i, v := range x {
			x[i] = toStringKeyMap(v)
		}
	}
	return i
}

// deepSearch scans deep maps, following the key indexes listed in the
// sequence "path".
// The last value is expected to be another map, and is returned.
//
// In case intermediate keys do not exist, or map to a non-map value,
// a new map is created and inserted, and the search continues from there:
// the initial map "m" may be modified!
func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			// intermediate key does not exist
			// => create it and continue from there
			m3 := make(map[string]interface{})
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]interface{})
		if !ok {
			// intermediate key is a value
			// => replace with a new map
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		// continue search from here
		m = m3
	}
	return m
}

// ConvertData .
func ConvertData(input, output interface{}, tag string) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           output,
		WeaklyTypedInput: true,
		TagName:          tag,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			mapstructure.StringToTimeHookFunc("2006-01-02 15:04:05"),
		),
	})
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}

// LoadFile .
func LoadFile(path string) ([]byte, error) {
	byts, err := ioutil.ReadFile(path)
	return byts, err
}

// LoadToMap .
func LoadToMap(path string, c map[string]interface{}) error {
	typ := filepath.Ext(path)
	if len(typ) <= 0 {
		return fmt.Errorf("%s unknown file extension", path)
	}
	byts, err := LoadFile(path)
	if err != nil {
		return err
	}
	return UnmarshalToMap(bytes.NewReader(byts), typ[1:], c)
}

// LoadEnvFileWithPath .
func LoadEnvFileWithPath(path string, override bool) {
	byts, err := ioutil.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		return
	}
	regex := regexp.MustCompile(`\s+\#`)
	content := reflectx.BytesToString(byts)
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, "#") {
			continue
		}
		loc := regex.FindIndex(reflectx.StringToBytes(line))
		if len(loc) > 0 {
			line = line[0:loc[0]]
		}
		idx := strings.Index(line, "=")
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		if len(key) <= 0 {
			continue
		}
		val := strings.TrimSpace(line[idx+1:])
		if override {
			os.Setenv(key, val)
		} else {
			_, ok := os.LookupEnv(key)
			if !ok {
				os.Setenv(key, val)
			}
		}

	}
}

// LoadEnvFile .
func LoadEnvFile(profile string) {
	wd, err := os.Getwd()
	if err != nil {
		return
	}
	path := filepath.Join(wd, ".env")
	LoadEnvFileWithPath(path, false)
	LoadEnvFileByProfile(path, profile)
}

// LoadEnvFileByProfile load env variables by profile file
func LoadEnvFileByProfile(path, profile string) {
	if profile != "" {
		path = fmt.Sprintf("%s-%s", path, profile)
		LoadEnvFileWithPath(path, true)
	}
}
