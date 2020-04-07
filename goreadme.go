// Package goreadme generates readme markdown file from go doc.
//
// The package can be used as a command line tool and as Github action, described below:
//
// Github Action
//
// Github actions can be configured to update the README.md automatically every time it is needed.
// Below there is an example that on every time a new change is pushed to the master branch, the
// action is trigerred, generates a new README file, and if there is a change - commits and pushes
// it to the master branch.
//
// Add the following content to `.github/workflows/goreadme.yml`:
//
// 	on:
// 	  push:
// 	    branches:
// 	      - master
// 	jobs:
// 	    goreadme:
// 	        runs-on: ubuntu-latest
// 	        steps:
// 	        - name: Check out repository
// 	          uses: actions/checkout@v2
// 	        - name: Update README.md according to Go doc
// 	          uses: posener/goreadme@<release>
// 	          with:
// 	            package_name: 'github.com/<your user>/<your project>'
// 	            badge-travisci: 'true'
// 	            badge-codecov: 'true'
// 	            badge-godoc: 'true'
// 	            badge-goreadme: 'true'
// 	        - name: Commit and push changes
// 	          run: |
// 	            git add README.md
// 	            if git diff --staged --exit-code; then exit 0; fi
// 	            git config user.name Goreadme
// 	            git config user.email posener@gmail.com
// 	            git commit -m "Update readme according to Go doc"
// 	            git push origin HEAD:$(echo "${GITHUB_REF}" | cut -d/ -f3)
//
// Use as a command line tool
//
// 	$ go get github.com/posener/goreadme/...
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
// (godoc page) https://blog.golang.org/godoc-documenting-go-code. Meaning:
//
// * A header is a single line that is separated from a paragraph above.
//
// * Code block is recognized by indentation.
//
// * Diff block is automatically detected:
//
// 	-removed
// 	 stay
// 	+added
//
// * Inline code is marked with backticks.
//
// * URLs will just automatically be converted to links.
//
// Additionally, some extra formatting was added.
//
// * Local paths will be automatically converted to links, for example: ./goreadme.go.
//
// * A URL and can have a title, as follows: (goreadme page) https://github.com/posener/goreadme.
//
// * A local path and can have a title, as follows: (goreadme main file) ./goreamde.go.
//
// * An image can be added: (image/goreadme icon) ./icon.png
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
	PackageName string `json:"package_name"`
	// Functions will make functions documentation to be added to the README.
	Functions bool `json:"functions"`
	// SkipExamples will omit the examples section from the README.
	SkipExamples bool `json:"skip_examples"`
	// SkipSubPackages will omit the sub packages section from the README.
	SkipSubPackages bool `json:"skip_sub_packages"`
	// RecursiveSubPackages will retrieved subpackages information recursively.
	// If false, only one level of subpackages will be retrieved.
	RecursiveSubPackages bool `json:"recursive_sub_packages"`
	Badges               struct {
		Goreadme     bool `json:"goreadme"`
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

	if packageName := r.config.PackageName; packageName != "" {
		p.ImportPath = packageName
	}

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

	if p.IsCmd {
		// TODO: make this better
		p.Name = filepath.Base(name)
		p.Doc = strings.TrimPrefix(p.Doc, "Package main is ")
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
	if os.Getenv("GOREADME_DEBUG") != "1" {
		return
	}

	d, _ := json.MarshalIndent(p, "  ", "  ")
	log.Printf("Package data: %s", string(d))
}
