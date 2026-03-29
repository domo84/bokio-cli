package api

import "context"

func (c *Client) ListItems(ctx context.Context, params ListParams) (*PaginatedResponse[Item], error) {
	var resp PaginatedResponse[Item]
	err := c.GetJSON(ctx, BuildListURL(c.companyURL("/items"), params), &resp)
	return &resp, err
}

func (c *Client) GetItem(ctx context.Context, id string) (*Item, error) {
	var item Item
	err := c.GetJSON(ctx, c.companyURL("/items/"+id), &item)
	return &item, err
}

func (c *Client) CreateItem(ctx context.Context, req CreateItemRequest) (*Item, error) {
	var item Item
	err := c.PostJSON(ctx, c.companyURL("/items"), req, &item)
	return &item, err
}

func (c *Client) UpdateItem(ctx context.Context, id string, req UpdateItemRequest) (*Item, error) {
	var item Item
	err := c.PutJSON(ctx, c.companyURL("/items/"+id), req, &item)
	return &item, err
}

func (c *Client) DeleteItem(ctx context.Context, id string) error {
	return c.DeleteEmpty(ctx, c.companyURL("/items/"+id))
}
