package tfstate

import (
	"testing"

	"github.com/snyk/driftctl-lite/internal/tfstate"
)

func verResource(id, version string) tfstate.Resource {
	return tfstate.Resource{
		Type: "aws_instance",
		Name: id,
		Attributes: map[string]string{"version": version},
	}
}

func TestFilterByVersion_Match(t *testing.T) {
	resources := []tfstate.Resource{
		verResource("a", "v1"),
		verResource("b", "v2"),
		verResource("c", "v1"),
	}
	got := FilterByVersion(resources, "v1", DefaultVersionFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByVersion_EmptyVersion_ReturnsAll(t *testing.T) {
	resources := []tfstate.Resource{
		verResource("a", "v1"),
		verResource("b", "v2"),
	}
	got := FilterByVersion(resources, "", DefaultVersionFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByVersion_CaseInsensitive(t *testing.T) {
	resources := []tfstate.Resource{
		verResource("a", "V1"),
	}
	got := FilterByVersion(resources, "v1", DefaultVersionFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByVersion_CaseSensitive_NoMatch(t *testing.T) {
	resources := []tfstate.Resource{
		verResource("a", "V1"),
	}
	opts := VersionFilterOptions{CaseSensitive: true}
	got := FilterByVersion(resources, "v1", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByVersions_ORSemantics(t *testing.T) {
	resources := []tfstate.Resource{
		verResource("a", "v1"),
		verResource("b", "v2"),
		verResource("c", "v3"),
	}
	got := FilterByVersions(resources, []string{"v1", "v3"}, DefaultVersionFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
