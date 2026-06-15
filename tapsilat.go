package tapsilat

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// TapsilatAPI is the main struct for the Tapsilat API
type API struct {
	EndPoint string `json:"end_point"`
	Token    string `json:"token"`
	Timeout  time.Duration
	client   *http.Client

	currencyRefsMu     sync.RWMutex
	currencyIDsByUnit  map[string]string
	currencyCacheReady bool
}

// NewAPI creates a new TapsilatAPI struct
func NewAPI(token string) *API {
	timeout := 30 * time.Second
	return &API{
		EndPoint: "https://panel.tapsilat.dev/api/v1",
		Token:    token,
		Timeout:  timeout,
		client:   &http.Client{Timeout: timeout},
	}
}

// NewCustomAPI creates a new TapsilatAPI struct with a custom endpoint
func NewCustomAPI(endpoint, token string) *API {
	timeout := 30 * time.Second
	return &API{
		EndPoint: endpoint,
		Token:    token,
		Timeout:  timeout,
		client:   &http.Client{Timeout: timeout},
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

func (t *API) patch(ctx context.Context, path string, payload any, response any) error {
	url := t.EndPoint + path
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return t.do(req, response)
}

func (t *API) get(ctx context.Context, path string, response any) error {
	url := t.EndPoint + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	return t.do(req, response)
}

func (t *API) getWithPayload(ctx context.Context, path string, payload any, response any) error {
	url := t.EndPoint + path
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return t.do(req, response)
}

func (t *API) delete(ctx context.Context, path string, response any) error {
	url := t.EndPoint + path
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	return t.do(req, response)
}

func (t *API) deleteWithPayload(ctx context.Context, path string, payload any, response any) error {
	url := t.EndPoint + path
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	return t.do(req, response)
}

func (t *API) do(req *http.Request, response any) error {
	req.Header.Set("Authorization", "Bearer "+t.Token)
	req.Header.Set("Accept", "application/json")

	resp, err := t.client.Do(req)
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
		return newAPIError(resp.StatusCode, resp.Status, body)
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
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%d", page))
	query.Set("per_page", fmt.Sprintf("%d", perPage))
	if startDate != "" {
		query.Set("start_date", startDate)
	}
	if endDate != "" {
		query.Set("end_date", endDate)
	}
	if organizationID != "" {
		query.Set("organization_id", organizationID)
	}
	if relatedReferenceID != "" {
		query.Set("related_reference_id", relatedReferenceID)
	}
	err := t.get(ctx, "/order/list?"+query.Encode(), &response)
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

func (t *API) GetOrderPaymentDetails(ctx context.Context, payload OrderPaymentDetailDTO) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/order/payment-details", payload, &response)
	return response, err
}

func (t *API) GetOrderTransactions(ctx context.Context, referenceID string) (map[string]any, error) {
	var response map[string]any
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
func (t *API) GetOrderTerm(ctx context.Context, termReferenceID string) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/order/term?term_reference_id="+termReferenceID, &response)
	return response, err
}

func (t *API) CreateOrderTerm(ctx context.Context, term OrderPaymentTermCreateDTO) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/order/term", term, &response)
	return response, err
}

func (t *API) DeleteOrderTerm(ctx context.Context, payload OrderPaymentTermDeleteDTO) (map[string]any, error) {
	var response map[string]any
	err := t.deleteWithPayload(ctx, "/order/term", payload, &response)
	return response, err
}

func (t *API) UpdateOrderTerm(ctx context.Context, term OrderPaymentTermUpdateDTO) (map[string]any, error) {
	var response map[string]any
	err := t.patch(ctx, "/order/term", term, &response)
	return response, err
}

func (t *API) RefundOrderTerm(ctx context.Context, term OrderTermRefundRequest) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/order/term/refund", term, &response)
	return response, err
}

func (t *API) OrderTerminate(ctx context.Context, referenceID string) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/order/terminate", map[string]string{
		"reference_id": referenceID,
	}, &response)
	return response, err
}

func (t *API) OrderManualCallback(ctx context.Context, payload OrderManualCallbackDTO) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/order/callback", payload, &response)
	return response, err
}

func (t *API) OrderRelatedUpdate(ctx context.Context, payload OrderRelatedReferenceDTO) (map[string]any, error) {
	var response map[string]any
	err := t.patch(ctx, "/order/releated", payload, &response)
	return response, err
}

func (t *API) OrderAccounting(ctx context.Context, payload OrderAccountingRequest) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/order/accounting", payload, &response)
	return response, err
}

func (t *API) OrderPostAuth(ctx context.Context, payload OrderPostAuthRequest) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/order/postauth", payload, &response)
	return response, err
}

func (t *API) GetOrderPaymentDetailsByID(ctx context.Context, referenceID string) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/order/"+referenceID+"/payment-details", &response)
	return response, err
}

func (t *API) OrderPaymentOptionsUpdate(ctx context.Context, payload OrderPaymentOptionsUpdateDTO) (map[string]any, error) {
	var response map[string]any
	err := t.patch(ctx, "/order/payment-options", payload, &response)
	return response, err
}

func (t *API) SplitOrderItemPayment(ctx context.Context, payload SplitOrderItemPaymentDTO) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/order/split", payload, &response)
	return response, err
}

func (t *API) OrderCallback(ctx context.Context, id string) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/orders/"+id+"/callback", &response)
	return response, err
}

func (t *API) OrderVposQuery(ctx context.Context, id string) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/orders/"+id+"/vpos-query", &response)
	return response, err
}

func (t *API) AddBasketItem(ctx context.Context, payload AddBasketItemRequest) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/order/basket-item", payload, &response)
	return response, err
}

func (t *API) RemoveBasketItem(ctx context.Context, payload RemoveBasketItemRequest) (map[string]any, error) {
	var response map[string]any
	err := t.deleteWithPayload(ctx, "/order/basket-item", payload, &response)
	return response, err
}

func (t *API) UpdateBasketItem(ctx context.Context, payload UpdateBasketItemRequest) (map[string]any, error) {
	var response map[string]any
	err := t.patch(ctx, "/order/basket-item", payload, &response)
	return response, err
}

func (t *API) GetSystemOrderStatuses(ctx context.Context) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/system/order-statuses", &response)
	return response, err
}

func (t *API) GetSystemBasketItemTypes(ctx context.Context) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/system/basket-item-types", &response)
	return response, err
}

func (t *API) GetSystemErrorCodes(ctx context.Context) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/system/error-codes", &response)
	return response, err
}

func (t *API) GetSystemPaymentTermStatuses(ctx context.Context) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/system/payment-term-statuses", &response)
	return response, err
}

func (t *API) GetSystemProductTypes(ctx context.Context) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/system/product-types", &response)
	return response, err
}

func (t *API) GetSystemShortcutTypes(ctx context.Context) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/system/shortcut-types", &response)
	return response, err
}

func (t *API) GetSystemTransactionPaymentTypes(ctx context.Context) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/system/transaction-payment-types", &response)
	return response, err
}

func (t *API) GetSystemTransactionPurposes(ctx context.Context) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/system/transaction-purposes", &response)
	return response, err
}

func (t *API) GetSystemTransactionStatuses(ctx context.Context) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/system/transaction-statuses", &response)
	return response, err
}

func (t *API) GetOrganizationCallback(ctx context.Context) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/organization/callback", &response)
	return response, err
}

func (t *API) UpdateOrganizationCallback(ctx context.Context, payload CallbackURLDTO) (map[string]any, error) {
	var response map[string]any
	err := t.patch(ctx, "/organization/callback", payload, &response)
	return response, err
}

func (t *API) CreateOrganizationBusiness(ctx context.Context, payload OrgCreateBusinessRequest) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/organization/business/create", payload, &response)
	return response, err
}

func (t *API) GetOrganizationLimitUser(ctx context.Context, payload GetUserLimitRequest) (map[string]any, error) {
	var response map[string]any
	err := t.getWithPayload(ctx, "/organization/limit/user", payload, &response)
	return response, err
}

func (t *API) SetOrganizationLimitUser(ctx context.Context, payload SetLimitUserRequest) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/organization/limit/user", payload, &response)
	return response, err
}

func (t *API) GetOrganizationLimits(ctx context.Context) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/organization/limits", &response)
	return response, err
}

func (t *API) GetOrganizationMeta(ctx context.Context, name string) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/organization/meta/"+name, &response)
	return response, err
}

func (t *API) GetOrganizationScopes(ctx context.Context) (map[string]any, error) {
	var response map[string]any
	err := t.get(ctx, "/organization/scopes", &response)
	return response, err
}

func (t *API) CreateOrganizationUser(ctx context.Context, payload OrgCreateUserReq) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/organization/user/create", payload, &response)
	return response, err
}

func (t *API) VerifyOrganizationUser(ctx context.Context, payload OrgUserVerifyReq) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/organization/user/verify", payload, &response)
	return response, err
}

func (t *API) VerifyOrganizationUserMobile(ctx context.Context, payload OrgUserMobileVerifyReq) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/organization/user/verify-mobile", payload, &response)
	return response, err
}

func VerifyWebhook(payload, signature, secret string) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(payload))
	expectedSignature := "sha256=" + hex.EncodeToString(h.Sum(nil))
	return expectedSignature == signature
}

func (t *API) GetOrganizationSettings(ctx context.Context) (OrganizationSettings, error) {
	var response OrganizationSettings
	err := t.get(ctx, "/organization/settings", &response)
	return response, err
}

func (t *API) GetOrganizationCurrencies(ctx context.Context) (OrganizationCurrenciesResponse, error) {
	var response OrganizationCurrenciesResponse
	err := t.get(ctx, "/organization/currencies", &response)
	return response, err
}

func (t *API) ListOrganizationCurrencyPresets(ctx context.Context) (OrganizationCurrencyPresetsResponse, error) {
	var response OrganizationCurrencyPresetsResponse
	err := t.get(ctx, "/organization/currency-presets", &response)
	return response, err
}

func (t *API) CreateOrganizationCurrency(ctx context.Context, currencyCode string) (CreateOrganizationCurrencyResponse, error) {
	var response CreateOrganizationCurrencyResponse
	payload := map[string]string{
		"currency_code": strings.ToUpper(strings.TrimSpace(currencyCode)),
	}
	err := t.post(ctx, "/organization/currencies", payload, &response)
	if err == nil {
		t.invalidateCurrencyCache()
	}
	return response, err
}

func (t *API) CreateSubmerchant(ctx context.Context, payload SubmerchantCreateRequest) (SubmerchantMutationResponse, error) {
	var response SubmerchantMutationResponse
	currencyID, err := t.normalizeCurrencyID(ctx, payload.CurrencyID)
	if err != nil {
		return response, err
	}
	payload.CurrencyID = currencyID
	err = t.post(ctx, "/submerchants", payload, &response)
	return response, err
}

func (t *API) GetSubmerchant(ctx context.Context, id string) (Submerchant, error) {
	var response Submerchant
	err := t.get(ctx, "/submerchants/"+id, &response)
	return response, err
}

func (t *API) ListSubmerchants(ctx context.Context, page, perPage int) (SubmerchantListResponse, error) {
	var response SubmerchantListResponse
	path := fmt.Sprintf("/submerchants?page=%d&per_page=%d", page, perPage)
	err := t.get(ctx, path, &response)
	return response, err
}

func (t *API) UpdateSubmerchant(ctx context.Context, id string, payload SubmerchantUpdateRequest) (SubmerchantMutationResponse, error) {
	var response SubmerchantMutationResponse
	currencyID, err := t.normalizeCurrencyID(ctx, payload.CurrencyID)
	if err != nil {
		return response, err
	}
	payload.CurrencyID = currencyID
	err = t.patch(ctx, "/submerchants/"+id, payload, &response)
	return response, err
}

func (t *API) DeleteSubmerchant(ctx context.Context, id string) (SubmerchantMutationResponse, error) {
	var response SubmerchantMutationResponse
	err := t.delete(ctx, "/submerchants/"+id, &response)
	return response, err
}

func (t *API) GetSuborganizations(ctx context.Context, page, perPage int) (SuborganizationListResponse, error) {
	var response SuborganizationListResponse
	path := fmt.Sprintf("/organization/suborganizations?page=%d&per_page=%d", page, perPage)
	err := t.get(ctx, path, &response)
	return response, err
}

func (t *API) GetSuborganization(ctx context.Context, id string) (SuborganizationListItem, error) {
	var response SuborganizationListItem
	err := t.get(ctx, "/organization/suborganizations/"+id, &response)
	return response, err
}

func (t *API) GetSuborganizationDetail(ctx context.Context, id string) (SuborganizationDetail, error) {
	var response SuborganizationDetail
	err := t.get(ctx, "/organization/suborganizations/"+id, &response)
	return response, err
}

func (t *API) ListVpos(ctx context.Context, page, perPage int) (VposListResponse, error) {
	return t.ListVposWithFilter(ctx, page, perPage, VposListFilter{})
}

func (t *API) ListVposWithFilter(ctx context.Context, page, perPage int, filter VposListFilter) (VposListResponse, error) {
	var response VposListResponse
	query := url.Values{}
	query.Set("page", fmt.Sprintf("%d", page))
	query.Set("per_page", fmt.Sprintf("%d", perPage))
	if filter.SuborganizationID != "" {
		query.Set("suborganization_id", filter.SuborganizationID)
	}
	path := "/vpos?" + query.Encode()
	err := t.get(ctx, path, &response)
	return response, err
}

func (t *API) CreateVpos(ctx context.Context, payload VposCreateRequest) (VposMutationResponse, error) {
	var response VposMutationResponse
	currencies, err := t.normalizeCurrencyIDs(ctx, payload.Currencies)
	if err != nil {
		return response, err
	}
	payload.Currencies = currencies
	err = t.post(ctx, "/vpos", payload, &response)
	return response, err
}

func (t *API) GetVpos(ctx context.Context, id string) (Vpos, error) {
	var response Vpos
	err := t.get(ctx, "/vpos/"+id, &response)
	return response, err
}

func (t *API) UpdateVpos(ctx context.Context, id string, payload VposUpdateRequest) (VposMutationResponse, error) {
	var response VposMutationResponse
	currencies, err := t.normalizeCurrencyIDs(ctx, payload.Currencies)
	if err != nil {
		return response, err
	}
	payload.Currencies = currencies
	err = t.patch(ctx, "/vpos/"+id, payload, &response)
	return response, err
}

func (t *API) DeleteVpos(ctx context.Context, id string) (VposMutationResponse, error) {
	var response VposMutationResponse
	err := t.delete(ctx, "/vpos/"+id, &response)
	return response, err
}

func (t *API) ListVposAcquirers(ctx context.Context) (VposAcquirerListResponse, error) {
	var response VposAcquirerListResponse
	err := t.get(ctx, "/vpos/acquirers", &response)
	return response, err
}

func (t *API) ListCardSchemes(ctx context.Context) (CardSchemeListResponse, error) {
	var response CardSchemeListResponse
	err := t.get(ctx, "/vpos/card-schemes", &response)
	return response, err
}

func (t *API) ListVposAcquirerTemplates(ctx context.Context) (VposAcquirerTemplateListResponse, error) {
	var response VposAcquirerTemplateListResponse
	err := t.get(ctx, "/vpos/acquirer-templates", &response)
	return response, err
}

func (t *API) ListVposSubmerchants(ctx context.Context, page, perPage int, vposID, externalReferenceID string) (VposSubmerchantListResponse, error) {
	var response VposSubmerchantListResponse
	path := fmt.Sprintf("/vpos-submerchant?page=%d&per_page=%d", page, perPage)
	if vposID != "" {
		path += "&vpos_id=" + vposID
	}
	if externalReferenceID != "" {
		path += "&external_reference_id=" + externalReferenceID
	}
	err := t.get(ctx, path, &response)
	return response, err
}

func (t *API) CreateVposSubmerchant(ctx context.Context, payload VposSubmerchantCreateRequest) (VposSubmerchantMutationResponse, error) {
	var response VposSubmerchantMutationResponse
	err := t.post(ctx, "/vpos-submerchant", payload, &response)
	return response, err
}

func (t *API) GetVposSubmerchant(ctx context.Context, id string) (VposSubmerchant, error) {
	var response VposSubmerchant
	err := t.get(ctx, "/vpos-submerchant/"+id, &response)
	return response, err
}

func (t *API) UpdateVposSubmerchant(ctx context.Context, id string, payload VposSubmerchantUpdateRequest) (VposSubmerchantMutationResponse, error) {
	var response VposSubmerchantMutationResponse
	err := t.patch(ctx, "/vpos-submerchant/"+id, payload, &response)
	return response, err
}

func (t *API) DeleteVposSubmerchant(ctx context.Context, id string) (VposSubmerchantMutationResponse, error) {
	var response VposSubmerchantMutationResponse
	err := t.delete(ctx, "/vpos-submerchant/"+id, &response)
	return response, err
}

func (t *API) GetSuborganizationBySubmerchant(ctx context.Context, submerchantID string) (SubmerchantSuborganizationMapping, error) {
	var response SubmerchantSuborganizationMapping
	err := t.get(ctx, "/submerchants/"+submerchantID+"/suborganization", &response)
	return response, err
}

func (t *API) GetSubmerchantBySuborganization(ctx context.Context, suborganizationID string) (SuborganizationSubmerchantMapping, error) {
	var response SuborganizationSubmerchantMapping
	err := t.get(ctx, "/organization/suborganizations/"+suborganizationID+"/submerchant", &response)
	return response, err
}

// Subscription methods

func (t *API) GetSubscription(ctx context.Context, payload SubscriptionGetRequest) (SubscriptionDetail, error) {
	var response SubscriptionDetail
	err := t.post(ctx, "/subscription", payload, &response)
	return response, err
}

func (t *API) CancelSubscription(ctx context.Context, payload SubscriptionCancelRequest) error {
	var response map[string]any
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
