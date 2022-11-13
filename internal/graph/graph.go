package graph

import "github.com/goeezi/gopkgdep/internal/set"

type Graph struct {
	Pkgs   map[string]set.Set[string]
	Depths map[string]int
}

func New() *Graph {
	return &Graph{
		Pkgs:   make(map[string]set.Set[string]),
		Depths: make(map[string]int),
	}
}

func (g *Graph) Depth(pkg string) int {
	if d, has := g.Depths[pkg]; has {
		return d
	}
	deps := g.Pkgs[pkg]
	if len(deps) == 0 {
		return 0
	}
	maxdepth := 0
	for dep := range g.Pkgs[pkg] {
		depth := g.Depth(dep)
		if maxdepth < depth {
			maxdepth = depth
		}
	}
	maxdepth++
	g.Depths[pkg] = maxdepth
	return maxdepth
}

func (g *Graph) Less(a, b string) bool {
	deptha := g.Depth(a)
	depthb := g.Depth(b)
	if deptha != depthb {
		return deptha < depthb
	}
	return a < b
}
