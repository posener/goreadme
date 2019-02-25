// Package main is an HTTP server that works with Github hooks.
//
// [goreadme](./goreadme) is a tool for creating README.md files from Go doc
// of a given package.
// This server provides Github automation on top of this tool, but creating
// PRs for your github repository, whenever the README file should be updated.
//
// ## Usage
//
// Go to [https://github.com/apps/goreadme](https://github.com/apps/goreadme)
// Press the "Configure" button, choose your account, and add the repositories
// you want goreadme to maintain for you.
//
// ## How does it Work?
//
// Once enabled, goreadme is registered on a Github hook, that calls goreadme
// server the repository default branch is modified.
// Goreadme then computes the new README.md file and compairs it to the exiting
// one. If a change is needed, Goreadme will create a PR with the new content
// of the README.md file.
package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"github.com/posener/goreadme/goreadme"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	githubAppURL      = "https://github.com/apps/goreadme"
	timeout           = time.Second * 60
	defaultReadmePath = "README.md"
	configPath        = "goreadme.json"

	goreadmeAuthor = "goreadme"
	goreadmeEmail  = "posener@gmail.com"
	goreadmeBranch = "goreadme"
	goreaedmeRef   = "refs/heads/" + goreadmeBranch
)

var (
	port         = os.Getenv("PORT")
	githubToken  = os.Getenv("GITHUB_TOKEN")
	githubSecret = []byte(os.Getenv("GITHUB_SECRET")) // Secret for github hooks
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)

	client := oauth2.NewClient(ctx, ts)
	h := &handler{
		github:   github.NewClient(client),
		goreadme: goreadme.New(client),
	}
	h.debugPR()
	m := mux.NewRouter()
	m.Methods("GET").Path("/").HandlerFunc(h.home)
	m.Methods("POST").Path("/github/hook").HandlerFunc(h.hook)
	logrus.Infof("Starting server...")
	http.ListenAndServe(":"+port, m)
}
