# tapsilat-go

Tapsilat Go is a Go client library for accessing the Tapsilat API. It provides convenient access to Tapsilat's API from applications written in the Go language. You can create an order, get order, get order list, get order status, cancel order, refund order, and manage subscriptions.

## Installation

```bash
go get github.com/tapsilat/tapsilat-go
```

## Configuration

```go
package main

import (
    "github.com/tapsilat/tapsilat-go"
)

func main() {
    // Create API client with your token
    api := tapsilat.NewAPI("your_token_here")

    // Or with custom endpoint
    api := tapsilat.NewCustomAPI("https://custom.endpoint.com/api/v1", "your_token")
}
```

## Context Usage

All API methods require a `context.Context` parameter. This allows you to control request timeouts and cancellation:

```go
import (
    "context"
    "time"
    "github.com/tapsilat/tapsilat-go"
)

// Simple usage with Background context
response, err := api.CreateOrder(context.Background(), order)

// With timeout - request will be cancelled after 10 seconds
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
response, err := api.CreateOrder(ctx, order)

// With cancellation - you can cancel the request manually
ctx, cancel := context.WithCancel(context.Background())
go func() {
    // Cancel after some condition
    time.Sleep(5 * time.Second)
    cancel()
}()
response, err := api.CreateOrder(ctx, order)
```

## Usage Examples

### Basic Order Creation

```go
package main

import (
	"context"
	"github.com/tapsilat/tapsilat-go"
)

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
	response, err := api.CreateOrder(context.Background(), order)
	if err != nil {
		panic(err)
	}
	println("Order created successfully!")
	println("Order ID:", response.OrderID)
	println("Reference ID:", response.ReferenceID)
	println("Checkout URL:", response.CheckoutURL)
}
```

### Order with Basket Items

```go
quantity := 2
basketItem := tapsilat.OrderBasketItem{
    Id:       "item_001",
    Name:     "Product Name",
    Price:    50.0,
    Quantity: &quantity,
    ItemType: "PHYSICAL",
}

order := tapsilat.Order{
    Locale:      "tr",
    Currency:    "TRY",
    Amount:      100.0,
    BasketItems: []tapsilat.OrderBasketItem{basketItem},
    Buyer: tapsilat.OrderBuyer{
        Name:    "John",
        Surname: "Doe",
        Email:   "john@doe.com",
    },
}
```

### Order with Payment Terms (Installments)

```go
amount1 := 50.0
required := true
sequence := 1

paymentTerm := tapsilat.OrderPaymentTerm{
    Amount:          &amount1,
    DueDate:         "2024-01-15",
    Required:        &required,
    TermSequence:    &sequence,
    Status:          "pending",
    TermReferenceID: "term_ref_1",
}

order := tapsilat.Order{
    Locale:       "tr",
    Currency:     "TRY",
    Amount:       100.0,
    PaymentTerms: []tapsilat.OrderPaymentTerm{paymentTerm},
    Buyer: tapsilat.OrderBuyer{
        Name:    "John",
        Surname: "Doe",
        Email:   "john@doe.com",
    },
}
```

### Validation

The SDK includes built-in validation for common fields:

```go
// GSM Number Validation
cleanGSM, err := tapsilat.ValidateGSMNumber("+90 555 123-45-67")
if err != nil {
    log.Fatal(err)
}
fmt.Println(cleanGSM) // Output: +905551234567

// Installments Validation
installments, err := tapsilat.ValidateInstallments("1,2,3,6")
if err != nil {
    log.Fatal(err)
}
fmt.Println(installments) // Output: [1 2 3 6]
```

### Checkout URLs

When you create an order, the response automatically includes a checkout URL that you can use to redirect customers for payment:

```go
response, err := api.CreateOrder(context.Background(), order)
if err != nil {
    log.Fatal(err)
}

// Checkout URL is automatically included in the response
fmt.Printf("Order ID: %s\n", response.OrderID)
fmt.Printf("Reference ID: %s\n", response.ReferenceID)
fmt.Printf("Checkout URL: %s\n", response.CheckoutURL)

// You can also get the checkout URL separately if needed
checkoutURL, err := api.GetCheckoutURL(context.Background(), response.ReferenceID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Checkout URL: %s\n", checkoutURL)
```

### Get Order

```go
package main

import (
	"context"
	"github.com/tapsilat/tapsilat-go"
)

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	order, err := api.GetOrder(context.Background(), "order_reference_id")
	if err != nil {
		panic(err)
	}
	println(order)
}
```

### Get Order List

```go
package main

import (
	"context"
	"github.com/tapsilat/tapsilat-go"
)

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	orders, err := api.GetOrders(context.Background(), "page", "limit", "")
	if err != nil {
		panic(err)
	}
	println(orders)
}
```

### Get Order Status

```go
package main

import (
	"context"
	"github.com/tapsilat/tapsilat-go"
)

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	status, err := api.GetOrderStatus(context.Background(), "order_reference_id")
	if err != nil {
		panic(err)
	}
	println(status)
}
```

### Cancel Order

```go
package main

import (
	"context"
	"github.com/tapsilat/tapsilat-go"
)

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	payload := tapsilat.CancelOrder{
		ReferenceID: "order_reference_id",
	}
	status, err := api.CancelOrder(context.Background(), payload)
	if err != nil {
		panic(err)
	}
	println(status)
}
```

### Refund Order

```go
package main

import (
	"context"
	"github.com/tapsilat/tapsilat-go"
)

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	payload := tapsilat.RefundOrder{
		ReferenceID: "order_reference_id",
		Amount:      "100",
	}
	status, err := api.RefundOrder(context.Background(), payload)
	if err != nil {
		panic(err)
	}
	println(status)
}
```

### Subscription Operations

#### Create Subscription

```go
package main

import (
	"context"
	"github.com/tapsilat/tapsilat-go"
)

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	
	subscription := tapsilat.SubscriptionCreateRequest{
		Amount:              100.0,
		Currency:            "TRY",
		Title:               "Monthly Subscription",
		Period:              30,
		Cycle:               1,
		PaymentDate:         1,
		ExternalReferenceID: "ext_sub_123",
		SuccessURL:          "https://example.com/success",
		FailureURL:          "https://example.com/failure",
		CardID:              "card_token_123",
		Billing: tapsilat.SubscriptionBilling{
			Address:     "Istanbul",
			City:        "Istanbul",
			Country:     "TR",
			ZipCode:     "34000",
			ContactName: "John Doe",
			VatNumber:   "1234567890",
		},
		User: tapsilat.SubscriptionUser{
			ID:             "user_123",
			FirstName:      "John",
			LastName:       "Doe",
			Email:          "john@doe.com",
			Phone:          "5555555555",
			IdentityNumber: "12345678901",
			Address:        "Istanbul",
			City:           "Istanbul",
			Country:        "TR",
			ZipCode:        "34000",
		},
	}
	
	response, err := api.CreateSubscription(context.Background(), subscription)
	if err != nil {
		panic(err)
	}
	println("Subscription created successfully!")
	println("Reference ID:", response.ReferenceID)
	println("Order Reference ID:", response.OrderReferenceID)
}
```

#### Get Subscription

```go
package main

import (
	"context"
	"github.com/tapsilat/tapsilat-go"
)

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	
	payload := tapsilat.SubscriptionGetRequest{
		ReferenceID: "subscription_reference_id",
		// Or use ExternalReferenceID: "ext_sub_123",
	}
	
	subscription, err := api.GetSubscription(context.Background(), payload)
	if err != nil {
		panic(err)
	}
	println("Subscription Title:", subscription.Title)
	println("Amount:", subscription.Amount)
	println("Is Active:", subscription.IsActive)
	println("Payment Status:", subscription.PaymentStatus)
}
```

#### List Subscriptions

```go
package main

import (
	"context"
	"encoding/json"
	"github.com/tapsilat/tapsilat-go"
)

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	
	// Get first page with 10 items per page
	subscriptions, err := api.ListSubscriptions(context.Background(), 1, 10)
	if err != nil {
		panic(err)
	}
	
	println("Total subscriptions:", subscriptions.Total)
	println("Total pages:", subscriptions.TotalPages)
	
	// Convert rows to subscription items
	if subscriptions.Rows != nil {
		rowsJSON, _ := json.Marshal(subscriptions.Rows)
		var items []tapsilat.SubscriptionListItem
		json.Unmarshal(rowsJSON, &items)
		
		for _, item := range items {
			println("Subscription:", item.Title, "- Amount:", item.Amount, "- Status:", item.PaymentStatus)
		}
	}
}
```

#### Cancel Subscription

```go
package main

import (
	"context"
	"github.com/tapsilat/tapsilat-go"
)

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	
	payload := tapsilat.SubscriptionCancelRequest{
		ReferenceID: "subscription_reference_id",
		// Or use ExternalReferenceID: "ext_sub_123",
	}
	
	err := api.CancelSubscription(context.Background(), payload)
	if err != nil {
		panic(err)
	}
	println("Subscription cancelled successfully!")
}
```

#### Redirect Subscription

```go
package main

import (
	"context"
	"github.com/tapsilat/tapsilat-go"
)

func main() {
	token := "TOKEN"
	api := tapsilat.NewAPI(token)
	
	payload := tapsilat.SubscriptionRedirectRequest{
		SubscriptionID: "subscription_id",
	}
	
	response, err := api.RedirectSubscription(context.Background(), payload)
	if err != nil {
		panic(err)
	}
	println("Redirect URL:", response.URL)
}
```

## API Methods

All API methods now require a `context.Context` as the first parameter for better control over request cancellation and timeouts.

### Order Operations

- `CreateOrder(ctx context.Context, order Order) (OrderResponse, error)`
- `GetOrder(ctx context.Context, referenceID string) (OrderDetail, error)`
- `GetOrderByConversationID(ctx context.Context, conversationID string) (OrderDetail, error)`
- `GetOrderStatus(ctx context.Context, referenceID string) (OrderStatus, error)`
- `GetOrders(ctx context.Context, page, perPage, buyerID string) (PaginatedData, error)`
- `GetOrderList(ctx context.Context, page, perPage int, startDate, endDate, organizationID, relatedReferenceID string) (PaginatedData, error)`
- `GetOrderSubmerchants(ctx context.Context, page, perPage int) (PaginatedData, error)`
- `GetCheckoutURL(ctx context.Context, referenceID string) (string, error)`

### Payment Operations

- `RefundOrder(ctx context.Context, refund RefundOrder) (RefundCancelOrderResponse, error)`
- `RefundAllOrder(ctx context.Context, referenceID string) (RefundCancelOrderResponse, error)`
- `CancelOrder(ctx context.Context, cancel CancelOrder) (RefundCancelOrderResponse, error)`

### Payment Terms Operations

- `CreateOrderTerm(ctx context.Context, term OrderPaymentTermCreateDTO) (map[string]interface{}, error)`
- `UpdateOrderTerm(ctx context.Context, term OrderPaymentTermUpdateDTO) (map[string]interface{}, error)`
- `GetOrderTerm(ctx context.Context, termReferenceID string) (map[string]interface{}, error)`
- `DeleteOrderTerm(ctx context.Context, orderID, termReferenceID string) (map[string]interface{}, error)`
- `RefundOrderTerm(ctx context.Context, term OrderTermRefundRequest) (map[string]interface{}, error)`

### Subscription Operations

- `CreateSubscription(ctx context.Context, subscription SubscriptionCreateRequest) (SubscriptionCreateResponse, error)`
- `GetSubscription(ctx context.Context, payload SubscriptionGetRequest) (SubscriptionDetail, error)`
- `ListSubscriptions(ctx context.Context, page, perPage int) (PaginatedData, error)`
- `CancelSubscription(ctx context.Context, payload SubscriptionCancelRequest) error`
- `RedirectSubscription(ctx context.Context, payload SubscriptionRedirectRequest) (SubscriptionRedirectResponse, error)`

### Utility Operations

- `GetOrderTransactions(ctx context.Context, referenceID string) (map[string]interface{}, error)`
- `GetOrderPaymentDetails(ctx context.Context, referenceID, conversationID string) (map[string]interface{}, error)`
- `OrderTerminate(ctx context.Context, referenceID string) (map[string]interface{}, error)`
- `OrderManualCallback(ctx context.Context, referenceID, conversationID string) (map[string]interface{}, error)`
- `OrderRelatedUpdate(ctx context.Context, referenceID, relatedReferenceID string) (map[string]interface{}, error)`
- `GetOrganizationSettings(ctx context.Context) (OrganizationSettings, error)`

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests (requires TAPSILAT_TOKEN)
make test-integration

# Run tests with coverage
make test-coverage

# Run specific test groups
make test-validators
make test-orders
make test-api
```

### Test Environment Setup

For integration tests, set your API token:

```bash
export TAPSILAT_TOKEN=your_test_token_here
make test-integration
```

### Development Setup

```bash
# Setup development environment
make dev-setup

# This will:
# 1. Download dependencies
# 2. Create .env file from .env.example
# 3. Set up the project for development
```

## Project Structure

```
tapsilat-go/
├── tapsilat.go          # Main API client
├── dtos.go              # Data transfer objects
├── validators.go        # Input validation functions
├── tests/
│   ├── unit/            # Unit tests
│   │   ├── validators_test.go
│   │   ├── order_test.go
│   │   └── api_test.go
│   └── integration/     # Integration tests
│       └── integration_test.go
├── Makefile             # Build and test commands
├── .env.example         # Environment variables template
└── README.md
```

## Error Handling

The SDK provides structured error handling:

```go
response, err := api.CreateOrder(context.Background(), order)
if err != nil {
    // Check if it's a validation error
    if validationErr, ok := err.(*tapsilat.ValidationError); ok {
        fmt.Printf("Validation Error: %s (Code: %d)\n",
            validationErr.Message, validationErr.Code)
        return
    }

    // Handle other errors
    fmt.Printf("API Error: %s\n", err.Error())
    return
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make test`)
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
