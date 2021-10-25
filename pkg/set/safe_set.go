package set

import "sync"

type syncSet struct {
	set set
	sync.RWMutex
}

func (ss *syncSet) Add(element interface{}) bool {
	ss.Lock()
	result := ss.set.Add(element)
	ss.Unlock()
	return result
}

func (ss *syncSet) Remove(element interface{}) {
	ss.Lock()
	ss.set.Remove(element)
	ss.Unlock()
}

func (ss *syncSet) Contains(elements ...interface{}) bool {
	ss.RLock()
	result := ss.set.Contains(elements...)
	ss.RUnlock()
	return result
}

func (ss *syncSet) Clear() {
	ss.Lock()
	ss.set = newSet()
	ss.Unlock()
}

func (ss *syncSet) Len() int {
	ss.RLock()
	result := ss.set.Len()
	ss.RUnlock()
	return result
}

func newSyncSet() syncSet {
	return syncSet{set: newSet()}
}
