package tfstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func roleResource(id, role string) Resource {
	return Resource{
		Type: "aws_iam_role",
		ID:   id,
		Attributes: map[string]interface{}{
			"role": role,
		},
	}
}

func TestFilterByRole_Match(t *testing.T) {
	res := []Resource{
		roleResource("r1", "admin"),
		roleResource("r2", "viewer"),
		roleResource("r3", "editor"),
	}
	got := FilterByRole(res, "admin", DefaultRoleFilterOptions())
	assert.Len(t, got, 1)
	assert.Equal(t, "r1", got[0].ID)
}

func TestFilterByRole_EmptyRole_ReturnsAll(t *testing.T) {
	res := []Resource{
		roleResource("r1", "admin"),
		roleResource("r2", "viewer"),
	}
	got := FilterByRole(res, "", DefaultRoleFilterOptions())
	assert.Len(t, got, 2)
}

func TestFilterByRole_CaseInsensitive(t *testing.T) {
	res := []Resource{
		roleResource("r1", "Admin"),
		roleResource("r2", "viewer"),
	}
	got := FilterByRole(res, "admin", DefaultRoleFilterOptions())
	assert.Len(t, got, 1)
	assert.Equal(t, "r1", got[0].ID)
}

func TestFilterByRole_CaseSensitive_NoMatch(t *testing.T) {
	res := []Resource{
		roleResource("r1", "Admin"),
	}
	opts := FilterOptions{CaseSensitive: true}
	got := FilterByRole(res, "admin", opts)
	assert.Empty(t, got)
}

func TestFilterByRoles_ORSemantics(t *testing.T) {
	res := []Resource{
		roleResource("r1", "admin"),
		roleResource("r2", "viewer"),
		roleResource("r3", "editor"),
	}
	got := FilterByRoles(res, []string{"admin", "viewer"}, DefaultRoleFilterOptions())
	assert.Len(t, got, 2)
}

func TestBuildRoleIndex_Lookup(t *testing.T) {
	res := []Resource{
		roleResource("r1", "admin"),
		roleResource("r2", "viewer"),
	}
	idx := BuildRoleIndex(res)
	got := idx.Lookup("admin")
	assert.Len(t, got, 1)
	assert.Equal(t, "r1", got[0].ID)
}

func TestBuildRoleIndex_LookupCaseInsensitive(t *testing.T) {
	res := []Resource{
		roleResource("r1", "ADMIN"),
	}
	idx := BuildRoleIndex(res)
	assert.Len(t, idx.Lookup("admin"), 1)
}

func TestBuildRoleIndex_LookupMissing(t *testing.T) {
	idx := BuildRoleIndex(nil)
	assert.Nil(t, idx.Lookup("admin"))
}

func TestBuildRoleIndex_Roles(t *testing.T) {
	res := []Resource{
		roleResource("r1", "admin"),
		roleResource("r2", "viewer"),
	}
	idx := BuildRoleIndex(res)
	roles := idx.Roles()
	assert.Len(t, roles, 2)
}
