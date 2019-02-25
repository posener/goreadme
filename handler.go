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
	"github.com/pkg/errors"
	"github.com/posener/goreadme/goreadme"
	"github.com/sirupsen/logrus"
	"github.com/src-d/go-git/plumbing"
)

type handler struct {
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

	log.WithFields(logrus.Fields{"sha": headSHA[:8], "default branch": defaultBranch}).Infof("Running PR process")

	b := bytes.NewBuffer(nil)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Get config
	cfg, err := h.getConfig(ctx, owner, repo)
	if err != nil {
		log.Errorf("Failed getting config: %s", err)
		return
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
	var currentContent string
	switch {
	case resp.StatusCode == http.StatusNotFound:
		log.Infof("No current readme, creating a new readme!")
	case err != nil:
		log.Errorf("Failed getting upstream readme: %s", err)
		return
	default:
		currentContent, err = readme.GetContent()
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
	_, resp, err = h.github.Repositories.GetBranch(ctx, owner, repo, goreadmeBranch)
	switch {
	case resp.StatusCode != http.StatusNotFound:
		log.Infof("Found existing branch, deleting...")
		_, err = h.github.Git.DeleteRef(ctx, owner, repo, goreaedmeRef)
		if err != nil {
			log.Errorf("Failed deleting existing branch: %s", err)
			return
		}
	case err != nil:
		log.Errorf("Failed getting branch: %s", err)
		return
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
		Message:   github.String("Update readme according to go doc"),
		SHA:       github.String(plumbing.ComputeHash(plumbing.BlobObject, []byte(currentContent)).String()),
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

func (h *handler) getConfig(ctx context.Context, repo, owner string) (goreadme.Config, error) {
	var cfg goreadme.Config
	cfgContent, _, resp, err := h.github.Repositories.GetContents(ctx, owner, repo, configPath, nil)
	switch {
	case resp.StatusCode == http.StatusNotFound:
	case err != nil:
		return cfg, errors.Wrap(err, "failed get config file")
	default:
		content, err := cfgContent.GetContent()
		if err != nil {
			return cfg, errors.Wrap(err, "failed get config content")
		}
		err = json.Unmarshal([]byte(content), &cfg)
		if err != nil {
			return cfg, errors.Wrapf(err, "unmarshaling config content %s", content)
		}
	}
	return cfg, err

}

func branch(push *github.PushEvent) string {
	return strings.TrimPrefix(push.GetRef(), "refs/heads/")
}

const credits = "\nCreated by [goreadme](" + githubAppURL + ")\n"

// debugPR runs in debug mode provide the required environment variables.
// Run with:
//
// 		DEBUG=1 REPO=repo OWNER=$USER HEAD=$(git rev-parse HEAD) go run .
//
func (h *handler) debugPR() {
	if os.Getenv("DEBUG") != "1" {
		return
	}
	h.runPR(logrus.StandardLogger(), &github.PushEvent{
		Repo: &github.PushEventRepository{
			Name:          github.String(os.Getenv("REPO")),
			Owner:         &github.PushEventRepoOwner{Name: github.String(os.Getenv("OWNER"))},
			DefaultBranch: github.String("master"),
		},
		HeadCommit: &github.PushEventCommit{ID: github.String(os.Getenv("HEAD"))},
	})
	os.Exit(0)
}
