package tapsilat

import "time"

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
}

type Order struct {
	Locale            string               `json:"locale"`
	Amount            float64              `json:"amount"`
	TaxAmount         float64              `json:"tax_amount"`
	Currency          string               `json:"currency"`
	ConversationID    string               `json:"conversation_id"`
	Buyer             OrderBuyer           `json:"buyer"`
	ShippingAddress   OrderShippingAddress `json:"shipping_address"`
	BillingAddress    OrderBillingAddress  `json:"billing_address"`
	BasketItems       []OrderBasketItem    `json:"basket_items"`
	Submerchants      []OrderSubmerchant   `json:"submerchants"`
	CheckoutDesign    OrderCheckoutDesign  `json:"checkout_design"`
	PaymentMethods    bool                 `json:"payment_methods"`
	PaymentFailureUrl string               `json:"payment_failure_url"`
	PaymentSuccessUrl string               `json:"payment_success_url"`
	PfSubMerchant     OrderPfSubMerchant   `json:"pf_sub_merchant"`
	ThreeDForce       bool                 `json:"three_d_force"`
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
}

type OrderDetail struct {
	Locale            string                 `json:"locale"`
	Error             string                 `json:"error"`
	Code              int                    `json:"code"`
	ReferenceID       string                 `json:"reference_id"`
	Amount            string                 `json:"amount"`
	Total             string                 `json:"total"`
	PaidAmount        string                 `json:"paid_amount"`
	RefundedAmount    string                 `json:"refunded_amount"`
	CreatedAt         string                 `json:"created_at"`
	Currency          string                 `json:"currency"`
	Status            int32                  `json:"status"`
	StatusEnum        string                 `json:"status_enum"`
	Buyer             OrderBuyer             `json:"buyer"`
	ShippingAddress   OrderShippingAddress   `json:"shipping_address"`
	CheckoutDesign    OrderCheckoutDesignDTO `json:"checkout_design"`
	BillingAddress    OrderBillingAddress    `json:"billing_address"`
	BasketItems       []OrderBasketItem      `json:"basket_items"`
	Submerchants      []OrderSubmerchant     `json:"submerchants"`
	PaymentTerms      []OrderPaymentTermDTO  `json:"payment_terms"`
	ItemPayments      []OrderItemPayment     `json:"item_payments"`
	PaymentFailureUrl string                 `json:"payment_failure_url" example:"https://www.example.com/payment/failure"`
	PaymentSuccessUrl string                 `json:"payment_success_url" example:"https://www.example.com/payment/success"`
	CheckoutURL       string                 `json:"checkout_url" example:"https://www.example.com/payment/checkout"`
	ConversationID    string                 `json:"conversation_id" example:"123456789"`
	PaymentOptions    []string               `json:"payment_options" example:"credit_card,bank_transfer,cash"`
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
	Id               string             `json:"id" example:"123456789"`
	Price            float64            `json:"price" example:"100.00"`
	Name             string             `json:"name" example:"Product Name"`
	Category1        string             `json:"category1" example:"Category 1"`
	Category2        string             `json:"category2" example:"Category 2"`
	ItemType         string             `json:"item_type" example:"PHYSICAL"`
	Status           uint64             `json:"status" example:"1"`
	RefundedAmount   float64            `json:"refunded_amount" example:"0.00"`
	RefundableAmount float64            `json:"refundable_amount" example:"0.00"`
	PaidAmount       float64            `json:"paid_amount" example:"0.00"`
	PaidableAmount   float64            `json:"paidable_amount" example:"0.00"`
	Coupon           string             `json:"coupon" example:"coupon"`
	CouponDiscount   float64            `json:"coupon_discount" example:"0.00"`
	ItemPayments     []OrderItemPayment `json:"item_payments"`
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
