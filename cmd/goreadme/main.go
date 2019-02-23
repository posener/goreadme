// Package main is a command line util that takes a Go repository and write to stdout
// the calculated README.md content.
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/posener/goreadme/goreadme"
	"golang.org/x/oauth2"
)

var githubToken = os.Getenv("GITHUB_TOKEN")

func main() {
	ctx := context.Background()
	if len(os.Args) < 2 {
		log.Fatal("Missing argument repository name")
	}

	client := http.DefaultClient
	if githubToken != "" {
		client = oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken},
		))
	}

	g := goreadme.GoReadme{Client: client}
	err := g.Create(ctx, os.Args[1], os.Stdout)
	if err != nil {
		log.Fatalf("Failed: %s", err)
	}
}
