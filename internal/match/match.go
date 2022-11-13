package match

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/goeezi/gopkgdep/internal/set"
	"golang.org/x/mod/modfile"
)

type Matcher struct {
	Paths  set.Set[string]
	InclRE *regexp.Regexp
	ExclRE *regexp.Regexp
	Module string
	Stdlib bool
}

func (m *Matcher) Match(pkg string) bool {
	switch {
	case !strings.Contains(pkg, ".") && !m.Stdlib:
		return false
	case !m.Paths.Has(m.Rel(pkg)):
		return false
	case !m.InclRE.MatchString(pkg):
		return false
	case m.ExclRE != nil:
		return !m.ExclRE.MatchString(pkg)
	default:
		return true
	}
}

func (m *Matcher) MatchUnlessExcludedStdLibPackage(pkg string) bool {
	return strings.Contains(pkg, ".") || m.Stdlib
}

func (m *Matcher) Resolve(pkg string) string {
	if pkg == "." {
		return m.Module
	}
	if modfile.IsDirectoryPath(pkg) && !filepath.IsAbs(pkg) {
		return filepath.Join(m.Module, pkg)
	}
	return pkg
}

func (m *Matcher) Rel(pkg string) string {
	if strings.HasPrefix(pkg, m.Module) {
		rel, err := filepath.Rel(m.Module, pkg)
		if err != nil {
			panic(err)
		}
		if rel != "." {
			rel = "./" + rel
		}
		return rel
	}
	return pkg
}
