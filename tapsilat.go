package tapsilat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// TapsilatAPI is the main struct for the Tapsilat API
type API struct {
	EndPoint string `json:"end_point"`
	Token    string `json:"token"`
	Timeout  time.Duration
}

// NewAPI creates a new TapsilatAPI struct
func NewAPI(token string) *API {
	return &API{
		EndPoint: "https://panel.tapsilat.dev/api/v1",
		Token:    token,
		Timeout:  30 * time.Second,
	}
}

// NewCustomAPI creates a new TapsilatAPI struct with a custom endpoint
func NewCustomAPI(endpoint, token string) *API {
	return &API{
		EndPoint: endpoint,
		Token:    token,
		Timeout:  30 * time.Second,
	}
}

func (t *API) post(ctx context.Context, path string, payload any, response any) error {
	url := t.EndPoint + path
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return t.do(req, response)
}

func (t *API) get(ctx context.Context, path string, response any) error {
	url := t.EndPoint + path
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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

func (t *API) CreateOrder(ctx context.Context, payload Order) (OrderResponse, error) {
	var response OrderResponse

	// Validate GSM number if provided
	if payload.Buyer.GsmNumber != "" {
		cleanedGSM, err := ValidateGSMNumber(payload.Buyer.GsmNumber)
		if err != nil {
			return response, err
		}
		payload.Buyer.GsmNumber = cleanedGSM
	}

	err := t.post(ctx, "/order/create", payload, &response)
	if err != nil {
		return response, err
	}

	// If order creation successful and we have a reference ID, get the checkout URL
	if response.ReferenceID != "" {
		checkoutURL, err := t.GetCheckoutURL(ctx, response.ReferenceID)
		if err == nil && checkoutURL != "" {
			response.CheckoutURL = checkoutURL
		}
		// Don't return error if checkout URL fetch fails, just continue without it
	}

	return response, nil
}

func (t *API) GetOrder(ctx context.Context, orderReferenceID string) (OrderDetail, error) {
	var response OrderDetail
	err := t.get(ctx, "/order/"+orderReferenceID, &response)
	return response, err
}

func (t *API) GetOrderByConversationID(ctx context.Context, conversationID string) (OrderDetail, error) {
	var response OrderDetail
	err := t.get(ctx, "/order/conversation/"+conversationID, &response)
	return response, err
}

func (t *API) GetOrders(ctx context.Context, page, perPage, buyerID string) (PaginatedData, error) {
	var response PaginatedData
	path := fmt.Sprintf("/order/list?page=%s&per_page=%s", page, perPage)
	if buyerID != "" {
		path += "&buyer_id=" + buyerID
	}
	err := t.get(ctx, path, &response)
	return response, err
}

func (t *API) GetOrderList(ctx context.Context, page, perPage int, startDate, endDate, organizationID, relatedReferenceID string) (PaginatedData, error) {
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

	err := t.get(ctx, path, &response)
	return response, err
}

func (t *API) GetOrderSubmerchants(ctx context.Context, page, perPage int) (PaginatedData, error) {
	var response PaginatedData
	path := fmt.Sprintf("/order/submerchants?page=%d&per_page=%d", page, perPage)
	err := t.get(ctx, path, &response)
	return response, err
}

func (t *API) GetCheckoutURL(ctx context.Context, referenceID string) (string, error) {
	order, err := t.GetOrder(ctx, referenceID)
	if err != nil {
		return "", err
	}
	return order.CheckoutURL, nil
}

func (t *API) GetOrderStatus(ctx context.Context, orderReferenceID string) (OrderStatus, error) {
	var orderStatus OrderStatus
	err := t.get(ctx, "/order/"+orderReferenceID+"/status", &orderStatus)
	return orderStatus, err
}

func (t *API) GetOrderPaymentDetails(ctx context.Context, referenceID, conversationID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	path := "/order/payment-details"

	if referenceID != "" {
		path += "?reference_id=" + referenceID
	} else if conversationID != "" {
		path += "?conversation_id=" + conversationID
	}

	err := t.get(ctx, path, &response)
	return response, err
}

func (t *API) GetOrderTransactions(ctx context.Context, referenceID string) ([]map[string]interface{}, error) {
	var response []map[string]interface{}
	err := t.get(ctx, "/order/"+referenceID+"/transactions", &response)
	return response, err
}

func (t *API) CancelOrder(ctx context.Context, payload CancelOrder) (RefundCancelOrderResponse, error) {
	var response RefundCancelOrderResponse
	err := t.post(ctx, "/order/cancel", payload, &response)
	return response, err
}

func (t *API) RefundOrder(ctx context.Context, payload RefundOrder) (RefundCancelOrderResponse, error) {
	var response RefundCancelOrderResponse
	err := t.post(ctx, "/order/refund", payload, &response)
	return response, err
}

func (t *API) RefundAllOrder(ctx context.Context, referenceID string) (RefundCancelOrderResponse, error) {
	payload := RefundOrder{
		ReferenceID: referenceID,
	}
	return t.RefundOrder(ctx, payload)
}

// Payment Terms methods
func (t *API) GetOrderTerm(ctx context.Context, termReferenceID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.get(ctx, "/order/term/"+termReferenceID, &response)
	return response, err
}

func (t *API) CreateOrderTerm(ctx context.Context, term OrderPaymentTermCreateDTO) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.post(ctx, "/order/term/create", term, &response)
	return response, err
}

func (t *API) DeleteOrderTerm(ctx context.Context, orderID, termReferenceID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.post(ctx, "/order/term/delete", map[string]string{
		"order_id":          orderID,
		"term_reference_id": termReferenceID,
	}, &response)
	return response, err
}

func (t *API) UpdateOrderTerm(ctx context.Context, term OrderPaymentTermUpdateDTO) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.post(ctx, "/order/term/update", term, &response)
	return response, err
}

func (t *API) RefundOrderTerm(ctx context.Context, term OrderTermRefundRequest) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.post(ctx, "/order/term/refund", term, &response)
	return response, err
}

func (t *API) OrderTerminate(ctx context.Context, referenceID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.post(ctx, "/order/terminate", map[string]string{
		"reference_id": referenceID,
	}, &response)
	return response, err
}

func (t *API) OrderManualCallback(ctx context.Context, referenceID, conversationID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	payload := map[string]string{
		"reference_id": referenceID,
	}
	if conversationID != "" {
		payload["conversation_id"] = conversationID
	}
	err := t.post(ctx, "/order/manual-callback", payload, &response)
	return response, err
}

func (t *API) OrderRelatedUpdate(ctx context.Context, referenceID, relatedReferenceID string) (map[string]interface{}, error) {
	var response map[string]interface{}
	err := t.post(ctx, "/order/related-update", map[string]string{
		"reference_id":         referenceID,
		"related_reference_id": relatedReferenceID,
	}, &response)
	return response, err
}

func (t *API) GetOrganizationSettings(ctx context.Context) (OrganizationSettings, error) {
	var response OrganizationSettings
	err := t.get(ctx, "/organization/settings", &response)
	return response, err
}

// Subscription methods

func (t *API) GetSubscription(ctx context.Context, payload SubscriptionGetRequest) (SubscriptionDetail, error) {
	var response SubscriptionDetail
	err := t.post(ctx, "/subscription", payload, &response)
	return response, err
}

func (t *API) CancelSubscription(ctx context.Context, payload SubscriptionCancelRequest) error {
	var response map[string]interface{}
	err := t.post(ctx, "/subscription/cancel", payload, &response)
	return err
}

func (t *API) CreateSubscription(ctx context.Context, payload SubscriptionCreateRequest) (SubscriptionCreateResponse, error) {
	var response SubscriptionCreateResponse
	err := t.post(ctx, "/subscription/create", payload, &response)
	return response, err
}

func (t *API) ListSubscriptions(ctx context.Context, page, perPage int) (PaginatedData, error) {
	var response PaginatedData
	path := fmt.Sprintf("/subscription/list?page=%d&per_page=%d", page, perPage)
	err := t.get(ctx, path, &response)
	return response, err
}

func (t *API) RedirectSubscription(ctx context.Context, payload SubscriptionRedirectRequest) (SubscriptionRedirectResponse, error) {
	var response SubscriptionRedirectResponse
	err := t.post(ctx, "/subscription/redirect", payload, &response)
	return response, err
}
