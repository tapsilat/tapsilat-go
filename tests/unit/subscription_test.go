package unit_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

func TestSubscriptionCreateRequestJSON(t *testing.T) {
	t.Run("SerializesIntervalCadenceAndOmitsEmptyCardID", func(t *testing.T) {
		payload := tapsilat.SubscriptionCreateRequest{
			Amount:              100,
			Currency:            "TRY",
			Cycle:               12,
			ExternalReferenceID: "ext_sub_123",
			Interval:            "month",
			IntervalCount:       1,
			Title:               "Monthly Subscription",
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"amount": 100,
			"billing": {},
			"currency": "TRY",
			"cycle": 12,
			"external_reference_id": "ext_sub_123",
			"interval": "month",
			"interval_count": 1,
			"title": "Monthly Subscription",
			"user": {}
		}`, string(body))
		assert.NotContains(t, string(body), `"card_id"`)
	})

	t.Run("PreservesLegacyPeriodContract", func(t *testing.T) {
		payload := tapsilat.SubscriptionCreateRequest{
			Amount:      100,
			CardID:      "card_token_123",
			Currency:    "TRY",
			Cycle:       1,
			PaymentDate: 1,
			Period:      30,
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"amount": 100,
			"billing": {},
			"card_id": "card_token_123",
			"currency": "TRY",
			"cycle": 1,
			"payment_date": 1,
			"period": 30,
			"user": {}
		}`, string(body))
	})
}

func TestSubscriptionResponseJSON(t *testing.T) {
	t.Run("DecodesDetailCadenceAndLifecycleProgress", func(t *testing.T) {
		var detail tapsilat.SubscriptionDetail
		err := json.Unmarshal([]byte(`{
			"amount": "100.00",
			"currency": "TRY",
			"cycle": 12,
			"cycles_completed": 3,
			"due_date": "2026-10-02",
			"external_reference_id": "ext_sub_123",
			"interval": "month",
			"interval_count": 1,
			"is_active": true,
			"payment_date": 2,
			"payment_status": "paid",
			"period": 30,
			"title": "Monthly Subscription"
		}`), &detail)

		require.NoError(t, err)
		assert.Equal(t, "month", detail.Interval)
		assert.Equal(t, 1, detail.IntervalCount)
		assert.Equal(t, 12, detail.Cycle)
		assert.Equal(t, 3, detail.CyclesCompleted)
		assert.Equal(t, 30, detail.Period)
		assert.Equal(t, 2, detail.PaymentDate)
		assert.Equal(t, "2026-10-02", detail.DueDate)
	})

	t.Run("DecodesListItemCadenceLifecycleAndNextPaymentDate", func(t *testing.T) {
		var item tapsilat.SubscriptionListItem
		err := json.Unmarshal([]byte(`{
			"amount": "100.00",
			"currency": "TRY",
			"cycle": 0,
			"cycles_completed": 3,
			"external_reference_id": "ext_sub_123",
			"interval": "month",
			"interval_count": 1,
			"is_active": true,
			"next_payment_date": "2026-10-02",
			"payment_date": 2,
			"payment_status": "paid",
			"period": 30,
			"reference_id": "sub_123",
			"title": "Monthly Subscription"
		}`), &item)

		require.NoError(t, err)
		assert.Equal(t, "month", item.Interval)
		assert.Equal(t, 1, item.IntervalCount)
		assert.Equal(t, 0, item.Cycle)
		assert.Equal(t, 3, item.CyclesCompleted)
		assert.Equal(t, "2026-10-02", item.NextPaymentDate)
		assert.Equal(t, 30, item.Period)
		assert.Equal(t, 2, item.PaymentDate)
		assert.Equal(t, "sub_123", item.ReferenceID)
	})
}
