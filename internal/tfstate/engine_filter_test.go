package tfstate

import (
	"testing"
)

func engineResource(engine string) Resource {
	attrs := map[string]interface{}{}
	if engine != "" {
		attrs["engine"] = engine
	}
	return Resource{
		Type:       "aws_db_instance",
		Name:       "test",
		Attributes: attrs,
	}
}

func TestFilterByEngine_Match(t *testing.T) {
	resources := []Resource{
		engineResource("mysql"),
		engineResource("postgres"),
		engineResource("mysql"),
	}
	got := FilterByEngine(resources, "mysql", DefaultEngineFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByEngine_EmptyEngine_ReturnsAll(t *testing.T) {
	resources := []Resource{
		engineResource("mysql"),
		engineResource("postgres"),
	}
	got := FilterByEngine(resources, "", DefaultEngineFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByEngine_CaseInsensitive(t *testing.T) {
	resources := []Resource{
		engineResource("MySQL"),
		engineResource("POSTGRES"),
	}
	got := FilterByEngine(resources, "mysql", DefaultEngineFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByEngine_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{
		engineResource("MySQL"),
	}
	opts := EngineFilterOptions{CaseSensitive: true}
	got := FilterByEngine(resources, "mysql", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByEngines_ORSemantics(t *testing.T) {
	resources := []Resource{
		engineResource("mysql"),
		engineResource("postgres"),
		engineResource("redis"),
	}
	got := FilterByEngines(resources, []string{"mysql", "redis"}, DefaultEngineFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByEngines_EmptyList_ReturnsAll(t *testing.T) {
	resources := []Resource{
		engineResource("mysql"),
		engineResource("postgres"),
	}
	got := FilterByEngines(resources, []string{}, DefaultEngineFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
