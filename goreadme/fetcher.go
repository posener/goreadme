package goreadme

import (
	"context"
	"net/http"
	"sort"
	"sync"

	"github.com/golang/gddo/doc"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// subpackagesFetcher fetches sub packages recursively.
type subpackagesFetcher struct {
	client     *http.Client
	importPath string
	recursive  bool

	wg       sync.WaitGroup
	mu       sync.Mutex
	errors   *multierror.Error
	packages []subPkg
}

func (f *subpackagesFetcher) Fetch(ctx context.Context, pkg *doc.Package) ([]subPkg, error) {
	for _, subDir := range pkg.Subdirectories {
		f.fetch(ctx, subDir)
	}
	f.wg.Wait()
	sort.Slice(f.packages, func(i, j int) bool { return f.packages[i].Path < f.packages[j].Path })
	return f.packages, f.errors.ErrorOrNil()
}

// Concurrently fetches information for all sub directories.
func (f *subpackagesFetcher) fetch(ctx context.Context, subDir string) {
	f.wg.Add(1)
	importPath := f.importPath + "/" + subDir

	go func() {
		defer f.wg.Done()
		sp, err := docGet(ctx, f.client, importPath, "")
		f.mu.Lock()
		defer f.mu.Unlock()
		if err != nil {
			f.errors = multierror.Append(f.errors, errors.Wrapf(err, "failed getting %s", importPath))
			return
		}
		// Append to packages only if this directory is a go package.
		if sp.Name != "" {
			f.packages = append(f.packages, subPkg{Path: subDir, Package: sp})
		}
		if f.recursive {
			for _, sd := range sp.Subdirectories {
				f.fetch(ctx, subDir+"/"+sd)
			}
		}
	}()
}
