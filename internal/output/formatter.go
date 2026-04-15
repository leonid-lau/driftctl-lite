package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/driftctl-lite/internal/drift"
)

// Format represents the output format type.
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

// Formatter writes drift results to a writer in the specified format.
type Formatter struct {
	format Format
	writer io.Writer
}

// NewFormatter creates a new Formatter with the given format and writer.
func NewFormatter(format Format, w io.Writer) (*Formatter, error) {
	switch format {
	case FormatText, FormatJSON:
		return &Formatter{format: format, writer: w}, nil
	default:
		return nil, fmt.Errorf("unsupported output format: %q", format)
	}
}

// Write outputs the drift report in the configured format.
func (f *Formatter) Write(report drift.Report) error {
	switch f.format {
	case FormatJSON:
		return f.writeJSON(report)
	default:
		return f.writeText(report)
	}
}

func (f *Formatter) writeJSON(report drift.Report) error {
	enc := json.NewEncoder(f.writer)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func (f *Formatter) writeText(report drift.Report) error {
	w := tabwriter.NewWriter(f.writer, 0, 0, 2, ' ', 0)
	defer w.Flush()

	summary := drift.Summarize(report)
	fmt.Fprintf(w, "Drift Summary:\t%s\n", summary)
	fmt.Fprintf(w, "Total Changes:\t%d\n", len(report.Diffs))

	if len(report.Diffs) == 0 {
		fmt.Fprintln(w, "No drift detected.")
		return nil
	}

	fmt.Fprintln(w, "\nDrifted Resources:")
	for _, d := range report.Diffs {
		fmt.Fprintf(w, "  [%s]\t%s\t-> %s\n", d.ChangeType, d.ResourceID, d.Detail)
	}
	return nil
}
