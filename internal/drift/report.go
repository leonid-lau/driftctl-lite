package drift

import (
	"encoding/json"
	"fmt"
	"io"
)

// Report holds the full drift detection result.
type Report struct {
	Diffs []Diff `json:"diffs"`
}

// Diff represents a single detected drift between state and live cloud.
type Diff struct {
	ResourceID string `json:"resource_id"`
	ChangeType string `json:"change_type"`
	Detail     string `json:"detail"`
}

// Summarize returns a human-readable one-line summary of the report.
func Summarize(r Report) string {
	if len(r.Diffs) == 0 {
		return "OK — no drift detected"
	}
	counts := map[string]int{}
	for _, d := range r.Diffs {
		counts[d.ChangeType]++
	}
	return fmt.Sprintf("DRIFT — missing: %d, extra: %d, modified: %d",
		counts["missing"], counts["extra"], counts["modified"])
}

// WriteReport serialises the report as indented JSON to the given writer.
func WriteReport(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}
