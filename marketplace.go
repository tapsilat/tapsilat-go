package tapsilat

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type MarketplaceSeller struct {
	ID                string    `json:"id"`
	RoutingReference  string    `json:"routing_reference"`
	Revision          time.Time `json:"revision"`
	Name              string    `json:"name"`
	Email             string    `json:"email"`
	GsmNumber         string    `json:"gsm_number"`
	Address           string    `json:"address"`
	IBAN              string    `json:"iban"`
	CurrencyID        string    `json:"currency_id"`
	Type              string    `json:"sub_merchant_type"`
	ExternalID        string    `json:"sub_merchant_external_id"`
	IdentityNumber    string    `json:"identity_number"`
	TaxNumber         string    `json:"tax_number"`
	TaxOffice         string    `json:"tax_office"`
	LegalCompanyTitle string    `json:"legal_company_title"`
	ContactName       string    `json:"contact_name"`
	ContactSurname    string    `json:"contact_surname"`
	Locale            string    `json:"locale"`
	ConversationID    string    `json:"conversation_id"`
	Status            string    `json:"status"`
	Labels            []string  `json:"labels"`
	ApprovalMode      *string   `json:"approval_mode"`
	ReleasePolicyID   *string   `json:"release_policy_id"`
}

type MarketplaceSellerList struct {
	Rows       []MarketplaceSeller `json:"rows"`
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	Total      int64               `json:"total"`
	TotalPages int                 `json:"total_pages"`
}

// MarketplaceProfileUpdate replaces the saved profile at Revision; account selection must be omitted.
type MarketplaceProfileUpdate struct {
	MarketplaceSubmerchantCreateRequest
	Revision time.Time `json:"revision"`
}
type MarketplaceProfileUpdateResponse struct {
	ApprovalRequired bool `json:"approval_required"`
}
type MarketplaceSellerContract struct {
	Supported       bool            `json:"supported"`
	Remote          bool            `json:"remote"`
	Fields          map[string]bool `json:"fields"`
	ImmutableFields []string        `json:"immutable_fields"`
	Currencies      []string        `json:"currencies"`
}
type MarketplaceContracts struct {
	Contracts map[string]MarketplaceSellerContract `json:"contracts"`
}
type MarketplaceContractFilter struct {
	CurrencyID string
	SellerType string
	Operation  string
	VposIDs    []string
}
type MarketplaceSynchronization struct {
	Status        string     `json:"status"`
	Operation     string     `json:"operation"`
	ErrorMessage  string     `json:"error_message,omitempty"`
	NeedsUpdate   bool       `json:"needs_update"`
	Attempts      uint64     `json:"attempts"`
	VerifiedAt    *time.Time `json:"verified_at,omitempty"`
	MissingFields []string   `json:"missing_fields"`
}
type MarketplaceSellerAccount struct {
	VposID                     string                               `json:"vpos_id"`
	Name                       string                               `json:"name"`
	Provider                   string                               `json:"provider"`
	Environment                string                               `json:"environment"`
	Eligible                   bool                                 `json:"eligible"`
	Mapped                     bool                                 `json:"mapped"`
	Status                     string                               `json:"status"`
	Attempts                   uint64                               `json:"attempts"`
	UpdatedAt                  *time.Time                           `json:"updated_at,omitempty"`
	ErrorMessage               string                               `json:"error_message,omitempty"`
	Retryable                  bool                                 `json:"retryable"`
	Revision                   string                               `json:"revision"`
	SellerRevision             string                               `json:"seller_revision"`
	AccountRevision            string                               `json:"account_revision"`
	Actions                    []string                             `json:"actions"`
	Contracts                  map[string]MarketplaceSellerContract `json:"contracts"`
	ContractError              string                               `json:"contract_error,omitempty"`
	Sync                       *MarketplaceSynchronization          `json:"sync,omitempty"`
	RegistrationProfileChanged bool                                 `json:"registration_profile_changed"`
	RegistrationMissingFields  []string                             `json:"registration_missing_fields,omitempty"`
}
type MarketplaceSellerAccounts struct {
	Rows []MarketplaceSellerAccount `json:"rows"`
}

// MarketplaceAccountAction confirms a listed action using the account's latest revision and environment.
// Retry reuses the original snapshot; resubmit registers the current saved profile.
type MarketplaceAccountAction struct {
	Action      string `json:"action"`
	Environment string `json:"environment"`
	Revision    string `json:"revision"`
}

func (t *API) marketplaceRequest(ctx context.Context, method, path string, payload, response any) error {
	endpoint, err := versionedEndpoint(t.EndPoint, "v2")
	if err != nil {
		return err
	}
	var body []byte
	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint+path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return t.do(req, response)
}

func (t *API) ListMarketplaceSubmerchants(ctx context.Context, page, perPage int) (MarketplaceSellerList, error) {
	var response MarketplaceSellerList
	err := t.marketplaceRequest(ctx, http.MethodGet, "/submerchants?page="+strconv.Itoa(page)+"&per_page="+strconv.Itoa(perPage), nil, &response)
	return response, err
}
func (t *API) GetMarketplaceSubmerchant(ctx context.Context, id string) (MarketplaceSeller, error) {
	var response MarketplaceSeller
	err := t.marketplaceRequest(ctx, http.MethodGet, "/submerchants/"+url.PathEscape(id), nil, &response)
	return response, err
}
func (t *API) UpdateMarketplaceSubmerchant(ctx context.Context, id string, input MarketplaceProfileUpdate) (MarketplaceProfileUpdateResponse, error) {
	var response MarketplaceProfileUpdateResponse
	currency, err := t.normalizeCurrencyID(ctx, input.CurrencyID)
	if err != nil {
		return response, err
	}
	input.CurrencyID = currency
	err = t.marketplaceRequest(ctx, http.MethodPatch, "/submerchants/"+url.PathEscape(id), input, &response)
	return response, err
}
func (t *API) GetMarketplaceSellerContracts(ctx context.Context, filter MarketplaceContractFilter) (MarketplaceContracts, error) {
	var response MarketplaceContracts
	currency, err := t.normalizeCurrencyID(ctx, filter.CurrencyID)
	if err != nil {
		return response, err
	}
	query := url.Values{"currency_id": {currency}, "seller_type": {filter.SellerType}, "operation": {filter.Operation}, "vpos_id": filter.VposIDs}
	err = t.marketplaceRequest(ctx, http.MethodGet, "/marketplace/seller-contracts?"+query.Encode(), nil, &response)
	return response, err
}
func (t *API) GetMarketplaceSellerAccounts(ctx context.Context, id string) (MarketplaceSellerAccounts, error) {
	var response MarketplaceSellerAccounts
	err := t.marketplaceRequest(ctx, http.MethodGet, "/submerchants/"+url.PathEscape(id)+"/providers", nil, &response)
	return response, err
}
func (t *API) RunMarketplaceSellerAccountAction(ctx context.Context, id, vposID string, input MarketplaceAccountAction) (MarketplaceSubmerchantProvisioning, error) {
	var response MarketplaceSubmerchantProvisioning
	err := t.marketplaceRequest(ctx, http.MethodPost, "/submerchants/"+url.PathEscape(id)+"/providers/"+url.PathEscape(vposID)+"/actions", input, &response)
	return response, err
}
