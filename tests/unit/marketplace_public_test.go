package unit_test

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMarketplacePublicMethods(t *testing.T) {
	paths := []string{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.Method+" "+r.URL.Path)
		require.Equal(t, "Bearer scoped-token", r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		if strings.HasSuffix(r.URL.Path, "/actions") {
			var input tapsilat.MarketplaceAccountAction
			require.NoError(t, json.NewDecoder(r.Body).Decode(&input))
			require.Equal(t, "resubmit", input.Action)
			require.Equal(t, "TEST", input.Environment)
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"status": "failed", "retryable": false, "error_message": "missing contact"}))
			return
		}
		require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"rows": []any{}, "id": "seller", "contracts": map[string]any{}}))
	}))
	defer server.Close()
	api := tapsilat.NewCustomAPI(server.URL+"/api/v1", "scoped-token")
	ctx := context.Background()
	_, err := api.ListMarketplaceSubmerchants(ctx, 1, 20)
	require.NoError(t, err)
	_, err = api.GetMarketplaceSubmerchant(ctx, "seller")
	require.NoError(t, err)
	_, err = api.GetMarketplaceSellerAccounts(ctx, "seller")
	require.NoError(t, err)
	result, err := api.RunMarketplaceSellerAccountAction(ctx, "seller", "account", tapsilat.MarketplaceAccountAction{Action: "resubmit", Environment: "TEST", Revision: strings.Repeat("a", 64)})
	require.NoError(t, err)
	require.Equal(t, "failed", result.Status)
	require.False(t, result.Retryable)
	require.Equal(t, []string{"GET /api/v2/submerchants", "GET /api/v2/submerchants/seller", "GET /api/v2/submerchants/seller/providers", "POST /api/v2/submerchants/seller/providers/account/actions"}, paths)
}
