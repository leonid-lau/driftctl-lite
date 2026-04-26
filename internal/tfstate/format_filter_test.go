package tfstate

import (
	"testing"
)

func formatResource(id, format string) Resource {
	return Resource{
		Type: "test_resource",
		Name: id,
		Attributes: map[string]string{"format": format},
	}
}

func TestFilterByFormat_Match(t *testing.T) {
	res := []Resource{
		formatResource("a", "json"),
		formatResource("b", "yaml"),
		formatResource("c", "json"),
	}
	got := FilterByFormat(res, "json", DefaultFormatFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByFormat_EmptyFormat_ReturnsAll(t *testing.T) {
	res := []Resource{
		formatResource("a", "json"),
		formatResource("b", "yaml"),
	}
	got := FilterByFormat(res, "", DefaultFormatFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByFormat_CaseInsensitive(t *testing.T) {
	res := []Resource{formatResource("a", "JSON")}
	got := FilterByFormat(res, "json", DefaultFormatFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByFormat_CaseSensitive_NoMatch(t *testing.T) {
	res := []Resource{formatResource("a", "JSON")}
	got := FilterByFormat(res, "json", FilterOptions{CaseInsensitive: false})
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByFormats_ORSemantics(t *testing.T) {
	res := []Resource{
		formatResource("a", "json"),
		formatResource("b", "yaml"),
		formatResource("c", "toml"),
	}
	got := FilterByFormats(res, []string{"json", "yaml"}, DefaultFormatFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildFormatIndex_Lookup(t *testing.T) {
	res := []Resource{
		formatResource("a", "json"),
		formatResource("b", "yaml"),
		formatResource("c", "JSON"),
	}
	idx := BuildFormatIndex(res)
	got := idx.Lookup("json")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildFormatIndex_LookupMissing(t *testing.T) {
	idx := BuildFormatIndex([]Resource{formatResource("a", "json")})
	if idx.Lookup("xml") != nil {
		t.Fatal("expected nil for missing format")
	}
}

func TestBuildFormatIndex_Formats(t *testing.T) {
	res := []Resource{
		formatResource("a", "json"),
		formatResource("b", "yaml"),
	}
	idx := BuildFormatIndex(res)
	if len(idx.Formats()) != 2 {
		t.Fatalf("expected 2 formats, got %d", len(idx.Formats()))
	}
}
