package tapsilat

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

// TapsilatAPI is the main struct for the Tapsilat API
type API struct {
	EndPoint string `json:"end_point"`
	Token    string `json:"token"`
}

// NewAPI creates a new TapsilatAPI struct
func NewAPI(token string) *API {
	return &API{
		EndPoint: "https://acquiring.tapsilat.com",
		Token:    token,
	}
}

// NewAPIWithEndpoint creates a new TapsilatAPI struct with a custom endpoint
func NewCustomAPI(endpoint, token string) *API {
	return &API{
		EndPoint: endpoint,
		Token:    token,
	}
}

func (t *API) post(path string, payload interface{}, response interface{}) error {
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

func (t *API) get(path string, payload interface{}, response interface{}) error {
	url := t.EndPoint + path
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	return t.do(req, response)
}

func (t *API) do(req *http.Request, response interface{}) error {
	req.Header.Set("Authorization", "Bearer "+t.Token)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(response)
}

func (t *API) CreateOrder(payload Order) (OrderResponse, error) {
	var response OrderResponse
	err := t.post("/order/create", payload, &response)
	return response, err
}

func (t *API) GetOrder(order_reference_id string) (Order, error) {
	var order Order
	err := t.get("/order/"+order_reference_id, nil, &order)
	return order, err
}

func (t *API) GetOrders(page, per_page string) (PaginatedData, error) {
	var response PaginatedData
	err := t.get("/order/list?page="+page+"&per_page="+per_page, nil, &response)
	return response, err
}

func (t *API) GetOrderStatus(order_reference_id string) (OrderStatus, error) {
	var orderStatus OrderStatus
	err := t.get("/order/"+order_reference_id+"/status", nil, &orderStatus)
	return orderStatus, err
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
