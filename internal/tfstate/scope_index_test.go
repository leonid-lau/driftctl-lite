package tfstate

import (
	"testing"
)

func idxScopeResource(scope string) Resource {
	return Resource{
		Type:       "aws_iam_policy",
		Attributes: map[string]string{"scope": scope},
	}
}

func TestBuildScopeIndex_Lookup(t *testing.T) {
	resources := []Resource{idxScopeResource("global"), idxScopeResource("regional")}
	idx := BuildScopeIndex(resources)
	got := idx.Lookup("global")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildScopeIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxScopeResource("Global")}
	idx := BuildScopeIndex(resources)
	if got := idx.Lookup("global"); len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildScopeIndex_LookupMissing(t *testing.T) {
	idx := BuildScopeIndex([]Resource{idxScopeResource("global")})
	if got := idx.Lookup("local"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestBuildScopeIndex_Scopes(t *testing.T) {
	resources := []Resource{idxScopeResource("global"), idxScopeResource("regional")}
	idx := BuildScopeIndex(resources)
	if len(idx.Scopes()) != 2 {
		t.Fatalf("expected 2 scopes, got %d", len(idx.Scopes()))
	}
}

func TestBuildScopeIndex_EmptyInput(t *testing.T) {
	idx := BuildScopeIndex(nil)
	if len(idx.Scopes()) != 0 {
		t.Fatal("expected empty")
	}
}
