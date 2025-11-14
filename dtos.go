package tapsilat

import (
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
	{23, "Partially paid"},
	{24, "Terminated"},
	{25, "Card tokenization"},
	{26, "Pre authorized"},
	{27, "Disputed"},
	{28, "Partially disputed"},
	{29, "Suspect"},
}

// OrderMetadata represents metadata key-value pairs
type OrderMetadata struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// OrderPaymentTerm represents payment term for installments
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
	PfSubMerchant       OrderPfSubMerchant   `json:"pf_sub_merchant"`
	ThreeDForce         bool                 `json:"three_d_force"`
	EnabledInstallments []int                `json:"enabled_installments"`
	PaymentOptions      []string             `json:"payment_options"`
	Metadata            []OrderMetadata      `json:"metadata"`
	PaymentTerms        []OrderPaymentTerm   `json:"payment_terms"`
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
}

type OrderDetail struct {
	Locale              string                 `json:"locale"`
	Error               string                 `json:"error"`
	Code                int                    `json:"code"`
	ReferenceID         string                 `json:"reference_id"`
	Amount              string                 `json:"amount"`
	Total               string                 `json:"total"`
	PaidAmount          string                 `json:"paid_amount"`
	RefundedAmount      string                 `json:"refunded_amount"`
	CreatedAt           string                 `json:"created_at"`
	Currency            string                 `json:"currency"`
	Status              int32                  `json:"status"`
	StatusEnum          string                 `json:"status_enum"`
	Buyer               OrderBuyer             `json:"buyer"`
	ShippingAddress     OrderShippingAddress   `json:"shipping_address"`
	CheckoutDesign      OrderCheckoutDesignDTO `json:"checkout_design"`
	BillingAddress      OrderBillingAddress    `json:"billing_address"`
	BasketItems         []OrderBasketItem      `json:"basket_items"`
	Submerchants        []OrderSubmerchant     `json:"submerchants"`
	PaymentTerms        []OrderPaymentTermDTO  `json:"payment_terms"`
	ItemPayments        []OrderItemPayment     `json:"item_payments"`
	PaymentFailureUrl   string                 `json:"payment_failure_url" example:"https://www.example.com/payment/failure"`
	PaymentSuccessUrl   string                 `json:"payment_success_url" example:"https://www.example.com/payment/success"`
	CheckoutURL         string                 `json:"checkout_url" example:"https://www.example.com/payment/checkout"`
	ConversationID      string                 `json:"conversation_id" example:"123456789"`
	PaymentOptions      []string               `json:"payment_options" example:"credit_card,bank_transfer,cash"`
	ExternalReferenceID string                 `json:"external_reference_id" example:"ext-123456"`
	MCC                 string                 `json:"mcc" example:"5411"`
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
	BirdthDate          string `json:"birdth_date"`
	Title               string `json:"title,omitempty"`
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
	Payer            *OrderBasketItemPayer `json:"payer,omitempty"`
}

// OrderBasketItemPayer represents the payer information for basket items
type OrderBasketItemPayer struct {
	Address     string `json:"address,omitempty"`
	ReferenceID string `json:"reference_id,omitempty"`
	TaxOffice   string `json:"tax_office,omitempty"`
	Title       string `json:"title,omitempty"`
	Type        string `json:"type,omitempty"`
	VAT         string `json:"vat,omitempty"`
}

type OrderSubmerchant struct {
	Amount              float64 `json:"amount"`
	OrderBasketItemID   string  `json:"order_basket_item_id"`
	MerchantReferenceID string  `json:"merchant_reference_id"`
}

type OrderCheckoutDesign struct {
	Logo                 string `json:"logo"`
	InputBackgroundColor string `json:"input_background_color"`
	PayButtonColor       string `json:"pay_button_color"`
	InputTextColor       string `json:"input_text_color"`
	LabelTextColor       string `json:"label_text_color"`
	LeftBackgroundColor  string `json:"left_background_color"`
	RightBackgroundColor string `json:"right_background_color"`
	TextColor            string `json:"text_color"`
	OrderDetailHtml      string `json:"order_detail_html"`
	RedirectUrl          string `json:"redirect_url"`
}
type OrderItemPayment struct {
	Id               string  `json:"id,omitempty"`
	Amount           float64 `json:"amount,omitempty"`
	PaidDate         string  `json:"paid_date,omitempty"`
	MaskedBin        string  `json:"masked_bin,omitempty"`
	CardBrand        string  `json:"card_brand,omitempty"`
	RefundedAmount   float64 `json:"refunded_amount,omitempty"`
	RefundableAmount float64 `json:"refundable_amount,omitempty"`
	Refunded         bool    `json:"refunded,omitempty"`
	Status           uint64  `json:"status,omitempty" example:"1"`
} // @name OrderItemPayment

type OrderBillingAddress struct {
	BillingType  string `json:"billing_type" example:"PERSONAL"`   // PERSONAL, BUSINESS
	Citizenship  string `json:"citizenship" example:"TR"`          // ISO 3166-1 alpha-2 country code
	Title        string `json:"title" example:"MonoPayments Inc."` // Legal title
	TaxOffice    string `json:"tax_office" example:"Merter"`
	Address      string `json:"address" example:"Istanbul"`
	ZipCode      string `json:"zip_code" example:"34000"`
	City         string `json:"city" example:"Istanbul"`
	District     string `json:"district" example:"Uskudar"`
	Country      string `json:"country" example:"Turkey"`
	ContactName  string `json:"contact_name" example:"John Doe"`
	ContactPhone string `json:"contact_phone" example:"+905555555555"`
	VatNumber    string `json:"vat_number" example:"1234567890"` // Identification number if billing type is PERSONAL
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
	ReferenceID string `json:"reference_id"`
	Amount      string `json:"amount"`
	Error       string `json:"error"`
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
	Page       int64       `json:"page,omitempty" example:"1"`
	PerPage    int64       `json:"per_page,omitempty" example:"10"`
	Total      int64       `json:"total,omitempty" example:"100"`
	TotalPages int         `json:"total_pages,omitempty" example:"10"`
	Rows       interface{} `json:"rows,omitempty" swaggertype:"array,string" example:"object,object2"`
	Error      string      `json:"error"`
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
	OrderDetailHtml      string `json:"order_detail_html" example:"<html><body><h1>Order Detail</h1></body></html>"`
} // @name OrderCheckoutDesignDTO

// Payment Term DTOs for create, update and refund operations
type OrderPaymentTermCreateDTO struct {
	OrderReferenceID string  `json:"order_reference_id"`
	Amount           float64 `json:"amount"`
	DueDate          string  `json:"due_date"`
	Required         bool    `json:"required"`
	Data             string  `json:"data,omitempty"`
	TermSequence     *int    `json:"term_sequence,omitempty"`
}

type OrderPaymentTermUpdateDTO struct {
	TermReferenceID string   `json:"term_reference_id"`
	Amount          *float64 `json:"amount,omitempty"`
	DueDate         string   `json:"due_date,omitempty"`
	Required        *bool    `json:"required,omitempty"`
	Data            string   `json:"data,omitempty"`
}

type OrderTermRefundRequest struct {
	TermReferenceID string   `json:"term_reference_id"`
	Amount          *float64 `json:"amount,omitempty"`
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
	Ttl                uint64 `json:"ttl,omitempty"`
	RetryCount         uint64 `json:"retry_count,omitempty"`
	AllowPayment       bool   `json:"allow_payment,omitempty"`
	SessionTtl         uint64 `json:"session_ttl,omitempty"`
	CustomCheckout     bool   `json:"custom_checkout,omitempty"`
	DomainAddress      string `json:"domain_address,omitempty"`
	CheckoutDomain     string `json:"checkout_domain,omitempty"`
	SubscriptionDomain string `json:"subscription_domain,omitempty"`
}
