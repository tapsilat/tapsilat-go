package unit_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

type elementsCall struct {
	method  string
	path    string
	rawPath string
	auth    string
	rawBody string
	body    map[string]any
}

func elementsServer(t *testing.T, status int, response string) (*httptest.Server, *elementsCall) {
	t.Helper()
	call := &elementsCall{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		call.method = r.Method
		call.path = r.URL.Path
		call.rawPath = r.URL.EscapedPath()
		call.auth = r.Header.Get("Authorization")
		raw, _ := io.ReadAll(r.Body)
		call.rawBody = string(raw)
		_ = json.Unmarshal(raw, &call.body)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", "no-store")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(response))
	}))
	t.Cleanup(server.Close)
	return server, call
}

func TestCreateElementsSession(t *testing.T) {
	server, call := elementsServer(t, http.StatusOK, `{
		"id": "es_test_8Kq2nV0rX4",
		"purpose": "pay",
		"status": "requires_payment_method",
		"order_reference_id": "order-ref-1",
		"attempts": 0,
		"max_attempts": 3,
		"created_at": "2026-10-22T09:00:00Z",
		"expires_at": "2026-10-22T09:30:00Z",
		"client_secret": "es_test_8Kq2nV0rX4_secret_q8Zr0bM1",
		"elements_base_url": "https://checkout.tapsilat.dev"
	}`)
	api := tapsilat.NewCustomAPI(server.URL, "token_123")

	session, err := api.CreateElementsSession(context.Background(), tapsilat.ElementsSessionCreateRequest{
		Purpose:          tapsilat.ElementsPurposePay,
		OrderReferenceID: "order-ref-1",
		SaveCard:         true,
	})
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, call.method)
	assert.Equal(t, "/elements/sessions", call.path)
	assert.Equal(t, "Bearer token_123", call.auth)
	assert.Equal(t, map[string]any{"purpose": "pay", "order_reference_id": "order-ref-1", "save_card": true}, call.body)

	assert.Equal(t, "es_test_8Kq2nV0rX4", session.ID)
	assert.Equal(t, tapsilat.ElementsPurposePay, session.Purpose)
	assert.Equal(t, tapsilat.ElementsStatusRequiresPaymentMethod, session.Status)
	assert.Equal(t, "order-ref-1", session.OrderReferenceID)
	assert.Equal(t, int32(3), session.MaxAttempts)
	assert.Equal(t, "2026-10-22T09:30:00Z", session.ExpiresAt)
	assert.Empty(t, session.CompletedAt)
	assert.Nil(t, session.Card)
	assert.Nil(t, session.LastError)
	assert.Equal(t, "es_test_8Kq2nV0rX4_secret_q8Zr0bM1", session.ClientSecret)
	assert.Equal(t, "https://checkout.tapsilat.dev", session.ElementsBaseURL)
}

func TestCreateElementsSessionTokenize(t *testing.T) {
	server, call := elementsServer(t, http.StatusOK, `{"id":"es_test_1","purpose":"tokenize","status":"requires_payment_method","user_id":"user_123","client_secret":"es_test_1_secret_x"}`)
	api := tapsilat.NewCustomAPI(server.URL, "token_123")

	session, err := api.CreateElementsSession(context.Background(), tapsilat.ElementsSessionCreateRequest{
		Purpose:   tapsilat.ElementsPurposeTokenize,
		UserID:    "user_123",
		UserEmail: "john@doe.com",
		Mode:      tapsilat.ElementsModeNonSecure,
	})
	require.NoError(t, err)

	assert.Equal(t, map[string]any{"purpose": "tokenize", "user_id": "user_123", "user_email": "john@doe.com", "mode": 3.0}, call.body)
	assert.Equal(t, "user_123", session.UserID)
	assert.Equal(t, "es_test_1_secret_x", session.ClientSecret)
}

func TestElementsSessionWithSecretRedactsSecretWhenPrinted(t *testing.T) {
	session := tapsilat.ElementsSessionWithSecret{
		ElementsSession: tapsilat.ElementsSession{ID: "es_test_1", Status: tapsilat.ElementsStatusRequiresPaymentMethod},
		ClientSecret:    "es_test_1_secret_x",
	}

	for _, format := range []string{"%v", "%+v", "%s", "%#v"} {
		printed := fmt.Sprintf(format, session)
		assert.NotContains(t, printed, "es_test_1_secret_x", format)
		assert.Contains(t, printed, "es_test_1", format)
	}

	encoded, err := json.Marshal(session)
	require.NoError(t, err)
	assert.Contains(t, string(encoded), `"client_secret":"es_test_1_secret_x"`)
}

func TestGetElementsSession(t *testing.T) {
	server, call := elementsServer(t, http.StatusOK, `{
		"id": "es_test_8Kq2nV0rX4",
		"purpose": "pay",
		"status": "failed",
		"order_reference_id": "order-ref-1",
		"attempts": 3,
		"max_attempts": 3,
		"completed_at": "2026-10-22T09:12:00Z",
		"card": {"id": "card-1", "brand": "visa", "last4": "4242", "exp_month": 12, "exp_year": 2030},
		"last_error": {"code": "card_declined", "message": "The card was declined.", "retryable": false}
	}`)
	api := tapsilat.NewCustomAPI(server.URL, "token_123")

	session, err := api.GetElementsSession(context.Background(), "es_test_8Kq2nV0rX4")
	require.NoError(t, err)

	assert.Equal(t, http.MethodGet, call.method)
	assert.Equal(t, "/elements/sessions/es_test_8Kq2nV0rX4", call.path)
	assert.Equal(t, "Bearer token_123", call.auth)
	assert.Empty(t, call.rawBody)

	assert.Equal(t, tapsilat.ElementsStatusFailed, session.Status)
	assert.Equal(t, int32(3), session.Attempts)
	assert.Equal(t, "2026-10-22T09:12:00Z", session.CompletedAt)
	require.NotNil(t, session.Card)
	assert.Equal(t, tapsilat.ElementsCardSummary{ID: "card-1", Brand: "visa", Last4: "4242", ExpMonth: 12, ExpYear: 2030}, *session.Card)
	require.NotNil(t, session.LastError)
	assert.Equal(t, tapsilat.ElementsError{Code: "card_declined", Message: "The card was declined."}, *session.LastError)
}

func TestGetElementsSessionEscapesID(t *testing.T) {
	server, call := elementsServer(t, http.StatusOK, `{"id":"es_test_1"}`)
	api := tapsilat.NewCustomAPI(server.URL, "token_123")

	_, err := api.GetElementsSession(context.Background(), "es_test_1/../x")
	require.NoError(t, err)
	assert.Equal(t, "/elements/sessions/es_test_1%2F..%2Fx", call.rawPath)
}

func TestGetElementsSessionReadsNullCardAndError(t *testing.T) {
	server, _ := elementsServer(t, http.StatusOK, `{"id":"es_mock_1","purpose":"pay","status":"requires_payment_method","card":null,"last_error":null}`)
	api := tapsilat.NewCustomAPI(server.URL, "token_123")

	session, err := api.GetElementsSession(context.Background(), "es_mock_1")
	require.NoError(t, err)
	assert.Nil(t, session.Card)
	assert.Nil(t, session.LastError)
}

func TestCancelElementsSession(t *testing.T) {
	server, call := elementsServer(t, http.StatusOK, `{"id":"es_test_8Kq2nV0rX4","purpose":"pay","status":"canceled","attempts":0,"max_attempts":3}`)
	api := tapsilat.NewCustomAPI(server.URL, "token_123")

	session, err := api.CancelElementsSession(context.Background(), "es_test_8Kq2nV0rX4")
	require.NoError(t, err)

	assert.Equal(t, http.MethodPost, call.method)
	assert.Equal(t, "/elements/sessions/es_test_8Kq2nV0rX4/cancel", call.path)
	assert.Equal(t, "Bearer token_123", call.auth)
	assert.Empty(t, call.rawBody)
	assert.Equal(t, tapsilat.ElementsStatusCanceled, session.Status)
}

func TestElementsSessionErrors(t *testing.T) {
	tests := []struct {
		name      string
		status    int
		body      string
		call      func(api *tapsilat.API) error
		code      string
		message   string
		retryable bool
	}{
		{
			name:   "CreateNotEnabled",
			status: http.StatusForbidden,
			body:   `{"error":{"code":"elements_not_enabled","message":"Elements is not enabled for this organization"}}`,
			call: func(api *tapsilat.API) error {
				_, err := api.CreateElementsSession(context.Background(), tapsilat.ElementsSessionCreateRequest{Purpose: tapsilat.ElementsPurposePay})
				return err
			},
			code:    "elements_not_enabled",
			message: "Elements is not enabled for this organization",
		},
		{
			name:   "GetNotFound",
			status: http.StatusNotFound,
			body:   `{"error":{"code":"session_not_found","message":"session not found"}}`,
			call: func(api *tapsilat.API) error {
				_, err := api.GetElementsSession(context.Background(), "es_test_missing")
				return err
			},
			code:    "session_not_found",
			message: "session not found",
		},
		{
			name:   "CancelUnavailable",
			status: http.StatusServiceUnavailable,
			body:   `{"error":{"code":"unavailable","message":"try again","retryable":true}}`,
			call: func(api *tapsilat.API) error {
				_, err := api.CancelElementsSession(context.Background(), "es_test_1")
				return err
			},
			code:      "unavailable",
			message:   "try again",
			retryable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, _ := elementsServer(t, tt.status, tt.body)
			err := tt.call(tapsilat.NewCustomAPI(server.URL, "token_123"))
			require.Error(t, err)

			var apiErr *tapsilat.APIError
			require.ErrorAs(t, err, &apiErr)
			assert.Equal(t, tt.status, apiErr.StatusCode)
			assert.Equal(t, tt.code, apiErr.Code)
			assert.Equal(t, tt.message, apiErr.Message)
			assert.Equal(t, tt.retryable, apiErr.Retryable)
			assert.Contains(t, apiErr.Error(), "["+tt.code+"] "+tt.message)
		})
	}
}
