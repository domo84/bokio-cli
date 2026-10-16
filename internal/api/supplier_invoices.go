package api

import (
	"context"
	"fmt"
	"net/url"
)

// ListSupplierInvoices lists supplier invoices. A non-empty projectID returns
// only invoices with a line item on that project.
func (c *Client) ListSupplierInvoices(ctx context.Context, params ListParams, projectID string) (*PaginatedResponse[SupplierInvoice], error) {
	base := c.companyURL("/supplier-invoices")
	if projectID != "" {
		base += "?projectId=" + url.QueryEscape(projectID)
	}
	var resp PaginatedResponse[SupplierInvoice]
	err := c.GetJSON(ctx, BuildListURL(base, params), &resp)
	return &resp, err
}

func (c *Client) GetSupplierInvoice(ctx context.Context, id string) (*SupplierInvoice, error) {
	var inv SupplierInvoice
	err := c.GetJSON(ctx, c.companyURL("/supplier-invoices/"+id), &inv)
	return &inv, err
}

func (c *Client) CreateSupplierInvoice(ctx context.Context, req CreateSupplierInvoiceRequest) (*SupplierInvoice, error) {
	var inv SupplierInvoice
	err := c.PostJSON(ctx, c.companyURL("/supplier-invoices"), req, &inv)
	return &inv, err
}

func (c *Client) UpdateSupplierInvoice(ctx context.Context, id string, req UpdateSupplierInvoiceRequest) (*SupplierInvoice, error) {
	var inv SupplierInvoice
	err := c.PutJSON(ctx, c.companyURL("/supplier-invoices/"+id), req, &inv)
	return &inv, err
}

func (c *Client) DeleteSupplierInvoice(ctx context.Context, id string) error {
	return c.DeleteEmpty(ctx, c.companyURL("/supplier-invoices/"+id))
}

func (c *Client) AddSupplierInvoiceUpload(ctx context.Context, invoiceID, uploadID string) (*Ref, error) {
	var ref Ref
	err := c.PostJSON(ctx, c.companyURL(fmt.Sprintf("/supplier-invoices/%s/uploads", invoiceID)), Ref{ID: uploadID}, &ref)
	return &ref, err
}

func (c *Client) AddSupplierInvoiceLineItem(ctx context.Context, invoiceID string, item SupplierInvoiceLineItem) (*SupplierInvoiceLineItem, error) {
	var li SupplierInvoiceLineItem
	err := c.PostJSON(ctx, c.companyURL(fmt.Sprintf("/supplier-invoices/%s/line-items", invoiceID)), item, &li)
	return &li, err
}
