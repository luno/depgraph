package depgraph

import (
	"path/filepath"
	"sort"
	"sync"

	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
)

// Node represents a package in a import/dependency graph.
type Node struct {
	// Name is the package path.
	Name string

	// Children are all direct downstream packages this package imports.
	Children []*Node

	// Descendants is union of all downstream packages.
	// Only internal packages (same module) are expanded.
	Descendants map[*Node]bool

	// Parents are direct upstream packages importing this package.
	Parents []*Node

	// Height is depth of deepest internal branch (chain of same module descendants).
	Height int

	// FileImports are the imports per go file (relative to module dir).
	// Only applicable files are included, so _test.go files and mismatching build tags are excluded.
	FileImports map[string][]string

	mu sync.Mutex
}

// AddParent adds a parent safely.
func (n *Node) AddParent(parent *Node) {
	n.mu.Lock()
	defer n.mu.Unlock()

	n.Parents = append(n.Parents, parent)
}

// getNode returns a Node struct for the package either from the cache
// or by making it from the filesystem.
func getNode(mod Mod, cache map[string]*Node, tags map[string]bool, pkg string, inclExt bool, ancestors map[string]bool) (*Node, error) {
	if c, ok := cache[pkg]; ok {
		return c, nil
	}

	res, err := makeNode(mod, cache, tags, pkg, inclExt, ancestors)
	if err != nil {
		return nil, err
	}

	cache[pkg] = res

	return res, err
}

// makeNode returns a Node struct that represents the package in the module.
//
// It recursively populates imported packages as child nodes if the package
// is part of the module. Non-module packages are therefore not expanded.
func makeNode(mod Mod, cache map[string]*Node, tags map[string]bool, pkg string, inclExt bool, ancestors map[string]bool) (*Node, error) {
	if !mod.MaybeHasPkg(pkg) {
		// Do not expand non-module packages.
		return &Node{
			Name: pkg,
		}, nil
	}

	files, err := listGoFiles(mod.PkgDir(pkg))
	if err != nil {
		return nil, err
	}

	imports := make(map[string]bool)
	fileImports := make(map[string][]string)
	for _, file := range files {
		il, ok, err := getGoImports(file, tags)
		if err != nil {
			return nil, err
		} else if !ok {
			continue
		}

		// Maybe filter external imports.
		if !inclExt {
			var temp []string
			for _, imprt := range il {
				if !inclExt && !mod.MaybeHasPkg(imprt) {
					continue
				}
				temp = append(temp, imprt)
			}
			il = temp
		}

		for _, imprt := range il {
			imports[imprt] = true
		}

		relFile, err := filepath.Rel(mod.Dir, file)
		if err != nil {
			return nil, err
		}
		fileImports[relFile] = il
	}

	node := &Node{
		Name:        pkg,
		Descendants: make(map[*Node]bool),
		FileImports: fileImports,
	}

	// Sort the imports
	var keys []string
	for imprt := range imports {
		keys = append(keys, imprt)
	}
	sort.Strings(keys)

	ancestors[node.Name] = true
	defer delete(ancestors, node.Name)

	var maxHeight int
	var hasModChild bool
	for _, imprt := range keys {
		if ancestors[imprt] {
			return nil, errors.New("cyclic import detected; descendant imports ancestor",
				j.MKS{"ancestor": imprt, "descendent": node.Name})
		}

		child, err := getNode(mod, cache, tags, imprt, inclExt, ancestors)
		if err != nil {
			return nil, err
		}

		// Add this direct child
		node.Children = append(node.Children, child)
		node.Descendants[child] = true

		// Node ancestor include union of all child descendants.
		for p := range child.Descendants {
			node.Descendants[p] = true
		}

		// Calculate height
		if mod.MaybeHasPkg(child.Name) {
			hasModChild = true
			if child.Height > maxHeight {
				maxHeight = child.Height
			}
		}

		// Add node as parent to child
		child.AddParent(node)
	}

	if hasModChild {
		node.Height = maxHeight + 1
	}

	return node, nil
}
