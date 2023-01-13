# tapsilat-go
Tapsilat Go is a Go client library for accessing the Tapsilat API. It provides convenient access to Tapsilat's API from applications written in the Go language.
You can create an order, get order, get order list, get order status, cancel order and refund order.

## Installation

```bash
go get github.com/tapsilat/tapsilat-go
```

## Usage

First, you need to create a client instance using the `tapsilat.NewAPI()` function. You need to pass your API token to the function.

If you have self-hosted Tapsilat, you need use the `tapsilat.NewCustomAPI()` function to create a client instance. You need to pass your tapilat url and API token to the function.
```go
package main

import "github.com/tapsilat/tapsilat-go-sdk"

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
}
```

```go
package main

import "github.com/tapsilat/tapsilat-go-sdk"

func main() {
	token := "TOKEN"
	selfHostedTapsilatURL := "https://your-tapsilat-url"
	api := tapsilat.NewCustomAPI(selfHostedTapsilatURL, token)
}
```

### Create a Order


```go
package main

import "github.com/tapsilat/tapsilat-go-sdk"

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	order := tapsilat.Order{
		Locale:   "tr",
		Currency: "TRY",
		Amount:   5,
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
		panic(err)
	}
	println(response)
}
```


### Get Order

```go
package main

import "github.com/tapsilat/tapsilat-go-sdk"

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	order, err := api.GetOrder("order_reference_id")
	if err != nil {
		panic(err)
	}
	println(order)
}

```

### Get Order List

```go
package main

import "github.com/tapsilat/tapsilat-go-sdk"

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	orders, err := api.GetOrders("page","limit")
	if err != nil {
		panic(err)
	}
	println(orders)
}

```

### Get Order Status

```go
package main

import "github.com/tapsilat/tapsilat-go-sdk"

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	status, err := api.GetOrderStatus("order_reference_id")
	if err != nil {
		panic(err)
	}
	println(status)
}

```

### Cancel Order

```go
package main

import "github.com/tapsilat/tapsilat-go-sdk"

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	payload := tapsilat.CancelOrder{
		ReferenceID: "order_reference_id",
	}
	status, err := api.CancelOrder(payload)
	if err != nil {
		panic(err)
	}
	println(status)
}

```

### Refund Order

```go
package main

import "github.com/tapsilat/tapsilat-go-sdk"

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	payload := tapsilat.RefundOrder{
		ReferenceID: "order_reference_id",
		Amount:      "100",
	}
	status, err := api.RefundOrder(payload)
	if err != nil {
		panic(err)
	}
	println(status)
}

```
