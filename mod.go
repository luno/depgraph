package depgraph

import (
	"path/filepath"
	"strings"
)

type Mod struct {
	// Path of the module: github.com/luno/reflex
	Path string

	// Dir of the module on local disk.
	Dir string
}

// MaybeHasPkg returns true if the package name matches the module.
// It might not actually exist though.
func (m Mod) MaybeHasPkg(pkg string) bool {
	if pkg == m.Path {
		return true
	}

	if !strings.HasPrefix(pkg, m.Path) {
		return false
	}

	return pkg[len(m.Path)] == '/'
}

// PkgDir returns the package directory path on local disk.
func (m Mod) PkgDir(pkg string) string {
	if !m.MaybeHasPkg(pkg) {
		return ""
	}
	if m.Path == pkg {
		return m.Dir
	}
	return filepath.Join(m.Dir, pkg[len(m.Path):])
}
