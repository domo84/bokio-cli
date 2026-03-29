package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/domo84/bokio-cli/internal/api"
)

const (
	AuthorizeURL = "https://api.bokio.se/v1/authorize"
	TokenURL     = "https://api.bokio.se/v1/token"
)

// OAuthFlow handles the OAuth 2.0 authorization code flow with PKCE.
type OAuthFlow struct {
	ClientID     string
	ClientSecret string
	RedirectPort int
}

// GeneratePKCE creates a code verifier and code challenge.
func GeneratePKCE() (verifier, challenge string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return
	}
	verifier = base64.RawURLEncoding.EncodeToString(b)
	h := sha256.Sum256([]byte(verifier))
	challenge = base64.RawURLEncoding.EncodeToString(h[:])
	return
}

// GenerateState creates a random state parameter.
func GenerateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// GetAuthURL returns the authorization URL for the browser.
func (f *OAuthFlow) GetAuthURL(state, codeChallenge string) string {
	params := url.Values{
		"client_id":             {f.ClientID},
		"response_type":        {"code"},
		"redirect_uri":         {fmt.Sprintf("http://localhost:%d/callback", f.RedirectPort)},
		"state":                {state},
		"code_challenge":       {codeChallenge},
		"code_challenge_method": {"S256"},
	}
	return AuthorizeURL + "?" + params.Encode()
}

// ExchangeCode exchanges an authorization code for tokens.
func (f *OAuthFlow) ExchangeCode(ctx context.Context, code, codeVerifier string) (*api.TokenResponse, error) {
	data := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {fmt.Sprintf("http://localhost:%d/callback", f.RedirectPort)},
		"code_verifier": {codeVerifier},
	}
	return f.tokenRequest(ctx, data)
}

// RefreshAccessToken uses a refresh token to get a new access token.
func (f *OAuthFlow) RefreshAccessToken(ctx context.Context, refreshToken string) (*api.TokenResponse, error) {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refreshToken},
	}
	return f.tokenRequest(ctx, data)
}

// ClientCredentials obtains a token using client credentials grant.
func (f *OAuthFlow) ClientCredentials(ctx context.Context) (*api.TokenResponse, error) {
	data := url.Values{
		"grant_type": {"client_credentials"},
	}
	return f.tokenRequest(ctx, data)
}

func (f *OAuthFlow) tokenRequest(ctx context.Context, data url.Values) (*api.TokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(f.ClientID, f.ClientSecret)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token request failed with status %d", resp.StatusCode)
	}

	var tokenResp api.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}
	return &tokenResp, nil
}

// GetToken resolves the current access token from env var or stored credentials.
func GetToken(store *TokenStore) (string, error) {
	if token := os.Getenv("BOKIO_TOKEN"); token != "" {
		return token, nil
	}

	creds, err := store.Load()
	if err != nil {
		return "", fmt.Errorf("not authenticated. Run 'bokio auth login' first")
	}

	if creds.AccessToken != "" {
		return creds.AccessToken, nil
	}

	return "", fmt.Errorf("no valid token found. Run 'bokio auth login' first")
}
