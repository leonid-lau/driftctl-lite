package tfstate

import (
	"sort"
	"testing"
)

func idxGroupResource(group string) Resource {
	return Resource{
		Type:       "aws_iam_group",
		ID:         group,
		Attributes: map[string]interface{}{"group": group},
	}
}

func TestBuildGroupIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxGroupResource("admins"),
		idxGroupResource("developers"),
	}
	idx := BuildGroupIndex(resources)
	got := idx.Lookup("admins")
	if len(got) != 1 || got[0].ID != "admins" {
		t.Fatalf("expected 1 result for 'admins', got %v", got)
	}
}

func TestBuildGroupIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxGroupResource("Admins")}
	idx := BuildGroupIndex(resources)
	got := idx.Lookup("admins")
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestBuildGroupIndex_LookupMissing(t *testing.T) {
	idx := BuildGroupIndex([]Resource{idxGroupResource("ops")})
	got := idx.Lookup("unknown")
	if got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestBuildGroupIndex_Groups(t *testing.T) {
	resources := []Resource{
		idxGroupResource("admins"),
		idxGroupResource("developers"),
		idxGroupResource("ops"),
	}
	idx := BuildGroupIndex(resources)
	groups := idx.Groups()
	sort.Strings(groups)
	expected := []string{"admins", "developers", "ops"}
	for i, g := range expected {
		if groups[i] != g {
			t.Fatalf("expected group %q at index %d, got %q", g, i, groups[i])
		}
	}
}

func TestBuildGroupIndex_EmptyInput(t *testing.T) {
	idx := BuildGroupIndex(nil)
	if len(idx.Groups()) != 0 {
		t.Fatal("expected empty index")
	}
}
