package commands

import (
	"fmt"

	"github.com/domo84/bokio-cli/internal/api"
	"github.com/spf13/cobra"
)

func newBankPaymentsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bank-payments",
		Short: "Manage bank payments",
	}

	cmd.AddCommand(newBankPaymentsCreateCmd())
	cmd.AddCommand(newBankPaymentsListCmd())
	cmd.AddCommand(newBankPaymentsGetCmd())

	return cmd
}

func newBankPaymentsCreateCmd() *cobra.Command {
	var fromFile string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a bank payment",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}

			var req api.CreateBankPaymentRequest
			if err := loadJSONFile(fromFile, &req); err != nil {
				return fmt.Errorf("reading file: %w", err)
			}

			payment, err := state.client.CreateBankPayment(cmd.Context(), req)
			if err != nil {
				return err
			}
			return state.formatter.Format(payment)
		},
	}

	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load request from JSON file (required)")

	return cmd
}

func newBankPaymentsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List bank payments",
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
				payments, err := api.AutoPaginate(cmd.Context(), func(p api.ListParams) (*api.PaginatedResponse[api.BankPayment], error) {
					return state.client.ListBankPayments(cmd.Context(), p)
				}, params)
				if err != nil {
					return err
				}
				return state.formatter.Format(payments)
			}

			resp, err := state.client.ListBankPayments(cmd.Context(), params)
			if err != nil {
				return err
			}
			return state.formatter.Format(resp.Items)
		},
	}
	addListFlags(cmd)
	return cmd
}

func newBankPaymentsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a bank payment by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			payment, err := state.client.GetBankPayment(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(payment)
		},
	}
}
