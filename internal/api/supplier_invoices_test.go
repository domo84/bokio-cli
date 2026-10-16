package api

import (
	"encoding/json"
	"strings"
	"testing"
)

// Payload taken from the supplierInvoiceGet schema example in company-api.yaml.
func TestSupplierInvoiceDecodesSpecExample(t *testing.T) {
	const example = `{
	  "id": "6f6b2f17-3c4a-4a69-9f2d-3f5d7e1b8c90",
	  "supplierRef": {"id": "6f6b2f17-3c4a-4a69-9f2d-3f5d7e1b8c90", "name": "Supplier AB"},
	  "invoiceNumber": "INV-001",
	  "invoiceDate": "2023-10-01",
	  "dueDate": "2023-10-15",
	  "totalAmount": 1000,
	  "remainingAmount": 1000,
	  "currency": "SEK",
	  "currencyRate": 1,
	  "journalEntryRef": {"id": "2d8ce6b4-1f3a-46a5-8a9d-5b7f2c4e9d10"},
	  "lineItems": [{
	    "id": "a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
	    "description": "Consulting",
	    "quantity": 10,
	    "unitPrice": 1000,
	    "unitType": "hour",
	    "taxRate": 25,
	    "projectRef": {"id": "4f8f6da6-4cbc-4bfc-8d23-67d5f6f2dd9b", "projectName": "Website Redesign", "projectNumber": "P-001"},
	    "costCenterRef": {"tagId": "7e5ee25f-1f5d-4e47-a26f-5d2d20fbef90", "tagName": "Nord"}
	  }],
	  "uploadRefs": [{"id": "24f3b7c1-8d6a-4e2f-9b31-5d7c2a8e4f90"}, {"id": "34e1c9a5-7b2d-4f86-8a31-6d2c5f9b7e40"}]
	}`

	var inv SupplierInvoice
	if err := json.Unmarshal([]byte(example), &inv); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if inv.SupplierRef.Name != "Supplier AB" || inv.RemainingAmount != 1000 {
		t.Errorf("supplier/remaining = %q/%v", inv.SupplierRef.Name, inv.RemainingAmount)
	}
	if len(inv.LineItems) != 1 || len(inv.UploadRefs) != 2 {
		t.Fatalf("lineItems/uploadRefs = %d/%d", len(inv.LineItems), len(inv.UploadRefs))
	}
	li := inv.LineItems[0]
	if li.TaxRate == nil || *li.TaxRate != 25 || li.ProjectRef.ProjectNumber != "P-001" || li.CostCenterRef.TagName != "Nord" {
		t.Errorf("line item = %+v", li)
	}

	// The update body keeps line items (a PUT removes any left out) and sends
	// only the supplier id.
	body, err := json.Marshal(inv.UpdateRequest())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	s := string(body)
	if !strings.Contains(s, `"supplierRef":{"id":"6f6b2f17-3c4a-4a69-9f2d-3f5d7e1b8c90"}`) {
		t.Errorf("supplierRef not reduced to id: %s", s)
	}
	if !strings.Contains(s, `"id":"a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d"`) {
		t.Errorf("line item id missing: %s", s)
	}
	for _, readOnly := range []string{"remainingAmount", "currencyRate", "journalEntryRef", "uploadRefs"} {
		if strings.Contains(s, readOnly) {
			t.Errorf("update body has read-only %s: %s", readOnly, s)
		}
	}
}
