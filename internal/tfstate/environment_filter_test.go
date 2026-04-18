package tfstate

import (
	"testing"
)

func envResource(id, env string) Resource {
	attrs := map[string]interface{}{}
	if env != "" {
		attrs["environment"] = env
	}
	return Resource{ID: id, Type: "aws_instance", Attributes: attrs}
}

func TestFilterByEnvironment_Match(t *testing.T) {
	resources := []Resource{
		envResource("r1", "production"),
		envResource("r2", "staging"),
		envResource("r3", "production"),
	}
	opts := DefaultEnvironmentFilterOptions()
	got := FilterByEnvironment(resources, "production", opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByEnvironment_EmptyEnv_ReturnsAll(t *testing.T) {
	resources := []Resource{envResource("r1", "prod"), envResource("r2", "staging")}
	opts := DefaultEnvironmentFilterOptions()
	got := FilterByEnvironment(resources, "", opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByEnvironment_CaseInsensitive(t *testing.T) {
	resources := []Resource{envResource("r1", "Production"), envResource("r2", "staging")}
	opts := DefaultEnvironmentFilterOptions()
	got := FilterByEnvironment(resources, "production", opts)
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByEnvironment_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{envResource("r1", "Production")}
	opts := EnvironmentFilterOptions{CaseSensitive: true}
	got := FilterByEnvironment(resources, "production", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByEnvironments_ORSemantics(t *testing.T) {
	resources := []Resource{
		envResource("r1", "production"),
		envResource("r2", "staging"),
		envResource("r3", "dev"),
	}
	opts := DefaultEnvironmentFilterOptions()
	got := FilterByEnvironments(resources, []string{"production", "dev"}, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByEnvironments_Empty_ReturnsAll(t *testing.T) {
	resources := []Resource{envResource("r1", "prod"), envResource("r2", "staging")}
	opts := DefaultEnvironmentFilterOptions()
	got := FilterByEnvironments(resources, nil, opts)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
