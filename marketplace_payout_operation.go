package tapsilat

import (
	"context"
	"net/http"
	"net/url"
	"time"
)

type MarketplacePayoutOperationInput struct {
	Kind            string    `json:"kind"`
	IdempotencyKey  string    `json:"idempotency_key"`
	Revision        time.Time `json:"revision"`
	SellerReference string    `json:"seller_reference,omitempty"`
	NetAmount       string    `json:"net_amount,omitempty"`
}

type MarketplacePayoutOperation struct {
	ID           string    `json:"id"`
	PayoutID     string    `json:"payout_id"`
	Kind         string    `json:"kind"`
	Status       string    `json:"status"`
	Attempts     uint64    `json:"attempts"`
	ErrorMessage string    `json:"error_message,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// RunMarketplacePayoutOperation changes allocation or withdraws approval; inspect Status before treating it as complete.
func (t *API) RunMarketplacePayoutOperation(ctx context.Context, payoutID string, input MarketplacePayoutOperationInput) (MarketplacePayoutOperation, error) {
	var response MarketplacePayoutOperation
	err := t.marketplaceRequest(ctx, http.MethodPost, "/submerchant-payouts/"+url.PathEscape(payoutID)+"/operations", input, &response)
	return response, err
}

// GetMarketplacePayoutOperation returns the persisted operation receipt scoped to the API token.
func (t *API) GetMarketplacePayoutOperation(ctx context.Context, id string) (MarketplacePayoutOperation, error) {
	var response MarketplacePayoutOperation
	err := t.marketplaceRequest(ctx, http.MethodGet, "/marketplace/payout-operations/"+url.PathEscape(id), nil, &response)
	return response, err
}

// ReconcileMarketplacePayoutOperation checks the provider result without repeating a mutation.
func (t *API) ReconcileMarketplacePayoutOperation(ctx context.Context, id string) (MarketplacePayoutOperation, error) {
	var response MarketplacePayoutOperation
	err := t.marketplaceRequest(ctx, http.MethodPost, "/marketplace/payout-operations/"+url.PathEscape(id)+"/reconcile", nil, &response)
	return response, err
}

// RetryMarketplacePayoutOperation repeats an unresolved change only after verifying the previous provider state.
func (t *API) RetryMarketplacePayoutOperation(ctx context.Context, id string) (MarketplacePayoutOperation, error) {
	var response MarketplacePayoutOperation
	err := t.marketplaceRequest(ctx, http.MethodPost, "/marketplace/payout-operations/"+url.PathEscape(id)+"/retry", nil, &response)
	return response, err
}
