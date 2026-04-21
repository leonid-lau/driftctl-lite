package tfstate

import (
	"testing"
)

func policyResource(policy string) Resource {
	attrs := map[string]string{}
	if policy != "" {
		attrs["policy"] = policy
	}
	return Resource{ID: policy + "-id", Type: "aws_iam_policy", Attributes: attrs}
}

func TestFilterByPolicy_Match(t *testing.T) {
	resources := []Resource{policyResource("ReadOnly"), policyResource("FullAccess"), policyResource("ReadOnly")}
	got := FilterByPolicy(resources, "ReadOnly", DefaultPolicyFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByPolicy_EmptyPolicy_ReturnsAll(t *testing.T) {
	resources := []Resource{policyResource("ReadOnly"), policyResource("FullAccess")}
	got := FilterByPolicy(resources, "", DefaultPolicyFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByPolicy_CaseInsensitive(t *testing.T) {
	resources := []Resource{policyResource("readonly")}
	got := FilterByPolicy(resources, "READONLY", DefaultPolicyFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByPolicy_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{policyResource("readonly")}
	opts := PolicyFilterOptions{CaseSensitive: true}
	got := FilterByPolicy(resources, "READONLY", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByPolicies_ORSemantics(t *testing.T) {
	resources := []Resource{policyResource("ReadOnly"), policyResource("FullAccess"), policyResource("WriteOnly")}
	got := FilterByPolicies(resources, []string{"ReadOnly", "WriteOnly"}, DefaultPolicyFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByPolicies_EmptyList_ReturnsAll(t *testing.T) {
	resources := []Resource{policyResource("ReadOnly"), policyResource("FullAccess")}
	got := FilterByPolicies(resources, []string{}, DefaultPolicyFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
