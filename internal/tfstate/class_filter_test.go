package tfstate

import (
	"testing"
)

func classResource(class string) Resource {
	attrs := map[string]interface{}{}
	if class != "" {
		attrs["class"] = class
	}
	return Resource{Type: "aws_instance", ID: class + "-id", Attributes: attrs}
}

func TestFilterByClass_Match(t *testing.T) {
	resources := []Resource{classResource("general-purpose"), classResource("compute"), classResource("memory")}
	got := FilterByClass(resources, "compute", DefaultClassFilterOptions())
	if len(got) != 1 || got[0].ID != "compute-id" {
		t.Fatalf("expected 1 compute resource, got %v", got)
	}
}

func TestFilterByClass_EmptyClass_ReturnsAll(t *testing.T) {
	resources := []Resource{classResource("general-purpose"), classResource("compute")}
	got := FilterByClass(resources, "", DefaultClassFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByClass_CaseInsensitive(t *testing.T) {
	resources := []Resource{classResource("General-Purpose"), classResource("compute")}
	got := FilterByClass(resources, "general-purpose", DefaultClassFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(got))
	}
}

func TestFilterByClass_CaseSensitive_NoMatch(t *testing.T) {
	opts := ClassFilterOptions{CaseSensitive: true}
	resources := []Resource{classResource("General-Purpose")}
	got := FilterByClass(resources, "general-purpose", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0 resources, got %d", len(got))
	}
}

func TestFilterByClasses_ORSemantics(t *testing.T) {
	resources := []Resource{classResource("general-purpose"), classResource("compute"), classResource("memory")}
	got := FilterByClasses(resources, []string{"compute", "memory"}, DefaultClassFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestBuildClassIndex_Lookup(t *testing.T) {
	resources := []Resource{classResource("general-purpose"), classResource("compute")}
	idx := BuildClassIndex(resources)
	got := idx.Lookup("compute")
	if len(got) != 1 || got[0].ID != "compute-id" {
		t.Fatalf("unexpected lookup result: %v", got)
	}
}

func TestBuildClassIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{classResource("General-Purpose")}
	idx := BuildClassIndex(resources)
	got := idx.Lookup("general-purpose")
	if len(got) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(got))
	}
}

func TestBuildClassIndex_LookupMissing(t *testing.T) {
	idx := BuildClassIndex([]Resource{classResource("compute")})
	if got := idx.Lookup("memory"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestBuildClassIndex_Classes(t *testing.T) {
	resources := []Resource{classResource("compute"), classResource("memory")}
	idx := BuildClassIndex(resources)
	classes := idx.Classes()
	if len(classes) != 2 {
		t.Fatalf("expected 2 classes, got %d", len(classes))
	}
}
