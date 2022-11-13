package set

import "golang.org/x/exp/maps"

type Set[T comparable] map[T]struct{}

func (s Set[T]) Has(t T) bool {
	_, has := s[t]
	return has
}

func (s Set[T]) Add(t T) {
	s[t] = struct{}{}
}

func (s Set[T]) Delete(t T) {
	delete(s, t)
}

func (s Set[T]) Elements() []T {
	return maps.Keys(s)
}
