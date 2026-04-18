package tfstate

import (
	"testing"
)

func providerResource(id, provider string) Resource {
	return Resource{ID: id, Type: "aws_instance", Provider: provider}
}

func TestFilterByProvider_Match(t *testing.T) {
	resources := []Resource{
		providerResource("a", "aws"),
		providerResource("b", "gcp"),
		providerResource("c", "aws"),
	}
	got := FilterByProvider(resources, "aws", DefaultProviderFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByProvider_EmptyProvider_ReturnsAll(t *testing.T) {
	resources := []Resource{
		providerResource("a", "aws"),
		providerResource("b", "gcp"),
	}
	got := FilterByProvider(resources, "", DefaultProviderFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByProvider_CaseInsensitive(t *testing.T) {
	resources := []Resource{providerResource("a", "AWS")}
	got := FilterByProvider(resources, "aws", DefaultProviderFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByProvider_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{providerResource("a", "AWS")}
	opts := ProviderFilterOptions{CaseSensitive: true}
	got := FilterByProvider(resources, "aws", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByProviders_ORSemantics(t *testing.T) {
	resources := []Resource{
		providerResource("a", "aws"),
		providerResource("b", "gcp"),
		providerResource("c", "azure"),
	}
	got := FilterByProviders(resources, []string{"aws", "gcp"}, DefaultProviderFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByProviders_EmptyList_ReturnsAll(t *testing.T) {
	resources := []Resource{
		providerResource("a", "aws"),
		providerResource("b", "gcp"),
	}
	got := FilterByProviders(resources, []string{}, DefaultProviderFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
