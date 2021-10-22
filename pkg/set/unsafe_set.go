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
