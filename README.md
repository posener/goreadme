# goreadme

    go get github.com/posener/goreadme

an HTTP server that works with Github hooks.

[goreadme](./goreadme) is a tool for creating README.md files from Go doc
of a given package.
This server provides Github automation on top of this tool, but creating
PRs for your github repository, whenever the README file should be updated.

## Usage

Go to [https://github.com/apps/goreadme](https://github.com/apps/goreadme)
Press the "Configure" button, choose your account, and add the repositories
you want goreadme to maintain for you.

## How does it Work?

Once enabled, goreadme is registered on a Github hook, that calls goreadme
server the repository default branch is modified.
Goreadme then computes the new README.md file and compairs it to the exiting
one. If a change is needed, Goreadme will create a PR with the new content
of the README.md file.

## Sub Packages

* [cmd/goreadme](./cmd/goreadme): Package main is a command line util that takes a Go repository and write to stdout the calculated README.md content.

* [goreadme](./goreadme): Package goreadme provides API to create readme markdown file from go doc.

* [goreadme/testdata/pkg1](./goreadme/testdata/pkg1): Package pkg1 is a testing package.

* [goreadme/testdata/pkg1/subpkg1](./goreadme/testdata/pkg1/subpkg1): Package subpkg1 is the first subpackage

* [goreadme/testdata/pkg1/subpkg1/subsubpkg](./goreadme/testdata/pkg1/subpkg1/subsubpkg): Package subsubpkg is the sub-subpackage

* [goreadme/testdata/pkg1/subpkg2](./goreadme/testdata/pkg1/subpkg2): Package subpkg1 is the second subpackage.

* [goreadme/testdata/pkg1/testdata](./goreadme/testdata/pkg1/testdata): Package testdata is a package that should be ignored in readme

* [goreadme/testdata/pkg1/testdata/subtd](./goreadme/testdata/pkg1/testdata/subtd): Package subtd is subpackage in testdata

* [goreadme/testdata/pkg2_recursive](./goreadme/testdata/pkg2_recursive): Package pkg2_recursive is a testing package.

* [goreadme/testdata/pkg2_recursive/subpkg1](./goreadme/testdata/pkg2_recursive/subpkg1): Package subpkg1 is the first subpackage

* [goreadme/testdata/pkg2_recursive/subpkg1/subsubpkg](./goreadme/testdata/pkg2_recursive/subpkg1/subsubpkg): Package subsubpkg is the sub-subpackage

* [goreadme/testdata/pkg2_recursive/subpkg2](./goreadme/testdata/pkg2_recursive/subpkg2): Package subpkg1 is the second subpackage.

* [goreadme/testdata/pkg3_skip_examples](./goreadme/testdata/pkg3_skip_examples): Package pkg3_skip_examples is a testing package.

* [goreadme/testdata/pkg3_skip_subpackages](./goreadme/testdata/pkg3_skip_subpackages): Package pkg3_skip_subpackages is a testing package.

* [goreadme/testdata/pkg3_skip_subpackages/subpkg1](./goreadme/testdata/pkg3_skip_subpackages/subpkg1): Package subpkg1 is the first subpackage

* [goreadme/testdata/pkg3_skip_subpackages/subpkg1/subsubpkg](./goreadme/testdata/pkg3_skip_subpackages/subpkg1/subsubpkg): Package subsubpkg is the sub-subpackage

* [goreadme/testdata/pkg3_skip_subpackages/subpkg2](./goreadme/testdata/pkg3_skip_subpackages/subpkg2): Package subpkg1 is the second subpackage.

Created by [goreadme](https://github.com/apps/goreadme)
