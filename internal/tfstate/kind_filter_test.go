package tfstate

import (
	"testing"
)

func kindResource(kind, id string) Resource {
	return Resource{Type: "aws_instance", Kind: kind, ID: id, Attributes: map[string]interface{}{}}
}

func TestFilterByKind_Match(t *testing.T) {
	resources := []Resource{
		kindResource("Deployment", "r1"),
		kindResource("Service", "r2"),
		kindResource("deployment", "r3"),
	}
	got := FilterByKind(resources, "deployment", DefaultKindFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByKind_EmptyKind_ReturnsAll(t *testing.T) {
	resources := []Resource{kindResource("Deployment", "r1"), kindResource("Service", "r2")}
	got := FilterByKind(resources, "", DefaultKindFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByKind_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{kindResource("Deployment", "r1")}
	opts := KindFilterOptions{CaseSensitive: true}
	got := FilterByKind(resources, "deployment", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByKinds_ORSemantics(t *testing.T) {
	resources := []Resource{
		kindResource("Deployment", "r1"),
		kindResource("Service", "r2"),
		kindResource("ConfigMap", "r3"),
	}
	got := FilterByKinds(resources, []string{"deployment", "service"}, DefaultKindFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildKindIndex_Lookup(t *testing.T) {
	resources := []Resource{kindResource("Deployment", "r1"), kindResource("Service", "r2")}
	idx := BuildKindIndex(resources)
	got := idx.Lookup("deployment")
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("unexpected result: %v", got)
	}
}

func TestBuildKindIndex_LookupMissing(t *testing.T) {
	idx := BuildKindIndex([]Resource{})
	if idx.Lookup("anything") != nil {
		t.Fatal("expected nil")
	}
}

func TestBuildKindIndex_Kinds(t *testing.T) {
	resources := []Resource{kindResource("Deployment", "r1"), kindResource("Service", "r2")}
	idx := BuildKindIndex(resources)
	kinds := idx.Kinds()
	if len(kinds) != 2 {
		t.Fatalf("expected 2 kinds, got %d", len(kinds))
	}
}
