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

func TestMarketplacePolicyMutationAPI(t *testing.T) {
	var requests []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r.Method+" "+r.URL.Path)
		require.Equal(t, "Bearer test", r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v2/marketplace/approval-settings" && r.Method == http.MethodGet:
			require.NoError(t, json.NewEncoder(w).Encode(map[string]string{"default_mode": "manual"}))
		case r.URL.Path == "/api/v2/marketplace/approval-settings" && r.Method == http.MethodPatch:
			var input tapsilat.MarketplaceApprovalSettings
			require.NoError(t, json.NewDecoder(r.Body).Decode(&input))
			require.Equal(t, "auto", input.DefaultMode)
			require.NoError(t, json.NewEncoder(w).Encode(input))
		case r.URL.Path == "/api/v2/marketplace/approval-policies" && r.Method == http.MethodPost:
			var input tapsilat.MarketplaceApprovalPolicyInput
			require.NoError(t, json.NewDecoder(r.Body).Decode(&input))
			require.Equal(t, "manual sellers", input.Name)
			w.WriteHeader(http.StatusCreated)
		case r.URL.Path == "/api/v2/marketplace/approval-policies/rule" && r.Method == http.MethodPatch:
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/api/v2/marketplace/approval-policies/rule" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/api/v2/marketplace/release-policies" && r.Method == http.MethodPost:
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"id": "release", "name": "next day"}))
		case r.URL.Path == "/api/v2/marketplace/release-policies/release" && r.Method == http.MethodPatch:
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"id": "release", "name": "updated"}))
		case r.URL.Path == "/api/v2/marketplace/release-policies/release" && r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	api := tapsilat.NewCustomAPI(server.URL+"/api/v1", "test")
	ctx := context.Background()
	settings, err := api.GetMarketplaceApprovalSettings(ctx)
	require.NoError(t, err)
	require.Equal(t, "manual", settings.DefaultMode)
	settings, err = api.UpdateMarketplaceApprovalSettings(ctx, tapsilat.MarketplaceApprovalSettings{DefaultMode: "auto"})
	require.NoError(t, err)
	require.Equal(t, "auto", settings.DefaultMode)

	enabled := true
	approval := tapsilat.MarketplaceApprovalPolicyInput{Name: "manual sellers", Enabled: &enabled, Mode: "manual"}
	require.NoError(t, api.CreateMarketplaceApprovalPolicy(ctx, approval))
	require.NoError(t, api.UpdateMarketplaceApprovalPolicy(ctx, "rule", approval))
	require.NoError(t, api.DeleteMarketplaceApprovalPolicy(ctx, "rule"))

	release := tapsilat.MarketplaceReleasePolicyInput{Name: "next day", Enabled: true}
	created, err := api.CreateMarketplaceReleasePolicy(ctx, release)
	require.NoError(t, err)
	require.Equal(t, "release", created.ID)
	updated, err := api.UpdateMarketplaceReleasePolicy(ctx, "release", release)
	require.NoError(t, err)
	require.Equal(t, "updated", updated.Name)
	require.NoError(t, api.DeleteMarketplaceReleasePolicy(ctx, "release"))
	require.Len(t, requests, 8)
}
