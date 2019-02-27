package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/posener/goreadme/templates"
	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/posener/goreadme/auth"
	"github.com/posener/goreadme/goreadme"
	"github.com/sirupsen/logrus"
)

var githubHookSecret = []byte(os.Getenv("GITHUB_HOOK_SECRET")) // Secret for github hooks

type handler struct {
	auth     *auth.Auth
	db       *gorm.DB
	github   *github.Client
	goreadme *goreadme.GoReadme
}

func (h *handler) home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, githubAppURL, http.StatusFound)
}

// hook is called by github when there is a push to repository.
func (h *handler) hook(w http.ResponseWriter, r *http.Request) {
	body, err := github.ValidatePayload(r, githubHookSecret)
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
	var data struct {
		templates.Base
		Jobs []Job
	}
	err := h.db.Model(&Job{}).Order("updated_at DESC").Scan(&data.Jobs).Error
	if err != nil {
		logrus.Errorf("Failed scanning jobs: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	data.User = h.auth.User(r)
	err = templates.JobsList.Execute(w, data)
	if err != nil {
		logrus.Errorf("Failed executing template: %s", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// debugPR runs in debug mode provide the required environment variables.
// Run with:
//
// 		DEBUG_HOOK=1 REPO=repo OWNER=$USER HEAD=$(git rev-parse HEAD) go run .
//
func (h *handler) debugPR() {
	if os.Getenv("DEBUG_HOOK") != "1" {
		return
	}
	logrus.Warnf("Debugging hook mode!")
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
