package tapsilat

import (
	"context"
	"fmt"
	"net/url"
)

// Purposes of an Elements session.
const (
	ElementsPurposeTokenize = "tokenize"
	ElementsPurposePay      = "pay"
)

// Tokenize modes of an Elements session.
const (
	ElementsModeSecure     int32 = 1 // 3D Secure
	ElementsModeHalfSecure int32 = 2
	ElementsModeNonSecure  int32 = 3
)

// Statuses of an Elements session.
const (
	ElementsStatusRequiresPaymentMethod = "requires_payment_method"
	ElementsStatusProcessing            = "processing"
	ElementsStatusRequiresAction        = "requires_action"
	ElementsStatusSucceeded             = "succeeded"
	ElementsStatusFailed                = "failed"
	ElementsStatusExpired               = "expired"
	ElementsStatusCanceled              = "canceled"
)

type ElementsSessionCreateRequest struct {
	Purpose string `json:"purpose"`
	// OrderReferenceID is required for ElementsPurposePay.
	OrderReferenceID string `json:"order_reference_id,omitempty"`
	// UserID is the card owner for ElementsPurposeTokenize, scoped to the caller's organization.
	UserID    string `json:"user_id,omitempty"`
	UserEmail string `json:"user_email,omitempty"`
	// Mode is the tokenize mode: ElementsModeSecure, ElementsModeHalfSecure or ElementsModeNonSecure.
	Mode int32 `json:"mode,omitempty"`
	// SaveCard lets the buyer save the card on a pay session.
	SaveCard bool `json:"save_card,omitempty"`
}

type ElementsCardSummary struct {
	ID       string `json:"id,omitempty"`
	Brand    string `json:"brand"`
	Last4    string `json:"last4"`
	ExpMonth int32  `json:"exp_month"`
	ExpYear  int32  `json:"exp_year"`
}

type ElementsError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Retryable bool   `json:"retryable,omitempty"`
}

// ElementsSession is the authoritative state of an Elements session. Timestamps are RFC 3339.
type ElementsSession struct {
	ID               string               `json:"id"`
	Purpose          string               `json:"purpose"`
	Status           string               `json:"status"`
	OrderReferenceID string               `json:"order_reference_id,omitempty"`
	UserID           string               `json:"user_id,omitempty"`
	Attempts         int32                `json:"attempts"`
	MaxAttempts      int32                `json:"max_attempts"`
	CreatedAt        string               `json:"created_at,omitempty"`
	ExpiresAt        string               `json:"expires_at,omitempty"`
	CompletedAt      string               `json:"completed_at,omitempty"`
	Card             *ElementsCardSummary `json:"card,omitempty"`
	LastError        *ElementsError       `json:"last_error,omitempty"`
}

// ElementsSessionWithSecret is a newly created session. Hand ClientSecret to the app that opens
// the card sheet; it is returned only once. Do not log or store it. String and GoString redact it.
type ElementsSessionWithSecret struct {
	ElementsSession
	ClientSecret    string `json:"client_secret"`
	ElementsBaseURL string `json:"elements_base_url,omitempty"`
}

func (s ElementsSessionWithSecret) String() string {
	return fmt.Sprintf("{ElementsSession:%+v ClientSecret:[redacted] ElementsBaseURL:%s}", s.ElementsSession, s.ElementsBaseURL)
}

func (s ElementsSessionWithSecret) GoString() string {
	return fmt.Sprintf("tapsilat.ElementsSessionWithSecret{ElementsSession:%#v, ClientSecret:\"[redacted]\", ElementsBaseURL:%q}", s.ElementsSession, s.ElementsBaseURL)
}

// CreateElementsSession opens a Tapsilat Elements session with POST /elements/sessions.
func (t *API) CreateElementsSession(ctx context.Context, payload ElementsSessionCreateRequest) (ElementsSessionWithSecret, error) {
	var response ElementsSessionWithSecret
	err := t.post(ctx, "/elements/sessions", payload, &response)
	return response, err
}

// GetElementsSession returns the session's authoritative state and outcome. Read it on your
// backend before fulfilling, rather than trusting what the app reports.
func (t *API) GetElementsSession(ctx context.Context, id string) (ElementsSession, error) {
	var response ElementsSession
	err := t.get(ctx, "/elements/sessions/"+url.PathEscape(id), &response)
	return response, err
}

// CancelElementsSession cancels a session that is still waiting for a card. A session that is
// confirming or has completed comes back unchanged.
func (t *API) CancelElementsSession(ctx context.Context, id string) (ElementsSession, error) {
	var response ElementsSession
	err := t.post(ctx, "/elements/sessions/"+url.PathEscape(id)+"/cancel", nil, &response)
	return response, err
}
