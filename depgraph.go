// Package depgraph provides an efficient golang dependency analysis
// tool for large modules. It builds the graph by inspecting file
// imports only, which is faster than go list or ast parsing.
package depgraph

import (
	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
)

// Tags represents golang build tags. The zero value is the default (no build tags).
type Tags map[string]bool

// Make returns a dependency graph with the provided packages
// as root nodes.
//
// Only internal imports (packages of the provided module) are expanded.
// Only direct external imports (packages of other modules) are included
// if inclExt is true.
//
// Build tags are respected and _test.go files are excluded.
func Make(mod Mod, tags Tags, inclExt bool, pkgs ...string) ([]*Node, error) {
	cache := make(map[string]*Node)

	var res []*Node
	for _, pkg := range pkgs {
		if !mod.MaybeHasPkg(pkg) {
			return nil, errors.New("package not of module", j.MKV{"pkg": pkg})
		}

		node, err := getNode(mod, cache, tags, pkg, inclExt, make(map[string]bool))
		if err != nil {
			return nil, err
		}
		res = append(res, node)
	}

	return res, nil
}
