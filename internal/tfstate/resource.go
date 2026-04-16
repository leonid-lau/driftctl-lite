package tfstate

// Resource represents a single Terraform-managed infrastructure resource.
type Resource struct {
	// Type is the Terraform resource type, e.g. "aws_s3_bucket".
	Type string `json:"type"`

	// ID is the unique identifier of the resource within its type.
	ID string `json:"id"`

	// Name is the logical Terraform resource name.
	Name string `json:"name"`

	// Attributes holds the resource's configuration attributes.
	Attributes map[string]interface{} `json:"attributes"`

	// Tags are cloud-provider tags attached to the resource.
	Tags map[string]string `json:"tags,omitempty"`

	// Labels are Kubernetes-style labels on the resource.
	Labels map[string]string `json:"labels,omitempty"`

	// Annotations are arbitrary metadata annotations on the resource.
	Annotations map[string]string `json:"annotations,omitempty"`
}

// Key returns a unique string key for the resource combining type and ID.
func (r Resource) Key() string {
	if r.ID != "" {
		return r.Type + "/" + r.ID
	}
	return r.Type + "/" + r.Name
}
