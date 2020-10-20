package template

import (
	"bytes"
	"io"
	"regexp"
	"strings"
	"text/template"

	"github.com/golang/gddo/doc"
	"github.com/posener/goreadme/internal/markdown"
)

// Execute is used to execute the README.md template.
func Execute(w io.Writer, data interface{}) error {
	return main.Execute(&multiNewLineEliminator{w: w}, data)
}

var base = template.New("base").Funcs(
	template.FuncMap{
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
			r := regexp.MustCompile(`\{.*\}`)
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
		"doc": func(s string) string {
			b := bytes.NewBuffer(nil)
			markdown.ToMarkdown(b, s, nil)
			return b.String()
		},
	},
)

var main = template.Must(base.Parse(`# {{.Package.Name}}

{{if .Config.Badges.TravisCI -}}
[![Build Status](https://travis-ci.org/{{fullName .Package}}.svg?branch=master)](https://travis-ci.org/{{fullName .Package}})
{{end -}}
{{if .Config.Badges.CodeCov -}}
[![codecov](https://codecov.io/gh/{{fullName .Package}}/branch/master/graph/badge.svg)](https://codecov.io/gh/{{fullName .Package}})
{{end -}}
{{if .Config.Badges.GolangCI -}}
[![golangci](https://golangci.com/badges/{{importPath .Package}}.svg)](https://golangci.com/r/{{importPath .Package}})
{{end -}}
{{if .Config.Badges.GoDoc -}}
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/{{importPath .Package}})
{{end -}}
{{if .Config.Badges.GoReportCard -}}
[![Go Report Card](https://goreportcard.com/badge/{{importPath .Package}})](https://goreportcard.com/report/{{importPath .Package}})
{{ end }}

{{ doc .Package.Doc }}

{{ if .Config.Functions }}
{{ template "functions" .Package }}
{{ end }}

{{ if .Config.Types }}
{{ template "types" .Package }}
{{ end }}

{{ if (not .Config.SkipSubPackages) }}
{{ template "subpackages" . }}
{{ end }}

{{ if (not .Config.SkipExamples) }}
{{ template "examples" .Package.Examples }}
{{ end }}
{{ if .Config.Credit }}
---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
{{ end }}
`))

var functions = template.Must(base.Parse(`
{{ define "functions" }}
{{ if .Funcs }}

## Functions

{{ range .Funcs }}

### func [{{ .Name }}]({{ urlOrName (index $.Files .Pos.File) }}#L{{ .Pos.Line }})

{{ inlineCode .Decl.Text }}

{{ doc .Doc }}

{{ template "examplesNoHeading" .Examples }}
{{ end }}

{{ end }}
{{ end }}
`))

var types = template.Must(base.Parse(`
{{ define "types" }}
{{ if .Types }}

## Types

{{ range .Types }}

### type [{{ .Name }}]({{ urlOrName (index $.Files .Pos.File) }}#L{{ .Pos.Line }})

{{ inlineCodeEllipsis .Decl.Text }}

{{ doc .Doc }}

{{ template "examplesNoHeading" .Examples }}
{{ end }}

{{ end }}
{{ end }}
`))

var examples = template.Must(base.Parse(`
{{ define "examples" }}
{{ if . }}

## Examples

{{ template "examplesNoHeading" . }}

{{ end }}
{{ end }}
`))

var examplesNoHeading = template.Must(base.Parse(`
{{ define "examplesNoHeading" }}
{{ if . }}

{{ range . }}

{{ if .Name }}### {{.Name}}{{ end }}

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
