package commands

import (
	"fmt"

	"github.com/domo84/bokio-cli/internal/api"
	"github.com/spf13/cobra"
)

func newCreditNotesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credit-notes",
		Short: "Manage credit notes",
	}

	cmd.AddCommand(newCreditNotesListCmd())
	cmd.AddCommand(newCreditNotesGetCmd())
	cmd.AddCommand(newCreditNotesUpdateCmd())
	cmd.AddCommand(newCreditNotesPublishCmd())
	cmd.AddCommand(newCreditNotesRecordCmd())

	return cmd
}

func newCreditNotesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List credit notes",
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
				notes, err := api.AutoPaginate(cmd.Context(), func(p api.ListParams) (*api.PaginatedResponse[api.CreditNote], error) {
					return state.client.ListCreditNotes(cmd.Context(), p)
				}, params)
				if err != nil {
					return err
				}
				return state.formatter.Format(notes)
			}

			resp, err := state.client.ListCreditNotes(cmd.Context(), params)
			if err != nil {
				return err
			}
			return state.formatter.Format(resp.Items)
		},
	}
	addListFlags(cmd)
	return cmd
}

func newCreditNotesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a credit note by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			note, err := state.client.GetCreditNote(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(note)
		},
	}
}

func newCreditNotesUpdateCmd() *cobra.Command {
	var fromFile string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a credit note",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			var req api.UpdateCreditNoteRequest
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			if err := loadJSONFile(fromFile, &req); err != nil {
				return fmt.Errorf("reading file: %w", err)
			}

			note, err := state.client.UpdateCreditNote(cmd.Context(), args[0], req)
			if err != nil {
				return err
			}
			return state.formatter.Format(note)
		},
	}

	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load request from JSON file")

	return cmd
}

func newCreditNotesPublishCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "publish <id>",
		Short: "Publish a credit note",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			note, err := state.client.PublishCreditNote(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(note)
		},
	}
}

func newCreditNotesRecordCmd() *cobra.Command {
	var (
		accountNumber int
		paymentDate   string
	)

	cmd := &cobra.Command{
		Use:   "record <id>",
		Short: "Record a credit note",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			req := api.RecordCreditNoteRequest{
				PaymentAccountNumber: accountNumber,
				PaymentDate:          paymentDate,
			}

			if err := state.client.RecordCreditNote(cmd.Context(), args[0], req); err != nil {
				return err
			}
			fmt.Println("Credit note recorded.")
			return nil
		},
	}

	cmd.Flags().IntVar(&accountNumber, "account-number", 0, "Payment account number")
	cmd.Flags().StringVar(&paymentDate, "payment-date", "", "Payment date (YYYY-MM-DD)")

	return cmd
}
