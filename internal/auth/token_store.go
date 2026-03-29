package auth

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/domo84/bokio-cli/internal/config"
)

// Credentials stores authentication tokens.
type Credentials struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	TenantID     string `json:"tenant_id,omitempty"`
	ConnectionID string `json:"connection_id,omitempty"`
}

// TokenStore handles reading and writing credentials.
type TokenStore struct {
	Dir string
}

// NewTokenStore creates a token store in the given directory.
func NewTokenStore(dir string) *TokenStore {
	if dir == "" {
		dir = config.DefaultConfigDir()
	}
	return &TokenStore{Dir: dir}
}

func (s *TokenStore) path() string {
	return filepath.Join(s.Dir, "credentials.json")
}

// Load reads stored credentials.
func (s *TokenStore) Load() (*Credentials, error) {
	data, err := os.ReadFile(s.path())
	if err != nil {
		return nil, err
	}
	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	return &creds, nil
}

// Save writes credentials to disk.
func (s *TokenStore) Save(creds *Credentials) error {
	if err := os.MkdirAll(s.Dir, 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(), data, 0600)
}

// Delete removes stored credentials.
func (s *TokenStore) Delete() error {
	return os.Remove(s.path())
}
