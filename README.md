# goreadme

[![Build Status](https://travis-ci.org/goreadme.svg?branch=master)](https://travis-ci.org/goreadme)[![codecov](https://codecov.io/gh/goreadme/branch/master/graph/badge.svg)](https://codecov.io/gh/goreadme)[![golangci](https://golangci.com/badges/github.com/goreadme.svg)](https://golangci.com/r/github.com/goreadme)[![GoDoc](https://godoc.org/github.com/goreadme?status.svg)](http://godoc.org/github.com/goreadme)[![Go Report Card](https://goreportcard.com/badge/github.com/goreadme)](https://goreportcard.com/report/github.com/goreadme)Package goreadme creates readme markdown file from go doc.

This package can be used as a web service, as a command line tool or as a library.

Try the web service: [https://gotreadme.herokuapp.com](https://gotreadme.herokuapp.com)

Integrate directly with Github: [https://github.com/apps/goreadme](https://github.com/apps/goreadme).

Use as a command line tool:

		$ go get github.com/posener/goreadme/...
		$ goreadme

## Sub Packages

* [cmd/goreadme](./cmd/goreadme): Package main is a command line util that takes a Go repository and write to stdout the calculated README.md content.

Created by [goreadme](https://github.com/apps/goreadme)
