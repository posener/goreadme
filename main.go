// Package main is an HTTP server that works with Github hooks.
//
// [goreadme](./goreadme) is a tool for creating README.md files from Go doc
// of a given package.
// This server provides Github automation on top of this tool, but creating
// PRs for your github repository, whenever the README file should be updated.
//
// ## Usage
//
// Go to [https://github.com/apps/goreadme](https://github.com/apps/goreadme)
// Press the "Configure" button, choose your account, and add the repositories
// you want goreadme to maintain for you.
//
// ## How does it Work?
//
// Once enabled, goreadme is registered on a Github hook, that calls goreadme
// server the repository default branch is modified.
// Goreadme then computes the new README.md file and compairs it to the exiting
// one. If a change is needed, Goreadme will create a PR with the new content
// of the README.md file.
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
	githubAppURL      = "https://github.com/apps/goreadme"
	timeout           = time.Second * 60
	defaultReadmePath = "README.md"
	configPath        = "goreadme.json"

	goreadmeAuthor = "goreadme"
	goreadmeEmail = "posener@gmail.com"
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
		goreadme: goreadme.New(client),
	}
	m := mux.NewRouter()
	m.Methods("GET").Path("/").HandlerFunc(h.home)
	m.Methods("POST").Path("/github/hook").Handler(auth(http.HandlerFunc(h.hook)))
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

	branch := branch(&push)
	log := logrus.WithField("repo", push.GetRepo().GetFullName())
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
	)

	log.WithFields(logrus.Fields{"sha": headSHA, "default branch": defaultBranch}).Infof("Running PR process")

	b := bytes.NewBuffer(nil)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Get config
	var cfg goreadme.Config
	cfgContent, _, resp, err := h.github.Repositories.GetContents(ctx, owner, repo, configPath, nil)
	switch {
	case resp.StatusCode == http.StatusNotFound:
	case err != nil:
		log.Errorf("Failed get config file: %s", err)
		return
	default:
		content, err := cfgContent.GetContent()
		if err != nil {
			log.Errorf("Failed get config content: %s", err)
			return
		}
		err = json.Unmarshal([]byte(content), &cfg)
		if err != nil {
			log.Errorf("Failed unmarshaling config content %s: %s", content, err)
			return
		}
	}

	// Create new readme for repository.
	err = h.goreadme.WithConfig(cfg).Create(ctx, "github.com/"+owner+"/"+repo, b)
	if err != nil {
		log.Errorf("Failed goreadme: %s", err)
		return
	}

	b.WriteString(credits)

	// Check for changes from current readme
	readmePath := defaultReadmePath
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
		Email: github.String(goreadmeEmail),
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

const credits = "\nCreated by [goreadme](" + githubAppURL + ")\n"
