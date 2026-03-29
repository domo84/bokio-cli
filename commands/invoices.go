package commands

import (
	"fmt"

	"github.com/domo84/bokio-cli/internal/api"
	"github.com/spf13/cobra"
)

func newInvoicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoices",
		Short: "Manage invoices",
	}

	cmd.AddCommand(newInvoicesListCmd())
	cmd.AddCommand(newInvoicesGetCmd())
	cmd.AddCommand(newInvoicesCreateCmd())
	cmd.AddCommand(newInvoicesUpdateCmd())
	cmd.AddCommand(newInvoicesPublishCmd())
	cmd.AddCommand(newInvoicesRecordCmd())

	// Sub-resource commands
	lineItems := &cobra.Command{Use: "line-items", Short: "Manage invoice line items"}
	lineItems.AddCommand(newInvoiceLineItemAddCmd())
	lineItems.AddCommand(newInvoiceLineItemUpdateCmd())
	lineItems.AddCommand(newInvoiceLineItemDeleteCmd())
	cmd.AddCommand(lineItems)

	attachments := &cobra.Command{Use: "attachments", Short: "Manage invoice attachments"}
	attachments.AddCommand(newInvoiceAttachmentsListCmd())
	attachments.AddCommand(newInvoiceAttachmentAddCmd())
	attachments.AddCommand(newInvoiceAttachmentDeleteCmd())
	cmd.AddCommand(attachments)

	payments := &cobra.Command{Use: "payments", Short: "Manage invoice payments"}
	payments.AddCommand(newInvoicePaymentsListCmd())
	payments.AddCommand(newInvoicePaymentAddCmd())
	cmd.AddCommand(payments)

	settlements := &cobra.Command{Use: "settlements", Short: "Manage invoice settlements"}
	settlements.AddCommand(newInvoiceSettlementAddCmd())
	cmd.AddCommand(settlements)

	return cmd
}

func newInvoicesListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List invoices",
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
				invoices, err := api.AutoPaginate(cmd.Context(), func(p api.ListParams) (*api.PaginatedResponse[api.Invoice], error) {
					return state.client.ListInvoices(cmd.Context(), p)
				}, params)
				if err != nil {
					return err
				}
				return state.formatter.Format(invoices)
			}

			resp, err := state.client.ListInvoices(cmd.Context(), params)
			if err != nil {
				return err
			}
			return state.formatter.Format(resp.Items)
		},
	}
	addListFlags(cmd)
	return cmd
}

func newInvoicesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get an invoice by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			inv, err := state.client.GetInvoice(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(inv)
		},
	}
}

func newInvoicesCreateCmd() *cobra.Command {
	var (
		customerID  string
		invoiceDate string
		dueDate     string
		currency    string
		fromFile    string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create an invoice",
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			var req api.CreateInvoiceRequest
			if fromFile != "" {
				if err := loadJSONFile(fromFile, &req); err != nil {
					return fmt.Errorf("reading file: %w", err)
				}
			}

			if customerID != "" {
				req.CustomerID = customerID
			}
			if invoiceDate != "" {
				req.InvoiceDate = invoiceDate
			}
			if dueDate != "" {
				req.DueDate = dueDate
			}
			if currency != "" {
				req.Currency = currency
			}

			if req.CustomerID == "" {
				return fmt.Errorf("--customer-id is required")
			}

			inv, err := state.client.CreateInvoice(cmd.Context(), req)
			if err != nil {
				return err
			}
			return state.formatter.Format(inv)
		},
	}

	cmd.Flags().StringVar(&customerID, "customer-id", "", "Customer ID (required)")
	cmd.Flags().StringVar(&invoiceDate, "invoice-date", "", "Invoice date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&dueDate, "due-date", "", "Due date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&currency, "currency", "", "Currency code")
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load request from JSON file")

	return cmd
}

func newInvoicesUpdateCmd() *cobra.Command {
	var fromFile string

	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update an invoice",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			var req api.UpdateInvoiceRequest
			if fromFile == "" {
				return fmt.Errorf("--from-file is required for update")
			}
			if err := loadJSONFile(fromFile, &req); err != nil {
				return fmt.Errorf("reading file: %w", err)
			}

			inv, err := state.client.UpdateInvoice(cmd.Context(), args[0], req)
			if err != nil {
				return err
			}
			return state.formatter.Format(inv)
		},
	}

	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load request from JSON file")

	return cmd
}

func newInvoicesPublishCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "publish <id>",
		Short: "Publish an invoice",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			inv, err := state.client.PublishInvoice(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(inv)
		},
	}
}

func newInvoicesRecordCmd() *cobra.Command {
	var (
		accountNumber int
		paymentDate   string
	)

	cmd := &cobra.Command{
		Use:   "record <id>",
		Short: "Record an invoice payment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			req := api.RecordInvoiceRequest{
				PaymentAccountNumber: accountNumber,
				PaymentDate:          paymentDate,
			}

			if err := state.client.RecordInvoice(cmd.Context(), args[0], req); err != nil {
				return err
			}
			fmt.Println("Invoice recorded.")
			return nil
		},
	}

	cmd.Flags().IntVar(&accountNumber, "account-number", 0, "Payment account number")
	cmd.Flags().StringVar(&paymentDate, "payment-date", "", "Payment date (YYYY-MM-DD)")

	return cmd
}

// Line Items sub-commands

func newInvoiceLineItemAddCmd() *cobra.Command {
	var (
		description   string
		quantity      float64
		unitPrice     float64
		vatRate       float64
		accountNumber int
		fromFile      string
	)

	cmd := &cobra.Command{
		Use:   "add <invoiceId>",
		Short: "Add a line item to an invoice",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			var item api.LineItem
			if fromFile != "" {
				if err := loadJSONFile(fromFile, &item); err != nil {
					return fmt.Errorf("reading file: %w", err)
				}
			}

			if description != "" {
				item.Description = description
			}
			if quantity != 0 {
				item.Quantity = &quantity
			}
			if unitPrice != 0 {
				item.UnitPrice = &unitPrice
			}
			if vatRate != 0 {
				item.TaxRate = &vatRate
			}
			if accountNumber != 0 {
				item.BookkeepingAccountNumber = &accountNumber
			}

			inv, err := state.client.AddInvoiceLineItem(cmd.Context(), args[0], item)
			if err != nil {
				return err
			}
			return state.formatter.Format(inv)
		},
	}

	cmd.Flags().StringVar(&description, "description", "", "Line item description")
	cmd.Flags().Float64Var(&quantity, "quantity", 0, "Quantity")
	cmd.Flags().Float64Var(&unitPrice, "unit-price", 0, "Unit price")
	cmd.Flags().Float64Var(&vatRate, "vat-rate", 0, "VAT rate")
	cmd.Flags().IntVar(&accountNumber, "account-number", 0, "Account number")
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load from JSON file")

	return cmd
}

func newInvoiceLineItemUpdateCmd() *cobra.Command {
	var fromFile string

	cmd := &cobra.Command{
		Use:   "update <invoiceId> <lineItemId>",
		Short: "Update a line item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			var item api.LineItem
			if fromFile == "" {
				return fmt.Errorf("--from-file is required")
			}
			if err := loadJSONFile(fromFile, &item); err != nil {
				return fmt.Errorf("reading file: %w", err)
			}

			inv, err := state.client.UpdateInvoiceLineItem(cmd.Context(), args[0], args[1], item)
			if err != nil {
				return err
			}
			return state.formatter.Format(inv)
		},
	}

	cmd.Flags().StringVar(&fromFile, "from-file", "", "Load from JSON file")

	return cmd
}

func newInvoiceLineItemDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <invoiceId> <lineItemId>",
		Short: "Delete a line item",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			if err := state.client.DeleteInvoiceLineItem(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			fmt.Println("Line item deleted.")
			return nil
		},
	}
}

// Attachments sub-commands

func newInvoiceAttachmentsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <invoiceId>",
		Short: "List invoice attachments",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			attachments, err := state.client.ListInvoiceAttachments(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(attachments)
		},
	}
}

func newInvoiceAttachmentAddCmd() *cobra.Command {
	var filePath string

	cmd := &cobra.Command{
		Use:   "add <invoiceId>",
		Short: "Add an attachment to an invoice",
		Args:  cobra.ExactArgs(1),
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

			// Use the uploads endpoint to upload, then attach
			upload, err := state.client.CreateUpload(cmd.Context(), filePath, "", "")
			if err != nil {
				return err
			}
			fmt.Printf("Attachment uploaded: %s\n", upload.ID)
			return nil
		},
	}

	cmd.Flags().StringVar(&filePath, "file", "", "File to attach")

	return cmd
}

func newInvoiceAttachmentDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <invoiceId> <attachmentId>",
		Short: "Delete an invoice attachment",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			if err := state.client.DeleteInvoiceAttachment(cmd.Context(), args[0], args[1]); err != nil {
				return err
			}
			fmt.Println("Attachment deleted.")
			return nil
		},
	}
}

// Payments sub-commands

func newInvoicePaymentsListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <invoiceId>",
		Short: "List invoice payments",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}
			payments, err := state.client.ListInvoicePayments(cmd.Context(), args[0])
			if err != nil {
				return err
			}
			return state.formatter.Format(payments)
		},
	}
}

func newInvoicePaymentAddCmd() *cobra.Command {
	var (
		amount        float64
		paymentDate   string
		accountNumber int
	)

	cmd := &cobra.Command{
		Use:   "add <invoiceId>",
		Short: "Add a payment to an invoice",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			req := api.CreatePaymentRequest{
				Amount:        amount,
				PaymentDate:   paymentDate,
				AccountNumber: accountNumber,
			}

			payment, err := state.client.AddInvoicePayment(cmd.Context(), args[0], req)
			if err != nil {
				return err
			}
			return state.formatter.Format([]api.InvoicePayment{*payment})
		},
	}

	cmd.Flags().Float64Var(&amount, "amount", 0, "Payment amount")
	cmd.Flags().StringVar(&paymentDate, "payment-date", "", "Payment date (YYYY-MM-DD)")
	cmd.Flags().IntVar(&accountNumber, "account-number", 0, "Account number")

	return cmd
}

// Settlements sub-commands

func newInvoiceSettlementAddCmd() *cobra.Command {
	var (
		amount        float64
		settleDate    string
		accountNumber int
	)

	cmd := &cobra.Command{
		Use:   "add <invoiceId>",
		Short: "Add a settlement to an invoice",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			state, err := initStateWithClient(cmd)
			if err != nil {
				return err
			}
			if err := requireCompanyID(state); err != nil {
				return err
			}

			req := api.CreateSettlementRequest{
				Amount:        amount,
				SettleDate:    settleDate,
				AccountNumber: accountNumber,
			}

			if err := state.client.AddInvoiceSettlement(cmd.Context(), args[0], req); err != nil {
				return err
			}
			fmt.Println("Settlement added.")
			return nil
		},
	}

	cmd.Flags().Float64Var(&amount, "amount", 0, "Settlement amount")
	cmd.Flags().StringVar(&settleDate, "settle-date", "", "Settlement date (YYYY-MM-DD)")
	cmd.Flags().IntVar(&accountNumber, "account-number", 0, "Account number")

	return cmd
}

