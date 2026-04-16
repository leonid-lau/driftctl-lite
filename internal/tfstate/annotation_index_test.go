package tfstate

import (
	"sort"
	"testing"
)

func idxAnnoResource(id, rtype string, annotations map[string]string) Resource {
	return Resource{
		ID:          id,
		Type:        rtype,
		Annotations: annotations,
	}
}

func TestBuildAnnotationIndex_LookupByKeyAndValue(t *testing.T) {
	resources := []Resource{
		idxAnnoResource("r1", "aws_s3_bucket", map[string]string{"env": "prod"}),
		idxAnnoResource("r2", "aws_s3_bucket", map[string]string{"env": "staging"}),
		idxAnnoResource("r3", "aws_instance", map[string]string{"env": "prod"}),
	}
	idx := BuildAnnotationIndex(resources)
	result := idx.Lookup("env", "prod")
	if len(result) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(result))
	}
}

func TestBuildAnnotationIndex_LookupByKeyOnly(t *testing.T) {
	resources := []Resource{
		idxAnnoResource("r1", "aws_s3_bucket", map[string]string{"team": "alpha"}),
		idxAnnoResource("r2", "aws_s3_bucket", map[string]string{"team": "beta"}),
		idxAnnoResource("r3", "aws_instance", map[string]string{"owner": "alice"}),
	}
	idx := BuildAnnotationIndex(resources)
	result := idx.Lookup("team", "")
	if len(result) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(result))
	}
}

func TestBuildAnnotationIndex_EmptyKey_ReturnsNil(t *testing.T) {
	resources := []Resource{
		idxAnnoResource("r1", "aws_s3_bucket", map[string]string{"env": "prod"}),
	}
	idx := BuildAnnotationIndex(resources)
	result := idx.Lookup("", "prod")
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestBuildAnnotationIndex_Keys(t *testing.T) {
	resources := []Resource{
		idxAnnoResource("r1", "aws_s3_bucket", map[string]string{"env": "prod", "team": "alpha"}),
		idxAnnoResource("r2", "aws_instance", map[string]string{"owner": "alice"}),
	}
	idx := BuildAnnotationIndex(resources)
	keys := idx.Keys()
	sort.Strings(keys)
	expected := []string{"env", "owner", "team"}
	if len(keys) != len(expected) {
		t.Fatalf("expected keys %v, got %v", expected, keys)
	}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("key[%d]: expected %q, got %q", i, expected[i], k)
		}
	}
}

func TestBuildAnnotationIndex_MissingKey_ReturnsNil(t *testing.T) {
	resources := []Resource{
		idxAnnoResource("r1", "aws_s3_bucket", map[string]string{"env": "prod"}),
	}
	idx := BuildAnnotationIndex(resources)
	result := idx.Lookup("nonexistent", "")
	if result != nil {
		t.Fatalf("expected nil for missing key, got %v", result)
	}
}
