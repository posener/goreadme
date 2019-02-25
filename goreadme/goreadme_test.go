package goreadme

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/gddo/gosrc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

var gr = New(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
	&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
)))

func init() {
	path, err := filepath.Abs("./goreadme/testdata")
	if err != nil {
		panic(err)
	}
	gosrc.SetLocalDevMode(path)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	for _, dir := range testDirs(t) {
		t.Run(dir, func(t *testing.T) {
			dir := dir
			t.Parallel()
			buf := bytes.NewBuffer(nil)
			err := gr.Create(context.Background(), "./" + dir, buf)
			require.NoError(t, err)
			assertReadme(t, dir, buf.String())
		})
	}
}

func assertReadme(t *testing.T, dir string, got string) {
	t.Helper()
	want, err := ioutil.ReadFile(dir + "/README.md")
	require.NoError(t, err)
	assert.Equal(t, string(want), got)
}

func testDirs(t *testing.T) []string {
	t.Helper()
	files, err := ioutil.ReadDir("testdata")
	require.NoError(t, err)

	dirs := make([]string, 0, len(files))

	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, "testdata/"+f.Name())
		}
	}
	return dirs
}
