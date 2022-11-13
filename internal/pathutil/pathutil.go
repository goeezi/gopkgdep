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

package pathutil

import "strings"

func Diff[T comparable](a, b []T) (head, atail, btail []T) {
	n := len(a)
	if n > len(b) {
		n = len(b)
	}
	i := 0
	for ; i < n; i++ {
		if a[i] != b[i] {
			break
		}
	}
	return a[:i], a[i:], b[i:]
}

func Join(s []string) string {
	return strings.Join(s, "")
}

func Split(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.SplitAfter(s, "/")

	// Fuse ["./", "foo/", ...] into ["./foo/", ...]
	if len(parts) > 1 && parts[0] == "./" {
		parts[1] = parts[0] + parts[1]
		parts = parts[1:]
	}

	return parts
}

func Last[T any](t []T) T {
	return t[len(t)-1]
}

func Trim[T any](t []T) []T {
	if len(t) == 0 {
		return t
	}
	return t[:len(t)-1]
}
