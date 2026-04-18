package tfstate

import (
	"sort"
	"testing"
)

func idxProviderResource(id, provider string) Resource {
	return Resource{ID: id, Provider: provider, Type: "aws_instance", Attributes: map[string]interface{}{}}
}

func TestBuildProviderIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxProviderResource("r1", "aws"),
		idxProviderResource("r2", "aws"),
		idxProviderResource("r3", "gcp"),
	}
	idx := BuildProviderIndex(resources)
	result := idx.Lookup("aws")
	if len(result) != 2 {
		t.Fatalf("expected 2, got %d", len(result))
	}
}

func TestBuildProviderIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxProviderResource("r1", "AWS")}
	idx := BuildProviderIndex(resources)
	result := idx.Lookup("aws")
	if len(result) != 1 {
		t.Fatalf("expected 1, got %d", len(result))
	}
}

func TestBuildProviderIndex_LookupMissing(t *testing.T) {
	idx := BuildProviderIndex([]Resource{})
	if idx.Lookup("aws") != nil {
		t.Fatal("expected nil")
	}
}

func TestBuildProviderIndex_Providers(t *testing.T) {
	resources := []Resource{
		idxProviderResource("r1", "aws"),
		idxProviderResource("r2", "gcp"),
		idxProviderResource("r3", "azure"),
	}
	idx := BuildProviderIndex(resources)
	providers := idx.Providers()
	sort.Strings(providers)
	if len(providers) != 3 {
		t.Fatalf("expected 3 providers, got %d", len(providers))
	}
	if providers[0] != "aws" || providers[1] != "azure" || providers[2] != "gcp" {
		t.Fatalf("unexpected providers: %v", providers)
	}
}

func TestBuildProviderIndex_EmptyInput(t *testing.T) {
	idx := BuildProviderIndex(nil)
	if len(idx.Providers()) != 0 {
		t.Fatal("expected empty providers")
	}
}
