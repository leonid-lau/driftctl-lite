package tfstate

import (
	"strings"
	"testing"
)

func TestValidateState_Nil(t *testing.T) {
	err := ValidateState(nil)
	if err == nil {
		t.Fatal("expected error for nil state")
	}
}

func TestValidateState_Valid(t *testing.T) {
	s := &State{
		Version: 4,
		Resources: []Resource{
			{Type: "aws_instance", Name: "web", Attributes: map[string]interface{}{"id": "i-123"}},
			{Type: "aws_s3_bucket", Name: "data", Attributes: map[string]interface{}{"id": "my-bucket"}},
		},
	}
	if err := ValidateState(s); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidateState_BadVersion(t *testing.T) {
	s := &State{
		Version:   0,
		Resources: []Resource{{Type: "aws_instance", Name: "web"}},
	}
	err := ValidateState(s)
	if err == nil {
		t.Fatal("expected error for version 0")
	}
	if !strings.Contains(err.Error(), "unsupported state version") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidateState_EmptyType(t *testing.T) {
	s := &State{
		Version:   4,
		Resources: []Resource{{Type: "", Name: "orphan"}},
	}
	err := ValidateState(s)
	if err == nil {
		t.Fatal("expected error for empty type")
	}
	if !strings.Contains(err.Error(), "empty type") {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestValidateState_DuplicateKey(t *testing.T) {
	s := &State{
		Version: 4,
		Resources: []Resource{
			{Type: "aws_instance", Name: "web"},
			{Type: "aws_instance", Name: "web"},
		},
	}
	err := ValidateState(s)
	if err == nil {
		t.Fatal("expected error for duplicate resource key")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d: %v", len(ve.Issues), ve.Issues)
	}
}

func TestValidateState_MultipleIssues(t *testing.T) {
	s := &State{
		Version: 0,
		Resources: []Resource{
			{Type: "", Name: ""},
		},
	}
	err := ValidateState(s)
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if len(ve.Issues) < 2 {
		t.Errorf("expected at least 2 issues, got %d", len(ve.Issues))
	}
}
