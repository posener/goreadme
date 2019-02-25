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

* [goreadme](./goreadme): Package goreadme provides API to create readme markdown file from go doc.

Created by [goreadme](https://github.com/apps/goreadme)
