package tfstate

import (
	"fmt"
	"sort"
	"time"
)

// HistoryEntry records a snapshot taken at a specific point in time.
type HistoryEntry struct {
	Timestamp time.Time
	Checksum  string
	Snapshot  *Snapshot
}

// History maintains an ordered list of snapshot entries for a given state key.
type History struct {
	entries []HistoryEntry
	maxSize int
}

// NewHistory creates a History that retains at most maxSize entries.
// If maxSize <= 0, it defaults to 10.
func NewHistory(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 10
	}
	return &History{maxSize: maxSize}
}

// Add appends a new entry to the history, evicting the oldest if necessary.
func (h *History) Add(snap *Snapshot) {
	entry := HistoryEntry{
		Timestamp: time.Now().UTC(),
		Checksum:  snap.Checksum,
		Snapshot:  snap,
	}
	h.entries = append(h.entries, entry)
	if len(h.entries) > h.maxSize {
		h.entries = h.entries[len(h.entries)-h.maxSize:]
	}
}

// Entries returns all history entries sorted oldest-first.
func (h *History) Entries() []HistoryEntry {
	sorted := make([]HistoryEntry, len(h.entries))
	copy(sorted, h.entries)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Timestamp.Before(sorted[j].Timestamp)
	})
	return sorted
}

// Latest returns the most recently added entry, or an error if empty.
func (h *History) Latest() (HistoryEntry, error) {
	if len(h.entries) == 0 {
		return HistoryEntry{}, fmt.Errorf("history is empty")
	}
	return h.entries[len(h.entries)-1], nil
}

// Len returns the number of entries currently stored.
func (h *History) Len() int {
	return len(h.entries)
}

// Changed reports whether the checksum of the latest entry differs from prev.
// Returns false if history is empty.
func (h *History) Changed(prev string) bool {
	latest, err := h.Latest()
	if err != nil {
		return false
	}
	return latest.Checksum != prev
}
