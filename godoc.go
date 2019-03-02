package goreadme

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/golang/gddo/doc"
	"github.com/pkg/errors"
)

// docGet is a wrapper around doc.Get function, that workarounds golang/gddo#600.
func docGet(ctx context.Context, client *http.Client, name, tag string) (*doc.Package, error) {
	p, err := doc.Get(ctx, client, name, tag)
	if err != nil {
		return nil, err
	}
	err = workaroundLocalSubdirs(p, name)
	return p, err
}

// workaroundLocalSubdirs adds subdireoctires for local load.
// Workaround for golang/gddo#600
func workaroundLocalSubdirs(p *doc.Package, pkg string) error {
	if !strings.HasPrefix(pkg, ".") {
		return nil // Not local
	}

	files, err := ioutil.ReadDir(p.ImportPath)
	if err != nil {
		errors.Wrap(err, "Failed reading import path")
	}

	for _, f := range files {
		if f.IsDir() {
			p.Subdirectories = append(p.Subdirectories, f.Name())
		}
	}
	return nil
}
