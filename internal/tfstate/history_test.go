package tfstate

import (
	"testing"
	"time"
)

func buildSnapshot(id, val string) *Snapshot {
	r := Resource{
		Type: "aws_instance",
		Name: id,
		Attributes: map[string]interface{}{"key": val},
	}
	snap, _ := NewSnapshot([]Resource{r})
	return snap
}

func TestHistory_AddAndLen(t *testing.T) {
	h := NewHistory(5)
	if h.Len() != 0 {
		t.Fatalf("expected 0, got %d", h.Len())
	}
	h.Add(buildSnapshot("a", "v1"))
	h.Add(buildSnapshot("b", "v2"))
	if h.Len() != 2 {
		t.Fatalf("expected 2, got %d", h.Len())
	}
}

func TestHistory_MaxSizeEviction(t *testing.T) {
	h := NewHistory(3)
	for i := 0; i < 5; i++ {
		h.Add(buildSnapshot("r", string(rune('a'+i))))
	}
	if h.Len() != 3 {
		t.Fatalf("expected 3 after eviction, got %d", h.Len())
	}
}

func TestHistory_Latest_Empty(t *testing.T) {
	h := NewHistory(5)
	_, err := h.Latest()
	if err == nil {
		t.Fatal("expected error for empty history")
	}
}

func TestHistory_Latest_ReturnsNewest(t *testing.T) {
	h := NewHistory(5)
	snap1 := buildSnapshot("a", "v1")
	time.Sleep(time.Millisecond)
	snap2 := buildSnapshot("b", "v2")
	h.Add(snap1)
	h.Add(snap2)
	latest, err := h.Latest()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if latest.Checksum != snap2.Checksum {
		t.Errorf("expected checksum %s, got %s", snap2.Checksum, latest.Checksum)
	}
}

func TestHistory_Entries_Sorted(t *testing.T) {
	h := NewHistory(5)
	h.Add(buildSnapshot("a", "v1"))
	time.Sleep(time.Millisecond)
	h.Add(buildSnapshot("b", "v2"))
	entries := h.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if !entries[0].Timestamp.Before(entries[1].Timestamp) {
		t.Error("entries not sorted oldest-first")
	}
}

func TestHistory_Changed(t *testing.T) {
	h := NewHistory(5)
	snap := buildSnapshot("a", "v1")
	h.Add(snap)
	if h.Changed(snap.Checksum) {
		t.Error("expected no change for same checksum")
	}
	if !h.Changed("different-checksum") {
		t.Error("expected change for different checksum")
	}
}

func TestHistory_DefaultMaxSize(t *testing.T) {
	h := NewHistory(0)
	for i := 0; i < 15; i++ {
		h.Add(buildSnapshot("r", string(rune('a'+i%26))))
	}
	if h.Len() != 10 {
		t.Fatalf("expected default max 10, got %d", h.Len())
	}
}
