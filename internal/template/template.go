package template

import (
	"bytes"
	"embed"
	"io"
	"regexp"
	"strings"
	"text/template"

	"github.com/golang/gddo/doc"
	"github.com/posener/goreadme/internal/config"
	"github.com/posener/goreadme/internal/markdown"
)

//go:embed *.md.gotmpl
var files embed.FS

// Execute is used to execute the README.md template.
func Execute(w io.Writer, data interface{}, cfg config.Config, options ...markdown.Option) error {
	templates, err := template.New("main.md.gotmpl").Funcs(funcs(cfg, options)).ParseFS(files, "*")
	if err != nil {
		return err
	}
	return templates.Execute(&multiNewLineEliminator{w: w}, data)
}

func funcs(cfg config.Config, options []markdown.Option) template.FuncMap {
	return template.FuncMap{
		"config": func() config.Config {
			return cfg
		},
		"doc": func(s string) string {
			b := bytes.NewBuffer(nil)
			markdown.ToMarkdown(b, s, options...)
			return b.String()
		},
		"gocode": func(s string) string {
			return "```golang\n" + s + "\n```\n"
		},
		"code": func(s string) string {
			if !strings.HasSuffix(s, "\n") {
				s = s + "\n"
			}
			return "```\n" + s + "```\n"
		},
		"inlineCode": func(s string) string {
			return "`" + s + "`"
		},
		"inlineCodeEllipsis": func(s string) string {
			r := regexp.MustCompile(`{(?s).*}`)
			s = r.ReplaceAllString(s, "{ ... }")
			return "`" + s + "`"
		},
		"importPath": func(p *doc.Package) string {
			return p.ImportPath
		},
		"fullName": func(p *doc.Package) string {
			return strings.TrimPrefix(p.ImportPath, "github.com/")
		},
		"urlOrName": func(f *doc.File) string {
			if f.URL != "" {
				return f.URL
			}
			return "/" + f.Name
		},
	}
}
