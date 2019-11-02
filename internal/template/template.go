package template

import (
	"bytes"
	"io"
	"strings"
	"text/template"

	"github.com/golang/gddo/doc"
	"github.com/posener/goreadme/internal/markdown"
)

// Execute is used to execute the README.md template
func Execute(w io.Writer, data interface{}) error {
	return main.Execute(&multiNewLineEliminator{w: w}, data)
}

var base = template.New("base").Funcs(
	template.FuncMap{
		"gocode": func(s string) string {
			return "```golang\n" + s + "\n```\n"
		},
		"code": func(s string) string {
			return "```\n" + s + "\n```\n"
		},
		"inlineCode": func(s string) string {
			return "`" + s + "`"
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
		"doc": func(s string) string {
			b := bytes.NewBuffer(nil)
			markdown.ToMarkdown(b, s, nil)
			return b.String()
		},
	},
)

var main = template.Must(base.Parse(`# {{.Package.Name}}

{{if .Config.Badges.TravicCI -}}
[![Build Status](https://travis-ci.org/{{fullName .Package}}.svg?branch=master)](https://travis-ci.org/{{fullName .Package}})
{{end -}}
{{if .Config.Badges.CodeCov -}}
[![codecov](https://codecov.io/gh/{{fullName .Package}}/branch/master/graph/badge.svg)](https://codecov.io/gh/{{fullName .Package}})
{{end -}}
{{if .Config.Badges.GolangCI -}}
[![golangci](https://golangci.com/badges/{{.Package.ImportPath}}.svg)](https://golangci.com/r/{{.Package.ImportPath}})
{{end -}}
{{if .Config.Badges.GoDoc -}}
[![GoDoc](https://godoc.org/{{.Package.ImportPath}}?status.svg)](http://godoc.org/{{.Package.ImportPath}})
{{end -}}
{{if .Config.Badges.GoReportCard -}}
[![Go Report Card](https://goreportcard.com/badge/{{.Package.ImportPath}})](https://goreportcard.com/report/{{.Package.ImportPath}})
{{end -}}
{{if .Config.Badges.Goreadme -}}
[![goreadme](https://goreadme.herokuapp.com/badge/{{fullName .Package}}.svg)](https://goreadme.herokuapp.com)
{{ end }}

{{ doc .Package.Doc }}

{{ if .Config.Functions }}
{{ template "functions" .Package }}
{{ end }}

{{ if (not .Config.SkipSubPackages) }}
{{ template "subpackages" . }}
{{ end }}

{{ if (not .Config.SkipExamples) }}
{{ template "examples" .Package.Examples }}
{{end }}
`))

var functions = template.Must(base.Parse(`
{{ define "functions" }}
{{ if .Funcs }}

## Functions

{{ range .Funcs }}

### func [{{ .Name }}]({{ urlOrName (index $.Files .Pos.File) }}#L{{ .Pos.Line }})

{{ inlineCode .Decl.Text }}

{{ doc .Doc }}

{{ template "examples" .Examples }}
{{ end }}

{{ end }}
{{ end }}
`))

var exmaples = template.Must(base.Parse(`
{{ define "examples" }}
{{ if . }}

#### Examples

{{ range . }}

{{ if .Name }}##### {{.Name}}{{ end }}

{{ doc .Doc }}

{{ if .Play }}{{gocode .Play}}{{ else }}{{gocode .Code.Text}}{{ end }}
{{ if .Output }} Output:

{{ code .Output }}{{ end }}
{{ end }}

{{ end }}
{{ end }}
`))

var subPackages = template.Must(base.Parse(`
{{ define "subpackages" }}
{{ if .SubPackages }}

## Sub Packages

{{ range .SubPackages }}
* [{{.Path}}](./{{.Path}}){{if .Package.Synopsis}}: {{.Package.Synopsis}}{{end}}
{{ end }}

{{ end }}
{{ end }}
`))
