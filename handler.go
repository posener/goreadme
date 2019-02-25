package main

import (
	"encoding/json"
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
