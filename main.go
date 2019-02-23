package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"github.com/posener/goreadme/auth"
	"github.com/posener/goreadme/goreadme"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	hookName = "goreadme"
	domain   = "https://goreadme.herokuapp.com"
	timeout  = time.Second * 60

	goreadmeAuthor = "goreadme"
	goreadmeBranch = "goreadme"
	goreaedmeRef   = "refs/heads/" + goreadmeBranch
)

var (
	port        = os.Getenv("PORT")
	githubToken = os.Getenv("GITHUB_TOKEN")
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)

	client := oauth2.NewClient(ctx, ts)
	h := &handler{
		client: client,
		github: github.NewClient(client),
	}
	callback, login := auth.Handlers(domain)
	m := mux.NewRouter()
	m.Methods("GET").Path("/").HandlerFunc(h.home)
	m.Methods("POST").Path("/github/hook").HandlerFunc(h.hook)
	m.Path("/github/callback").Handler(callback)
	m.Path("/github/login").Handler(login)
	logrus.Infof("Starting server...")
	http.ListenAndServe(":"+port, m)
}

type handler struct {
	client *http.Client
	github *github.Client
}

func (h *handler) home(w http.ResponseWriter, r *http.Request) {
	if !auth.IsAuthenticated(r) {
		http.Redirect(w, r, "/github/login", http.StatusFound)
		return
	}
	fmt.Fprintf(w, "<body><h1>Logged in as %s</h1></body>", auth.ID(r))
}

func (h *handler) hook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var push github.PushEvent
	err := json.NewDecoder(r.Body).Decode(&push)
	if err != nil {
		logrus.Errorf("Failed decoding push event: %s", err)
		http.Error(w, "Failed", 500)
		return
	}

	fullName := push.GetRepo().GetFullName()
	branch := branch(&push)
	log := logrus.WithField("repo", fullName)
	log.Infof("Got push event to %s", branch)
	if branch != push.GetRepo().GetDefaultBranch() {
		log.Infof("Skipping push to non-default branch %s", branch)
		return
	}

	if author := push.GetHeadCommit().GetAuthor().GetName(); author == goreadmeAuthor {
		log.Infof("Skipping check of change by %s", author)
		return
	}

	log.Info("Running goreadme in background...")
	go h.runPR(log, &push)
}

func (h *handler) runPR(log logrus.FieldLogger, push *github.PushEvent) {

	var (
		owner         = push.GetRepo().GetOwner().GetName()
		repo          = push.GetRepo().GetName()
		headSHA       = push.GetHeadCommit().GetID()
		defaultBranch = push.GetRepo().GetDefaultBranch()
	)

	log.WithFields(logrus.Fields{"sha": headSHA, "default branch": defaultBranch}).Infof("Running PR process")

	b := bytes.NewBuffer(nil)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := goreadme.Create(ctx, h.client, "github.com/"+owner+"/"+repo, b)
	if err != nil {
		log.Errorf("Failed goreadme: %s", err)
		return
	}

	readmePath := "README.md"
	readme, resp, err := h.github.Repositories.GetReadme(ctx, owner, repo, nil)
	switch {
	case resp.StatusCode == http.StatusNotFound:
		log.Infof("No current readme")
	case err != nil:
		log.Errorf("Failed getting upstream readme: %s", err)
		return
	default:
		currentContent, err := readme.GetContent()
		if err != nil {
			log.Errorf("Failed get readme content: %s", err)
			return
		}
		if currentContent == b.String() {
			log.Infof("Done! Nothing to change...")
			return
		}
		readmePath = readme.GetPath()
	}

	date := time.Now()
	author := &github.CommitAuthor{
		Name:  github.String(goreadmeAuthor),
		Email: github.String(goreadmeAuthor + "@gmail.com"),
		Date:  &date,
	}

	_, resp, _ = h.github.Repositories.GetBranch(ctx, owner, repo, goreadmeBranch)
	if resp.StatusCode != http.StatusNotFound {
		_, err = h.github.Git.DeleteRef(ctx, owner, repo, goreaedmeRef)
		if err != nil {
			log.Errorf("Failed deleting existing branch: %s", err)
			return
		}
	}

	_, _, err = h.github.Git.CreateRef(ctx, owner, repo, &github.Reference{
		Ref:    github.String(goreaedmeRef),
		Object: &github.GitObject{SHA: github.String(headSHA)},
	})
	if err != nil {
		log.Errorf("Failed creating branch: %s", err)
		return
	}

	_, _, err = h.github.Repositories.UpdateFile(ctx, owner, repo, readmePath, &github.RepositoryContentFileOptions{
		Author:    author,
		Committer: author,
		Branch:    github.String(goreadmeBranch),
		Content:   b.Bytes(),
		Message:   github.String("update readme according to go doc"),
		SHA:       github.String(headSHA),
	})
	if err != nil {
		log.Errorf("Failed updating readme content: %s", err)
		return
	}

	pr, _, err := h.github.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
		Title: github.String("readme: Update according to go doc"),
		Base:  github.String(defaultBranch),
		Head:  github.String(goreadmeBranch),
	})
	if err != nil {
		log.Errorf("Failed creatring PR: %s", err)
		return
	}

	log.Infof("Created PR: %s", *pr.URL)
}

func branch(push *github.PushEvent) string {
	return strings.TrimPrefix(push.GetRef(), "refs/heads/")
}

func (h *handler) createHook(ctx context.Context, owner, repo string) {
	_, _, err := h.github.Repositories.CreateHook(ctx, owner, repo, &github.Hook{
		Name:   github.String(hookName),
		Active: github.Bool(true),
		Events: []string{"push"},
		URL:    github.String(domain + "/github/hook"),
	})
	if err != nil {
		logrus.Errorf("Failed creating hook %s/%s: %s", owner, repo, err)
	}
}

const home = `
<html>
  <head>
  </head>
  <body>
    <p>
      Well, hello there!
    </p>
    <p>
      We're going to now talk to the GitHub API. Ready?
      <a href="https://github.com/login/oauth/authorize?scope=user:email&client_id=%s">Click here</a> to begin!</a>
    </p>
    <p>
      If that link doesn't work, remember to provide your own <a href="/apps/building-oauth-apps/authorizing-oauth-apps/">Client ID</a>!
    </p>
  </body>
</html>
`
