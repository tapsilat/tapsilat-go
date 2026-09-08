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

func TestMarketplaceAccountSelectionEncoding(t *testing.T) {
	for _, mode := range []string{"omitted", "empty", "explicit", "legacy"} {
		t.Run(mode, func(t *testing.T) {
			request := tapsilat.MarketplaceSubmerchantCreateRequest{CurrencyID: "9f4050e8-1111-4f25-b4ef-aaaaaaaaaaaa"}
			if mode == "empty" {
				accounts := []tapsilat.MarketplaceAccountTarget{}
				request.ProviderAccounts = &accounts
			}
			if mode == "explicit" {
				accounts := []tapsilat.MarketplaceAccountTarget{{VposID: "8f4050e8-1111-4f25-b4ef-aaaaaaaaaaaa", Environment: "TEST"}}
				request.ProviderAccounts = &accounts
			}
			if mode == "legacy" {
				request.VposID = "8f4050e8-1111-4f25-b4ef-aaaaaaaaaaaa"
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, "/api/v2/submerchants", r.URL.Path)
				var body map[string]json.RawMessage
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				_, accounts := body["provider_accounts"]
				_, legacy := body["vpos_id"]
				require.Equal(t, mode == "empty" || mode == "explicit", accounts)
				require.Equal(t, mode == "legacy", legacy)
				w.Header().Set("Content-Type", "application/json")
				require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"id": "seller", "provisionings": []any{}}))
			}))
			defer server.Close()
			_, err := tapsilat.NewCustomAPI(server.URL+"/api/v1", "test").CreateMarketplaceSubmerchant(context.Background(), request)
			require.NoError(t, err)
		})
	}
}
