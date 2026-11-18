package unit_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

func TestOrderBookingJSON(t *testing.T) {
	guests, refundable := int32(2), false
	body, err := json.Marshal(tapsilat.Order{Amount: 2500, Currency: "TRY", Booking: &tapsilat.OrderBooking{
		BookingReference: "BK-2026-104233", BookingType: "hotel", PropertyCountry: "TR",
		CheckInDate: "2026-12-20", CheckOutDate: "2026-12-23", GuestCount: &guests, IsRefundable: &refundable, BookingTotal: "12500.00",
	}})
	require.NoError(t, err)
	var wire map[string]json.RawMessage
	require.NoError(t, json.Unmarshal(body, &wire))
	require.JSONEq(t, `{"booking_reference":"BK-2026-104233","booking_type":"hotel","property_country":"TR","check_in_date":"2026-12-20",
		"check_out_date":"2026-12-23","guest_count":2,"is_refundable":false,"booking_total":"12500.00"}`, string(wire["booking"]),
		"an unset count stays out; a false refundability is sent")

	body, err = json.Marshal(tapsilat.Order{Amount: 2500, Currency: "TRY"})
	require.NoError(t, err)
	require.NotContains(t, string(body), `"booking"`)
}

func TestOrderDetailBookingAndFraudReviewJSON(t *testing.T) {
	var detail tapsilat.OrderDetail
	require.NoError(t, json.Unmarshal([]byte(`{"reference_id":"r","booking":{"booking_reference":"BK-1","room_count":1,"booking_total":"12500"}}`), &detail))
	require.NotNil(t, detail.Booking)
	require.Equal(t, "BK-1", detail.Booking.BookingReference)
	require.Equal(t, int32(1), *detail.Booking.RoomCount)
	require.Nil(t, detail.Booking.GuestCount)

	var status tapsilat.OrderStatus
	require.NoError(t, json.Unmarshal([]byte(`{"status":"On hold","fraud_review":"pending"}`), &status))
	require.Equal(t, tapsilat.FraudReviewPending, status.FraudReview)
}
