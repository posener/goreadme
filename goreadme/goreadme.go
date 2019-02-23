// Package goreadme provides API to create readme markdown file from go doc.
package goreadme

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"text/template"

	"github.com/golang/gddo/doc"
	"github.com/hashicorp/go-multierror"
)

// pkg contains information about a go package, to be used in the template.
type pkg struct {
	Package     *doc.Package
	SubPackages []subPkg
}

// subPkg is information about sub package, to be used in the template.
type subPkg struct {
	Path string
	Doc  string
}

// GoReadme enables gettring readme.md text from a go package.
type GoReadme struct {
	// Client is an HTTP client used to perform the requests. It can be used
	// to authenticate github requests, for example, a github client can be used:
	//
	//		oauth2.NewClient(ctx, oauth2.StaticTokenSource(
	//			&oauth2.Token{AccessToken: "...github access token..."
	//		))
	Client *http.Client
}

// Create writes the content of readme.md to w, with the default client.
// name should be a Go repository name, such as "github.com/posener/goreadme".
func Create(ctx context.Context, name string, w io.Writer) error {
	g := GoReadme{Client: http.DefaultClient}
	return g.Create(ctx, name, w)
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
	pkg := &pkg{Package: p}
	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		errors *multierror.Error
	)

	// Concurrently get information for all sub directories.
	wg.Add(len(p.Subdirectories))
	for _, subDir := range p.Subdirectories {
		go func(subDir string) {
			importPath := name + "/" + subDir
			defer wg.Done()
			log.Printf("Getting %s", importPath)
			sp, err := doc.Get(ctx, r.Client, importPath, "")
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors = multierror.Append(errors, fmt.Errorf("failed getting %s: %s", importPath, err))
				return
			}
			if sp.Name == "" {
				return
			}
			pkg.SubPackages = append(pkg.SubPackages, subPkg{
				Path: subDir,
				Doc:  sp.Synopsis,
			})
		}(subDir)
	}
	wg.Wait()
	return pkg, errors.ErrorOrNil()
}

var tmpl = template.Must(template.New("readme").Funcs(
	template.FuncMap{
		"code": func(s string) string { return "```golang\n" + s + "\n```\n" },
	},
).Parse(`
# Package {{.Package.Name}}

	go get {{.Package.ImportPath}}

{{.Package.Doc}}

{{if .SubPackages}}
## Sub Packages
{{range .SubPackages}}
* [{{.Path}}](./{{.Path}}){{if .Doc}}: {{.Doc}}{{end}}
{{end}}
{{end}}

{{if .Package.Examples}}
## Examples
{{range .Package.Examples}}
### {{.Name}}

{{.Doc}}

{{code .Play}}
{{end}}
{{end}}
`))
