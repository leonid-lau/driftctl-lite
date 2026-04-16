package tfstate

import (
	"sort"
	"testing"
)

func indexResource(id, rtype string, attrs map[string]string) Resource {
	return Resource{ID: id, Type: rtype, Attributes: attrs}
}

func TestBuildTagIndex_LookupByKeyAndValue(t *testing.T) {
	resources := []Resource{
		indexResource("r1", "aws_instance", map[string]string{"tags.env": "prod", "tags.team": "infra"}),
		indexResource("r2", "aws_instance", map[string]string{"tags.env": "staging"}),
		indexResource("r3", "aws_s3_bucket", map[string]string{"tags.env": "prod"}),
	}
	idx := BuildTagIndex(resources)

	got := idx.LookupByTag("env", "prod")
	sort.Strings(got)
	if len(got) != 2 || got[0] != "r1" || got[1] != "r3" {
		t.Errorf("expected [r1 r3], got %v", got)
	}
}

func TestBuildTagIndex_LookupByKeyOnly(t *testing.T) {
	resources := []Resource{
		indexResource("r1", "aws_instance", map[string]string{"tags.env": "prod"}),
		indexResource("r2", "aws_instance", map[string]string{"tags.env": "staging"}),
	}
	idx := BuildTagIndex(resources)

	got := idx.LookupByTag("env", "")
	sort.Strings(got)
	if len(got) != 2 {
		t.Errorf("expected 2 results, got %v", got)
	}
}

func TestBuildTagIndex_EmptyKey_ReturnsNil(t *testing.T) {
	idx := BuildTagIndex([]Resource{})
	if got := idx.LookupByTag("", "prod"); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestBuildTagIndex_Keys(t *testing.T) {
	resources := []Resource{
		indexResource("r1", "aws_instance", map[string]string{"tags.env": "prod", "tags.team": "infra"}),
	}
	idx := BuildTagIndex(resources)
	keys := idx.Keys()
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "env" || keys[1] != "team" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestBuildTagIndex_NonTagAttributes_Ignored(t *testing.T) {
	resources := []Resource{
		indexResource("r1", "aws_instance", map[string]string{"instance_type": "t2.micro", "tags.env": "prod"}),
	}
	idx := BuildTagIndex(resources)
	keys := idx.Keys()
	if len(keys) != 1 || keys[0] != "env" {
		t.Errorf("expected only 'env' key, got %v", keys)
	}
}
