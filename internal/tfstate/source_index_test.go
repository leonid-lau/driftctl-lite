package tfstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func idxSrcResource(id, source string) Resource {
	return Resource{ID: id, Type: "aws_s3_bucket", Source: source}
}

func TestBuildSourceIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxSrcResource("r1", "terraform"),
		idxSrcResource("r2", "manual"),
		idxSrcResource("r3", "terraform"),
	}
	idx := BuildSourceIndex(resources)
	got := idx.Lookup("terraform")
	assert.Len(t, got, 2)
}

func TestBuildSourceIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxSrcResource("r1", "Terraform")}
	idx := BuildSourceIndex(resources)
	got := idx.Lookup("TERRAFORM")
	assert.Len(t, got, 1)
}

func TestBuildSourceIndex_LookupMissing(t *testing.T) {
	idx := BuildSourceIndex([]Resource{idxSrcResource("r1", "manual")})
	got := idx.Lookup("import")
	assert.Nil(t, got)
}

func TestBuildSourceIndex_Sources(t *testing.T) {
	resources := []Resource{
		idxSrcResource("r1", "terraform"),
		idxSrcResource("r2", "manual"),
	}
	idx := BuildSourceIndex(resources)
	sources := idx.Sources()
	assert.Len(t, sources, 2)
	assert.ElementsMatch(t, []string{"terraform", "manual"}, sources)
}

func TestBuildSourceIndex_EmptyInput(t *testing.T) {
	idx := BuildSourceIndex(nil)
	assert.Empty(t, idx.Sources())
}
