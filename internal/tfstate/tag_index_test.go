package tfstate

import (
	"sort"
	"testing"
)

func indexResource(id, tagKey, tagVal string) Resource {
	attrs := map[string]string{}
	if tagKey != "" {
		attrs["tags."+tagKey] = tagVal
	}
	return Resource{ID: id, Type: "aws_instance", Attributes: attrs}
}

func TestBuildTagIndex_LookupByKeyAndValue(t *testing.T) {
	resources := []Resource{
		indexResource("r1", "env", "prod"),
		indexResource("r2", "env", "dev"),
	}
	idx := BuildTagIndex(resources)
	got := idx.Lookup("env", "prod")
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected r1, got %+v", got)
	}
}

func TestBuildTagIndex_LookupByKeyOnly(t *testing.T) {
	resources := []Resource{
		indexResource("r1", "env", "prod"),
		indexResource("r2", "env", "dev"),
		indexResource("r3", "team", "sre"),
	}
	idx := BuildTagIndex(resources)
	got := idx.Lookup("env", "")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildTagIndex_EmptyKey_ReturnsNil(t *testing.T) {
	resources := []Resource{indexResource("r1", "env", "prod")}
	idx := BuildTagIndex(resources)
	got := idx.Lookup("", "")
	if got != nil {
		t.Fatalf("expected nil, got %+v", got)
	}
}

func TestBuildTagIndex_Keys(t *testing.T) {
	resources := []Resource{
		indexResource("r1", "env", "prod"),
		indexResource("r2", "team", "sre"),
	}
	idx := BuildTagIndex(resources)
	keys := idx.Keys()
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "env" || keys[1] != "team" {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestBuildTagIndex_FallbackTagPrefix(t *testing.T) {
	r := Resource{ID: "r1", Type: "aws_instance", Attributes: map[string]string{
		"tag.env": "staging",
	}}
	idx := BuildTagIndex([]Resource{r})
	got := idx.Lookup("env", "staging")
	if len(got) != 1 || got[0].ID != "r1" {
		t.Fatalf("expected r1 via tag. prefix, got %+v", got)
	}
}
