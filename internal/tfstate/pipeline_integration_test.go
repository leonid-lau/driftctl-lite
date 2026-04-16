package tfstate_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"driftctl-lite/internal/tfstate"
)

func TestPipeline_FullCycle_MultipleFiles(t *testing.T) {
	dir := t.TempDir()

	resources := [][]tfstate.Resource{
		{{Type: "aws_s3_bucket", Name: "a", Instances: []tfstate.Instance{{Attributes: map[string]interface{}{"id": "bucket-a"}}}}},
		{{Type: "aws_instance", Name: "b", Instances: []tfstate.Instance{{Attributes: map[string]interface{}{"id": "i-123"}}}}},
	}

	for i, res := range resources {
		s := &tfstate.State{Version: 4, Resources: res}
		data, _ := json.Marshal(s)
		name := filepath.Join(dir, filepath.Join("mod"+string(rune('a'+i)), "terraform.tfstate"))
		_ = os.MkdirAll(filepath.Dir(name), 0755)
		_ = os.WriteFile(name, data, 0644)
	}

	opts := tfstate.DefaultPipelineOptions()
	opts.LoadOptions = tfstate.LoadOptions{Recursive: true}
	p := tfstate.NewPipeline(opts)

	state, err := p.Run(dir)
	if err != nil {
		t.Fatalf("pipeline error: %v", err)
	}
	if len(state.Resources) != 2 {
		t.Errorf("expected 2 resources, got %d", len(state.Resources))
	}

	byType := tfstate.FilterByType(state.Resources, "aws_s3_bucket")
	if len(byType) != 1 {
		t.Errorf("expected 1 s3 bucket, got %d", len(byType))
	}
}
