package tfstate

// Resource represents a single resource entry from a Terraform state file.
type Resource struct {
	// Type is the Terraform resource type, e.g. "aws_instance".
	Type string `json:"type"`
	// Name is the logical name given in the Terraform configuration.
	Name string `json:"name"`
	// ProviderName identifies the provider, e.g. "registry.terraform.io/hashicorp/aws".
	ProviderName string `json:"provider_name"`
	// Attributes holds the resource's attribute key-value pairs as stored in state.
	Attributes map[string]interface{} `json:"attributes"`
}

// ID returns the value of the "id" attribute if present, otherwise falls back
// to the resource Name. This mirrors the behaviour of IndexByID in filter.go.
func (r Resource) ID() string {
	if id, ok := r.Attributes["id"].(string); ok && id != "" {
		return id
	}
	return r.Name
}

// HasAttribute reports whether the resource has a non-nil value for the given
// attribute key.
func (r Resource) HasAttribute(key string) bool {
	_, ok := r.Attributes[key]
	return ok
}

// AttributeString returns the string representation of an attribute value.
// Returns an empty string when the key is absent.
func (r Resource) AttributeString(key string) string {
	v, ok := r.Attributes[key]
	if !ok || v == nil {
		return ""
	}
	switch s := v.(type) {
	case string:
		return s
	default:
		return ""
	}
}
