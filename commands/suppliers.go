package commands

import (
	"fmt"

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
	cmd.AddCommand(newSuppliersCreateCmd())
	cmd.AddCommand(newSuppliersUpdateCmd())

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

func newSuppliersUpdateCmd() *cobra.Command {
	var (
		name     string
		orgNum   string
		vatNum   string
		currency string
		bankgiro string
		plusgiro string
		fromFile string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a supplier",
		Long: `Update a supplier. The current supplier is fetched first and the changes are
applied on top, so fields you don't set are kept. --from-file is merged over
the current supplier, then flags are applied over that.

--bankgiro and --plusgiro replace the supplier's payment details.`,
		Example: `  bokio suppliers update <id> --bankgiro 5097-1282
  bokio suppliers update <id> --from-file supplier.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if bankgiro != "" && plusgiro != "" {
				return fmt.Errorf("--bankgiro and --plusgiro are mutually exclusive")
			}

			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			current, err := state.client.GetSupplier(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			req := current.UpdateRequest()

			if fromFile != "" {
				if err := loadJSONFile(fromFile, &req); err != nil {
					return fmt.Errorf("reading file: %w", err)
				}
			}

			if name != "" {
				req.Name = name
			}
			if orgNum != "" {
				req.OrgNumber = orgNum
			}
			if vatNum != "" {
				req.VatNumber = vatNum
			}
			if currency != "" {
				req.Currency = currency
			}
			if bankgiro != "" {
				req.PaymentDetails = &api.SupplierPaymentDetails{Type: "bankgiro", BankgiroNumber: bankgiro}
			}
			if plusgiro != "" {
				req.PaymentDetails = &api.SupplierPaymentDetails{Type: "plusgiro", PlusgiroNumber: plusgiro}
			}

			s, err := state.client.UpdateSupplier(cmd.Context(), args[0], req)
			if err != nil {
				return err
			}
			return state.formatter.Format(s)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Supplier name")
	cmd.Flags().StringVar(&orgNum, "org-number", "", "Organisation number")
	cmd.Flags().StringVar(&vatNum, "vat-number", "", "VAT number")
	cmd.Flags().StringVar(&currency, "currency", "", "ISO 4217 currency code")
	cmd.Flags().StringVar(&bankgiro, "bankgiro", "", "Bankgiro number, e.g. 5097-1282")
	cmd.Flags().StringVar(&plusgiro, "plusgiro", "", "Plusgiro number, e.g. 12345-6")
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Merge a JSON supplier body over the current supplier")

	return cmd
}

func newSuppliersCreateCmd() *cobra.Command {
	var (
		name     string
		orgNum   string
		vatNum   string
		currency string
		address  string
		city     string
		zipCode  string
		country  string
		bankgiro string
		plusgiro string
		fromFile string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a supplier",
		Example: `  bokio suppliers create --name "Acme AB" --org-number 556677-8899 --bankgiro 1234-5678
  bokio suppliers create --from-file supplier.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if bankgiro != "" && plusgiro != "" {
				return fmt.Errorf("--bankgiro and --plusgiro are mutually exclusive")
			}

			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			var req api.CreateSupplierRequest
			if fromFile != "" {
				if err := loadJSONFile(fromFile, &req); err != nil {
					return fmt.Errorf("reading file: %w", err)
				}
			}

			if name != "" {
				req.Name = name
			}
			if orgNum != "" {
				req.OrgNumber = orgNum
			}
			if vatNum != "" {
				req.VatNumber = vatNum
			}
			if currency != "" {
				req.Currency = currency
			}
			if address != "" || city != "" || zipCode != "" || country != "" {
				if req.Address == nil {
					req.Address = &api.SupplierAddress{}
				}
				if address != "" {
					req.Address.Line1 = address
				}
				if city != "" {
					req.Address.City = city
				}
				if zipCode != "" {
					req.Address.PostalCode = zipCode
				}
				if country != "" {
					req.Address.Country = country
				}
			}
			if bankgiro != "" {
				req.PaymentDetails = &api.SupplierPaymentDetails{Type: "bankgiro", BankgiroNumber: bankgiro}
			}
			if plusgiro != "" {
				req.PaymentDetails = &api.SupplierPaymentDetails{Type: "plusgiro", PlusgiroNumber: plusgiro}
			}

			if req.Name == "" {
				return fmt.Errorf("--name is required")
			}

			s, err := state.client.CreateSupplier(cmd.Context(), req)
			if err != nil {
				return err
			}
			return state.formatter.Format(s)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Supplier name (required)")
	cmd.Flags().StringVar(&orgNum, "org-number", "", "Organisation number")
	cmd.Flags().StringVar(&vatNum, "vat-number", "", "VAT number")
	cmd.Flags().StringVar(&currency, "currency", "", "ISO 4217 currency code (default SEK)")
	cmd.Flags().StringVar(&address, "address", "", "Street address")
	cmd.Flags().StringVar(&city, "city", "", "City")
	cmd.Flags().StringVar(&zipCode, "zip-code", "", "Postal code")
	cmd.Flags().StringVar(&country, "country", "", "ISO 3166-1 alpha-2 country code (default SE)")
	cmd.Flags().StringVar(&bankgiro, "bankgiro", "", "Bankgiro number, e.g. 5097-1282")
	cmd.Flags().StringVar(&plusgiro, "plusgiro", "", "Plusgiro number, e.g. 12345-6")
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load request from JSON file")

	return cmd
}
