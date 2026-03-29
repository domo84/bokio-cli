package commands

import (
	"fmt"

	"github.com/domo84/bokio-cli/internal/auth"
	"github.com/spf13/cobra"
)

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
	}

	cmd.AddCommand(newAuthLoginCmd())
	cmd.AddCommand(newAuthLogoutCmd())
	cmd.AddCommand(newAuthStatusCmd())
	cmd.AddCommand(newAuthTokenCmd())

	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	var token string
	var useOAuth bool

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Bokio",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initState(cmd)
			if err != nil {
				return err
			}

			if token != "" {
				// Private token mode
				creds := &auth.Credentials{
					AccessToken: token,
					TokenType:   "bearer",
				}
				if err := state.store.Save(creds); err != nil {
					return fmt.Errorf("saving credentials: %w", err)
				}
				fmt.Println("Token saved successfully.")
				return nil
			}

			if useOAuth {
				if state.cfg.ClientID == "" || state.cfg.ClientSecret == "" {
					return fmt.Errorf("OAuth requires client_id and client_secret in config or environment")
				}

				flow := &auth.OAuthFlow{
					ClientID:     state.cfg.ClientID,
					ClientSecret: state.cfg.ClientSecret,
					RedirectPort: state.cfg.RedirectPort,
				}

				verifier, challenge, err := auth.GeneratePKCE()
				if err != nil {
					return fmt.Errorf("generating PKCE: %w", err)
				}

				oauthState, err := auth.GenerateState()
				if err != nil {
					return fmt.Errorf("generating state: %w", err)
				}

				resultCh, shutdown := auth.StartCallbackServer(state.cfg.RedirectPort)
				defer shutdown()

				authURL := flow.GetAuthURL(oauthState, challenge)
				fmt.Printf("Open this URL in your browser to authenticate:\n\n%s\n\nWaiting for callback...\n", authURL)

				result := <-resultCh
				if result.Error != "" {
					return fmt.Errorf("authentication failed: %s", result.Error)
				}
				if result.State != oauthState {
					return fmt.Errorf("state mismatch in OAuth callback")
				}

				tokenResp, err := flow.ExchangeCode(cmd.Context(), result.Code, verifier)
				if err != nil {
					return fmt.Errorf("exchanging code: %w", err)
				}

				creds := &auth.Credentials{
					AccessToken:  tokenResp.AccessToken,
					RefreshToken: tokenResp.RefreshToken,
					TokenType:    tokenResp.TokenType,
					ExpiresIn:    tokenResp.ExpiresIn,
					TenantID:     tokenResp.TenantID,
					ConnectionID: tokenResp.ConnectionID,
				}
				if err := state.store.Save(creds); err != nil {
					return fmt.Errorf("saving credentials: %w", err)
				}
				fmt.Println("Authentication successful.")
				return nil
			}

			return fmt.Errorf("provide --token for private token auth or --oauth for OAuth flow")
		},
	}

	cmd.Flags().StringVar(&token, "token", "", "Private API token")
	cmd.Flags().BoolVar(&useOAuth, "oauth", false, "Use OAuth 2.0 flow")

	return cmd
}

func newAuthLogoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Remove stored credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initState(cmd)
			if err != nil {
				return err
			}
			if err := state.store.Delete(); err != nil {
				return fmt.Errorf("removing credentials: %w", err)
			}
			fmt.Println("Logged out successfully.")
			return nil
		},
	}
}

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initState(cmd)
			if err != nil {
				return err
			}
			creds, err := state.store.Load()
			if err != nil {
				fmt.Println("Not authenticated.")
				return nil
			}
			fmt.Println("Authenticated.")
			if creds.TenantID != "" {
				fmt.Printf("Tenant ID: %s\n", creds.TenantID)
			}
			if creds.ConnectionID != "" {
				fmt.Printf("Connection ID: %s\n", creds.ConnectionID)
			}
			if creds.RefreshToken != "" {
				fmt.Println("Auth mode: OAuth")
			} else {
				fmt.Println("Auth mode: Private token")
			}
			return nil
		},
	}
}

func newAuthTokenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "token",
		Short: "Print current access token",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initState(cmd)
			if err != nil {
				return err
			}
			token, err := auth.GetToken(state.store)
			if err != nil {
				return err
			}
			fmt.Print(token)
			return nil
		},
	}
}
