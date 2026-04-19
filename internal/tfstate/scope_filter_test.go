package tfstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func scopeResource(id, scope string) Resource {
	return Resource{
		Type: "aws_s3_bucket",
		Name: id,
		Attributes: map[string]string{"scope": scope},
	}
}

func TestFilterByScope_Match(t *testing.T) {
	resources := []Resource{
		scopeResource("a", "global"),
		scopeResource("b", "regional"),
		scopeResource("c", "global"),
	}
	got := FilterByScope(resources, "global", DefaultScopeFilterOptions())
	assert.Len(t, got, 2)
}

func TestFilterByScope_EmptyScope_ReturnsAll(t *testing.T) {
	resources := []Resource{scopeResource("a", "global"), scopeResource("b", "regional")}
	got := FilterByScope(resources, "", DefaultScopeFilterOptions())
	assert.Len(t, got, 2)
}

func TestFilterByScope_CaseInsensitive(t *testing.T) {
	resources := []Resource{scopeResource("a", "Global"), scopeResource("b", "regional")}
	got := FilterByScope(resources, "global", DefaultScopeFilterOptions())
	assert.Len(t, got, 1)
}

func TestFilterByScope_CaseSensitive_NoMatch(t *testing.T) {
	opts := ScopeFilterOptions{CaseSensitive: true}
	resources := []Resource{scopeResource("a", "Global")}
	got := FilterByScope(resources, "global", opts)
	assert.Empty(t, got)
}

func TestFilterByScopes_ORSemantics(t *testing.T) {
	resources := []Resource{
		scopeResource("a", "global"),
		scopeResource("b", "regional"),
		scopeResource("c", "local"),
	}
	got := FilterByScopes(resources, []string{"global", "local"}, DefaultScopeFilterOptions())
	assert.Len(t, got, 2)
}

func TestBuildScopeIndex_Lookup(t *testing.T) {
	resources := []Resource{scopeResource("a", "global"), scopeResource("b", "regional")}
	idx := BuildScopeIndex(resources)
	assert.Len(t, idx.Lookup("global"), 1)
	assert.Len(t, idx.Lookup("GLOBAL"), 1)
	assert.Empty(t, idx.Lookup("zonal"))
}

func TestBuildScopeIndex_Scopes(t *testing.T) {
	resources := []Resource{scopeResource("a", "global"), scopeResource("b", "regional")}
	idx := BuildScopeIndex(resources)
	assert.ElementsMatch(t, []string{"global", "regional"}, idx.Scopes())
}
