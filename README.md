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

## Local End-to-End Validation (Panel + SDK)

Use this flow to validate newly added submerchant/vpos-related SDK APIs against local `panel/backend`.

- Start local stack from `panel/backend`:

```bash
make compose
```

- Open `http://localhost:8080`, login, and generate API key/secret.

- In `tapsilat-go`, run local E2E script:

```bash
chmod +x scripts/local_e2e.sh
TAPSILAT_API_KEY="<ui_api_key>" \
TAPSILAT_API_SECRET="<ui_api_secret>" \
TAPSILAT_IT_SUBMERCHANT_ID="<submerchant_id>" \
TAPSILAT_IT_SUBORGANIZATION_ID="<suborganization_id>" \
scripts/local_e2e.sh
```

Optional:

- Reuse an existing token: set `TAPSILAT_API_TOKEN`.
- Start stack from script: `scripts/local_e2e.sh --start-stack`.
- Run only smoke: `scripts/local_e2e.sh --smoke-only`.
- Run only integration: `scripts/local_e2e.sh --integration-only`.
- Auto-create submerchant + resolve suborganization IDs before tests:

```bash
TAPSILAT_API_KEY="<ui_api_key>" \
TAPSILAT_API_SECRET="<ui_api_secret>" \
scripts/local_e2e.sh --bootstrap-submerchant
```

- Auto-create VPOS + attach it to submerchant via `vpos-submerchant` mapping:

```bash
TAPSILAT_API_KEY="<ui_api_key>" \
TAPSILAT_API_SECRET="<ui_api_secret>" \
scripts/local_e2e.sh --bootstrap-submerchant --bootstrap-vpos
```

This writes IDs to `/tmp/tapsilat_local_e2e_ids.env` (override with `BOOTSTRAP_OUTPUT_FILE`).

- `BOOTSTRAP_CREATED_SUBMERCHANT_ID` / `BOOTSTRAP_CREATED_SUBORGANIZATION_ID`: raw IDs created/discovered during bootstrap.
- `BOOTSTRAP_CREATED_VPOS_ID` / `BOOTSTRAP_CREATED_VPOS_SUBMERCHANT_ID`: raw VPOS and mapping IDs created during VPOS bootstrap.
- `TAPSILAT_SMOKE_*` and `TAPSILAT_IT_*`: effective IDs used for tests.

If reverse mapping is not yet consistent in backend, script continues with submerchant-only smoke assertions and leaves integration IDs empty instead of failing immediately.

Notes:

- Script sets `User-Agent: Go-http-client/1.1` while generating token to match SDK request context and avoid local auth mismatch.
- Integration test requires `TAPSILAT_IT_SUBMERCHANT_ID` and `TAPSILAT_IT_SUBORGANIZATION_ID`; otherwise test code skips it.

## Marketplace payout and policy APIs

Marketplace methods use the version 2 Panel routes even when the client endpoint ends in `/api/v1`.
They include exact seller lookup by `external_id`, payout list/read/bulk approval, organization
approval settings, and approval/release policy create, update and delete operations. Delete methods
accept an empty successful response body.

Bulk payout approval accepts at most 100 IDs. Callers should approve only payouts whose
`eligibility_status` is `eligible` and whose `release_status` is `awaiting_approval`; the server
revalidates both conditions.

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

### Marketplace Payout Release

Set an optional release time on each order item and report a business event when the service is completed:

```go
releaseAt := time.Date(2026, time.December, 15, 10, 0, 0, 0, time.UTC)
order.BasketItems[0].PayoutReleaseAt = &releaseAt

event, err := api.RecordSubmerchantPayoutEvent(ctx, tapsilat.SubmerchantPayoutEventRequest{
    IdempotencyKey: "booking-123-service-completed",
    OrderReference: orderReference,
    ItemID:          order.BasketItems[0].Id,
    EventType:       "service_completed",
    OccurredAt:      &releaseAt,
    Payload:         map[string]any{"booking_id": "booking-123"},
})
```

Create a logical marketplace seller and provision explicitly selected VPOS accounts:

```go
seller, err := api.CreateMarketplaceSubmerchant(ctx, tapsilat.MarketplaceSubmerchantCreateRequest{
    IdempotencyKey:    "seller-create-123",
    Locale:            "tr",
    ConversationID:    "seller-create-123",
    Name:              "Example Travel",
    Email:             "finance@example.test",
    GsmNumber:         "905551112233",
    CurrencyID:        "TRY",
    SubmerchantType:   "PRIVATE_COMPANY",
    TaxNumber:         "1234567890",
    ProviderAccounts: &[]tapsilat.MarketplaceAccountTarget{
        {VposID: marketplaceVposID, Environment: "TEST"},
    },
})
```

Use `seller.RoutingReference` as the basket item's `SubMerchantKey`. This is an opaque, provider-neutral seller reference.

After checking `err`, inspect `seller.Provisionings`: each entry identifies a VPOS account and reports `status`, `retryable`, and `error_message`. A successful API response confirms logical seller creation, not that every account is ready. Partial failures preserve completed mappings; recovery handles eligible unfinished operations. Only accounts with completed mappings can route this seller's payments. An absent or empty list from an older server does not confirm readiness. Reuse the same `IdempotencyKey` and request data when retrying creation.

Omitting both account selectors, or providing an empty `ProviderAccounts` slice, creates only the logical seller. Legacy `VposID` selects exactly that account and uses its configured environment; it never registers other accounts. Do not combine `VposID` and `ProviderAccounts`. Explicit selections must belong to the token's organization and share one environment. Use `GetMarketplaceSellerContracts` to discover account requirements before submitting a profile.

### Marketplace lifecycle API

All marketplace operations use the organization and actor from the API token. Callers cannot override tenant scope in the request.

| Operation | SDK methods |
| --- | --- |
| Seller directory and profile | `ListMarketplaceSubmerchants`, `GetMarketplaceSubmerchant`, `UpdateMarketplaceSubmerchant` |
| Account contracts and state | `GetMarketplaceSellerContracts`, `GetMarketplaceSellerAccounts` |
| Account lifecycle | `RunMarketplaceSellerAccountAction` |
| Payout tracking and approval | `ListMarketplacePayouts`, `GetMarketplacePayout`, `ApproveMarketplacePayouts` |
| Policy inspection | `ListMarketplaceApprovalPolicies`, `GetMarketplaceApprovalPolicy`, `ListMarketplaceReleasePolicies`, `GetMarketplaceReleasePolicy` |
| Item business events | `RecordSubmerchantPayoutEvent` |

Profile updates replace the profile and require the latest `Revision` returned by `GetMarketplaceSubmerchant`. Preserve fields you are not changing and omit account selectors. A stale revision is rejected. Inspect `ApprovalRequired`: when true, an administrator must review the proposed change. Account synchronization remains separately observable in `GetMarketplaceSellerAccounts`.

For account operations, read the current row and choose an operation from its `Actions`. Send its exact `Revision` and `Environment`; never automatically select a production account. `retry` uses the original registration snapshot, while `resubmit` uses the current saved profile. The available operations are capability-dependent and can include `create`, `retry`, `resubmit`, `update`, `retrieve`, and `retry_sync`. Inspect the returned provisioning `Status` even when the HTTP request succeeds.

Payout records are created by successful marketplace payment processing. There is no standalone create-payout API. A `service_completed` event applies only to the referenced order item and satisfies its business-event gate. It does not override a future release date, minimum delay, provider restriction, or manual approval requirement. Repeating an event with the same idempotency key and payload is safe; a conflicting payload is rejected.

After reporting the event, inspect the payout before requesting approval:

```go
payout, err := api.GetMarketplacePayout(ctx, payoutID)
if err != nil {
    return err
}
if payout.EligibilityStatus != "eligible" || payout.ReleaseStatus != "awaiting_approval" {
    return fmt.Errorf("payout is not awaiting eligible approval: %s / %s", payout.EligibilityStatus, payout.ReleaseStatus)
}
accepted, err := api.ApproveMarketplacePayouts(ctx, tapsilat.MarketplacePayoutApproval{
    IDs: []string{payout.ID},
})
if err != nil {
    return err
}
_ = accepted // "accepted" acknowledges local approval; poll GetMarketplacePayout for settlement.
```

Approval supports up to 100 payout IDs. All release gates are checked again by Admin. A repeated approval of an already released payout is rejected rather than creating another payout. Eligibility overrides remain Admin-only.

The legacy `ApproveSubmerchantPayment` endpoint now resolves an existing organization-scoped payout from the provider transaction reference and enters the same approval workflow. Unknown or ambiguous references are rejected. Prefer payout IDs, especially when an organization has multiple accounts. Legacy profile updates enforce marketplace immutable-field checks, but cannot express the caller's revision or review outcome; prefer `UpdateMarketplaceSubmerchant`.

Legacy `DisapproveSubmerchantPayment` and `UpdateSubmerchantPaymentItem` resolve the stored, organization-scoped payout and use the same durable operation service as the new API. They select the payout's exact VPOS account and translate the seller reference through that account's mapping. Responses include `OperationID`; inspect `Status`, since `unknown` or `processing` does not mean success. Use a new `ConversationID` for each logical change. Repeating it with the same payload returns the original receipt; changing its payload is rejected. Without a conversation ID, identical requests share a deterministic key, so a later separate change requires a new conversation ID.

### Payout allocation changes and approval withdrawal

Read the payout's current revision, then submit `disapprove` or `update_item`:

```go
current, err := api.GetMarketplacePayout(ctx, payoutID)
if err != nil {
    return err
}
operation, err := api.RunMarketplacePayoutOperation(ctx, current.ID, tapsilat.MarketplacePayoutOperationInput{
    Kind:            "update_item",
    IdempotencyKey:  "booking-123-allocation-change-1",
    Revision:        current.Revision,
    SellerReference: seller.RoutingReference,
    NetAmount:       "60.25",
})
if err != nil {
    return err
}
_ = operation // Inspect Status, and retain ID for verification or explicit retry.
```

For `disapprove`, omit `SellerReference` and `NetAmount`. `update_item` requires a registered seller on the payout's existing account, matching currency, and a nonnegative amount with at most two decimal places within the captured item amount. Refunded/cancelled items and allocations with applied offsets cannot be reassigned. The original paid basket remains a historical snapshot; use the payout API for its current seller and allocation.

Both operations preserve the payment-time release policy, dates and required business event. They require a fresh manual approval afterward. Operation `succeeded` confirms the local/provider change, not seller settlement.

| Provider | Supported change window |
| --- | --- |
| iyzico | Update an unapproved paid item; withdraw a provider approval through the disapprove service. Withdraw approval before changing an already approved allocation. |
| PayTR | Change the local transfer instruction before submission. Seller/beneficiary and refund allocation bookkeeping update in one transaction. Processing, submitted, settled and indeterminate transfers cannot be recalled through these operations. |

The service persists an operation receipt before remote dispatch and sets `ActiveOperationID` and `ReleaseStatus = "processing"`. Workers, approvals and eligibility overrides cannot release that payout during an unresolved operation. Identical idempotent replays never redispatch the mutation. A known provider rejection leaves the original allocation intact and requires manual approval; submit a corrected operation with a new key and fresh payout revision.

Use `GetMarketplacePayoutOperation` to read a receipt. For a timeout or lost response, `ReconcileMarketplacePayoutOperation` queries the provider and completes local bookkeeping only when the requested result is confirmed. `RetryMarketplacePayoutOperation` explicitly retries only after verifying that the provider still has the previous allocation/approval. An operation still within its one-minute in-flight window is not redispatched. Unverifiable results remain held. Changing account configuration during an unresolved operation requires restoring/verifying that account before recovery.

The provider constraints are described in the official [iyzico approval documentation](https://docs.iyzico.com/en/products/marketplace/marketplace-implementation/approval), [iyzico item update documentation](https://docs.iyzico.com/en/products/marketplace/marketplace-implementation/submerchant/alt-uye-sorgulama-1) and [PayTR transfer documentation](https://dev.paytr.com/platform-transfer-talebi/transfer-talimatinin-verilmesi).

### Local marketplace stack verification

The opt-in test starts isolated Admin gRPC and Panel HTTP servers and calls them through this SDK; Monolog persistence is disabled by the test environment:

```sh
TAPSILAT_STACK_ADMIN_ROOT=/absolute/path/to/admin/backend \
TAPSILAT_STACK_PANEL_ROOT=/absolute/path/to/panel/backend \
CONFIG_PATH=/absolute/path/to/local/admin/config.yaml \
go test -v ./tests/integration -run '^TestMarketplaceAdminPanelSDKStack$' -count=1
```

Use matching feature-stack revisions in all three repositories. The test uses an in-memory database, a temporary token verifier, and a local gateway contract fixture. It exercises real SDK HTTP, Panel routing, Admin gRPC handlers and domain services. It covers explicit account selection, idempotency, profile revisions, provider contracts/actions, payout materialization for both providers from seeded paid orders, item events, eligibility, approval, policy reads, legacy disapprove/item updates, versioned payout operations and lost-response reconciliation. It does not execute an order-create request, a real JWT login, remote provider registration, a card payment or provider settlement. Sandbox payment verification is a separate test.

Apply the Admin migration for `marketplace_payout_operations` and `marketplace_payouts.active_operation_id`, deploy matching Admin handlers/workers, then Panel, and verify both before publishing a new SDK release. These methods require the matching backend stack; an older published SDK or server does not provide this API parity.

### Order with Payment Terms (Installments)

```go
amount1 := 50.0
required := true
sequence := 1

paymentTerm := tapsilat.OrderPaymentTerm{
    Amount:          &amount1,
    DueDate:         "2024-01-15T00:00:00Z", // RFC3339; required on every term
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

Payment terms on an existing order use their own endpoints. `OrderID` is the order's ID (`OrderResponse.OrderID`), not its reference ID, and term dates use `YYYY-MM-DD`:

```go
sequence := 2
_, err := api.CreateOrderTerm(ctx, tapsilat.OrderPaymentTermCreateDTO{
    OrderID:      orderResponse.OrderID,
    TermSequence: &sequence,
    Required:     true,
    Amount:       25,
    DueDate:      "2026-02-15",
})

// Refunds need the order reference and the term's ID (OrderPaymentTermDTO.ID).
_, err = api.RefundOrderTerm(ctx, tapsilat.OrderTermRefundRequest{
    ReferenceID: orderResponse.ReferenceID,
    TermID:      term.ID,
    Amount:      10,
})
```

`UpdateOrderTerm` replaces the term, so send every field rather than only the changed ones.

### Card Installments and Charging

```go
options, err := api.GetCardInstallments(ctx, tapsilat.CardInstallmentsRequest{
    Bin:      "54066200",
    Amount:   100,
    Currency: "TRY",
})
// options.Installments lists each installment count with its commission rate and monthly payment.

table, err := api.GetInstallmentTable(ctx, tapsilat.InstallmentTableRequest{Amount: 100, Currency: "TRY"})
// table.Groups lists installment rows per bank card program.

charge, err := api.ChargeOrder(ctx, tapsilat.OrderChargeRequest{OrderReferenceID: referenceID})
if err == nil && !charge.IsSuccess {
    // The request succeeded but the charge did not; inspect charge.ErrorCode and charge.Message.
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

### Car Loan Limit and Payment Plan

```go
limit, err := api.CarLoanCheckLimit(ctx, tapsilat.CarLoanCheckLimitRequest{
	ReferenceID:     "vehicle-ref",
	Tckn:            "11111111111",
	PhoneNumber:     "5551112233",
	Dob:             "1990-01-01",
	RequestedAmount: 400000,
	SellerID:        "dealer-seller-code",
	EmployeeID:      "dealer-employee-code",
	Vehicle:         &tapsilat.CarLoanVehicle{Model: "ADRIA CORAL XL 600 DP", Year: 2024, Plate: "34ABC123"},
})
if err != nil {
	return err
}

offer := limit.Offers[0]
plan, err := api.CarLoanPaymentPlan(ctx, tapsilat.CarLoanPaymentPlanRequest{
	CredentialID: offer.CredentialID,
	OfferID:      offer.ID,
	Term:         36,
})
```

### Subscription Operations

#### Create Subscription

Use `Interval` with one of `day`, `week`, `month`, or `year`, together with a
positive `IntervalCount`, to define the recurring cadence. `Cycle: 0` renews the
subscription until it is canceled; a positive value limits the total number of
settled billing cycles.

Set `CardID` to an existing scoped saved card. For supported hosted mandate
setups, leave `CardID` empty and it will be omitted from the request. Applications
do not need to collect or send raw card data through this SDK.

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
		Interval:            "month",
		IntervalCount:       1,
		Cycle:               0,
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

### Webhook Handling

#### Verify Webhook Signature

Tapsilat signs every callback it sends to your callback URLs, and to an order's `payment_success_url` and `payment_failure_url`, with the webhook secret from the panel's callback settings. The `Tapsilat-Signature` header looks like `t=1726480000,v1=5257a8…`: `v1` is the hex HMAC-SHA256 of the timestamp, a dot and the raw body.

`VerifyWebhook` checks it in constant time and rejects a timestamp more than 5 minutes from the server clock, so a captured callback cannot be replayed. Pass the raw body, not re-encoded JSON.

```go
package main

import (
	"io"
	"log"
	"net/http"

	"github.com/tapsilat/tapsilat-go"
)

func main() {
	secret := "your_webhook_secret"

	http.HandleFunc("/tapsilat/callback", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil || !tapsilat.VerifyWebhook(body, r.Header.Get(tapsilat.WebhookSignatureHeader), secret) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		event, err := tapsilat.ParseWebhookEvent(body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// event.Type is "order_lite" (the default), "order", "order_payment_detail", ...
		// ReferenceID reads the order's reference ID whatever the callback type; Status is
		// e.g. "Paid" when the callback type carries a status.
		log.Printf("%s callback for order %s (conversation %s): %s",
			event.Type, event.ReferenceID(), event.ConversationID(), event.Status())

		if event.IsPaid() {
			// fulfil the order
		}
		w.WriteHeader(http.StatusOK)
	})
	_ = http.ListenAndServe(":8080", nil)
}
```

`VerifyWebhookAt(payload, signature, secret, tolerance, now)` takes a different tolerance (`0` disables the timestamp check) and the time to compare with. A bare hex HMAC of the body, with or without a `sha256=` prefix, is still accepted.

#### Order Statuses

The `OrderStatus*` constants are the numeric status codes that callbacks carry, and `OrderStatusName` returns the names `GetOrderStatus` returns. The helpers accept a code, a numeric string or a name (case-insensitive):

```go
status, err := api.GetOrderStatus(context.Background(), referenceID) // status.Status is e.g. "Waiting for payment"

tapsilat.IsPaidOrderStatus(status.Status)     // Paid, Completed, Sub merchant payment approved
tapsilat.IsFailedOrderStatus(status.Status)   // Cancelled, Fraud, Rejected, Failure, Expired, Terminated
tapsilat.IsRefundedOrderStatus(status.Status) // Refunded, Partially refunded
tapsilat.OrderStatusCode(status.Status)       // 7; also accepts "7" or 7
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
- `GetOrderListWithFilter(ctx context.Context, filter OrderListFilter) (PaginatedData, error)`
- `GetOrderSubmerchants(ctx context.Context, page, perPage int) (OrderSubmerchantListResponse, error)`
- `GetCheckoutURL(ctx context.Context, referenceID string) (string, error)`

### Payment Operations

- `RefundOrder(ctx context.Context, refund RefundOrder) (RefundCancelOrderResponse, error)`
- `RefundAllOrder(ctx context.Context, referenceID string) (RefundCancelOrderResponse, error)`
- `CancelOrder(ctx context.Context, cancel CancelOrder) (RefundCancelOrderResponse, error)`
- `GetOrderPayments(ctx context.Context, payload GetOrderPaymentsRequest) (GetOrderPaymentsResponse, error)`
- `ChargeOrder(ctx context.Context, payload OrderChargeRequest) (OrderChargeResponse, error)`
- `GetCardInstallments(ctx context.Context, req CardInstallmentsRequest) (CardInstallmentsResponse, error)`
- `GetInstallmentTable(ctx context.Context, req InstallmentTableRequest) (InstallmentTableResponse, error)`
- `ApproveSubmerchantPayment(ctx context.Context, payload SubmerchantPaymentAction) (SubmerchantPaymentActionResponse, error)`
- `DisapproveSubmerchantPayment(ctx context.Context, payload SubmerchantPaymentAction) (SubmerchantPaymentActionResponse, error)`
- `UpdateSubmerchantPaymentItem(ctx context.Context, payload SubmerchantPaymentItemUpdate) (SubmerchantPaymentItemUpdateResponse, error)`

### Payment Terms Operations

- `CreateOrderTerm(ctx context.Context, term OrderPaymentTermCreateDTO) (map[string]any, error)`
- `UpdateOrderTerm(ctx context.Context, term OrderPaymentTermUpdateDTO) (map[string]any, error)`
- `GetOrderTerm(ctx context.Context, termReferenceID string) (map[string]any, error)`
- `DeleteOrderTerm(ctx context.Context, orderID, termReferenceID string) (map[string]any, error)`
- `RefundOrderTerm(ctx context.Context, term OrderTermRefundRequest) (map[string]any, error)`

### Subscription Operations

- `CreateSubscription(ctx context.Context, subscription SubscriptionCreateRequest) (SubscriptionCreateResponse, error)`
- `GetSubscription(ctx context.Context, payload SubscriptionGetRequest) (SubscriptionDetail, error)`
- `ListSubscriptions(ctx context.Context, page, perPage int) (PaginatedData, error)`
- `CancelSubscription(ctx context.Context, payload SubscriptionCancelRequest) error`
- `RedirectSubscription(ctx context.Context, payload SubscriptionRedirectRequest) (SubscriptionRedirectResponse, error)`

### Car Loan Operations

- `CarLoanProviders(ctx context.Context) (CarLoanProvidersResponse, error)`
- `CarLoanCheckLimit(ctx context.Context, payload CarLoanCheckLimitRequest) (CarLoanCheckLimitResponse, error)`
- `CarLoanPaymentPlan(ctx context.Context, payload CarLoanPaymentPlanRequest) (CarLoanPaymentPlanResponse, error)`
- `CarLoanApply(ctx context.Context, payload CarLoanApplyRequest) (CarLoanApplyResponse, error)`
- `CarLoanStatus(ctx context.Context, payload CarLoanStatusRequest) (CarLoanStatusResponse, error)`

These call the panel's `/loan/car/*` routes. The panel runs them on the mono-loan instance the token's organization is set up with (its Mono / Car Loan provider in admin), so quote, payment plan and application always reach the same instance. The organization comes from the token; it is never sent in the request.

### Utility Operations

- `GetOrderTransactions(ctx context.Context, referenceID string) ([]map[string]any, error)`
- `GetOrderPaymentDetails(ctx context.Context, referenceID string) (map[string]any, error)`
- `GetOrderPaymentDetailsByConversationID(ctx context.Context, conversationID string) (map[string]any, error)`
- `OrderTerminate(ctx context.Context, referenceID string) (map[string]any, error)`
- `OrderManualCallback(ctx context.Context, referenceID, conversationID string) (map[string]any, error)`
- `OrderRelatedUpdate(ctx context.Context, referenceID, relatedReferenceID string) (map[string]any, error)`
- `GetOrganizationCurrencies(ctx context.Context) (OrganizationCurrenciesResponse, error)`
- `GetOrganizationSettings(ctx context.Context) (OrganizationSettings, error)`
- `CreateOrganizationUser(ctx context.Context, payload OrgCreateUserRequest) (OrgCreateUserResponse, error)`
- `CreateOrganizationUserToken(ctx context.Context, payload OrgUserTokenCreateRequest) (OrgUserTokenCreateResponse, error)`

### Webhook Helpers

- `VerifyWebhook(payload []byte, signature, secret string) bool`
- `VerifyWebhookAt(payload []byte, signature, secret string, tolerance time.Duration, now time.Time) bool`
- `ParseWebhookEvent(payload []byte) (WebhookEvent, error)` and `NewWebhookEvent(payload map[string]any) WebhookEvent`
- `WebhookEvent.ReferenceID()`, `ConversationID()`, `StatusCode()`, `Status()`, `IsPaid()`, `IsFailed()`, `IsRefunded()`
- `OrderStatusCode(status any) int`, `OrderStatusName(code int) string`
- `IsPaidOrderStatus`, `IsFailedOrderStatus`, `IsRefundedOrderStatus` with `PaidOrderStatuses`, `FailedOrderStatuses`, `RefundedOrderStatuses`

### Management Operations

- `CreateSubmerchant(ctx context.Context, payload SubmerchantCreateRequest) (SubmerchantMutationResponse, error)`
- `CreateMarketplaceSubmerchant(ctx context.Context, payload MarketplaceSubmerchantCreateRequest) (MarketplaceSubmerchantCreateResponse, error)`
- `RecordSubmerchantPayoutEvent(ctx context.Context, payload SubmerchantPayoutEventRequest) (SubmerchantPayoutEventResponse, error)`
- `GetSubmerchant(ctx context.Context, id string) (Submerchant, error)`
- `ListSubmerchants(ctx context.Context, page, perPage int) (SubmerchantListResponse, error)`
- `UpdateSubmerchant(ctx context.Context, id string, payload SubmerchantUpdateRequest) (SubmerchantMutationResponse, error)`
- `DeleteSubmerchant(ctx context.Context, id string) (SubmerchantMutationResponse, error)`
- `GetSuborganizations(ctx context.Context, page, perPage int) (SuborganizationListResponse, error)`
- `GetSuborganization(ctx context.Context, id string) (SuborganizationListItem, error)`
- `GetSuborganizationDetail(ctx context.Context, id string) (SuborganizationDetail, error)`
- `GetSuborganizationBySubmerchant(ctx context.Context, submerchantID string) (SubmerchantSuborganizationMapping, error)`
- `GetSubmerchantBySuborganization(ctx context.Context, suborganizationID string) (SuborganizationSubmerchantMapping, error)`
- `ListVpos(ctx context.Context, page, perPage int) (VposListResponse, error)`
- `ListVposWithFilter(ctx context.Context, page, perPage int, filter VposListFilter) (VposListResponse, error)`
- `CreateVpos(ctx context.Context, payload VposCreateRequest) (VposMutationResponse, error)`
- `GetVpos(ctx context.Context, id string) (Vpos, error)`
- `UpdateVpos(ctx context.Context, id string, payload VposUpdateRequest) (VposMutationResponse, error)`
- `DeleteVpos(ctx context.Context, id string) (VposMutationResponse, error)`
- `ListVposAcquirers(ctx context.Context) (VposAcquirerListResponse, error)`
- `ListCardSchemes(ctx context.Context) (CardSchemeListResponse, error)`
- `ListVposSubmerchants(ctx context.Context, page, perPage int, vposID, externalReferenceID string) (VposSubmerchantListResponse, error)`
- `CreateVposSubmerchant(ctx context.Context, payload VposSubmerchantCreateRequest) (VposSubmerchantMutationResponse, error)`
- `GetVposSubmerchant(ctx context.Context, id string) (VposSubmerchant, error)`
- `UpdateVposSubmerchant(ctx context.Context, id string, payload VposSubmerchantUpdateRequest) (VposSubmerchantMutationResponse, error)`
- `DeleteVposSubmerchant(ctx context.Context, id string) (VposSubmerchantMutationResponse, error)`

Currency fields on submerchants (v1 and marketplace), VPOS `Currencies` and seller contracts accept either currency UUIDs or organization `currency_unit` values such as `TRY`/`USD`; the server resolves units, so the SDK sends them unchanged. Marketplace approval policies are the exception: their `CurrencyID` must be a UUID on the wire, so `CreateMarketplaceApprovalPolicy` and `UpdateMarketplaceApprovalPolicy` resolve a unit first. `ResolveCurrencyID` exposes the same lookup.

`UpdateVpos` and `UpdateSubmerchant` replace the stored configuration. Read the current record first and send every field back; for example, omitting `Failover` or `TLS` stores `false`.

`OrderCheckoutDesign.OrderDetailHtml` and `RedirectUrl` are deprecated and ignored by checkout. Use `Order.RedirectSuccessUrl` and `Order.RedirectFailureUrl` for redirects.

## Testing

### Running Tests

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run only integration tests (requires TAPSILAT_TOKEN)
make test-integration

# Run smoke tests against local/sandbox panel
make test-smoke

# Run tests with coverage
make test-coverage

# Run specific test groups
make test-validators
make test-orders
make test-api
```

### Test Environment Setup

For integration tests, set backend verification environment variables:

```bash
export TAPSILAT_IT_ENDPOINT=http://localhost:3001/api/v1
export TAPSILAT_IT_TOKEN=your_token_here
export TAPSILAT_IT_SUBMERCHANT_ID=your_submerchant_id
export TAPSILAT_IT_SUBORGANIZATION_ID=your_suborganization_id
export TAPSILAT_IT_VPOS_ID=optional_vpos_id
go test -v ./tests/integration -run TestBackendConsistency_SubmerchantSuborganizationAndScopedVpos
```

This integration test validates real backend behavior via SDK calls:

- submerchant <-> suborganization mapping consistency (both directions)
- suborganization-scoped VPOS listing consistency
- optional `vpos_id` read + `ListVposSubmerchants` scope validation

For smoke tests, set endpoint/token (and optionally sample ids):

```bash
export TAPSILAT_SMOKE_ENDPOINT=http://localhost:3001/api/v1
export TAPSILAT_SMOKE_TOKEN=your_token_here
export TAPSILAT_SMOKE_SUBMERCHANT_ID=
export TAPSILAT_SMOKE_VPOS_ID=
export TAPSILAT_SMOKE_SUBORGANIZATION_ID=
make test-smoke
```

When optional IDs are provided, smoke coverage also verifies:

- submerchant <-> suborganization mapping reads
- scoped `ListVposWithFilter(..., VposListFilter{SuborganizationID: ...})`
- `ListVposSubmerchants` filtered by `vpos_id`

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

```text
tapsilat-go/
├── tapsilat.go          # Main API client
├── dtos.go              # Data transfer objects
├── validators.go        # Input validation functions
├── webhook.go           # Callback signature verification and parsing
├── tests/
│   ├── unit/            # Unit tests
│   │   ├── validators_test.go
│   │   ├── order_test.go
│   │   └── api_test.go
│   └── integration/     # Integration tests
│       └── integration_test.go
│   └── smoke/           # Local/sandbox smoke tests
│       └── smoke_test.go
├── Makefile             # Build and test commands
├── .env.example         # Environment variables template
└── README.md
```

## Error Handling

The SDK provides structured error handling:

```go
response, err := api.CreateOrder(context.Background(), order)
if err != nil {
	var validationErr *tapsilat.ValidationError
	if errors.As(err, &validationErr) {
		fmt.Printf("Validation Error: %s (Code: %d)\n", validationErr.Message, validationErr.Code)
		return
	}

	var apiErr *tapsilat.APIError
	if errors.As(err, &apiErr) {
		fmt.Printf("API Error: status=%d code=%s message=%s raw=%s\n",
			apiErr.StatusCode,
			apiErr.Code,
			apiErr.Message,
			apiErr.RawBody,
		)
		return
	}

	fmt.Printf("Unexpected Error: %s\n", err.Error())
	return
}
```

`APIError` normalizes HTTP/API failures into these fields:

- `StatusCode`: HTTP status code
- `Status`: API status field if present, otherwise HTTP status text
- `Code`: API error code if present
- `Message`: API `message` or `error` field if present
- `Errors`: per-field validation messages, when the server returns an `errors` object
- `RawBody`: original response body for fallback debugging

Most "not found" errors arrive as HTTP 400, and rejected state transitions (for example approving an ineligible payout) as HTTP 424. A successful status with an HTML body, which a deployment serving the panel UI returns for unknown routes, is reported as an `APIError` whose `Message` says the response was not JSON.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make test`)
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
