package auth

import (
	"context"
	"fmt"
	"net/http"
)

// CallbackResult holds the OAuth callback response.
type CallbackResult struct {
	Code  string
	State string
	Error string
}

// StartCallbackServer starts a local HTTP server to receive the OAuth callback.
func StartCallbackServer(port int) (resultCh chan CallbackResult, shutdown func()) {
	ch := make(chan CallbackResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		result := CallbackResult{
			Code:  r.URL.Query().Get("code"),
			State: r.URL.Query().Get("state"),
			Error: r.URL.Query().Get("error"),
		}

		if result.Error != "" {
			fmt.Fprintf(w, "<html><body><h1>Authentication Failed</h1><p>%s</p><p>You can close this window.</p></body></html>", result.Error)
		} else {
			fmt.Fprint(w, "<html><body><h1>Authentication Successful</h1><p>You can close this window and return to the terminal.</p></body></html>")
		}

		ch <- result
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go srv.ListenAndServe()

	return ch, func() {
		srv.Shutdown(context.Background())
	}
}
