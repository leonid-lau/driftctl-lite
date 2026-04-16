package tfstate

import (
	"sort"
	"testing"
)

func idxLabelResource(id, rtype string, labels map[string]string) Resource {
	return Resource{ID: id, Type: rtype, Labels: labels}
}

func TestBuildLabelIndex_LookupByKeyAndValue(t *testing.T) {
	resources := []Resource{
		idxLabelResource("a", "aws_s3_bucket", map[string]string{"env": "prod"}),
		idxLabelResource("b", "aws_instance", map[string]string{"env": "dev"}),
	}
	idx := BuildLabelIndex(resources)
	got := idx.Lookup("env", "prod")
	if len(got) != 1 || got[0].ID != "a" {
		t.Fatalf("expected resource a, got %+v", got)
	}
}

func TestBuildLabelIndex_LookupByKeyOnly(t *testing.T) {
	resources := []Resource{
		idxLabelResource("a", "aws_s3_bucket", map[string]string{"env": "prod"}),
		idxLabelResource("b", "aws_instance", map[string]string{"env": "dev"}),
		idxLabelResource("c", "aws_vpc", map[string]string{"team": "infra"}),
	}
	idx := BuildLabelIndex(resources)
	got := idx.Lookup("env", "")
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestBuildLabelIndex_EmptyKey_ReturnsNil(t *testing.T) {
	resources := []Resource{
		idxLabelResource("a", "aws_s3_bucket", map[string]string{"env": "prod"}),
	}
	idx := BuildLabelIndex(resources)
	got := idx.Lookup("", "prod")
	if got != nil {
		t.Fatalf("expected nil for empty key, got %+v", got)
	}
}

func TestBuildLabelIndex_Keys(t *testing.T) {
	resources := []Resource{
		idxLabelResource("a", "aws_s3_bucket", map[string]string{"env": "prod", "team": "infra"}),
	}
	idx := BuildLabelIndex(resources)
	keys := idx.Keys()
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "env" || keys[1] != "team" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestBuildLabelIndex_MissingKey_ReturnsNil(t *testing.T) {
	resources := []Resource{
		idxLabelResource("a", "aws_s3_bucket", map[string]string{"env": "prod"}),
	}
	idx := BuildLabelIndex(resources)
	got := idx.Lookup("nonexistent", "")
	if got != nil {
		t.Fatalf("expected nil for missing key, got %+v", got)
	}
}
