# goreadme

[![Build Status](https://travis-ci.org/posener/goreadme.svg?branch=master)](https://travis-ci.org/posener/goreadme)
[![codecov](https://codecov.io/gh/posener/goreadme/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/goreadme)
[![golangci](https://golangci.com/badges/github.com/posener/goreadme.svg)](https://golangci.com/r/github.com/posener/goreadme)
[![GoDoc](https://godoc.org/github.com/posener/goreadme?status.svg)](http://godoc.org/github.com/posener/goreadme)
[![goreadme](https://goreadme.herokuapp.com/badge/posener/goreadme.svg)](https://goreadme.herokuapp.com)

Package goreadme creates readme markdown file from go doc.

This package can be used as a web service, as a command line tool or as a library.

Try the web service: [https://gotreadme.herokuapp.com](https://gotreadme.herokuapp.com).

Integrate directly with Github: [https://github.com/apps/goreadme](https://github.com/apps/goreadme).

Use as a command line tool:

```go
$ go get github.com/posener/goreadme/...
$ goreadme -h
```

Why should you use it?

Both go doc and readme files are important. Go doc to be used by your user's
library, and README file to welcome users to use your library. They share
common content, which is usually duplicated from the doc to the readme or vice versa
once the library is ready. The problem is that keeping documentation updated
is important, and hard enough - keeping both updated is twice as hard.

This library provides an easy way to create the one from the other. Using the
goreadme Github App: [https://github.com/apps/goreadme](https://github.com/apps/goreadme) makes it even easier.

#### Go Doc Instructions

The formatting of the README.md is done by the go doc parser. This makes the
Result README.md a bit more limited.
Currently, `goreadme` supports the formatting as explained in: [https://blog.golang.org/godoc-documenting-go-code](https://blog.golang.org/godoc-documenting-go-code).
Meaning:
* A header is a single line that is separated from a paragraph above.
* Code block is recognized by indentation.
* URLs will just automatically converted to links

```go
(Added support to local links - prefixed by dot and slash, for example: [./cmd](./cmd)).
```

* Lists are also available with "* " prefix or a "1. " prefix.

## Sub Packages

* [cmd/goreadme](./cmd/goreadme): Package main is a command line util that takes a Go repository and write to stdout the calculated README.md content.

Created by [goreadme](https://github.com/apps/goreadme)
