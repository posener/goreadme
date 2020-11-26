// Package goreadme generates readme markdown file from go doc.
//
// The package can be used as a command line tool and as Github action, described below:
//
// Github Action
//
// Github actions can be configured to update the README file automatically every time it is needed.
// Below there is an example that on every time a new change is pushed to the master branch, the
// action is trigerred, generates a new README file, and if there is a change - commits and pushes
// it to the master branch. In pull requests that affect the README content, if the `github-token`
// is given, the action will post a comment on the pull request with chnages that will be made to
// the README file.
//
// To use this with Github actions, add the following content to `.github/workflows/goreadme.yml`.
// See ./actions.yml for all available input options.
//
// 	on:
// 	  push:
// 	    branches: [master]
// 	  pull_request:
// 	    branches: [master]
// 	jobs:
// 	    goreadme:
// 	        runs-on: ubuntu-latest
// 	        steps:
// 	        - name: Check out repository
// 	          uses: actions/checkout@v2
// 	        - name: Update readme according to Go doc
// 	          uses: posener/goreadme@v1
// 	          with:
// 	            badge-travisci: 'true'
// 	            badge-codecov: 'true'
// 	            badge-godoc: 'true'
// 	            badge-goreadme: 'true'
// 	            # Optional: Token allows goreadme to comment the PR with diff preview.
// 	            github-token: '${{ secrets.GITHUB_TOKEN }}'
//
// Use as a command line tool
//
// 	$ GO111MODULE=on go get github.com/posener/goreadme/cmd/goreadme
// 	$ goreadme -h
//
// Why Should You Use It
//
// Both Go doc and readme files are important. Go doc to be used by your user's library, and README
// file to welcome users to use your library. They share common content, which is usually duplicated
// from the doc to the readme or vice versa once the library is ready. The problem is that keeping
// documentation updated is important, and hard enough - keeping both updated is twice as hard.
//
// Go Doc Instructions
//
// The formatting of the README.md is done by the go doc parser. This makes the result README.md a
// bit more limited. Currently, `goreadme` supports the formatting as explained in
// (godoc page) https://blog.golang.org/godoc-documenting-go-code, or
// (here) https://pkg.go.dev/github.com/fluhus/godoc-tricks. Meaning:
//
// * A header is a single line that is separated from a paragraph above.
//
// * Code block is recognized by indentation as Go code.
//
// 	func main() {
// 		...
// 	}
//
// * Inline code is marked with `backticks`.
//
// * URLs will just automatically be converted to links: https://github.com/posener/goreadme
//
// Additionally, the syntax was extended to include some more markdown features while keeping the Go
// doc readable:
//
// * Bulleted and numbered lists are possible when each bullet item is followed by an empty line.
//
// * Diff blocks are automatically detected when each line in a code block starts with a `' '`,
// `'-'` or `'+'`:
//
// 	-removed line starts with '-'
// 	 remained line starts with ' '
// 	+added line starts with '+'
//
// * A repository file can be linked when providing a path that start with `./`: ./goreadme.go.
//
// * A link can have a link text by prefixing it with parenthesised text:
// (goreadme page) https://github.com/posener/goreadme.
//
// * A link to repository file and can have a link text: (goreadme main file) ./goreamde.go.
//
// * An image can be added by prefixing a link to an image with `(image/<image title>)`:
//
// (image/title of image) https://github.githubassets.com/images/icons/emoji/unicode/1f44c.png
package goreadme

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/golang/gddo/doc"
	"github.com/pkg/errors"
	"github.com/posener/goreadme/internal/template"
)

// New returns a GoReadme object with a custom client.
// client is an HTTP client used to perform the requests. It can be used
// to authenticate github requests, for example, a github client can be used:
//
//		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
//			&oauth2.Token{AccessToken: "...github access token..."
//		))
func New(c *http.Client) *GoReadme {
	return &GoReadme{client: c}
}

// GoReadme enables getting readme.md text from a go package.
type GoReadme struct {
	client *http.Client
	config Config
}

type Config struct {
	// Override readme title. Default is package name.
	Title string `json:"title"`
	// ImportPath is used to override the import path. For example: github.com/user/project,
	// github.com/user/project/package or github.com/user/project/version.
	ImportPath string `json:"import_path"`
	// DestinationPath will be where README.md will be written
	DestinationPath string `json:"destination_path"`
	// Functions will make functions documentation to be added to the README.
	Functions bool `json:"functions"`
	// Types will make types documentation to be added to the README.
	Types bool `json:"types"`
	// SkipExamples will omit the examples section from the README.
	SkipExamples bool `json:"skip_examples"`
	// SkipSubPackages will omit the sub packages section from the README.
	SkipSubPackages bool `json:"skip_sub_packages"`
	// RecursiveSubPackages will retrieved subpackages information recursively.
	// If false, only one level of subpackages will be retrieved.
	RecursiveSubPackages bool `json:"recursive_sub_packages"`
	Badges               struct {
		TravisCI     bool `json:"travis_ci"`
		CodeCov      bool `json:"code_cov"`
		GolangCI     bool `json:"golang_ci"`
		GoDoc        bool `json:"go_doc"`
		GoReportCard bool `json:"go_report_card"`
	} `json:"badges"`
	Credit bool `json:"credit"`
}

// Create writes the content of readme.md to w, with the default client.
// name should be a Go repository name, such as "github.com/posener/goreadme".
func Create(ctx context.Context, name string, w io.Writer) error {
	g := GoReadme{client: http.DefaultClient}
	return g.Create(ctx, name, w)
}

// WithConfig returns a copy of the converter with the given configuration.
func (r GoReadme) WithConfig(cfg Config) *GoReadme {
	r.config = cfg
	return &r
}

// Create writes the content of readme.md to w, with r's HTTP client.
// name should be a Go repository name, such as "github.com/posener/goreadme".
func (r *GoReadme) Create(ctx context.Context, name string, w io.Writer) error {
	p, err := r.get(ctx, name)
	if err != nil {
		return err
	}
	return template.Execute(w, p)
}

// pkg contains information about a go package, to be used in the template.
type pkg struct {
	Package     *doc.Package
	SubPackages []subPkg
	Config      Config
}

// subPkg is information about sub package, to be used in the template.
type subPkg struct {
	Path    string
	Package *doc.Package
}

func (r *GoReadme) get(ctx context.Context, name string) (*pkg, error) {
	log.Printf("Getting %s", name)
	p, err := docGet(ctx, r.client, name, "")
	if err != nil {
		return nil, errors.Wrapf(err, "failed getting %s", name)
	}
	sort.Strings(p.Subdirectories)

	// If functions were not requested to be added to the readme, add their
	// examples to the main readme.
	if !r.config.Functions {
		for _, f := range p.Funcs {
			for _, e := range f.Examples {
				if e.Name == "" {
					e.Name = f.Name
				}
				if e.Doc == "" {
					e.Doc = f.Doc
				}
				p.Examples = append(p.Examples, e)
			}
		}
	}

	// If types were not requested to be added to the readme, add their
	// examples to the main readme.
	if !r.config.Types {
		for _, f := range p.Types {
			for _, e := range f.Examples {
				if e.Name == "" {
					e.Name = f.Name
				}
				if e.Doc == "" {
					e.Doc = f.Doc
				}
				p.Examples = append(p.Examples, e)
			}
		}
	}

	if p.IsCmd {
		// TODO: make this better
		p.Name = filepath.Base(name)
		p.Doc = strings.TrimPrefix(p.Doc, "Package main is ")
	}

	if override := r.config.Title; override != "" {
		p.Name = override
	}

	if override := r.config.ImportPath; override != "" {
		p.ImportPath = override
	}

	pkg := &pkg{
		Package: p,
		Config:  r.config,
	}

	if !r.config.SkipSubPackages {
		f := subpackagesFetcher{
			importPath: name,
			client:     r.client,
			recursive:  r.config.RecursiveSubPackages,
		}
		pkg.SubPackages, err = f.Fetch(ctx, p)
		if err != nil {
			return nil, err
		}
	}
	debug(pkg)
	return pkg, nil
}

func debug(p *pkg) {
	if os.Getenv("debug") == "" {
		return
	}

	d, _ := json.MarshalIndent(p, "  ", "  ")
	log.Printf("Package data: %s", string(d))
}
