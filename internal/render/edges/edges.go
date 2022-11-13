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
