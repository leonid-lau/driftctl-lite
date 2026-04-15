package cloud

import (
	"context"
	"fmt"
)

// MockProvider is an in-memory Provider implementation for use in tests.
type MockProvider struct {
	// Resources maps "resourceType/resourceID" to its attributes.
	// A nil value signals that the resource does not exist in the cloud.
	Resources map[string]ResourceAttributes
}

// NewMockProvider initialises a MockProvider with the given resource map.
func NewMockProvider(resources map[string]ResourceAttributes) *MockProvider {
	if resources == nil {
		resources = make(map[string]ResourceAttributes)
	}
	return &MockProvider{Resources: resources}
}

// FetchResource returns the pre-configured attributes for the given resource,
// or an error if the resource is not present in the mock map.
func (m *MockProvider) FetchResource(_ context.Context, resourceType, resourceID string) (ResourceAttributes, error) {
	key := resourceType + "/" + resourceID
	attrs, ok := m.Resources[key]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", key)
	}
	return attrs, nil
}
