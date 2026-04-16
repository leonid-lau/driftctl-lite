package tfstate

import (
	"fmt"
	"sync"
)

// HistoryStore manages per-key History instances.
type HistoryStore struct {
	mu      sync.RWMutex
	records map[string]*History
	maxSize int
}

// NewHistoryStore creates a HistoryStore where each key retains at most
// maxSize snapshots.
func NewHistoryStore(maxSize int) *HistoryStore {
	return &HistoryStore{
		records: make(map[string]*History),
		maxSize: maxSize,
	}
}

// Record appends snap to the history for the given key.
func (s *HistoryStore) Record(key string, snap *Snapshot) {
	s.mu.Lock()
	defer s.mu.Unlock()
	h, ok := s.records[key]
	if !ok {
		h = NewHistory(s.maxSize)
		s.records[key] = h
	}
	h.Add(snap)
}

// Get returns the History for key, or an error if not found.
func (s *HistoryStore) Get(key string) (*History, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	h, ok := s.records[key]
	if !ok {
		return nil, fmt.Errorf("no history for key %q", key)
	}
	return h, nil
}

// Keys returns all keys currently tracked.
func (s *HistoryStore) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.records))
	for k := range s.records {
		keys = append(keys, k)
	}
	return keys
}

// Delete removes the history for key.
func (s *HistoryStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.records, key)
}

// HasChanged reports whether the latest snapshot for key has a different
// checksum than prev. Returns false if the key has no history.
func (s *HistoryStore) HasChanged(key, prev string) bool {
	h, err := s.Get(key)
	if err != nil {
		return false
	}
	return h.Changed(prev)
}
