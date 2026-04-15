package tfstate

import (
	"fmt"
	"os"
	"path/filepath"
)

// LoadOptions configures state file discovery.
type LoadOptions struct {
	// Recursive enables walking subdirectories.
	Recursive bool
	// Validate runs ValidateState on each parsed state when true.
	Validate bool
}

// DefaultLoadOptions returns sensible defaults.
func DefaultLoadOptions() LoadOptions {
	return LoadOptions{Recursive: false, Validate: true}
}

// LoadAll discovers and parses all *.tfstate files under root.
// It returns the merged list of resources and any fatal error.
func LoadAll(root string, opts LoadOptions) ([]Resource, error) {
	files, err := discoverStateFiles(root, opts.Recursive)
	if err != nil {
		return nil, fmt.Errorf("discover state files: %w", err)
	}
	if len(files) == 0 {
		return nil, nil
	}

	var states []*State
	for _, f := range files {
		s, err := ParseFile(f)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", f, err)
		}
		if opts.Validate {
			if verr := ValidateState(s); verr != nil {
				return nil, fmt.Errorf("validate %s: %w", f, verr)
			}
		}
		states = append(states, s)
	}

	merged, err := MergeStates(states, DefaultMergeOptions())
	if err != nil {
		return nil, fmt.Errorf("merge states: %w", err)
	}
	return merged.Resources, nil
}

// discoverStateFiles returns paths to *.tfstate files under root.
func discoverStateFiles(root string, recursive bool) ([]string, error) {
	var found []string

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && path != root && !recursive {
			return filepath.SkipDir
		}
		if !info.IsDir() && filepath.Ext(path) == ".tfstate" {
			found = append(found, path)
		}
		return nil
	}

	if err := filepath.Walk(root, walkFn); err != nil {
		return nil, err
	}
	return found, nil
}
