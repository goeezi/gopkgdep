package edges

import (
	"fmt"
	"io"

	"github.com/goeezi/gopkgdep/internal/graph"
)

func Render(w io.Writer, g *graph.Graph) error {
	for k, v := range g.Pkgs {
		for _, d := range v {
			if _, err := fmt.Fprintf(w, "%s %s\n", k, d); err != nil {
				return err
			}
		}
	}
	return nil
}
