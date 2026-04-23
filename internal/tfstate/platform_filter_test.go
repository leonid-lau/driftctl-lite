package tfstate

import (
	"testing"
)

func platformResource(platform string) Resource {
	attrs := map[string]interface{}{}
	if platform != "" {
		attrs["platform"] = platform
	}
	return Resource{Type: "aws_instance", ID: platform, Attributes: attrs}
}

func TestFilterByPlatform_Match(t *testing.T) {
	resources := []Resource{
		platformResource("linux"),
		platformResource("windows"),
		platformResource("darwin"),
	}
	got := FilterByPlatform(resources, "linux", DefaultPlatformFilterOptions())
	if len(got) != 1 || got[0].Attributes["platform"] != "linux" {
		t.Fatalf("expected 1 linux resource, got %v", got)
	}
}

func TestFilterByPlatform_EmptyPlatform_ReturnsAll(t *testing.T) {
	resources := []Resource{platformResource("linux"), platformResource("windows")}
	got := FilterByPlatform(resources, "", DefaultPlatformFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByPlatform_CaseInsensitive(t *testing.T) {
	resources := []Resource{platformResource("Linux"), platformResource("WINDOWS")}
	got := FilterByPlatform(resources, "linux", DefaultPlatformFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(got))
	}
}

func TestFilterByPlatform_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{platformResource("Linux")}
	opts := PlatformFilterOptions{CaseSensitive: true}
	got := FilterByPlatform(resources, "linux", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0 resources, got %d", len(got))
	}
}

func TestFilterByPlatforms_ORSemantics(t *testing.T) {
	resources := []Resource{
		platformResource("linux"),
		platformResource("windows"),
		platformResource("darwin"),
	}
	got := FilterByPlatforms(resources, []string{"linux", "darwin"}, DefaultPlatformFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestBuildPlatformIndex_Lookup(t *testing.T) {
	resources := []Resource{platformResource("linux"), platformResource("windows")}
	idx := BuildPlatformIndex(resources)
	result := idx.Lookup("linux")
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestBuildPlatformIndex_LookupMissing(t *testing.T) {
	idx := BuildPlatformIndex([]Resource{platformResource("linux")})
	if idx.Lookup("freebsd") != nil {
		t.Fatal("expected nil for missing platform")
	}
}

func TestBuildPlatformIndex_Platforms(t *testing.T) {
	resources := []Resource{platformResource("linux"), platformResource("windows")}
	idx := BuildPlatformIndex(resources)
	if len(idx.Platforms()) != 2 {
		t.Fatalf("expected 2 platforms, got %d", len(idx.Platforms()))
	}
}
