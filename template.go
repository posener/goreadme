package goreadme

import (
	"strings"
	"text/template"

	"github.com/golang/gddo/doc"
)

var tmpl = template.Must(template.New("readme").Funcs(
	template.FuncMap{
		"code": func(s string) string {
			return "```golang\n" + s + "\n```\n"
		},
		"inlineCode": func(s string) string {
			return "`" + s + "`"
		},
		"fullName": func(p *doc.Package) string {
			return strings.TrimPrefix(p.ImportPath, "github.com/")
		},
	},
).Parse(`# {{.Package.Name}}

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
{{ .Package.Doc }}

{{ if (and .Config.Functions .Package.Funcs) }}
## Functions

{{ range .Package.Funcs }}

### {{ .Name }}

{{ inlineCode .Decl.Text }}

{{ .Doc }}

{{ if .Examples }}
#### Examples

{{ range .Examples }}

{{ if .Name }}
##### {{.Name}}
{{ end }}

{{ .Doc }}

{{ if .Play }}{{code .Play}}{{ else }}{{code .Code.Text}}
{{ end }}
{{ end }}
{{ end }}
{{ end }}
{{ end }}

{{ if (and .SubPackages (not .Config.SkipSubPackages)) }}
## Sub Packages

{{ range .SubPackages }}
* [{{.Path}}](./{{.Path}}){{if .Package.Synopsis}}: {{.Package.Synopsis}}{{end}}
{{ end }}
{{ end }}

{{ if (and .Package.Examples (not .Config.SkipExamples)) }}
## Examples

{{ range .Package.Examples }}
### {{.Name}}

{{ .Doc }}

{{ if .Play }}{{code .Play}}{{ else }}{{code .Code.Text}}
{{end }}
{{end }}
{{end }}
`))
