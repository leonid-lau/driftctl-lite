package tfstate

import (
	"sort"
	"testing"
)

func idxClusterResource(id, cluster string) Resource {
	return Resource{
		ID:       id,
		Type:     "aws_instance",
		Metadata: map[string]string{"cluster": cluster},
	}
}

func TestBuildClusterIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxClusterResource("r1", "prod"),
		idxClusterResource("r2", "staging"),
		idxClusterResource("r3", "prod"),
	}
	idx := BuildClusterIndex(resources)
	got := idx.Lookup("prod")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestBuildClusterIndex_LookupMissing(t *testing.T) {
	idx := BuildClusterIndex([]Resource{idxClusterResource("r1", "prod")})
	if got := idx.Lookup("unknown"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestBuildClusterIndex_Clusters(t *testing.T) {
	resources := []Resource{
		idxClusterResource("r1", "prod"),
		idxClusterResource("r2", "staging"),
	}
	idx := BuildClusterIndex(resources)
	clusters := idx.Clusters()
	sort.Strings(clusters)
	if len(clusters) != 2 || clusters[0] != "prod" || clusters[1] != "staging" {
		t.Fatalf("unexpected clusters: %v", clusters)
	}
}

func TestBuildClusterIndex_EmptyInput(t *testing.T) {
	idx := BuildClusterIndex(nil)
	if len(idx.Clusters()) != 0 {
		t.Fatal("expected empty index")
	}
}
