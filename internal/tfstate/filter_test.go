package tfstate_test

import (
	"testing"

	"github.com/example/driftctl-lite/internal/tfstate"
)

func makeResources() []tfstate.Resource {
	return []tfstate.Resource{
		{Type: "aws_s3_bucket", Name: "bucket1", Attributes: map[string]interface{}{"id": "b1"}},
		{Type: "aws_s3_bucket", Name: "bucket2", Attributes: map[string]interface{}{"id": "b2"}},
		{Type: "aws_instance", Name: "web", Attributes: map[string]interface{}{"id": "i-abc"}},
	}
}

func TestFilterByType(t *testing.T) {
	resources := makeResources()
	buckets := tfstate.FilterByType(resources, "aws_s3_bucket")
	if len(buckets) != 2 {
		t.Errorf("expected 2 buckets, got %d", len(buckets))
	}
	instances := tfstate.FilterByType(resources, "aws_instance")
	if len(instances) != 1 {
		t.Errorf("expected 1 instance, got %d", len(instances))
	}
	none := tfstate.FilterByType(resources, "aws_lambda_function")
	if len(none) != 0 {
		t.Errorf("expected 0 results, got %d", len(none))
	}
}

func TestIndexByID_WithID(t *testing.T) {
	resources := makeResources()
	index := tfstate.IndexByID(resources)
	if _, ok := index["b1"]; !ok {
		t.Error("expected key b1 in index")
	}
	if _, ok := index["i-abc"]; !ok {
		t.Error("expected key i-abc in index")
	}
}

func TestIndexByID_FallbackKey(t *testing.T) {
	resources := []tfstate.Resource{
		{Type: "aws_s3_bucket", Name: "orphan", Attributes: map[string]interface{}{}},
	}
	index := tfstate.IndexByID(resources)
	if _, ok := index["aws_s3_bucket.orphan"]; !ok {
		t.Error("expected fallback key aws_s3_bucket.orphan")
	}
}
