package cloud_test

import (
	"testing"

	"github.com/example/driftctl-lite/internal/cloud"
)

func TestGetProvider_AWS(t *testing.T) {
	p, ok := cloud.GetProvider("aws")
	if !ok {
		t.Fatal("expected aws provider to be registered")
	}
	if p == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestGetProvider_Mock(t *testing.T) {
	p, ok := cloud.GetProvider("mock")
	if !ok {
		t.Fatal("expected mock provider to be registered")
	}
	if p == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestGetProvider_Unknown(t *testing.T) {
	p, ok := cloud.GetProvider("gcp")
	if ok {
		t.Fatal("expected unknown provider to return false")
	}
	if p != nil {
		t.Fatal("expected nil provider for unknown name")
	}
}

func TestGetProvider_ImplementsInterface(t *testing.T) {
	attrs := cloud.ResourceAttributes{
		"id":   "res-1",
		"name": "test",
	}
	mock := cloud.NewMockProvider(map[string]cloud.ResourceAttributes{
		"aws_s3_bucket:res-1": attrs,
	})

	var _ cloud.Provider = mock

	got, err := mock.FetchResource("aws_s3_bucket", "res-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["id"] != "res-1" {
		t.Errorf("expected id=res-1, got %v", got["id"])
	}
}

func TestGetProvider_NotFound_ReturnsNilAttrs(t *testing.T) {
	mock := cloud.NewMockProvider(nil)
	got, err := mock.FetchResource("aws_instance", "i-missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil attributes for missing resource, got %v", got)
	}
}
