package tfstate

import "fmt"

// MergeOptions controls how multiple states are merged.
type MergeOptions struct {
	// OnConflict determines behaviour when two states contain the same resource ID.
	// "error" (default) returns an error; "last-wins" keeps the later entry.
	OnConflict string
}

// DefaultMergeOptions returns sensible defaults for MergeOptions.
func DefaultMergeOptions() MergeOptions {
	return MergeOptions{OnConflict: "error"}
}

// MergeStates combines multiple parsed States into a single flat resource list.
// Duplicate resource IDs are handled according to opts.OnConflict.
func MergeStates(states []State, opts MergeOptions) ([]Resource, error) {
	if opts.OnConflict == "" {
		opts.OnConflict = "error"
	}

	seen := make(map[string]struct{})
	var merged []Resource

	for _, s := range states {
		for _, r := range s.Resources {
			key := fallbackKey(r)
			if _, exists := seen[key]; exists {
				switch opts.OnConflict {
				case "error":
					return nil, fmt.Errorf("merge conflict: duplicate resource key %q", key)
				case "last-wins":
					// replace the previous entry
					for i, existing := range merged {
						if fallbackKey(existing) == key {
							merged[i] = r
							break
						}
					}
					continue
				default:
					return nil, fmt.Errorf("unknown OnConflict strategy: %q", opts.OnConflict)
				}
			}
			seen[key] = struct{}{}
			merged = append(merged, r)
		}
	}
	return merged, nil
}
