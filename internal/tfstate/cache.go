package tfstate

import (
	"sync"
	"time"
)

// CacheEntry holds a parsed state along with its expiry metadata.
type CacheEntry struct {
	State     *State
	CachedAt  time.Time
	ExpiresAt time.Time
}

// IsExpired reports whether the cache entry has passed its TTL.
func (e *CacheEntry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// StateCache is a thread-safe in-memory cache for parsed Terraform states
// keyed by file path.
type StateCache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
	ttl     time.Duration
}

// NewStateCache creates a new StateCache with the given TTL.
func NewStateCache(ttl time.Duration) *StateCache {
	return &StateCache{
		entries: make(map[string]*CacheEntry),
		ttl:     ttl,
	}
}

// Set stores a State in the cache under the given key.
func (c *StateCache) Set(key string, s *State) {
	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	c.entries[key] = &CacheEntry{
		State:     s,
		CachedAt:  now,
		ExpiresAt: now.Add(c.ttl),
	}
}

// Get retrieves a State from the cache. Returns nil, false if not found or expired.
func (c *StateCache) Get(key string) (*State, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, ok := c.entries[key]
	if !ok || entry.IsExpired() {
		return nil, false
	}
	return entry.State, true
}

// Delete removes a single entry from the cache.
func (c *StateCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, key)
}

// Flush removes all expired entries from the cache.
func (c *StateCache) Flush() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	removed := 0
	for k, e := range c.entries {
		if e.IsExpired() {
			delete(c.entries, k)
			removed++
		}
	}
	return removed
}

// Len returns the number of entries currently in the cache (including expired).
func (c *StateCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}
