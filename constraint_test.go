package depgraph_test

import (
	"go/build/constraint"
	"testing"

	"github.com/luno/jettison/jtest"
	"github.com/stretchr/testify/require"
)

// TestBuildConstraints demonstrates how the constraint package can be used to
// parse build tags.
func TestBuildConstraints(t *testing.T) {
	testCases := []struct {
		name         string
		sourceTags   string
		providedTags []string
		exp          bool
	}{
		{
			name:         "simple case, tag not provided",
			sourceTags:   "//go:build foo",
			providedTags: nil,
			exp:          false,
		},
		{
			name:         "simple case, tag provided",
			sourceTags:   "//go:build foo",
			providedTags: []string{"foo"},
			exp:          true,
		},
		{
			name:         "and case, neither provided",
			sourceTags:   "//go:build foo && bar",
			providedTags: nil,
			exp:          false,
		},
		{
			name:         "and case, one provided",
			sourceTags:   "//go:build foo && bar",
			providedTags: []string{"foo"},
			exp:          false,
		},
		{
			name:         "and case, both provided",
			sourceTags:   "//go:build foo && bar",
			providedTags: []string{"foo", "bar"},
			exp:          true,
		},
		{
			name:         "or case, neither provided",
			sourceTags:   "//go:build foo || bar",
			providedTags: nil,
			exp:          false,
		},
		{
			name:         "or case, one provided",
			sourceTags:   "//go:build foo || bar",
			providedTags: []string{"bar"},
			exp:          true,
		},
		{
			name:         "or case, both provided",
			sourceTags:   "//go:build foo || bar",
			providedTags: []string{"foo", "bar"},
			exp:          true,
		},
		{
			name:         "not case, not provided",
			sourceTags:   "//go:build !foo",
			providedTags: nil,
			exp:          true,
		},
		{
			name:         "not case, provided",
			sourceTags:   "//go:build !foo",
			providedTags: []string{"foo"},
			exp:          false,
		},
		{
			name:         "more complex case",
			sourceTags:   "//go:build (!foo && bar) || (foo && baz)",
			providedTags: []string{"foo", "bar"},
			exp:          false,
		},
		{
			name:         "more complex case",
			sourceTags:   "//go:build (!foo && bar) || (foo && baz)",
			providedTags: []string{"foo", "baz"},
			exp:          true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c, err := constraint.Parse(testCase.sourceTags)
			jtest.RequireNil(t, err)

			tagMap := make(map[string]bool)
			for _, tag := range testCase.providedTags {
				tagMap[tag] = true
			}
			act := c.Eval(func(tag string) bool {
				return tagMap[tag]
			})
			require.Equal(t, testCase.exp, act)
		})
	}
}
