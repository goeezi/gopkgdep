package edges

import (
	"fmt"
	"io"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/goeezi/gopkgdep/internal/graph"
)

func Render(w io.Writer, g *graph.Graph) error {
	keys := maps.Keys(g.Pkgs)
	slices.Sort(keys)
	for _, k := range keys {
		dkeys := maps.Keys(g.Pkgs[k])
		slices.Sort(dkeys)
		for _, d := range dkeys {
			if _, err := fmt.Fprintf(w, "%s %s\n", k, d); err != nil {
				return err
			}
		}
	}
	return nil
}
