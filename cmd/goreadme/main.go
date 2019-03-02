// Package main is a command line util that takes a Go repository and write to stdout
// the calculated README.md content.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang/gddo/gosrc"
	"github.com/posener/goreadme"
	"golang.org/x/oauth2"
)

func init() {
	flag.Usage = func() {
		fmt.Println(`goreadme: Create markdown file from go doc.
Usage:
	goreadme -h
		Show this help.
	goreadme [import path]
		Create a readme file. Omitting import path will create
		a readme for the package in CWD.

For accessing private github repositories, a suitable Github token can be
stored in the GITHUB_TOKEN environment variable.

		export GITHUB_TOKEN="<Your github token>"`)
	}
	flag.Parse()
}

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
