package tfstate

import (
	"testing"
)

func buildState(resources []Resource) State {
	return State{Resources: resources}
}

func TestMergeStates_Empty(t *testing.T) {
	result, err := MergeStates(nil, DefaultMergeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 resources, got %d", len(result))
	}
}

func TestMergeStates_NoConflict(t *testing.T) {
	s1 := buildState([]Resource{
		{Type: "aws_s3_bucket", Name: "a", Instances: []Instance{{Attributes: map[string]interface{}{"id": "bucket-a"}}}},
	})
	s2 := buildState([]Resource{
		{Type: "aws_s3_bucket", Name: "b", Instances: []Instance{{Attributes: map[string]interface{}{"id": "bucket-b"}}}},
	})

	result, err := MergeStates([]State{s1, s2}, DefaultMergeOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 resources, got %d", len(result))
	}
}

func TestMergeStates_ConflictError(t *testing.T) {
	res := Resource{
		Type: "aws_s3_bucket", Name: "dup",
		Instances: []Instance{{Attributes: map[string]interface{}{"id": "same-id"}}},
	}
	s1 := buildState([]Resource{res})
	s2 := buildState([]Resource{res})

	_, err := MergeStates([]State{s1, s2}, DefaultMergeOptions())
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
}

func TestMergeStates_ConflictLastWins(t *testing.T) {
	res1 := Resource{
		Type: "aws_s3_bucket", Name: "dup",
		Instances: []Instance{{Attributes: map[string]interface{}{"id": "same-id", "region": "us-east-1"}}},
	}
	res2 := Resource{
		Type: "aws_s3_bucket", Name: "dup",
		Instances: []Instance{{Attributes: map[string]interface{}{"id": "same-id", "region": "eu-west-1"}}},
	}

	opts := MergeOptions{OnConflict: "last-wins"}
	result, err := MergeStates([]State{buildState([]Resource{res1}), buildState([]Resource{res2})}, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 resource after last-wins, got %d", len(result))
	}
	got := result[0].Instances[0].Attributes["region"]
	if got != "eu-west-1" {
		t.Errorf("expected last-wins region eu-west-1, got %v", got)
	}
}

func TestMergeStates_UnknownStrategy(t *testing.T) {
	res := Resource{Type: "aws_s3_bucket", Name: "dup"}
	s := buildState([]Resource{res, res})
	_, err := MergeStates([]State{s}, MergeOptions{OnConflict: "unknown"})
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}
