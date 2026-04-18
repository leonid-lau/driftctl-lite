package tfstate

import (
	"testing"
)

func idxSevResource(id, severity string) Resource {
	return Resource{
		Type:       "aws_instance",
		ID:         id,
		Attributes: map[string]interface{}{"severity": severity},
	}
}

func TestBuildSeverityIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxSevResource("r1", "high"),
		idxSevResource("r2", "low"),
		idxSevResource("r3", "high"),
	}
	idx := BuildSeverityIndex(resources)
	got := idx.Lookup("high")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildSeverityIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxSevResource("r1", "HIGH")}
	idx := BuildSeverityIndex(resources)
	got := idx.Lookup("high")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildSeverityIndex_LookupMissing(t *testing.T) {
	idx := BuildSeverityIndex([]Resource{idxSevResource("r1", "low")})
	got := idx.Lookup("critical")
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestBuildSeverityIndex_Severities(t *testing.T) {
	resources := []Resource{
		idxSevResource("r1", "low"),
		idxSevResource("r2", "high"),
		idxSevResource("r3", "low"),
	}
	idx := BuildSeverityIndex(resources)
	sevs := idx.Severities()
	if len(sevs) != 2 {
		t.Fatalf("expected 2 distinct severities, got %d", len(sevs))
	}
}

func TestBuildSeverityIndex_EmptyInput(t *testing.T) {
	idx := BuildSeverityIndex(nil)
	if len(idx.Severities()) != 0 {
		t.Fatal("expected empty index")
	}
}
