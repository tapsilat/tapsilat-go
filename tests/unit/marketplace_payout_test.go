package unit_test

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMarketplacePayoutAPI(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodGet {
			require.Equal(t, "order", r.URL.Query().Get("order_id"))
			require.Equal(t, "waiting", r.URL.Query().Get("eligibility"))
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"rows": []any{map[string]any{"id": "payout", "amount": "80.00", "eligibility_status": "waiting", "required_event": "service_completed"}}}))
		} else {
			require.Equal(t, "/api/v2/submerchant-payouts/approve", r.URL.Path)
			var input tapsilat.MarketplacePayoutApproval
			require.NoError(t, json.NewDecoder(r.Body).Decode(&input))
			require.Equal(t, []string{"payout"}, input.IDs)
			w.WriteHeader(http.StatusAccepted)
			require.NoError(t, json.NewEncoder(w).Encode(map[string]string{"status": "accepted"}))
		}
	}))
	defer server.Close()
	api := tapsilat.NewCustomAPI(server.URL+"/api/v1", "test")
	ctx := context.Background()
	rows, err := api.ListMarketplacePayouts(ctx, tapsilat.MarketplacePayoutFilter{OrderID: "order", Eligibility: "waiting"})
	require.NoError(t, err)
	require.Len(t, rows.Rows, 1)
	require.Equal(t, "service_completed", rows.Rows[0].RequiredEvent)
	accepted, err := api.ApproveMarketplacePayouts(ctx, tapsilat.MarketplacePayoutApproval{IDs: []string{"payout"}})
	require.NoError(t, err)
	require.Equal(t, "accepted", accepted.Status)
}
