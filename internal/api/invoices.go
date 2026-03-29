package api

import (
	"context"
	"fmt"
)

func (c *Client) ListInvoices(ctx context.Context, params ListParams) (*PaginatedResponse[Invoice], error) {
	var resp PaginatedResponse[Invoice]
	err := c.GetJSON(ctx, BuildListURL(c.companyURL("/invoices"), params), &resp)
	return &resp, err
}

func (c *Client) GetInvoice(ctx context.Context, id string) (*Invoice, error) {
	var inv Invoice
	err := c.GetJSON(ctx, c.companyURL("/invoices/"+id), &inv)
	return &inv, err
}

func (c *Client) CreateInvoice(ctx context.Context, req CreateInvoiceRequest) (*Invoice, error) {
	var inv Invoice
	err := c.PostJSON(ctx, c.companyURL("/invoices"), req, &inv)
	return &inv, err
}

func (c *Client) UpdateInvoice(ctx context.Context, id string, req UpdateInvoiceRequest) (*Invoice, error) {
	var inv Invoice
	err := c.PutJSON(ctx, c.companyURL("/invoices/"+id), req, &inv)
	return &inv, err
}

func (c *Client) PublishInvoice(ctx context.Context, id string) (*Invoice, error) {
	var inv Invoice
	err := c.PostJSON(ctx, c.companyURL("/invoices/"+id+"/publish"), nil, &inv)
	return &inv, err
}

func (c *Client) RecordInvoice(ctx context.Context, id string, req RecordInvoiceRequest) error {
	return c.PostEmpty(ctx, c.companyURL("/invoices/"+id+"/record"), req)
}

// Line Items
func (c *Client) AddInvoiceLineItem(ctx context.Context, invoiceID string, item LineItem) (*Invoice, error) {
	var inv Invoice
	err := c.PostJSON(ctx, c.companyURL(fmt.Sprintf("/invoices/%s/line-items", invoiceID)), item, &inv)
	return &inv, err
}

func (c *Client) UpdateInvoiceLineItem(ctx context.Context, invoiceID, lineItemID string, item LineItem) (*Invoice, error) {
	var inv Invoice
	err := c.PutJSON(ctx, c.companyURL(fmt.Sprintf("/invoices/%s/line-items/%s", invoiceID, lineItemID)), item, &inv)
	return &inv, err
}

func (c *Client) DeleteInvoiceLineItem(ctx context.Context, invoiceID, lineItemID string) error {
	return c.DeleteEmpty(ctx, c.companyURL(fmt.Sprintf("/invoices/%s/line-items/%s", invoiceID, lineItemID)))
}

// Attachments
func (c *Client) ListInvoiceAttachments(ctx context.Context, invoiceID string) ([]InvoiceAttachment, error) {
	var attachments []InvoiceAttachment
	err := c.GetJSON(ctx, c.companyURL(fmt.Sprintf("/invoices/%s/attachments", invoiceID)), &attachments)
	return attachments, err
}

func (c *Client) DeleteInvoiceAttachment(ctx context.Context, invoiceID, attachmentID string) error {
	return c.DeleteEmpty(ctx, c.companyURL(fmt.Sprintf("/invoices/%s/attachments/%s", invoiceID, attachmentID)))
}

// Payments
func (c *Client) ListInvoicePayments(ctx context.Context, invoiceID string) ([]InvoicePayment, error) {
	var payments []InvoicePayment
	err := c.GetJSON(ctx, c.companyURL(fmt.Sprintf("/invoices/%s/payments", invoiceID)), &payments)
	return payments, err
}

func (c *Client) AddInvoicePayment(ctx context.Context, invoiceID string, req CreatePaymentRequest) (*InvoicePayment, error) {
	var payment InvoicePayment
	err := c.PostJSON(ctx, c.companyURL(fmt.Sprintf("/invoices/%s/payments", invoiceID)), req, &payment)
	return &payment, err
}

// Settlements
func (c *Client) AddInvoiceSettlement(ctx context.Context, invoiceID string, req CreateSettlementRequest) error {
	return c.PostEmpty(ctx, c.companyURL(fmt.Sprintf("/invoices/%s/settlements", invoiceID)), req)
}
