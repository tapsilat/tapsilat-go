package tapsilat

type Order struct {
	Locale            string               `json:"locale" binding:"required"`
	Amount            float64              `json:"amount" binding:"required"`
	Currency          string               `json:"currency" binding:"required"`
	ConservationID    string               `json:"conservation_id"`
	Buyer             OrderBuyer           `json:"buyer"`
	ShippingAddress   OrderShippingAddress `json:"shipping_address"`
	BillingAddress    OrderBillingAddress  `json:"billing_address"`
	BasketItems       []OrderBaseketItem   `json:"basket_items"`
	Submerchants      []OrderSubmerchant   `json:"submerchants"`
	CheckoutDesign    OrderCheckoutDesign  `json:"checkout_design"`
	PaymentMethods    bool                 `json:"payment_methods"`
	PaymentFailureUrl string               `json:"payment_failure_url"`
	PaymentSuccessUrl string               `json:"payment_success_url"`
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

type OrderBaseketItem struct {
	Id               string  `json:"id"`
	Price            float32 `json:"price"`
	Name             string  `json:"name"`
	Category1        string  `json:"category1"`
	Category2        string  `json:"category2"`
	ItemType         string  `json:"item_type"`
	SubMerchantKey   string  `json:"sub_merchant_key"`
	SubMerchantPrice string  `json:"sub_merchant_price"`
}

type OrderSubmerchant struct {
	Amount              float32 `json:"amount"`
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
}

type OrderStatus struct {
	Status string `json:"status"`
}

type RefundOrder struct {
	ReferenceID string `json:"reference_id" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
}

type CancelOrder struct {
	ReferenceID string `json:"reference_id" binding:"required"`
}

type RefundCancelOrderResponse struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	IsSuccess bool   `json:"is_success"`
}

type PaginatedData struct {
	Page       int64       `json:"page,omitempty" example:"1"`
	PerPage    int64       `json:"per_page,omitempty" example:"10"`
	Total      int64       `json:"total,omitempty" example:"100"`
	TotalPages int         `json:"total_pages,omitempty" example:"10"`
	Rows       interface{} `json:"rows,omitempty" swaggertype:"array,string" example:"object,object2"`
}
