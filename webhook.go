package tapsilat

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"
)

// WebhookSignatureHeader is the header that carries the signature of every callback Tapsilat sends.
const WebhookSignatureHeader = "Tapsilat-Signature"

// DefaultWebhookTolerance is how far a callback's signature timestamp may be from the current time.
const DefaultWebhookTolerance = 5 * time.Minute

// Callback types. Each organization chooses one in the panel; WebhookTypeOrderLite is the default.
const (
	WebhookTypeOrder              = "order"
	WebhookTypeOrderLite          = "order_lite"
	WebhookTypeOrderPaymentDetail = "order_payment_detail"
	WebhookTypeExtendedOrder      = "extended_order"
	WebhookTypePaymentFacilitator = "payment_facilitator"
)

const zeroUUID = "00000000-0000-0000-0000-000000000000"

// VerifyWebhook verifies the Tapsilat-Signature header of a callback against the organization's
// webhook secret, rejecting a timestamp more than DefaultWebhookTolerance from the current time.
// payload must be the raw request body, before any JSON decoding.
func VerifyWebhook(payload []byte, signature, secret string) bool {
	return VerifyWebhookAt(payload, signature, secret, DefaultWebhookTolerance, time.Now())
}

// VerifyWebhookAt is VerifyWebhook with the tolerance (0 disables the timestamp check) and the
// time to compare the timestamp with.
//
// Tapsilat sends "t=<unix timestamp>,v1=<hex HMAC-SHA256 of '<timestamp>.<raw body>'>"; any of
// several v1 values may match, and the comparison takes constant time. A bare hex HMAC-SHA256 of
// the payload, with or without a "sha256=" prefix, is still accepted.
func VerifyWebhookAt(payload []byte, signature, secret string, tolerance time.Duration, now time.Time) bool {
	signature = strings.TrimSpace(signature)
	if signature == "" || secret == "" {
		return false
	}

	if strings.HasPrefix(signature, "sha256=") || !strings.Contains(signature, "=") {
		expected := webhookHMAC(secret, payload)
		return hmac.Equal([]byte(expected), []byte(strings.ToLower(strings.TrimPrefix(signature, "sha256="))))
	}

	var timestamp string
	var candidates []string
	for _, part := range strings.Split(signature, ",") {
		key, value, ok := strings.Cut(strings.TrimSpace(part), "=")
		if !ok {
			continue
		}
		switch key {
		case "t":
			timestamp = value
		case "v1":
			candidates = append(candidates, strings.ToLower(value))
		}
	}
	if !isDigits(timestamp) || len(candidates) == 0 {
		return false
	}
	unix, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return false
	}
	if tolerance > 0 {
		// Whole seconds on both sides, as the header carries them; Sub saturates on overflow.
		age := time.Unix(now.Unix(), 0).Sub(time.Unix(unix, 0))
		if age > tolerance || age < -tolerance {
			return false
		}
	}

	expected := []byte(webhookHMAC(secret, append([]byte(timestamp+"."), payload...)))
	for _, candidate := range candidates {
		if hmac.Equal(expected, []byte(candidate)) {
			return true
		}
	}
	return false
}

func webhookHMAC(secret string, message []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(message)
	return hex.EncodeToString(mac.Sum(nil))
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// WebhookEvent is a callback Tapsilat posts to an organization's callback URLs or to an order's
// payment_success_url and payment_failure_url: {"type": "...", "data": {...}, "token": "..."}.
//
// The shape of Data depends on Type. The methods read the order's reference ID, conversation ID
// and status from whichever shape arrived. Verify the signature with VerifyWebhook before
// trusting a callback.
type WebhookEvent struct {
	// Type is the callback type, such as WebhookTypeOrderLite; empty when the payload has none.
	Type string
	// Token is unique per delivery.
	Token string
	// Data is the callback's "data" object, or the whole payload when it has none.
	Data map[string]any
	// Payload is the decoded body.
	Payload map[string]any
}

// ParseWebhookEvent decodes a raw callback body. Numbers are decoded as json.Number.
func ParseWebhookEvent(payload []byte) (WebhookEvent, error) {
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.UseNumber()
	var decoded map[string]any
	if err := decoder.Decode(&decoded); err != nil || decoded == nil {
		return WebhookEvent{}, errors.New("the webhook payload is not a JSON object")
	}
	if _, err := decoder.Token(); err != io.EOF {
		return WebhookEvent{}, errors.New("the webhook payload is not a JSON object")
	}
	return NewWebhookEvent(decoded), nil
}

// NewWebhookEvent wraps an already decoded callback body.
func NewWebhookEvent(payload map[string]any) WebhookEvent {
	event := WebhookEvent{Payload: payload, Data: payload}
	event.Type, _ = payload["type"].(string)
	event.Token, _ = payload["token"].(string)
	if data, ok := payload["data"].(map[string]any); ok {
		event.Data = data
	}
	return event
}

// ReferenceID returns the order's reference ID, or "" when the callback carries none.
func (e WebhookEvent) ReferenceID() string {
	order, _ := e.Data["order"].(map[string]any)
	return firstID(e.Data["orderReferenceId"], e.Data["order_reference_id"], e.Data["reference_id"], order["reference_id"])
}

// ConversationID returns the conversation ID the order was created with, or "".
func (e WebhookEvent) ConversationID() string {
	return firstID(e.Data["conversationId"], e.Data["conversation_id"])
}

// StatusCode returns the order status code, such as OrderStatusPaid, or 0 when the callback type
// carries no status.
func (e WebhookEvent) StatusCode() int {
	if code := OrderStatusCode(e.Data["status"]); code != 0 {
		return code
	}
	order, _ := e.Data["order"].(map[string]any)
	return OrderStatusCode(order["status"])
}

// Status returns the order status name, such as "Paid", or "".
func (e WebhookEvent) Status() string {
	return OrderStatusName(e.StatusCode())
}

// IsPaid reports whether the order's payment was captured.
func (e WebhookEvent) IsPaid() bool {
	return IsPaidOrderStatus(e.StatusCode())
}

// IsFailed reports whether the order will not be paid.
func (e WebhookEvent) IsFailed() bool {
	return IsFailedOrderStatus(e.StatusCode())
}

// IsRefunded reports whether the order's payment was refunded in full or in part.
func (e WebhookEvent) IsRefunded() bool {
	return IsRefundedOrderStatus(e.StatusCode())
}

func firstID(values ...any) string {
	for _, value := range values {
		if s, ok := value.(string); ok {
			if s != "" && s != zeroUUID {
				return s
			}
			continue
		}
		if n, ok := integer(value); ok {
			return strconv.FormatInt(n, 10)
		}
	}
	return ""
}
