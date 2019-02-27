package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/posener/goreadme/auth"
	"github.com/posener/goreadme/goreadme"
	"github.com/posener/goreadme/internal/templates"
	"github.com/sirupsen/logrus"
)

var githubHookSecret = []byte(os.Getenv("GITHUB_HOOK_SECRET")) // Secret for github hooks

type handler struct {
	auth     *auth.Auth
	db       *gorm.DB
	github   *github.Client
	goreadme *goreadme.GoReadme
}

type templateData struct {
	User  *github.User
	Repos []*github.Repository
	States []State
	Jobs  []Job
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
		log.Infof("Skipping push to non-default branch %s", br)
		return
	}

	log.Info("Running goreadme in background...")
	go h.runJob(&Job{
		State: State{
			Owner:   push.GetRepo().GetOwner().GetName(),
			Repo:    push.GetRepo().GetName(),
			HeadSHA: push.GetHeadCommit().GetID(),
		},
		defaultBranch: push.GetRepo().GetDefaultBranch(),
	})
}

func (h *handler) runJob(j *Job) {
	j.db = h.db
	j.github = h.github
	j.goreadme = h.goreadme
	j.Run()
}

func (h *handler) home(w http.ResponseWriter, r *http.Request) {
	var data templateData
	data.User = h.auth.User(r)
	login := data.User.GetLogin()
	owners, err := h.getOwners(r.Context(), login)
	if err != nil {
		h.doError(err, w, r)
		return
	}

	err = h.db.Model(&State{}).Order("updated_at DESC").Where("owner IN ( ? )", owners).Scan(&data.States).Error
	if err != nil {
		h.doError(errors.Wrap(err, "failed scanning states"), w, r)
		return
	}
	
	err = templates.Home.Execute(w, data)
	if err != nil {
		h.doError(errors.Wrap(err, "failed executing template"), w, r)
	}
}

func (h *handler) login(w http.ResponseWriter, r *http.Request) {
	err := templates.Login.Execute(w, templateData{})
	if err != nil {
		h.doError(errors.Wrap(err, "failed executing template"), w, r)
	}
}

func (h *handler) jobsList(w http.ResponseWriter, r *http.Request) {
	var data templateData
	data.User = h.auth.User(r)

	login := data.User.GetLogin()
	owners, err := h.getOwners(r.Context(), login)
	if err != nil {
		h.doError(err, w, r)
		return
	}
	err = h.db.Model(&Job{}).Order("updated_at DESC").Where("owner IN ( ? )", owners).Scan(&data.Jobs).Error
	if err != nil {
		h.doError(errors.Wrap(err, "failed scanning jobs"), w, r)
		return
	}
	err = templates.JobsList.Execute(w, data)
	if err != nil {
		h.doError(errors.Wrap(err, "failed executing template"), w, r)
	}
}

func (h *handler) addRepo(w http.ResponseWriter, r *http.Request) {
	// repos, _, err := h.github.Repositories.List(r.Context(), data.User.GetLogin(), nil)
	// if err != nil {
	// 	h.doError(errors.Wrapf(err, "failed getting repo for user %s", data.User.GetLogin()), w, r)
	// 	return
	// }
	// data.Repos = repos
}

func (h *handler) doError(err error, w http.ResponseWriter, r *http.Request) {
	logrus.Error(err)
	http.Redirect(w, r, "/?error=internal%20server%error", http.StatusFound)
}

func (h *handler) getOwners(ctx context.Context, login string) ([]string, error) {
	orgs, _, err := h.github.Organizations.List(ctx, login, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed listing orgs")
	}
	owners := make([]string, 0, len(orgs)+1)
	owners = append(owners, login)
	for _, o := range orgs {
		owners = append(owners, o.GetName())
	}
	return owners, nil
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
		State: State{
			Owner:   os.Getenv("OWNER"),
			Repo:    os.Getenv("REPO"),
			HeadSHA: os.Getenv("HEAD"),
		},
		defaultBranch: "master",
	})
	os.Exit(0)
}

func branch(push *github.PushEvent) string {
	return strings.TrimPrefix(push.GetRef(), "refs/heads/")
}
