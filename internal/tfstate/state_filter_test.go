package tfstate

import (
	"testing"
)

func stateResource(state string) Resource {
	attrs := map[string]interface{}{}
	if state != "" {
		attrs["state"] = state
	}
	return Resource{
		Type:       "aws_instance",
		Name:       "test",
		Attributes: attrs,
	}
}

func TestFilterByState_Match(t *testing.T) {
	resources := []Resource{
		stateResource("running"),
		stateResource("stopped"),
		stateResource("running"),
	}
	got := FilterByState(resources, "running", DefaultStateFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByState_EmptyState_ReturnsAll(t *testing.T) {
	resources := []Resource{stateResource("running"), stateResource("stopped")}
	got := FilterByState(resources, "", DefaultStateFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByState_CaseInsensitive(t *testing.T) {
	resources := []Resource{stateResource("Running"), stateResource("STOPPED")}
	got := FilterByState(resources, "running", DefaultStateFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByState_CaseSensitive_NoMatch(t *testing.T) {
	opts := StateFilterOptions{CaseSensitive: true}
	resources := []Resource{stateResource("Running")}
	got := FilterByState(resources, "running", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByStates_ORSemantics(t *testing.T) {
	resources := []Resource{
		stateResource("running"),
		stateResource("stopped"),
		stateResource("pending"),
	}
	got := FilterByStates(resources, []string{"running", "stopped"}, DefaultStateFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByStates_EmptyTargets_ReturnsAll(t *testing.T) {
	resources := []Resource{stateResource("running"), stateResource("stopped")}
	got := FilterByStates(resources, nil, DefaultStateFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByState_MissingAttribute(t *testing.T) {
	resources := []Resource{stateResource("")}
	got := FilterByState(resources, "running", DefaultStateFilterOptions())
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}
