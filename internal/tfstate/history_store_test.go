package tfstate

import (
	"testing"
)

func TestHistoryStore_RecordAndGet(t *testing.T) {
	store := NewHistoryStore(5)
	snap := buildSnapshot("res1", "v1")
	store.Record("key1", snap)

	h, err := store.Get("key1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if h.Len() != 1 {
		t.Errorf("expected 1 entry, got %d", h.Len())
	}
}

func TestHistoryStore_Get_NotFound(t *testing.T) {
	store := NewHistoryStore(5)
	_, err := store.Get("missing")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestHistoryStore_Keys(t *testing.T) {
	store := NewHistoryStore(5)
	store.Record("a", buildSnapshot("r", "v1"))
	store.Record("b", buildSnapshot("r", "v2"))
	keys := store.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
}

func TestHistoryStore_Delete(t *testing.T) {
	store := NewHistoryStore(5)
	store.Record("key1", buildSnapshot("r", "v1"))
	store.Delete("key1")
	_, err := store.Get("key1")
	if err == nil {
		t.Fatal("expected error after delete")
	}
}

func TestHistoryStore_MultipleRecords(t *testing.T) {
	store := NewHistoryStore(3)
	for i := 0; i < 5; i++ {
		store.Record("key1", buildSnapshot("r", string(rune('a'+i))))
	}
	h, _ := store.Get("key1")
	if h.Len() != 3 {
		t.Errorf("expected 3 (maxSize), got %d", h.Len())
	}
}

func TestHistoryStore_HasChanged_True(t *testing.T) {
	store := NewHistoryStore(5)
	snap := buildSnapshot("r", "v1")
	store.Record("key1", snap)
	if !store.HasChanged("key1", "old-checksum") {
		t.Error("expected HasChanged to return true")
	}
}

func TestHistoryStore_HasChanged_False(t *testing.T) {
	store := NewHistoryStore(5)
	snap := buildSnapshot("r", "v1")
	store.Record("key1", snap)
	if store.HasChanged("key1", snap.Checksum) {
		t.Error("expected HasChanged to return false for same checksum")
	}
}

func TestHistoryStore_HasChanged_MissingKey(t *testing.T) {
	store := NewHistoryStore(5)
	if store.HasChanged("nonexistent", "any") {
		t.Error("expected false for missing key")
	}
}
