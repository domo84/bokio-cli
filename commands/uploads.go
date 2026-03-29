package commands

import (
	"fmt"
	"os"

	"github.com/domo84/bokio-cli/internal/api"
	"github.com/spf13/cobra"
)

func newUploadsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uploads",
		Short: "Manage file uploads",
	}

	cmd.AddCommand(newUploadsCreateCmd())
	cmd.AddCommand(newUploadsListCmd())
	cmd.AddCommand(newUploadsGetCmd())
	cmd.AddCommand(newUploadsDownloadCmd())

	return cmd
}

func newUploadsCreateCmd() *cobra.Command {
	var (
		filePath       string
		description    string
		journalEntryID string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Upload a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			if filePath == "" {
				return fmt.Errorf("--file is required")
			}

			upload, err := state.client.CreateUpload(cmd.Context(), filePath, description, journalEntryID)
			if err != nil {
				return err
			}
			return state.formatter.Format(upload)
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "File to upload (required)")
	cmd.Flags().StringVar(&description, "description", "", "File description")
	cmd.Flags().StringVar(&journalEntryID, "journal-entry-id", "", "Associated journal entry ID")

	return cmd
}

func newUploadsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List uploads",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			params := listParamsFromFlags()

			if flagAll {
				uploads, err := api.AutoPaginate(cmd.Context(), func(p api.ListParams) (*api.PaginatedResponse[api.Upload], error) {
					return state.client.ListUploads(cmd.Context(), p)
				}, params)
				if err != nil {
					return err
				}
				return state.formatter.Format(uploads)
			}

			resp, err := state.client.ListUploads(cmd.Context(), params)
			if err != nil {
				return err
			}
			return state.formatter.Format(resp.Items)
		},
	}
	addListFlags(cmd)
	return cmd
}

func newUploadsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get upload details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			upload, err := state.client.GetUpload(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(upload)
		},
	}
}

func newUploadsDownloadCmd() *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "download <id>",
		Short: "Download an uploaded file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			var w *os.File
			if outputPath != "" {
				w, err = os.Create(outputPath)
				if err != nil {
					return fmt.Errorf("creating output file: %w", err)
				}
				defer w.Close()
			} else {
				w = os.Stdout
			}

			return state.client.DownloadUpload(cmd.Context(), args[0], w)
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output-file", "O", "", "Output file path (default: stdout)")

	return cmd
}
