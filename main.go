package main

import (
	"flag"
	"go/build"
	"log"
	"os"
	"regexp"

	"golang.org/x/mod/modfile"

	"github.com/goeezi/gopkgdep/internal/graph"
	"github.com/goeezi/gopkgdep/internal/match"
	"github.com/goeezi/gopkgdep/internal/pathutil"
	"github.com/goeezi/gopkgdep/internal/render/dot"
	"github.com/goeezi/gopkgdep/internal/render/edges/edges"
	"github.com/goeezi/gopkgdep/internal/render/tsort"
	"github.com/goeezi/gopkgdep/internal/walk"
)

func main() {
	// Filters
	incl := flag.String("include", ".*",
		"include packages whose fully-qualified import path matches pattern")
	excl := flag.String("exclude", "",
		"exclude packages whose fully-qualified import path matches pattern "+
			"(takes priority over -include and -stdlib)")
	closed := flag.Bool("closed", false,
		"exclude neighbors of matching packages (takes priority over -include and -stdlib)")
	stdlib := flag.Bool("stdlib", false, "don't exclude standard library packages")

	// Output: tsort
	outtsort := flag.Bool("depthsort", false, "sort by graph depth")
	compact := flag.Bool("compact", false, "radix-compresses dependency list (implies -depthsort)")

	// Output: dot
	outdot := flag.Bool("dot", false, "output graphviz (graphviz must be installed)")
	strata := flag.Bool("strata", false, "layer graphviz nodes according to graph depth (implies -dot)")

	flag.Parse()

	*outtsort = *outtsort || *compact
	*outdot = *outdot || *strata

	if *outtsort && *outdot {
		log.Fatal("-dsort and -dot cannot be used together")
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	c := build.Default
	c.Dir = wd

	module, err := goModule()
	if err != nil {
		log.Fatal(err)
	}

	inclRE, err := regexp.Compile(*incl)
	if err != nil {
		log.Fatal(err)
	}

	var exclRE *regexp.Regexp
	if *excl != "" {
		var err error
		exclRE, err = regexp.Compile(*excl)
		if err != nil {
			log.Fatal(err)
		}
	}

	g := graph.New()

	m := &match.Matcher{
		InclRE: inclRE,
		ExclRE: exclRE,
		Module: module,
		Stdlib: *stdlib,
	}

	w := &walk.Walker{
		Import: func(path string) (*build.Package, error) {
			return c.Import(path, wd, 0)
		},
		Open: !*closed,
	}

	paths, err := pathutil.ResolvePaths(flag.Args())
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Walk(g, m, paths)
	if err != nil {
		log.Fatal(err)
	}

	focus := func(pkg string) int {
		switch {
		case paths.Has(pkg):
			return 2
		case m.Match(m.Resolve(pkg)):
			return 1
		default:
			return 0
		}
	}

	f := os.Stdout
	switch {
	case *outtsort:
		err = tsort.Render(f, g, focus, *compact)
	case *outdot:
		err = dot.Render(f, g, focus, *strata)
	default:
		err = edges.Render(f, g)
	}
	if err != nil {
		log.Fatal(err)
	}
}

func goModule() (string, error) {
	const gomod = "go.mod"
	data, err := os.ReadFile(gomod)
	if err != nil {
		return "", err
	}
	f, err := modfile.Parse(gomod, data, nil)
	if err != nil {
		return "", err
	}
	return f.Module.Mod.Path, nil
}
