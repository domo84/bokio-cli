package commands

import (
	"github.com/spf13/cobra"
)

func newFiscalYearsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fiscal-years",
		Short: "View fiscal years",
	}

	cmd.AddCommand(newFiscalYearsListCmd())
	cmd.AddCommand(newFiscalYearsGetCmd())

	return cmd
}

func newFiscalYearsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List fiscal years",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			resp, err := state.client.ListFiscalYears(cmd.Context())
			if err != nil {
				return err
			}
			return state.formatter.Format(resp.Items)
		},
	}
}

func newFiscalYearsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a fiscal year by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			year, err := state.client.GetFiscalYear(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(year)
		},
	}
}
