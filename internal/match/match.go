package match

import (
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/mod/modfile"
)

type Matcher struct {
	InclRE *regexp.Regexp
	ExclRE *regexp.Regexp
	Module string
	Stdlib bool
}

func (m *Matcher) Match(pkg string) bool {
	switch {
	case !strings.Contains(pkg, ".") && !m.Stdlib:
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

func (m *Matcher) Rel(pkg string) (string, error) {
	if strings.HasPrefix(pkg, m.Module) {
		rel, err := filepath.Rel(m.Module, pkg)
		if rel != "." {
			rel = "./" + rel
		}
		return rel, err
	}
	return pkg, nil
}
