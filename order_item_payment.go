package tapsilat

import (
	"context"
	"encoding/json"
	"errors"
)

const (
	OrderItemPaymentTypeVpos          uint64 = 1
	OrderItemPaymentTypeMoneyTransfer uint64 = 2
	OrderItemPaymentTypeWallet        uint64 = 3
	OrderItemPaymentTypeLoan          uint64 = 4
	OrderItemPaymentTypePin           uint64 = 5
	OrderItemPaymentTypeRemittance    uint64 = 6
	OrderItemPaymentTypeMetropol      uint64 = 7
	OrderItemPaymentTypeTgn           uint64 = 8
	OrderItemPaymentTypeVIBAN         uint64 = 9
	OrderItemPaymentTypePersonalLoan  uint64 = 10
	OrderItemPaymentTypeStripe        uint64 = 11
	OrderItemPaymentTypeStripeVBAN    uint64 = 12
)

var ErrOrderIDNotFound = errors.New("tapsilat: order id not found for reference")

type OrderOIPRequest struct {
	OrderID      string  `json:"order_id"`
	BasketItemID string  `json:"basket_item_id"`
	Amount       float64 `json:"amount"`
	Type         uint64  `json:"type"`
	WalletID     string  `json:"wallet_id,omitempty"`
	PaidDate     string  `json:"paid_date,omitempty"`
}

func (t *API) AddOrderOIP(ctx context.Context, payload OrderOIPRequest) (map[string]any, error) {
	var response map[string]any
	err := t.post(ctx, "/order/oip", payload, &response)
	return response, err
}

func (t *API) GetOrderID(ctx context.Context, referenceID string) (string, error) {
	detail, err := t.GetOrder(ctx, referenceID)
	if err != nil {
		return "", err
	}
	if detail.ConversationID == "" {
		return "", ErrOrderIDNotFound
	}

	list, err := t.GetOrderListWithFilter(ctx, OrderListFilter{Page: 1, PerPage: 50, ConversationID: detail.ConversationID})
	if err != nil {
		return "", err
	}

	raw, err := json.Marshal(list.Rows)
	if err != nil {
		return "", err
	}
	var rows []struct {
		ID          string `json:"id"`
		ReferenceID string `json:"reference_id"`
	}
	if err := json.Unmarshal(raw, &rows); err != nil {
		return "", err
	}
	for _, row := range rows {
		if row.ReferenceID == referenceID && row.ID != "" {
			return row.ID, nil
		}
	}
	return "", ErrOrderIDNotFound
}
