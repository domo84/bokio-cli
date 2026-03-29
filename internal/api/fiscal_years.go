package api

import "context"

func (c *Client) ListFiscalYears(ctx context.Context) (*PaginatedResponse[FiscalYear], error) {
	var resp PaginatedResponse[FiscalYear]
	err := c.GetJSON(ctx, c.companyURL("/fiscal-years"), &resp)
	return &resp, err
}

func (c *Client) GetFiscalYear(ctx context.Context, id string) (*FiscalYear, error) {
	var year FiscalYear
	err := c.GetJSON(ctx, c.companyURL("/fiscal-years/"+id), &year)
	return &year, err
}
