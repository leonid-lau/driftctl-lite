package tfstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func srcResource(id, source string) Resource {
	return Resource{ID: id, Type: "aws_instance", Source: source}
}

func TestFilterBySource_Match(t *testing.T) {
	resources := []Resource{
		srcResource("a", "terraform"),
		srcResource("b", "manual"),
		srcResource("c", "terraform"),
	}
	got := FilterBySource(resources, "terraform", DefaultSourceFilterOptions())
	assert.Len(t, got, 2)
}

func TestFilterBySource_EmptySource_ReturnsAll(t *testing.T) {
	resources := []Resource{srcResource("a", "terraform"), srcResource("b", "manual")}
	got := FilterBySource(resources, "", DefaultSourceFilterOptions())
	assert.Len(t, got, 2)
}

func TestFilterBySource_CaseInsensitive(t *testing.T) {
	resources := []Resource{srcResource("a", "Terraform"), srcResource("b", "Manual")}
	got := FilterBySource(resources, "terraform", DefaultSourceFilterOptions())
	assert.Len(t, got, 1)
	assert.Equal(t, "a", got[0].ID)
}

func TestFilterBySource_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{srcResource("a", "Terraform")}
	opts := SourceFilterOptions{CaseSensitive: true}
	got := FilterBySource(resources, "terraform", opts)
	assert.Empty(t, got)
}

func TestFilterBySources_ORSemantics(t *testing.T) {
	resources := []Resource{
		srcResource("a", "terraform"),
		srcResource("b", "manual"),
		srcResource("c", "import"),
	}
	got := FilterBySources(resources, []string{"terraform", "manual"}, DefaultSourceFilterOptions())
	assert.Len(t, got, 2)
}

func TestFilterBySources_EmptyList_ReturnsAll(t *testing.T) {
	resources := []Resource{srcResource("a", "terraform")}
	got := FilterBySources(resources, nil, DefaultSourceFilterOptions())
	assert.Len(t, got, 1)
}
