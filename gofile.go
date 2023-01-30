package depgraph

import (
	"bufio"
	"go/build/constraint"
	"os"
	"path/filepath"
	"strings"

	"github.com/luno/jettison/errors"
	"github.com/luno/jettison/j"
)

// listGoFiles returns a slice of go files in the folder ignoring _test.go files.
func listGoFiles(dir string) ([]string, error) {
	files, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		return nil, err
	}
	var res []string
	for _, f := range files {
		if strings.Contains(f, "_test.go") {
			continue
		}
		res = append(res, f)
	}
	return res, nil
}

// getGoImports returns the imports of the go file and true if it matches the build tags.
//
// It only scans the head of the file.
func getGoImports(file string, tags map[string]bool) ([]string, bool, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var (
		res          []string
		imports      bool
		singleImport bool
		seenPkg      bool
	)
	for scanner.Scan() {
		line := scanner.Text()

		if !seenPkg && strings.HasPrefix(line, "//go:build") {
			c, err := constraint.Parse(line)
			if err != nil {
				return nil, false, errors.Wrap(err, "", j.KS("file", file))
			}

			if !c.Eval(func(tag string) bool {
				return tags[tag]
			}) {
				return nil, false, nil
			}
		} else if strings.HasPrefix(line, "package ") {
			seenPkg = true
		}

		if line == "import (" {
			imports = true
			continue
		} else if strings.HasPrefix(line, "import \"") {
			imports = true
			singleImport = true
		}

		if !imports {
			continue
		}
		if line == ")" {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Handle line comments
		if strings.HasPrefix(line, "//") {
			continue
		}
		// Handle block comments (only the simple case)
		if strings.HasPrefix(line, "/*") && strings.HasSuffix(line, "*/") {
			continue
		}

		if strings.Contains(line, " \"") {
			// Import aliased
			line = strings.Split(line, " ")[1]
		}
		res = append(res, strings.Trim(line, "\""))

		if singleImport {
			break
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, false, err
	}

	return res, true, nil
}
