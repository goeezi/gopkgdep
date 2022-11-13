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
