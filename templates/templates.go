package templates

import (
	"html/template"

	"github.com/google/go-github/github"
)

type Base struct {
	User *github.User
}

var base = template.Must(template.New("base").Parse(`
{{define "base"}}
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Goreadme</title>
  <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
</head>
<body>

<nav class="navbar navbar-expand-md navbar-light">
    <a class="navbar-brand abs" href="#">Goreadme</a>
    <button class="navbar-toggler" type="button" data-toggle="collapse" data-target="#collapsingNavbar">
        <span class="navbar-toggler-icon"></span>
    </button>
    <div class="navbar-collapse collapse" id="collapsingNavbar">
        <ul class="navbar-nav">
            <li class="nav-item {{if .Jobs}}active{{end}}">
                <a class="nav-link" href="/jobs">Jobs</a>
            </li>
        </ul>
        <ul class="navbar-nav ml-auto">
			<li class="nav-item dropdown">
				<a class="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
					<img src="{{.User.GetAvatarURL}}" width="30" height="30" class="d-inline-block align-top" alt="">
					{{.User.GetLogin}}
				</a>
				<div class="dropdown-menu" aria-labelledby="navbarDropdown">
					<a class="dropdown-item" href="{{.User.GetHTMLURL}}">Github page</a>
					<a class="dropdown-item" href="/github/logout">Logout</a>
				</div>
			</li>
        </ul>
    </div>
</nav>

  <!-- <h1 class="text-center">{{block "title" .}}{{end}}</h1> -->

  {{template "content" .}}

  <script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.14.7/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
  <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.3.1/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>
</body>
</html>
{{end}}
`))

var JobsList = template.Must(template.Must(base.Clone()).Parse(`
{{define "title"}}Jobs List{{end}}
{{define "content"}}
<table class="table">
<thead>
	<tr>
	<th scope="col">Repository</th>
	<th scope="col">Job #</th>
	<th scope="col">Status</th>
	<th scope="col">Message</th>
	<th scope="col">Created</th>
	<th scope="col">Updated</th>
	</tr>
</thead>
<tbody>
	{{range .Jobs}}
	<tr>
		<th scope="row">{{.Owner}}/{{.Repo}}</th>
		<td>{{.Num}}</td>
		<td>{{if .PRURL}}<a href="{{.PRURL}}">{{end}}{{.Status}}{{if .PRURL}}</a>{{end}}</td>
		<td>{{.Message}}</td>
		<td>{{.CreatedAt}}</td>
		<td>{{.UpdatedAt}}</td>
	</tr>
	{{end}}
</tbody>
</table>
{{end}}
`))
