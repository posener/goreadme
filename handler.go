package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/posener/goreadme/goreadme"
	"github.com/sirupsen/logrus"
)

type handler struct {
	db       *gorm.DB
	github   *github.Client
	goreadme *goreadme.GoReadme
}

func (h *handler) home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, githubAppURL, http.StatusFound)
}

// hook is called by github when there is a push to repository.
func (h *handler) hook(w http.ResponseWriter, r *http.Request) {
	body, err := github.ValidatePayload(r, githubSecret)
	if err != nil {
		logrus.Warnf("Unauthorized request: %s", err)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var push github.PushEvent
	err = json.Unmarshal(body, &push)
	if err != nil {
		logrus.Errorf("Failed decoding push event: %s", err)
		http.Error(w, "Failed", 500)
		return
	}

	br := branch(&push)
	log := logrus.WithField("repo", push.GetRepo().GetFullName())
	log.Infof("Got push event to %s", br)
	if br != push.GetRepo().GetDefaultBranch() {
		log.Infof("Skipping push to non-default branch %s", branch)
		return
	}

	log.Info("Running goreadme in background...")
	go h.runJob(&Job{
		Owner:         push.GetRepo().GetOwner().GetName(),
		Repo:          push.GetRepo().GetName(),
		HeadSHA:       push.GetHeadCommit().GetID(),
		defaultBranch: push.GetRepo().GetDefaultBranch(),
	})
}

func (h *handler) runJob(j *Job) {
	j.db = h.db
	j.github = h.github
	j.goreadme = h.goreadme
	j.Run()
}

func (h *handler) jobsList(w http.ResponseWriter, r *http.Request) {
	var j []Job
	err := h.db.Model(&Job{}).Order("updated_at DESC").Scan(&j).Error
	if err != nil {
		logrus.Errorf("Failed scanning jobs: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = jobsList.Execute(w, j)
	if err != nil {
		logrus.Errorf("Failed executing template: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// debugPR runs in debug mode provide the required environment variables.
// Run with:
//
// 		DEBUG=1 REPO=repo OWNER=$USER HEAD=$(git rev-parse HEAD) go run .
//
func (h *handler) debugPR() {
	if os.Getenv("DEBUG") != "1" {
		return
	}
	h.runJob(&Job{
		Owner:         os.Getenv("OWNER"),
		Repo:          os.Getenv("REPO"),
		defaultBranch: "master",
		HeadSHA:       os.Getenv("HEAD"),
	})
	os.Exit(0)
}

func branch(push *github.PushEvent) string {
	return strings.TrimPrefix(push.GetRef(), "refs/heads/")
}

var jobsList = template.Must(template.New("jobs-list").Parse(`<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>ListenTo</title>
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/3.4.1/css/bootstrap.min.css" integrity="sha384-HSMxcRTRxnN+Bdg0JdbxYKrThecOKuH5zCYotlSAcp1+c8xmyTe9GYg1l9a69psu" crossorigin="anonymous">
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.3/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
  </head>
  <body>
	<h1 class="text-center">Jobs List</h1>
	<table class="table table-dark">
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
		{{range .}}
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
    <script src="https://code.jquery.com/jquery-1.12.4.min.js" integrity="sha384-nvAa0+6Qg9clwYCGGPpDQLVpLNn0fRaROjHqs13t4Ggj3Ez50XnGQqc/r8MhnRDZ" crossorigin="anonymous"></script>
    <script src="https://stackpath.bootstrapcdn.com/bootstrap/3.4.1/js/bootstrap.min.js" integrity="sha384-aJ21OjlMXNL5UyIl/XNwTMqvzeRMZH2w8c5cRVpzpU8Y5bApTppSuUkhZXN0VxHd" crossorigin="anonymous"></script>
  </body>
</html>`))
