package util

type Set struct {
	m map[interface{}]interface{}
}

func NewSet() *Set {
	return &Set{
		m: make(map[interface{}]interface{}),
	}
}

func (s *Set) Add(item interface{}) {
	s.m[item] = nil
}

func (s *Set) Remove(item interface{}) {
	delete(s.m, item)
}

func (s *Set) Contains(item interface{}) bool {
	_, ok := s.m[item]
	return ok
}
