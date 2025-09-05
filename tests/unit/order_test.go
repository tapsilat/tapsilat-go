package unit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

func TestOrderCreation(t *testing.T) {
	t.Run("BasicOrderCreation", func(t *testing.T) {
		buyer := tapsilat.OrderBuyer{
			Name:    "John",
			Surname: "Doe",
			Email:   "test@example.com",
		}

		order := tapsilat.Order{
			Amount:   100.0,
			Currency: "TRY",
			Locale:   "tr",
			Buyer:    buyer,
		}

		assert.Equal(t, 100.0, order.Amount)
		assert.Equal(t, "TRY", order.Currency)
		assert.Equal(t, "tr", order.Locale)
		assert.Equal(t, "John", order.Buyer.Name)
		assert.Equal(t, "Doe", order.Buyer.Surname)
		assert.Equal(t, "test@example.com", order.Buyer.Email)
	})

	t.Run("OrderWithBasketItems", func(t *testing.T) {
		buyer := tapsilat.OrderBuyer{
			Name:    "John",
			Surname: "Doe",
			Email:   "test@example.com",
		}

		payer1 := &tapsilat.OrderBasketItemPayer{
			ReferenceID: "payer_ref0_item1",
			Type:        "PERSONAL",
		}

		payer2 := &tapsilat.OrderBasketItemPayer{
			ReferenceID: "payer_ref1_item2",
			Type:        "BUSINESS",
		}

		quantity1 := 1
		quantity2 := 2

		basketItem1 := tapsilat.OrderBasketItem{
			Id:       "B001",
			Name:     "Item 1",
			Price:    10.00,
			Quantity: &quantity1,
			ItemType: "PHYSICAL",
			Payer:    payer1,
		}

		basketItem2 := tapsilat.OrderBasketItem{
			Id:       "B002",
			Name:     "Item 2",
			Price:    20.49,
			Quantity: &quantity2,
			ItemType: "PHYSICAL",
			Payer:    payer2,
		}

		order := tapsilat.Order{
			Amount:      30.49,
			Currency:    "TRY",
			Locale:      "tr",
			Buyer:       buyer,
			BasketItems: []tapsilat.OrderBasketItem{basketItem1, basketItem2},
		}

		assert.Equal(t, 30.49, order.Amount)
		assert.Len(t, order.BasketItems, 2)
		assert.Equal(t, "B001", order.BasketItems[0].Id)
		assert.Equal(t, "Item 1", order.BasketItems[0].Name)
		assert.Equal(t, 10.00, order.BasketItems[0].Price)
		assert.Equal(t, "PERSONAL", order.BasketItems[0].Payer.Type)
		assert.Equal(t, "B002", order.BasketItems[1].Id)
		assert.Equal(t, "BUSINESS", order.BasketItems[1].Payer.Type)
	})

	t.Run("OrderWithAddresses", func(t *testing.T) {
		buyer := tapsilat.OrderBuyer{
			Name:    "John",
			Surname: "Doe",
			Email:   "test@example.com",
		}

		billingAddress := tapsilat.OrderBillingAddress{
			Address:     "uskudar",
			City:        "Istanbul",
			Country:     "TR",
			ContactName: "John Doe",
			ZipCode:     "34000",
		}

		shippingAddress := tapsilat.OrderShippingAddress{
			Address:     "kadikoy",
			City:        "Istanbul",
			Country:     "TR",
			ContactName: "Jane Doe",
			ZipCode:     "34001",
		}

		order := tapsilat.Order{
			Amount:          25.00,
			Currency:        "TRY",
			Locale:          "tr",
			Buyer:           buyer,
			BillingAddress:  billingAddress,
			ShippingAddress: shippingAddress,
		}

		assert.Equal(t, 25.00, order.Amount)
		assert.Equal(t, "uskudar", order.BillingAddress.Address)
		assert.Equal(t, "Istanbul", order.BillingAddress.City)
		assert.Equal(t, "TR", order.BillingAddress.Country)
		assert.Equal(t, "kadikoy", order.ShippingAddress.Address)
		assert.Equal(t, "Jane Doe", order.ShippingAddress.ContactName)
	})

	t.Run("OrderWithInstallmentsAndPaymentMethods", func(t *testing.T) {
		buyer := tapsilat.OrderBuyer{
			Name:    "John",
			Surname: "Doe",
			Email:   "test@example.com",
		}

		order := tapsilat.Order{
			Amount:              1200.00,
			Currency:            "TRY",
			Locale:              "tr",
			Buyer:               buyer,
			EnabledInstallments: []int{2, 3, 6, 9},
			PaymentMethods:      true,
			PaymentOptions:      []string{"credit_card", "cash"},
			PaymentSuccessUrl:   "https://example.com/install_success_s8",
			PaymentFailureUrl:   "https://example.com/install_failure_s8",
		}

		assert.Equal(t, 1200.00, order.Amount)
		assert.Equal(t, []int{2, 3, 6, 9}, order.EnabledInstallments)
		assert.True(t, order.PaymentMethods)
		assert.Equal(t, []string{"credit_card", "cash"}, order.PaymentOptions)
		assert.Equal(t, "https://example.com/install_success_s8", order.PaymentSuccessUrl)
		assert.Equal(t, "https://example.com/install_failure_s8", order.PaymentFailureUrl)
	})

	t.Run("OrderWithCheckoutDesign", func(t *testing.T) {
		buyer := tapsilat.OrderBuyer{
			Name:    "John",
			Surname: "Doe",
			Email:   "test@example.com",
		}

		design := tapsilat.OrderCheckoutDesign{
			PayButtonColor:       "#FF0000",
			Logo:                 "http://example.com/logo.png",
			InputBackgroundColor: "#EEEEEE",
			InputTextColor:       "#333333",
			RightBackgroundColor: "#FAFAFA",
		}

		order := tapsilat.Order{
			Amount:         55.00,
			Currency:       "TRY",
			Locale:         "tr",
			Buyer:          buyer,
			CheckoutDesign: design,
		}

		assert.Equal(t, 55.00, order.Amount)
		assert.Equal(t, "#FF0000", order.CheckoutDesign.PayButtonColor)
		assert.Equal(t, "http://example.com/logo.png", order.CheckoutDesign.Logo)
		assert.Equal(t, "#EEEEEE", order.CheckoutDesign.InputBackgroundColor)
		assert.Equal(t, "#333333", order.CheckoutDesign.InputTextColor)
		assert.Equal(t, "#FAFAFA", order.CheckoutDesign.RightBackgroundColor)
	})

	t.Run("OrderWithMetadata", func(t *testing.T) {
		buyer := tapsilat.OrderBuyer{
			Name:    "John",
			Surname: "Doe",
			Email:   "test@example.com",
		}

		metadata := []tapsilat.OrderMetadata{
			{Key: "customer_id", Value: "12345"},
			{Key: "source", Value: "mobile_app"},
		}

		order := tapsilat.Order{
			Amount:   100.0,
			Currency: "TRY",
			Locale:   "tr",
			Buyer:    buyer,
			Metadata: metadata,
		}

		assert.Equal(t, 100.0, order.Amount)
		assert.Len(t, order.Metadata, 2)
		assert.Equal(t, "customer_id", order.Metadata[0].Key)
		assert.Equal(t, "12345", order.Metadata[0].Value)
		assert.Equal(t, "source", order.Metadata[1].Key)
		assert.Equal(t, "mobile_app", order.Metadata[1].Value)
	})

	t.Run("OrderWithPaymentTerms", func(t *testing.T) {
		buyer := tapsilat.OrderBuyer{
			Name:    "John",
			Surname: "Doe",
			Email:   "test@example.com",
		}

		amount1 := 50.0
		required1 := true
		sequence1 := 1

		amount2 := 50.0
		required2 := false
		sequence2 := 2

		paymentTerms := []tapsilat.OrderPaymentTerm{
			{
				Amount:          &amount1,
				DueDate:         "2024-01-15",
				Required:        &required1,
				TermSequence:    &sequence1,
				Status:          "pending",
				TermReferenceID: "term_ref_1",
			},
			{
				Amount:          &amount2,
				DueDate:         "2024-02-15",
				Required:        &required2,
				TermSequence:    &sequence2,
				Status:          "pending",
				TermReferenceID: "term_ref_2",
			},
		}

		order := tapsilat.Order{
			Amount:       100.0,
			Currency:     "TRY",
			Locale:       "tr",
			Buyer:        buyer,
			PaymentTerms: paymentTerms,
		}

		assert.Equal(t, 100.0, order.Amount)
		assert.Len(t, order.PaymentTerms, 2)
		assert.Equal(t, 50.0, *order.PaymentTerms[0].Amount)
		assert.Equal(t, "2024-01-15", order.PaymentTerms[0].DueDate)
		assert.True(t, *order.PaymentTerms[0].Required)
		assert.Equal(t, 1, *order.PaymentTerms[0].TermSequence)
		assert.Equal(t, "pending", order.PaymentTerms[0].Status)
		assert.Equal(t, "term_ref_1", order.PaymentTerms[0].TermReferenceID)
	})
}

func TestPaymentTermDTOs(t *testing.T) {
	t.Run("OrderPaymentTermCreateDTO", func(t *testing.T) {
		sequence := 1
		createDTO := tapsilat.OrderPaymentTermCreateDTO{
			OrderReferenceID: "order_123",
			Amount:           100.0,
			DueDate:          "2024-01-15",
			Required:         true,
			Data:             "test_data",
			TermSequence:     &sequence,
		}

		assert.Equal(t, "order_123", createDTO.OrderReferenceID)
		assert.Equal(t, 100.0, createDTO.Amount)
		assert.Equal(t, "2024-01-15", createDTO.DueDate)
		assert.True(t, createDTO.Required)
		assert.Equal(t, "test_data", createDTO.Data)
		assert.Equal(t, 1, *createDTO.TermSequence)
	})

	t.Run("OrderPaymentTermUpdateDTO", func(t *testing.T) {
		amount := 150.0
		required := false

		updateDTO := tapsilat.OrderPaymentTermUpdateDTO{
			TermReferenceID: "term_123",
			Amount:          &amount,
			DueDate:         "2024-02-15",
			Required:        &required,
			Data:            "updated_data",
		}

		assert.Equal(t, "term_123", updateDTO.TermReferenceID)
		assert.Equal(t, 150.0, *updateDTO.Amount)
		assert.Equal(t, "2024-02-15", updateDTO.DueDate)
		assert.False(t, *updateDTO.Required)
		assert.Equal(t, "updated_data", updateDTO.Data)
	})

	t.Run("OrderTermRefundRequest", func(t *testing.T) {
		amount := 75.0
		refundRequest := tapsilat.OrderTermRefundRequest{
			TermReferenceID: "term_123",
			Amount:          &amount,
		}

		assert.Equal(t, "term_123", refundRequest.TermReferenceID)
		assert.Equal(t, 75.0, *refundRequest.Amount)
	})
}

func TestBasketItemPayer(t *testing.T) {
	t.Run("BasketItemPayerCreation", func(t *testing.T) {
		payer := tapsilat.OrderBasketItemPayer{
			Address:     "uskudar",
			Type:        "PERSONAL",
			ReferenceID: "123456789",
			TaxOffice:   "Test Tax Office",
			Title:       "Test Title",
			VAT:         "1234567890",
		}

		assert.Equal(t, "uskudar", payer.Address)
		assert.Equal(t, "PERSONAL", payer.Type)
		assert.Equal(t, "123456789", payer.ReferenceID)
		assert.Equal(t, "Test Tax Office", payer.TaxOffice)
		assert.Equal(t, "Test Title", payer.Title)
		assert.Equal(t, "1234567890", payer.VAT)
	})
}

func TestBillingAddress(t *testing.T) {
	t.Run("BillingAddressCreation", func(t *testing.T) {
		billingAddress := tapsilat.OrderBillingAddress{
			BillingType:  "BUSINESS",
			Citizenship:  "TR",
			Title:        "Test Company Ltd.",
			TaxOffice:    "Merter",
			Address:      "Test Address",
			ZipCode:      "34000",
			City:         "Istanbul",
			District:     "Uskudar",
			Country:      "Turkey",
			ContactName:  "John Doe",
			ContactPhone: "+905555555555",
			VatNumber:    "1234567890",
		}

		assert.Equal(t, "BUSINESS", billingAddress.BillingType)
		assert.Equal(t, "TR", billingAddress.Citizenship)
		assert.Equal(t, "Test Company Ltd.", billingAddress.Title)
		assert.Equal(t, "Merter", billingAddress.TaxOffice)
		assert.Equal(t, "Test Address", billingAddress.Address)
		assert.Equal(t, "34000", billingAddress.ZipCode)
		assert.Equal(t, "Istanbul", billingAddress.City)
		assert.Equal(t, "Uskudar", billingAddress.District)
		assert.Equal(t, "Turkey", billingAddress.Country)
		assert.Equal(t, "John Doe", billingAddress.ContactName)
		assert.Equal(t, "+905555555555", billingAddress.ContactPhone)
		assert.Equal(t, "1234567890", billingAddress.VatNumber)
	})
}
