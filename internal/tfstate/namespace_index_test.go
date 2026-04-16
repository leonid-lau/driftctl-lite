package tfstate

import (
	"sort"
	"testing"
)

func idxNSResource(id, ns, rtype string) Resource {
	return Resource{ID: id, Type: rtype, Namespace: ns, Attributes: map[string]interface{}{}}
}

func TestBuildNamespaceIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxNSResource("a", "prod", "aws_instance"),
		idxNSResource("b", "prod", "aws_s3_bucket"),
		idxNSResource("c", "staging", "aws_instance"),
	}
	idx := BuildNamespaceIndex(resources)
	got := idx.Lookup("prod")
	if len(got) != 2 {
		t.Fatalf("expected 2 prod resources, got %d", len(got))
	}
}

func TestBuildNamespaceIndex_LookupMissing(t *testing.T) {
	idx := BuildNamespaceIndex([]Resource{})
	got := idx.Lookup("nonexistent")
	if got != nil && len(got) != 0 {
		t.Fatalf("expected nil/empty, got %v", got)
	}
}

func TestBuildNamespaceIndex_Namespaces(t *testing.T) {
	resources := []Resource{
		idxNSResource("a", "prod", "aws_instance"),
		idxNSResource("b", "staging", "aws_instance"),
		idxNSResource("c", "dev", "aws_s3_bucket"),
	}
	idx := BuildNamespaceIndex(resources)
	ns := idx.Namespaces()
	sort.Strings(ns)
	if len(ns) != 3 {
		t.Fatalf("expected 3 namespaces, got %d", len(ns))
	}
	if ns[0] != "dev" || ns[1] != "prod" || ns[2] != "staging" {
		t.Fatalf("unexpected namespaces: %v", ns)
	}
}

func TestBuildNamespaceIndex_EmptyInput(t *testing.T) {
	idx := BuildNamespaceIndex(nil)
	if len(idx.Namespaces()) != 0 {
		t.Fatal("expected no namespaces")
	}
}
