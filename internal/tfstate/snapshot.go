package tfstate

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// Snapshot represents a point-in-time capture of parsed Terraform state resources.
type Snapshot struct {
	CreatedAt time.Time          `json:"created_at"`
	Checksum  string             `json:"checksum"`
	Resources []Resource         `json:"resources"`
	Meta      map[string]string  `json:"meta,omitempty"`
}

// NewSnapshot creates a Snapshot from a slice of resources, computing a deterministic checksum.
func NewSnapshot(resources []Resource, meta map[string]string) (*Snapshot, error) {
	checksum, err := computeChecksum(resources)
	if err != nil {
		return nil, fmt.Errorf("snapshot: compute checksum: %w", err)
	}
	return &Snapshot{
		CreatedAt: time.Now().UTC(),
		Checksum:  checksum,
		Resources: resources,
		Meta:      meta,
	}, nil
}

// Equal reports whether two snapshots contain identical resource sets by comparing checksums.
func (s *Snapshot) Equal(other *Snapshot) bool {
	if s == nil || other == nil {
		return s == other
	}
	return s.Checksum == other.Checksum
}

// ResourceCount returns the number of resources in the snapshot.
func (s *Snapshot) ResourceCount() int {
	if s == nil {
		return 0
	}
	return len(s.Resources)
}

// computeChecksum produces a SHA-256 hash over a deterministically sorted JSON encoding of resources.
func computeChecksum(resources []Resource) (string, error) {
	sorted := make([]Resource, len(resources))
	copy(sorted, resources)
	sort.Slice(sorted, func(i, j int) bool {
		ki := fallbackKey(sorted[i])
		kj := fallbackKey(sorted[j])
		return ki < kj
	})

	data, err := json.Marshal(sorted)
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
