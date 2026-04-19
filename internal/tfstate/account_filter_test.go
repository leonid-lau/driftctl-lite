package tfstate

import (
	"testing"
)

func accountResource(id, account string) Resource {
	return Resource{
		Type: "aws_s3_bucket",
		Name: id,
		Attributes: map[string]interface{}{"account": account},
	}
}

func accountIDResource(id, accountID string) Resource {
	return Resource{
		Type: "aws_s3_bucket",
		Name: id,
		Attributes: map[string]interface{}{"account_id": accountID},
	}
}

func TestFilterByAccount_Match(t *testing.T) {
	resources := []Resource{
		accountResource("a", "123456789012"),
		accountResource("b", "999999999999"),
	}
	got := FilterByAccount(resources, "123456789012", DefaultAccountFilterOptions())
	if len(got) != 1 || got[0].Name != "a" {
		t.Fatalf("expected 1 match, got %v", got)
	}
}

func TestFilterByAccount_EmptyAccount_ReturnsAll(t *testing.T) {
	resources := []Resource{accountResource("a", "111"), accountResource("b", "222")}
	got := FilterByAccount(resources, "", DefaultAccountFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilterByAccount_CaseInsensitive(t *testing.T) {
	resources := []Resource{accountResource("a", "ABCDEF")}
	got := FilterByAccount(resources, "abcdef", DefaultAccountFilterOptions())
	if len(got) != 1 {
		t.Fatal("expected case-insensitive match")
	}
}

func TestFilterByAccount_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{accountResource("a", "ABCDEF")}
	opts := AccountFilterOptions{CaseSensitive: true}
	got := FilterByAccount(resources, "abcdef", opts)
	if len(got) != 0 {
		t.Fatal("expected no match with case-sensitive filter")
	}
}

func TestFilterByAccount_FallbackToAccountID(t *testing.T) {
	resources := []Resource{accountIDResource("a", "777777777777")}
	got := FilterByAccount(resources, "777777777777", DefaultAccountFilterOptions())
	if len(got) != 1 {
		t.Fatal("expected fallback to account_id attribute")
	}
}

func TestFilterByAccounts_ORSemantics(t *testing.T) {
	resources := []Resource{
		accountResource("a", "111"),
		accountResource("b", "222"),
		accountResource("c", "333"),
	}
	got := FilterByAccounts(resources, []string{"111", "333"}, DefaultAccountFilterOptions())
	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(got))
	}
}
