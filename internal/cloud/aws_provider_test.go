package cloud_test

import (
	"context"
	"testing"

	"driftctl-lite/internal/cloud"
)

func TestMockProvider_FetchResource_Found(t *testing.T) {
	provider := cloud.NewMockProvider(map[string]cloud.ResourceAttributes{
		"aws_s3_bucket/my-bucket": {
			"bucket": "my-bucket",
			"region": "us-east-1",
		},
	})

	attrs, err := provider.FetchResource(context.Background(), "aws_s3_bucket", "my-bucket")
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if attrs["bucket"] != "my-bucket" {
		t.Errorf("expected bucket=my-bucket, got %q", attrs["bucket"])
	}
	if attrs["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %q", attrs["region"])
	}
}

func TestMockProvider_FetchResource_NotFound(t *testing.T) {
	provider := cloud.NewMockProvider(nil)

	_, err := provider.FetchResource(context.Background(), "aws_s3_bucket", "missing-bucket")
	if err == nil {
		t.Fatal("expected an error for missing resource, got nil")
	}
}

func TestMockProvider_FetchResource_MultipleTypes(t *testing.T) {
	provider := cloud.NewMockProvider(map[string]cloud.ResourceAttributes{
		"aws_s3_bucket/bucket-a": {"bucket": "bucket-a"},
		"aws_s3_bucket/bucket-b": {"bucket": "bucket-b"},
	})

	for _, id := range []string{"bucket-a", "bucket-b"} {
		attrs, err := provider.FetchResource(context.Background(), "aws_s3_bucket", id)
		if err != nil {
			t.Errorf("unexpected error for %q: %v", id, err)
		}
		if attrs["bucket"] != id {
			t.Errorf("expected bucket=%s, got %q", id, attrs["bucket"])
		}
	}
}
