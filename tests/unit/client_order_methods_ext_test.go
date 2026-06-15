package unit_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

func mockServer(t *testing.T, expectedMethod, expectedPath, responseBody string, validateBody func(body string)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, expectedMethod, r.Method)
		assert.Equal(t, expectedPath, r.URL.Path)
		if validateBody != nil {
			bodyBytes, _ := io.ReadAll(r.Body)
			validateBody(string(bodyBytes))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseBody))
	}))
}

func TestNewOrderMethods(t *testing.T) {
	t.Run("GetOrderPaymentDetails", func(t *testing.T) {
		server := mockServer(t, http.MethodPost, "/order/payment-details", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.GetOrderPaymentDetails(context.Background(), tapsilat.OrderPaymentDetailDTO{ReferenceID: "ref1"})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("OrderManualCallback", func(t *testing.T) {
		server := mockServer(t, http.MethodPost, "/order/callback", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.OrderManualCallback(context.Background(), tapsilat.OrderManualCallbackDTO{ReferenceID: "ref1"})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("DeleteOrderTerm", func(t *testing.T) {
		server := mockServer(t, http.MethodDelete, "/order/term", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.DeleteOrderTerm(context.Background(), tapsilat.OrderPaymentTermDeleteDTO{OrderID: "ref1", TermReferenceID: "ref2"})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("UpdateOrderTerm", func(t *testing.T) {
		server := mockServer(t, http.MethodPatch, "/order/term", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.UpdateOrderTerm(context.Background(), tapsilat.OrderPaymentTermUpdateDTO{TermReferenceID: "ref1"})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("GetOrderTerm", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/order/term", r.URL.Path)
			assert.Equal(t, "ref1", r.URL.Query().Get("term_reference_id"))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"success"}`))
		}))
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.GetOrderTerm(context.Background(), "ref1")
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("OrderRelatedUpdate", func(t *testing.T) {
		server := mockServer(t, http.MethodPatch, "/order/releated", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.OrderRelatedUpdate(context.Background(), tapsilat.OrderRelatedReferenceDTO{ReferenceID: "ref1", RelatedReferenceID: "ref2"})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("OrderAccounting", func(t *testing.T) {
		server := mockServer(t, http.MethodPost, "/order/accounting", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.OrderAccounting(context.Background(), tapsilat.OrderAccountingRequest{OrderReferenceID: "ref1"})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("OrderPostAuth", func(t *testing.T) {
		server := mockServer(t, http.MethodPost, "/order/postauth", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.OrderPostAuth(context.Background(), tapsilat.OrderPostAuthRequest{ReferenceID: "ref1", Amount: 100})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("GetOrderPaymentDetailsByID", func(t *testing.T) {
		server := mockServer(t, http.MethodGet, "/order/ref1/payment-details", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.GetOrderPaymentDetailsByID(context.Background(), "ref1")
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("OrderPaymentOptionsUpdate", func(t *testing.T) {
		server := mockServer(t, http.MethodPatch, "/order/payment-options", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.OrderPaymentOptionsUpdate(context.Background(), tapsilat.OrderPaymentOptionsUpdateDTO{ReferenceID: "ref1"})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("SplitOrderItemPayment", func(t *testing.T) {
		server := mockServer(t, http.MethodPost, "/order/split", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.SplitOrderItemPayment(context.Background(), tapsilat.SplitOrderItemPaymentDTO{OrderID: "ref1"})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("OrderCallback", func(t *testing.T) {
		server := mockServer(t, http.MethodGet, "/orders/ref1/callback", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.OrderCallback(context.Background(), "ref1")
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("OrderVposQuery", func(t *testing.T) {
		server := mockServer(t, http.MethodGet, "/orders/ref1/vpos-query", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.OrderVposQuery(context.Background(), "ref1")
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("AddBasketItem", func(t *testing.T) {
		server := mockServer(t, http.MethodPost, "/order/basket-item", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.AddBasketItem(context.Background(), tapsilat.AddBasketItemRequest{})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("RemoveBasketItem", func(t *testing.T) {
		server := mockServer(t, http.MethodDelete, "/order/basket-item", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.RemoveBasketItem(context.Background(), tapsilat.RemoveBasketItemRequest{})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("UpdateBasketItem", func(t *testing.T) {
		server := mockServer(t, http.MethodPatch, "/order/basket-item", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.UpdateBasketItem(context.Background(), tapsilat.UpdateBasketItemRequest{})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})
}
