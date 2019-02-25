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

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/posener/goreadme/goreadme"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	port         = os.Getenv("PORT")
	dbURL        = os.Getenv("DATABASE_URL")
	githubToken  = os.Getenv("GITHUB_TOKEN")
	githubSecret = []byte(os.Getenv("GITHUB_SECRET")) // Secret for github hooks
)

func main() {
	ctx := context.Background()

	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		logrus.Fatalf("Connect to DB on %s: %v", dbURL, err)
	}
	defer db.Close()

	db.LogMode(true)

	if err := db.AutoMigrate(&Job{}).Error; err != nil {
		logrus.Fatalf("Migrate database: %s", err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)

	client := oauth2.NewClient(ctx, ts)
	h := &handler{
		db:       db,
		github:   github.NewClient(client),
		goreadme: goreadme.New(client),
	}
	h.debugPR()
	m := mux.NewRouter()
	m.Methods("GET").Path("/").HandlerFunc(h.home)
	m.Methods("POST").Path("/github/hook").HandlerFunc(h.hook)

	m.Methods("GET").Path("/jobs").HandlerFunc(h.jobsList)

	logrus.Infof("Starting server...")
	http.ListenAndServe(":"+port, m)
}
