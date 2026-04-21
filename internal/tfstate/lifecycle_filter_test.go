package tfstate

import (
	"testing"
)

func lifecycleResource(lc string) Resource {
	attrs := map[string]interface{}{}
	if lc != "" {
		attrs["lifecycle"] = lc
	}
	return Resource{Type: "aws_instance", Name: "r_" + lc, Attributes: attrs}
}

func TestFilterByLifecycle_Match(t *testing.T) {
	resources := []Resource{
		lifecycleResource("active"),
		lifecycleResource("deprecated"),
		lifecycleResource("active"),
	}
	got := FilterByLifecycle(resources, "active", DefaultLifecycleFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByLifecycle_EmptyLifecycle_ReturnsAll(t *testing.T) {
	resources := []Resource{lifecycleResource("active"), lifecycleResource("deprecated")}
	got := FilterByLifecycle(resources, "", DefaultLifecycleFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByLifecycle_CaseInsensitive(t *testing.T) {
	resources := []Resource{lifecycleResource("Active"), lifecycleResource("ACTIVE")}
	got := FilterByLifecycle(resources, "active", DefaultLifecycleFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByLifecycle_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{lifecycleResource("Active")}
	opts := LifecycleFilterOptions{CaseSensitive: true}
	got := FilterByLifecycle(resources, "active", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByLifecycles_ORSemantics(t *testing.T) {
	resources := []Resource{
		lifecycleResource("active"),
		lifecycleResource("deprecated"),
		lifecycleResource("retired"),
	}
	got := FilterByLifecycles(resources, []string{"active", "deprecated"}, DefaultLifecycleFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildLifecycleIndex_Lookup(t *testing.T) {
	resources := []Resource{lifecycleResource("active"), lifecycleResource("deprecated")}
	idx := BuildLifecycleIndex(resources)
	got := idx.Lookup("active")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildLifecycleIndex_LookupMissing(t *testing.T) {
	idx := BuildLifecycleIndex([]Resource{})
	if idx.Lookup("active") != nil {
		t.Fatal("expected nil for missing lifecycle")
	}
}

func TestBuildLifecycleIndex_Lifecycles(t *testing.T) {
	resources := []Resource{lifecycleResource("active"), lifecycleResource("deprecated")}
	idx := BuildLifecycleIndex(resources)
	if len(idx.Lifecycles()) != 2 {
		t.Fatalf("expected 2 lifecycles, got %d", len(idx.Lifecycles()))
	}
}
