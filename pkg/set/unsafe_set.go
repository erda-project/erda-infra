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

package set

type set map[interface{}]struct{}

func (set set) Add(element interface{}) bool {
	if _, ok := set[element]; ok {
		return false
	}
	set[element] = struct{}{}
	return true
}

func (set set) Remove(element interface{}) {
	delete(set, element)
}

func (set set) Contains(elements ...interface{}) bool {
	for _, e := range elements {
		if _, ok := (set)[e]; !ok {
			return false
		}
	}
	return true
}

func (set set) Clear() {
	for k := range set {
		delete(set, k)
	}
}

func (set set) Len() int {
	return len(set)
}

func newSet() set {
	return make(set)
}
