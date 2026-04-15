package tfstate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeCachedState(t *testing.T, dir string, name string, s *State) string {
	t.Helper()
	path := filepath.Join(dir, name)
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	return path
}

func TestCachedLoader_Load_ReturnsMergedState(t *testing.T) {
	dir := t.TempDir()
	writeCachedState(t, dir, "terraform.tfstate", &State{
		Version: 4,
		Resources: []Resource{
			{Type: "aws_instance", Name: "web", Attributes: map[string]interface{}{"id": "i-123"}},
		},
	})

	cl := NewCachedLoader(5*time.Second, DefaultLoadOptions())
	s, err := cl.Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Resources) != 1 {
		t.Errorf("expected 1 resource, got %d", len(s.Resources))
	}
}

func TestCachedLoader_Load_CacheHit(t *testing.T) {
	dir := t.TempDir()
	writeCachedState(t, dir, "terraform.tfstate", &State{
		Version:   4,
		Resources: []Resource{{Type: "aws_vpc", Name: "main", Attributes: map[string]interface{}{"id": "vpc-1"}}},
	})

	cl := NewCachedLoader(5*time.Second, DefaultLoadOptions())
	first, _ := cl.Load(dir)
	second, _ := cl.Load(dir)

	if first != second {
		t.Error("expected second call to return cached pointer")
	}
}

func TestCachedLoader_Invalidate(t *testing.T) {
	dir := t.TempDir()
	writeCachedState(t, dir, "terraform.tfstate", &State{
		Version:   4,
		Resources: []Resource{{Type: "aws_s3_bucket", Name: "assets", Attributes: map[string]interface{}{"id": "bucket-1"}}},
	})

	cl := NewCachedLoader(5*time.Second, DefaultLoadOptions())
	first, _ := cl.Load(dir)
	cl.Invalidate(dir)
	second, _ := cl.Load(dir)

	if first == second {
		t.Error("expected fresh load after invalidation")
	}
}

func TestCachedLoader_FlushExpired(t *testing.T) {
	dir := t.TempDir()
	writeCachedState(t, dir, "terraform.tfstate", &State{
		Version:   4,
		Resources: []Resource{{Type: "aws_instance", Name: "db", Attributes: map[string]interface{}{"id": "i-999"}}},
	})

	cl := NewCachedLoader(10*time.Millisecond, DefaultLoadOptions())
	cl.Load(dir) //nolint:errcheck

	time.Sleep(20 * time.Millisecond)
	removed := cl.FlushExpired()
	if removed != 1 {
		t.Errorf("expected 1 flushed, got %d", removed)
	}
}
