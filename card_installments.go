package tapsilat

import (
	"context"
	"net/url"
	"strconv"
)

// CardInstallmentsRequest selects installment options for a card BIN.
type CardInstallmentsRequest struct {
	// Bin is the card's first digits (at least 8).
	Bin    string
	Amount float64
	// Currency is a currency unit or name, for example "TRY".
	Currency string
	IsThreeD bool
}

// CardInstallment is one installment option with its commission.
type CardInstallment struct {
	Installment    uint64  `json:"installment,omitempty"`
	Amount         string  `json:"amount,omitempty"`
	TotalAmount    string  `json:"total_amount,omitempty"`
	CommissionRate float64 `json:"commission_rate,omitempty"`
	MonthlyPayment string  `json:"monthly_payment,omitempty"`
}

// CardInstallmentsResponse lists the installment options available for a card.
type CardInstallmentsResponse struct {
	IsInstallmentAvailable bool              `json:"is_installment_available,omitempty"`
	Installments           []CardInstallment `json:"installments,omitempty"`
}

// InstallmentTableRequest selects the installment table for an amount.
type InstallmentTableRequest struct {
	Amount float64
	// Currency is a currency unit or name, for example "TRY".
	Currency string
	IsThreeD bool
}

// InstallmentTableRow is one installment count within a card program.
type InstallmentTableRow struct {
	Installment       uint64  `json:"installment,omitempty"`
	InstallmentAmount string  `json:"installment_amount,omitempty"`
	TotalAmount       string  `json:"total_amount,omitempty"`
	CommissionRate    float64 `json:"commission_rate,omitempty"`
	IsInterestFree    bool    `json:"is_interest_free,omitempty"`
}

// InstallmentTableGroup groups installment rows by bank card program.
type InstallmentTableGroup struct {
	Bank        string                `json:"bank,omitempty"`
	Program     string                `json:"program,omitempty"`
	ProgramCode string                `json:"program_code,omitempty"`
	Rows        []InstallmentTableRow `json:"rows,omitempty"`
}

// InstallmentTableResponse is the installment table for every card program.
type InstallmentTableResponse struct {
	Groups []InstallmentTableGroup `json:"groups,omitempty"`
}

// GetCardInstallments returns the installment options and commission rates for a card BIN and amount.
func (t *API) GetCardInstallments(ctx context.Context, req CardInstallmentsRequest) (CardInstallmentsResponse, error) {
	var response CardInstallmentsResponse
	query := url.Values{
		"bin":      {req.Bin},
		"amount":   {strconv.FormatFloat(req.Amount, 'f', -1, 64)},
		"currency": {req.Currency},
	}
	if req.IsThreeD {
		query.Set("is_three_d", "true")
	}
	err := t.get(ctx, "/card/installments?"+query.Encode(), &response)
	return response, err
}

// GetInstallmentTable returns installment options grouped by bank card program, with an interest-free flag per row.
func (t *API) GetInstallmentTable(ctx context.Context, req InstallmentTableRequest) (InstallmentTableResponse, error) {
	var response InstallmentTableResponse
	query := url.Values{
		"amount":   {strconv.FormatFloat(req.Amount, 'f', -1, 64)},
		"currency": {req.Currency},
	}
	if req.IsThreeD {
		query.Set("is_three_d", "true")
	}
	err := t.get(ctx, "/card/installment-table?"+query.Encode(), &response)
	return response, err
}
