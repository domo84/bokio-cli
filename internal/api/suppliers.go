package api

import "context"

func (c *Client) ListSuppliers(ctx context.Context, params ListParams) (*PaginatedResponse[Supplier], error) {
	var resp PaginatedResponse[Supplier]
	err := c.GetJSON(ctx, BuildListURL(c.companyURL("/suppliers"), params), &resp)
	return &resp, err
}

func (c *Client) GetSupplier(ctx context.Context, id string) (*Supplier, error) {
	var s Supplier
	err := c.GetJSON(ctx, c.companyURL("/suppliers/"+id), &s)
	return &s, err
}

func (c *Client) UpdateSupplier(ctx context.Context, id string, req UpdateSupplierRequest) (*Supplier, error) {
	var s Supplier
	err := c.PutJSON(ctx, c.companyURL("/suppliers/"+id), req, &s)
	return &s, err
}

func (c *Client) CreateSupplier(ctx context.Context, req CreateSupplierRequest) (*Supplier, error) {
	var s Supplier
	err := c.PostJSON(ctx, c.companyURL("/suppliers"), req, &s)
	return &s, err
}
