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

package safe

import (
	"testing"
)

func TestDo(t *testing.T) {
	panicFunc := func() { panic("panic now") }

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic should be covered already, err: %v", r)
		}
	}()
	Do(panicFunc)
}

func TestGo(t *testing.T) {
	done := make(chan bool)
	panicFunc := func() {
		defer func() { done <- true }()
		panic("panic now")
	}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("panic should be covered already, err: %v", r)
		}
	}()
	Go(panicFunc)
	<-done
}
