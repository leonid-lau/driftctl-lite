package tfstate

import (
	"testing"
)

func archResource(id, arch string) Resource {
	return Resource{
		ID:         id,
		Type:       "aws_instance",
		Attributes: map[string]interface{}{"arch": arch},
	}
}

func TestFilterByArch_Match(t *testing.T) {
	resources := []Resource{
		archResource("r1", "x86_64"),
		archResource("r2", "arm64"),
		archResource("r3", "x86_64"),
	}
	got := FilterByArch(resources, "x86_64", DefaultArchFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByArch_EmptyArch_ReturnsAll(t *testing.T) {
	resources := []Resource{
		archResource("r1", "x86_64"),
		archResource("r2", "arm64"),
	}
	got := FilterByArch(resources, "", DefaultArchFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByArch_CaseInsensitive(t *testing.T) {
	resources := []Resource{
		archResource("r1", "ARM64"),
		archResource("r2", "x86_64"),
	}
	got := FilterByArch(resources, "arm64", DefaultArchFilterOptions())
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected r1, got %+v", got)
	}
}

func TestFilterByArch_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{
		archResource("r1", "ARM64"),
	}
	opts := ArchFilterOptions{CaseSensitive: true}
	got := FilterByArch(resources, "arm64", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByArchs_ORSemantics(t *testing.T) {
	resources := []Resource{
		archResource("r1", "x86_64"),
		archResource("r2", "arm64"),
		archResource("r3", "s390x"),
	}
	got := FilterByArchs(resources, []string{"x86_64", "arm64"}, DefaultArchFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByArch_FallbackToArchitectureKey(t *testing.T) {
	r := Resource{
		ID:         "r1",
		Type:       "aws_instance",
		Attributes: map[string]interface{}{"architecture": "arm64"},
	}
	got := FilterByArch([]Resource{r}, "arm64", DefaultArchFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}
