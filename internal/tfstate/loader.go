package tfstate

import (
	"fmt"
	"os"
	"path/filepath"
)

// LoadOptions configures the behaviour of LoadAll.
type LoadOptions struct {
	Recursive bool
	Merge     MergeOptions
}

// DefaultLoadOptions returns safe defaults.
func DefaultLoadOptions() LoadOptions {
	return LoadOptions{
		Recursive: false,
		Merge:     DefaultMergeOptions(),
	}
}

// LoadAll discovers and parses all *.tfstate files under root, then merges them.
func LoadAll(root string, opts LoadOptions) ([]Resource, error) {
	files, err := discoverStateFiles(root, opts.Recursive)
	if err != nil {
		return nil, fmt.Errorf("discovering state files: %w", err)
	}
	if len(files) == 0 {
		return nil, nil
	}

	states := make([]State, 0, len(files))
	for _, f := range files {
		s, err := ParseFile(f)
		if err != nil {
			return nil, fmt.Errorf("parsing %s: %w", f, err)
		}
		states = append(states, *s)
	}

	return MergeStates(states, opts.Merge)
}

// discoverStateFiles walks root and collects paths ending in .tfstate.
func discoverStateFiles(root string, recursive bool) ([]string, error) {
	var found []string

	if recursive {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && filepath.Ext(path) == ".tfstate" {
				found = append(found, path)
			}
			return nil
		})
		return found, err
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() && filepath.Ext(e.Name()) == ".tfstate" {
			found = append(found, filepath.Join(root, e.Name()))
		}
	}
	return found, nil
}
