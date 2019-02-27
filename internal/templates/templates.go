package templates

import (
	"html/template"
	"time"
)

var html = template.Must(
	template.New("html").Funcs(
		template.FuncMap{
			"formatDate": func(t time.Time) string { return t.Format("2006-01-02T15:04:05") },
		}).Parse(`
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Goreadme</title>
  <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
</head>
<body>

{{template "body" .}}

  <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
  <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
</body>
</html>
`))

var base = template.Must(html.Parse(`
{{define "body"}}
<nav class="navbar navbar-expand-md navbar-light bg-light">
	<a class="navbar-brand abs" href="/">
		<img src="https://avatars3.githubusercontent.com/in/25929?s=30&u=0a3756b6a47f20c14b650528b9a477a81ca5dd15&v=4" width="30" height="30" alt="">
		Goreadme
	</a>
    <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#collapsingNavbar">
        <span class="navbar-toggler-icon"></span>
    </button>
    <div class="navbar-collapse collapse" id="collapsingNavbar">
        <ul class="navbar-nav">
            <li class="nav-item {{if .Jobs}}active{{end}}">
                <a class="nav-link" href="/jobs">Jobs</a>
            </li>
		</ul>
		{{if .User}}
        <ul class="navbar-nav ml-auto">
			<li class="nav-item dropdown">
				<a class="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
					<img src="{{.User.GetAvatarURL}}" width="30" height="30" class="d-inline-block align-top" alt="">
					{{.User.GetLogin}}
				</a>
				<div class="dropdown-menu" aria-labelledby="navbarDropdown">
					<a class="dropdown-item" href="{{.User.GetHTMLURL}}">Github page</a>
					<a class="dropdown-item" href="/auth/logout">Logout</a>
				</div>
			</li>
		</ul>
		{{end}}
    </div>
</nav>

  <!-- <h1 class="text-center">{{block "title" .}}{{end}}</h1> -->

  <div class="container">
  {{template "content" .}}
  </div>
{{end}}
`))

var Home = template.Must(base.Parse(`
{{define "title"}}Home{{end}}
{{define "content"}}
<h1>Your Readmes</h1>

<table class="table">
<thead>
	<tr>
	<th scope="col">Repository</th>
	<th scope="col">Last Status</th>
	<th scope="col">Message</th>
	<th scope="col">Job #</th>
	<th scope="col">Latest SHA</th>
	<th scope="col">Created</th>
	<th scope="col">Updated</th>
	<th scope="col">Trigger</th>
	</tr>
</thead>
<tbody>
	{{range .States}}
	<tr>
		<th scope="row">{{.Owner}}/{{.Repo}}</th>
		<td>{{if .PRURL}}<a href="{{.PRURL}}">{{end}}{{.Status}}{{if .PRURL}}</a>{{end}}</td>
		<td>{{.Message}}</td>
		<td>{{.Num}}</td>
		<td>{{.HeadSHA}}</td>
		<td>{{formatDate .CreatedAt}}</td>
		<td>{{formatDate .UpdatedAt}}</td>
		<td><a href="">Trigger</a></td>
	</tr>
	{{end}}
</tbody>
</table>


{{end}}
`))

var AddRepo = template.Must(base.Parse(`
{{define "title"}}Home{{end}}
{{define "content"}}
<div class="row">
	<div class="col-sm">
		<h1>Your Repositories</h1>
		<ul class="list-group list-group-flush">
			{{range .Repos}}
			<li class="list-group-item">
				<a href="{{.GetHTMLURL}}">{{.GetFullName}}</a>
				{{if .GetPrivate}}
				<span>P</span>
				{{end}}
			</li>
			{{end}}
		</ul>
	</div>
</div>
{{end}}
`))

var JobsList = template.Must(template.Must(base.Clone()).Parse(`
{{define "title"}}Jobs List{{end}}
{{define "content"}}
<table class="table">
<thead>
	<tr>
	<th scope="col">Repository</th>
	<th scope="col">Status</th>
	<th scope="col">Message</th>
	<th scope="col">Job #</th>
	<th scope="col">Created</th>
	<th scope="col">Updated</th>
	</tr>
</thead>
<tbody>
	{{range .Jobs}}
	<tr>
		<th scope="row">{{.Owner}}/{{.Repo}}</th>
		<td>{{if .PRURL}}<a href="{{.PRURL}}">{{end}}{{.Status}}{{if .PRURL}}</a>{{end}}</td>
		<td>{{.Message}}</td>
		<td>{{.Num}}</td>
		<td>{{formatDate .CreatedAt}}</td>
		<td>{{formatDate .UpdatedAt}}</td>
	</tr>
	{{end}}
</tbody>
</table>
{{end}}
`))

var Login = template.Must(base.Parse(`
{{define "title"}}Login{{end}}
{{define "content"}}
<a href="/auth/login">Please Login</a>
{{end}}
`))
