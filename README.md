# DepGraph

[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://github.com/luno/depgraph/blob/main/LICENSE)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=luno_depgraph&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=luno_depgraph)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=luno_depgraph&metric=coverage)](https://sonarcloud.io/summary/new_code?id=luno_depgraph)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=luno_depgraph&metric=bugs)](https://sonarcloud.io/summary/new_code?id=luno_depgraph)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=luno_depgraph&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=luno_depgraph)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=luno_depgraph&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=luno_depgraph)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=luno_depgraph&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=luno_depgraph)
[![Go Report Card](https://goreportcard.com/badge/github.com/luno/depgraph)](https://goreportcard.com/report/github.com/luno/depgraph)
[![GoDoc](https://godoc.org/github.com/luno/depgraph?status.png)](https://godoc.org/github.com/luno/depgraph)


DepGraph is a package to figure out an efficient dependency tree for a Go package.

The dependency graph is built by inspecting file imports only, which is faster than `go list`, or AST parsing.

It can be used to figure out why a particular package is imported into a service.

Examples
1. Why is `modulename/currency` imported in `modulename/services/fe`?
```shell
depgraph modulename/services/fe modulename/currency
```
```
modulename/services/fe
└── modulename/locale/allstrings
    └── modulename/locale
        └── modulename/currency
```

2. Which apps is `modulename/services/withdrawals/ops` imported in?
```shell
depgraph modulename/services/withdrawals/ops
```
```
modulename/services/fe/website/modulename_website
└── modulename/services/fe/website
    └── modulename/services/fe/api/base
        └── modulename/services/withdrawals/ops/send
            └── modulename/services/withdrawals/ops

modulename/services/withdrawals/withdrawals
└── modulename/services/withdrawals/state
    └── modulename/services/withdrawals/client/logical
        └── modulename/services/withdrawals/ops
```
