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

// Set defined set func, Set is used to store a collection of unique elements.
type Set interface {

	// Add element to set
	Add(element interface{}) bool
	// Remove remove element to set
	Remove(element interface{})
	// Contains Judge whether the set contains the specified elements, if anyone does not contain, return false, otherwise return true
	Contains(elements ...interface{}) bool
	// Clear remove all from set
	Clear()
	// Len return the number of element
	Len() int
}

// NewSet Initialize a set
func NewSet(e ...interface{}) Set {
	set := newSet()
	for _, item := range e {
		set.Add(item)
	}
	return &set
}

func NewSyncSet(e ...interface{}) Set {
	set := newSyncSet()
	for _, item := range e {
		set.Add(item)
	}
	return &set
}
