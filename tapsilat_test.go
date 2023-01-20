package tapsilat_test

import (
	"os"
	"testing"

	"github.com/tapsilat/tapsilat-go"
)

func TestCreateOrder(t *testing.T) {
	token := os.Getenv("TAPSILAT_TOKEN")
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
	if err != nil {
		t.Error(err)
	}

	if response.ReferenceID == "" {
		t.Error("ReferenceID is empty")
	}
	if response.OrderID == "" {
		t.Error("OrderID is empty")
	}

}

func TestGetOrder(t *testing.T) {
	token := os.Getenv("TAPSILAT_TOKEN")
	api := tapsilat.NewAPI(token)

	order, err := api.GetOrder("e0176c98-fb41-4f08-aa03-55bb8d7bb9d6")
	if err != nil {
		t.Error(err)
	}

	if order.Locale == "" {
		t.Error("Locale is empty")
	}
	if order.Currency == "" {
		t.Error("Currency is empty")
	}

}

func TestGetOrderStatus(t *testing.T) {
	token := os.Getenv("TAPSILAT_TOKEN")
	api := tapsilat.NewAPI(token)

	order, err := api.GetOrderStatus("e0176c98-fb41-4f08-aa03-55bb8d7bb9d6")
	if err != nil {
		t.Error(err)
	}

	if order.Status == "" {
		t.Error("Status is can not be empty")
	}
	if order.Status == "Waiting For Payment" {
		t.Log("Waiting For Payment")
	}
	t.Log("status", order.Status)
}
