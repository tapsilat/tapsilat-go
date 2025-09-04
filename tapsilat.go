package tapsilat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// TapsilatAPI is the main struct for the Tapsilat API
type API struct {
	EndPoint string `json:"end_point"`
	Token    string `json:"token"`
	Timeout  time.Duration
}

// NewAPI creates a new TapsilatAPI struct with environment variable support
func NewAPI(token string) *API {
	// Load .env file if it exists
	_ = godotenv.Load()

	// Use provided token or fall back to environment variable
	if token == "" {
		token = os.Getenv("TAPSILAT_TOKEN")
	}

	endpoint := os.Getenv("TAPSILAT_BASE_URL")
	if endpoint == "" {
		endpoint = "https://acquiring.tapsilat.com/api/v1"
	}

	return &API{
		EndPoint: endpoint,
		Token:    token,
		Timeout:  30 * time.Second,
	}
}

// NewAPIFromEnv creates a new API instance using only environment variables
func NewAPIFromEnv() *API {
	return NewAPI("")
}

// NewAPIWithEndpoint creates a new TapsilatAPI struct with a custom endpoint
func NewCustomAPI(endpoint, token string) *API {
	return &API{
		EndPoint: endpoint,
		Token:    token,
		Timeout:  30 * time.Second,
	}
}

func (t *API) post(path string, payload any, response any) error {
	url := t.EndPoint + path
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return t.do(req, response)
}

func (t *API) get(path string, response any) error {
	url := t.EndPoint + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	return t.do(req, response)
}

func (t *API) do(req *http.Request, response any) error {
	req.Header.Set("Authorization", "Bearer "+t.Token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: t.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	decode := json.NewDecoder(bytes.NewReader(body))
	decode.UseNumber()
	if err := decode.Decode(response); err != nil {
		return err
	}
	return nil
}

func (t *API) CreateOrder(payload Order) (OrderResponse, error) {
	var response OrderResponse

	// Validate GSM number if provided
	if payload.Buyer.GsmNumber != "" {
		cleanedGSM, err := ValidateGSMNumber(payload.Buyer.GsmNumber)
		if err != nil {
			return response, err
		}
		payload.Buyer.GsmNumber = cleanedGSM
	}

	err := t.post("/order/create", payload, &response)
	if err != nil {
		return response, err
	}

	// If order creation successful and we have a reference ID, get the checkout URL
	if response.ReferenceID != "" {
		checkoutURL, err := t.GetCheckoutURL(response.ReferenceID)
		if err == nil && checkoutURL != "" {
			response.CheckoutURL = checkoutURL
		}
		// Don't return error if checkout URL fetch fails, just continue without it
	}

	return response, nil
}

func (t *API) GetOrder(orderReferenceID string) (OrderDetail, error) {
	var response OrderDetail
	err := t.get("/order/"+orderReferenceID, &response)
	return response, err
}

func (t *API) GetOrderByConversationID(conversationID string) (OrderDetail, error) {
	var response OrderDetail
	err := t.get("/order/conversation/"+conversationID, &response)
	return response, err
}

func (t *API) GetOrders(page, perPage, buyerID string) (PaginatedData, error) {
	var response PaginatedData
	path := fmt.Sprintf("/order/list?page=%s&per_page=%s", page, perPage)
	if buyerID != "" {
		path += "&buyer_id=" + buyerID
	}
	err := t.get(path, &response)
	return response, err
}

func (t *API) GetOrderList(page, perPage int, startDate, endDate, organizationID, relatedReferenceID string) (PaginatedData, error) {
	var response PaginatedData
	path := fmt.Sprintf("/order/list?page=%d&per_page=%d", page, perPage)

	if startDate != "" {
		path += "&start_date=" + startDate
	}
	if endDate != "" {
		path += "&end_date=" + endDate
	}
	if organizationID != "" {
		path += "&organization_id=" + organizationID
	}
	if relatedReferenceID != "" {
		path += "&related_reference_id=" + relatedReferenceID
	}

	err := t.get(path, &response)
	return response, err
}

func (t *API) GetOrderSubmerchants(page, perPage int) (PaginatedData, error) {
	var response PaginatedData
	path := fmt.Sprintf("/order/submerchants?page=%d&per_page=%d", page, perPage)
	err := t.get(path, &response)
	return response, err
}

func (t *API) GetCheckoutURL(referenceID string) (string, error) {
	order, err := t.GetOrder(referenceID)
	if err != nil {
		return "", err
	}
	return order.CheckoutURL, nil
}

func (t *API) GetOrderStatus(orderReferenceID string) (OrderStatus, error) {
	var orderStatus OrderStatus
	err := t.get("/order/"+orderReferenceID+"/status", &orderStatus)
	return orderStatus, err
}

func (t *API) GetOrderPaymentDetails(referenceID, conversationID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	path := "/order/payment-details"

	if referenceID != "" {
		path += "?reference_id=" + referenceID
	} else if conversationID != "" {
		path += "?conversation_id=" + conversationID
	}

	err := t.get(path, &response)
	return response, err
}

func (t *API) GetOrderTransactions(referenceID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.get("/order/"+referenceID+"/transactions", &response)
	return response, err
}

func (t *API) CancelOrder(payload CancelOrder) (RefundCancelOrderResponse, error) {
	var response RefundCancelOrderResponse
	err := t.post("/order/cancel", payload, &response)
	return response, err
}

func (t *API) RefundOrder(payload RefundOrder) (RefundCancelOrderResponse, error) {
	var response RefundCancelOrderResponse
	err := t.post("/order/refund", payload, &response)
	return response, err
}

func (t *API) RefundAllOrder(referenceID string) (RefundCancelOrderResponse, error) {
	payload := RefundOrder{
		ReferenceID: referenceID,
	}
	return t.RefundOrder(payload)
}

// Payment Terms methods
func (t *API) GetOrderTerm(termReferenceID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.get("/order/term/"+termReferenceID, &response)
	return response, err
}

func (t *API) CreateOrderTerm(term OrderPaymentTermCreateDTO) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.post("/order/term/create", term, &response)
	return response, err
}

func (t *API) DeleteOrderTerm(orderID, termReferenceID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.post("/order/term/delete", map[string]string{
		"order_id":          orderID,
		"term_reference_id": termReferenceID,
	}, &response)
	return response, err
}

func (t *API) UpdateOrderTerm(term OrderPaymentTermUpdateDTO) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.post("/order/term/update", term, &response)
	return response, err
}

func (t *API) RefundOrderTerm(term OrderTermRefundRequest) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.post("/order/term/refund", term, &response)
	return response, err
}

func (t *API) OrderTerminate(referenceID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.post("/order/terminate", map[string]string{
		"reference_id": referenceID,
	}, &response)
	return response, err
}

func (t *API) OrderManualCallback(referenceID, conversationID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	payload := map[string]string{
		"reference_id": referenceID,
	}
	if conversationID != "" {
		payload["conversation_id"] = conversationID
	}
	err := t.post("/order/manual-callback", payload, &response)
	return response, err
}

func (t *API) OrderRelatedUpdate(referenceID, relatedReferenceID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.post("/order/related-update", map[string]string{
		"reference_id":         referenceID,
		"related_reference_id": relatedReferenceID,
	}, &response)
	return response, err
}
