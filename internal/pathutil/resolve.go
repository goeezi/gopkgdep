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
