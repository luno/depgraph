# DepGraph

[![PkgGoDev](https://pkg.go.dev/badge/github.com/luno/depgraph)](https://pkg.go.dev/github.com/luno/depgraph)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://github.com/luno/depgraph/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/luno/depgraph)](https://goreportcard.com/report/github.com/luno/depgraph)

DepGraph is a package to figure out an efficient dependency tree for a Go package.

The dependency graph is built by inspecting file imports only, which is faster than `go list`, or ast parsing.

It can be used to figure out why a particular package is imported into a service.

## Example

Why does service A import service B?

```shell
depgraph src/serviceA src/serviceB
```

```
src/serviceA
└── src/fe/api/base
    └── src/fe/api/pkg
        └── src/fe/pkg
            └── src/serviceB
```
