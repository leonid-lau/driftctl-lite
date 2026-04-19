package tfstate

import (
	"testing"
)

func idxAccountResource(account string) Resource {
	return Resource{
		Type: "aws_s3_bucket",
		Attributes: map[string]string{"account": account},
	}
}

func idxAccountIDResource(accountID string) Resource {
	return Resource{
		Type: "aws_s3_bucket",
		Attributes: map[string]string{"account_id": accountID},
	}
}

func TestBuildAccountIndex_Lookup(t *testing.T) {
	resources := []Resource{idxAccountResource("123456789"), idxAccountResource("999")}
	idx := BuildAccountIndex(resources)
	got := idx.Lookup("123456789")
	if len(got) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(got))
	}
}

func TestBuildAccountIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{idxAccountResource("ProdAccount")}
	idx := BuildAccountIndex(resources)
	got := idx.Lookup("prodaccount")
	if len(got) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(got))
	}
}

func TestBuildAccountIndex_FallbackToAccountID(t *testing.T) {
	resources := []Resource{idxAccountIDResource("777")}
	idx := BuildAccountIndex(resources)
	got := idx.Lookup("777")
	if len(got) != 1 {
		t.Fatalf("expected 1 resource via account_id fallback, got %d", len(got))
	}
}

func TestBuildAccountIndex_LookupMissing(t *testing.T) {
	idx := BuildAccountIndex([]Resource{idxAccountResource("123")})
	if got := idx.Lookup("999"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestBuildAccountIndex_Accounts(t *testing.T) {
	resources := []Resource{idxAccountResource("aaa"), idxAccountResource("bbb")}
	idx := BuildAccountIndex(resources)
	if len(idx.Accounts()) != 2 {
		t.Fatalf("expected 2 accounts, got %d", len(idx.Accounts()))
	}
}

func TestBuildAccountIndex_EmptyInput(t *testing.T) {
	idx := BuildAccountIndex(nil)
	if len(idx.Accounts()) != 0 {
		t.Fatal("expected empty index")
	}
}
