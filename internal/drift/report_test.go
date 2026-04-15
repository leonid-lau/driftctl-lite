package drift_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/driftctl-lite/internal/drift"
)

func TestSummarize_Empty(t *testing.T) {
	s := drift.Summarize(nil)
	if s.Total != 0 {
		t.Errorf("expected total 0, got %d", s.Total)
	}
}

func TestSummarize_Mixed(t *testing.T) {
	results := []drift.DriftResult{
		{Type: drift.DriftMissing, ResourceID: "a"},
		{Type: drift.DriftMissing, ResourceID: "b"},
		{Type: drift.DriftExtra, ResourceID: "c"},
		{Type: drift.DriftModified, ResourceID: "d", Attribute: "x"},
	}
	s := drift.Summarize(results)
	if s.Missing != 2 {
		t.Errorf("expected 2 missing, got %d", s.Missing)
	}
	if s.Extra != 1 {
		t.Errorf("expected 1 extra, got %d", s.Extra)
	}
	if s.Modified != 1 {
		t.Errorf("expected 1 modified, got %d", s.Modified)
	}
	if s.Total != 4 {
		t.Errorf("expected total 4, got %d", s.Total)
	}
}

func TestWriteReport_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	if err := drift.WriteReport(&buf, nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift") {
		t.Errorf("expected 'No drift' in output, got: %q", buf.String())
	}
}

func TestWriteReport_WithDrift(t *testing.T) {
	results := []drift.DriftResult{
		{Type: drift.DriftMissing, ResourceID: "i-123"},
		{Type: drift.DriftModified, ResourceID: "i-456", Attribute: "ami", Expected: "ami-old", Actual: "ami-new"},
	}
	var buf bytes.Buffer
	if err := drift.WriteReport(&buf, results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"MISSING", "MODIFIED", "i-123", "i-456", "ami", "Summary"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in report output", want)
		}
	}
}
