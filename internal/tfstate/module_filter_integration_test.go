package tfstate_test

import (
	"testing"

	"github.com/example/driftctl-lite/internal/tfstate"
)

func TestModuleFilterAndIndex_RoundTrip(t *testing.T) {
	resources := []tfstate.Resource{
		{Type: "aws_instance", Attributes: map[string]interface{}{"module": "vpc"}},
		{Type: "aws_subnet", Attributes: map[string]interface{}{"module": "vpc"}},
		{Type: "aws_ecs_service", Attributes: map[string]interface{}{"module": "ecs"}},
		{Type: "aws_rds_instance", Attributes: map[string]interface{}{"module": "rds"}},
	}

	// Filter then index should be consistent.
	filtered := tfstate.FilterByModule(resources, "vpc", tfstate.DefaultModuleFilterOptions())
	if len(filtered) != 2 {
		t.Fatalf("filter: expected 2, got %d", len(filtered))
	}

	idx := tfstate.BuildModuleIndex(resources)
	looked := idx.Lookup("vpc")
	if len(looked) != len(filtered) {
		t.Fatalf("index lookup count %d != filter count %d", len(looked), len(filtered))
	}

	// OR filter across multiple modules.
	multi := tfstate.FilterByModules(resources, []string{"ecs", "rds"}, tfstate.DefaultModuleFilterOptions())
	if len(multi) != 2 {
		t.Fatalf("multi filter: expected 2, got %d", len(multi))
	}
}
