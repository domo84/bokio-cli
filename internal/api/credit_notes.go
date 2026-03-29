package api

import "context"

func (c *Client) ListCreditNotes(ctx context.Context, params ListParams) (*PaginatedResponse[CreditNote], error) {
	var resp PaginatedResponse[CreditNote]
	err := c.GetJSON(ctx, BuildListURL(c.companyURL("/credit-notes"), params), &resp)
	return &resp, err
}

func (c *Client) GetCreditNote(ctx context.Context, id string) (*CreditNote, error) {
	var note CreditNote
	err := c.GetJSON(ctx, c.companyURL("/credit-notes/"+id), &note)
	return &note, err
}

func (c *Client) UpdateCreditNote(ctx context.Context, id string, req UpdateCreditNoteRequest) (*CreditNote, error) {
	var note CreditNote
	err := c.PutJSON(ctx, c.companyURL("/credit-notes/"+id), req, &note)
	return &note, err
}

func (c *Client) PublishCreditNote(ctx context.Context, id string) (*CreditNote, error) {
	var note CreditNote
	err := c.PostJSON(ctx, c.companyURL("/credit-notes/"+id+"/publish"), nil, &note)
	return &note, err
}

func (c *Client) RecordCreditNote(ctx context.Context, id string, req RecordCreditNoteRequest) error {
	return c.PostEmpty(ctx, c.companyURL("/credit-notes/"+id+"/record"), req)
}
