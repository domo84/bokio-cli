package commands

import (
	"github.com/domo84/bokio-cli/internal/api"
	"github.com/spf13/cobra"
)

func newSuppliersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "suppliers",
		Short: "Manage suppliers (leverantörer)",
	}

	cmd.AddCommand(newSuppliersListCmd())
	cmd.AddCommand(newSuppliersGetCmd())

	return cmd
}

func newSuppliersListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List suppliers",
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
				suppliers, err := api.AutoPaginate(cmd.Context(), func(p api.ListParams) (*api.PaginatedResponse[api.Supplier], error) {
					return state.client.ListSuppliers(cmd.Context(), p)
				}, params)
				if err != nil {
					return err
				}
				return state.formatter.Format(suppliers)
			}

			resp, err := state.client.ListSuppliers(cmd.Context(), params)
			if err != nil {
				return err
			}
			return state.formatter.Format(resp.Items)
		},
	}
	addListFlags(cmd)
	return cmd
}

func newSuppliersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a supplier by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			s, err := state.client.GetSupplier(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(s)
		},
	}
}
