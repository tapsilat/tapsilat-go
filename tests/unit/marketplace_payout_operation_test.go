package unit_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

func TestMarketplacePayoutOperationTransport(t *testing.T) {
	revision := time.Now().UTC()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "Bearer test", r.Header.Get("Authorization"))
		switch r.URL.Path {
		case "/api/v2/submerchant-payouts/payout/operations":
			require.Equal(t, http.MethodPost, r.Method)
			var input tapsilat.MarketplacePayoutOperationInput
			require.NoError(t, json.NewDecoder(r.Body).Decode(&input))
			require.Equal(t, "update_item", input.Kind)
			require.Equal(t, "60.25", input.NetAmount)
			require.True(t, revision.Equal(input.Revision))
		case "/api/v2/marketplace/payout-operations/op":
			require.Equal(t, http.MethodGet, r.Method)
		case "/api/v2/marketplace/payout-operations/op/reconcile", "/api/v2/marketplace/payout-operations/op/retry":
			require.Equal(t, http.MethodPost, r.Method)
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"id": "op", "payout_id": "payout", "status": "unknown", "attempts": 1}))
	}))
	defer server.Close()
	api := tapsilat.NewCustomAPI(server.URL+"/api/v1", "test")
	ctx := context.Background()
	input := tapsilat.MarketplacePayoutOperationInput{Kind: "update_item", IdempotencyKey: "change-1", Revision: revision, SellerReference: "seller", NetAmount: "60.25"}
	op, err := api.RunMarketplacePayoutOperation(ctx, "payout", input)
	require.NoError(t, err)
	require.Equal(t, "unknown", op.Status)
	_, err = api.GetMarketplacePayoutOperation(ctx, op.ID)
	require.NoError(t, err)
	_, err = api.ReconcileMarketplacePayoutOperation(ctx, op.ID)
	require.NoError(t, err)
	_, err = api.RetryMarketplacePayoutOperation(ctx, op.ID)
	require.NoError(t, err)
}
