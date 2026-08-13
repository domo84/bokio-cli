package api

import (
	"encoding/json"
	"testing"
)

// exampleEntry mirrors the shape of a real GET /journal-entries/{id} response: a
// reverse-charge EU service purchase booked from a bank event. Identifiers, amounts
// and the counterparty are synthetic; the structure is what matters here, including
// item ids that do not follow the order of the items array.
const exampleEntry = `{
  "id": "11111111-2222-4333-8444-555555555555",
  "title": "Example Vendor AB - Inköp tjänster inom EU 25%",
  "journalEntryNumber": "V42",
  "date": "2026-01-15",
  "items": [
    { "id": 1001, "debit": 0.0,   "credit": 100.00, "account": 1930, "tags": [] },
    { "id": 1004, "debit": 0.0,   "credit": 25.00,  "account": 2614, "tags": [] },
    { "id": 1003, "debit": 25.00, "credit": 0.0,    "account": 2645, "tags": [] },
    { "id": 1002, "debit": 100.00, "credit": 0.0,   "account": 4535, "tags": [] }
  ],
  "tags": [],
  "reversingJournalEntryId": null,
  "reversedByJournalEntryId": null
}`

func TestJournalEntryDecodesEntry(t *testing.T) {
	var e JournalEntry
	if err := json.Unmarshal([]byte(exampleEntry), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if e.Title != "Example Vendor AB - Inköp tjänster inom EU 25%" {
		t.Errorf("Title = %q", e.Title)
	}
	if e.JournalEntryNumber != "V42" {
		t.Errorf("JournalEntryNumber = %q", e.JournalEntryNumber)
	}
	if e.Date != "2026-01-15" {
		t.Errorf("Date = %q", e.Date)
	}
	if len(e.Items) != 4 {
		t.Fatalf("Items = %d, want 4", len(e.Items))
	}
	if e.IsReversed() {
		t.Error("IsReversed() = true, want false")
	}

	debit, credit := e.Balance()
	if debit != credit {
		t.Errorf("unbalanced: debit %v != credit %v", debit, credit)
	}
	if debit != 125.00 {
		t.Errorf("debit total = %v, want 125.00", debit)
	}

	// Line identity must come from account + amount, not slice position: the server
	// assigns ids in an order unrelated to the array order.
	byAccount := map[int]JournalEntryItem{}
	for _, item := range e.Items {
		byAccount[item.Account] = item
	}
	if got := byAccount[1930].Credit; got != 100.00 {
		t.Errorf("1930 credit = %v, want 100.00", got)
	}
	if got := byAccount[4535].Debit; got != 100.00 {
		t.Errorf("4535 debit = %v, want 100.00", got)
	}
	if got := byAccount[2614].Credit; got != 25.00 {
		t.Errorf("2614 credit = %v, want 25.00", got)
	}
	if got := byAccount[2645].Debit; got != 25.00 {
		t.Errorf("2645 debit = %v, want 25.00", got)
	}
}

func TestJournalEntryReversalFields(t *testing.T) {
	const reversed = `{"id":"a","date":"2026-01-01","items":[],
	  "reversedByJournalEntryId":"b0000000-0000-0000-0000-000000000000"}`

	var e JournalEntry
	if err := json.Unmarshal([]byte(reversed), &e); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if !e.IsReversed() {
		t.Error("IsReversed() = false, want true")
	}
}

// The wire format is title/date/items with debit/credit/account. Guards against
// regressing to the description/rows naming, which the API silently ignores.
func TestCreateJournalEntryRequestWireFormat(t *testing.T) {
	req := CreateJournalEntryRequest{
		Title: "Test entry",
		Date:  "2026-01-15",
		Items: []JournalEntryItem{
			{Account: 1930, Credit: 100.00},
			{Account: 4535, Debit: 100.00},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	for _, key := range []string{"title", "date", "items"} {
		if _, ok := got[key]; !ok {
			t.Errorf("missing key %q in %s", key, data)
		}
	}
	for _, key := range []string{"description", "rows"} {
		if _, ok := got[key]; ok {
			t.Errorf("unexpected legacy key %q in %s", key, data)
		}
	}

	items, ok := got["items"].([]any)
	if !ok || len(items) != 2 {
		t.Fatalf("items = %v", got["items"])
	}
	first, ok := items[0].(map[string]any)
	if !ok {
		t.Fatalf("item 0 = %v", items[0])
	}
	for _, key := range []string{"debit", "credit", "account"} {
		if _, ok := first[key]; !ok {
			t.Errorf("item missing key %q in %s", key, data)
		}
	}
	// id is server-assigned and must not be sent for new lines.
	if _, ok := first["id"]; ok {
		t.Errorf("item should omit id: %s", data)
	}
}
