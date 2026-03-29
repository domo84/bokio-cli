package api

import "context"

func (c *Client) ListBankPayments(ctx context.Context, params ListParams) (*PaginatedResponse[BankPayment], error) {
	var resp PaginatedResponse[BankPayment]
	err := c.GetJSON(ctx, BuildListURL(c.companyURL("/bank-payments"), params), &resp)
	return &resp, err
}

func (c *Client) GetBankPayment(ctx context.Context, id string) (*BankPayment, error) {
	var payment BankPayment
	err := c.GetJSON(ctx, c.companyURL("/bank-payments/"+id), &payment)
	return &payment, err
}

func (c *Client) CreateBankPayment(ctx context.Context, req CreateBankPaymentRequest) (*BankPayment, error) {
	var payment BankPayment
	err := c.PostJSON(ctx, c.companyURL("/bank-payments"), req, &payment)
	return &payment, err
}
