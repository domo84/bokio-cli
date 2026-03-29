package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func (c *Client) ListUploads(ctx context.Context, params ListParams) (*PaginatedResponse[Upload], error) {
	var resp PaginatedResponse[Upload]
	err := c.GetJSON(ctx, BuildListURL(c.companyURL("/uploads"), params), &resp)
	return &resp, err
}

func (c *Client) GetUpload(ctx context.Context, id string) (*Upload, error) {
	var upload Upload
	err := c.GetJSON(ctx, c.companyURL("/uploads/"+id), &upload)
	return &upload, err
}

func (c *Client) CreateUpload(ctx context.Context, filePath, description, journalEntryID string) (*Upload, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, err
	}

	if description != "" {
		writer.WriteField("description", description)
	}
	if journalEntryID != "" {
		writer.WriteField("journalEntryId", journalEntryID)
	}

	writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.companyURL("/uploads"), &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.PostRaw(ctx, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var upload Upload
	if err := json.NewDecoder(resp.Body).Decode(&upload); err != nil {
		return nil, err
	}
	return &upload, nil
}

func (c *Client) DownloadUpload(ctx context.Context, id string, w io.Writer) error {
	resp, err := c.Get(ctx, c.companyURL("/uploads/"+id+"/download"))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	return err
}
