package commands

import (
	"fmt"

	"github.com/domo84/bokio-cli/internal/api"
	"github.com/spf13/cobra"
)

func newItemsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "items",
		Short: "Manage items",
	}

	cmd.AddCommand(newItemsListCmd())
	cmd.AddCommand(newItemsGetCmd())
	cmd.AddCommand(newItemsCreateCmd())
	cmd.AddCommand(newItemsUpdateCmd())
	cmd.AddCommand(newItemsDeleteCmd())

	return cmd
}

func newItemsListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List items",
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
				items, err := api.AutoPaginate(cmd.Context(), func(p api.ListParams) (*api.PaginatedResponse[api.Item], error) {
					return state.client.ListItems(cmd.Context(), p)
				}, params)
				if err != nil {
					return err
				}
				return state.formatter.Format(items)
			}

			resp, err := state.client.ListItems(cmd.Context(), params)
			if err != nil {
				return err
			}
			return state.formatter.Format(resp.Items)
		},
	}
	addListFlags(cmd)
	return cmd
}

func newItemsGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an item by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			item, err := state.client.GetItem(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(item)
		},
	}
}

func newItemsCreateCmd() *cobra.Command {
	var (
		description string
		unitPrice   float64
		unit        string
		vatRate     float64
		fromFile    string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an item",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			var req api.CreateItemRequest
			if fromFile != "" {
				if err := loadJSONFile(fromFile, &req); err != nil {
					return fmt.Errorf("reading file: %w", err)
				}
			}

			if description != "" {
				req.Description = description
			}
			if unitPrice != 0 {
				req.UnitPrice = unitPrice
			}
			if unit != "" {
				req.Unit = unit
			}
			if vatRate != 0 {
				req.VatRate = vatRate
			}

			if req.Description == "" {
				return fmt.Errorf("--description is required")
			}

			item, err := state.client.CreateItem(cmd.Context(), req)
			if err != nil {
				return err
			}
			return state.formatter.Format(item)
		},
	}

	cmd.Flags().StringVar(&description, "description", "", "Item description (required)")
	cmd.Flags().Float64Var(&unitPrice, "unit-price", 0, "Unit price")
	cmd.Flags().StringVar(&unit, "unit", "", "Unit (e.g., pcs, hours)")
	cmd.Flags().Float64Var(&vatRate, "vat-rate", 0, "VAT rate")
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load request from JSON file")

	return cmd
}

func newItemsUpdateCmd() *cobra.Command {
	var fromFile string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			var req api.UpdateItemRequest
			if fromFile == "" {
				return fmt.Errorf("--from-file is required for update")
			}
			if err := loadJSONFile(fromFile, &req); err != nil {
				return fmt.Errorf("reading file: %w", err)
			}

			item, err := state.client.UpdateItem(cmd.Context(), args[0], req)
			if err != nil {
				return err
			}
			return state.formatter.Format(item)
		},
	}

	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load request from JSON file")

	return cmd
}

func newItemsDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete an item",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			if err := state.client.DeleteItem(cmd.Context(), args[0]); err != nil {
				return err
			}
			fmt.Println("Item deleted.")
			return nil
		},
	}
}
