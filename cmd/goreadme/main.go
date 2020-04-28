// Goreadme command line tool and Github action
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/golang/gddo/gosrc"
	"github.com/posener/goaction"
	"github.com/posener/goaction/actionutil"
	"github.com/posener/goreadme"
	"golang.org/x/oauth2"
)

var (
	// Holds configuration for Goreadme invokation.
	cfg goreadme.Config

	// Write readme output
	out io.WriteCloser = os.Stdout

	// Github action variables.
	path        = goaction.Getenv("readme-file", "README.md", "Name of readme file")
	debug       = goaction.Getenv("debug", "", "Print Goredme debug output") != ""
	email       = goaction.Getenv("email", "posener@gmail.com", "Email for commit message")
	githubToken = goaction.Getenv("github-token", "", "Github token for PR comments. Optional.")
)

func init() {
	log.SetFlags(log.Lshortfile)

	flag.StringVar(&cfg.ImportPath, "import-path", "", "Override package import path.")
	flag.BoolVar(&cfg.RecursiveSubPackages, "recursive", false, "Load docs recursively.")
	flag.BoolVar(&cfg.Functions, "functions", false, "Write functions section.")
	flag.BoolVar(&cfg.SkipExamples, "skip-examples", false, "Skip the examples section.")
	flag.BoolVar(&cfg.SkipSubPackages, "skip-sub-packages", false, "Skip the sub packages section.")
	flag.BoolVar(&cfg.Badges.TravisCI, "badge-travisci", false, "Show TravisCI badge.")
	flag.BoolVar(&cfg.Badges.CodeCov, "badge-codecov", false, "Show CodeCov badge.")
	flag.BoolVar(&cfg.Badges.GolangCI, "badge-golangci", false, "Show GolangCI badge.")
	flag.BoolVar(&cfg.Badges.GoDoc, "badge-godoc", false, "Show GoDoc badge.")
	flag.BoolVar(&cfg.Badges.GoReportCard, "badge-goreportcard", false, "Show GoReportCard badge.")
	flag.BoolVar(&cfg.Credit, "credit", true, "Add credit line.")
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
	// Steps to do only in Github Action mode.
	if goaction.CI {
		// Setup debug mode.
		if debug {
			os.Setenv("GOREADME_DEBUG", "1")
		}

		// Setup output file.
		var err error
		out, err = os.Create(path)
		if err != nil {
			log.Fatalf("Failed opening file %s: %s", path, err)
		}
		defer out.Close()

		// Fix import path if it was not overridden by the user.
		if cfg.ImportPath == "" {
			cfg.ImportPath = "github.com/" + goaction.Repository
		}
	}

	ctx := context.Background()
	gr := goreadme.New(
		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: githubToken})),
	)

	err := gr.WithConfig(cfg).Create(ctx, pkg(flag.Args()), out)
	if err != nil {
		log.Fatalf("Failed: %s", err)
	}

	if !goaction.CI {
		return
	}

	// Runs only in Github CI mode.

	diff := gitDiff()

	log.Printf("Diff:\n\n%s\n", diff)

	switch {
	case goaction.IsPush():
		if diff == "" {
			log.Println("No changes were made. Skipping push.")
			break
		}
		push()
	case goaction.IsPR():
		pr(diff)
	default:
		log.Fatalf("unexpected action mode.")
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

func gitDiff() string {
	// Add files to git, in case it does not exists
	d, err := actionutil.GitDiff(path)
	if err != nil {
		log.Fatal(err)
	}
	if d == "" {
		return ""
	}
	return fmt.Sprintf("Path: %s\n\n```diff\n%s\n```\n\n", path, d)
}

// Commit and push chnages to upstream branch.
func push() {
	err := actionutil.GitConfig("goreadme", email)
	if err != nil {
		log.Fatal(err)
	}

	err = actionutil.GitCommitPush([]string{path}, "Update readme accoridng to godoc")
	if err != nil {
		log.Fatal(err)
	}
}

// Post a pull request comment with the expected diff.
func pr(diff string) {
	if githubToken == "" {
		log.Println("In order to add request comment, set the GITHUB_TOKEN input.")
		return
	}

	body := "[goreadme](https://github.com/posener/goreadme) will not make any changes in this PR"
	if diff != "" {
		body = fmt.Sprintf(
			"[goreadme](https://github.com/posener/goreadme) diff for %s file for this PR:\n\n%s",
			path,
			diff)
	}

	ctx := context.Background()
	err := actionutil.PRComment(ctx, githubToken, "goreadme", body)
	if err != nil {
		log.Fatal(err)
	}
}
