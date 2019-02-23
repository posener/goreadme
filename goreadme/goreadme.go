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

type Package struct {
	Package     *doc.Package
	SubPackages []SubPackage
}

type SubPackage struct {
	Path string
	Doc  string
}

func Create(ctx context.Context, client *http.Client, name string, w io.Writer) error {
	p, err := get(ctx, client, name)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, p)
}

func get(ctx context.Context, c *http.Client, name string) (*Package, error) {
	log.Printf("Getting %s", name)
	p, err := doc.Get(ctx, c, name, "")
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
	pkg := &Package{Package: p}
	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		errors *multierror.Error
	)
	wg.Add(len(p.Subdirectories))
	for _, subPkg := range p.Subdirectories {
		go func(subPkg string) {
			importPath := name + "/" + subPkg
			defer wg.Done()
			log.Printf("Getting %s", importPath)
			sp, err := doc.Get(ctx, c, importPath, "")
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors = multierror.Append(errors, fmt.Errorf("failed getting %s: %s", importPath, err))
				return
			}
			if sp.Name == "" {
				return
			}
			pkg.SubPackages = append(pkg.SubPackages, SubPackage{
				Path: subPkg,
				Doc:  sp.Synopsis,
			})
		}(subPkg)
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
