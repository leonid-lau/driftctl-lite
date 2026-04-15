package drift

import (
	"fmt"

	"github.com/driftctl-lite/internal/tfstate"
)

// DriftType indicates the kind of drift detected.
type DriftType string

const (
	DriftMissing  DriftType = "MISSING"   // resource in state but not in live
	DriftExtra    DriftType = "EXTRA"     // resource in live but not in state
	DriftModified DriftType = "MODIFIED" // resource exists in both but attributes differ
)

// DriftResult represents a single drift finding.
type DriftResult struct {
	Type       DriftType
	ResourceID string
	Attribute  string
	Expected   interface{}
	Actual     interface{}
}

func (d DriftResult) String() string {
	switch d.Type {
	case DriftMissing:
		return fmt.Sprintf("[%s] resource %q exists in state but not in live cloud", d.Type, d.ResourceID)
	case DriftExtra:
		return fmt.Sprintf("[%s] resource %q exists in live cloud but not in state", d.Type, d.ResourceID)
	case DriftModified:
		return fmt.Sprintf("[%s] resource %q attribute %q: expected=%v actual=%v",
			d.Type, d.ResourceID, d.Attribute, d.Expected, d.Actual)
	}
	return fmt.Sprintf("[UNKNOWN] %q", d.ResourceID)
}

// Detector compares Terraform state resources against live cloud resources.
type Detector struct{}

// NewDetector creates a new Detector.
func NewDetector() *Detector {
	return &Detector{}
}

// Detect compares stateIndex (id -> resource from tfstate) with liveIndex
// (id -> attribute map from live cloud) and returns all drift findings.
func (d *Detector) Detect(
	stateIndex map[string]tfstate.Resource,
	liveIndex map[string]map[string]interface{},
) []DriftResult {
	var results []DriftResult

	// Check for MISSING and MODIFIED resources.
	for id, stateRes := range stateIndex {
		liveAttrs, found := liveIndex[id]
		if !found {
			results = append(results, DriftResult{
				Type:       DriftMissing,
				ResourceID: id,
			})
			continue
		}
		// Compare attributes present in state.
		for k, expectedVal := range stateRes.Attributes {
			actualVal, exists := liveAttrs[k]
			if !exists || fmt.Sprintf("%v", actualVal) != fmt.Sprintf("%v", expectedVal) {
				results = append(results, DriftResult{
					Type:       DriftModified,
					ResourceID: id,
					Attribute:  k,
					Expected:   expectedVal,
					Actual:     actualVal,
				})
			}
		}
	}

	// Check for EXTRA resources.
	for id := range liveIndex {
		if _, found := stateIndex[id]; !found {
			results = append(results, DriftResult{
				Type:       DriftExtra,
				ResourceID: id,
			})
		}
	}

	return results
}
