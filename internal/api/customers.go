package api

import "context"

func (c *Client) ListCustomers(ctx context.Context, params ListParams) (*PaginatedResponse[Customer], error) {
	var resp PaginatedResponse[Customer]
	err := c.GetJSON(ctx, BuildListURL(c.companyURL("/customers"), params), &resp)
	return &resp, err
}

func (c *Client) GetCustomer(ctx context.Context, id string) (*Customer, error) {
	var cust Customer
	err := c.GetJSON(ctx, c.companyURL("/customers/"+id), &cust)
	return &cust, err
}

func (c *Client) CreateCustomer(ctx context.Context, req CreateCustomerRequest) (*Customer, error) {
	var cust Customer
	err := c.PostJSON(ctx, c.companyURL("/customers"), req, &cust)
	return &cust, err
}

func (c *Client) UpdateCustomer(ctx context.Context, id string, req UpdateCustomerRequest) (*Customer, error) {
	var cust Customer
	err := c.PutJSON(ctx, c.companyURL("/customers/"+id), req, &cust)
	return &cust, err
}

func (c *Client) DeleteCustomer(ctx context.Context, id string) error {
	return c.DeleteEmpty(ctx, c.companyURL("/customers/"+id))
}
