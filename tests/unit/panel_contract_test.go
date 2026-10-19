package unit_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

type recordedRequest struct {
	Method string
	Path   string
	Query  string
	Body   string
}

// newContractServer answers every request with status and body and records what it received.
func newContractServer(t *testing.T, status int, contentType, body string) (*httptest.Server, *recordedRequest) {
	t.Helper()
	recorded := &recordedRequest{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		recorded.Method = r.Method
		recorded.Path = r.URL.Path
		recorded.Query = r.URL.RawQuery
		recorded.Body = string(data)
		if contentType != "" {
			w.Header().Set("Content-Type", contentType)
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	return server, recorded
}

func TestOrderTermRoutesMatchPanel(t *testing.T) {
	ctx := context.Background()

	t.Run("GetOrderTerm", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"term_sequence":1,"amount":50,"due_date":{"seconds":1767225600},"status":"pending"}`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token").GetOrderTerm(ctx, "term ref/1")
		require.NoError(t, err)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "/order/term", req.Path)
		assert.Equal(t, "term_reference_id=term+ref%2F1", req.Query)
		assert.Equal(t, "pending", res["status"])
	})

	t.Run("CreateOrderTerm", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"message":"ok","code":200}`)
		sequence := 2
		_, err := tapsilat.NewCustomAPI(server.URL, "token").CreateOrderTerm(ctx, tapsilat.OrderPaymentTermCreateDTO{
			OrderID: "order_1", TermSequence: &sequence, Required: true, Amount: 25, DueDate: "2026-01-15",
		})
		require.NoError(t, err)
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, "/order/term", req.Path)
		assert.JSONEq(t, `{"order_id":"order_1","term_sequence":2,"required":true,"amount":25,"due_date":"2026-01-15"}`, req.Body)
	})

	t.Run("UpdateOrderTerm", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"message":"ok","code":200}`)
		amount := 30.0
		_, err := tapsilat.NewCustomAPI(server.URL, "token").UpdateOrderTerm(ctx, tapsilat.OrderPaymentTermUpdateDTO{
			TermReferenceID: "term_1", Amount: &amount, DueDate: "2026-02-01", Status: "pending",
		})
		require.NoError(t, err)
		assert.Equal(t, http.MethodPatch, req.Method)
		assert.Equal(t, "/order/term", req.Path)
		assert.JSONEq(t, `{"term_reference_id":"term_1","amount":30,"due_date":"2026-02-01","status":"pending"}`, req.Body)
	})

	t.Run("DeleteOrderTermSendsBody", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"message":"deleted","code":200}`)
		_, err := tapsilat.NewCustomAPI(server.URL, "token").DeleteOrderTerm(ctx, "order_1", "term_1")
		require.NoError(t, err)
		assert.Equal(t, http.MethodDelete, req.Method)
		assert.Equal(t, "/order/term", req.Path)
		assert.JSONEq(t, `{"order_id":"order_1","term_reference_id":"term_1"}`, req.Body)
	})

	t.Run("RefundOrderTerm", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"message":"ok","code":200,"is_success":true}`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token").RefundOrderTerm(ctx, tapsilat.OrderTermRefundRequest{
			ReferenceID: "order_ref_1", TermID: "term_id_1", Amount: 10,
		})
		require.NoError(t, err)
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, "/order/term/refund", req.Path)
		assert.JSONEq(t, `{"reference_id":"order_ref_1","term_id":"term_id_1","amount":10}`, req.Body)
		assert.Equal(t, true, res["is_success"])
	})
}

func TestOrderActionRoutesMatchPanel(t *testing.T) {
	ctx := context.Background()

	t.Run("OrderManualCallback", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"message":"sent","code":200}`)
		_, err := tapsilat.NewCustomAPI(server.URL, "token").OrderManualCallback(ctx, "ref_1", "conv_1")
		require.NoError(t, err)
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, "/order/callback", req.Path)
		assert.JSONEq(t, `{"reference_id":"ref_1","conversation_id":"conv_1"}`, req.Body)
	})

	t.Run("OrderRelatedUpdate", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"status":200,"message":"updated","is_success":true}`)
		_, err := tapsilat.NewCustomAPI(server.URL, "token").OrderRelatedUpdate(ctx, "ref_1", "ref_2")
		require.NoError(t, err)
		assert.Equal(t, http.MethodPatch, req.Method)
		assert.Equal(t, "/order/releated", req.Path)
		assert.JSONEq(t, `{"reference_id":"ref_1","related_reference_id":"ref_2"}`, req.Body)
	})

	t.Run("RefundAllOrder", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"status":"success","message":"refunded","is_success":true}`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token").RefundAllOrder(ctx, "ref_1")
		require.NoError(t, err)
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, "/order/refund-all", req.Path)
		assert.JSONEq(t, `{"reference_id":"ref_1"}`, req.Body)
		assert.True(t, res.IsSuccess)
	})

	t.Run("ChargeOrderDecodesCamelCaseResponse", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"isSuccess":true,"message":"charged","statusCode":"00","errorCode":""}`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token").ChargeOrder(ctx, tapsilat.OrderChargeRequest{OrderReferenceID: "ref_1", Retry: true})
		require.NoError(t, err)
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, "/order/charge", req.Path)
		assert.JSONEq(t, `{"order_reference_id":"ref_1","retry":true}`, req.Body)
		assert.True(t, res.IsSuccess)
		assert.Equal(t, "00", res.StatusCode)
	})

	t.Run("GetOrderPaymentDetailsByConversationID", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"conversation_id":"conv_1"}`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token").GetOrderPaymentDetailsByConversationID(ctx, "conv_1")
		require.NoError(t, err)
		assert.Equal(t, http.MethodPost, req.Method)
		assert.Equal(t, "/order/payment-details", req.Path)
		assert.JSONEq(t, `{"conversation_id":"conv_1"}`, req.Body)
		assert.Equal(t, "conv_1", res["conversation_id"])
	})
}

func TestOrderReadsMatchPanel(t *testing.T) {
	ctx := context.Background()

	t.Run("GetOrderTransactionsDecodesArray", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `[{"id":"tx_1","amount":"10.00","status":"completed"}]`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token").GetOrderTransactions(ctx, "ref_1")
		require.NoError(t, err)
		assert.Equal(t, "/order/ref_1/transactions", req.Path)
		require.Len(t, res, 1)
		assert.Equal(t, "tx_1", res[0]["id"])
	})

	t.Run("GetOrderTransactionsDecodesNull", func(t *testing.T) {
		server, _ := newContractServer(t, http.StatusOK, "application/json", `null`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token").GetOrderTransactions(ctx, "ref_1")
		require.NoError(t, err)
		assert.Empty(t, res)
	})

	t.Run("GetOrderListWithFilter", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"rows":[],"page":1,"per_page":20,"total":0,"total_pages":0}`)
		status := 3
		_, err := tapsilat.NewCustomAPI(server.URL, "token").GetOrderListWithFilter(ctx, tapsilat.OrderListFilter{
			Page: 1, PerPage: 20, ConversationID: "conv 1", BuyerID: "buyer&1", Status: &status,
		})
		require.NoError(t, err)
		assert.Equal(t, "/order/list", req.Path)
		assert.Equal(t, "buyer_id=buyer%261&conversation_id=conv+1&page=1&per_page=20&status=3", req.Query)
	})

	t.Run("GetOrdersEscapesBuyerID", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"rows":[]}`)
		_, err := tapsilat.NewCustomAPI(server.URL, "token").GetOrders(ctx, "1", "10", "a&b")
		require.NoError(t, err)
		assert.Equal(t, "buyer_id=a%26b&page=1&per_page=10", req.Query)
	})

	t.Run("GetOrderDecodesPanelFields", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{
			"reference_id":"ref_1","redirect_success_url":"https://shop.test/ok","redirect_failure_url":"https://shop.test/fail",
			"metadata":[{"key":"k","value":"v"}],"customer_paid_amount":100.5,
			"payment_requests":[{"id":"pr_1","amount":"100.50","paid":true}],
			"refund_cancel_requests":[{"id":"rc_1","status":1,"type":2}],
			"basket_items":[{"id":"bi_1","mcc":"5411","item_payments":[{"id":"ip_1","type":"card","refunded_date":"2026-01-02"}]}]
		}`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token").GetOrder(ctx, "ref_1")
		require.NoError(t, err)
		assert.Equal(t, "/order/ref_1", req.Path)
		assert.Equal(t, "https://shop.test/ok", res.RedirectSuccessUrl)
		assert.Equal(t, 100.5, res.CustomerPaidAmount)
		require.Len(t, res.PaymentRequests, 1)
		assert.True(t, res.PaymentRequests[0].Paid)
		require.Len(t, res.RefundCancelRequests, 1)
		assert.Equal(t, uint64(2), res.RefundCancelRequests[0].Type)
		require.Len(t, res.BasketItems, 1)
		assert.Equal(t, "5411", res.BasketItems[0].MCC)
		assert.Equal(t, "card", res.BasketItems[0].ItemPayments[0].Type)
	})
}

func TestCardInstallmentRoutes(t *testing.T) {
	ctx := context.Background()

	t.Run("GetCardInstallments", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"is_installment_available":true,"installments":[{"installment":3,"amount":"100.00","total_amount":"104.50","commission_rate":4.5,"monthly_payment":"34.83"}]}`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token").GetCardInstallments(ctx, tapsilat.CardInstallmentsRequest{
			Bin: "54066200", Amount: 100.5, Currency: "TRY", IsThreeD: true,
		})
		require.NoError(t, err)
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "/card/installments", req.Path)
		assert.Equal(t, "amount=100.5&bin=54066200&currency=TRY&is_three_d=true", req.Query)
		assert.True(t, res.IsInstallmentAvailable)
		require.Len(t, res.Installments, 1)
		assert.Equal(t, uint64(3), res.Installments[0].Installment)
		assert.Equal(t, "34.83", res.Installments[0].MonthlyPayment)
	})

	t.Run("GetInstallmentTable", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"groups":[{"bank":"Bank","program":"MAXIMUM","program_code":"max","rows":[{"installment":2,"installment_amount":"50.00","total_amount":"100.00","is_interest_free":true}]}]}`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token").GetInstallmentTable(ctx, tapsilat.InstallmentTableRequest{Amount: 100, Currency: "TRY"})
		require.NoError(t, err)
		assert.Equal(t, "/card/installment-table", req.Path)
		assert.Equal(t, "amount=100&currency=TRY", req.Query)
		require.Len(t, res.Groups, 1)
		assert.Equal(t, "MAXIMUM", res.Groups[0].Program)
		assert.True(t, res.Groups[0].Rows[0].IsInterestFree)
	})
}

func TestPanelErrorHandling(t *testing.T) {
	ctx := context.Background()

	t.Run("ValidationErrorKeepsIntegerCodeAndFieldErrors", func(t *testing.T) {
		server, _ := newContractServer(t, http.StatusBadRequest, "application/json", `{"message":"INVALID_REQUEST_BODY","code":4000000,"errors":{"email":"REQUIRED"}}`)
		_, err := tapsilat.NewCustomAPI(server.URL, "token").CreateSubmerchant(ctx, tapsilat.SubmerchantCreateRequest{})
		var apiErr *tapsilat.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, "4000000", apiErr.Code)
		assert.Equal(t, "INVALID_REQUEST_BODY", apiErr.Message)
		assert.Equal(t, "REQUIRED", apiErr.Errors["email"])
	})

	t.Run("HTMLPageBecomesAPIError", func(t *testing.T) {
		server, _ := newContractServer(t, http.StatusOK, "text/html; charset=utf-8", `<!doctype html><html><body>panel</body></html>`)
		_, err := tapsilat.NewCustomAPI(server.URL, "token").GetOrder(ctx, "ref_1")
		var apiErr *tapsilat.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusOK, apiErr.StatusCode)
		assert.Contains(t, apiErr.Message, "non-JSON response")
		assert.Contains(t, apiErr.RawBody, "<html>")
	})

	t.Run("IDsArePathEscaped", func(t *testing.T) {
		server, _ := newContractServer(t, http.StatusOK, "application/json", `{}`)
		var rawPath string
		server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawPath = r.URL.EscapedPath()
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{}`))
		})
		_, err := tapsilat.NewCustomAPI(server.URL, "token").GetVpos(ctx, "v/1")
		require.NoError(t, err)
		assert.Equal(t, "/vpos/v%2F1", rawPath)
	})
}

func TestApprovalPolicyCurrencyResolution(t *testing.T) {
	ctx := context.Background()
	const currencyID = "9f4050e8-1111-4f25-b4ef-aaaaaaaaaaaa"

	t.Run("ResolvesUnitToUUID", func(t *testing.T) {
		var policyBody string
		lookups := 0
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch r.URL.Path {
			case "/organization/currencies":
				lookups++
				_, _ = w.Write([]byte(`{"currencies":[{"id":"` + currencyID + `","currency_unit":"TRY"}]}`))
			case "/api/v2/marketplace/approval-policies":
				data, _ := io.ReadAll(r.Body)
				policyBody = string(data)
				w.WriteHeader(http.StatusCreated)
				_, _ = w.Write([]byte(`{"status":"created"}`))
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()

		err := tapsilat.NewCustomAPI(server.URL, "token").CreateMarketplaceApprovalPolicy(ctx, tapsilat.MarketplaceApprovalPolicyInput{Name: "Manual", Priority: 1, Mode: "manual", CurrencyID: "try"})
		require.NoError(t, err)
		assert.Equal(t, 1, lookups)
		assert.Contains(t, policyBody, `"currency_id":"`+currencyID+`"`)
	})

	t.Run("LeavesUUIDAndEmptyUnchanged", func(t *testing.T) {
		for _, ref := range []string{currencyID, ""} {
			var policyBody string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.NotEqual(t, "/organization/currencies", r.URL.Path)
				data, _ := io.ReadAll(r.Body)
				policyBody = string(data)
				w.WriteHeader(http.StatusOK)
			}))
			err := tapsilat.NewCustomAPI(server.URL, "token").UpdateMarketplaceApprovalPolicy(ctx, "policy_1", tapsilat.MarketplaceApprovalPolicyInput{Name: "Manual", Priority: 1, Mode: "manual", CurrencyID: ref})
			server.Close()
			require.NoError(t, err)
			assert.Contains(t, policyBody, `"currency_id":"`+ref+`"`)
		}
	})
}

func TestResolveCurrencyIDAcceptsUUIDForms(t *testing.T) {
	const canonical = "9f4050e8-1111-4f25-b4ef-aaaaaaaaaaaa"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("a UUID must not trigger a currency lookup, got %s", r.URL.Path)
	}))
	defer server.Close()
	api := tapsilat.NewCustomAPI(server.URL, "token")

	for _, ref := range []string{
		canonical,
		"9F4050E8-1111-4F25-B4EF-AAAAAAAAAAAA",
		"{9f4050e8-1111-4f25-b4ef-aaaaaaaaaaaa}",
		"urn:uuid:9f4050e8-1111-4f25-b4ef-aaaaaaaaaaaa",
		"9f4050e811114f25b4efaaaaaaaaaaaa",
		"  " + canonical + "  ",
	} {
		resolved, err := api.ResolveCurrencyID(context.Background(), ref)
		require.NoError(t, err, ref)
		assert.Equal(t, canonical, resolved, ref)
	}
}

func TestPanelResponseFieldsDecode(t *testing.T) {
	ctx := context.Background()

	t.Run("SellerAccountProviderSubmerchantKey", func(t *testing.T) {
		server, req := newContractServer(t, http.StatusOK, "application/json", `{"rows":[{"vpos_id":"v_1","status":"completed","provider_submerchant_key":"psk_1"}]}`)
		res, err := tapsilat.NewCustomAPI(server.URL+"/api/v1", "token").GetMarketplaceSellerAccounts(ctx, "seller_1")
		require.NoError(t, err)
		assert.Equal(t, "/api/v2/submerchants/seller_1/providers", req.Path)
		require.Len(t, res.Rows, 1)
		assert.Equal(t, "psk_1", res.Rows[0].ProviderSubmerchantKey)
	})

	t.Run("VposReadIncludesConnectionSettings", func(t *testing.T) {
		server, _ := newContractServer(t, http.StatusOK, "application/json", `{"id":"v_1","failover":true,"pf_id":"pf_1","redirect_uri":"https://pos.test/cb","tls":true,"root_ca":"ca","main_installment":true}`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token").GetVpos(ctx, "v_1")
		require.NoError(t, err)
		assert.True(t, res.Failover)
		assert.True(t, res.TLS)
		assert.Equal(t, "pf_1", res.PfID)
		assert.True(t, res.MainInstallment)
	})

	t.Run("SuborganizationListItemFields", func(t *testing.T) {
		server, _ := newContractServer(t, http.StatusOK, "application/json", `{"rows":[{"id":"so_1","name":"Branch","city":"Istanbul","currency":"TRY","logo":"https://cdn.test/l.png"}]}`)
		res, err := tapsilat.NewCustomAPI(server.URL, "token").GetSuborganizations(ctx, 1, 10)
		require.NoError(t, err)
		require.Len(t, res.Rows, 1)
		assert.Equal(t, "Istanbul", res.Rows[0].City)
		assert.Equal(t, "TRY", res.Rows[0].Currency)
	})
}
