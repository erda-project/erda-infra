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

// set Set data structure
type set map[interface{}]struct{}

func (s *set) Add(element interface{}) bool {
	if _, ok := (*s)[element]; ok {
		return false
	}
	(*s)[element] = struct{}{}
	return true
}

func (s *set) Remove(element interface{}) {
	delete(*s, element)
}

func (s *set) Contains(elements ...interface{}) bool {
	for _, e := range elements {
		if _, ok := (*s)[e]; !ok {
			return false
		}
	}
	return true
}

func (s *set) Clear() {
	*s = newSet()
}

func (s *set) Len() int {
	return len(*s)
}

func newSet() set {
	return make(set)
}
