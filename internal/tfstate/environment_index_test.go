package tfstate

import (
	"sort"
	"testing"
)

func idxEnvResource(id, env string) Resource {
	return Resource{ID: id, Type: "aws_instance", Attributes: map[string]interface{}{"environment": env}}
}

func TestBuildEnvironmentIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxEnvResource("r1", "production"),
		idxEnvResource("r2", "staging"),
		idxEnvResource("r3", "production"),
	}
	idx := BuildEnvironmentIndex(resources)
	got := idx.Lookup("production")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildEnvironmentIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxEnvResource("r1", "Production")}
	idx := BuildEnvironmentIndex(resources)
	got := idx.Lookup("production")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestBuildEnvironmentIndex_LookupMissing(t *testing.T) {
	idx := BuildEnvironmentIndex([]Resource{})
	got := idx.Lookup("production")
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestBuildEnvironmentIndex_Environments(t *testing.T) {
	resources := []Resource{
		idxEnvResource("r1", "production"),
		idxEnvResource("r2", "staging"),
		idxEnvResource("r3", "dev"),
	}
	idx := BuildEnvironmentIndex(resources)
	envs := idx.Environments()
	sort.Strings(envs)
	if len(envs) != 3 {
		t.Fatalf("expected 3 environments, got %d", len(envs))
	}
	if envs[0] != "dev" || envs[1] != "production" || envs[2] != "staging" {
		t.Fatalf("unexpected environments: %v", envs)
	}
}

func TestBuildEnvironmentIndex_EmptyInput(t *testing.T) {
	idx := BuildEnvironmentIndex(nil)
	if len(idx.Environments()) != 0 {
		t.Fatal("expected empty index")
	}
}
