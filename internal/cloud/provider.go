package cloud

// ResourceAttributes is a map of attribute key-value pairs for a cloud resource.
type ResourceAttributes map[string]interface{}

// Provider defines the interface for fetching live cloud resource state.
type Provider interface {
	// FetchResource retrieves the live attributes of a resource by type and ID.
	// Returns (nil, nil) when the resource does not exist in the cloud.
	FetchResource(resourceType, resourceID string) (ResourceAttributes, error)
}

// ProviderRegistry maps provider names to their constructor functions.
var ProviderRegistry = map[string]func() Provider{
	"aws": func() Provider { return NewAWSProvider() },
	"mock": func() Provider { return NewMockProvider(nil) },
}

// GetProvider returns a Provider for the given name, or nil if not found.
func GetProvider(name string) (Provider, bool) {
	constructor, ok := ProviderRegistry[name]
	if !ok {
		return nil, false
	}
	return constructor(), true
}
