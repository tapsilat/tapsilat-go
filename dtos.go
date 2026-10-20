package tapsilat

import (
	"strings"
	"time"
)

const (
	OrderStatusReceived = iota + 1
	OrderStatusUnpaid
	OrderStatusPaid
	OrderStatusProcessing
	OrderStatusShipped
	OrderStatusOnHold
	OrderStatusPayment
	OrderStatusCancelled
	OrderStatusCompleted
	OrderStatusRefunded
	OrderStatusFraud
	OrderStatusRejected
	OrderStatusFailure
	OrderStatusRetrying
	OrderStatusPartiallyRefunded
	OrderStatusSubMerchantPaymentApproved
	OrderStatusSubMerchantPaymentDisapproved
	OrderStatusSubMerchantPaymentErrored
	OrderStatusStillHasUnpaidInstallments
	OrderStatusStillHasUnpaidTerms
	OrderStatusExpired
	OrderStatusStillHasUnpaidSubMerchantPayments
	OrderStatusPartiallyPaid
	OrderStatusTerminated
	OrderStatusCardTokenization
	OrderStatusPreAuthorized
	OrderStatusDisputed
	OrderStatusPartiallyDisputed
	OrderStatusSuspect

	OrderTypePhysical       uint64 = 1
	OrderTypeVirtual        uint64 = 2
	OrderTypeMarketplace    uint64 = 3
	OrderTypeSubscription   uint64 = 4
	OrderTypeDeposit        uint64 = 5
	OrderTypeMailOrder      uint64 = 6
	OrderTypeTelephoneOrder uint64 = 7
	OrderTypeInvoice        uint64 = 8
	OrderTypeLoan           uint64 = 9
	OrderTypeRemittance     uint64 = 10
	OrderTypeCarLoan        uint64 = 11

	SubscriptionStatusSuccess = "success"
	SubscriptionStatusFailure = "failure"
	SubscriptionStatusPending = "pending"
)

var OrderStatuesMap = []struct {
	Id     int
	Status string
}{
	{1, "Received"},
	{2, "Unpaid"},
	{3, "Paid"},
	{4, "Processing"},
	{5, "Shipped"},
	{6, "On hold"},
	{7, "Waiting for payment"},
	{8, "Cancelled"},
	{9, "Completed"},
	{10, "Refunded"},
	{11, "Fraud"},
	{12, "Rejected"},
	{13, "Failure"},
	{14, "Retrying"},
	{15, "Partially refunded"},
	{16, "Sub merchant payment approved"},
	{17, "Sub merchant payment disapproved"},
	{18, "Sub merchant payment errored"},
	{19, "Still has unpaid installments"},
	{20, "Still has unpaid terms"},
	{21, "Expired"},
	{22, "Still has unpaid sub merchant payments"},
	{23, "Partially Paid"},
	{24, "Terminated"},
	{25, "Card Tokenization"},
	{26, "Pre Authorized"},
	{27, "Disputed"},
	{28, "Partially disputed"},
	{29, "Suspect"},
}

// GetOrderStatusByStr returns the order status id for a status name, ignoring case, or 0
func GetOrderStatusByStr(status string) int {
	for _, v := range OrderStatuesMap {
		if strings.EqualFold(v.Status, status) {
			return v.Id
		}
	}
	return 0
}

// OrderMetadata represents metadata key-value pairs
type OrderMetadata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// OrderPaymentTerm represents payment term for installments.
// DueDate and PaidDate must be RFC3339 timestamps (for example
// "2026-01-15T00:00:00Z"); the order-create endpoint rejects other formats and
// requires DueDate on every term.
type OrderPaymentTerm struct {
	Amount          *float64 `json:"amount,omitempty"`
	Data            string   `json:"data,omitempty"`
	DueDate         string   `json:"due_date,omitempty"`
	PaidDate        string   `json:"paid_date,omitempty"`
	Required        *bool    `json:"required,omitempty"`
	Status          string   `json:"status,omitempty"`
	TermReferenceID string   `json:"term_reference_id,omitempty"`
	TermSequence    *int     `json:"term_sequence,omitempty"`
}

type Order struct {
	Locale              string               `json:"locale"`
	Amount              float64              `json:"amount"`
	TaxAmount           float64              `json:"tax_amount"`
	Currency            string               `json:"currency"`
	ConversationID      string               `json:"conversation_id"`
	Buyer               OrderBuyer           `json:"buyer"`
	ShippingAddress     OrderShippingAddress `json:"shipping_address"`
	BillingAddress      OrderBillingAddress  `json:"billing_address"`
	BasketItems         []OrderBasketItem    `json:"basket_items"`
	Submerchants        []OrderSubmerchant   `json:"submerchants"`
	CheckoutDesign      OrderCheckoutDesign  `json:"checkout_design"`
	PaymentMethods      bool                 `json:"payment_methods"`
	PaymentFailureUrl   string               `json:"payment_failure_url"`
	PaymentSuccessUrl   string               `json:"payment_success_url"`
	RedirectSuccessUrl  string               `json:"redirect_success_url"`
	RedirectFailureUrl  string               `json:"redirect_failure_url"`
	PfSubMerchant       OrderPfSubMerchant   `json:"pf_sub_merchant"`
	ThreeDForce         bool                 `json:"three_d_force"`
	EnabledInstallments []int                `json:"enabled_installments"`
	PaymentOptions      []string             `json:"payment_options"`
	Metadata            []OrderMetadata      `json:"metadata"`
	PaymentTerms        []OrderPaymentTerm   `json:"payment_terms"`
	PaidAmount          float64              `json:"paid_amount,omitempty"`
	ExternalReferenceID string               `json:"external_reference_id,omitempty"`
	PartialPayment      bool                 `json:"partial_payment,omitempty"`
	// PaymentMode is "auth" (default) or "preauth".
	PaymentMode     string                `json:"payment_mode,omitempty"`
	Consents        []OrderConsent        `json:"consents,omitempty"`
	OrderCards      []OrderCard           `json:"order_cards,omitempty"`
	SubOrganization *OrderSubOrganization `json:"sub_organization,omitempty"`
	OrderVposID     string                `json:"order_vpos_id,omitempty"`
}

// OrderConsent is a consent text the buyer must accept on the checkout page.
type OrderConsent struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// OrderCard selects a saved card for the order.
type OrderCard struct {
	CardID       string `json:"card_id"`
	CardSequence uint64 `json:"card_sequence,omitempty"`
}

// OrderSubOrganization describes the sub-organization an order is created for.
type OrderSubOrganization struct {
	OrganizationName      string `json:"organization_name,omitempty"`
	Email                 string `json:"email,omitempty"`
	Iban                  string `json:"iban,omitempty"`
	Address               string `json:"address,omitempty"`
	Currency              string `json:"currency,omitempty"`
	Acquirer              string `json:"acquirer,omitempty"`
	TaxOffice             string `json:"tax_office,omitempty"`
	TaxNumber             string `json:"tax_number,omitempty"`
	GsmNumber             string `json:"gsm_number,omitempty"`
	ContactFirstName      string `json:"contact_first_name,omitempty"`
	ContactLastName       string `json:"contact_last_name,omitempty"`
	IdentityNumber        string `json:"identity_number,omitempty"`
	LegalCompanyTitle     string `json:"legal_company_title,omitempty"`
	SubMerchantKey        string `json:"sub_merchant_key,omitempty"`
	SubMerchantType       string `json:"sub_merchant_type,omitempty"`
	SubMerchantExternalID string `json:"sub_merchant_external_id,omitempty"`
}

type OrderPfSubMerchant struct {
	Address        string `json:"address"`
	City           string `json:"city"`
	Country        string `json:"country"`
	CountryISOCode string `json:"country_iso_code"`
	ID             string `json:"id"`
	MCC            string `json:"mcc"`
	Name           string `json:"name"`
	OrgID          string `json:"org_id"`
	PostalCode     string `json:"postal_code"`
	TerminalNo     string `json:"terminal_no"`
	SubmerchantNIN string `json:"submerchant_nin,omitempty"`
	SubmerchantURL string `json:"submerchant_url,omitempty"`
	SwitchID       string `json:"switch_id,omitempty"`
	NationalID     string `json:"national_id,omitempty"`
}

type OrderDetail struct {
	Locale          string                 `json:"locale"`
	Error           string                 `json:"error"`
	Code            int                    `json:"code"`
	ReferenceID     string                 `json:"reference_id"`
	Amount          string                 `json:"amount"`
	Total           string                 `json:"total"`
	PaidAmount      string                 `json:"paid_amount"`
	RefundedAmount  string                 `json:"refunded_amount"`
	CreatedAt       string                 `json:"created_at"`
	Currency        string                 `json:"currency"`
	Status          int32                  `json:"status"`
	StatusEnum      string                 `json:"status_enum"`
	Buyer           OrderBuyer             `json:"buyer"`
	ShippingAddress OrderShippingAddress   `json:"shipping_address"`
	CheckoutDesign  OrderCheckoutDesignDTO `json:"checkout_design"`
	BillingAddress  OrderBillingAddress    `json:"billing_address"`
	BasketItems     []OrderBasketItem      `json:"basket_items"`
	// Deprecated: the order detail endpoint does not return submerchants.
	Submerchants []OrderSubmerchant    `json:"submerchants"`
	PaymentTerms []OrderPaymentTermDTO `json:"payment_terms"`
	// Deprecated: item payments are returned per basket item in BasketItems.
	ItemPayments         []OrderItemPayment         `json:"item_payments"`
	PaymentFailureUrl    string                     `json:"payment_failure_url" example:"https://www.example.com/payment/failure"`
	PaymentSuccessUrl    string                     `json:"payment_success_url" example:"https://www.example.com/payment/success"`
	RedirectSuccessUrl   string                     `json:"redirect_success_url"`
	RedirectFailureUrl   string                     `json:"redirect_failure_url"`
	CheckoutURL          string                     `json:"checkout_url" example:"https://www.example.com/payment/checkout"`
	ConversationID       string                     `json:"conversation_id" example:"123456789"`
	PaymentOptions       []string                   `json:"payment_options" example:"credit_card,bank_transfer,cash"`
	ExternalReferenceID  string                     `json:"external_reference_id" example:"ext-123456"`
	Metadata             []OrderMetadata            `json:"metadata"`
	PaymentRequests      []OrderPaymentRequest      `json:"payment_requests"`
	RefundCancelRequests []OrderRefundCancelRequest `json:"refund_cancel_requests"`
	CustomerPaidAmount   float64                    `json:"customer_paid_amount"`
	// Deprecated: the merchant category code is returned per basket item.
	MCC string `json:"mcc" example:"5411"`
}

// OrderPaymentRequest is a payment attempt recorded on an order.
type OrderPaymentRequest struct {
	ID          string `json:"id"`
	Date        string `json:"date"`
	PaymentMode string `json:"payment_mode"`
	Amount      string `json:"amount"`
	CardHash    string `json:"card_hash"`
	MaskedCard  string `json:"masked_card"`
	Paid        bool   `json:"paid"`
	Type        string `json:"type"`
}

// OrderRefundCancelRequest is a refund or cancellation request recorded on an order.
type OrderRefundCancelRequest struct {
	ID        string `json:"id"`
	Amount    string `json:"amount"`
	Status    uint64 `json:"status"`
	Reason    string `json:"reason"`
	Type      uint64 `json:"type"`
	TermID    string `json:"term_id"`
	ItemID    string `json:"item_id"`
	CreatedAt string `json:"created_at"`
}
type OrderPaymentTermDTO struct {
	ID              string             `json:"id" example:"123456789"`
	HashID          string             `json:"hash_id" example:"123456789"`
	TermSequence    uint64             `json:"term_sequence" example:"1"`
	Required        bool               `json:"required" example:"true"`
	DueDate         time.Time          `json:"due_date" example:"2019-01-01 00:00:00"`
	PaidDate        time.Time          `json:"paid_date" example:"2019-01-01 00:00:00"`
	Amount          float64            `json:"amount" example:"100.00"`
	TermReferenceID string             `json:"term_reference_id" example:"41f8fce7-71a7-4d55-a603-6a4bd2f30d07"`
	Status          string             `json:"status" example:"pending"`
	Payments        []OrderTermPayment `json:"payments"`
	Data            string             `json:"data" example:"data"`
} // @name OrderPaymentTermDTO
type OrderTermPayment struct {
	Id               string  `json:"id,omitempty"`
	TermID           string  `json:"term_id,omitempty"`
	Amount           float64 `json:"amount,omitempty"`
	PaidDate         string  `json:"paid_date,omitempty"`
	MaskedBin        string  `json:"masked_bin,omitempty"`
	CardBrand        string  `json:"card_brand,omitempty"`
	RefundedAmount   float64 `json:"refunded_amount,omitempty"`
	RefundableAmount float64 `json:"refundable_amount,omitempty"`
	Refunded         bool    `json:"refunded,omitempty"`
	Status           uint64  `json:"status,omitempty" example:"1"`
	Type             uint64  `json:"type,omitempty" example:"1"` //  Credit Card,  Bank Transfer etc.
} //
type OrderBuyer struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	Surname             string `json:"surname"`
	Email               string `json:"email"`
	GsmNumber           string `json:"gsm_number"`
	IdentityNumber      string `json:"identity_number"`
	RegistrationDate    string `json:"registration_date"`
	RegistrationAddress string `json:"registration_address"`
	LastLoginDate       string `json:"last_login_date"`
	City                string `json:"city"`
	Country             string `json:"country"`
	ZipCode             string `json:"zip_code"`
	Ip                  string `json:"ip"`
	// BirdthDate is the buyer's birth date (YYYY-MM-DD). The field name keeps
	// its historical spelling; it is sent as birth_date.
	BirdthDate string `json:"birth_date"`
	Title      string `json:"title,omitempty"`
	IncomeType string `json:"income_type,omitempty"`
	Education  string `json:"education,omitempty"`
	Occupation string `json:"occupation,omitempty"`
}

type OrderShippingAddress struct {
	Address      string `json:"address"`
	ZipCode      string `json:"zip_code"`
	City         string `json:"city"`
	Country      string `json:"country"`
	ContactName  string `json:"contact_name"`
	TrackingCode string `json:"tracking_code"`
	ShippingDate string `json:"shipping_date" example:"2019-01-01 00:00:00"`
}

type OrderBasketItem struct {
	Id               string                `json:"id" example:"123456789"`
	Price            float64               `json:"price" example:"100.00"`
	Name             string                `json:"name" example:"Product Name"`
	Category1        string                `json:"category1" example:"Category 1"`
	Category2        string                `json:"category2" example:"Category 2"`
	ItemType         string                `json:"item_type" example:"PHYSICAL"`
	Status           uint64                `json:"status" example:"1"`
	RefundedAmount   float64               `json:"refunded_amount" example:"0.00"`
	RefundableAmount float64               `json:"refundable_amount" example:"0.00"`
	PaidAmount       float64               `json:"paid_amount" example:"0.00"`
	PaidableAmount   float64               `json:"paidable_amount" example:"0.00"`
	Coupon           string                `json:"coupon" example:"coupon"`
	CouponDiscount   float64               `json:"coupon_discount" example:"0.00"`
	ItemPayments     []OrderItemPayment    `json:"item_payments"`
	Quantity         *int                  `json:"quantity,omitempty"`
	QuantityFloat    *float64              `json:"quantity_float,omitempty"`
	QuantityUnit     string                `json:"quantity_unit,omitempty"`
	Data             string                `json:"data,omitempty"`
	CommissionAmount *float64              `json:"commission_amount,omitempty"`
	SubMerchantKey   string                `json:"sub_merchant_key,omitempty"`
	SubMerchantPrice string                `json:"sub_merchant_price,omitempty"`
	PayoutReleaseAt  *time.Time            `json:"payout_release_at,omitempty"`
	Payer            *OrderBasketItemPayer `json:"payer,omitempty"`
	// MCC is the item's merchant category code.
	MCC string `json:"mcc,omitempty"`
}

// OrderBasketItemPayer represents the payer information for basket items
type OrderBasketItemPayer struct {
	Address        string `json:"address,omitempty"`
	ReferenceID    string `json:"reference_id,omitempty"`
	TaxOffice      string `json:"tax_office,omitempty"`
	Title          string `json:"title,omitempty"`
	Type           string `json:"type,omitempty"`
	VAT            string `json:"vat,omitempty"`
	Name           string `json:"name,omitempty"`
	Surname        string `json:"surname,omitempty"`
	IdentityNumber string `json:"identity_number,omitempty"`
	Email          string `json:"email,omitempty"`
	Phone          string `json:"phone,omitempty"`
}

type OrderSubmerchant struct {
	Amount              float64 `json:"amount"`
	OrderBasketItemID   string  `json:"order_basket_item_id"`
	MerchantReferenceID string  `json:"merchant_reference_id"`
}

// OrderCheckoutDesign customizes the checkout page branding for one order.
type OrderCheckoutDesign struct {
	Logo                 string            `json:"logo"`
	InputBackgroundColor string            `json:"input_background_color"`
	PayButtonColor       string            `json:"pay_button_color"`
	InputTextColor       string            `json:"input_text_color"`
	LabelTextColor       string            `json:"label_text_color"`
	LeftBackgroundColor  string            `json:"left_background_color"`
	RightBackgroundColor string            `json:"right_background_color"`
	TextColor            string            `json:"text_color"`
	PayButtonTextColor   string            `json:"pay_button_text_color,omitempty"`
	PageBackgroundColor  string            `json:"page_background_color,omitempty"`
	BorderColor          string            `json:"border_color,omitempty"`
	InputBorderColor     string            `json:"input_border_color,omitempty"`
	AccentColor          string            `json:"accent_color,omitempty"`
	TextMutedColor       string            `json:"text_muted_color,omitempty"`
	SurfaceMutedColor    string            `json:"surface_muted_color,omitempty"`
	DarkModeEnabled      bool              `json:"dark_mode_enabled,omitempty"`
	DarkColors           map[string]string `json:"dark_colors,omitempty"`
	// Deprecated: checkout no longer renders custom order detail HTML; the
	// value is ignored.
	OrderDetailHtml string `json:"order_detail_html,omitempty"`
	// Deprecated: checkout design no longer carries a redirect URL; the value
	// is ignored. Use Order.RedirectSuccessUrl and Order.RedirectFailureUrl.
	RedirectUrl string `json:"redirect_url,omitempty"`
}
type OrderItemPayment struct {
	Id               string  `json:"id,omitempty"`
	Amount           float64 `json:"amount,omitempty"`
	PaidDate         string  `json:"paid_date,omitempty"`
	MaskedBin        string  `json:"masked_bin,omitempty"`
	CardBrand        string  `json:"card_brand,omitempty"`
	RefundedAmount   float64 `json:"refunded_amount,omitempty"`
	RefundableAmount float64 `json:"refundable_amount,omitempty"`
	RefundedDate     string  `json:"refunded_date,omitempty"`
	Refunded         bool    `json:"refunded,omitempty"`
	Status           uint64  `json:"status,omitempty" example:"1"`
	Type             string  `json:"type,omitempty"`
} // @name OrderItemPayment

type OrderBillingAddress struct {
	BillingType   string `json:"billing_type" example:"PERSONAL"`   // PERSONAL, BUSINESS
	Citizenship   string `json:"citizenship" example:"TR"`          // ISO 3166-1 alpha-2 country code
	Title         string `json:"title" example:"MonoPayments Inc."` // Legal title
	TaxOffice     string `json:"tax_office" example:"Merter"`
	Address       string `json:"address" example:"Istanbul"`
	ZipCode       string `json:"zip_code" example:"34000"`
	City          string `json:"city" example:"Istanbul"`
	District      string `json:"district" example:"Uskudar"`
	Country       string `json:"country" example:"Turkey"`
	ContactName   string `json:"contact_name" example:"John Doe"`
	ContactPhone  string `json:"contact_phone" example:"+905555555555"`
	VatNumber     string `json:"vat_number" example:"1234567890"` // Identification number if billing type is PERSONAL
	Neighbourhood string `json:"neighbourhood,omitempty"`
	Street1       string `json:"street1,omitempty"`
	Street2       string `json:"street2,omitempty"`
	Street3       string `json:"street3,omitempty"`
}

type OrderResponse struct {
	OrderID     string `json:"order_id"`
	ReferenceID string `json:"reference_id"`
	CheckoutURL string `json:"checkout_url"`
	Error       string `json:"error"`
}

type OrderStatus struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

type RefundOrder struct {
	ReferenceID string  `json:"reference_id"`
	Amount      float64 `json:"amount"`
	Error       string  `json:"error"`
	// OrderItemID and OrderItemPaymentID optionally limit the refund to one
	// basket item or item payment.
	OrderItemID        string `json:"order_item_id,omitempty"`
	OrderItemPaymentID string `json:"order_item_payment_id,omitempty"`
}

type CancelOrder struct {
	ReferenceID string `json:"reference_id"`
	Error       string `json:"error"`
}

type RefundCancelOrderResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	IsSuccess bool   `json:"is_success"`
	Error     string `json:"error"`
}

type PaginatedData struct {
	Page       int64  `json:"page,omitempty" example:"1"`
	PerPage    int64  `json:"per_page,omitempty" example:"10"`
	Total      int64  `json:"total,omitempty" example:"100"`
	TotalPages int    `json:"total_pages,omitempty" example:"10"`
	Rows       any    `json:"rows,omitempty" swaggertype:"array,string" example:"object,object2"`
	Error      string `json:"error"`
}

type OrderCheckoutDesignDTO struct {
	Logo                 string `json:"logo" example:"https://www.example.com/logo.png"`
	InputBackgroundColor string `json:"input_background_color" example:"#ffffff"`
	InputTextColor       string `json:"input_text_color" example:"#000000"`
	LabelTextColor       string `json:"label_text_color" example:"#000000"`
	LeftBackgroundColor  string `json:"left_background_color" example:"#ffffff"`
	RightBackgroundColor string `json:"right_background_color" example:"#ffffff"`
	TextColor            string `json:"text_color" example:"#000000"`
	PlaceholderColor     string `json:"placeholder_color" example:"#000000"`
	PayButtonColor       string `json:"pay_button_color,omitempty"`
	PayButtonTextColor   string `json:"pay_button_text_color,omitempty"`
	PageBackgroundColor  string `json:"page_background_color,omitempty"`
	BorderColor          string `json:"border_color,omitempty"`
	InputBorderColor     string `json:"input_border_color,omitempty"`
	AccentColor          string `json:"accent_color,omitempty"`
	TextMutedColor       string `json:"text_muted_color,omitempty"`
	SurfaceMutedColor    string `json:"surface_muted_color,omitempty"`
	DarkModeEnabled      bool   `json:"dark_mode_enabled,omitempty"`
	// Deprecated: checkout no longer renders custom order detail HTML.
	OrderDetailHtml string            `json:"order_detail_html" example:"<html><body><h1>Order Detail</h1></body></html>"`
	DarkColors      map[string]string `json:"dark_colors,omitempty"`
} // @name OrderCheckoutDesignDTO

// OrderPaymentTermCreateDTO adds a payment term to an existing order.
// DueDate and PaidDate use the YYYY-MM-DD format.
type OrderPaymentTermCreateDTO struct {
	// OrderID is the order's ID (OrderResponse.OrderID), not its reference ID.
	OrderID         string  `json:"order_id"`
	TermSequence    *int    `json:"term_sequence,omitempty"`
	Required        bool    `json:"required"`
	Amount          float64 `json:"amount"`
	DueDate         string  `json:"due_date"`
	PaidDate        string  `json:"paid_date,omitempty"`
	TermReferenceID string  `json:"term_reference_id,omitempty"`
	Status          string  `json:"status,omitempty"`
	Data            string  `json:"data,omitempty"`
}

// OrderPaymentTermUpdateDTO replaces a payment term. The server stores every
// field as sent, so populate all of them, not only the ones that change.
// DueDate and PaidDate use the YYYY-MM-DD format.
type OrderPaymentTermUpdateDTO struct {
	TermReferenceID string   `json:"term_reference_id"`
	TermSequence    *int     `json:"term_sequence,omitempty"`
	Required        *bool    `json:"required,omitempty"`
	Amount          *float64 `json:"amount,omitempty"`
	DueDate         string   `json:"due_date,omitempty"`
	PaidDate        string   `json:"paid_date,omitempty"`
	Status          string   `json:"status,omitempty"`
	// Deprecated: the server ignores data on term updates.
	Data string `json:"data,omitempty"`
}

// OrderTermRefundRequest refunds a payment made against a payment term.
type OrderTermRefundRequest struct {
	// ReferenceID is the order's reference ID; it is required.
	ReferenceID string `json:"reference_id"`
	// TermID is the payment term's ID (OrderPaymentTermDTO.ID).
	TermID        string  `json:"term_id"`
	TermPaymentID string  `json:"term_payment_id,omitempty"`
	Amount        float64 `json:"amount"`
}

type SubmerchantCreateRequest struct {
	Locale                string `json:"locale"`
	ConversationID        string `json:"conversation_id"`
	Name                  string `json:"name"`
	Email                 string `json:"email"`
	GsmNumber             string `json:"gsm_number"`
	Address               string `json:"address"`
	Iban                  string `json:"iban"`
	TaxOffice             string `json:"tax_office"`
	LegalCompanyTitle     string `json:"legal_company_title"`
	CurrencyID            string `json:"currency_id"`
	SubmerchantExternalID string `json:"sub_merchant_external_id"`
	IdentityNumber        string `json:"identity_number"`
	SubmerchantType       string `json:"sub_merchant_type"`
	TaxNumber             string `json:"tax_number"`
	SubmerchantKey        string `json:"sub_merchant_key"`
	OrganizationID        string `json:"organization_id"`
	Status                string `json:"status"`
	SystemTime            int64  `json:"system_time"`
	ContactName           string `json:"contact_name"`
	ContactSurname        string `json:"contact_surname"`
}

type MarketplaceSubmerchantCreateRequest struct {
	// ProviderAccounts explicitly selects accounts in one environment; an empty list creates only the seller.
	// It cannot be combined with VposID.
	ProviderAccounts      *[]MarketplaceAccountTarget `json:"provider_accounts,omitempty"`
	Locale                string                      `json:"locale"`
	ConversationID        string                      `json:"conversation_id"`
	Name                  string                      `json:"name"`
	Email                 string                      `json:"email"`
	GsmNumber             string                      `json:"gsm_number"`
	Address               string                      `json:"address"`
	Iban                  string                      `json:"iban"`
	TaxOffice             string                      `json:"tax_office"`
	LegalCompanyTitle     string                      `json:"legal_company_title"`
	CurrencyID            string                      `json:"currency_id"`
	SubmerchantExternalID string                      `json:"sub_merchant_external_id"`
	IdentityNumber        string                      `json:"identity_number"`
	SubmerchantType       string                      `json:"sub_merchant_type"`
	TaxNumber             string                      `json:"tax_number"`
	Labels                []string                    `json:"labels,omitempty"`
	Status                string                      `json:"status"`
	SystemTime            int64                       `json:"system_time"`
	ContactName           string                      `json:"contact_name"`
	ContactSurname        string                      `json:"contact_surname"`
	// VposID targets only this account, using its configured environment; it never selects other accounts.
	// When VposID and ProviderAccounts are omitted, only the seller is created.
	VposID          string  `json:"vpos_id,omitempty"`
	ApprovalMode    *string `json:"approval_mode,omitempty"`
	ReleasePolicyID *string `json:"release_policy_id,omitempty"`
	IdempotencyKey  string  `json:"idempotency_key,omitempty"`
}

type MarketplaceAccountTarget struct {
	VposID      string `json:"vpos_id"`
	Environment string `json:"environment"`
}

type MarketplaceSubmerchantCreateResponse struct {
	ID                string                               `json:"id"`
	SuborganizationID string                               `json:"suborganization_id"`
	VposSubmerchantID string                               `json:"vpos_submerchant_id"`
	RoutingReference  string                               `json:"routing_reference"`
	Provisionings     []MarketplaceSubmerchantProvisioning `json:"provisionings"`
}

type MarketplaceSubmerchantProvisioning struct {
	VposID                 string `json:"vpos_id"`
	VposSubmerchantID      string `json:"vpos_submerchant_id,omitempty"`
	Provider               string `json:"provider"`
	Status                 string `json:"status"`
	ProviderSubmerchantKey string `json:"provider_submerchant_key,omitempty"`
	Retryable              bool   `json:"retryable"`
	ErrorMessage           string `json:"error_message,omitempty"`
}

type SubmerchantPayoutEventRequest struct {
	IdempotencyKey string         `json:"idempotency_key"`
	OrderReference string         `json:"order_reference"`
	ItemID         string         `json:"item_id"`
	EventType      string         `json:"event_type"`
	OccurredAt     *time.Time     `json:"occurred_at,omitempty"`
	Payload        map[string]any `json:"payload,omitempty"`
}

type SubmerchantPayoutEventResponse struct {
	ID             string         `json:"id"`
	IdempotencyKey string         `json:"idempotency_key"`
	OrderReference string         `json:"order_reference"`
	ItemID         string         `json:"item_id"`
	EventType      string         `json:"event_type"`
	OccurredAt     time.Time      `json:"occurred_at"`
	Payload        map[string]any `json:"payload,omitempty"`
}

type SubmerchantUpdateRequest struct {
	Locale                string `json:"locale"`
	ConversationID        string `json:"conversation_id"`
	Name                  string `json:"name"`
	Email                 string `json:"email"`
	GsmNumber             string `json:"gsm_number"`
	Address               string `json:"address"`
	Iban                  string `json:"iban"`
	TaxOffice             string `json:"tax_office"`
	LegalCompanyTitle     string `json:"legal_company_title"`
	CurrencyID            string `json:"currency_id"`
	SubmerchantExternalID string `json:"sub_merchant_external_id"`
	IdentityNumber        string `json:"identity_number"`
	SubmerchantType       string `json:"sub_merchant_type"`
	TaxNumber             string `json:"tax_number"`
	SubmerchantKey        string `json:"sub_merchant_key"`
	OrganizationID        string `json:"organization_id"`
	Status                string `json:"status"`
	SystemTime            int64  `json:"system_time"`
	ContactName           string `json:"contact_name"`
	ContactSurname        string `json:"contact_surname"`
}

type Submerchant struct {
	ID                    string `json:"id,omitempty"`
	Locale                string `json:"locale,omitempty"`
	ConversationID        string `json:"conversation_id,omitempty"`
	Name                  string `json:"name,omitempty"`
	Email                 string `json:"email,omitempty"`
	GsmNumber             string `json:"gsm_number,omitempty"`
	Address               string `json:"address,omitempty"`
	Iban                  string `json:"iban,omitempty"`
	TaxOffice             string `json:"tax_office,omitempty"`
	LegalCompanyTitle     string `json:"legal_company_title,omitempty"`
	CurrencyID            string `json:"currency_id,omitempty"`
	SubmerchantExternalID string `json:"sub_merchant_external_id,omitempty"`
	IdentityNumber        string `json:"identity_number,omitempty"`
	SubmerchantType       string `json:"sub_merchant_type,omitempty"`
	TaxNumber             string `json:"tax_number,omitempty"`
	SubmerchantKey        string `json:"sub_merchant_key,omitempty"`
	OrganizationID        string `json:"organization_id,omitempty"`
	Labels                string `json:"labels,omitempty"`
	Status                string `json:"status,omitempty"`
	SystemTime            int64  `json:"system_time,omitempty"`
	Acquirer              string `json:"acquirer,omitempty"`
	ContactName           string `json:"contact_name,omitempty"`
	ContactSurname        string `json:"contact_surname,omitempty"`
}

type SubmerchantListItem struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Email           string `json:"email,omitempty"`
	SubmerchantType string `json:"submerchant_type,omitempty"`
	SubmerchantKey  string `json:"submerchant_key,omitempty"`
	Labels          string `json:"labels,omitempty"`
	Status          string `json:"status,omitempty"`
	Acquirer        string `json:"acquirer,omitempty"`
}

type SubmerchantListResponse struct {
	Page       int64                 `json:"page,omitempty"`
	PerPage    int64                 `json:"per_page,omitempty"`
	Total      int64                 `json:"total,omitempty"`
	TotalPages int64                 `json:"total_pages,omitempty"`
	Rows       []SubmerchantListItem `json:"row,omitempty"`
}

// SubmerchantMutationResponse is returned by submerchant create, update and
// delete. The server only sends Code and Message; read the submerchant with
// GetSubmerchant for its key and status.
type SubmerchantMutationResponse struct {
	Code    uint64 `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	// Deprecated: never returned by the server.
	Status string `json:"status,omitempty"`
	// Deprecated: never returned by the server.
	Locale string `json:"locale,omitempty"`
	// Deprecated: never returned by the server.
	SystemTime int64 `json:"system_time,omitempty"`
	// Deprecated: never returned by the server.
	ConversationID string `json:"conversation_id,omitempty"`
	// Deprecated: never returned by the server; use GetSubmerchant.
	SubmerchantKey string `json:"sub_merchant_key,omitempty"`
}

type SuborganizationListItem struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Address  string `json:"address,omitempty"`
	City     string `json:"city,omitempty"`
	Country  string `json:"country,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Currency string `json:"currency,omitempty"`
	Logo     string `json:"logo,omitempty"`
}

type SuborganizationDetail struct {
	ID                 string `json:"id,omitempty"`
	Name               string `json:"name,omitempty"`
	ParentID           string `json:"parent_id,omitempty"`
	PublicStatus       int64  `json:"public_status,omitempty"`
	AvailabilityStatus int64  `json:"availability_status,omitempty"`
	CreatedAt          string `json:"created_at,omitempty"`
	UpdatedAt          string `json:"updated_at,omitempty"`
}

type SuborganizationListResponse struct {
	Page       int64                     `json:"page,omitempty"`
	PerPage    int64                     `json:"per_page,omitempty"`
	Total      int64                     `json:"total,omitempty"`
	TotalPages int64                     `json:"total_pages,omitempty"`
	Rows       []SuborganizationListItem `json:"rows,omitempty"`
}

type SubmerchantSuborganizationMapping struct {
	SubmerchantID     string `json:"submerchant_id,omitempty"`
	SuborganizationID string `json:"suborganization_id,omitempty"`
}

type SuborganizationSubmerchantMapping struct {
	SuborganizationID string `json:"suborganization_id,omitempty"`
	SubmerchantID     string `json:"submerchant_id,omitempty"`
}

type VposCreateRequest struct {
	Name                   string   `json:"name"`
	BankName               string   `json:"bank_name"`
	EnvMode                string   `json:"env_mode"`
	Merchant               string   `json:"merchant"`
	MerchantCode           string   `json:"merchant_code"`
	MerchantKey            string   `json:"merchant_key"`
	AuthKey                string   `json:"auth_key"`
	Terminal               string   `json:"terminal"`
	Company                string   `json:"company"`
	Username               string   `json:"username"`
	Password               string   `json:"password"`
	StoreKey               string   `json:"store_key"`
	ApiKey                 string   `json:"api_key"`
	ApiSecret              string   `json:"api_secret"`
	GUID                   string   `json:"guid"`
	PID                    string   `json:"pid"`
	ClientID               string   `json:"client_id"`
	ClientCode             string   `json:"client_code"`
	Type                   string   `json:"type"`
	BlockDate              uint64   `json:"block_date"`
	Priority               int64    `json:"priority"`
	Main                   bool     `json:"main"`
	ForceThreeD            bool     `json:"force_three_d"`
	Marketplace            bool     `json:"marketplace"`
	Pf                     bool     `json:"pf"`
	PaymentMode            string   `json:"payment_mode"`
	Prefix                 string   `json:"prefix"`
	AcquirerID             string   `json:"acquirer_id"`
	CardSchemes            []string `json:"card_schemes"`
	Currencies             []string `json:"currencies"`
	Failover               bool     `json:"failover"`
	PfID                   string   `json:"pf_id"`
	RedirectURI            string   `json:"redirect_uri"`
	TLS                    bool     `json:"tls"`
	TLSKey                 string   `json:"tls_key"`
	TLSCert                string   `json:"tls_cert"`
	RootCA                 string   `json:"root_ca"`
	MainInstallment        bool     `json:"main_installment"`
	MainWithoutInstallment bool     `json:"main_without_installment"`
}

// VposUpdateRequest replaces the VPOS configuration: fields left at their zero
// value are stored as zero (for example Failover and TLS become false). Read
// the VPOS with GetVpos and send every field back.
type VposUpdateRequest struct {
	Name                   string   `json:"name"`
	BankName               string   `json:"bank_name"`
	EnvMode                string   `json:"env_mode"`
	Merchant               string   `json:"merchant"`
	MerchantCode           string   `json:"merchant_code"`
	MerchantKey            string   `json:"merchant_key"`
	AuthKey                string   `json:"auth_key"`
	Terminal               string   `json:"terminal"`
	Company                string   `json:"company"`
	Username               string   `json:"username"`
	Password               string   `json:"password"`
	StoreKey               string   `json:"store_key"`
	ApiKey                 string   `json:"api_key"`
	ApiSecret              string   `json:"api_secret"`
	GUID                   string   `json:"guid"`
	PID                    string   `json:"pid"`
	ClientID               string   `json:"client_id"`
	ClientCode             string   `json:"client_code"`
	Type                   string   `json:"type"`
	BlockDate              uint64   `json:"block_date"`
	Priority               int64    `json:"priority"`
	Main                   bool     `json:"main"`
	ForceThreeD            bool     `json:"force_three_d"`
	Marketplace            bool     `json:"marketplace"`
	Pf                     bool     `json:"pf"`
	PaymentMode            string   `json:"payment_mode"`
	Prefix                 string   `json:"prefix"`
	AcquirerID             string   `json:"acquirer_id"`
	CardSchemes            []string `json:"card_schemes"`
	Currencies             []string `json:"currencies"`
	Failover               bool     `json:"failover"`
	PfID                   string   `json:"pf_id"`
	RedirectURI            string   `json:"redirect_uri"`
	TLS                    bool     `json:"tls"`
	TLSKey                 string   `json:"tls_key"`
	TLSCert                string   `json:"tls_cert"`
	RootCA                 string   `json:"root_ca"`
	MainInstallment        bool     `json:"main_installment"`
	MainWithoutInstallment bool     `json:"main_without_installment"`
}

type Vpos struct {
	ID                     string   `json:"id,omitempty"`
	Name                   string   `json:"name,omitempty"`
	BankName               string   `json:"bank_name,omitempty"`
	EnvMode                string   `json:"env_mode,omitempty"`
	Provider               string   `json:"provider,omitempty"`
	PaymentMode            string   `json:"payment_mode,omitempty"`
	Merchant               string   `json:"merchant,omitempty"`
	MerchantCode           string   `json:"merchant_code,omitempty"`
	MerchantKey            string   `json:"merchant_key,omitempty"`
	AuthKey                string   `json:"auth_key,omitempty"`
	Terminal               string   `json:"terminal,omitempty"`
	Company                string   `json:"company,omitempty"`
	Username               string   `json:"username,omitempty"`
	Password               string   `json:"password,omitempty"`
	StoreKey               string   `json:"store_key,omitempty"`
	ApiKey                 string   `json:"api_key,omitempty"`
	ApiSecret              string   `json:"api_secret,omitempty"`
	GUID                   string   `json:"guid,omitempty"`
	PID                    string   `json:"pid,omitempty"`
	ClientID               string   `json:"client_id,omitempty"`
	ClientCode             string   `json:"client_code,omitempty"`
	Type                   string   `json:"type,omitempty"`
	BlockDate              uint64   `json:"block_date,omitempty"`
	Priority               int64    `json:"priority,omitempty"`
	Main                   bool     `json:"main,omitempty"`
	ForceThreeD            bool     `json:"force_three_d,omitempty"`
	Marketplace            bool     `json:"marketplace,omitempty"`
	Pf                     bool     `json:"pf,omitempty"`
	Prefix                 string   `json:"prefix,omitempty"`
	AcquirerID             string   `json:"acquirer_id,omitempty"`
	CardSchemes            []string `json:"card_schemes,omitempty"`
	Currencies             []string `json:"currencies,omitempty"`
	Failover               bool     `json:"failover,omitempty"`
	PfID                   string   `json:"pf_id,omitempty"`
	RedirectURI            string   `json:"redirect_uri,omitempty"`
	TLS                    bool     `json:"tls,omitempty"`
	TLSKey                 string   `json:"tls_key,omitempty"`
	TLSCert                string   `json:"tls_cert,omitempty"`
	RootCA                 string   `json:"root_ca,omitempty"`
	MainInstallment        bool     `json:"main_installment,omitempty"`
	MainWithoutInstallment bool     `json:"main_without_installment,omitempty"`
}

type VposListItem struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	BankName    string `json:"bank_name,omitempty"`
	EnvMode     string `json:"env_mode,omitempty"`
	Provider    string `json:"provider,omitempty"`
	PaymentMode string `json:"payment_mode,omitempty"`
}

type VposListFilter struct {
	SuborganizationID string `json:"suborganization_id,omitempty"`
}

type VposListResponse struct {
	Page       int64          `json:"page,omitempty"`
	PerPage    int64          `json:"per_page,omitempty"`
	Total      int64          `json:"total,omitempty"`
	TotalPages int64          `json:"total_pages,omitempty"`
	Rows       []VposListItem `json:"rows,omitempty"`
}

type VposMutationResponse struct {
	Code    uint64 `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type VposAcquirer struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Prefix string `json:"prefix,omitempty"`
}

type VposAcquirerListResponse struct {
	Items []VposAcquirer `json:"items,omitempty"`
}

type CardScheme struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type CardSchemeListResponse struct {
	Items []CardScheme `json:"items,omitempty"`
}

type VposSubmerchantCreateRequest struct {
	ExternalReferenceID string `json:"external_reference_id"`
	SubmerchantID       string `json:"submerchant_id"`
	TerminalNo          string `json:"terminal_no"`
	VposID              string `json:"vpos_id"`
	MCC                 string `json:"mcc"`
	TaxID               string `json:"tax_id"`
	NationalID          string `json:"national_id"`
	Title               string `json:"title"`
	SwitchID            string `json:"switch_id"`
	City                string `json:"city"`
	Country             string `json:"country"`
	CountryISOCode      string `json:"country_isocode"`
	PostalCode          string `json:"postal_code"`
	Address             string `json:"address"`
	SubmerchantURL      string `json:"submerchant_url"`
	SubmerchantNIN      string `json:"submerchant_nin"`
}

type VposSubmerchantUpdateRequest struct {
	ExternalReferenceID string `json:"external_reference_id"`
	SubmerchantID       string `json:"submerchant_id"`
	TerminalNo          string `json:"terminal_no"`
	MCC                 string `json:"mcc"`
	TaxID               string `json:"tax_id"`
	NationalID          string `json:"national_id"`
	Title               string `json:"title"`
	SwitchID            string `json:"switch_id"`
}

type VposSubmerchant struct {
	ID                  string `json:"id,omitempty"`
	ExternalReferenceID string `json:"external_reference_id,omitempty"`
	SubmerchantID       string `json:"submerchant_id,omitempty"`
	TerminalNo          string `json:"terminal_no,omitempty"`
	VposID              string `json:"vpos_id,omitempty"`
	MCC                 string `json:"mcc,omitempty"`
	TaxID               string `json:"tax_id,omitempty"`
	NationalID          string `json:"national_id,omitempty"`
	Title               string `json:"title,omitempty"`
	SwitchID            string `json:"switch_id,omitempty"`
}

type VposSubmerchantListItem struct {
	ID                  string `json:"id,omitempty"`
	ExternalReferenceID string `json:"external_reference_id,omitempty"`
	SubmerchantID       string `json:"submerchant_id,omitempty"`
	TerminalNo          string `json:"terminal_no,omitempty"`
	VposID              string `json:"vpos_id,omitempty"`
}

type VposSubmerchantListResponse struct {
	Page       int64                     `json:"page,omitempty"`
	PerPage    int64                     `json:"per_page,omitempty"`
	Total      int64                     `json:"total,omitempty"`
	TotalPages int64                     `json:"total_pages,omitempty"`
	Rows       []VposSubmerchantListItem `json:"rows,omitempty"`
}

type VposSubmerchantMutationResponse struct {
	Code    uint64 `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Sub Organization DTO
type SubOrganizationDTO struct {
	Acquirer         string `json:"acquirer,omitempty"`
	Address          string `json:"address,omitempty"`
	ContactFirstName string `json:"contact_first_name,omitempty"`
	ContactLastName  string `json:"contact_last_name,omitempty"`
	ContactEmail     string `json:"contact_email,omitempty"`
	ContactPhone     string `json:"contact_phone,omitempty"`
	Name             string `json:"name,omitempty"`
	TaxOffice        string `json:"tax_office,omitempty"`
	TaxNumber        string `json:"tax_number,omitempty"`
	Type             string `json:"type,omitempty"`
}

// Organization Settings DTO
type OrganizationSettings struct {
	Ttl                    uint64 `json:"ttl,omitempty"`
	RetryCount             uint64 `json:"retry_count,omitempty"`
	AllowPayment           bool   `json:"allow_payment,omitempty"`
	SessionTtl             uint64 `json:"session_ttl,omitempty"`
	CustomCheckout         bool   `json:"custom_checkout,omitempty"`
	DomainAddress          string `json:"domain_address,omitempty"`
	CheckoutDomain         string `json:"checkout_domain,omitempty"`
	SubscriptionDomain     string `json:"subscription_domain,omitempty"`
	HasKycKyxTemplate      bool   `json:"has_kyc_kyx_template,omitempty"`
	HasKybKyxTemplate      bool   `json:"has_kyb_kyx_template,omitempty"`
	IsMoneyTransferPartner bool   `json:"is_money_transfer_partner,omitempty"`
}

type OrganizationCurrenciesResponse struct {
	Currencies []OrganizationCurrency `json:"currencies,omitempty"`
}

type OrganizationCurrencyPresetsResponse struct {
	Items []OrganizationCurrencyPreset `json:"items,omitempty"`
}

type OrganizationCurrencyPreset struct {
	CurrencyCode string `json:"currency_code,omitempty"`
	CurrencyUnit string `json:"currency_unit,omitempty"`
	Name         string `json:"name,omitempty"`
	MinorUnit    int    `json:"minor_unit,omitempty"`
}

type CreateOrganizationCurrencyResponse struct {
	Code     uint64               `json:"code,omitempty"`
	Message  string               `json:"message,omitempty"`
	Created  bool                 `json:"created,omitempty"`
	Currency OrganizationCurrency `json:"currency,omitempty"`
}

type VposAcquirerTemplateListResponse struct {
	Items []VposAcquirerTemplate `json:"items,omitempty"`
}

type VposAcquirerTemplate struct {
	AcquirerID     string                       `json:"acquirer_id,omitempty"`
	Name           string                       `json:"name,omitempty"`
	Prefix         string                       `json:"prefix,omitempty"`
	RequiredFields []string                     `json:"required_fields,omitempty"`
	OptionalFields []string                     `json:"optional_fields,omitempty"`
	Defaults       VposAcquirerTemplateDefaults `json:"defaults,omitempty"`
}

type VposAcquirerTemplateDefaults struct {
	PaymentMode string            `json:"payment_mode,omitempty"`
	Main        bool              `json:"main,omitempty"`
	Marketplace bool              `json:"marketplace,omitempty"`
	ForceThreeD bool              `json:"force_three_d,omitempty"`
	Credentials map[string]string `json:"credentials,omitempty"`
}

type OrganizationCurrency struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	Code         string `json:"code,omitempty"`
	Symbol       string `json:"symbol,omitempty"`
	CurrencyUnit string `json:"currency_unit,omitempty"`
}

// When you select order lite for callbacks, your callback request body will be like this.
type OrderCallbackLiteDTO struct {
	ID             string                `json:"id" example:"f0a0a1e9-69bd-4bef-b8c6-4e8c0d3a1212"`
	ReferenceID    string                `json:"orderReferenceId" example:"f0a0a1e9-69bd-4bef-b8c6-4e8c0d3a1212"`
	ConversationID string                `json:"conversationId" example:"f0a0a1e9-69bd-4bef-b8c6-4e8c0d3a1212"`
	OrderNote      string                `json:"orderNote" example:"order note"`
	IdentityNumber string                `json:"identityNumber" example:"12345678901"`
	Transactions   []string              `json:"transactions,omitempty" example:"[\"transaction_id1\", \"transaction_id2\"]"`
	OrderTerms     []string              `json:"orderTerms,omitempty" example:"[\"term1\", \"term2\"]"`
	Status         string                `json:"status" example:"1"`
	OrderPayments  []OrderPaymentItemDTO `json:"orderPayments,omitempty"`
}

type OrderPaymentItemDTO struct {
	Type       string `json:"type"`
	PaidAt     string `json:"paid_at" example:"2021-01-01 00:00:00"`
	PaidAmount string `json:"paid_amount" example:"100.00"`
}

// When you select order extended for callbacks, your callback request body will be like this.
type OrderCallbackExtendedDTO struct {
	ID          string                 `json:"id"`
	ReferenceID string                 `json:"reference_id"`
	Terms       []OrderExtendedTermDTO `json:"terms"`
}

type OrderExtendedTermDTO struct {
	ReferenceID string                        `json:"reference_id"`
	Amount      float64                       `json:"amount"`
	Required    bool                          `json:"required"`
	Status      uint64                        `json:"status"`
	Sequence    uint64                        `json:"sequence"`
	Payments    []OrderExtendedTermPaymentDTO `json:"payments"`
}

type OrderExtendedTermPaymentDTO struct {
	ID                string    `json:"id"`
	ReferenceID       string    `json:"reference_id"`
	Amount            float64   `json:"amount"`
	Date              time.Time `json:"date"`
	Status            uint64    `json:"status"`
	PaymentType       uint64    `json:"payment_type"`
	PaymentTypeString string    `json:"payment_type_string"`
}

// When you select order detail for callbacks, your callback request body will be like this.
type OrderCallbackDetailDTO struct {
	ID                 string                       `json:"id"`
	OrganizationID     string                       `json:"organization_id" example:"f0a0a1e9-69bd-4bef-b8c6-4e8c0d3a1212"`
	OrderReferenceID   string                       `json:"order_reference_id" example:"f0a0a1e9-69bd-4bef-b8c6-4e8c0d3a1212"`
	ConversationID     string                       `json:"conversation_id" example:"123456789"`
	Rule               OrderDetailRuleDTO           `json:"rule"`
	Vpos               OrderDetailVposDTO           `json:"vpos"`
	Order              OrderDetailInfoDTO           `json:"order"`
	PaymentDetails     OrderDetailPaymentDetailsDTO `json:"paymentDetails"`
	Response           string                       `json:"response" example:"example response"`
	ResponseJSON       map[string]any               `json:"response_json"`
	ResponseCallbacks  map[string]any               `json:"response_callbacks"`
	OrderPaymentStatus OrderDetailPaymentStatusDTO  `json:"order_payment_status"`
	CreatedAt          time.Time                    `json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt          time.Time                    `json:"updated_at" example:"2021-01-01T00:00:00Z"`
}

type OrderDetailInfoDTO struct {
	ReferenceID string    `json:"reference_id"`
	PaidDate    time.Time `json:"paid_date"`
	ID          string    `json:"id"`
	Status      string    `json:"status"`
}

type OrderDetailPaymentDetailsDTO struct {
	AuthCode             string `json:"auth_code"`
	BatchNo              string `json:"batch_no"`
	OrderID              string `json:"order_id"`
	ReferenceID          string `json:"reference_id"`
	PaymentTransactionId string `json:"payment_transaction_id"`
	IsThreeDS            bool   `json:"is_three_ds"`
	Installment          string `json:"installment"`
}

type OrderDetailRuleDTO struct {
	Name string `json:"name"`
	ID   string `json:"id" example:"f0a0a1e9-69bd-4bef-b8c6-4e8c0d3a1212"`
}

type OrderDetailVposDTO struct {
	ID             string `json:"id" example:"f0a0a1e9-69bd-4bef-b8c6-4e8c0d3a1212"`
	Name           string `json:"name"`
	CommissionRate string `json:"commission_rate"`
}

type OrderDetailPaymentStatusDTO struct {
	Code             string `json:"code" example:"16"`
	Message          string `json:"message" example:"NO_SUFFICIENT_FUNDS"`
	IsError          bool   `json:"is_error" example:"false"`
	MaskedPan        string `json:"masked_pan" example:"411111******1111"`
	ExpiryYear       string `json:"expiry_year" example:"2022"`
	ExpiryMonth      string `json:"expiry_month" example:"12"`
	AcquirerResponse string `json:"acquirer_response" example:"Insufficient funds"`
}

// When you select order for callbacks, your callback request body will be like this.
type OrderCallbackDTO struct {
	ID                  string               `json:"id" example:"f0a0a1e9-69bd-4bef-b8c6-4e8c0d3a1212"`
	PartialPayment      *bool                `json:"partial_payment"`
	ThreeDForce         *bool                `json:"three_d_force"`
	Locale              string               `json:"locale" example:"en"`
	ExternalReferenceID string               `json:"external_reference_id" example:"123456789"`
	ConversationID      string               `json:"conversation_id" example:"f4f4f4f4-f4f4-f4f4-f4f4-f4f4f4f4f4f4"`
	Amount              float64              `json:"amount" example:"100"`
	Fee                 float64              `json:"fee"`
	TaxAmount           float64              `json:"tax_amount" example:"100"`
	RefundAmount        float64              `json:"refund_amount" example:"100"`
	BasketID            string               `json:"basket_id" example:"13fwefsa"`
	PaymentGroup        string               `json:"payment_group" example:"2f3wdqac"`
	Buyer               OrderBuyer           `json:"buyer"`
	ShippingAddress     OrderShippingAddress `json:"shipping_address"`
	BillingAddress      OrderBillingAddress  `json:"billing_address"`
	BasketItems         []OrderBasketItem    `json:"basket_items"`
	PaidAmount          float64              `json:"paid_amount" example:"100"`
	PaidDate            time.Time            `json:"paid_date" example:"2020-01-01 00:00:00"`
	CancelDate          time.Time            `json:"cancel_date" example:"2020-01-01 00:00:00"`
	RefundDate          time.Time            `json:"refund_date" example:"2020-01-01 00:00:00"`
	EnabledInstallments []int                `json:"enabled_installments"`
	Currency            string               `json:"currency" example:"TRY"`
	Latitude            float64              `json:"latitude" example:"41.01234567"`
	Longitude           float64              `json:"longitude" example:"29.01234567"`
	Status              int                  `json:"status" example:"1"`
	ReferenceID         string               `json:"reference_id" example:"f0a0a1e9-69bd-4bef-b8c6-4e8c0d3a1212"`
	OrganizationID      string               `json:"organization_id" example:"f0a0a1e9-69bd-4bef-b8c6-4e8c0d3a1212"`
	UserID              string               `json:"user_id" example:"f0a0a1e9-69bd-4bef-b8c6-4e8c0d3a1212"`
	AcquirerID          string               `json:"acquirer_id" example:"f0a0a1e9-69bd-4bef-b8c6-4e8c0d3a1212"`
	PaymentPageID       string               `json:"payment_page_id" example:"f0a0a1e9-69bd-4bef-b8c6-4e8c0d3a1212"`
	SubMerchants        []OrderSubmerchant   `json:"sub_merchants"`
	PaymentSuccessURL   string               `json:"payment_success_url" example:"https://example.com/success"`
	PaymentFailureURL   string               `json:"payment_failure_url" example:"https://example.com/failure"`
	RedirectSuccessURL  string               `json:"redirect_success_url" example:"https://example.com/success"`
	RedirectFailureURL  string               `json:"redirect_failure_url" example:"https://example.com/failure"`
	SubOrganization     []SubOrganizationDTO `json:"sub_organization"`
	OrderType           uint64               `json:"order_type" example:"1"`
	RelatedReferenceID  string               `json:"related_reference_id" example:"f0a0a1e9-69bd-4bef-b8c6-4e8c0d3a1212"`
	SurchargeAmount     float64              `json:"surcharge_amount"`
	Descriptor          uint64               `json:"descriptor" example:"1"`
	Metadata            []OrderMetadata      `json:"metadata"`
	Note                string               `json:"note" example:"note"`
	PaymentOptions      []string             `json:"payment_options"`
	CommissionAmount    float64              `json:"commission_amount"`
	ThreeDInitializedAt time.Time            `json:"three_d_initialized_at" example:"2020-01-01 00:00:00"`
	Installment         string               `json:"installment" example:"1"`
	ScheduledAt         time.Time            `json:"scheduled_at" example:"2020-01-01 00:00:00"` //running query after this date
	PaymentMode         string               `json:"payment_mode" example:"auth or preauth"`     // auth or preauth
}

// Subscription DTOs

// SubscriptionGetRequest represents the request payload for getting a subscription
type SubscriptionGetRequest struct {
	ExternalReferenceID string `json:"external_reference_id,omitempty"`
	ReferenceID         string `json:"reference_id,omitempty"`
}

// SubscriptionCancelRequest represents the request payload for canceling a subscription
type SubscriptionCancelRequest struct {
	ExternalReferenceID string `json:"external_reference_id,omitempty"`
	ReferenceID         string `json:"reference_id,omitempty"`
}

// SubscriptionBilling represents billing information for subscription
type SubscriptionBilling struct {
	Address     string `json:"address,omitempty"`
	City        string `json:"city,omitempty"`
	ContactName string `json:"contact_name,omitempty"`
	Country     string `json:"country,omitempty"`
	VatNumber   string `json:"vat_number,omitempty"`
	ZipCode     string `json:"zip_code,omitempty"`
}

// SubscriptionUser represents user information for subscription
type SubscriptionUser struct {
	Address        string `json:"address,omitempty"`
	City           string `json:"city,omitempty"`
	Country        string `json:"country,omitempty"`
	Email          string `json:"email,omitempty"`
	FirstName      string `json:"first_name,omitempty"`
	ID             string `json:"id,omitempty"`
	IdentityNumber string `json:"identity_number,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	Phone          string `json:"phone,omitempty"`
	ZipCode        string `json:"zip_code,omitempty"`
}

// SubscriptionCreateRequest represents the request payload for creating a subscription
type SubscriptionCreateRequest struct {
	Amount              float64             `json:"amount,omitempty"`
	Billing             SubscriptionBilling `json:"billing"`
	CardID              string              `json:"card_id,omitempty"`
	Currency            string              `json:"currency,omitempty"`
	Cycle               int                 `json:"cycle,omitempty"`
	ExternalReferenceID string              `json:"external_reference_id,omitempty"`
	FailureURL          string              `json:"failure_url,omitempty"`
	Interval            string              `json:"interval,omitempty"`
	IntervalCount       int                 `json:"interval_count,omitempty"`
	// Deprecated: PaymentDate is ignored by the current server. The subscription
	// start date determines the billing anchor.
	PaymentDate int `json:"payment_date,omitempty"`
	// Deprecated: Period is the legacy day-count cadence. Use Interval and
	// IntervalCount for new subscriptions.
	Period     int              `json:"period,omitempty"`
	SuccessURL string           `json:"success_url,omitempty"`
	Title      string           `json:"title,omitempty"`
	User       SubscriptionUser `json:"user"`
}

// SubscriptionRedirectRequest represents the request payload for redirecting a subscription
type SubscriptionRedirectRequest struct {
	SubscriptionID string `json:"subscription_id,omitempty"`
}

// SubscriptionOrder represents an order within a subscription
type SubscriptionOrder struct {
	Amount      string `json:"amount,omitempty"`
	Currency    string `json:"currency,omitempty"`
	PaymentDate string `json:"payment_date,omitempty"`
	PaymentURL  string `json:"payment_url,omitempty"`
	ReferenceID string `json:"reference_id,omitempty"`
	Status      string `json:"status,omitempty"`
}

// SubscriptionDetail represents the detailed subscription information
type SubscriptionDetail struct {
	Amount              string              `json:"amount,omitempty"`
	Cycle               int                 `json:"cycle,omitempty"`
	CyclesCompleted     int                 `json:"cycles_completed,omitempty"`
	Currency            string              `json:"currency,omitempty"`
	DueDate             string              `json:"due_date,omitempty"`
	ExternalReferenceID string              `json:"external_reference_id,omitempty"`
	Interval            string              `json:"interval,omitempty"`
	IntervalCount       int                 `json:"interval_count,omitempty"`
	IsActive            bool                `json:"is_active,omitempty"`
	Orders              []SubscriptionOrder `json:"orders,omitempty"`
	PaymentDate         int                 `json:"payment_date,omitempty"`
	PaymentStatus       string              `json:"payment_status,omitempty"`
	Period              int                 `json:"period,omitempty"`
	Title               string              `json:"title,omitempty"`
}

// SubscriptionCreateResponse represents the response from creating a subscription
type SubscriptionCreateResponse struct {
	Code             int    `json:"code,omitempty"`
	Message          string `json:"message,omitempty"`
	OrderReferenceID string `json:"order_reference_id,omitempty"`
	ReferenceID      string `json:"reference_id,omitempty"`
}

// SubscriptionListItem represents a single subscription item in the list
type SubscriptionListItem struct {
	Amount              string `json:"amount,omitempty"`
	Cycle               int    `json:"cycle,omitempty"`
	CyclesCompleted     int    `json:"cycles_completed,omitempty"`
	Currency            string `json:"currency,omitempty"`
	ExternalReferenceID string `json:"external_reference_id,omitempty"`
	Interval            string `json:"interval,omitempty"`
	IntervalCount       int    `json:"interval_count,omitempty"`
	IsActive            bool   `json:"is_active,omitempty"`
	NextPaymentDate     string `json:"next_payment_date,omitempty"`
	PaymentDate         int    `json:"payment_date,omitempty"`
	PaymentStatus       string `json:"payment_status,omitempty"`
	Period              int    `json:"period,omitempty"`
	ReferenceID         string `json:"reference_id,omitempty"`
	Title               string `json:"title,omitempty"`
}

// SubscriptionRedirectResponse represents the response from redirecting a subscription
type SubscriptionRedirectResponse struct {
	URL string `json:"url,omitempty"`
}

// CardTokenizeRequest represents the request to tokenize a card
type CardTokenizeRequest struct {
	CardNumber         string `json:"card_number"`
	HolderName         string `json:"holder_name"`
	ExpiryMonth        string `json:"expiry_month"`
	ExpiryYear         string `json:"expiry_year"`
	CVV                string `json:"cvv"`
	Name               string `json:"name,omitempty"`
	Default            bool   `json:"default,omitempty"`
	Mode               int    `json:"mode,omitempty"`
	UserId             string `json:"user_id,omitempty"`
	UserEmail          string `json:"user_email,omitempty"`
	RedirectSuccessUrl string `json:"redirect_success_url,omitempty"`
	RedirectFailureUrl string `json:"redirect_failure_url,omitempty"`
}

// CardTokenizeResponse represents the response from tokenizing a card
type CardTokenizeResponse struct {
	ThreedformURL    string `json:"threedform_url,omitempty"`
	ThreedformHTML   string `json:"threedform_html,omitempty"`
	IsSuccess        bool   `json:"is_success,omitempty"`
	Message          string `json:"message,omitempty"`
	StatusCode       string `json:"status_code,omitempty"`
	ErrorCode        string `json:"error_code,omitempty"`
	OrderReferenceID string `json:"order_reference_id,omitempty"`
	CardID           string `json:"card_id,omitempty"`
}

// SavedCard represents a tokenized card stored on the account
type SavedCard struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	HolderName   string `json:"holder_name,omitempty"`
	MaskedNumber string `json:"masked_number,omitempty"`
	LastFour     string `json:"last_four,omitempty"`
	Brand        string `json:"brand,omitempty"`
	Bin          string `json:"bin,omitempty"`
	IsDefault    bool   `json:"is_default,omitempty"`
}

// ListSavedCardsResponse represents a paginated list of saved cards
type ListSavedCardsResponse struct {
	Page       int64       `json:"page,omitempty"`
	PerPage    int64       `json:"per_page,omitempty"`
	Total      int64       `json:"total,omitempty"`
	TotalPages int64       `json:"total_pages,omitempty"`
	Rows       []SavedCard `json:"rows,omitempty"`
}

// DeleteSavedCardResponse represents the response from deleting a saved card
type DeleteSavedCardResponse struct {
	Success bool   `json:"success,omitempty"`
	Message string `json:"message,omitempty"`
}

// GetOrderPaymentsRequest represents the request to list payments of an order.
// At least one of the identifiers should be provided.
type GetOrderPaymentsRequest struct {
	OrderID          string `json:"order_id,omitempty"`
	OrderReferenceID string `json:"order_reference_id,omitempty"`
	ConversationID   string `json:"conversation_id,omitempty"`
}

// OrderPayment represents a single payment made against an order
type OrderPayment struct {
	ID                            string  `json:"id,omitempty"`
	Date                          string  `json:"date,omitempty"`
	PaymentMode                   string  `json:"payment_mode,omitempty"`
	Amount                        float64 `json:"amount,omitempty"`
	CardID                        string  `json:"card_id,omitempty"`
	MaskedCard                    string  `json:"masked_card,omitempty"`
	CardHolderName                string  `json:"card_holder_name,omitempty"`
	CardType                      string  `json:"card_type,omitempty"`
	CardBrand                     string  `json:"card_brand,omitempty"`
	CardBank                      string  `json:"card_bank,omitempty"`
	CardNetwork                   string  `json:"card_network,omitempty"`
	Paid                          bool    `json:"paid,omitempty"`
	Type                          string  `json:"type,omitempty"`
	AcquirerResponse              string  `json:"acquirer_response,omitempty"`
	VposCommissionAmount          float64 `json:"vpos_commission_amount,omitempty"`
	VposCommissionRate            float64 `json:"vpos_commission_rate,omitempty"`
	VposCommissionFeeType         string  `json:"vpos_commission_fee_type,omitempty"`
	OrganizationCommissionAmount  float64 `json:"organization_commission_amount,omitempty"`
	OrganizationCommissionRate    float64 `json:"organization_commission_rate,omitempty"`
	OrganizationCommissionFeeType string  `json:"organization_commission_fee_type,omitempty"`
	TotalAmount                   float64 `json:"total_amount,omitempty"`
	AcquirerCommissionAmount      float64 `json:"acquirer_commission_amount,omitempty"`
	AcquirerCommissionRate        float64 `json:"acquirer_commission_rate,omitempty"`
	AcquirerCommissionFeeType     string  `json:"acquirer_commission_fee_type,omitempty"`
	AcquirerName                  string  `json:"acquirer_name,omitempty"`
	CustomerPaidCommissionAmount  float64 `json:"customer_paid_commission_amount,omitempty"`
	MerchantPaidCommissionAmount  float64 `json:"merchant_paid_commission_amount,omitempty"`
	MerchantNetAmount             float64 `json:"merchant_net_amount,omitempty"`
}

// GetOrderPaymentsResponse represents the list of payments made against an order
type GetOrderPaymentsResponse struct {
	Payments                          []OrderPayment `json:"payments,omitempty"`
	TotalAmount                       float64        `json:"total_amount,omitempty"`
	TotalCustomerPaidCommissionAmount float64        `json:"total_customer_paid_commission_amount,omitempty"`
	TotalMerchantPaidCommissionAmount float64        `json:"total_merchant_paid_commission_amount,omitempty"`
	TotalMerchantNetAmount            float64        `json:"total_merchant_net_amount,omitempty"`
	OrderAmount                       float64        `json:"order_amount,omitempty"`
}

// OrderListFilter selects orders for GetOrderListWithFilter. Page defaults to
// 1 and PerPage to 10 (at most 100) on the server when left at zero.
type OrderListFilter struct {
	Page               int
	PerPage            int
	StartDate          string
	EndDate            string
	OrganizationID     string
	ConversationID     string
	RelatedReferenceID string
	BuyerID            string
	// Status filters by numeric order status (see the OrderStatus* constants).
	Status *int
}

// OrderSubmerchantListItem is a submerchant row returned by GetOrderSubmerchants.
type OrderSubmerchantListItem struct {
	ID              string `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	Email           string `json:"email,omitempty"`
	SubmerchantType string `json:"submerchant_type,omitempty"`
	SubmerchantKey  string `json:"submerchant_key,omitempty"`
	Labels          string `json:"labels,omitempty"`
	Status          string `json:"status,omitempty"`
	Acquirer        string `json:"acquirer,omitempty"`
}

// OrderSubmerchantListResponse is the paginated result of GetOrderSubmerchants.
type OrderSubmerchantListResponse struct {
	Page       int64                      `json:"page,omitempty"`
	PerPage    int64                      `json:"per_page,omitempty"`
	Total      int64                      `json:"total,omitempty"`
	TotalPages int64                      `json:"total_pages,omitempty"`
	Rows       []OrderSubmerchantListItem `json:"row,omitempty"`
}

// OrderChargeRequest charges an order, for example with its saved card.
type OrderChargeRequest struct {
	OrderReferenceID string `json:"order_reference_id"`
	Retry            bool   `json:"retry,omitempty"`
}

// OrderChargeResponse is the charge result. The server sends these keys in
// camelCase, unlike most responses.
type OrderChargeResponse struct {
	IsSuccess  bool   `json:"isSuccess,omitempty"`
	Message    string `json:"message,omitempty"`
	StatusCode string `json:"statusCode,omitempty"`
	ErrorCode  string `json:"errorCode,omitempty"`
}

// SubmerchantPaymentAction requests settlement approval or disapproval for an item payment.
type SubmerchantPaymentAction struct {
	Locale               string `json:"locale"`
	ConversationID       string `json:"conversation_id"`
	PaymentTransactionID string `json:"payment_transaction_id"`
	Acquirer             string `json:"acquirer"`
}

// SubmerchantPaymentActionResponse is the settlement action result returned by the payment provider.
type SubmerchantPaymentActionResponse struct {
	OperationID          string `json:"operation_id,omitempty"`
	Status               string `json:"status,omitempty"`
	Locale               string `json:"locale,omitempty"`
	SystemTime           int64  `json:"system_time,omitempty"`
	ConversationID       string `json:"conversation_id,omitempty"`
	PaymentTransactionID string `json:"payment_transaction_id,omitempty"`
	Acquirer             string `json:"acquirer,omitempty"`
}

// SubmerchantPaymentItemUpdate changes the seller or payout for a completed marketplace item.
type SubmerchantPaymentItemUpdate struct {
	Locale               string `json:"locale"`
	ConversationID       string `json:"conversation_id"`
	PaymentTransactionID string `json:"payment_transaction_id"`
	SubMerchantKey       string `json:"sub_merchant_key"`
	SubMerchantPrice     string `json:"sub_merchant_price"`
	Acquirer             string `json:"acquirer"`
}

type SubmerchantPaymentItemUpdateResponse struct {
	OperationID          string `json:"operation_id,omitempty"`
	Status               string `json:"status,omitempty"`
	Locale               string `json:"locale,omitempty"`
	SystemTime           int64  `json:"system_time,omitempty"`
	ConversationID       string `json:"conversation_id,omitempty"`
	PaymentTransactionID string `json:"payment_transaction_id,omitempty"`
}

// OrgCreateUserRequest represents the request to create an organization user.
type OrgCreateUserRequest struct {
	Email          string `json:"email,omitempty"`
	IdentityNumber string `json:"identity_number,omitempty"`
	Phone          string `json:"phone,omitempty"`
	FirstName      string `json:"first_name,omitempty"`
	LastName       string `json:"last_name,omitempty"`
	ReferenceID    string `json:"reference_id,omitempty"`
	ConversationID string `json:"conversation_id,omitempty"`
	IsMailVerified bool   `json:"is_mail_verified,omitempty"`
}

// OrgCreateUserResponse represents the response from creating an organization user.
type OrgCreateUserResponse struct {
	Code    uint64 `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	UserID  string `json:"user_id,omitempty"`
}

// OrgUserTokenCreateRequest represents the request to create a token for an
// organization user. Expire is given in minutes.
type OrgUserTokenCreateRequest struct {
	Email               string `json:"email,omitempty"`
	Expire              uint32 `json:"expire,omitempty"`
	InvalidateOldTokens bool   `json:"invalidate_old_tokens,omitempty"`
	UserID              string `json:"user_id,omitempty"`
}

// OrgUserTokenCreateResponse represents the response from creating an
// organization user token.
type OrgUserTokenCreateResponse struct {
	Code    uint64 `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Token   string `json:"token,omitempty"`
	UserID  string `json:"user_id,omitempty"`
}
