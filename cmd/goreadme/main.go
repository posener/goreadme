// Package main is a command line util that takes a Go repository and write to stdout
// the calculated README.md content.
//
// It can create the README.md from a remote Github repository or from a local Go module.
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

var cfg goreadme.Config

func init() {
	flag.BoolVar(&cfg.RecursiveSubPackages, "recursive", false, "Load docs recursively.")
	flag.BoolVar(&cfg.Functions, "functions", false, "Write functions section.")
	flag.BoolVar(&cfg.SkipExamples, "skip-examples", false, "Skip the examples section.")
	flag.BoolVar(&cfg.SkipSubPackages, "skip-sub-packages", false, "Skip the sub packages section.")
	flag.BoolVar(&cfg.Badges.Goreadme, "badge-goreadme", false, "Show Goreadme badge.")
	flag.BoolVar(&cfg.Badges.TravicCI, "badge-travisci", false, "Show TravicCI badge.")
	flag.BoolVar(&cfg.Badges.CodeCov, "badge-codecov", false, "Show CodeCov badge.")
	flag.BoolVar(&cfg.Badges.GolangCI, "badge-golangci", false, "Show GolangCI badge.")
	flag.BoolVar(&cfg.Badges.GoDoc, "badge-godoc", false, "Show GoDoc badge.")
	flag.BoolVar(&cfg.Badges.GoReportCard, "badge-goreportcard", false, "Show GoReportCard badge.")
	flag.Usage = func() {
		fmt.Fprint(
			flag.CommandLine.Output(),
			`goreadme: Create markdown file from go doc.

Usage:
	goreadme [flags] [import path]

import path (optional): Create a readme file for a package from github.
 Omitting import path will create a readme for the package in CWD.
Flags:
`)
		flag.PrintDefaults()
	}
	flag.Parse()
}

func main() {
	ctx := context.Background()
	gr := goreadme.New(
		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")})),
	)

	err := gr.WithConfig(cfg).Create(ctx, pkg(flag.Args()), os.Stdout)
	if err != nil {
		log.Fatalf("Failed: %s", err)
	}
}

func pkg(args []string) string {
	if len(args) > 0 {
		return args[0]
	}

	path, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}
	gosrc.SetLocalDevMode(path)
	return "."
}
