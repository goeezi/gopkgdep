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

package pathutil

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/goeezi/gopkgdep/internal/set"
)

func ResolvePaths(paths []string) (set.Set[string], error) {
	if len(paths) == 0 {
		paths = []string{"./..."}
	}
	resolvedPaths := set.Set[string]{}
	for _, pathArg := range paths {
		dir := strings.TrimSuffix(pathArg, "/...")
		if len(dir) == len(pathArg) {
			resolvedPaths.Add(pathArg)
		} else {
			err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
				if !d.IsDir() && filepath.Ext(path) == ".go" {
					resolvedPaths.Add("./" + filepath.Dir(path))
					return nil
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
		}
	}
	if resolvedPaths.Has("./.") {
		resolvedPaths.Delete("./.")
		resolvedPaths.Add(".")
	}
	return resolvedPaths, nil
}
