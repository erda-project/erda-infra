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

import "sync"

// syncSet syncSet data structure, contains sync.RWMutex
type syncSet struct {
	set set
	sync.RWMutex
}

func (ss *syncSet) Add(element interface{}) bool {
	ss.Lock()
	defer ss.Unlock()
	result := ss.set.Add(element)
	return result
}

func (ss *syncSet) Remove(element interface{}) {
	ss.Lock()
	defer ss.Unlock()
	ss.set.Remove(element)
}

func (ss *syncSet) Contains(elements ...interface{}) bool {
	ss.RLock()
	defer ss.RUnlock()
	result := ss.set.Contains(elements...)
	return result
}

func (ss *syncSet) Clear() {
	ss.Lock()
	defer ss.Unlock()
	ss.set = newSet()
}

func (ss *syncSet) Len() int {
	ss.RLock()
	defer ss.RUnlock()
	result := ss.set.Len()
	return result
}

func newSyncSet() syncSet {
	return syncSet{set: newSet()}
}
