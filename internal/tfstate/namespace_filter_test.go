package tfstate

import (
	"testing"
)

func nsResource(id, ns, rtype string) Resource {
	return Resource{ID: id, Type: rtype, Namespace: ns, Attributes: map[string]interface{}{}}
}

func TestFilterByNamespace_Match(t *testing.T) {
	resources := []Resource{
		nsResource("a", "prod", "aws_instance"),
		nsResource("b", "staging", "aws_instance"),
		nsResource("c", "prod", "aws_s3_bucket"),
	}
	opts := DefaultNamespaceFilterOptions()
	got := FilterByNamespace(resources, "prod", opts)
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByNamespace_EmptyNS_ReturnsAll(t *testing.T) {
	resources := []Resource{
		nsResource("a", "prod", "aws_instance"),
		nsResource("b", "staging", "aws_instance"),
	}
	opts := DefaultNamespaceFilterOptions()
	got := FilterByNamespace(resources, "", opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByNamespace_CaseInsensitive(t *testing.T) {
	resources := []Resource{
		nsResource("a", "Prod", "aws_instance"),
		nsResource("b", "staging", "aws_instance"),
	}
	opts := NamespaceFilterOptions{CaseSensitive: false}
	got := FilterByNamespace(resources, "prod", opts)
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByNamespaces_MultipleMatch(t *testing.T) {
	resources := []Resource{
		nsResource("a", "prod", "aws_instance"),
		nsResource("b", "staging", "aws_instance"),
		nsResource("c", "dev", "aws_s3_bucket"),
	}
	opts := DefaultNamespaceFilterOptions()
	got := FilterByNamespaces(resources, []string{"prod", "dev"}, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByNamespaces_Empty_ReturnsAll(t *testing.T) {
	resources := []Resource{
		nsResource("a", "prod", "aws_instance"),
	}
	opts := DefaultNamespaceFilterOptions()
	got := FilterByNamespaces(resources, nil, opts)
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}
