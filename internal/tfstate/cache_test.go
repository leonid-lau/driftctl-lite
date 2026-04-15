package tfstate

import (
	"testing"
	"time"
)

func makeState(resourceType string) *State {
	return &State{
		Version: 4,
		Resources: []Resource{
			{Type: resourceType, Name: "example", Attributes: map[string]interface{}{"id": "abc"}},
		},
	}
}

func TestStateCache_SetAndGet(t *testing.T) {
	c := NewStateCache(5 * time.Second)
	s := makeState("aws_s3_bucket")
	c.Set("path/to/state.tfstate", s)

	got, ok := c.Get("path/to/state.tfstate")
	if !ok {
		t.Fatal("expected cache hit, got miss")
	}
	if got != s {
		t.Errorf("expected same state pointer")
	}
}

func TestStateCache_Miss(t *testing.T) {
	c := NewStateCache(5 * time.Second)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected cache miss")
	}
}

func TestStateCache_Expiry(t *testing.T) {
	c := NewStateCache(10 * time.Millisecond)
	c.Set("key", makeState("aws_instance"))

	time.Sleep(20 * time.Millisecond)

	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected expired entry to be a miss")
	}
}

func TestStateCache_Delete(t *testing.T) {
	c := NewStateCache(5 * time.Second)
	c.Set("key", makeState("aws_vpc"))
	c.Delete("key")

	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected miss after delete")
	}
}

func TestStateCache_Flush(t *testing.T) {
	c := NewStateCache(10 * time.Millisecond)
	c.Set("a", makeState("aws_instance"))
	c.Set("b", makeState("aws_s3_bucket"))

	time.Sleep(20 * time.Millisecond)
	c.Set("c", makeState("aws_vpc")) // fresh entry

	removed := c.Flush()
	if removed != 2 {
		t.Errorf("expected 2 removed, got %d", removed)
	}
	if c.Len() != 1 {
		t.Errorf("expected 1 remaining, got %d", c.Len())
	}
}

func TestStateCache_Len(t *testing.T) {
	c := NewStateCache(5 * time.Second)
	if c.Len() != 0 {
		t.Errorf("expected empty cache")
	}
	c.Set("x", makeState("aws_instance"))
	c.Set("y", makeState("aws_vpc"))
	if c.Len() != 2 {
		t.Errorf("expected 2, got %d", c.Len())
	}
}
