package commands

import (
	"fmt"

	"github.com/domo84/bokio-cli/internal/api"
	"github.com/spf13/cobra"
)

func newCustomersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "customers",
		Short: "Manage customers",
	}

	cmd.AddCommand(newCustomersListCmd())
	cmd.AddCommand(newCustomersGetCmd())
	cmd.AddCommand(newCustomersCreateCmd())
	cmd.AddCommand(newCustomersUpdateCmd())
	cmd.AddCommand(newCustomersDeleteCmd())

	return cmd
}

func newCustomersListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List customers",
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
				customers, err := api.AutoPaginate(cmd.Context(), func(p api.ListParams) (*api.PaginatedResponse[api.Customer], error) {
					return state.client.ListCustomers(cmd.Context(), p)
				}, params)
				if err != nil {
					return err
				}
				return state.formatter.Format(customers)
			}

			resp, err := state.client.ListCustomers(cmd.Context(), params)
			if err != nil {
				return err
			}
			return state.formatter.Format(resp.Items)
		},
	}
	addListFlags(cmd)
	return cmd
}

func newCustomersGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a customer by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			cust, err := state.client.GetCustomer(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(cust)
		},
	}
}

func newCustomersCreateCmd() *cobra.Command {
	var (
		name     string
		email    string
		phone    string
		address  string
		city     string
		zipCode  string
		country  string
		orgNum   string
		fromFile string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a customer",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			var req api.CreateCustomerRequest
			if fromFile != "" {
				if err := loadJSONFile(fromFile, &req); err != nil {
					return fmt.Errorf("reading file: %w", err)
				}
			}

			// Flag overrides
			if name != "" {
				req.Name = name
			}
			if email != "" {
				req.Email = email
			}
			if phone != "" {
				req.Phone = phone
			}
			if address != "" {
				req.Address = address
			}
			if city != "" {
				req.City = city
			}
			if zipCode != "" {
				req.ZipCode = zipCode
			}
			if country != "" {
				req.Country = country
			}
			if orgNum != "" {
				req.OrganisationNumber = orgNum
			}

			if req.Name == "" {
				return fmt.Errorf("--name is required")
			}

			cust, err := state.client.CreateCustomer(cmd.Context(), req)
			if err != nil {
				return err
			}
			return state.formatter.Format(cust)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Customer name (required)")
	cmd.Flags().StringVar(&email, "email", "", "Email address")
	cmd.Flags().StringVar(&phone, "phone", "", "Phone number")
	cmd.Flags().StringVar(&address, "address", "", "Street address")
	cmd.Flags().StringVar(&city, "city", "", "City")
	cmd.Flags().StringVar(&zipCode, "zip-code", "", "Zip code")
	cmd.Flags().StringVar(&country, "country", "", "Country")
	cmd.Flags().StringVar(&orgNum, "org-number", "", "Organisation number")
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load request from JSON file")

	return cmd
}

func newCustomersUpdateCmd() *cobra.Command {
	var fromFile string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a customer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			var req api.UpdateCustomerRequest
			if fromFile == "" {
				return fmt.Errorf("--from-file is required for update")
			}
			if err := loadJSONFile(fromFile, &req); err != nil {
				return fmt.Errorf("reading file: %w", err)
			}

			cust, err := state.client.UpdateCustomer(cmd.Context(), args[0], req)
			if err != nil {
				return err
			}
			return state.formatter.Format(cust)
		},
	}

	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load request from JSON file")

	return cmd
}

func newCustomersDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a customer",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			if err := state.client.DeleteCustomer(cmd.Context(), args[0]); err != nil {
				return err
			}
			fmt.Println("Customer deleted.")
			return nil
		},
	}
}
