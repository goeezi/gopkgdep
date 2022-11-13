package walk

import (
	"go/build"
	"log"
	"path"
	"strings"

	"github.com/goeezi/gopkgdep/internal/graph"
	"github.com/goeezi/gopkgdep/internal/match"
	"github.com/goeezi/gopkgdep/internal/set"
)

type Walker struct {
	Import func(path string) (*build.Package, error)
	Open   bool
}

func (w *Walker) Walk(
	g *graph.Graph,
	m *match.Matcher,
	pkgs set.Set[string],
) (set.Set[string], error) {
	filtered := set.Set[string]{}
	for p := range pkgs {
		if p == "C" {
			// C isn't really a package.
			g.Pkgs["C"] = nil
			continue
		}
		p := m.Resolve(p)
		match := m.Match(p)
		p, err := m.Rel(p)
		if err != nil {
			return nil, err
		}
		filtered.Add(p)
		if _, ok := g.Pkgs[p]; ok {
			// already seen
			continue
		}
		if strings.HasPrefix(p, "golang_org") {
			p = path.Join("vendor", p)
		}

		pkg, err := w.Import(p)
		if err != nil {
			return nil, err
		}
		log.Printf("%s => %#v", p, pkg.Imports)

		deps := set.Set[string]{}
		candidates, err := w.candidates(m, match, pkg.Imports)
		if err != nil {
			return nil, err
		}
		for _, imp := range candidates {
			dep, err := m.Rel(m.Resolve(imp))
			if err != nil {
				return nil, err
			}
			deps.Add(dep)
		}
		deps, err = w.Walk(g, m, deps)
		if err != nil {
			return nil, err
		}
		g.Pkgs[p] = deps
	}
	return filtered, nil
}

func (w *Walker) candidates(
	m *match.Matcher,
	match bool,
	deps []string,
) ([]string, error) {
	switch {
	case match && w.Open:
		var c []string
		for _, dep := range deps {
			if m.MatchUnlessExcludedStdLibPackage(dep) {
				c = append(c, dep)
			}
		}
		return c, nil
	case match || w.Open:
		var c []string
		for _, dep := range deps {
			dep := m.Resolve(dep)
			if m.Match(dep) {
				c = append(c, dep)
			}
		}
		return c, nil
	default:
		return nil, nil
	}
}
