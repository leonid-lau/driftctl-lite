package tfstate

import (
	"testing"
)

func makeResource(rtype, name, id string) Resource {
	return Resource{
		Type: rtype,
		Name: name,
		Instances: []ResourceInstance{
			{Attributes: map[string]interface{}{"id": id}},
		},
	}
}

func TestNewSnapshot_Basic(t *testing.T) {
	res := []Resource{
		makeResource("aws_s3_bucket", "my_bucket", "bucket-1"),
	}
	s, err := NewSnapshot(res, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
	if s.ResourceCount() != 1 {
		t.Errorf("expected 1 resource, got %d", s.ResourceCount())
	}
}

func TestNewSnapshot_DeterministicChecksum(t *testing.T) {
	res1 := []Resource{
		makeResource("aws_s3_bucket", "a", "id-a"),
		makeResource("aws_instance", "b", "id-b"),
	}
	res2 := []Resource{
		makeResource("aws_instance", "b", "id-b"),
		makeResource("aws_s3_bucket", "a", "id-a"),
	}
	s1, _ := NewSnapshot(res1, nil)
	s2, _ := NewSnapshot(res2, nil)
	if s1.Checksum != s2.Checksum {
		t.Errorf("checksums should be equal regardless of input order: %s vs %s", s1.Checksum, s2.Checksum)
	}
}

func TestSnapshot_Equal(t *testing.T) {
	res := []Resource{makeResource("aws_s3_bucket", "b", "id-1")}
	s1, _ := NewSnapshot(res, nil)
	s2, _ := NewSnapshot(res, nil)
	if !s1.Equal(s2) {
		t.Error("expected snapshots to be equal")
	}
}

func TestSnapshot_NotEqual(t *testing.T) {
	s1, _ := NewSnapshot([]Resource{makeResource("aws_s3_bucket", "a", "id-1")}, nil)
	s2, _ := NewSnapshot([]Resource{makeResource("aws_s3_bucket", "a", "id-2")}, nil)
	if s1.Equal(s2) {
		t.Error("expected snapshots to differ")
	}
}

func TestSnapshot_NilEqual(t *testing.T) {
	var s1 *Snapshot
	var s2 *Snapshot
	if !s1.Equal(s2) {
		t.Error("two nil snapshots should be equal")
	}
	s3, _ := NewSnapshot(nil, nil)
	if s1.Equal(s3) {
		t.Error("nil and non-nil snapshot should not be equal")
	}
}

func TestSnapshot_WithMeta(t *testing.T) {
	meta := map[string]string{"source": "terraform.tfstate", "workspace": "prod"}
	s, err := NewSnapshot(nil, meta)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Meta["workspace"] != "prod" {
		t.Errorf("expected meta workspace=prod, got %s", s.Meta["workspace"])
	}
}
