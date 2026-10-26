package unit_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

// An order_lite callback and a known-good Tapsilat-Signature for this body and secret; the expected
// signature was computed with openssl, independently of the SDK.
const (
	webhookPayload = `{"data":{"id":"5f0c6f1e-1111-4a5b-9c1d-2e3f4a5b6c7d",` +
		`"orderReferenceId":"0b6a4f9e-2222-4c3d-8e9f-a1b2c3d4e5f6","conversationId":"wc-1001","status":"3"},` +
		`"type":"order_lite","token":"9d8c7b6a-3333-4e5f-a6b7-c8d9e0f1a2b3"}`
	webhookSecret    = "whsec_test"
	webhookTimestamp = 1726480000
	webhookSig       = "3e003945e8122172a03ae1aa56fd099b56e6ac11a8472c03e33809c7c480aeb6"
	webhookBareSig   = "402b9d6616b32be533a4fc0300e8db19c304a56137d3fb274c614f30f9be6fbb"
)

func webhookHeader(signature, timestamp string) string {
	return "t=" + timestamp + ",v1=" + signature
}

func verifyAt(payload, header, secret string, tolerance time.Duration, now int64) bool {
	return tapsilat.VerifyWebhookAt([]byte(payload), header, secret, tolerance, time.Unix(now, 0))
}

func TestVerifyWebhook(t *testing.T) {
	header := webhookHeader(webhookSig, "1726480000")

	t.Run("VerifiesTheTapsilatSignatureHeader", func(t *testing.T) {
		assert.Equal(t, "Tapsilat-Signature", tapsilat.WebhookSignatureHeader)
		assert.True(t, verifyAt(webhookPayload, header, webhookSecret, 300*time.Second, webhookTimestamp+10))
		assert.True(t, verifyAt(webhookPayload, " t=1726480000, v1="+strings.ToUpper(webhookSig)+" ", webhookSecret,
			300*time.Second, webhookTimestamp))
	})

	t.Run("RejectsATamperedBodyOrTheWrongSecret", func(t *testing.T) {
		tampered := strings.Replace(webhookPayload, `"status":"3"`, `"status":"13"`, 1)

		assert.False(t, verifyAt(tampered, header, webhookSecret, 300*time.Second, webhookTimestamp))
		assert.False(t, verifyAt(webhookPayload, header, "other", 300*time.Second, webhookTimestamp))
		assert.False(t, verifyAt(webhookPayload, header, "", 300*time.Second, webhookTimestamp))
		assert.False(t, verifyAt(webhookPayload, "", webhookSecret, 300*time.Second, webhookTimestamp))
	})

	t.Run("RejectsAReplayedOrMalformedHeader", func(t *testing.T) {
		assert.False(t, verifyAt(webhookPayload, header, webhookSecret, 300*time.Second, webhookTimestamp+301))
		assert.False(t, verifyAt(webhookPayload, header, webhookSecret, 300*time.Second, webhookTimestamp-301))
		assert.True(t, verifyAt(webhookPayload, header, webhookSecret, 300*time.Second, webhookTimestamp+300))
		assert.True(t, verifyAt(webhookPayload, header, webhookSecret, 0, webhookTimestamp+86400))
		assert.False(t, verifyAt(webhookPayload, "v1="+webhookSig, webhookSecret, 0, webhookTimestamp))
		assert.False(t, verifyAt(webhookPayload, "t=1726480000", webhookSecret, 0, webhookTimestamp))
		assert.False(t, verifyAt(webhookPayload, webhookHeader(webhookSig, "1726480000.5"), webhookSecret, 0, webhookTimestamp))
		assert.False(t, verifyAt(webhookPayload, webhookHeader(webhookSig, "-1726480000"), webhookSecret, 0, webhookTimestamp))
		assert.False(t, verifyAt(webhookPayload, webhookHeader(webhookSig, "99999999999999999999"), webhookSecret, 0, webhookTimestamp))
		assert.False(t, verifyAt(webhookPayload, webhookHeader(webhookSig, "9223372036854775807"), webhookSecret,
			300*time.Second, webhookTimestamp))
	})

	t.Run("AcceptsAnyOfSeveralV1Signatures", func(t *testing.T) {
		multi := "t=1726480000,v1=" + strings.Repeat("0", 64) + ",v1=" + webhookSig

		assert.True(t, verifyAt(webhookPayload, multi, webhookSecret, 300*time.Second, webhookTimestamp))
	})

	t.Run("UsesTheCurrentTimeByDefault", func(t *testing.T) {
		now := time.Now().Unix()
		mac := hmac.New(sha256.New, []byte(webhookSecret))
		mac.Write([]byte(strconv.FormatInt(now, 10) + "." + webhookPayload))
		signature := hex.EncodeToString(mac.Sum(nil))

		assert.True(t, tapsilat.VerifyWebhook([]byte(webhookPayload), webhookHeader(signature, strconv.FormatInt(now, 10)), webhookSecret))
		assert.False(t, tapsilat.VerifyWebhook([]byte(webhookPayload), header, webhookSecret))
	})

	t.Run("AcceptsTheLegacyBareForm", func(t *testing.T) {
		assert.True(t, tapsilat.VerifyWebhook([]byte(webhookPayload), webhookBareSig, webhookSecret))
		assert.True(t, tapsilat.VerifyWebhook([]byte(webhookPayload), "sha256="+webhookBareSig, webhookSecret))
		assert.True(t, tapsilat.VerifyWebhook([]byte(webhookPayload), strings.ToUpper(webhookBareSig), webhookSecret))
		assert.False(t, tapsilat.VerifyWebhook([]byte(webhookPayload), "sha256=invalid", webhookSecret))
		assert.False(t, tapsilat.VerifyWebhook([]byte("tampered"), webhookBareSig, webhookSecret))
	})
}

func TestWebhookEvent(t *testing.T) {
	t.Run("ParsesAnOrderLiteCallback", func(t *testing.T) {
		event, err := tapsilat.ParseWebhookEvent([]byte(webhookPayload))
		require.NoError(t, err)

		assert.Equal(t, tapsilat.WebhookTypeOrderLite, event.Type)
		assert.Equal(t, "9d8c7b6a-3333-4e5f-a6b7-c8d9e0f1a2b3", event.Token)
		assert.Equal(t, "0b6a4f9e-2222-4c3d-8e9f-a1b2c3d4e5f6", event.ReferenceID())
		assert.Equal(t, "wc-1001", event.ConversationID())
		assert.Equal(t, tapsilat.OrderStatusPaid, event.StatusCode())
		assert.Equal(t, "Paid", event.Status())
		assert.True(t, event.IsPaid())
		assert.False(t, event.IsFailed())
	})

	t.Run("ParsesAnOrderCallback", func(t *testing.T) {
		event, err := tapsilat.ParseWebhookEvent([]byte(`{"type":"order","token":"t-1",` +
			`"data":{"id":"o-1","reference_id":"ref-1","conversation_id":"c-1","status":13}}`))
		require.NoError(t, err)

		assert.Equal(t, "ref-1", event.ReferenceID())
		assert.Equal(t, "c-1", event.ConversationID())
		assert.Equal(t, "Failure", event.Status())
		assert.True(t, event.IsFailed())
	})

	t.Run("ParsesAnOrderPaymentDetailCallback", func(t *testing.T) {
		event, err := tapsilat.ParseWebhookEvent([]byte(`{"type":"order_payment_detail","token":"t-2",` +
			`"data":{"id":"d-1","order_reference_id":"ref-2","conversation_id":"c-2",` +
			`"order":{"reference_id":"ref-2","id":"o-2","status":"Paid"},"order_payment_status":{"code":0}}}`))
		require.NoError(t, err)

		assert.Equal(t, "ref-2", event.ReferenceID())
		assert.Equal(t, "c-2", event.ConversationID())
		assert.Equal(t, tapsilat.OrderStatusPaid, event.StatusCode())
	})

	t.Run("ParsesCallbacksWithoutAStatus", func(t *testing.T) {
		extended, err := tapsilat.ParseWebhookEvent([]byte(`{"type":"extended_order",` +
			`"data":{"id":"o-3","reference_id":"ref-3","terms":[{"reference_id":"term-1","status":2}]}}`))
		require.NoError(t, err)
		assert.Equal(t, "ref-3", extended.ReferenceID())
		assert.Equal(t, "", extended.ConversationID())
		assert.Equal(t, 0, extended.StatusCode())
		assert.Equal(t, "", extended.Status())
		assert.False(t, extended.IsPaid())

		facilitator, err := tapsilat.ParseWebhookEvent([]byte(`{"type":"payment_facilitator",` +
			`"data":{"id":"p-1","payment_id":"pay-1","reference_id":"ref-4"}}`))
		require.NoError(t, err)
		assert.Equal(t, "ref-4", facilitator.ReferenceID())
	})

	t.Run("IgnoresAZeroUUID", func(t *testing.T) {
		event := tapsilat.NewWebhookEvent(map[string]any{
			"type": "order_payment_detail",
			"data": map[string]any{
				"order_reference_id": "00000000-0000-0000-0000-000000000000",
				"order":              map[string]any{"reference_id": "ref-5"},
			},
		})

		assert.Equal(t, "ref-5", event.ReferenceID())
	})

	t.Run("UsesThePayloadWhenThereIsNoData", func(t *testing.T) {
		event := tapsilat.NewWebhookEvent(map[string]any{"reference_id": "ref-6", "status": float64(9)})

		assert.Equal(t, "", event.Type)
		assert.Equal(t, "ref-6", event.ReferenceID())
		assert.Equal(t, "Completed", event.Status())
	})

	t.Run("RejectsAPayloadThatIsNotAJSONObject", func(t *testing.T) {
		for _, payload := range []string{`"order"`, `null`, `[1]`, `{"type":"order"} {}`, `not json`} {
			_, err := tapsilat.ParseWebhookEvent([]byte(payload))
			assert.Error(t, err, payload)
		}
	})
}

func TestOrderStatusHelpers(t *testing.T) {
	assert.Equal(t, tapsilat.OrderStatusPayment, tapsilat.OrderStatusCode("waiting for payment"))
	assert.Equal(t, tapsilat.OrderStatusExpired, tapsilat.OrderStatusCode("21"))
	assert.Equal(t, tapsilat.OrderStatusFailure, tapsilat.OrderStatusCode(13))
	assert.Equal(t, tapsilat.OrderStatusPartiallyPaid, tapsilat.OrderStatusCode("Partially Paid"))
	assert.Equal(t, "Partially refunded", tapsilat.OrderStatusName(15))
	assert.Equal(t, 0, tapsilat.OrderStatusCode("SUCCESS"))
	assert.Equal(t, 0, tapsilat.OrderStatusCode(99))
	assert.Equal(t, 0, tapsilat.OrderStatusCode(nil))
	assert.Equal(t, 0, tapsilat.OrderStatusCode(3.5))
	assert.Equal(t, "", tapsilat.OrderStatusName(0))

	assert.True(t, tapsilat.IsPaidOrderStatus("Sub merchant payment approved"))
	assert.True(t, tapsilat.IsPaidOrderStatus(tapsilat.OrderStatusCompleted))
	assert.False(t, tapsilat.IsPaidOrderStatus("Waiting for payment"))
	assert.True(t, tapsilat.IsFailedOrderStatus("Expired"))
	assert.False(t, tapsilat.IsFailedOrderStatus("Paid"))
	assert.True(t, tapsilat.IsRefundedOrderStatus("10"))
	assert.False(t, tapsilat.IsRefundedOrderStatus(nil))

	assert.Equal(t, tapsilat.OrderStatusPaid, tapsilat.GetOrderStatusByStr("paid"))
}
