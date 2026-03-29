package output

import (
	"fmt"
	"io"
	"strconv"

	"github.com/domo84/bokio-cli/internal/api"
	"github.com/olekukonko/tablewriter"
)

// TableFormatter outputs data as an ASCII table.
type TableFormatter struct {
	Writer io.Writer
}

func (f *TableFormatter) Format(data any) error {
	switch v := data.(type) {
	case *api.CompanyInfo:
		return f.formatCompanyInfo(v)
	case []api.Customer:
		return f.formatCustomers(v)
	case *api.Customer:
		return f.formatCustomers([]api.Customer{*v})
	case []api.Item:
		return f.formatItems(v)
	case *api.Item:
		return f.formatItems([]api.Item{*v})
	case []api.Invoice:
		return f.formatInvoices(v)
	case *api.Invoice:
		return f.formatInvoices([]api.Invoice{*v})
	case []api.JournalEntry:
		return f.formatJournalEntries(v)
	case *api.JournalEntry:
		return f.formatJournalEntries([]api.JournalEntry{*v})
	case []api.CreditNote:
		return f.formatCreditNotes(v)
	case *api.CreditNote:
		return f.formatCreditNotes([]api.CreditNote{*v})
	case []api.Upload:
		return f.formatUploads(v)
	case *api.Upload:
		return f.formatUploads([]api.Upload{*v})
	case []api.BankPayment:
		return f.formatBankPayments(v)
	case *api.BankPayment:
		return f.formatBankPayments([]api.BankPayment{*v})
	case []api.Account:
		return f.formatAccounts(v)
	case *api.Account:
		return f.formatAccounts([]api.Account{*v})
	case []api.FiscalYear:
		return f.formatFiscalYears(v)
	case *api.FiscalYear:
		return f.formatFiscalYears([]api.FiscalYear{*v})
	case []api.Connection:
		return f.formatConnections(v)
	case *api.Connection:
		return f.formatConnections([]api.Connection{*v})
	case []api.InvoicePayment:
		return f.formatInvoicePayments(v)
	case []api.InvoiceAttachment:
		return f.formatInvoiceAttachments(v)
	default:
		return fmt.Errorf("unsupported type for table output: %T", data)
	}
}

func newTable(w io.Writer) *tablewriter.Table {
	return tablewriter.NewWriter(w)
}

func (f *TableFormatter) formatCompanyInfo(c *api.CompanyInfo) error {
	t := newTable(f.Writer)
	t.Header("Field", "Value")
	t.Append("ID", c.ID)
	t.Append("Name", c.Name)
	t.Append("Type", c.CompanyType)
	t.Append("Org Number", c.OrganizationNumber)
	t.Append("Email", c.Email)
	t.Append("Phone", c.Phone)
	t.Append("Address", c.Address.Line1)
	if c.Address.Line2 != "" {
		t.Append("Address 2", c.Address.Line2)
	}
	t.Append("City", c.Address.City)
	t.Append("Postal Code", c.Address.PostalCode)
	t.Append("Country", c.Address.Country)
	return t.Render()
}

func (f *TableFormatter) formatCustomers(customers []api.Customer) error {
	t := newTable(f.Writer)
	t.Header("ID", "Name", "Email", "Customer #", "City")
	for _, c := range customers {
		t.Append(c.ID, c.Name, c.Email, c.CustomerNumber, c.City)
	}
	return t.Render()
}

func (f *TableFormatter) formatItems(items []api.Item) error {
	t := newTable(f.Writer)
	t.Header("ID", "Description", "Unit Price", "Unit", "VAT Rate", "Article #")
	for _, i := range items {
		t.Append(i.ID, i.Description, fmtFloat(i.UnitPrice), i.Unit, fmtFloat(i.VatRate), i.ArticleNumber)
	}
	return t.Render()
}

func (f *TableFormatter) formatInvoices(invoices []api.Invoice) error {
	t := newTable(f.Writer)
	t.Header("ID", "Invoice #", "Customer", "Date", "Due", "Status", "Total")
	for _, inv := range invoices {
		invNum := ""
		if inv.InvoiceNumber != nil {
			invNum = *inv.InvoiceNumber
		}
		custName := ""
		if inv.CustomerRef != nil {
			custName = inv.CustomerRef.Name
		}
		t.Append(inv.ID, invNum, custName, inv.InvoiceDate, inv.DueDate, inv.Status, fmtFloat(inv.TotalAmount))
	}
	return t.Render()
}

func (f *TableFormatter) formatJournalEntries(entries []api.JournalEntry) error {
	t := newTable(f.Writer)
	t.Header("ID", "Date", "Description", "Rows", "Reversed")
	for _, e := range entries {
		reversed := "No"
		if e.IsReversed {
			reversed = "Yes"
		}
		t.Append(e.ID, e.Date, e.Description, strconv.Itoa(len(e.Rows)), reversed)
	}
	return t.Render()
}

func (f *TableFormatter) formatCreditNotes(notes []api.CreditNote) error {
	t := newTable(f.Writer)
	t.Header("ID", "Credit Note #", "Invoice ID", "Status", "Total")
	for _, n := range notes {
		t.Append(n.ID, n.CreditNoteNumber, n.InvoiceID, n.Status, fmtFloat(n.TotalAmount))
	}
	return t.Render()
}

func (f *TableFormatter) formatUploads(uploads []api.Upload) error {
	t := newTable(f.Writer)
	t.Header("ID", "File Name", "Content Type", "Description", "Created")
	for _, u := range uploads {
		t.Append(u.ID, u.FileName, u.ContentType, u.Description, u.CreatedAt.Format("2006-01-02"))
	}
	return t.Render()
}

func (f *TableFormatter) formatBankPayments(payments []api.BankPayment) error {
	t := newTable(f.Writer)
	t.Header("ID", "Amount", "Currency", "Recipient", "Status")
	for _, p := range payments {
		t.Append(p.ID, fmtFloat(p.Amount), p.Currency, p.RecipientName, p.Status)
	}
	return t.Render()
}

func (f *TableFormatter) formatAccounts(accounts []api.Account) error {
	t := newTable(f.Writer)
	t.Header("Account #", "Name", "Type")
	for _, a := range accounts {
		t.Append(strconv.Itoa(a.AccountNumber), a.Name, a.AccountType)
	}
	return t.Render()
}

func (f *TableFormatter) formatFiscalYears(years []api.FiscalYear) error {
	t := newTable(f.Writer)
	t.Header("ID", "Start Date", "End Date", "Status", "Method")
	for _, y := range years {
		t.Append(y.ID, y.StartDate, y.EndDate, y.Status, y.AccountingMethod)
	}
	return t.Render()
}

func (f *TableFormatter) formatConnections(connections []api.Connection) error {
	t := newTable(f.Writer)
	t.Header("ID", "Tenant ID", "Status", "Created")
	for _, c := range connections {
		t.Append(c.ID, c.TenantID, c.Status, c.CreatedAt.Format("2006-01-02"))
	}
	return t.Render()
}

func (f *TableFormatter) formatInvoicePayments(payments []api.InvoicePayment) error {
	t := newTable(f.Writer)
	t.Header("ID", "Amount", "Payment Date", "Account #")
	for _, p := range payments {
		t.Append(p.ID, fmtFloat(p.Amount), p.PaymentDate, strconv.Itoa(p.AccountNumber))
	}
	return t.Render()
}

func (f *TableFormatter) formatInvoiceAttachments(attachments []api.InvoiceAttachment) error {
	t := newTable(f.Writer)
	t.Header("ID", "File Name", "Content Type")
	for _, a := range attachments {
		t.Append(a.ID, a.FileName, a.ContentType)
	}
	return t.Render()
}

func fmtFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
