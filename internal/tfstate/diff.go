package tfstate

import (
	"encoding/json"
	"fmt"
	"sort"
)

// AttributeDiff represents a change in a single attribute between state and live.
type AttributeDiff struct {
	Key      string
	OldValue interface{}
	NewValue interface{}
}

// ResourceDiff holds all attribute-level differences for a single resource.
type ResourceDiff struct {
	ResourceID string
	Type       string
	Diffs      []AttributeDiff
}

// String returns a human-readable summary of the resource diff.
func (rd ResourceDiff) String() string {
	return fmt.Sprintf("resource %q (%s): %d attribute(s) changed", rd.ResourceID, rd.Type, len(rd.Diffs))
}

// DiffAttributes compares two attribute maps and returns a slice of AttributeDiff.
// Both maps are expected to be map[string]interface{} encoded as json.RawMessage values.
func DiffAttributes(stateAttrs, liveAttrs map[string]interface{}) []AttributeDiff {
	var diffs []AttributeDiff

	keys := unionKeys(stateAttrs, liveAttrs)
	sort.Strings(keys)

	for _, k := range keys {
		sv := stateAttrs[k]
		lv := liveAttrs[k]
		if !attributeEqual(sv, lv) {
			diffs = append(diffs, AttributeDiff{
				Key:      k,
				OldValue: sv,
				NewValue: lv,
			})
		}
	}
	return diffs
}

// unionKeys returns the union of keys from both maps.
func unionKeys(a, b map[string]interface{}) []string {
	seen := make(map[string]struct{}, len(a)+len(b))
	for k := range a {
		seen[k] = struct{}{}
	}
	for k := range b {
		seen[k] = struct{}{}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	return keys
}

// attributeEqual compares two attribute values for equality using JSON normalisation.
func attributeEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	aj, err1 := json.Marshal(a)
	bj, err2 := json.Marshal(b)
	if err1 != nil || err2 != nil {
		return false
	}
	return string(aj) == string(bj)
}
