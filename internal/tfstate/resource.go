package tfstate

// Resource represents a single Terraform-managed cloud resource.
type Resource struct {
	// Type is the Terraform resource type, e.g. "aws_instance".
	Type string `json:"type"`

	// ID is the primary identifier of the resource.
	ID string `json:"id"`

	// Name is the logical Terraform resource name.
	Name string `json:"name"`

	// Namespace groups resources into logical scopes (e.g. Kubernetes namespace
	// or a custom organisational boundary).
	Namespace string `json:"namespace,omitempty"`

	// Attributes holds the raw attribute map from the state.
	Attributes map[string]interface{} `json:"attributes"`

	// Tags are AWS-style key/value metadata pairs.
	Tags map[string]string `json:"tags,omitempty"`

	// Labels are Kubernetes-style key/value metadata pairs.
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations are Kubernetes-style annotation pairs.
	Annotations map[string]string `json:"annotations,omitempty"`
}
