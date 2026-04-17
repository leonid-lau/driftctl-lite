package tfstate

import (
	"testing"
)

func statusResource(id, status string) Resource {
	return Resource{ID: id, Type: "aws_instance", Status: status}
}

func TestFilterByStatus_Match(t *testing.T) {
	resources := []Resource{
		statusResource("a", "created"),
		statusResource("b", "tainted"),
		statusResource("c", "created"),
	}
	got := FilterByStatus(resources, StatusCreated, DefaultStatusFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByStatus_EmptyStatus_ReturnsAll(t *testing.T) {
	resources := []Resource{statusResource("a", "created"), statusResource("b", "tainted")}
	got := FilterByStatus(resources, "", DefaultStatusFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByStatus_CaseInsensitive(t *testing.T) {
	resources := []Resource{statusResource("a", "Created"), statusResource("b", "tainted")}
	got := FilterByStatus(resources, StatusCreated, DefaultStatusFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByStatuses_ORSemantics(t *testing.T) {
	resources := []Resource{
		statusResource("a", "created"),
		statusResource("b", "tainted"),
		statusResource("c", "destroyed"),
	}
	got := FilterByStatuses(resources, []ResourceStatus{StatusCreated, StatusTainted}, DefaultStatusFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildStatusIndex_Lookup(t *testing.T) {
	resources := []Resource{
		statusResource("a", "created"),
		statusResource("b", "tainted"),
		statusResource("c", "created"),
	}
	idx := BuildStatusIndex(resources)
	if got := idx.Lookup("created"); len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
	if got := idx.Lookup("missing"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestBuildStatusIndex_Statuses(t *testing.T) {
	resources := []Resource{
		statusResource("a", "created"),
		statusResource("b", "tainted"),
	}
	idx := BuildStatusIndex(resources)
	if len(idx.Statuses()) != 2 {
		t.Fatalf("expected 2 statuses, got %d", len(idx.Statuses()))
	}
}
