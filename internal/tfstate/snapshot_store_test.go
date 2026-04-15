package tfstate

import (
	"errors"
	"os"
	"testing"
)

func TestSnapshotStore_SaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	store, err := NewSnapshotStore(dir)
	if err != nil {
		t.Fatalf("NewSnapshotStore: %v", err)
	}

	res := []Resource{makeResource("aws_vpc", "main", "vpc-123")}
	snap, _ := NewSnapshot(res, map[string]string{"env": "test"})

	if err := store.Save("baseline", snap); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := store.Load("baseline")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Checksum != snap.Checksum {
		t.Errorf("checksum mismatch: got %s, want %s", loaded.Checksum, snap.Checksum)
	}
	if loaded.Meta["env"] != "test" {
		t.Errorf("meta not preserved: got %v", loaded.Meta)
	}
}

func TestSnapshotStore_Load_NotFound(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewSnapshotStore(dir)

	_, err := store.Load("nonexistent")
	if !errors.Is(err, ErrSnapshotNotFound) {
		t.Errorf("expected ErrSnapshotNotFound, got %v", err)
	}
}

func TestSnapshotStore_Delete(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewSnapshotStore(dir)
	snap, _ := NewSnapshot(nil, nil)

	_ = store.Save("tmp", snap)
	if err := store.Delete("tmp"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := store.Load("tmp")
	if !errors.Is(err, ErrSnapshotNotFound) {
		t.Errorf("expected ErrSnapshotNotFound after delete, got %v", err)
	}
}

func TestSnapshotStore_Delete_NonExistent(t *testing.T) {
	dir := t.TempDir()
	store, _ := NewSnapshotStore(dir)
	if err := store.Delete("ghost"); err != nil {
		t.Errorf("Delete of non-existent should not error, got %v", err)
	}
}

func TestNewSnapshotStore_CreatesDir(t *testing.T) {
	base := t.TempDir()
	newDir := base + "/nested/store"
	_, err := NewSnapshotStore(newDir)
	if err != nil {
		t.Fatalf("expected dir creation, got: %v", err)
	}
	if _, err := os.Stat(newDir); os.IsNotExist(err) {
		t.Error("directory was not created")
	}
}
