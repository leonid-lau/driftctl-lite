package tfstate

import (
	"fmt"
	"time"
)

// DefaultCacheTTL is the default time-to-live for cached state files.
const DefaultCacheTTL = 30 * time.Second

// CachedLoader wraps LoadAll with an in-memory StateCache so that repeated
// calls for the same set of paths avoid redundant disk I/O and JSON parsing.
type CachedLoader struct {
	cache *StateCache
	opts  LoadOptions
}

// NewCachedLoader creates a CachedLoader with the given TTL and load options.
func NewCachedLoader(ttl time.Duration, opts LoadOptions) *CachedLoader {
	return &CachedLoader{
		cache: NewStateCache(ttl),
		opts:  opts,
	}
}

// Load returns the merged State for the given root directory, using the cache
// when a valid (non-expired) entry exists.
func (cl *CachedLoader) Load(root string) (*State, error) {
	if s, ok := cl.cache.Get(root); ok {
		return s, nil
	}

	states, err := LoadAll(root, cl.opts)
	if err != nil {
		return nil, fmt.Errorf("cached loader: %w", err)
	}

	merged, err := MergeStates(states, DefaultMergeOptions())
	if err != nil {
		return nil, fmt.Errorf("cached loader merge: %w", err)
	}

	cl.cache.Set(root, merged)
	return merged, nil
}

// Invalidate removes the cached entry for the given root, forcing a fresh load
// on the next call to Load.
func (cl *CachedLoader) Invalidate(root string) {
	cl.cache.Delete(root)
}

// FlushExpired removes all expired entries from the underlying cache and
// returns the number of entries removed.
func (cl *CachedLoader) FlushExpired() int {
	return cl.cache.Flush()
}
