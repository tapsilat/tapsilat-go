package unit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

func TestAPICreation(t *testing.T) {
	t.Run("NewAPIWithToken", func(t *testing.T) {
		api := tapsilat.NewAPI("test_token")

		assert.NotNil(t, api)
		assert.Equal(t, "test_token", api.Token)
		assert.Equal(t, "https://acquiring.tapsilat.com/api/v1", api.EndPoint)
		assert.NotZero(t, api.Timeout)
	})

	t.Run("NewCustomAPI", func(t *testing.T) {
		customEndpoint := "https://custom.endpoint.com/api/v1"
		api := tapsilat.NewCustomAPI(customEndpoint, "custom_token")

		assert.NotNil(t, api)
		assert.Equal(t, "custom_token", api.Token)
		assert.Equal(t, customEndpoint, api.EndPoint)
		assert.NotZero(t, api.Timeout)
	})
}

func TestOrderStatusMap(t *testing.T) {
	t.Run("OrderStatusMapExists", func(t *testing.T) {
		statusMap := tapsilat.OrderStatuesMap

		assert.NotNil(t, statusMap)
		assert.Greater(t, len(statusMap), 0)

		// Check some known statuses
		var receivedFound, paidFound, cancelledFound bool
		for _, status := range statusMap {
			switch status.Status {
			case "Received":
				receivedFound = true
				assert.Equal(t, 1, status.Id)
			case "Paid":
				paidFound = true
				assert.Equal(t, 3, status.Id)
			case "Cancelled":
				cancelledFound = true
				assert.Equal(t, 8, status.Id)
			}
		}

		assert.True(t, receivedFound, "Received status should exist")
		assert.True(t, paidFound, "Paid status should exist")
		assert.True(t, cancelledFound, "Cancelled status should exist")
	})
}

func TestRefundOrder(t *testing.T) {
	t.Run("RefundOrderCreation", func(t *testing.T) {
		refund := tapsilat.RefundOrder{
			ReferenceID: "test_ref_123",
			Amount:      "100.50",
		}

		assert.Equal(t, "test_ref_123", refund.ReferenceID)
		assert.Equal(t, "100.50", refund.Amount)
	})
}

func TestCancelOrder(t *testing.T) {
	t.Run("CancelOrderCreation", func(t *testing.T) {
		cancel := tapsilat.CancelOrder{
			ReferenceID: "test_ref_123",
		}

		assert.Equal(t, "test_ref_123", cancel.ReferenceID)
	})
}

func TestOrderResponse(t *testing.T) {
	t.Run("OrderResponseCreation", func(t *testing.T) {
		response := tapsilat.OrderResponse{
			OrderID:     "order_123",
			ReferenceID: "ref_123",
			CheckoutURL: "https://checkout.example.com/order_123",
		}

		assert.Equal(t, "order_123", response.OrderID)
		assert.Equal(t, "ref_123", response.ReferenceID)
		assert.Equal(t, "https://checkout.example.com/order_123", response.CheckoutURL)
	})
}

func TestOrderStatus(t *testing.T) {
	t.Run("OrderStatusCreation", func(t *testing.T) {
		status := tapsilat.OrderStatus{
			Status: "Paid",
		}

		assert.Equal(t, "Paid", status.Status)
	})
}

func TestPaginatedData(t *testing.T) {
	t.Run("PaginatedDataCreation", func(t *testing.T) {
		data := tapsilat.PaginatedData{
			Page:       1,
			PerPage:    10,
			Total:      100,
			TotalPages: 10,
		}

		assert.Equal(t, int64(1), data.Page)
		assert.Equal(t, int64(10), data.PerPage)
		assert.Equal(t, int64(100), data.Total)
		assert.Equal(t, 10, data.TotalPages)
	})
}

func TestSubOrganizationDTO(t *testing.T) {
	t.Run("SubOrganizationDTOCreation", func(t *testing.T) {
		subOrg := tapsilat.SubOrganizationDTO{
			Acquirer:         "test_acquirer",
			Address:          "test_address",
			ContactFirstName: "John",
			ContactLastName:  "Doe",
			ContactEmail:     "john@doe.com",
			ContactPhone:     "+905551234567",
			Name:             "Test Organization",
			TaxOffice:        "Test Tax Office",
			TaxNumber:        "1234567890",
			Type:             "BUSINESS",
		}

		assert.Equal(t, "test_acquirer", subOrg.Acquirer)
		assert.Equal(t, "test_address", subOrg.Address)
		assert.Equal(t, "John", subOrg.ContactFirstName)
		assert.Equal(t, "Doe", subOrg.ContactLastName)
		assert.Equal(t, "john@doe.com", subOrg.ContactEmail)
		assert.Equal(t, "+905551234567", subOrg.ContactPhone)
		assert.Equal(t, "Test Organization", subOrg.Name)
		assert.Equal(t, "Test Tax Office", subOrg.TaxOffice)
		assert.Equal(t, "1234567890", subOrg.TaxNumber)
		assert.Equal(t, "BUSINESS", subOrg.Type)
	})
}

func TestValidationError(t *testing.T) {
	t.Run("ValidationErrorCreation", func(t *testing.T) {
		err := &tapsilat.ValidationError{
			StatusCode: 400,
			Code:       0,
			Message:    "Invalid input",
		}

		assert.Equal(t, 400, err.StatusCode)
		assert.Equal(t, 0, err.Code)
		assert.Equal(t, "Invalid input", err.Message)
		assert.Contains(t, err.Error(), "Tapsilat Validation Error")
		assert.Contains(t, err.Error(), "status_code:400")
		assert.Contains(t, err.Error(), "code:0")
		assert.Contains(t, err.Error(), "error:Invalid input")
	})
}
