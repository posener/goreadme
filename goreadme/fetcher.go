package goreadme

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/golang/gddo/doc"
	"github.com/hashicorp/go-multierror"
)

type fetcher struct {
	client     *http.Client
	importPath string
	recursive  bool

	wg       sync.WaitGroup
	mu       sync.Mutex
	errors   *multierror.Error
	packages []subPkg
}

func (f *fetcher) Fetch(ctx context.Context, pkg *doc.Package) ([]subPkg, error) {
	for _, subDir := range pkg.Subdirectories {
		f.fetch(ctx, subDir)
	}
	f.wg.Wait()
	return f.packages, f.errors.ErrorOrNil()
}

// Concurrently fetches information for all sub directories.
func (f *fetcher) fetch(ctx context.Context, subDir string) {
	f.wg.Add(1)
	importPath := f.importPath + "/" + subDir

	go func() {
		defer f.wg.Done()
		log.Printf("Getting %s", importPath)
		sp, err := doc.Get(ctx, f.client, importPath, "")
		f.mu.Lock()
		defer f.mu.Unlock()
		if err != nil {
			f.errors = multierror.Append(f.errors, fmt.Errorf("failed getting %s: %s", importPath, err))
			return
		}
		// Append to packages only if this directory is a go package.
		if sp.Name != "" {
			f.packages = append(f.packages, subPkg{Path: subDir, Package: sp})
		}
		if f.recursive {
			for _, subSubDir := range sp.Subdirectories {
				f.fetch(ctx, subDir+"/"+subSubDir)
			}
		}
	}()
}
