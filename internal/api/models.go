package api

import "time"

// Address represents a structured address.
type Address struct {
	Line1      string `json:"line1,omitempty"`
	Line2      string `json:"line2,omitempty"`
	City       string `json:"city,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
	Country    string `json:"country,omitempty"`
}

// Company Information
type CompanyInfo struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	CompanyType          string  `json:"companyType,omitempty"`
	OrganizationNumber   string  `json:"organizationNumber,omitempty"`
	Email                string  `json:"email,omitempty"`
	Phone                string  `json:"phone,omitempty"`
	HasBBA               bool    `json:"hasBBA,omitempty"`
	Address              Address `json:"address,omitempty"`
}

// Customer
type Customer struct {
	ID                 string            `json:"id"`
	Name               string            `json:"name"`
	OrganisationNumber string            `json:"organisationNumber,omitempty"`
	VatNumber          string            `json:"vatNumber,omitempty"`
	Email              string            `json:"email,omitempty"`
	Phone              string            `json:"phone,omitempty"`
	Address            string            `json:"address,omitempty"`
	Address2           string            `json:"address2,omitempty"`
	City               string            `json:"city,omitempty"`
	ZipCode            string            `json:"zipCode,omitempty"`
	Country            string            `json:"country,omitempty"`
	CustomerNumber     string            `json:"customerNumber,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
	CreatedAt          time.Time         `json:"createdAt"`
	UpdatedAt          time.Time         `json:"updatedAt"`
}

type CreateCustomerRequest struct {
	Name               string            `json:"name"`
	OrganisationNumber string            `json:"organisationNumber,omitempty"`
	VatNumber          string            `json:"vatNumber,omitempty"`
	Email              string            `json:"email,omitempty"`
	Phone              string            `json:"phone,omitempty"`
	Address            string            `json:"address,omitempty"`
	Address2           string            `json:"address2,omitempty"`
	City               string            `json:"city,omitempty"`
	ZipCode            string            `json:"zipCode,omitempty"`
	Country            string            `json:"country,omitempty"`
	CustomerNumber     string            `json:"customerNumber,omitempty"`
	Metadata           map[string]string `json:"metadata,omitempty"`
}

type UpdateCustomerRequest = CreateCustomerRequest

// Item
type Item struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	UnitPrice   float64 `json:"unitPrice"`
	Unit        string  `json:"unit,omitempty"`
	VatRate     float64 `json:"vatRate,omitempty"`
	AccountNumber int   `json:"accountNumber,omitempty"`
	ArticleNumber string `json:"articleNumber,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type CreateItemRequest struct {
	Description   string  `json:"description"`
	UnitPrice     float64 `json:"unitPrice"`
	Unit          string  `json:"unit,omitempty"`
	VatRate       float64 `json:"vatRate,omitempty"`
	AccountNumber int     `json:"accountNumber,omitempty"`
	ArticleNumber string  `json:"articleNumber,omitempty"`
}

type UpdateItemRequest = CreateItemRequest

// Ref is a reference to another entity.
type Ref struct {
	ID    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}

// Invoice
type Invoice struct {
	ID                    string            `json:"id"`
	Type                  string            `json:"type,omitempty"`
	InvoiceNumber         *string           `json:"invoiceNumber,omitempty"`
	OrderNumberReference  *string           `json:"orderNumberReference,omitempty"`
	CustomerRef           *Ref              `json:"customerRef,omitempty"`
	ContactDetailRef      *Ref              `json:"contactDetailRef,omitempty"`
	Currency              string            `json:"currency,omitempty"`
	CurrencyRate          float64           `json:"currencyRate,omitempty"`
	TotalAmount           float64           `json:"totalAmount"`
	TotalTax              float64           `json:"totalTax"`
	PaidAmount            float64           `json:"paidAmount"`
	Status                string            `json:"status"`
	InvoiceDate           string            `json:"invoiceDate"`
	DueDate               string            `json:"dueDate"`
	PublishedDateTime     *string           `json:"publishedDateTime,omitempty"`
	LineItems             []LineItem        `json:"lineItems,omitempty"`
	Metadata              map[string]string `json:"metadata,omitempty"`
}

type LineItem struct {
	ID                       int      `json:"id,omitempty"`
	ItemRef                  *Ref     `json:"itemRef,omitempty"`
	Description              string   `json:"description"`
	ItemType                 string   `json:"itemType,omitempty"`
	ProductType              *string  `json:"productType,omitempty"`
	UnitType                 *string  `json:"unitType,omitempty"`
	Quantity                 *float64 `json:"quantity,omitempty"`
	UnitPrice                *float64 `json:"unitPrice,omitempty"`
	TaxRate                  *float64 `json:"taxRate,omitempty"`
	BookkeepingAccountNumber *int     `json:"bookkeepingAccountNumber,omitempty"`
}

type CreateInvoiceRequest struct {
	CustomerID  string            `json:"customerId"`
	InvoiceDate string            `json:"invoiceDate"`
	DueDate     string            `json:"dueDate"`
	Currency    string            `json:"currency,omitempty"`
	LineItems   []LineItem        `json:"lineItems,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type UpdateInvoiceRequest struct {
	CustomerID  string            `json:"customerId,omitempty"`
	InvoiceDate string            `json:"invoiceDate,omitempty"`
	DueDate     string            `json:"dueDate,omitempty"`
	Currency    string            `json:"currency,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type RecordInvoiceRequest struct {
	PaymentAccountNumber int    `json:"paymentAccountNumber"`
	PaymentDate          string `json:"paymentDate"`
}

type InvoiceAttachment struct {
	ID          string `json:"id"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
}

type InvoicePayment struct {
	ID            string  `json:"id"`
	Amount        float64 `json:"amount"`
	PaymentDate   string  `json:"paymentDate"`
	AccountNumber int     `json:"accountNumber"`
}

type CreatePaymentRequest struct {
	Amount        float64 `json:"amount"`
	PaymentDate   string  `json:"paymentDate"`
	AccountNumber int     `json:"accountNumber"`
}

type CreateSettlementRequest struct {
	Amount        float64 `json:"amount"`
	SettleDate    string  `json:"settleDate"`
	AccountNumber int     `json:"accountNumber"`
}

// Journal Entry
type JournalEntry struct {
	ID          string             `json:"id"`
	Date        string             `json:"date"`
	Description string             `json:"description,omitempty"`
	Rows        []JournalEntryRow  `json:"rows"`
	IsReversed  bool               `json:"isReversed"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

type JournalEntryRow struct {
	AccountNumber int     `json:"accountNumber"`
	DebitAmount   float64 `json:"debitAmount"`
	CreditAmount  float64 `json:"creditAmount"`
	Description   string  `json:"description,omitempty"`
}

type CreateJournalEntryRequest struct {
	Date        string            `json:"date"`
	Description string            `json:"description,omitempty"`
	Rows        []JournalEntryRow `json:"rows"`
}

// Credit Note
type CreditNote struct {
	ID             string    `json:"id"`
	CreditNoteNumber string  `json:"creditNoteNumber,omitempty"`
	InvoiceID      string    `json:"invoiceId,omitempty"`
	Status         string    `json:"status"`
	TotalAmount    float64   `json:"totalAmount"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type UpdateCreditNoteRequest struct {
	Metadata map[string]string `json:"metadata,omitempty"`
}

type RecordCreditNoteRequest struct {
	PaymentAccountNumber int    `json:"paymentAccountNumber"`
	PaymentDate          string `json:"paymentDate"`
}

// Upload
type Upload struct {
	ID          string    `json:"id"`
	FileName    string    `json:"fileName"`
	ContentType string    `json:"contentType"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Bank Payment
type BankPayment struct {
	ID              string    `json:"id"`
	Amount          float64   `json:"amount"`
	Currency        string    `json:"currency"`
	RecipientName   string    `json:"recipientName,omitempty"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
}

type CreateBankPaymentRequest struct {
	Amount        float64 `json:"amount"`
	Currency      string  `json:"currency"`
	RecipientName string  `json:"recipientName,omitempty"`
}

// Chart of Accounts
type Account struct {
	AccountNumber int    `json:"account"`
	Name          string `json:"name"`
	AccountType   string `json:"accountType,omitempty"`
}

// Fiscal Year
type FiscalYear struct {
	ID               string `json:"id"`
	AccountingMethod string `json:"accountingMethod,omitempty"`
	StartDate        string `json:"startDate"`
	EndDate          string `json:"endDate"`
	Status           string `json:"status"`
}

// Connection (General API)
type Connection struct {
	ID        string    `json:"id"`
	TenantID  string    `json:"tenantId"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

// Token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TenantID     string `json:"tenant_id,omitempty"`
	TenantType   string `json:"tenant_type,omitempty"`
	ConnectionID string `json:"connection_id,omitempty"`
}
