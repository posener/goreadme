package goreadme

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/gddo/gosrc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

// writeReadmes is used to write the expected output instead of asserting equality.
var writeReadmes = os.Getenv("WRITE_READMES") == "1"

var gr = New(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
	&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
)))

func init() {
	path, err := filepath.Abs("./testdata")
	if err != nil {
		panic(err)
	}
	gosrc.SetLocalDevMode(path)
}

func TestCreate(t *testing.T) {
	t.Parallel()

	for _, dir := range testDirs(t) {
		t.Run(dir, func(t *testing.T) {
			dir := "./" + dir
			t.Parallel()

			buf := bytes.NewBuffer(nil)
			cfg := loadConfig(t, dir)
			t.Logf("Running with config: %+v", cfg)
			err := gr.WithConfig(cfg).Create(context.Background(), dir, buf)
			require.NoError(t, err)
			if writeReadmes {
				// Helper with writing the README files.
				require.NoError(t, ioutil.WriteFile(readmeFileName(dir), buf.Bytes(), 0664))
			}
			assertReadme(t, dir, buf.String())
		})
	}
}

func assertReadme(t *testing.T, dir string, got string) {
	t.Helper()

	want, err := ioutil.ReadFile(readmeFileName(dir))
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

func loadConfig(t *testing.T, dir string) Config {
	t.Helper()
	c := Config{}
	b, err := ioutil.ReadFile(dir + "/goreadme.json")
	if err != nil {
		return c
	}
	err = json.Unmarshal(b, &c)
	require.NoError(t, err)
	return c
}

func readmeFileName(dir string) string {
	return dir + "/README.md"
}
