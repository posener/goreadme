// Package main is a command line util that takes a Go repository and write to stdout
// the calculated README.md content.
package main

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/golang/gddo/gosrc"

	"github.com/posener/goreadme/goreadme"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()

	gr := goreadme.New(
		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})),
	)

	err := gr.Create(ctx, pkg(), os.Stdout)
	if err != nil {
		log.Fatalf("Failed: %s", err)
	}
}

func pkg() string {
	if len(os.Args) > 1 {
		return os.Args[1]
	}

	path, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}
	gosrc.SetLocalDevMode(path)
	return "."
}
