package cloud_test

import (
	"testing"

	"github.com/example/driftctl-lite/internal/cloud"
)

// TestProviderRegistry_AllEntriesReturnNonNil verifies every registered
// provider constructor returns a usable, non-nil Provider.
func TestProviderRegistry_AllEntriesReturnNonNil(t *testing.T) {
	for name, constructor := range cloud.ProviderRegistry {
		t.Run(name, func(t *testing.T) {
			p := constructor()
			if p == nil {
				t.Errorf("provider %q: constructor returned nil", name)
			}
		})
	}
}

// TestProviderRegistry_KnownProviders ensures expected providers are present.
func TestProviderRegistry_KnownProviders(t *testing.T) {
	expected := []string{"aws", "mock"}
	for _, name := range expected {
		if _, ok := cloud.ProviderRegistry[name]; !ok {
			t.Errorf("expected provider %q to be registered", name)
		}
	}
}

// TestGetProvider_RoundTrip checks that GetProvider is consistent with the registry.
func TestGetProvider_RoundTrip(t *testing.T) {
	for name := range cloud.ProviderRegistry {
		t.Run(name, func(t *testing.T) {
			p, ok := cloud.GetProvider(name)
			if !ok {
				t.Errorf("GetProvider(%q) returned ok=false", name)
			}
			if p == nil {
				t.Errorf("GetProvider(%q) returned nil provider", name)
			}
		})
	}
}
