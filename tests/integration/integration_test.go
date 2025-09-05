package integration_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tapsilat/tapsilat-go"
)

func TestCreateOrder(t *testing.T) {
	// Integration testler için gerçek bir token gerekiyor
	// Bu test sadece gerçek token ile çalışır
	token := "your_test_token_here"
	if token == "your_test_token_here" {
		t.Skip("Please set a real token for integration tests")
	}

	api := tapsilat.NewAPI(token)
	order := tapsilat.Order{
		Locale:    "tr",
		Currency:  "TRY",
		Amount:    5,
		TaxAmount: 0.18,
		Buyer: tapsilat.OrderBuyer{
			Id:                  "123456789",
			Name:                "John",
			Surname:             "Doe",
			Email:               "john@doe.com",
			GsmNumber:           "5555555555",
			IdentityNumber:      "12345678901",
			RegistrationDate:    "2023-01-01",
			RegistrationAddress: "Istanbul",
			LastLoginDate:       "2023-01-01",
			City:                "Istanbul",
			Country:             "Türkiye",
			ZipCode:             "34000",
			Ip:                  "127.0.0.1",
			BirdthDate:          "1990-01-01",
		},
		ShippingAddress: tapsilat.OrderShippingAddress{
			Address:      "Istanbul",
			ZipCode:      "34000",
			City:         "Istanbul",
			Country:      "Türkiye",
			ContactName:  "John Doe",
			TrackingCode: "123456789",
		},
		BillingAddress: tapsilat.OrderBillingAddress{
			Address:     "Istanbul",
			ZipCode:     "34000",
			City:        "Istanbul",
			Country:     "Türkiye",
			ContactName: "John Doe",
		},
		BasketItems: []tapsilat.OrderBasketItem{
			{
				Id:        "1",
				Name:      "Product 1",
				Price:     5,
				Category1: "Category 1",
				Category2: "Category 2",
				ItemType:  "VIRTUAL",
			},
			{
				Id:        "2",
				Name:      "Product 2",
				Price:     5,
				Category1: "Category 1",
				Category2: "Category 2",
				ItemType:  "VIRTUAL",
			},
		},
	}
	response, err := api.CreateOrder(order)
	require.NoError(t, err, "CreateOrder should not return an error")

	assert.NotEmpty(t, response.ReferenceID, "ReferenceID should not be empty")
	assert.NotEmpty(t, response.OrderID, "OrderID should not be empty")
	assert.NotEmpty(t, response.CheckoutURL, "CheckoutURL should not be empty")

	t.Logf("Order created - ID: %s, Reference: %s, CheckoutURL: %s", response.OrderID, response.ReferenceID, response.CheckoutURL)
}

func TestGetOrder(t *testing.T) {
	// Integration testler için gerçek bir token gerekiyor
	token := "your_test_token_here"
	if token == "your_test_token_here" {
		t.Skip("Please set a real token for integration tests")
	}

	api := tapsilat.NewAPI(token)

	order, err := api.GetOrder("e0176c98-fb41-4f08-aa03-55bb8d7bb9d6")
	require.NoError(t, err, "GetOrder should not return an error")

	assert.NotEmpty(t, order.Locale, "Locale should not be empty")
	assert.NotEmpty(t, order.Currency, "Currency should not be empty")
}

func TestGetOrderStatus(t *testing.T) {
	// Integration testler için gerçek bir token gerekiyor
	token := "your_test_token_here"
	if token == "your_test_token_here" {
		t.Skip("Please set a real token for integration tests")
	}

	api := tapsilat.NewAPI(token)

	order, err := api.GetOrderStatus("e0176c98-fb41-4f08-aa03-55bb8d7bb9d6")
	require.NoError(t, err, "GetOrderStatus should not return an error")

	assert.NotEmpty(t, order.Status, "Status should not be empty")

	if order.Status == "Waiting For Payment" {
		t.Log("Status: Waiting For Payment")
	}
	t.Logf("Status: %s", order.Status)
}
