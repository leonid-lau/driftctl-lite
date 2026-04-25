package tfstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func profileResource(id, profile string) Resource {
	attrs := map[string]string{"id": id}
	if profile != "" {
		attrs["profile"] = profile
	}
	return Resource{Type: "aws_instance", ID: id, Attributes: attrs}
}

func TestFilterByProfile_Match(t *testing.T) {
	resources := []Resource{
		profileResource("r1", "default"),
		profileResource("r2", "production"),
		profileResource("r3", "staging"),
	}
	got := FilterByProfile(resources, "production", DefaultProfileFilterOptions())
	assert.Len(t, got, 1)
	assert.Equal(t, "r2", got[0].ID)
}

func TestFilterByProfile_EmptyProfile_ReturnsAll(t *testing.T) {
	resources := []Resource{
		profileResource("r1", "default"),
		profileResource("r2", "production"),
	}
	got := FilterByProfile(resources, "", DefaultProfileFilterOptions())
	assert.Len(t, got, 2)
}

func TestFilterByProfile_CaseInsensitive(t *testing.T) {
	resources := []Resource{
		profileResource("r1", "Default"),
		profileResource("r2", "production"),
	}
	got := FilterByProfile(resources, "default", DefaultProfileFilterOptions())
	assert.Len(t, got, 1)
	assert.Equal(t, "r1", got[0].ID)
}

func TestFilterByProfile_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{
		profileResource("r1", "Default"),
	}
	opts := ProfileFilterOptions{CaseSensitive: true}
	got := FilterByProfile(resources, "default", opts)
	assert.Empty(t, got)
}

func TestFilterByProfiles_ORSemantics(t *testing.T) {
	resources := []Resource{
		profileResource("r1", "default"),
		profileResource("r2", "production"),
		profileResource("r3", "staging"),
	}
	got := FilterByProfiles(resources, []string{"default", "staging"}, DefaultProfileFilterOptions())
	assert.Len(t, got, 2)
}

func TestBuildProfileIndex_Lookup(t *testing.T) {
	resources := []Resource{
		profileResource("r1", "default"),
		profileResource("r2", "production"),
		profileResource("r3", "default"),
	}
	idx := BuildProfileIndex(resources)
	got := idx.Lookup("default")
	assert.Len(t, got, 2)
}

func TestBuildProfileIndex_LookupMissing(t *testing.T) {
	idx := BuildProfileIndex([]Resource{profileResource("r1", "default")})
	assert.Nil(t, idx.Lookup("nonexistent"))
}

func TestBuildProfileIndex_Profiles(t *testing.T) {
	resources := []Resource{
		profileResource("r1", "default"),
		profileResource("r2", "production"),
	}
	idx := BuildProfileIndex(resources)
	profiles := idx.Profiles()
	assert.Len(t, profiles, 2)
}
