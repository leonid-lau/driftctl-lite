package tfstate

import (
	"testing"
)

func clusterResource(id, cluster string) Resource {
	return Resource{
		ID:   id,
		Type: "aws_instance",
		Metadata: map[string]string{"cluster": cluster},
	}
}

func TestFilterByCluster_Match(t *testing.T) {
	resources := []Resource{
		clusterResource("r1", "prod"),
		clusterResource("r2", "staging"),
		clusterResource("r3", "prod"),
	}
	got := FilterByCluster(resources, "prod", DefaultClusterFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByCluster_EmptyCluster_ReturnsAll(t *testing.T) {
	resources := []Resource{clusterResource("r1", "prod")}
	got := FilterByCluster(resources, "", DefaultClusterFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByCluster_CaseInsensitive(t *testing.T) {
	resources := []Resource{clusterResource("r1", "Prod")}
	got := FilterByCluster(resources, "prod", DefaultClusterFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByCluster_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{clusterResource("r1", "Prod")}
	opts := ClusterFilterOptions{CaseSensitive: true}
	got := FilterByCluster(resources, "prod", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByClusters_ORSemantics(t *testing.T) {
	resources := []Resource{
		clusterResource("r1", "prod"),
		clusterResource("r2", "staging"),
		clusterResource("r3", "dev"),
	}
	got := FilterByClusters(resources, []string{"prod", "dev"}, DefaultClusterFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
