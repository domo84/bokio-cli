package api

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"testing"
)

// The API rejects the file part unless its Content-Type is one it accepts, so an
// octet-stream default fails with a 400. Guard the sniffing that prevents that.
func TestDetectUploadContentType(t *testing.T) {
	cases := []struct {
		name    string
		content []byte
		want    string
		wantErr bool
	}{
		{name: "pdf", content: []byte("%PDF-1.4\n1 0 obj\n"), want: "application/pdf"},
		{name: "png", content: []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR"), want: "image/png"},
		{name: "jpeg", content: []byte("\xff\xd8\xff\xe0\x00\x10JFIF"), want: "image/jpeg"},
		{name: "plain text is rejected", content: []byte("not a receipt"), wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "upload")
			if err := os.WriteFile(path, tc.content, 0o600); err != nil {
				t.Fatalf("write: %v", err)
			}
			f, err := os.Open(path)
			if err != nil {
				t.Fatalf("open: %v", err)
			}
			defer f.Close()

			got, err := detectUploadContentType(f)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("detect: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}

			// The body is streamed after sniffing, so the offset must be back at 0.
			offset, err := f.Seek(0, io.SeekCurrent)
			if err != nil {
				t.Fatalf("seek: %v", err)
			}
			if offset != 0 {
				t.Errorf("offset = %d, want 0", offset)
			}
		})
	}
}

// Payloads taken verbatim from the upload schema example in company-api.yaml.
func TestUploadDecodesSpecExample(t *testing.T) {
	const attached = `{
	  "id": "a419cf69-db6f-4de9-992c-b1a60942a443",
	  "description": "example.png",
	  "contentType": "img/png",
	  "journalEntryId": "835ba700-b306-4bd9-8447-59207b6b0002"
	}`

	var u Upload
	if err := json.Unmarshal([]byte(attached), &u); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if u.ID != "a419cf69-db6f-4de9-992c-b1a60942a443" {
		t.Errorf("ID = %q", u.ID)
	}
	if u.Description != "example.png" {
		t.Errorf("Description = %q", u.Description)
	}
	if u.ContentType != "img/png" {
		t.Errorf("ContentType = %q", u.ContentType)
	}
	if u.JournalEntryID == nil || *u.JournalEntryID != "835ba700-b306-4bd9-8447-59207b6b0002" {
		t.Errorf("JournalEntryID = %v", u.JournalEntryID)
	}

	// journalEntryId is nullable: an unbookkept upload must be distinguishable.
	var orphan Upload
	if err := json.Unmarshal([]byte(`{"id":"x","journalEntryId":null}`), &orphan); err != nil {
		t.Fatalf("unmarshal orphan: %v", err)
	}
	if orphan.JournalEntryID != nil {
		t.Errorf("expected nil JournalEntryID, got %v", *orphan.JournalEntryID)
	}
}

func TestUploadPagedEnvelope(t *testing.T) {
	const paged = `{"totalItems":2,"totalPages":1,"currentPage":1,"items":[
	  {"id":"a","description":"one.pdf","contentType":"application/pdf","journalEntryId":null},
	  {"id":"b","description":"two.png","contentType":"image/png","journalEntryId":"835ba700-b306-4bd9-8447-59207b6b0002"}
	]}`

	var resp PaginatedResponse[Upload]
	if err := json.Unmarshal([]byte(paged), &resp); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if resp.TotalItems != 2 || resp.TotalPages != 1 || resp.CurrentPage != 1 {
		t.Errorf("envelope = %+v", resp)
	}
	if len(resp.Items) != 2 {
		t.Fatalf("items = %d, want 2", len(resp.Items))
	}
	if resp.Items[0].JournalEntryID != nil {
		t.Errorf("item 0 should be unbookkept")
	}
	if resp.Items[1].JournalEntryID == nil {
		t.Errorf("item 1 should be linked")
	}
}
