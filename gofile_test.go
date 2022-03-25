package depgraph

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImports(t *testing.T) {
	files, err := listGoFiles(".")
	require.NoError(t, err)

	require.Len(t, files, 4)
	for _, file := range files {
		imports, _, err := getGoImports(file, nil)
		require.NoError(t, err)

		fmt.Printf("%s imports=%v\n", file, imports)
	}
}
