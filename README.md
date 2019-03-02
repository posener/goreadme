# goreadme


[![Build Status](https://travis-ci.org/posener/goreadme.svg?branch=master)](https://travis-ci.org/posener/goreadme)
[![codecov](https://codecov.io/gh/posener/goreadme/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/goreadme)
[![golangci](https://golangci.com/badges/github.com/posener/goreadme.svg)](https://golangci.com/r/github.com/posener/goreadme)
[![GoDoc](https://godoc.org/github.com/posener/goreadme?status.svg)](http://godoc.org/github.com/posener/goreadme)
[![Go Report Card](https://goreportcard.com/badge/github.com/posener/goreadme)](https://goreportcard.com/report/github.com/posener/goreadme)Package goreadme creates readme markdown file from go doc.

This package can be used as a web service, as a command line tool or as a library.

Try the web service: [https://gotreadme.herokuapp.com](https://gotreadme.herokuapp.com)

Integrate directly with Github: [https://github.com/apps/goreadme](https://github.com/apps/goreadme).

Use as a command line tool:

		$ go get github.com/posener/goreadme/...
		$ goreadme

## Sub Packages

* [cmd/goreadme](./cmd/goreadme): Package main is a command line util that takes a Go repository and write to stdout the calculated README.md content.

Created by [goreadme](https://github.com/apps/goreadme)
