package goreadme

import "text/template"

var tmpl = template.Must(template.New("readme").Funcs(
	template.FuncMap{
		"code": func(s string) string { return "```golang\n" + s + "\n```\n" },
	},
).Parse(`# {{.Package.Name}}

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
