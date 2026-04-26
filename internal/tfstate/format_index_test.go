package tfstate

import (
	"testing"
)

func idxFormatResource(id, format string) Resource {
	return Resource{
		ID:   id,
		Type: "test_resource",
		Attributes: map[string]interface{}{
			"format": format,
		},
	}
}

func TestBuildFormatIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxFormatResource("r1", "json"),
		idxFormatResource("r2", "yaml"),
		idxFormatResource("r3", "json"),
	}

	idx := BuildFormatIndex(resources)

	got := idx.Lookup("json")
	if len(got) != 2 {
		t.Fatalf("expected 2 resources for format 'json', got %d", len(got))
	}
}

func TestBuildFormatIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{
		idxFormatResource("r1", "JSON"),
		idxFormatResource("r2", "Yaml"),
	}

	idx := BuildFormatIndex(resources)

	got := idx.Lookup("json")
	if len(got) != 1 {
		t.Fatalf("expected 1 resource for format 'json' (case-insensitive), got %d", len(got))
	}
	if got[0].ID != "r1" {
		t.Errorf("expected resource ID 'r1', got %q", got[0].ID)
	}
}

func TestBuildFormatIndex_LookupMissing(t *testing.T) {
	resources := []Resource{
		idxFormatResource("r1", "json"),
	}

	idx := BuildFormatIndex(resources)

	got := idx.Lookup("xml")
	if len(got) != 0 {
		t.Errorf("expected 0 resources for missing format 'xml', got %d", len(got))
	}
}

func TestBuildFormatIndex_Formats(t *testing.T) {
	resources := []Resource{
		idxFormatResource("r1", "json"),
		idxFormatResource("r2", "yaml"),
		idxFormatResource("r3", "toml"),
		idxFormatResource("r4", "json"),
	}

	idx := BuildFormatIndex(resources)

	formats := idx.Formats()
	if len(formats) != 3 {
		t.Errorf("expected 3 distinct formats, got %d: %v", len(formats), formats)
	}

	formatSet := make(map[string]bool)
	for _, f := range formats {
		formatSet[f] = true
	}
	for _, expected := range []string{"json", "yaml", "toml"} {
		if !formatSet[expected] {
			t.Errorf("expected format %q in index formats, but not found", expected)
		}
	}
}

func TestBuildFormatIndex_EmptyInput(t *testing.T) {
	idx := BuildFormatIndex(nil)

	got := idx.Lookup("json")
	if len(got) != 0 {
		t.Errorf("expected 0 resources for empty index, got %d", len(got))
	}

	formats := idx.Formats()
	if len(formats) != 0 {
		t.Errorf("expected 0 formats for empty index, got %d", len(formats))
	}
}
