package tfstate_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"driftctl-lite/internal/tfstate"
)

func writePipelineState(t *testing.T, dir string, name string, resources []tfstate.Resource) {
	t.Helper()
	s := &tfstate.State{Version: 4, Resources: resources}
	data, _ := json.Marshal(s)
	_ = os.WriteFile(filepath.Join(dir, name), data, 0644)
}

func TestPipeline_Run_Success(t *testing.T) {
	dir := t.TempDir()
	writePipelineState(t, dir, "terraform.tfstate", []tfstate.Resource{
		{Type: "aws_s3_bucket", Name: "main", Instances: []tfstate.Instance{{Attributes: map[string]interface{}{"id": "bucket-1"}}}},
	})

	p := tfstate.NewPipeline(tfstate.DefaultPipelineOptions())
	state, err := p.Run(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(state.Resources) != 1 {
		t.Errorf("expected 1 resource, got %d", len(state.Resources))
	}
}

func TestPipeline_Run_NoFiles(t *testing.T) {
	dir := t.TempDir()
	p := tfstate.NewPipeline(tfstate.DefaultPipelineOptions())
	state, err := p.Run(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(state.Resources) != 0 {
		t.Errorf("expected 0 resources, got %d", len(state.Resources))
	}
}

func TestPipeline_Run_ValidationFails(t *testing.T) {
	dir := t.TempDir()
	// Write a state with an invalid version to trigger validation failure
	s := &tfstate.State{Version: 0, Resources: []tfstate.Resource{
		{Type: "", Name: "bad", Instances: []tfstate.Instance{}},
	}}
	data, _ := json.Marshal(s)
	_ = os.WriteFile(filepath.Join(dir, "terraform.tfstate"), data, 0644)

	opts := tfstate.DefaultPipelineOptions()
	opts.Validate = true
	p := tfstate.NewPipeline(opts)
	_, err := p.Run(dir)
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestPipeline_DefaultOptions(t *testing.T) {
	opts := tfstate.DefaultPipelineOptions()
	if !opts.Validate {
		t.Error("expected Validate to be true by default")
	}
	if opts.UseCache {
		t.Error("expected UseCache to be false by default")
	}
	if opts.CacheDir == "" {
		t.Error("expected non-empty CacheDir")
	}
}
