package tfstate

import "testing"

func idxArchResource(arch string) Resource {
	return Resource{
		Type: "aws_instance",
		ID:   "i-" + arch,
		Attributes: map[string]string{
			"arch": arch,
		},
	}
}

func idxArchAltResource(arch string) Resource {
	return Resource{
		Type: "aws_instance",
		ID:   "i-alt-" + arch,
		Attributes: map[string]string{
			"architecture": arch,
		},
	}
}

func TestBuildArchIndex_Lookup(t *testing.T) {
	resources := []Resource{
		idxArchResource("x86_64"),
		idxArchResource("arm64"),
	}
	idx := BuildArchIndex(resources)

	got := idx.Lookup("x86_64")
	if len(got) != 1 {
		t.Fatalf("expected 1 result, got %d", len(got))
	}
	if got[0].ID != "i-x86_64" {
		t.Errorf("unexpected ID: %s", got[0].ID)
	}
}

func TestBuildArchIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxArchResource("X86_64")}
	idx := BuildArchIndex(resources)

	if got := idx.Lookup("x86_64"); len(got) != 1 {
		t.Errorf("expected 1 result for lower-case lookup, got %d", len(got))
	}
}

func TestBuildArchIndex_FallbackToArchitecture(t *testing.T) {
	resources := []Resource{idxArchAltResource("arm64")}
	idx := BuildArchIndex(resources)

	if got := idx.Lookup("arm64"); len(got) != 1 {
		t.Errorf("expected fallback to 'architecture' attribute, got %d results", len(got))
	}
}

func TestBuildArchIndex_LookupMissing(t *testing.T) {
	idx := BuildArchIndex([]Resource{idxArchResource("x86_64")})
	if got := idx.Lookup("arm64"); got != nil {
		t.Errorf("expected nil for missing arch, got %v", got)
	}
}

func TestBuildArchIndex_Archs(t *testing.T) {
	resources := []Resource{
		idxArchResource("x86_64"),
		idxArchResource("arm64"),
		idxArchResource("x86_64"),
	}
	idx := BuildArchIndex(resources)
	if len(idx.Archs()) != 2 {
		t.Errorf("expected 2 distinct archs, got %d", len(idx.Archs()))
	}
}
