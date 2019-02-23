// Package main is a command line util that takes a Go repository and write to stdout
// the calculated README.md content.
package main

import (
	"context"
	"log"
	"os"

	"github.com/posener/goreadme/goreadme"
	"golang.org/x/oauth2"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Missing argument repository name")
	}
	ctx := context.Background()

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	))

	err := goreadme.Create(ctx, client, os.Args[1], os.Stdout)
	if err != nil {
		log.Fatalf("Failed: %s", err)
	}
}
