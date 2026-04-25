package tfstate

import (
	"sort"
	"testing"
)

func idxImageResource(image string) Resource {
	return Resource{
		Type:       "aws_instance",
		Attributes: map[string]interface{}{"image": image},
	}
}

func TestBuildImageIndex_Lookup(t *testing.T) {
	resources := []Resource{idxImageResource("ami-abc"), idxImageResource("ami-xyz")}
	idx := BuildImageIndex(resources)
	got := idx.Lookup("ami-abc")
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestBuildImageIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxImageResource("AMI-ABC")}
	idx := BuildImageIndex(resources)
	got := idx.Lookup("ami-abc")
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
}

func TestBuildImageIndex_LookupMissing(t *testing.T) {
	idx := BuildImageIndex([]Resource{idxImageResource("ami-111")})
	got := idx.Lookup("ami-999")
	if len(got) != 0 {
		t.Fatalf("expected 0 results, got %d", len(got))
	}
}

func TestBuildImageIndex_Images(t *testing.T) {
	resources := []Resource{idxImageResource("ami-a"), idxImageResource("ami-b"), idxImageResource("ami-a")}
	idx := BuildImageIndex(resources)
	images := idx.Images()
	sort.Strings(images)
	if len(images) != 2 {
		t.Fatalf("expected 2 distinct images, got %d: %v", len(images), images)
	}
}

func TestBuildImageIndex_EmptyInput(t *testing.T) {
	idx := BuildImageIndex(nil)
	if len(idx.Images()) != 0 {
		t.Fatal("expected empty index")
	}
}
