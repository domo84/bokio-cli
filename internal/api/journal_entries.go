package api

import "context"

func (c *Client) ListJournalEntries(ctx context.Context, params ListParams) (*PaginatedResponse[JournalEntry], error) {
	var resp PaginatedResponse[JournalEntry]
	err := c.GetJSON(ctx, BuildListURL(c.companyURL("/journal-entries"), params), &resp)
	return &resp, err
}

func (c *Client) GetJournalEntry(ctx context.Context, id string) (*JournalEntry, error) {
	var entry JournalEntry
	err := c.GetJSON(ctx, c.companyURL("/journal-entries/"+id), &entry)
	return &entry, err
}

func (c *Client) CreateJournalEntry(ctx context.Context, req CreateJournalEntryRequest) (*JournalEntry, error) {
	var entry JournalEntry
	err := c.PostJSON(ctx, c.companyURL("/journal-entries"), req, &entry)
	return &entry, err
}

func (c *Client) ReverseJournalEntry(ctx context.Context, id string) (*JournalEntry, error) {
	var entry JournalEntry
	err := c.PostJSON(ctx, c.companyURL("/journal-entries/"+id+"/reverse"), nil, &entry)
	return &entry, err
}
