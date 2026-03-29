package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newSIECmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sie",
		Short: "SIE file operations",
	}

	cmd.AddCommand(newSIEDownloadCmd())
	return cmd
}

func newSIEDownloadCmd() *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "download <fiscalYearId>",
		Short: "Download SIE export file",
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

			return state.client.DownloadSIE(cmd.Context(), args[0], w)
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output-file", "O", "", "Output file path (default: stdout)")

	return cmd
}
