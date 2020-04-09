# goreadme

[![Build Status](https://travis-ci.org/posener/goreadme.svg?branch=master)](https://travis-ci.org/posener/goreadme)
[![codecov](https://codecov.io/gh/posener/goreadme/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/goreadme)
[![GoDoc](https://godoc.org/github.com/posener/goreadme?status.svg)](http://godoc.org/github.com/posener/goreadme)

Package goreadme generates readme markdown file from go doc.

The package can be used as a command line tool and as Github action, described below:

## Github Action

Github actions can be configured to update the README.md automatically every time it is needed.
Below there is an example that on every time a new change is pushed to the master branch, the
action is trigerred, generates a new README file, and if there is a change - commits and pushes
it to the master branch.

Add the following content to `.github/workflows/goreadme.yml`:

```go
on:
  push:
    branches: [master]
  pull_request:
    branches: [master]
jobs:
    goreadme:
        runs-on: ubuntu-latest
        steps:
        - name: Check out repository
          uses: actions/checkout@v2
        - name: Update readme according to Go doc
          uses: posener/goreadme@<release>
          with:
            badge-travisci: 'true'
            badge-codecov: 'true'
            badge-godoc: 'true'
            badge-goreadme: 'true'
            github-token: '${{ secrets.GITHUB_TOKEN }}'
```

Use as a command line tool

```go
$ GO111MODULE=on go get github.com/posener/goreadme/cmd/goreadme
$ goreadme -h
```

## Why Should You Use It

Both Go doc and readme files are important. Go doc to be used by your user's library, and README
file to welcome users to use your library. They share common content, which is usually duplicated
from the doc to the readme or vice versa once the library is ready. The problem is that keeping
documentation updated is important, and hard enough - keeping both updated is twice as hard.

## Go Doc Instructions

The formatting of the README.md is done by the go doc parser. This makes the result README.md a
bit more limited. Currently, `goreadme` supports the formatting as explained in
[godoc page](https://blog.golang.org/godoc-documenting-go-code). Meaning:

* A header is a single line that is separated from a paragraph above.

* Code block is recognized by indentation as Go code.

```go
func main() {
	...
}
```

* Inline code is marked with `backticks`.

* URLs will just automatically be converted to links: [https://github.com/posener/goreadme](https://github.com/posener/goreadme)

Additionally, some extra formatting was added.

* Bullets are recognized when each bullet item is followed by an empty line.

* Diff block is automatically detected:

```diff
-removed
 stay
+added
```

* Local paths will be automatically converted to links: [./goreadme.go](./goreadme.go).

* A URL and can have a title: [goreadme page](https://github.com/posener/goreadme).

* A local path and can have a title: [goreadme main file](./goreamde.go).

* An image can be added:

![title of image](https://github.githubassets.com/images/icons/emoji/unicode/1f44c.png)

## Sub Packages

* [cmd/goreadme](./cmd/goreadme): Package main is a command line util that takes a Go repository and write to stdout the calculated README.md content.

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
