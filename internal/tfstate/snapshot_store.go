package tfstate

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// ErrSnapshotNotFound is returned when no snapshot file exists at the given path.
var ErrSnapshotNotFound = errors.New("snapshot not found")

// SnapshotStore handles persistence of Snapshots to the filesystem.
type SnapshotStore struct {
	Dir string
}

// NewSnapshotStore creates a SnapshotStore rooted at dir, creating it if necessary.
func NewSnapshotStore(dir string) (*SnapshotStore, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("snapshot store: mkdir %s: %w", dir, err)
	}
	return &SnapshotStore{Dir: dir}, nil
}

// Save writes a snapshot to <dir>/<name>.json, overwriting any existing file.
func (ss *SnapshotStore) Save(name string, s *Snapshot) error {
	path := ss.snapshotPath(name)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("snapshot store: create %s: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("snapshot store: encode: %w", err)
	}
	return nil
}

// Load reads a snapshot from <dir>/<name>.json.
func (ss *SnapshotStore) Load(name string) (*Snapshot, error) {
	path := ss.snapshotPath(name)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrSnapshotNotFound
		}
		return nil, fmt.Errorf("snapshot store: open %s: %w", path, err)
	}
	defer f.Close()

	var s Snapshot
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("snapshot store: decode: %w", err)
	}
	return &s, nil
}

// Delete removes the snapshot file for the given name, if it exists.
func (ss *SnapshotStore) Delete(name string) error {
	path := ss.snapshotPath(name)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("snapshot store: delete %s: %w", path, err)
	}
	return nil
}

func (ss *SnapshotStore) snapshotPath(name string) string {
	return filepath.Join(ss.Dir, name+".json")
}
