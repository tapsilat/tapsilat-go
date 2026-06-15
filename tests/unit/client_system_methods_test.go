package unit_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

func TestSystemMethods(t *testing.T) {
	endpoints := []struct {
		name     string
		path     string
		callFunc func(api *tapsilat.API) (map[string]any, error)
	}{
		{"GetSystemOrderStatuses", "/system/order-statuses", func(api *tapsilat.API) (map[string]any, error) { return api.GetSystemOrderStatuses(context.Background()) }},
		{"GetSystemBasketItemTypes", "/system/basket-item-types", func(api *tapsilat.API) (map[string]any, error) { return api.GetSystemBasketItemTypes(context.Background()) }},
		{"GetSystemErrorCodes", "/system/error-codes", func(api *tapsilat.API) (map[string]any, error) { return api.GetSystemErrorCodes(context.Background()) }},
		{"GetSystemPaymentTermStatuses", "/system/payment-term-statuses", func(api *tapsilat.API) (map[string]any, error) { return api.GetSystemPaymentTermStatuses(context.Background()) }},
		{"GetSystemProductTypes", "/system/product-types", func(api *tapsilat.API) (map[string]any, error) { return api.GetSystemProductTypes(context.Background()) }},
		{"GetSystemShortcutTypes", "/system/shortcut-types", func(api *tapsilat.API) (map[string]any, error) { return api.GetSystemShortcutTypes(context.Background()) }},
		{"GetSystemTransactionPaymentTypes", "/system/transaction-payment-types", func(api *tapsilat.API) (map[string]any, error) { return api.GetSystemTransactionPaymentTypes(context.Background()) }},
		{"GetSystemTransactionPurposes", "/system/transaction-purposes", func(api *tapsilat.API) (map[string]any, error) { return api.GetSystemTransactionPurposes(context.Background()) }},
		{"GetSystemTransactionStatuses", "/system/transaction-statuses", func(api *tapsilat.API) (map[string]any, error) { return api.GetSystemTransactionStatuses(context.Background()) }},
	}

	for _, ep := range endpoints {
		t.Run(ep.name, func(t *testing.T) {
			server := mockServer(t, http.MethodGet, ep.path, `{"status":"success"}`, nil)
			defer server.Close()
			api := tapsilat.NewCustomAPI(server.URL, "token")
			res, err := ep.callFunc(api)
			require.NoError(t, err)
			assert.Equal(t, "success", res["status"])
		})
	}
}
