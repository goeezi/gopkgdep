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

package tsort

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/mattn/go-isatty"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/goeezi/gopkgdep/internal/graph"
	"github.com/goeezi/gopkgdep/internal/trie"
)

func Render(w io.Writer, g *graph.Graph, focus func(string) int, compact bool) error {
	tty := isatty.IsTerminal(os.Stdout.Fd())

	punctFormat := "%s"
	if tty {
		punctFormat = "\x1b[38;2;190;70;190m%s\x1b[0m"
	}

	printf := func(format string, args ...any) error {
		_, err := fmt.Fprintf(w, format, args...)
		return err
	}

	nodes := maps.Keys(g.Pkgs)
	slices.SortFunc(nodes, g.Less)
	for i, node := range nodes {
		deps := g.Pkgs[node]
		depth := g.Depth(node)
		depthFmt := strconv.Itoa(depth)
		if tty && i+1 < len(nodes) && depth != g.Depth(nodes[i+1]) {
			depthFmt = fmt.Sprintf("\x1b[4m" + depthFmt)
		}
		depthFmt += " "

		pkgFormat := [3]string{
			"\x1b[2m%s\x1b[0m",
			"\x1b[32m%s\x1b[0m",
			"\x1b[1;32m%s\x1b[0m",
		}[focus(node)]
		headFormat := "\x1b[38;2;103;103;103m%s\x1b[0m" + pkgFormat + " \x1b[2m:\x1b[0m"

		if err := printf(headFormat, depthFmt, node); err != nil {
			return err
		}
		if compact {
			var sb strings.Builder
			sb.WriteByte(' ')
			t := trie.Build(deps)
			if err := t.Write(&sb, punctFormat, 0); err != nil {
				return err
			}
			if err := printf("%s", sb.String()); err != nil {
				return err
			}
		} else {
			for _, dep := range deps {
				if err := printf(" %s", dep); err != nil {
					return err
				}
			}
		}
		if err := printf("\n"); err != nil {
			return err
		}
	}
	return nil
}
