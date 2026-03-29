package api

import (
	"context"
	"io"
)

func (c *Client) DownloadSIE(ctx context.Context, fiscalYearID string, w io.Writer) error {
	resp, err := c.Get(ctx, c.companyURL("/sie/"+fiscalYearID+"/download"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	return err
}
