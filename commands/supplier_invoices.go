package commands

import (
	"fmt"

	"github.com/domo84/bokio-cli/internal/api"
	"github.com/spf13/cobra"
)

func newSupplierInvoicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "supplier-invoices",
		Short: "Manage supplier invoices (leverantörsfakturor)",
	}

	cmd.AddCommand(newSupplierInvoicesListCmd())
	cmd.AddCommand(newSupplierInvoicesGetCmd())
	cmd.AddCommand(newSupplierInvoicesCreateCmd())
	cmd.AddCommand(newSupplierInvoicesUpdateCmd())
	cmd.AddCommand(newSupplierInvoicesDeleteCmd())

	lineItems := &cobra.Command{Use: "line-items", Short: "Manage supplier invoice line items"}
	lineItems.AddCommand(newSupplierInvoiceLineItemAddCmd())
	cmd.AddCommand(lineItems)

	uploads := &cobra.Command{Use: "uploads", Short: "Manage supplier invoice uploads"}
	uploads.AddCommand(newSupplierInvoiceUploadAddCmd())
	cmd.AddCommand(uploads)

	return cmd
}

func newSupplierInvoicesListCmd() *cobra.Command {
	var projectID string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List supplier invoices",
		Long: `List supplier invoices.

Filter fields for --query: supplierRef.id, invoiceNumber, totalAmount, currency,
invoiceDate, dueDate.`,
		Example: `  bokio supplier-invoices list -q "dueDate>=2025-01-01"
  bokio supplier-invoices list --project-id <id>`,
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
				invoices, err := api.AutoPaginate(cmd.Context(), func(p api.ListParams) (*api.PaginatedResponse[api.SupplierInvoice], error) {
					return state.client.ListSupplierInvoices(cmd.Context(), p, projectID)
				}, params)
				if err != nil {
					return err
				}
				return state.formatter.Format(invoices)
			}

			resp, err := state.client.ListSupplierInvoices(cmd.Context(), params, projectID)
			if err != nil {
				return err
			}
			return state.formatter.Format(resp.Items)
		},
	}
	addListFlags(cmd)
	cmd.Flags().StringVar(&projectID, "project-id", "", "Only invoices with a line item on this project")
	return cmd
}

func newSupplierInvoicesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a supplier invoice by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			inv, err := state.client.GetSupplierInvoice(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(inv)
		},
	}
}

func newSupplierInvoicesCreateCmd() *cobra.Command {
	var (
		supplierID    string
		invoiceNumber string
		invoiceDate   string
		dueDate       string
		totalAmount   float64
		uploadID      string
		fromFile      string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a supplier invoice",
		Long: `Create a supplier invoice. --from-file is loaded first, then flags are applied
on top. Line items can only be added once an upload is connected, so pass
--upload-id (or add one later with "supplier-invoices uploads add").`,
		Example: `  bokio supplier-invoices create --supplier-id <id> --invoice-number INV-001 \
    --invoice-date 2025-01-15 --due-date 2025-02-15 --total-amount 1250 --upload-id <uploadId>
  bokio supplier-invoices create --from-file supplier-invoice.json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			var req api.CreateSupplierInvoiceRequest
			if fromFile != "" {
				if err := loadJSONFile(fromFile, &req); err != nil {
					return fmt.Errorf("reading file: %w", err)
				}
			}

			if supplierID != "" {
				req.SupplierRef = api.Ref{ID: supplierID}
			}
			if invoiceNumber != "" {
				req.InvoiceNumber = invoiceNumber
			}
			if invoiceDate != "" {
				req.InvoiceDate = invoiceDate
			}
			if dueDate != "" {
				req.DueDate = dueDate
			}
			if cmd.Flags().Changed("total-amount") {
				req.TotalAmount = totalAmount
			}
			if uploadID != "" {
				req.UploadRef = &api.Ref{ID: uploadID}
			}

			if req.SupplierRef.ID == "" {
				return fmt.Errorf("--supplier-id is required")
			}
			if req.InvoiceDate == "" || req.DueDate == "" {
				return fmt.Errorf("--invoice-date and --due-date are required")
			}

			inv, err := state.client.CreateSupplierInvoice(cmd.Context(), req)
			if err != nil {
				return err
			}
			return state.formatter.Format(inv)
		},
	}

	cmd.Flags().StringVar(&supplierID, "supplier-id", "", "Supplier ID (required)")
	cmd.Flags().StringVar(&invoiceNumber, "invoice-number", "", "The supplier's invoice number")
	cmd.Flags().StringVar(&invoiceDate, "invoice-date", "", "Invoice date (YYYY-MM-DD, required)")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date (YYYY-MM-DD, required)")
	cmd.Flags().Float64Var(&totalAmount, "total-amount", 0, "Total amount including VAT (required)")
	cmd.Flags().StringVar(&uploadID, "upload-id", "", "ID of an existing upload to connect")
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load request from JSON file")

	return cmd
}

func newSupplierInvoicesUpdateCmd() *cobra.Command {
	var (
		supplierID    string
		invoiceNumber string
		invoiceDate   string
		dueDate       string
		totalAmount   float64
		fromFile      string
	)

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a supplier invoice",
		Long: `Update a supplier invoice. The current invoice is fetched first and the changes
are applied on top, so fields you don't set are kept. --from-file is merged over
the current invoice, then flags are applied over that.

A lineItems array in --from-file replaces the line items: include "id" to keep
or change an existing line item; line items left out are removed.

Only supplier invoices created through the API by the same integration can be
updated.`,
		Example: `  bokio supplier-invoices update <id> --due-date 2025-03-01
  bokio supplier-invoices update <id> --from-file supplier-invoice.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			current, err := state.client.GetSupplierInvoice(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			req := current.UpdateRequest()

			if fromFile != "" {
				if err := loadJSONFile(fromFile, &req); err != nil {
					return fmt.Errorf("reading file: %w", err)
				}
			}

			if supplierID != "" {
				req.SupplierRef = api.Ref{ID: supplierID}
			}
			if invoiceNumber != "" {
				req.InvoiceNumber = invoiceNumber
			}
			if invoiceDate != "" {
				req.InvoiceDate = invoiceDate
			}
			if dueDate != "" {
				req.DueDate = dueDate
			}
			if cmd.Flags().Changed("total-amount") {
				req.TotalAmount = totalAmount
			}

			inv, err := state.client.UpdateSupplierInvoice(cmd.Context(), args[0], req)
			if err != nil {
				return err
			}
			return state.formatter.Format(inv)
		},
	}

	cmd.Flags().StringVar(&supplierID, "supplier-id", "", "Supplier ID")
	cmd.Flags().StringVar(&invoiceNumber, "invoice-number", "", "The supplier's invoice number")
	cmd.Flags().StringVar(&invoiceDate, "invoice-date", "", "Invoice date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date (YYYY-MM-DD)")
	cmd.Flags().Float64Var(&totalAmount, "total-amount", 0, "Total amount including VAT")
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Merge a JSON supplier invoice body over the current invoice")

	return cmd
}

func newSupplierInvoicesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a supplier invoice",
		Long: `Delete a supplier invoice. Only supplier invoices created through the API by
the same integration can be deleted.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			if err := state.client.DeleteSupplierInvoice(cmd.Context(), args[0]); err != nil {
				return err
			}
			fmt.Println("Supplier invoice deleted.")
			return nil
		},
	}
}

func newSupplierInvoiceLineItemAddCmd() *cobra.Command {
	var (
		description  string
		quantity     float64
		unitPrice    float64
		unitType     string
		taxRate      float64
		projectID    string
		costCenterID string
		fromFile     string
	)

	cmd := &cobra.Command{
		Use:   "add <supplierInvoiceId>",
		Short: "Add a line item to a supplier invoice",
		Long: `Add a line item to a supplier invoice. The invoice must have an upload
connected first.`,
		Example: `  bokio supplier-invoices line-items add <id> --description Consulting \
    --quantity 10 --unit-price 100 --unit-type hour --tax-rate 25`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			var item api.SupplierInvoiceLineItem
			if fromFile != "" {
				if err := loadJSONFile(fromFile, &item); err != nil {
					return fmt.Errorf("reading file: %w", err)
				}
			}

			if description != "" {
				item.Description = description
			}
			if cmd.Flags().Changed("quantity") {
				item.Quantity = &quantity
			}
			if cmd.Flags().Changed("unit-price") {
				item.UnitPrice = unitPrice
			}
			if unitType != "" {
				item.UnitType = unitType
			}
			if cmd.Flags().Changed("tax-rate") {
				item.TaxRate = &taxRate
			}
			if projectID != "" {
				item.ProjectRef = &api.ProjectRef{ID: projectID}
			}
			if costCenterID != "" {
				item.CostCenterRef = &api.CostCenterRef{TagID: costCenterID}
			}

			if item.Description == "" {
				return fmt.Errorf("--description is required")
			}

			li, err := state.client.AddSupplierInvoiceLineItem(cmd.Context(), args[0], item)
			if err != nil {
				return err
			}
			return state.formatter.Format(li)
		},
	}

	cmd.Flags().StringVar(&description, "description", "", "Line item description (required)")
	cmd.Flags().Float64Var(&quantity, "quantity", 1, "Quantity")
	cmd.Flags().Float64Var(&unitPrice, "unit-price", 0, "Unit price (required)")
	cmd.Flags().StringVar(&unitType, "unit-type", "", "Unit, e.g. hour or piece")
	cmd.Flags().Float64Var(&taxRate, "tax-rate", 0, "VAT rate in percent, e.g. 25")
	cmd.Flags().StringVar(&projectID, "project-id", "", "Project ID")
	cmd.Flags().StringVar(&costCenterID, "cost-center-id", "", "Cost center tag value ID (dimension 1)")
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load from JSON file")

	return cmd
}

func newSupplierInvoiceUploadAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <supplierInvoiceId> <uploadId>",
		Short: "Connect an existing upload to a supplier invoice",
		Long: `Connect an existing upload to a supplier invoice. Create the upload first with
"bokio uploads create".`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			if _, err := state.client.AddSupplierInvoiceUpload(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			fmt.Println("Upload added.")
			return nil
		},
	}
}
