package api

import (
	"encoding/json"
	"testing"
)

// Payload taken verbatim from the customer schema example in company-api.yaml.
func TestCustomerDecodesSpecExample(t *testing.T) {
	const example = `{
	  "id": "55c899c5-82b2-47fa-9c51-e35fc9b26443",
	  "name": "customer 1",
	  "type": "company",
	  "vatNumber": "SE1234567890",
	  "orgNumber": "123456-7890",
	  "paymentTerms": "30",
	  "contactsDetails": [
	    {"name": "John Doe", "email": "john@email.com", "phone": "0927-5631505", "isDefault": true}
	  ],
	  "address": {
	    "line1": "Älvsborgsvägen 10",
	    "line2": null,
	    "city": "Göteborg",
	    "postalCode": "123 45",
	    "country": "SE"
	  },
	  "language": "sv",
	  "modifiedDateTime": "2024-10-10T00:00:00Z"
	}`

	var c Customer
	if err := json.Unmarshal([]byte(example), &c); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if c.Type != "company" || c.OrgNumber != "123456-7890" {
		t.Errorf("type/orgNumber = %q/%q", c.Type, c.OrgNumber)
	}
	if c.Address == nil || c.Address.City != "Göteborg" || c.Address.Line2 != nil {
		t.Errorf("address = %+v", c.Address)
	}
	contact := c.DefaultContact()
	if contact == nil || contact.Email != "john@email.com" {
		t.Errorf("default contact = %+v", contact)
	}

	// modifiedDateTime is null for customers changed before it was introduced.
	var old Customer
	if err := json.Unmarshal([]byte(`{"id":"x","name":"n","type":"private","modifiedDateTime":null}`), &old); err != nil {
		t.Fatalf("unmarshal old: %v", err)
	}
	if old.ModifiedDateTime != nil || old.DefaultContact() != nil {
		t.Errorf("old = %+v", old)
	}
}
