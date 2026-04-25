package tfstate

import (
	"testing"
)

func groupResource(group string) Resource {
	attrs := map[string]interface{}{}
	if group != "" {
		attrs["group"] = group
	}
	return Resource{Type: "aws_iam_group", ID: group, Attributes: attrs}
}

func TestFilterByGroup_Match(t *testing.T) {
	resources := []Resource{
		groupResource("admins"),
		groupResource("developers"),
		groupResource("ops"),
	}
	got := FilterByGroup(resources, "admins", DefaultGroupFilterOptions())
	if len(got) != 1 || got[0].ID != "admins" {
		t.Fatalf("expected 1 resource 'admins', got %v", got)
	}
}

func TestFilterByGroup_EmptyGroup_ReturnsAll(t *testing.T) {
	resources := []Resource{groupResource("admins"), groupResource("ops")}
	got := FilterByGroup(resources, "", DefaultGroupFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByGroup_CaseInsensitive(t *testing.T) {
	resources := []Resource{groupResource("Admins"), groupResource("ops")}
	got := FilterByGroup(resources, "admins", DefaultGroupFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(got))
	}
}

func TestFilterByGroup_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{groupResource("Admins")}
	opts := GroupFilterOptions{CaseSensitive: true}
	got := FilterByGroup(resources, "admins", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0 resources, got %d", len(got))
	}
}

func TestFilterByGroups_ORSemantics(t *testing.T) {
	resources := []Resource{
		groupResource("admins"),
		groupResource("developers"),
		groupResource("ops"),
	}
	got := FilterByGroups(resources, []string{"admins", "ops"}, DefaultGroupFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByGroups_EmptySlice_ReturnsAll(t *testing.T) {
	resources := []Resource{groupResource("admins"), groupResource("ops")}
	got := FilterByGroups(resources, nil, DefaultGroupFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}
