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

// GetOrderStatusByStr returns the order status id by string
func GetOrderStatusByStr(status string) int {
	for _, v := range OrderStatuesMap {
		if v.Status == status {
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
	ResponseJSON       map[string]interface{}       `json:"response_json"`
	ResponseCallbacks  map[string]interface{}       `json:"response_callbacks"`
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
	Billing             SubscriptionBilling `json:"billing,omitempty"`
	CardID              string              `json:"card_id,omitempty"`
	Currency            string              `json:"currency,omitempty"`
	Cycle               int                 `json:"cycle,omitempty"`
	ExternalReferenceID string              `json:"external_reference_id,omitempty"`
	FailureURL          string              `json:"failure_url,omitempty"`
	PaymentDate         int                 `json:"payment_date,omitempty"`
	Period              int                 `json:"period,omitempty"`
	SuccessURL          string              `json:"success_url,omitempty"`
	Title               string              `json:"title,omitempty"`
	User                SubscriptionUser    `json:"user,omitempty"`
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
	Currency            string              `json:"currency,omitempty"`
	DueDate             string              `json:"due_date,omitempty"`
	ExternalReferenceID string              `json:"external_reference_id,omitempty"`
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
	Currency            string `json:"currency,omitempty"`
	ExternalReferenceID string `json:"external_reference_id,omitempty"`
	IsActive            bool   `json:"is_active,omitempty"`
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
