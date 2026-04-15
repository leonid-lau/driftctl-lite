package drift

import (
	"fmt"
	"io"
	"text/tabwriter"
)

// Summary holds aggregated counts of drift findings.
type Summary struct {
	Missing  int
	Extra    int
	Modified int
	Total    int
}

// Summarize computes a Summary from a slice of DriftResults.
func Summarize(results []DriftResult) Summary {
	var s Summary
	for _, r := range results {
		switch r.Type {
		case DriftMissing:
			s.Missing++
		case DriftExtra:
			s.Extra++
		case DriftModified:
			s.Modified++
		}
	}
	s.Total = s.Missing + s.Extra + s.Modified
	return s
}

// WriteReport writes a human-readable drift report to w.
func WriteReport(w io.Writer, results []DriftResult) error {
	if len(results) == 0 {
		_, err := fmt.Fprintln(w, "✓ No drift detected.")
		return err
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "TYPE\tRESOURCE ID\tATTRIBUTE\tEXPECTED\tACTUAL")
	fmt.Fprintln(tw, "----\t-----------\t---------\t--------\t------")
	for _, r := range results {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%v\t%v\n",
			r.Type, r.ResourceID, r.Attribute, r.Expected, r.Actual)
	}
	if err := tw.Flush(); err != nil {
		return err
	}

	s := Summarize(results)
	fmt.Fprintf(w, "\nSummary: %d missing, %d extra, %d modified (%d total)\n",
		s.Missing, s.Extra, s.Modified, s.Total)
	return nil
}
