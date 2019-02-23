// Package goreadme provides API to create readme markdown file from go doc.
package goreadme

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"

	"github.com/golang/gddo/doc"
)

// GoReadme enables getting readme.md text from a go package.
type GoReadme struct {
	// Client is an HTTP client used to perform the requests. It can be used
	// to authenticate github requests, for example, a github client can be used:
	//
	//		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
	//			&oauth2.Token{AccessToken: "...github access token..."
	//		))
	Client *http.Client

	config Config
}

type Config struct {
	// SkipExamples will omit the examples section from the README.
	SkipExamples bool `json:"sub_packages"`
	// SkipSubPackages will omit the sub packages section from the README.
	SkipSubPackages bool `json:"skip_sub_packages"`
	// RecursiveSubPackages will retrived subpackages information recursively.
	// If false, only one level of subpackages will be retrived.
	RecursiveSubPackages bool `json:"recursive_sub_packages"`
}

// Create writes the content of readme.md to w, with the default client.
// name should be a Go repository name, such as "github.com/posener/goreadme".
func Create(ctx context.Context, name string, w io.Writer) error {
	g := GoReadme{Client: http.DefaultClient}
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
	return tmpl.Execute(w, p)
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
	p, err := doc.Get(ctx, r.Client, name, "")
	if err != nil {
		return nil, fmt.Errorf("failed getting %s: %s", name, err)
	}
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

	if p.IsCmd {
		p.Name = p.ProjectName
		p.Doc = strings.TrimPrefix(p.Doc, "Package main is ")
	}

	pkg := &pkg{
		Package: p,
		Config:  r.config,
	}

	if !r.config.SkipSubPackages {
		f := fetcher{
			importPath: name,
			client:     r.Client,
			recursive:  r.config.RecursiveSubPackages,
		}
		pkg.SubPackages, err = f.Fetch(ctx, p)
		if err != nil {
			return nil, err
		}
	}
	return pkg, nil
}

var tmpl = template.Must(template.New("readme").Funcs(
	template.FuncMap{
		"code": func(s string) string { return "```golang\n" + s + "\n```\n" },
	},
).Parse(`
# {{.Package.Name}}

	go get {{.Package.ImportPath}}

{{.Package.Doc}}

{{if (and .SubPackages (not .Config.SkipSubPackages)) -}}
## Sub Packages
{{range .SubPackages}}
* [{{.Path}}](./{{.Path}}){{if .Package.Synopsis}}: {{.Package.Synopsis}}{{end}}
{{end -}}
{{end -}}

{{if (and .Package.Examples (not .Config.SkipExamples)) -}}
## Examples
{{range .Package.Examples}}
### {{.Name}}

{{.Doc}}

{{code .Play}}
{{end -}}
{{end -}}
`))
