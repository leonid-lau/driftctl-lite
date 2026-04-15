package tfstate

import (
	"errors"
	"fmt"
)

// ValidationError holds all issues found during state validation.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	if len(e.Issues) == 1 {
		return fmt.Sprintf("state validation failed: %s", e.Issues[0])
	}
	return fmt.Sprintf("state validation failed with %d issues: %v", len(e.Issues), e.Issues)
}

// ValidateState checks a parsed State for common integrity problems.
// It returns a *ValidationError if any issues are found, or nil on success.
func ValidateState(s *State) error {
	if s == nil {
		return errors.New("state is nil")
	}

	var issues []string

	if s.Version < 1 {
		issues = append(issues, fmt.Sprintf("unsupported state version: %d", s.Version))
	}

	seen := make(map[string]int) // key -> first index
	for i, r := range s.Resources {
		if r.Type == "" {
			issues = append(issues, fmt.Sprintf("resource[%d] has empty type", i))
		}
		if r.Name == "" {
			issues = append(issues, fmt.Sprintf("resource[%d] has empty name", i))
		}
		key := fallbackKey(r)
		if prev, ok := seen[key]; ok {
			issues = append(issues, fmt.Sprintf("duplicate resource key %q at index %d (first seen at %d)", key, i, prev))
		} else {
			seen[key] = i
		}
	}

	if len(issues) > 0 {
		return &ValidationError{Issues: issues}
	}
	return nil
}
