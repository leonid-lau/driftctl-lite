package tfstate

import (
	"sort"
	"testing"
)

func idxModuleResource(module string) Resource {
	return Resource{
		Type:       "aws_instance",
		Attributes: map[string]interface{}{"module": module},
	}
}

func TestBuildModuleIndex_Lookup(t *testing.T) {
	resources := []Resource{idxModuleResource("vpc"), idxModuleResource("ecs"), idxModuleResource("vpc")}
	idx := BuildModuleIndex(resources)
	got := idx.Lookup("vpc")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildModuleIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxModuleResource("VPC")}
	idx := BuildModuleIndex(resources)
	got := idx.Lookup("vpc")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildModuleIndex_LookupMissing(t *testing.T) {
	idx := BuildModuleIndex([]Resource{idxModuleResource("ecs")})
	got := idx.Lookup("vpc")
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestBuildModuleIndex_Modules(t *testing.T) {
	resources := []Resource{idxModuleResource("vpc"), idxModuleResource("ecs")}
	idx := BuildModuleIndex(resources)
	modules := idx.Modules()
	sort.Strings(modules)
	if len(modules) != 2 || modules[0] != "ecs" || modules[1] != "vpc" {
		t.Fatalf("unexpected modules: %v", modules)
	}
}

func TestBuildModuleIndex_EmptyInput(t *testing.T) {
	idx := BuildModuleIndex(nil)
	if len(idx.Modules()) != 0 {
		t.Fatal("expected empty index")
	}
}
