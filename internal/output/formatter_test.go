package output_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/driftctl-lite/internal/drift"
	"github.com/driftctl-lite/internal/output"
)

func makeReport(diffs []drift.Diff) drift.Report {
	return drift.Report{Diffs: diffs}
}

func TestNewFormatter_ValidFormats(t *testing.T) {
	for _, fmt := range []output.Format{output.FormatText, output.FormatJSON} {
		_, err := output.NewFormatter(fmt, &bytes.Buffer{})
		if err != nil {
			t.Errorf("expected no error for format %q, got %v", fmt, err)
		}
	}
}

func TestNewFormatter_InvalidFormat(t *testing.T) {
	_, err := output.NewFormatter("xml", &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_TextNoDrift(t *testing.T) {
	var buf bytes.Buffer
	f, _ := output.NewFormatter(output.FormatText, &buf)
	report := makeReport(nil)
	if err := f.Write(report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No drift detected") {
		t.Errorf("expected 'No drift detected' in output, got: %s", buf.String())
	}
}

func TestWrite_TextWithDrift(t *testing.T) {
	var buf bytes.Buffer
	f, _ := output.NewFormatter(output.FormatText, &buf)
	report := makeReport([]drift.Diff{
		{ResourceID: "aws_s3_bucket.my_bucket", ChangeType: "modified", Detail: "tags changed"},
	})
	if err := f.Write(report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "aws_s3_bucket.my_bucket") {
		t.Errorf("expected resource ID in output, got: %s", out)
	}
}

func TestWrite_JSONOutput(t *testing.T) {
	var buf bytes.Buffer
	f, _ := output.NewFormatter(output.FormatJSON, &buf)
	report := makeReport([]drift.Diff{
		{ResourceID: "aws_instance.web", ChangeType: "missing", Detail: "not found in cloud"},
	})
	if err := f.Write(report); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var decoded drift.Report
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if len(decoded.Diffs) != 1 {
		t.Errorf("expected 1 diff, got %d", len(decoded.Diffs))
	}
}
