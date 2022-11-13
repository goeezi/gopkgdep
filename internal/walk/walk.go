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

package walk

import (
	"go/build"
	"path"
	"strings"

	"github.com/goeezi/gopkgdep/internal/graph"
	"github.com/goeezi/gopkgdep/internal/match"
	"github.com/goeezi/gopkgdep/internal/set"
)

type Walker struct {
	Import  func(path string) (*build.Package, error)
	Closed  bool
	Graph   *graph.Graph
	Matcher *match.Matcher
}

func (w *Walker) Walk(
	pkgs set.Set[string],
) (set.Set[string], error) {
	filtered := set.Set[string]{}
	for p := range pkgs {
		if p == "C" {
			// C isn't really a package.
			w.Graph.Pkgs["C"] = nil
			continue
		}
		p := w.Matcher.Resolve(p)
		match := w.Matcher.Match(p)
		p = w.Matcher.Rel(p)
		filtered.Add(p)
		if _, ok := w.Graph.Pkgs[p]; ok {
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

		deps := set.Set[string]{}
		// log.Printf("CANDIDATES FOR %v:", p)
		candidates, err := w.candidates(w.Matcher, match, pkg.Imports)
		if err != nil {
			return nil, err
		}
		for _, imp := range candidates {
			deps.Add(w.Matcher.Rel(w.Matcher.Resolve(imp)))
		}
		deps, err = w.Walk(deps)
		if err != nil {
			return nil, err
		}
		w.Graph.Pkgs[p] = deps
	}
	return filtered, nil
}

func (w *Walker) candidates(
	m *match.Matcher,
	match bool,
	deps []string,
) ([]string, error) {
	switch {
	case match && !w.Closed:
		// log.Printf("  MATCH && OPEN: %v", deps)
		var c []string
		for _, dep := range deps {
			if m.MatchUnlessExcludedStdLibPackage(dep) {
				c = append(c, dep)
			}
		}
		return c, nil
	case match || !w.Closed:
		// log.Printf("  MATCH || OPEN: %v", deps)
		var c []string
		for _, dep := range deps {
			dep := m.Resolve(dep)
			if m.Match(dep) {
				// log.Printf("    MATCHED: %v", dep)
				c = append(c, dep)
			}
		}
		return c, nil
	default:
		return nil, nil
	}
}
