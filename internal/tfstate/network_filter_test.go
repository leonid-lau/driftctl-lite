package tfstate

import "testing"

func networkResource(network string) Resource {
	return Resource{
		Type: "aws_instance",
		Name: "test",
		Attributes: map[string]interface{}{"network": network},
	}
}

func vpcResource(vpc string) Resource {
	return Resource{
		Type: "aws_instance",
		Name: "test",
		Attributes: map[string]interface{}{"vpc": vpc},
	}
}

func TestFilterByNetwork_Match(t *testing.T) {
	resources := []Resource{networkResource("prod-net"), networkResource("dev-net")}
	got := FilterByNetwork(resources, "prod-net", DefaultNetworkFilterOptions())
	if len(got) != 1 || got[0].Attributes["network"] != "prod-net" {
		t.Fatalf("expected 1 match, got %v", got)
	}
}

func TestFilterByNetwork_EmptyNetwork_ReturnsAll(t *testing.T) {
	resources := []Resource{networkResource("prod-net"), networkResource("dev-net")}
	got := FilterByNetwork(resources, "", DefaultNetworkFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByNetwork_CaseInsensitive(t *testing.T) {
	resources := []Resource{networkResource("PROD-NET")}
	got := FilterByNetwork(resources, "prod-net", DefaultNetworkFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByNetwork_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{networkResource("PROD-NET")}
	opts := NetworkFilterOptions{CaseSensitive: true}
	got := FilterByNetwork(resources, "prod-net", opts)
	if len(got) != 0 {
		t.Fatalf("expected 0, got %d", len(got))
	}
}

func TestFilterByNetwork_FallbackToVPC(t *testing.T) {
	resources := []Resource{vpcResource("vpc-123")}
	got := FilterByNetwork(resources, "vpc-123", DefaultNetworkFilterOptions())
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}
}

func TestFilterByNetworks_ORSemantics(t *testing.T) {
	resources := []Resource{
		networkResource("prod-net"),
		networkResource("dev-net"),
		networkResource("staging-net"),
	}
	got := FilterByNetworks(resources, []string{"prod-net", "dev-net"}, DefaultNetworkFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}
