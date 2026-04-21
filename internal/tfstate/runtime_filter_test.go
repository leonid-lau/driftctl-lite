package tfstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func runtimeResource(id, runtime string) Resource {
	return Resource{
		Type: "aws_lambda_function",
		Name: id,
		Attributes: map[string]interface{}{"runtime": runtime},
	}
}

func TestFilterByRuntime_Match(t *testing.T) {
	resources := []Resource{
		runtimeResource("fn1", "python3.9"),
		runtimeResource("fn2", "nodejs18.x"),
		runtimeResource("fn3", "python3.9"),
	}
	opts := DefaultRuntimeFilterOptions()
	got := FilterByRuntime(resources, "python3.9", opts)
	assert.Len(t, got, 2)
}

func TestFilterByRuntime_EmptyRuntime_ReturnsAll(t *testing.T) {
	resources := []Resource{
		runtimeResource("fn1", "python3.9"),
		runtimeResource("fn2", "nodejs18.x"),
	}
	opts := DefaultRuntimeFilterOptions()
	got := FilterByRuntime(resources, "", opts)
	assert.Len(t, got, 2)
}

func TestFilterByRuntime_CaseInsensitive(t *testing.T) {
	resources := []Resource{
		runtimeResource("fn1", "Python3.9"),
	}
	opts := DefaultRuntimeFilterOptions()
	got := FilterByRuntime(resources, "python3.9", opts)
	assert.Len(t, got, 1)
}

func TestFilterByRuntime_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{
		runtimeResource("fn1", "Python3.9"),
	}
	opts := FilterOptions{CaseSensitive: true}
	got := FilterByRuntime(resources, "python3.9", opts)
	assert.Empty(t, got)
}

func TestFilterByRuntimes_ORSemantics(t *testing.T) {
	resources := []Resource{
		runtimeResource("fn1", "python3.9"),
		runtimeResource("fn2", "nodejs18.x"),
		runtimeResource("fn3", "go1.x"),
	}
	opts := DefaultRuntimeFilterOptions()
	got := FilterByRuntimes(resources, []string{"python3.9", "go1.x"}, opts)
	assert.Len(t, got, 2)
}

func TestBuildRuntimeIndex_Lookup(t *testing.T) {
	resources := []Resource{
		runtimeResource("fn1", "python3.9"),
		runtimeResource("fn2", "nodejs18.x"),
		runtimeResource("fn3", "python3.9"),
	}
	idx := BuildRuntimeIndex(resources)
	got := LookupByRuntime(idx, "python3.9")
	assert.Len(t, got, 2)
}

func TestBuildRuntimeIndex_LookupMissing(t *testing.T) {
	idx := BuildRuntimeIndex([]Resource{runtimeResource("fn1", "go1.x")})
	got := LookupByRuntime(idx, "ruby2.7")
	assert.Empty(t, got)
}

func TestBuildRuntimeIndex_Runtimes(t *testing.T) {
	resources := []Resource{
		runtimeResource("fn1", "python3.9"),
		runtimeResource("fn2", "nodejs18.x"),
	}
	idx := BuildRuntimeIndex(resources)
	keys := Runtimes(idx)
	assert.Len(t, keys, 2)
}
