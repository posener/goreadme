package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/posener/goreadme/goreadme"
	"github.com/sirupsen/logrus"
	"github.com/src-d/go-git/plumbing"
)

const (
	githubAppURL      = "https://github.com/apps/goreadme"
	timeout           = time.Second * 60
	configPath        = "goreadme.json"
	defaultReadmePath = "README.md"

	goreadmeAuthor = "goreadme"
	goreadmeEmail  = "posener@gmail.com"
	goreadmeBranch = "goreadme"
	goreaedmeRef   = "refs/heads/" + goreadmeBranch
)

type Job struct {
	Num       int    `gorm:"primary_key"`
	Repo      string `gorm:"primary_key"`
	Owner     string `gorm:"primary_key"`
	HeadSHA   string
	PRURL     string
	Message   string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time

	defaultBranch string

	db       *gorm.DB
	github   *github.Client
	goreadme *goreadme.GoReadme
	log      logrus.FieldLogger
}

// Run runs the pull request flow
func (j *Job) Run() {
	err := j.init()
	if err != nil {
		j.log.Errorf("Failed creating job entry in database: %s", err)
		return
	}
	j.log.Infof("Starting PR process")

	b := bytes.NewBuffer(nil)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	saveError := func(format string, args ...interface{}) {
		s := fmt.Sprintf(format, args...)
		j.log.Error(s)
		j.Message = s
		j.Status = "Failed"
		err = j.db.Save(j).Error
		if err != nil {
			j.log.Errorf("Failed saving failed job: %s", err)
		}
	}

	saveSuccess := func(format string, args ...interface{}) {
		s := fmt.Sprintf(format, args...)
		j.log.Info(s)
		j.Message = s
		j.Status = "Success"
		err = j.db.Save(j).Error
		if err != nil {
			j.log.Errorf("Failed saving successful job: %s", err)
		}
	}

	// Get config
	cfg, err := j.getConfig(ctx)
	if err != nil {
		saveError("Failed getting config: %s", err)
		return
	}

	// Create new readme for repository.
	err = j.goreadme.WithConfig(cfg).Create(ctx, j.githubURL(), b)
	if err != nil {
		saveError("Failed goreadme: %s", err)
		return
	}
	b.WriteString(credits)

	// Check for changes from current readme
	readmePath := defaultReadmePath
	readme, resp, err := j.github.Repositories.GetReadme(ctx, j.Owner, j.Repo, nil)
	var currentContent string
	switch {
	case resp.StatusCode == http.StatusNotFound:
		j.log.Infof("No current readme, creating a new readme!")
	case err != nil:
		saveError("Failed getting upstream readme: %s", err)
		return
	default:
		currentContent, err = readme.GetContent()
		if err != nil {
			saveError("Failed get readme content: %s", err)
			return
		}
		if currentContent == b.String() {
			saveSuccess("No change needed")
			return
		}
		readmePath = readme.GetPath()
	}

	// Reset goreadme branch - delete it if exists and then create it.
	_, resp, err = j.github.Repositories.GetBranch(ctx, j.Owner, j.Repo, goreadmeBranch)
	switch {
	case resp.StatusCode == http.StatusNotFound:
		// Branch does not exist, we will create it later
	case err != nil:
		saveError("Failed getting branch: %s", err)
		return
	default:
		// Branch exist, delete it
		j.log.Infof("Found existing branch, deleting...")
		_, err = j.github.Git.DeleteRef(ctx, j.Owner, j.Repo, goreaedmeRef)
		if err != nil {
			saveError("Failed deleting existing branch: %s", err)
			return
		}
	}
	_, _, err = j.github.Git.CreateRef(ctx, j.Owner, j.Repo, &github.Reference{
		Ref:    github.String(goreaedmeRef),
		Object: &github.GitObject{SHA: github.String(j.HeadSHA)},
	})
	if err != nil {
		saveError("Failed creating branch: %s", err)
		return
	}

	// Commit changes to readme file.
	date := time.Now()
	author := &github.CommitAuthor{
		Name:  github.String(goreadmeAuthor),
		Email: github.String(goreadmeEmail),
		Date:  &date,
	}
	_, _, err = j.github.Repositories.UpdateFile(ctx, j.Owner, j.Repo, readmePath, &github.RepositoryContentFileOptions{
		Author:    author,
		Committer: author,
		Branch:    github.String(goreadmeBranch),
		Content:   b.Bytes(),
		Message:   github.String("Update readme according to go doc"),
		SHA:       github.String(plumbing.ComputeHash(plumbing.BlobObject, []byte(currentContent)).String()),
	})
	if err != nil {
		saveError("Failed updating readme content: %s", err)
		return
	}

	// Create pull request
	pr, _, err := j.github.PullRequests.Create(ctx, j.Owner, j.Repo, &github.NewPullRequest{
		Title: github.String("readme: Update according to go doc"),
		Base:  github.String(j.defaultBranch),
		Head:  github.String(goreadmeBranch),
	})
	if err != nil {
		saveError("Failed creatring PR: %s", err)
		return
	}

	j.log.Infof("Created PR: %s", pr.GetHTMLURL())

	j.PRURL = pr.GetHTMLURL()
	saveSuccess("Created PR")
}

func (j *Job) getConfig(ctx context.Context) (goreadme.Config, error) {
	var cfg goreadme.Config
	cfgContent, _, resp, err := j.github.Repositories.GetContents(ctx, j.Owner, j.Repo, configPath, nil)
	switch {
	case resp.StatusCode == http.StatusNotFound:
		return cfg, nil
	case err != nil:
		return cfg, errors.Wrap(err, "failed get config file")
	}
	content, err := cfgContent.GetContent()
	if err != nil {
		return cfg, errors.Wrap(err, "failed get config content")
	}
	err = json.Unmarshal([]byte(content), &cfg)
	if err != nil {
		return cfg, errors.Wrapf(err, "unmarshaling config content %s", content)
	}
	return cfg, nil
}

func (j *Job) init() error {
	tx := j.db.Begin()

	var maxNum struct{ Num int }
	err := tx.Table("jobs").Select("MAX(num) as num").Where("owner = ? AND repo = ?", j.Owner, j.Repo).First(&maxNum).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	j.Num = maxNum.Num + 1
	j.Status = "Started"
	j.log = logrus.WithFields(logrus.Fields{
		"sha":  j.HeadSHA[:8],
		"job#": j.Num,
		"repo": j.Owner + "/" + j.Repo,
	})
	err = tx.Create(j).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (j *Job) setNextNum() error {
	return nil
}

func (j *Job) githubURL() string {
	return "github.com/" + j.Owner + "/" + j.Repo
}

const credits = "\nCreated by [goreadme](" + githubAppURL + ")\n"
