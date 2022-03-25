package depgraph

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/luno/jettison/jtest"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

//go:generate go test -update

func TestSelf(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	mod := Mod{
		Path: "github.com/luno/depgraph",
		Dir:  strings.TrimSuffix(wd, "/lib/depgraph"),
	}

	res, err := Make(mod, nil, true, "github.com/luno/depgraph")
	jtest.RequireNil(t, err)
	goldie.New(t).AssertJson(t, t.Name(), makeTestNodes(res))
}

func TestTags(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	root := "github.com/luno/depgraph"
	mod := Mod{
		Path: root,
		Dir:  path.Join(wd, "testdata"),
	}

	pkgs := []string{root, path.Join(root, "childa"), path.Join(root, "childb")}

	res, err := Make(mod, map[string]bool{"tag": true}, true, pkgs...)
	jtest.RequireNil(t, err)
	goldie.New(t).AssertJson(t, t.Name()+"True", makeTestNodes(res))

	res, err = Make(mod, map[string]bool{"tag": false}, true, pkgs...)
	jtest.RequireNil(t, err)
	goldie.New(t).AssertJson(t, t.Name()+"False", makeTestNodes(res))

	res, err = Make(mod, nil, true, pkgs...)
	jtest.RequireNil(t, err)
	goldie.New(t).AssertJson(t, t.Name()+"None", makeTestNodes(res))

	res, err = Make(mod, nil, false, pkgs...)
	jtest.RequireNil(t, err)
	goldie.New(t).AssertJson(t, t.Name()+"NoneInternal", makeTestNodes(res))
}

func TestCircular(t *testing.T) {
	wd, err := os.Getwd()
	jtest.RequireNil(t, err)

	root := "github.com/luno/depgraph/testdata"
	mod := Mod{
		Path: root,
		Dir:  path.Join(wd, "testdata"),
	}

	pkgs := []string{root, path.Join(root, "circular_a")}
	_, err = Make(mod, nil, true, pkgs...)
	require.EqualError(t, err, "cyclic import detected; descendant imports ancestor")
}

func makeTestNodes(nl []*Node) []testnode {
	var res []testnode
	for _, n := range nl {
		res = append(res, makeTestNode(n))
	}
	return res
}

func makeTestNode(n *Node) testnode {
	res := testnode{
		Name:        n.Name,
		Parents:     len(n.Parents),
		FileImports: n.FileImports,
		Height:      n.Height,
	}
	for _, child := range n.Children {
		c := makeTestNode(child)
		res.Children = append(res.Children, c)
	}
	res.Descendants = make(map[string]bool)
	for desc := range n.Descendants {
		res.Descendants[desc.Name] = true
	}

	return res
}

type testnode struct {
	Name string

	// Children are direct downstream packages this package imports.
	Children []testnode `json:",omitempty"`

	// Descendants is union of all downstream packages by total number of occurrences.
	Descendants map[string]bool `json:",omitempty"`

	FileImports map[string][]string `json:",omitempty"`

	Height int

	Parents int
}
