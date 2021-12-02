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

package hook

import (
	"log"
	"reflect"
	"runtime"
	"strings"

	"github.com/brahma-adshonor/gohook"
)

// Hook .
func Hook(target, replacement, trampoline interface{}) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[ERROR] failed to hook %T : %v\n", target, err)
		}
	}()
	err := gohook.Hook(target, replacement, trampoline)
	if err == nil {
		funcDecl := reflect.TypeOf(target).String()
		funcDecl = strings.TrimPrefix(funcDecl, "func")
		log.Printf("hook func %s%s", getFunctionName(target), funcDecl)
	}
	return err
}

func getFunctionName(i interface{}, seps ...rune) string {
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	return fn
}
