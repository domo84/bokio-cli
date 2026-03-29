package api

import "context"

// GetCompanyInfo retrieves company information.
func (c *Client) GetCompanyInfo(ctx context.Context) (*CompanyInfo, error) {
	var info CompanyInfo
	err := c.GetJSON(ctx, c.companyURL("/company-information"), &info)
	return &info, err
}
