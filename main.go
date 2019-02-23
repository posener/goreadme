package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/gorilla/mux"
	"github.com/posener/goreadme/goreadme"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	hookName          = "goreadme"
	domain            = "https://goreadme.herokuapp.com"
	githubAppURL      = "https://github.com/apps/goreadme"
	timeout           = time.Second * 60
	defaultReadmePath = "README.md"

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
		github:   github.NewClient(client),
		goreadme: &goreadme.GoReadme{Client: client},
	}
	m := mux.NewRouter()
	m.Methods("GET").Path("/").HandlerFunc(h.home)
	m.Methods("POST").Path("/github/hook").HandlerFunc(h.hook)
	logrus.Infof("Starting server...")
	http.ListenAndServe(":"+port, m)
}

type handler struct {
	github   *github.Client
	goreadme *goreadme.GoReadme
}

func (h *handler) home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, githubAppURL, http.StatusFound)
}

// hook is called by github when there is a push to repository.
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

	log.Info("Running goreadme in background...")
	go h.runPR(log, &push)
}

func (h *handler) runPR(log logrus.FieldLogger, push *github.PushEvent) {
	var (
		owner         = push.GetRepo().GetOwner().GetName()
		repo          = push.GetRepo().GetName()
		headSHA       = push.GetHeadCommit().GetID()
		defaultBranch = push.GetRepo().GetDefaultBranch()
		ref           = push.GetRef()
	)

	log.WithFields(logrus.Fields{"sha": headSHA, "default branch": defaultBranch}).Infof("Running PR process")

	b := bytes.NewBuffer(nil)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Create new readme for repository.
	err := h.goreadme.Create(ctx, "github.com/"+owner+"/"+repo, b)
	if err != nil {
		log.Errorf("Failed goreadme: %s", err)
		return
	}

	// Check for changes from current readme
	readmePath := defaultReadmePath
	readme, resp, err := h.github.Repositories.GetReadme(ctx, owner, repo, &github.RepositoryContentGetOptions{
		Ref: ref,
	})
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

	// Reset goreadme branch - delete it if exists and then create it.
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

	// Commit changes to readme file.
	date := time.Now()
	author := &github.CommitAuthor{
		Name:  github.String(goreadmeAuthor),
		Email: github.String(goreadmeAuthor + "@gmail.com"),
		Date:  &date,
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

	// Create pull request
	pr, _, err := h.github.PullRequests.Create(ctx, owner, repo, &github.NewPullRequest{
		Title: github.String("readme: Update according to go doc"),
		Base:  github.String(defaultBranch),
		Head:  github.String(goreadmeBranch),
	})
	if err != nil {
		log.Errorf("Failed creatring PR: %s", err)
		return
	}

	log.Infof("Created PR: %s", pr.GetHTMLURL())
}

func branch(push *github.PushEvent) string {
	return strings.TrimPrefix(push.GetRef(), "refs/heads/")
}