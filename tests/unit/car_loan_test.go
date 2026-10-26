package unit_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

type carLoanCall struct {
	method string
	path   string
	auth   string
	body   map[string]any
}

func carLoanServer(t *testing.T, status int, response string) (*httptest.Server, *carLoanCall) {
	t.Helper()
	call := &carLoanCall{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		call.method = r.Method
		call.path = r.URL.Path
		call.auth = r.Header.Get("Authorization")
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &call.body)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(server.Close)
	return server, call
}

func TestCarLoanCheckLimit(t *testing.T) {
	server, call := carLoanServer(t, http.StatusOK, `{
		"response_code": "00",
		"response_message": "OK",
		"organization_id": "org-1",
		"reference_id": "ref-1",
		"trace_id": "",
		"vehicle_value": 750000,
		"banks": [{"bank_name": "Bank", "has_error": false, "offers": []}],
		"offers": [{"id": "offer-1", "credential_id": "cred-1", "available_limit": 420000, "min_amount": 50000, "max_term": 36, "loan_to_value_ratio": 0.7, "has_error": false}]
	}`)
	api := tapsilat.NewCustomAPI(server.URL, "token_123")

	score := int32(1500)
	res, err := api.CarLoanCheckLimit(context.Background(), tapsilat.CarLoanCheckLimitRequest{
		ProviderType:    3,
		ReferenceID:     "ref-1",
		Tckn:            "11111111111",
		PhoneNumber:     "5551112233",
		Dob:             "1990-01-01",
		RequestedAmount: 400000,
		SellerID:        "seller-1",
		EmployeeID:      "employee-1",
		FindeksScore:    &score,
		Vehicle:         &tapsilat.CarLoanVehicle{Model: "ADRIA CORAL", Year: 2024, Plate: "34ABC123"},
	})
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, call.method)
	assert.Equal(t, "/loan/car/limit", call.path)
	assert.Equal(t, "Bearer token_123", call.auth)
	assert.Equal(t, "ref-1", call.body["reference_id"])
	assert.Equal(t, 400000.0, call.body["requested_amount"])
	assert.Equal(t, "seller-1", call.body["seller_id"])
	assert.Equal(t, 1500.0, call.body["findeks_score"])
	assert.NotContains(t, call.body, "organization_id")
	assert.NotContains(t, call.body, "marital_status")
	vehicle := call.body["vehicle"].(map[string]any)
	assert.Equal(t, "ADRIA CORAL", vehicle["model"])
	assert.Equal(t, "34ABC123", vehicle["plate"])

	assert.Equal(t, "00", res.ResponseCode)
	assert.Equal(t, 750000.0, res.VehicleValue)
	require.Len(t, res.Banks, 1)
	require.Len(t, res.Offers, 1)
	assert.Equal(t, "offer-1", res.Offers[0].ID)
	assert.Equal(t, "cred-1", res.Offers[0].CredentialID)
	assert.Equal(t, 420000.0, res.Offers[0].AvailableLimit)
	assert.Equal(t, int32(36), res.Offers[0].MaxTerm)
	assert.Equal(t, 0.7, res.Offers[0].LoanToValueRatio)
}

func TestCarLoanPaymentPlanReadsInt64TermsWrittenAsStrings(t *testing.T) {
	server, call := carLoanServer(t, http.StatusOK, `{
		"response_code": "00",
		"loan_amount": 350000,
		"plans": [{"payment_plan_id": "plan-1", "term": "36", "installment_amount": 12000.5, "total_payment_amount": 432018}]
	}`)
	api := tapsilat.NewCustomAPI(server.URL, "token_123")

	amount := 350000.0
	res, err := api.CarLoanPaymentPlan(context.Background(), tapsilat.CarLoanPaymentPlanRequest{
		CredentialID:    "cred-1",
		OfferID:         "offer-1",
		Term:            36,
		RequestedAmount: &amount,
	})
	require.NoError(t, err)
	assert.Equal(t, "/loan/car/payment-plan", call.path)
	assert.Equal(t, 350000.0, call.body["requested_amount"])
	require.Len(t, res.Plans, 1)
	assert.Equal(t, int64(36), res.Plans[0].Term)
	assert.Equal(t, 12000.5, res.Plans[0].InstallmentAmount)
}

func TestCarLoanPaymentPlanOmitsUnsetAmount(t *testing.T) {
	server, call := carLoanServer(t, http.StatusOK, `{"response_code":"00","plans":[]}`)
	api := tapsilat.NewCustomAPI(server.URL, "token_123")

	_, err := api.CarLoanPaymentPlan(context.Background(), tapsilat.CarLoanPaymentPlanRequest{CredentialID: "cred-1", OfferID: "offer-1", Term: 36})
	require.NoError(t, err)
	assert.NotContains(t, call.body, "requested_amount")
}

func TestCarLoanProvidersApplyAndStatus(t *testing.T) {
	t.Run("Providers", func(t *testing.T) {
		server, call := carLoanServer(t, http.StatusOK, `{"response_code":"00","providers":[{"provider_id":"p1","provider_type":3,"credential_id":"cred-1","is_active":true,"allowed_installment_counts":[12,24,36]}]}`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token_123").CarLoanProviders(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "/loan/car/providers", call.path)
		require.Len(t, res.Providers, 1)
		assert.Equal(t, []int32{12, 24, 36}, res.Providers[0].AllowedInstallmentCounts)
		assert.True(t, res.Providers[0].IsActive)
	})

	t.Run("Apply", func(t *testing.T) {
		server, call := carLoanServer(t, http.StatusOK, `{"response_code":"00","application_id":"app-1","redirect_url":"https://bank","form_data":{"token":"abc"}}`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token_123").CarLoanApply(context.Background(), tapsilat.CarLoanApplyRequest{
			CredentialID:  "cred-1",
			OfferID:       "offer-1",
			Term:          36,
			PaymentPlanID: "plan-1",
			CallbackURL:   "https://checkout/cb",
			OrderID:       "order-1",
		})
		require.NoError(t, err)
		assert.Equal(t, "/loan/car/apply", call.path)
		assert.Equal(t, "https://checkout/cb", call.body["callback_url"])
		assert.Equal(t, "app-1", res.ApplicationID)
		assert.Equal(t, map[string]string{"token": "abc"}, res.FormData)
	})

	t.Run("Status", func(t *testing.T) {
		server, call := carLoanServer(t, http.StatusOK, `{"response_code":"00","status":"APPROVED","is_approved":true,"payment_plan":{"term":"36","loan_amount":350000}}`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token_123").CarLoanStatus(context.Background(), tapsilat.CarLoanStatusRequest{CredentialID: "cred-1", ApplicationID: "app-1"})
		require.NoError(t, err)
		assert.Equal(t, "/loan/car/status", call.path)
		assert.True(t, res.IsApproved)
		require.NotNil(t, res.PaymentPlan)
		assert.Equal(t, int64(36), res.PaymentPlan.Term)
	})
}

func TestCarLoanReturnsAPIError(t *testing.T) {
	server, _ := carLoanServer(t, http.StatusBadRequest, `{"error":"LOAN_PROVIDER_NOT_FOUND","code":411001}`)
	_, err := tapsilat.NewCustomAPI(server.URL, "token_123").CarLoanCheckLimit(context.Background(), tapsilat.CarLoanCheckLimitRequest{})
	require.Error(t, err)

	var apiErr *tapsilat.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
	assert.Equal(t, "LOAN_PROVIDER_NOT_FOUND", apiErr.Message)
}
