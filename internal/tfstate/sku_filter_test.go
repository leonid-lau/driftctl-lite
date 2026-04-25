package tfstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func skuResource(sku string) Resource {
	attrs := map[string]interface{}{}
	if sku != "" {
		attrs["sku"] = sku
	}
	return Resource{Type: "azurerm_vm", Attributes: attrs}
}

func TestFilterBySKU_Match(t *testing.T) {
	resources := []Resource{skuResource("Standard_D2s_v3"), skuResource("Basic_A1"), skuResource("")}
	got := FilterBySKU(resources, "Standard_D2s_v3", DefaultSKUFilterOptions())
	assert.Len(t, got, 1)
	assert.Equal(t, "Standard_D2s_v3", got[0].Attributes["sku"])
}

func TestFilterBySKU_EmptySKU_ReturnsAll(t *testing.T) {
	resources := []Resource{skuResource("Standard_D2s_v3"), skuResource("Basic_A1")}
	got := FilterBySKU(resources, "", DefaultSKUFilterOptions())
	assert.Len(t, got, 2)
}

func TestFilterBySKU_CaseInsensitive(t *testing.T) {
	resources := []Resource{skuResource("Standard_D2s_v3")}
	got := FilterBySKU(resources, "standard_d2s_v3", DefaultSKUFilterOptions())
	assert.Len(t, got, 1)
}

func TestFilterBySKU_CaseSensitive_NoMatch(t *testing.T) {
	resources := []Resource{skuResource("Standard_D2s_v3")}
	opts := SKUFilterOptions{CaseSensitive: true}
	got := FilterBySKU(resources, "standard_d2s_v3", opts)
	assert.Empty(t, got)
}

func TestFilterBySKUs_ORSemantics(t *testing.T) {
	resources := []Resource{skuResource("Standard_D2s_v3"), skuResource("Basic_A1"), skuResource("Premium_LRS")}
	got := FilterBySKUs(resources, []string{"Standard_D2s_v3", "Premium_LRS"}, DefaultSKUFilterOptions())
	assert.Len(t, got, 2)
}

func TestBuildSKUIndex_Lookup(t *testing.T) {
	resources := []Resource{skuResource("Standard_D2s_v3"), skuResource("Basic_A1")}
	idx := BuildSKUIndex(resources)
	result := idx.Lookup("Standard_D2s_v3")
	assert.Len(t, result, 1)
}

func TestBuildSKUIndex_LookupCaseInsensitive(t *testing.T) {
	resources := []Resource{skuResource("Standard_D2s_v3")}
	idx := BuildSKUIndex(resources)
	assert.Len(t, idx.Lookup("STANDARD_D2S_V3"), 1)
}

func TestBuildSKUIndex_LookupMissing(t *testing.T) {
	idx := BuildSKUIndex([]Resource{skuResource("Basic_A1")})
	assert.Nil(t, idx.Lookup("nonexistent"))
}

func TestBuildSKUIndex_SKUs(t *testing.T) {
	resources := []Resource{skuResource("Standard_D2s_v3"), skuResource("Basic_A1"), skuResource("Basic_A1")}
	idx := BuildSKUIndex(resources)
	assert.Len(t, idx.SKUs(), 2)
}
