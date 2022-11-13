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
	"github.com/goeezi/gopkgdep/internal/render/edges"
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
	outTsort := flag.Bool("depthsort", false, "sort by graph depth")
	outCompact := flag.Bool("compact", false, "radix-compresses dependency list (implies -depthsort)")

	// Output: dot
	outDot := flag.Bool("dot", false, "output graphviz (graphviz must be installed)")
	outLayers := flag.Bool("layers", false, "layer graphviz nodes according to graph depth (implies -dot)")

	flag.Parse()

	*outTsort = *outTsort || *outCompact
	*outDot = *outDot || *outLayers

	if *outTsort && *outDot {
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

	paths, err := pathutil.ResolvePaths(flag.Args())
	if err != nil {
		log.Fatal(err)
	}

	g := graph.New()

	m := &match.Matcher{
		InclRE: inclRE,
		ExclRE: exclRE,
		Module: module,
		Stdlib: *stdlib,
		Paths:  paths,
	}

	w := &walk.Walker{
		Import: func(path string) (*build.Package, error) {
			return c.Import(path, wd, 0)
		},
		Closed:  *closed,
		Graph:   g,
		Matcher: m,
	}

	_, err = w.Walk(paths)
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
	case *outTsort:
		err = tsort.Render(f, g, focus, *outCompact)
	case *outDot:
		err = dot.Render(f, g, focus, *outLayers)
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
