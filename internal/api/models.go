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
	ID               string            `json:"id"`
	Name             string            `json:"name"`
	Type             string            `json:"type"`
	VatNumber        string            `json:"vatNumber,omitempty"`
	OrgNumber        string            `json:"orgNumber,omitempty"`
	PaymentTerms     string            `json:"paymentTerms,omitempty"`
	ContactsDetails  []CustomerContact `json:"contactsDetails,omitempty"`
	Address          *CustomerAddress  `json:"address,omitempty"`
	Language         string            `json:"language,omitempty"`
	ModifiedDateTime *time.Time        `json:"modifiedDateTime,omitempty"`
}

// DefaultContact returns the contact flagged as default, falling back to the
// first one, or nil when the customer has no contacts.
func (c *Customer) DefaultContact() *CustomerContact {
	for i := range c.ContactsDetails {
		if c.ContactsDetails[i].IsDefault {
			return &c.ContactsDetails[i]
		}
	}
	if len(c.ContactsDetails) > 0 {
		return &c.ContactsDetails[0]
	}
	return nil
}

type CustomerContact struct {
	ID        *string `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Email     string  `json:"email,omitempty"`
	Phone     string  `json:"phone,omitempty"`
	IsDefault bool    `json:"isDefault"`
}

// CustomerAddress is the spec's addressWithCountrySubdivision. When sent,
// line1, city, postalCode and country are required.
type CustomerAddress struct {
	Line1              string  `json:"line1"`
	Line2              *string `json:"line2,omitempty"`
	City               string  `json:"city"`
	PostalCode         string  `json:"postalCode"`
	Country            string  `json:"country"`
	CountrySubdivision *string `json:"countrySubdivision,omitempty"`
}

type CreateCustomerRequest struct {
	Name            string            `json:"name"`
	Type            string            `json:"type"`
	VatNumber       string            `json:"vatNumber,omitempty"`
	OrgNumber       string            `json:"orgNumber,omitempty"`
	PaymentTerms    string            `json:"paymentTerms,omitempty"`
	ContactsDetails []CustomerContact `json:"contactsDetails,omitempty"`
	Address         *CustomerAddress  `json:"address,omitempty"`
	Language        string            `json:"language,omitempty"`
}

type UpdateCustomerRequest = CreateCustomerRequest

// Supplier
type Supplier struct {
	ID             string                  `json:"id"`
	Name           string                  `json:"name"`
	OrgNumber      string                  `json:"orgNumber,omitempty"`
	VatNumber      string                  `json:"vatNumber,omitempty"`
	Currency       string                  `json:"currency,omitempty"`
	Address        *SupplierAddress        `json:"address,omitempty"`
	PaymentDetails *SupplierPaymentDetails `json:"paymentDetails,omitempty"`
}

// UpdateSupplierRequest is the supplier body without the read-only id.
type UpdateSupplierRequest struct {
	Name           string                  `json:"name"`
	OrgNumber      string                  `json:"orgNumber,omitempty"`
	VatNumber      string                  `json:"vatNumber,omitempty"`
	Currency       string                  `json:"currency,omitempty"`
	Address        *SupplierAddress        `json:"address,omitempty"`
	PaymentDetails *SupplierPaymentDetails `json:"paymentDetails,omitempty"`
}

type CreateSupplierRequest = UpdateSupplierRequest

// UpdateRequest returns the supplier as an update body, so an update can
// start from the current record rather than blanking omitted fields.
func (s *Supplier) UpdateRequest() UpdateSupplierRequest {
	return UpdateSupplierRequest{
		Name:           s.Name,
		OrgNumber:      s.OrgNumber,
		VatNumber:      s.VatNumber,
		Currency:       s.Currency,
		Address:        s.Address,
		PaymentDetails: s.PaymentDetails,
	}
}

type SupplierAddress struct {
	Line1      string  `json:"line1,omitempty"`
	Line2      *string `json:"line2,omitempty"`
	City       string  `json:"city,omitempty"`
	PostalCode string  `json:"postalCode,omitempty"`
	Country    string  `json:"country,omitempty"`
}

// SupplierPaymentDetails flattens the API's type-discriminated union; only the
// fields matching Type are set.
type SupplierPaymentDetails struct {
	Type           string `json:"type"`
	BankgiroNumber string `json:"bankgiroNumber,omitempty"`
	PlusgiroNumber string `json:"plusgiroNumber,omitempty"`
	ClearingNumber string `json:"clearingNumber,omitempty"`
	AccountNumber  string `json:"accountNumber,omitempty"`
	IBAN           string `json:"iban,omitempty"`
	BIC            string `json:"bic,omitempty"`
}

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

// JournalEntry is a verifikat: the posted bookkeeping record. Debits and credits
// across Items must balance. Entries are immutable once posted — corrections are
// made by reversing, not editing.
type JournalEntry struct {
	ID                 string             `json:"id,omitempty"`
	Title              string             `json:"title,omitempty"`
	JournalEntryNumber string             `json:"journalEntryNumber,omitempty"`
	Date               string             `json:"date"`
	Items              []JournalEntryItem `json:"items"`
	// Tags applies to the whole entry (verification level). Mutually exclusive
	// with per-item tags: use one level or the other, never both.
	Tags                     []JournalEntryTag `json:"tags,omitempty"`
	ReversingJournalEntryID  *string           `json:"reversingJournalEntryId,omitempty"`
	ReversedByJournalEntryID *string           `json:"reversedByJournalEntryId,omitempty"`
}

// JournalEntryItem is a single debit/credit line (kontering).
//
// ID is assigned by the server and is NOT guaranteed to follow the order the
// items were submitted in, nor the order they are returned in. To identify a
// line — for example to tag it via PUT /journal-entries/{id}/tags — match on
// Account plus amount rather than on slice index.
type JournalEntryItem struct {
	ID      int64             `json:"id,omitempty"`
	Debit   float64           `json:"debit"`
	Credit  float64           `json:"credit"`
	Account int               `json:"account"`
	Tags    []JournalEntryTag `json:"tags,omitempty"`
}

// JournalEntryTag references a tag value within a tag group (a dimension such as
// cost center or project). Weight allocates the amount across tags in a group and
// must be greater than 0 and at most 1; weights within a group must sum to 1 or less.
type JournalEntryTag struct {
	TagID        string  `json:"tagId"`
	TagGroupID   string  `json:"tagGroupId,omitempty"`
	TagName      string  `json:"tagName,omitempty"`
	TagGroupName string  `json:"tagGroupName,omitempty"`
	Weight       float64 `json:"weight"`
}

// CreateJournalEntryRequest is the body for POST /journal-entries. The API accepts
// the full journalEntry schema; read-only fields are omitted here.
type CreateJournalEntryRequest struct {
	Title string             `json:"title,omitempty"`
	Date  string             `json:"date"`
	Items []JournalEntryItem `json:"items"`
	Tags  []JournalEntryTag  `json:"tags,omitempty"`
}

// IsReversed reports whether this entry has been reversed by another entry.
func (e JournalEntry) IsReversed() bool {
	return e.ReversedByJournalEntryID != nil && *e.ReversedByJournalEntryID != ""
}

// Balance returns the summed debits and credits. They must be equal for the entry
// to be valid.
func (e JournalEntry) Balance() (debit, credit float64) {
	for _, item := range e.Items {
		debit += item.Debit
		credit += item.Credit
	}
	return debit, credit
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

// Upload is a file stored in Bokio, optionally linked to a journal entry.
// Description defaults to the file name when not supplied at upload time.
// JournalEntryID is nil for uploads that are not yet bookkept.
type Upload struct {
	ID             string  `json:"id"`
	Description    string  `json:"description,omitempty"`
	ContentType    string  `json:"contentType,omitempty"`
	JournalEntryID *string `json:"journalEntryId,omitempty"`
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
