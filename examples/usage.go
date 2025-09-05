package main

import (
	"fmt"
	"log"

	"github.com/tapsilat/tapsilat-go"
)

func main() {
	// API client oluşturma - token direkt verilmeli
	token := "your_token_here"
	api := tapsilat.NewAPI(token)

	// Eğer token yoksa hata veriyor
	if api.Token == "" {
		log.Fatal("TAPSILAT_TOKEN environment variable is required")
	}

	// Scenario 1: Basit order oluşturma
	fmt.Println("=== Scenario 1: Basic Order ===")
	runBasicOrderExample(api)

	// Scenario 2: Basket items ile order
	fmt.Println("\n=== Scenario 2: Order with Basket Items ===")
	runOrderWithBasketItemsExample(api)

	// Scenario 3: Addresses ile order
	fmt.Println("\n=== Scenario 3: Order with Addresses ===")
	runOrderWithAddressesExample(api)

	// Scenario 4: Installments ve payment methods
	fmt.Println("\n=== Scenario 4: Order with Installments ===")
	runOrderWithInstallmentsExample(api)

	// Scenario 5: Checkout design
	fmt.Println("\n=== Scenario 5: Order with Checkout Design ===")
	runOrderWithCheckoutDesignExample(api)

	// Scenario 6: Validation examples
	fmt.Println("\n=== Scenario 6: Validation Examples ===")
	runValidationExamples()
}

func runBasicOrderExample(api *tapsilat.API) {
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

	response, err := api.CreateOrder(order)
	if err != nil {
		log.Printf("Error creating order: %v", err)
		return
	}

	fmt.Printf("Order created successfully!\n")
	fmt.Printf("Order ID: %s\n", response.OrderID)
	fmt.Printf("Reference ID: %s\n", response.ReferenceID)
	fmt.Printf("Checkout URL: %s\n", response.CheckoutURL)
}

func runOrderWithBasketItemsExample(api *tapsilat.API) {
	buyer := tapsilat.OrderBuyer{
		Name:    "John",
		Surname: "Doe",
		Email:   "test@example.com",
	}

	// Basket item payers
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

	response, err := api.CreateOrder(order)
	if err != nil {
		log.Printf("Error creating order with basket items: %v", err)
		return
	}

	fmt.Printf("Order with basket items created successfully!\n")
	fmt.Printf("Order ID: %s\n", response.OrderID)
	fmt.Printf("Reference ID: %s\n", response.ReferenceID)
	fmt.Printf("Checkout URL: %s\n", response.CheckoutURL)
}

func runOrderWithAddressesExample(api *tapsilat.API) {
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

	response, err := api.CreateOrder(order)
	if err != nil {
		log.Printf("Error creating order with addresses: %v", err)
		return
	}

	fmt.Printf("Order with addresses created successfully!\n")
	fmt.Printf("Order ID: %s\n", response.OrderID)
	fmt.Printf("Reference ID: %s\n", response.ReferenceID)
	fmt.Printf("Checkout URL: %s\n", response.CheckoutURL)
}

func runOrderWithInstallmentsExample(api *tapsilat.API) {
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
		PaymentSuccessUrl:   "https://example.com/success",
		PaymentFailureUrl:   "https://example.com/failure",
	}

	response, err := api.CreateOrder(order)
	if err != nil {
		log.Printf("Error creating order with installments: %v", err)
		return
	}

	fmt.Printf("Order with installments created successfully!\n")
	fmt.Printf("Order ID: %s\n", response.OrderID)
	fmt.Printf("Reference ID: %s\n", response.ReferenceID)
	fmt.Printf("Checkout URL: %s\n", response.CheckoutURL)
}

func runOrderWithCheckoutDesignExample(api *tapsilat.API) {
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

	response, err := api.CreateOrder(order)
	if err != nil {
		log.Printf("Error creating order with checkout design: %v", err)
		return
	}

	fmt.Printf("Order with checkout design created successfully!\n")
	fmt.Printf("Order ID: %s\n", response.OrderID)
	fmt.Printf("Reference ID: %s\n", response.ReferenceID)
	fmt.Printf("Checkout URL: %s\n", response.CheckoutURL)
}

func runValidationExamples() {
	// GSM Number validation examples
	fmt.Println("GSM Number Validation Examples:")

	gsmNumbers := []string{
		"+905551234567",
		"00905551234567",
		"05551234567",
		"5551234567",
		"+90 555 123-45(67)",
		"invalid_phone",
		"+90123", // too short
	}

	for _, gsm := range gsmNumbers {
		cleaned, err := tapsilat.ValidateGSMNumber(gsm)
		if err != nil {
			fmt.Printf("  ❌ %s -> Error: %s\n", gsm, err.Error())
		} else {
			fmt.Printf("  ✅ %s -> %s\n", gsm, cleaned)
		}
	}

	// Installments validation examples
	fmt.Println("\nInstallments Validation Examples:")

	installmentStrings := []string{
		"1,2,3,6",
		"1, 2, 3, 6",
		"2,4,8,12",
		"",
		"1,15,3",  // invalid
		"1,abc,3", // invalid
	}

	for _, installmentStr := range installmentStrings {
		installments, err := tapsilat.ValidateInstallments(installmentStr)
		if err != nil {
			fmt.Printf("  ❌ %s -> Error: %s\n", installmentStr, err.Error())
		} else {
			fmt.Printf("  ✅ %s -> %v\n", installmentStr, installments)
		}
	}
}
