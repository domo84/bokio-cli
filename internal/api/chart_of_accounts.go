package api

import (
	"context"
	"fmt"
)

func (c *Client) ListAccounts(ctx context.Context) ([]Account, error) {
	var accounts []Account
	err := c.GetJSON(ctx, c.companyURL("/chart-of-accounts"), &accounts)
	return accounts, err
}

func (c *Client) GetAccount(ctx context.Context, accountNumber int) (*Account, error) {
	var account Account
	err := c.GetJSON(ctx, c.companyURL(fmt.Sprintf("/chart-of-accounts/%d", accountNumber)), &account)
	return &account, err
}
