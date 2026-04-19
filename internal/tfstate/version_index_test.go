package tfstate

import (
	"testing"
)

func idxVerResource(id, version string) Resource {
	return Resource{
		ID:         id,
		Type:       "aws_lambda_function",
		Attributes: map[string]string{"version": version},
	}
}

func TestBuildVersionIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxVerResource("r1", "v1.0.0"),
		idxVerResource("r2", "v2.0.0"),
		idxVerResource("r3", "v1.0.0"),
	}
	idx := BuildVersionIndex(resources)
	got := idx.Lookup("v1.0.0")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildVersionIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxVerResource("r1", "V1.0.0")}
	idx := BuildVersionIndex(resources)
	got := idx.Lookup("v1.0.0")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildVersionIndex_LookupMissing(t *testing.T) {
	idx := BuildVersionIndex([]Resource{idxVerResource("r1", "v1.0.0")})
	if got := idx.Lookup("v9.9.9"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestBuildVersionIndex_Versions(t *testing.T) {
	resources := []Resource{
		idxVerResource("r1", "v1.0.0"),
		idxVerResource("r2", "v2.0.0"),
	}
	idx := BuildVersionIndex(resources)
	if len(idx.Versions()) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(idx.Versions()))
	}
}
