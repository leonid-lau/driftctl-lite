package tfstate

import (
	"testing"
)

func idxPolicyResource(policy string) Resource {
	attrs := map[string]string{}
	if policy != "" {
		attrs["policy"] = policy
	}
	return Resource{ID: policy + "-id", Type: "aws_iam_policy", Attributes: attrs}
}

func TestBuildPolicyIndex_Lookup(t *testing.T) {
	resources := []Resource{idxPolicyResource("ReadOnly"), idxPolicyResource("FullAccess")}
	idx := BuildPolicyIndex(resources)
	got := idx.Lookup("ReadOnly")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildPolicyIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxPolicyResource("readonly")}
	idx := BuildPolicyIndex(resources)
	got := idx.Lookup("READONLY")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildPolicyIndex_LookupMissing(t *testing.T) {
	resources := []Resource{idxPolicyResource("ReadOnly")}
	idx := BuildPolicyIndex(resources)
	got := idx.Lookup("nonexistent")
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestBuildPolicyIndex_Policies(t *testing.T) {
	resources := []Resource{idxPolicyResource("ReadOnly"), idxPolicyResource("FullAccess")}
	idx := BuildPolicyIndex(resources)
	policies := idx.Policies()
	if len(policies) != 2 {
		t.Fatalf("expected 2 policies, got %d", len(policies))
	}
}

func TestBuildPolicyIndex_EmptyInput(t *testing.T) {
	idx := BuildPolicyIndex([]Resource{})
	if len(idx.Policies()) != 0 {
		t.Fatal("expected empty index")
	}
}
