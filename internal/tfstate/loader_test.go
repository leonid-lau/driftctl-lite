package tfstate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func writeTempState(t *testing.T, dir string, name string, resources []Resource) string {
	t.Helper()
	state := State{
		Version:   4,
		Resources: resources,
	}
	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("marshal state: %v", err)
	}
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatalf("write state file: %v", err)
	}
	return path
}

func TestLoadAll_SingleFile(t *testing.T) {
	dir := t.TempDir()
	writeTempState(t, dir, "terraform.tfstate", []Resource{
		{Type: "aws_s3_bucket", Name: "my_bucket", Instances: []Instance{{Attributes: map[string]interface{}{"id": "bucket-1"}}}},
	})

	resources, err := LoadAll(LoadOptions{SearchDirs: []string{dir}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 1 {
		t.Errorf("expected 1 resource, got %d", len(resources))
	}
}

func TestLoadAll_MultipleFiles(t *testing.T) {
	dir := t.TempDir()
	writeTempState(t, dir, "a.tfstate", []Resource{
		{Type: "aws_s3_bucket", Name: "bucket_a"},
	})
	writeTempState(t, dir, "b.tfstate", []Resource{
		{Type: "aws_instance", Name: "web"},
		{Type: "aws_instance", Name: "db"},
	})

	resources, err := LoadAll(LoadOptions{SearchDirs: []string{dir}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 3 {
		t.Errorf("expected 3 resources, got %d", len(resources))
	}
}

func TestLoadAll_NoFiles(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadAll(LoadOptions{SearchDirs: []string{dir}})
	if err == nil {
		t.Fatal("expected error when no .tfstate files found")
	}
}

func TestLoadAll_Recursive(t *testing.T) {
	root := t.TempDir()
	sub := filepath.Join(root, "module")
	if err := os.Mkdir(sub, 0750); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	writeTempState(t, root, "root.tfstate", []Resource{{Type: "aws_vpc", Name: "main"}})
	writeTempState(t, sub, "module.tfstate", []Resource{{Type: "aws_subnet", Name: "private"}})

	resources, err := LoadAll(LoadOptions{SearchDirs: []string{root}, Recursive: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 2 {
		t.Errorf("expected 2 resources, got %d", len(resources))
	}
}

func TestDefaultLoadOptions(t *testing.T) {
	opts := DefaultLoadOptions()
	if len(opts.SearchDirs) != 1 || opts.SearchDirs[0] != "." {
		t.Errorf("unexpected default SearchDirs: %v", opts.SearchDirs)
	}
	if opts.Recursive {
		t.Error("expected Recursive to be false by default")
	}
}
