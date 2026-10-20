package tapsilat

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
)

// PaidOrderStatuses are the statuses of an order whose payment was captured.
var PaidOrderStatuses = []int{OrderStatusPaid, OrderStatusCompleted, OrderStatusSubMerchantPaymentApproved}

// FailedOrderStatuses are the statuses of an order that will not be paid.
var FailedOrderStatuses = []int{
	OrderStatusCancelled, OrderStatusFraud, OrderStatusRejected, OrderStatusFailure, OrderStatusExpired, OrderStatusTerminated,
}

// RefundedOrderStatuses are the statuses of an order whose payment was refunded in full or in part.
var RefundedOrderStatuses = []int{OrderStatusRefunded, OrderStatusPartiallyRefunded}

// OrderStatusCode returns the status code for a code (any integer type, an integral float64 or a
// json.Number), a numeric string such as order_lite callbacks send, or a status name such as
// "Paid" (case-insensitive). It returns 0 when the value is not a status.
func OrderStatusCode(status any) int {
	if s, ok := status.(string); ok {
		s = strings.TrimSpace(s)
		if !isDigits(s) {
			return GetOrderStatusByStr(s)
		}
		code, err := strconv.Atoi(s)
		if err != nil {
			return 0
		}
		status = code
	}
	code, ok := integer(status)
	if !ok || code < 1 || code > math.MaxInt32 || OrderStatusName(int(code)) == "" {
		return 0
	}
	return int(code)
}

// OrderStatusName returns the name of a status code, such as "Paid", or "" when it is not a status.
func OrderStatusName(code int) string {
	for _, v := range OrderStatuesMap {
		if v.Id == code {
			return v.Status
		}
	}
	return ""
}

// IsPaidOrderStatus reports whether status (see OrderStatusCode) is in PaidOrderStatuses.
func IsPaidOrderStatus(status any) bool {
	return hasOrderStatus(PaidOrderStatuses, status)
}

// IsFailedOrderStatus reports whether status (see OrderStatusCode) is in FailedOrderStatuses.
func IsFailedOrderStatus(status any) bool {
	return hasOrderStatus(FailedOrderStatuses, status)
}

// IsRefundedOrderStatus reports whether status (see OrderStatusCode) is in RefundedOrderStatuses.
func IsRefundedOrderStatus(status any) bool {
	return hasOrderStatus(RefundedOrderStatuses, status)
}

func hasOrderStatus(group []int, status any) bool {
	code := OrderStatusCode(status)
	for _, c := range group {
		if c == code && code != 0 {
			return true
		}
	}
	return false
}

func integer(value any) (int64, bool) {
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint32:
		return int64(v), true
	case uint64:
		if v > math.MaxInt64 {
			return 0, false
		}
		return int64(v), true
	case float64:
		if v != math.Trunc(v) || v < math.MinInt64 || v >= math.MaxInt64 {
			return 0, false
		}
		return int64(v), true
	case json.Number:
		if n, err := v.Int64(); err == nil {
			return n, true
		}
		f, err := v.Float64()
		if err != nil {
			return 0, false
		}
		return integer(f)
	}
	return 0, false
}
