package util

type Set[T comparable] struct {
	m map[T]interface{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		m: make(map[T]interface{}),
	}
}

func NewSetFromItr[T comparable](arr ...T) *Set[T] {
	s := NewSet[T]()
	for _, item := range arr {
		s.Add(item)
	}
	return s
}

func (s *Set[T]) Add(item T) {
	s.m[item] = nil
}

func (s *Set[T]) Remove(item T) {
	delete(s.m, item)
}

func (s *Set[T]) Contains(item T) bool {
	_, ok := s.m[item]
	return ok
}
