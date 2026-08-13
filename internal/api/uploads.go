package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
)

// uploadContentTypes are the file types POST /uploads accepts.
var uploadContentTypes = map[string]bool{
	"application/pdf": true,
	"image/jpeg":      true,
	"image/png":       true,
}

// detectUploadContentType sniffs f's content type and rejects types the API does
// not accept, so an unsupported file fails locally instead of as a 400. The read
// offset is restored before returning.
func detectUploadContentType(f *os.File) (string, error) {
	head := make([]byte, 512)
	n, err := f.Read(head)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("reading file: %w", err)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("rewinding file: %w", err)
	}

	contentType := http.DetectContentType(head[:n])
	if !uploadContentTypes[contentType] {
		return "", fmt.Errorf("unsupported file type %s: only application/pdf, image/jpeg and image/png are allowed", contentType)
	}
	return contentType, nil
}

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

	contentType, err := detectUploadContentType(f)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// The API validates the Content-Type of the file part. CreateFormFile would
	// label it application/octet-stream, which Bokio rejects as an invalid file
	// type, so the part header is built explicitly.
	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="file"; filename=%q`, filepath.Base(filePath)))
	header.Set("Content-Type", contentType)

	part, err := writer.CreatePart(header)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, f); err != nil {
		return nil, err
	}

	// Field names are defined by the multipart/form-data schema for POST /uploads.
	if description != "" {
		if err := writer.WriteField("description", description); err != nil {
			return nil, err
		}
	}
	if journalEntryID != "" {
		if err := writer.WriteField("journalEntryId", journalEntryID); err != nil {
			return nil, err
		}
	}

	// Close writes the trailing boundary; without it the body is malformed.
	if err := writer.Close(); err != nil {
		return nil, err
	}

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
