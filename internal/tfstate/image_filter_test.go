package tfstate

import "testing"

func imageResource(image string) Resource {
	return Resource{
		Type: "aws_instance",
		Attributes: map[string]interface{}{"image": image},
	}
}

func TestFilterByImage_Match(t *testing.T) {
	resources := []Resource{imageResource("ami-12345"), imageResource("ami-99999")}
	got := FilterByImage(resources, "ami-12345", DefaultImageFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(got))
	}
}

func TestFilterByImage_EmptyImage_ReturnsAll(t *testing.T) {
	resources := []Resource{imageResource("ami-1"), imageResource("ami-2")}
	got := FilterByImage(resources, "", DefaultImageFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByImage_CaseInsensitive(t *testing.T) {
	resources := []Resource{imageResource("AMI-UPPER")}
	got := FilterByImage(resources, "ami-upper", DefaultImageFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(got))
	}
}

func TestFilterByImage_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{imageResource("AMI-UPPER")}
	opts := ImageFilterOptions{CaseSensitive: true}
	got := FilterByImage(resources, "ami-upper", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0 resources, got %d", len(got))
	}
}

func TestFilterByImages_ORSemantics(t *testing.T) {
	resources := []Resource{imageResource("ami-1"), imageResource("ami-2"), imageResource("ami-3")}
	got := FilterByImages(resources, []string{"ami-1", "ami-3"}, DefaultImageFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}

func TestFilterByImages_Empty_ReturnsAll(t *testing.T) {
	resources := []Resource{imageResource("ami-1"), imageResource("ami-2")}
	got := FilterByImages(resources, nil, DefaultImageFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(got))
	}
}
