package api

import "context"

func (c *Client) ListConnections(ctx context.Context, params ListParams) (*PaginatedResponse[Connection], error) {
	var resp PaginatedResponse[Connection]
	err := c.GetJSON(ctx, BuildListURL(c.generalURL("/connections"), params), &resp)
	return &resp, err
}

func (c *Client) GetConnection(ctx context.Context, id string) (*Connection, error) {
	var conn Connection
	err := c.GetJSON(ctx, c.generalURL("/connections/"+id), &conn)
	return &conn, err
}

func (c *Client) DeleteConnection(ctx context.Context, id string) error {
	return c.DeleteEmpty(ctx, c.generalURL("/connections/"+id))
}
