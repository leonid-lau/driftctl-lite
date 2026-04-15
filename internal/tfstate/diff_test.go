package tfstate

import (
	"testing"
)

func TestDiffAttributes_NoDiff(t *testing.T) {
	state := map[string]interface{}{"id": "abc", "region": "us-east-1"}
	live := map[string]interface{}{"id": "abc", "region": "us-east-1"}

	diffs := DiffAttributes(state, live)
	if len(diffs) != 0 {
		t.Errorf("expected 0 diffs, got %d", len(diffs))
	}
}

func TestDiffAttributes_ModifiedValue(t *testing.T) {
	state := map[string]interface{}{"id": "abc", "region": "us-east-1"}
	live := map[string]interface{}{"id": "abc", "region": "eu-west-1"}

	diffs := DiffAttributes(state, live)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Key != "region" {
		t.Errorf("expected diff on 'region', got %q", diffs[0].Key)
	}
	if diffs[0].OldValue != "us-east-1" {
		t.Errorf("unexpected OldValue: %v", diffs[0].OldValue)
	}
	if diffs[0].NewValue != "eu-west-1" {
		t.Errorf("unexpected NewValue: %v", diffs[0].NewValue)
	}
}

func TestDiffAttributes_MissingInLive(t *testing.T) {
	state := map[string]interface{}{"id": "abc", "tags": "env=prod"}
	live := map[string]interface{}{"id": "abc"}

	diffs := DiffAttributes(state, live)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Key != "tags" {
		t.Errorf("expected diff on 'tags', got %q", diffs[0].Key)
	}
	if diffs[0].NewValue != nil {
		t.Errorf("expected nil NewValue, got %v", diffs[0].NewValue)
	}
}

func TestDiffAttributes_ExtraInLive(t *testing.T) {
	state := map[string]interface{}{"id": "abc"}
	live := map[string]interface{}{"id": "abc", "extra": "value"}

	diffs := DiffAttributes(state, live)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Key != "extra" {
		t.Errorf("expected diff on 'extra', got %q", diffs[0].Key)
	}
}

func TestResourceDiff_String(t *testing.T) {
	rd := ResourceDiff{
		ResourceID: "i-123",
		Type:       "aws_instance",
		Diffs:      []AttributeDiff{{Key: "region", OldValue: "us-east-1", NewValue: "eu-west-1"}},
	}
	s := rd.String()
	if s == "" {
		t.Error("expected non-empty string from ResourceDiff.String()")
	}
}
