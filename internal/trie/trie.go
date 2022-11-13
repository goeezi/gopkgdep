package trie

import (
	"fmt"
	"io"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/goeezi/gopkgdep/internal/pathutil"
	"github.com/goeezi/gopkgdep/internal/set"
)

type Trie map[string]Trie

func (t Trie) Write(w io.Writer, punctFormat string, depth int) error {
	if len(t) == 0 {
		return nil
	}
	if depth > 0 {
		fmt.Fprintf(w, punctFormat, "{")
	}
	i := 0
	keys := maps.Keys(t)
	slices.Sort(keys)
	for _, k := range keys {
		if i > 0 {
			if depth > 0 {
				fmt.Fprintf(w, punctFormat, ",")
			} else {
				if _, err := w.Write([]byte{' '}); err != nil {
					return err
				}
			}
		}
		i++
		if _, err := w.Write([]byte(k)); err != nil {
			return err
		}
		if err := t[k].Write(w, punctFormat, depth+1); err != nil {
			return err
		}
	}
	if depth > 0 {
		fmt.Fprintf(w, punctFormat, "}")
	}
	return nil
}

func Build(paths set.Set[string]) Trie {
	if len(paths) == 1 {
		for s := range paths {
			return Trie{s: {}}
		}
	}
	tree := map[string]set.Set[string]{}
	for s := range paths {
		path := pathutil.Split(s)
		head, tail := "", []string{}
		if len(s) > 0 {
			head, tail = path[0], path[1:]
		}
		c, has := tree[head]
		if !has {
			c = set.Set[string]{}
			tree[head] = c
		}
		c.Add(pathutil.Join(tail))
	}

	result := Trie{}
	for p, s := range tree {
		t := Build(s)
		if len(t) == 1 {
			for k, v := range t {
				result[pathutil.Join([]string{p, k})] = v
			}
		} else {
			result[p] = t
		}
	}
	return result
}
