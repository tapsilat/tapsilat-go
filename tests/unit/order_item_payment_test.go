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

func TestAddOrderOIP(t *testing.T) {
	var body map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/order/oip", r.URL.Path)
		assert.Equal(t, "Bearer token_123", r.Header.Get("Authorization"))
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	_, err := tapsilat.NewCustomAPI(server.URL, "token_123").AddOrderOIP(context.Background(), tapsilat.OrderOIPRequest{
		OrderID:      "order-1",
		BasketItemID: "item-1",
		Amount:       150000,
		Type:         tapsilat.OrderItemPaymentTypeMoneyTransfer,
		PaidDate:     "2026-10-23T10:00:00Z",
	})
	require.NoError(t, err)
	assert.Equal(t, "order-1", body["order_id"])
	assert.Equal(t, "item-1", body["basket_item_id"])
	assert.Equal(t, 150000.0, body["amount"])
	assert.Equal(t, 2.0, body["type"])
	assert.Equal(t, "2026-10-23T10:00:00Z", body["paid_date"])
	assert.NotContains(t, body, "wallet_id")
}

func TestAddOrderOIPReturnsAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":"ORDER_ADD_OIP_PERMISSION_DENIED","code":403}`))
	}))
	defer server.Close()

	_, err := tapsilat.NewCustomAPI(server.URL, "token_123").AddOrderOIP(context.Background(), tapsilat.OrderOIPRequest{OrderID: "o", BasketItemID: "i", Amount: 1, Type: tapsilat.OrderItemPaymentTypeMoneyTransfer})
	var apiErr *tapsilat.APIError
	require.ErrorAs(t, err, &apiErr)
	assert.Equal(t, http.StatusForbidden, apiErr.StatusCode)
}

func TestGetOrderIDResolvesTheInternalIDFromTheReference(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/order/ref-1":
			_, _ = w.Write([]byte(`{"reference_id":"ref-1","conversation_id":"conv-1"}`))
		case "/order/list":
			assert.Equal(t, "conv-1", r.URL.Query().Get("conversation_id"))
			_, _ = w.Write([]byte(`{"page":1,"per_page":50,"total":2,"rows":[{"id":"other-id","reference_id":"ref-2"},{"id":"order-id-1","reference_id":"ref-1"}]}`))
		default:
			t.Errorf("unexpected path %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	id, err := tapsilat.NewCustomAPI(server.URL, "token_123").GetOrderID(context.Background(), "ref-1")
	require.NoError(t, err)
	assert.Equal(t, "order-id-1", id)
}

func TestGetOrderIDFailsWhenTheOrderIsNotListed(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/order/ref-1":
			_, _ = w.Write([]byte(`{"reference_id":"ref-1","conversation_id":"conv-1"}`))
		case "/order/list":
			_, _ = w.Write([]byte(`{"rows":[{"id":"other-id","reference_id":"ref-2"}]}`))
		}
	}))
	defer server.Close()

	_, err := tapsilat.NewCustomAPI(server.URL, "token_123").GetOrderID(context.Background(), "ref-1")
	require.ErrorIs(t, err, tapsilat.ErrOrderIDNotFound)
}

func TestGetOrderIDFailsWithoutAConversation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"reference_id":"ref-1"}`))
	}))
	defer server.Close()

	_, err := tapsilat.NewCustomAPI(server.URL, "token_123").GetOrderID(context.Background(), "ref-1")
	require.ErrorIs(t, err, tapsilat.ErrOrderIDNotFound)
}
