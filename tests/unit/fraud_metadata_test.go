package unit_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

func TestFraudMetadata(t *testing.T) {
	metadata := tapsilat.FraudMetadata(map[string]string{"check_in_date": "2026-12-20", "booking_reference": "BK-1"})
	require.Equal(t, []tapsilat.OrderMetadata{
		{Key: "fraud.booking_reference", Value: "BK-1"},
		{Key: "fraud.check_in_date", Value: "2026-12-20"},
	}, metadata, "prefixed and in key order")
	require.Empty(t, tapsilat.FraudMetadata(nil))

	body, err := json.Marshal(tapsilat.Order{Amount: 2500, Currency: "TRY", Metadata: append([]tapsilat.OrderMetadata{{Key: "invoice_no", Value: "INV-1"}}, metadata...)})
	require.NoError(t, err)
	var wire struct {
		Metadata []tapsilat.OrderMetadata `json:"metadata"`
	}
	require.NoError(t, json.Unmarshal(body, &wire))
	require.Len(t, wire.Metadata, 3)
}

func TestOrderStatusFraudReviewJSON(t *testing.T) {
	var status tapsilat.OrderStatus
	require.NoError(t, json.Unmarshal([]byte(`{"status":"On hold","fraud_review":"pending"}`), &status))
	require.Equal(t, tapsilat.FraudReviewPending, status.FraudReview)

	var paid tapsilat.OrderStatus
	require.NoError(t, json.Unmarshal([]byte(`{"status":"Paid"}`), &paid))
	require.Empty(t, paid.FraudReview)
}

func TestFraudItemData(t *testing.T) {
	data, err := tapsilat.FraudItemData(map[string]string{"booking_reference": "BK-1", "check_in_date": "2026-12-20"})
	require.NoError(t, err)
	require.JSONEq(t, `{"fraud":{"booking_reference":"BK-1","check_in_date":"2026-12-20"}}`, data)

	body, err := json.Marshal(tapsilat.OrderBasketItem{Id: "hotel-1", Data: data})
	require.NoError(t, err)
	require.Contains(t, string(body), `"data":"{\"fraud\":{`)
}
