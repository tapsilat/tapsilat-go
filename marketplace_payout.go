package tapsilat

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type MarketplacePayout struct {
	ID                 string     `json:"id"`
	Provider           string     `json:"provider"`
	OrganizationID     string     `json:"organization_id"`
	SubmerchantID      string     `json:"submerchant_id"`
	SubmerchantName    string     `json:"submerchant_name"`
	OrderID            string     `json:"order_id"`
	OrderReferenceID   string     `json:"order_reference_id"`
	OrderDisplay       string     `json:"order_display"`
	Amount             string     `json:"amount"`
	Currency           string     `json:"currency"`
	ApprovalMode       string     `json:"approval_mode"`
	ApprovalRuleID     string     `json:"approval_rule_id,omitempty"`
	ApprovalRuleName   string     `json:"approval_rule_name,omitempty"`
	ReleasePolicyID    string     `json:"release_policy_id,omitempty"`
	ReleasePolicyName  string     `json:"release_policy_name,omitempty"`
	EligibilityStatus  string     `json:"eligibility_status"`
	RequiredEvent      string     `json:"required_event,omitempty"`
	ItemID             string     `json:"item_id,omitempty"`
	BusinessEligibleAt *time.Time `json:"business_eligible_at,omitempty"`
	EligibleAt         *time.Time `json:"eligible_at,omitempty"`
	ReleaseStatus      string     `json:"release_status"`
	Status             string     `json:"status"`
	ProviderReference  string     `json:"provider_reference,omitempty"`
	Attempts           uint64     `json:"attempts"`
	ReleasedAt         *time.Time `json:"released_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	ErrorMessage       string     `json:"error_message,omitempty"`
}
type MarketplacePayoutFilter struct {
	OrderID, SubmerchantID, VposID, Provider, Status, ReleaseStatus, ApprovalMode, Currency, Eligibility, Search string
	MinAmount, MaxAmount, StartDate, EndDate                                                                     string
	Page, PerPage                                                                                                int
}
type MarketplacePayoutList struct {
	Rows       []MarketplacePayout `json:"rows"`
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	Total      int64               `json:"total"`
	TotalPages int                 `json:"total_pages"`
}
type MarketplacePayoutApproval struct {
	IDs []string `json:"ids"`
}

// MarketplacePayoutAcceptance acknowledges local approval; payout completion is reported by the payout status.
type MarketplacePayoutAcceptance struct {
	Status string `json:"status"`
}
type MarketplaceApprovalPolicy struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Priority       uint64   `json:"priority"`
	Enabled        bool     `json:"enabled"`
	Provider       string   `json:"provider"`
	MinAmount      *string  `json:"min_amount"`
	MaxAmount      *string  `json:"max_amount"`
	CurrencyID     *string  `json:"currency_id"`
	SubmerchantIDs []string `json:"submerchant_ids"`
	LabelsAny      []string `json:"labels_any"`
	Mode           string   `json:"mode"`
}
type MarketplaceReleasePolicy struct {
	ID                    string `json:"id"`
	Name                  string `json:"name"`
	Enabled               bool   `json:"enabled"`
	Default               bool   `json:"default"`
	MinimumDelaySeconds   int64  `json:"minimum_delay_seconds"`
	ExplicitNotBefore     string `json:"explicit_not_before"`
	RequiredExternalEvent string `json:"required_external_event"`
}
type MarketplaceApprovalPolicies struct {
	Rows []MarketplaceApprovalPolicy `json:"rows"`
}
type MarketplaceReleasePolicies struct {
	Rows []MarketplaceReleasePolicy `json:"rows"`
}

func (t *API) ListMarketplacePayouts(ctx context.Context, f MarketplacePayoutFilter) (MarketplacePayoutList, error) {
	var response MarketplacePayoutList
	q := url.Values{}
	for key, value := range map[string]string{"order_id": f.OrderID, "submerchant_id": f.SubmerchantID, "vpos_id": f.VposID, "provider": f.Provider,
		"status": f.Status, "release_status": f.ReleaseStatus, "approval_mode": f.ApprovalMode, "currency": f.Currency, "eligibility": f.Eligibility,
		"search": f.Search, "min_amount": f.MinAmount, "max_amount": f.MaxAmount, "start_date": f.StartDate, "end_date": f.EndDate} {
		if value != "" {
			q.Set(key, value)
		}
	}
	q.Set("page", strconv.Itoa(f.Page))
	q.Set("per_page", strconv.Itoa(f.PerPage))
	err := t.marketplaceRequest(ctx, http.MethodGet, "/submerchant-payouts?"+q.Encode(), nil, &response)
	return response, err
}
func (t *API) GetMarketplacePayout(ctx context.Context, id string) (MarketplacePayout, error) {
	var response MarketplacePayout
	err := t.marketplaceRequest(ctx, http.MethodGet, "/submerchant-payouts/"+url.PathEscape(id), nil, &response)
	return response, err
}

// ApproveMarketplacePayouts submits eligible payouts to the existing approval workflow without bypassing release gates.
func (t *API) ApproveMarketplacePayouts(ctx context.Context, input MarketplacePayoutApproval) (MarketplacePayoutAcceptance, error) {
	var response MarketplacePayoutAcceptance
	err := t.marketplaceRequest(ctx, http.MethodPost, "/submerchant-payouts/approve", input, &response)
	return response, err
}
func (t *API) ListMarketplaceApprovalPolicies(ctx context.Context) (MarketplaceApprovalPolicies, error) {
	var response MarketplaceApprovalPolicies
	err := t.marketplaceRequest(ctx, http.MethodGet, "/marketplace/approval-policies", nil, &response)
	return response, err
}
func (t *API) GetMarketplaceApprovalPolicy(ctx context.Context, id string) (MarketplaceApprovalPolicy, error) {
	var response MarketplaceApprovalPolicy
	err := t.marketplaceRequest(ctx, http.MethodGet, "/marketplace/approval-policies/"+url.PathEscape(id), nil, &response)
	return response, err
}
func (t *API) ListMarketplaceReleasePolicies(ctx context.Context) (MarketplaceReleasePolicies, error) {
	var response MarketplaceReleasePolicies
	err := t.marketplaceRequest(ctx, http.MethodGet, "/marketplace/release-policies", nil, &response)
	return response, err
}
func (t *API) GetMarketplaceReleasePolicy(ctx context.Context, id string) (MarketplaceReleasePolicy, error) {
	var response MarketplaceReleasePolicy
	err := t.marketplaceRequest(ctx, http.MethodGet, "/marketplace/release-policies/"+url.PathEscape(id), nil, &response)
	return response, err
}
