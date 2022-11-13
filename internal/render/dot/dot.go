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

package dot

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/linqgo/linq"
	"golang.org/x/exp/slices"

	"github.com/goeezi/gopkgdep/internal/graph"
	"github.com/goeezi/gopkgdep/internal/pathutil"
)

func Render(w io.Writer, g *graph.Graph, focus func(string) int, layers bool) error {
	nodes, err := linq.ToMapKV(linq.Select(
		linq.Index(linq.Order(linq.SelectKeys(linq.FromMap(g.Pkgs)))),
		func(e linq.KV[int, string]) linq.KV[string, int] {
			return linq.NewKV(e.Value, e.Key)
		},
	))
	if err != nil {
		log.Fatal(err)
	}

	in := 0
	prefix := ""
	nl := false
	indent := func(n int) {
		in += n
		prefix = fmt.Sprintf("%*s", 2*in, "")
	}
	printf := func(format string, args ...any) {
		s := fmt.Sprintf(format, args...)
		if nl {
			s = prefix + s
		}
		nl = strings.HasSuffix(s, "\n")
		s = strings.TrimSuffix(s, "\n")
		if in > 0 {
			s = strings.ReplaceAll(s, "\n", "\n"+prefix)
		}
		fmt.Fprint(w, s)
		if nl {
			fmt.Fprintln(w)
		}
	}
	inf := func(format string, args ...any) {
		printf(format, args...)
		indent(1)
	}
	outf := func(format string, args ...any) {
		indent(-1)
		printf(format, args...)
	}

	inf("digraph G {\n")

	keys := linq.SelectKeys(linq.FromMap(nodes)).Select(func(t string) string {
		if !strings.Contains(t, ".") {
			return "<>/" + t
		}
		return t
	}).ToSlice()
	slices.Sort(keys)
	var prevKey []string
	for i, k := range keys {
		path := pathutil.Split(k)
		if path[0] == "<>/" {
			path[0] = "<stdlib>/"
		}
		_, a, b := pathutil.Diff(pathutil.Trim(prevKey), pathutil.Trim(path))
		prevKey = path
		for range a {
			outf("}\n")
		}
		for j, part := range b {
			inf("subgraph cluster%d_%d {\n", i, j)
			c := 255 - 15*(in-2)
			if part == "<stdlib>/" {
				part = "<stdlib>"
			}
			printf("label=%q;\n", part)
			printf("bgcolor=\"#%02x%02[1]x%02[1]x\";\n", c)
		}
		var label, fill string
		if layers {
			label = fmt.Sprintf(
				`<<FONT>%s<BR/><FONT POINT-SIZE="16"><B>%s</B></FONT></FONT>>`,
				strings.Join(path[:len(path)-1], "<BR/>"),
				pathutil.Last(path),
			)
			fill = `"#eeeeee"`
		} else {
			label = fmt.Sprintf("%q", pathutil.Last(path))
			fill = "white"
		}

		type style struct {
			width float64
			color string
			fill  string
		}
		s := [3]style{
			{0.625, "gray", `"#f8f8f8"`},
			{1, "black", fill},
			{3, "black", fill},
		}[focus(k)]
		printf(
			`%d [`+
				`label=%s;`+
				`style=filled;`+
				`penwidth=%f;`+
				`fontname="Helvetica";`+
				`color=%s;`+
				`fontcolor=%[4]s;`+
				`fillcolor=%s];`+
				"\n",
			i, label, s.width, s.color, s.fill,
		)
	}
	for range prevKey[1:] {
		outf("}\n")
	}

	if layers {
		printf("\n")
		pkgs := linq.SelectKeys(linq.FromMap(g.Pkgs)).OrderComp(g.Less)
		for _, depthPkgs := range linq.GroupBy(pkgs, g.Depth).ToSlice() {
			printf("cluster { rank = same;")
			for _, pkg := range depthPkgs.Value.ToSlice() {
				printf(" %d;", nodes[pkg])
			}
			printf(" }\n")
		}
	}

	printf("\n")
	for i, k := range keys {
		printf("%d -> {", i)
		j := 0
		for d := range g.Pkgs[k] {
			if j > 0 {
				printf(" ")
			}
			j++
			printf("%d", nodes[d])
		}
		printf("}\n")
	}
	outf("}\n")

	return nil
}
