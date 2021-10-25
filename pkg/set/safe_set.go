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
