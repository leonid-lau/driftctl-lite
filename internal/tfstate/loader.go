package tfstate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// LoadOptions configures how state files are discovered and loaded.
type LoadOptions struct {
	// SearchDirs are directories to search for .tfstate files.
	SearchDirs []string
	// Recursive enables recursive directory traversal.
	Recursive bool
}

// DefaultLoadOptions returns sensible defaults for loading state files.
func DefaultLoadOptions() LoadOptions {
	return LoadOptions{
		SearchDirs: []string{"."},
		Recursive:  false,
	}
}

// LoadAll discovers and parses all Terraform state files under the given options.
// It returns a merged slice of resources from all discovered state files.
func LoadAll(opts LoadOptions) ([]Resource, error) {
	paths, err := discoverStateFiles(opts)
	if err != nil {
		return nil, fmt.Errorf("discovering state files: %w", err)
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no .tfstate files found in search directories: %v", opts.SearchDirs)
	}

	var all []Resource
	for _, p := range paths {
		resources, err := ParseFile(p)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", p, err)
		}
		all = append(all, resources...)
	}
	return all, nil
}

// discoverStateFiles walks the configured directories and collects paths to
// files ending in ".tfstate".
func discoverStateFiles(opts LoadOptions) ([]string, error) {
	var found []string
	for _, dir := range opts.SearchDirs {
		if opts.Recursive {
			err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if !info.IsDir() && strings.HasSuffix(info.Name(), ".tfstate") {
					found = append(found, path)
				}
				return nil
			})
			if err != nil {
				return nil, fmt.Errorf("walking directory %s: %w", dir, err)
			}
		} else {
			entries, err := os.ReadDir(dir)
			if err != nil {
				return nil, fmt.Errorf("reading directory %s: %w", dir, err)
			}
			for _, e := range entries {
				if !e.IsDir() && strings.HasSuffix(e.Name(), ".tfstate") {
					found = append(found, filepath.Join(dir, e.Name()))
				}
			}
		}
	}
	return found, nil
}
