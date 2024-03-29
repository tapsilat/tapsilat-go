package tapsilat

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
}

type OrderDetail struct {
	Locale          string               `json:"locale"`
	ReferenceID     string               `json:"reference_id"`
	Amount          string               `json:"amount"`
	CreatedAt       string               `json:"created_at"`
	Currency        string               `json:"currency"`
	Status          int32                `json:"status"`
	Buyer           OrderBuyer           `json:"buyer"`
	ShippingAddress OrderShippingAddress `json:"shipping_address"`
	BillingAddress  OrderBillingAddress  `json:"billing_address"`
	BasketItems     []OrderBasketItem    `json:"basket_items"`
	Submerchants    []OrderSubmerchant   `json:"submerchants"`
}

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
}

type OrderBasketItem struct {
	Id               string  `json:"id" example:"BI101"`
	Price            float64 `json:"price" example:"0.3"`
	Quantity         uint64  `json:"quantity" example:"1"`
	Name             string  `json:"name" example:"Binocular"`
	Category1        string  `json:"category1" example:"Collectibles"`
	Category2        string  `json:"category2" example:"Accessories"`
	ItemType         string  `json:"item_type" example:"PHYSICAL"`
	SubMerchantKey   string  `json:"sub_merchant_key" example:"sub merchant key"`
	SubMerchantPrice string  `json:"sub_merchant_price" example:"0.27"`
	Coupon           string  `json:"coupon" example:"coupon"`
	CouponDiscount   float64 `json:"coupon_discount" example:"0.1"`
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

type OrderBillingAddress struct {
	Address     string `json:"address"`
	ZipCode     string `json:"zip_code"`
	City        string `json:"city"`
	Country     string `json:"country"`
	ContactName string `json:"contact_name"`
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
