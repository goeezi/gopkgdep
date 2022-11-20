// Copyright 2022 Marcelo Cantos, Melbourne, Australia
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package set

import "golang.org/x/exp/maps"

type Set[T comparable] map[T]struct{}

func New[T comparable](args ...T) Set[T] {
	s := Set[T]{}
	for _, a := range args {
		s.Add(a)
	}
	return s
}

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
