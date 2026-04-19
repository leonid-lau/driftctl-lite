package tfstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func serviceResource(id, service string) Resource {
	return Resource{
		Type: "aws_instance",
		Name: id,
		Attributes: map[string]interface{}{"service": service},
	}
}

func TestFilterByService_Match(t *testing.T) {
	resources := []Resource{
		serviceResource("r1", "payments"),
		serviceResource("r2", "auth"),
		serviceResource("r3", "payments"),
	}
	got := FilterByService(resources, "payments", DefaultServiceFilterOptions())
	assert.Len(t, got, 2)
}

func TestFilterByService_EmptyService_ReturnsAll(t *testing.T) {
	resources := []Resource{
		serviceResource("r1", "payments"),
		serviceResource("r2", "auth"),
	}
	got := FilterByService(resources, "", DefaultServiceFilterOptions())
	assert.Len(t, got, 2)
}

func TestFilterByService_CaseInsensitive(t *testing.T) {
	resources := []Resource{serviceResource("r1", "Payments")}
	got := FilterByService(resources, "payments", DefaultServiceFilterOptions())
	assert.Len(t, got, 1)
}

func TestFilterByService_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{serviceResource("r1", "Payments")}
	opts := ServiceFilterOptions{CaseSensitive: true}
	got := FilterByService(resources, "payments", opts)
	assert.Len(t, got, 0)
}

func TestFilterByServices_ORSemantics(t *testing.T) {
	resources := []Resource{
		serviceResource("r1", "payments"),
		serviceResource("r2", "auth"),
		serviceResource("r3", "storage"),
	}
	got := FilterByServices(resources, []string{"payments", "auth"}, DefaultServiceFilterOptions())
	assert.Len(t, got, 2)
}

func TestBuildServiceIndex_Lookup(t *testing.T) {
	resources := []Resource{
		serviceResource("r1", "payments"),
		serviceResource("r2", "auth"),
		serviceResource("r3", "Payments"),
	}
	idx := BuildServiceIndex(resources)
	assert.Len(t, idx.Lookup("payments"), 2)
	assert.Len(t, idx.Lookup("auth"), 1)
	assert.Len(t, idx.Lookup("unknown"), 0)
}

func TestBuildServiceIndex_Services(t *testing.T) {
	resources := []Resource{
		serviceResource("r1", "payments"),
		serviceResource("r2", "auth"),
	}
	idx := BuildServiceIndex(resources)
	svcs := idx.Services()
	assert.Len(t, svcs, 2)
}
