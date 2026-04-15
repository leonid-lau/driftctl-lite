package tfstate_test

import (
	"testing"

	"github.com/example/driftctl-lite/internal/tfstate"
)

const sampleState = `{
  "version": 4,
  "resources": [
    {
      "type": "aws_s3_bucket",
      "name": "my_bucket",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "instances": [
        {
          "attributes": {
            "id": "my-bucket-123",
            "region": "us-east-1"
          }
        }
      ]
    }
  ]
}`

func TestParse_ValidState(t *testing.T) {
	state, err := tfstate.Parse([]byte(sampleState))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Version != 4 {
		t.Errorf("expected version 4, got %d", state.Version)
	}
	if len(state.Resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(state.Resources))
	}
	r := state.Resources[0]
	if r.Type != "aws_s3_bucket" {
		t.Errorf("expected type aws_s3_bucket, got %s", r.Type)
	}
	if r.Name != "my_bucket" {
		t.Errorf("expected name my_bucket, got %s", r.Name)
	}
	if r.Attributes["id"] != "my-bucket-123" {
		t.Errorf("expected id my-bucket-123, got %v", r.Attributes["id"])
	}
}

func TestParse_EmptyResources(t *testing.T) {
	state, err := tfstate.Parse([]byte(`{"version":4,"resources":[]}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(state.Resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(state.Resources))
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	_, err := tfstate.Parse([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}
