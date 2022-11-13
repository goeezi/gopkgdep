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
		punctFormat = "\x1b[38;2;180;60;180m%s\x1b[0m"
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

		headFormat := "%d %s :"
		pkgFormat := [3]string{
			"\x1b[2m%s\x1b[0m",
			"\x1b[32m%s\x1b[0m",
			"\x1b[1;32m%s\x1b[0m",
		}[focus(node)]
		headFormat = "\x1b[38;2;103;103;103m%s\x1b[0m" + pkgFormat + " \x1b[2m:\x1b[0m"

		printf(headFormat, depthFmt, node)
		if compact {
			var sb strings.Builder
			sb.WriteByte(' ')
			t := trie.Build(deps)
			if err := t.Write(&sb, punctFormat, 0); err != nil {
				return err
			}
			printf("%s", sb.String())
		} else {
			for _, dep := range deps {
				printf(" %s", dep)
			}
		}
		printf("\n")
	}
	return nil
}
