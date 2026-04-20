package tfstate

import (
	"sort"
	"testing"
)

func idxDepResource(dep string) Resource {
	return Resource{
		Type:       "aws_lambda_function",
		Attributes: map[string]interface{}{"dependency": dep},
	}
}

func TestBuildDependencyIndex_Lookup(t *testing.T) {
	resources := []Resource{idxDepResource("moduleA"), idxDepResource("moduleB"), idxDepResource("moduleA")}
	idx := BuildDependencyIndex(resources)
	got := idx.Lookup("moduleA")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildDependencyIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxDepResource("ModuleX")}
	idx := BuildDependencyIndex(resources)
	got := idx.Lookup("modulex")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildDependencyIndex_LookupMissing(t *testing.T) {
	idx := BuildDependencyIndex([]Resource{idxDepResource("moduleA")})
	if got := idx.Lookup("moduleZ"); len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestBuildDependencyIndex_Dependencies(t *testing.T) {
	resources := []Resource{idxDepResource("moduleA"), idxDepResource("moduleB"), idxDepResource("moduleA")}
	idx := BuildDependencyIndex(resources)
	deps := idx.Dependencies()
	sort.Strings(deps)
	if len(deps) != 2 || deps[0] != "modulea" || deps[1] != "moduleb" {
		t.Fatalf("unexpected deps: %v", deps)
	}
}

func TestBuildDependencyIndex_EmptyInput(t *testing.T) {
	idx := BuildDependencyIndex(nil)
	if got := idx.Lookup("anything"); len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
	if len(idx.Dependencies()) != 0 {
		t.Fatal("expected no dependencies")
	}
}
