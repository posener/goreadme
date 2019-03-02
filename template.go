package goreadme

import "text/template"

var tmpl = template.Must(template.New("readme").Funcs(
	template.FuncMap{
		"code": func(s string) string { return "```golang\n" + s + "\n```\n" },
	},
).Parse(`# {{.Package.Name}}

{{if .Config.Badges.TravicCI}}
[![Build Status](https://travis-ci.org/{{.Package.ProjectName}}.svg?branch=master)](https://travis-ci.org/{{.Package.ProjectName}}){{end -}}
{{if .Config.Badges.CodeCov}}
[![codecov](https://codecov.io/gh/{{.Package.ProjectName}}/branch/master/graph/badge.svg)](https://codecov.io/gh/{{.Package.ProjectName}}){{end -}}
{{if .Config.Badges.GolangCI}}
[![golangci](https://golangci.com/badges/github.com/{{.Package.ProjectName}}.svg)](https://golangci.com/r/github.com/{{.Package.ProjectName}}){{end -}}
{{if .Config.Badges.GoDoc}}
[![GoDoc](https://godoc.org/github.com/{{.Package.ProjectName}}?status.svg)](http://godoc.org/github.com/{{.Package.ProjectName}}){{end -}}
{{if .Config.Badges.GoReportCard}}
[![Go Report Card](https://goreportcard.com/badge/github.com/{{.Package.ProjectName}})](https://goreportcard.com/report/github.com/{{.Package.ProjectName}}){{end -}}

{{ .Package.Doc -}}
{{if (and .SubPackages (not .Config.SkipSubPackages)) }}

## Sub Packages
{{range .SubPackages}}
* [{{.Path}}](./{{.Path}}){{if .Package.Synopsis}}: {{.Package.Synopsis}}{{end}}
{{end -}}
{{end -}}
{{if (and .Package.Examples (not .Config.SkipExamples)) }}

## Examples

{{range .Package.Examples -}}
### {{.Name}}

{{ if .Doc }}{{ .Doc }}
{{ end -}}
{{ if .Play }}{{code .Play}}{{ else }}{{code .Code.Text}}
{{end -}}
{{end -}}
{{end -}}
`))
