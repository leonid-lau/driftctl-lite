package tfstate

// Resource represents a single Terraform-managed cloud resource.
type Resource struct {
	// ID is the unique cloud identifier (e.g. ARN, resource ID).
	ID string
	// Type is the Terraform resource type (e.g. aws_instance).
	Type string
	// Name is the Terraform logical name.
	Name string
	// Namespace groups resources (e.g. Kubernetes namespace or AWS account).
	Namespace string
	// Attributes holds the resource's configuration key-value pairs.
	Attributes map[string]interface{}
	// Tags are cloud provider tags attached to the resource.
	Tags map[string]string
	// Labels are arbitrary key-value metadata (Kubernetes-style).
	Labels map[string]string
	// Annotations are non-identifying metadata.
	Annotations map[string]string
	// Owner identifies the team or individual responsible.
	Owner string
	// Status is the current lifecycle status of the resource.
	Status string
	// Region is the cloud region where the resource lives.
	Region string
	// Environment identifies the deployment environment (e.g. prod, staging).
	Environment string
	// Severity indicates the drift severity level if applicable.
	Severity string
	// Priority indicates the remediation priority.
	Priority string
	// Source indicates the origin of the resource record (e.g. terraform, manual).
	Source string
}
