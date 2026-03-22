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

func TestGetOrderSubmerchants(t *testing.T) {
	t.Run("SendsExpectedRequestAndParsesResponse", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/order/submerchants", r.URL.Path)
			assert.Equal(t, "2", r.URL.Query().Get("page"))
			assert.Equal(t, "5", r.URL.Query().Get("per_page"))
			assert.Equal(t, "Bearer token_123", r.Header.Get("Authorization"))
			assert.Equal(t, "application/json", r.Header.Get("Accept"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"page":2,"per_page":5,"total":1,"total_pages":1,"rows":[{"id":"sm_1"}]}`))
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_123")
		res, err := api.GetOrderSubmerchants(context.Background(), 2, 5)
		require.NoError(t, err)
		assert.Equal(t, int64(2), res.Page)
		assert.Equal(t, int64(5), res.PerPage)
		assert.Equal(t, int64(1), res.Total)
		assert.Equal(t, 1, res.TotalPages)
	})

	t.Run("ReturnsErrorForHttpFailure", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"bad_request"}`))
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_123")
		_, err := api.GetOrderSubmerchants(context.Background(), 1, 10)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "status 400")

		var apiErr *tapsilat.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
		assert.Equal(t, "bad_request", apiErr.Message)
	})
}

func TestAPIErrorNormalization(t *testing.T) {
	t.Run("ParsesMessageCodeAndStatus", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"status":"error","code":401001,"message":"unauthorized"}`))
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_401")
		_, err := api.GetOrganizationSettings(context.Background())
		require.Error(t, err)

		var apiErr *tapsilat.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
		assert.Equal(t, "error", apiErr.Status)
		assert.Equal(t, "401001", apiErr.Code)
		assert.Equal(t, "unauthorized", apiErr.Message)
	})

	t.Run("FallsBackToRawBodyForNonJsonErrors", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadGateway)
			_, _ = w.Write([]byte(`upstream exploded`))
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_502")
		_, err := api.GetOrganizationSettings(context.Background())
		require.Error(t, err)

		var apiErr *tapsilat.APIError
		require.ErrorAs(t, err, &apiErr)
		assert.Equal(t, http.StatusBadGateway, apiErr.StatusCode)
		assert.Equal(t, `upstream exploded`, apiErr.RawBody)
		assert.Empty(t, apiErr.Code)
	})
}

func TestGetOrganizationSettings(t *testing.T) {
	t.Run("ParsesSettingsResponse", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/organization/settings", r.URL.Path)
			assert.Equal(t, "Bearer token_456", r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"ttl": 3600,
				"retry_count": 3,
				"allow_payment": true,
				"session_ttl": 1800,
				"custom_checkout": false,
				"domain_address": "example.com",
				"checkout_domain": "checkout.example.com",
				"subscription_domain": "sub.example.com"
			}`))
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_456")
		res, err := api.GetOrganizationSettings(context.Background())
		require.NoError(t, err)
		assert.Equal(t, uint64(3600), res.Ttl)
		assert.Equal(t, uint64(3), res.RetryCount)
		assert.True(t, res.AllowPayment)
		assert.Equal(t, uint64(1800), res.SessionTtl)
		assert.False(t, res.CustomCheckout)
		assert.Equal(t, "example.com", res.DomainAddress)
		assert.Equal(t, "checkout.example.com", res.CheckoutDomain)
		assert.Equal(t, "sub.example.com", res.SubscriptionDomain)
	})
}

func TestCreateSubmerchant(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/submerchants", r.URL.Path)
		assert.Equal(t, "Bearer token_sm", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.JSONEq(t, `{
			"locale":"tr",
			"conversation_id":"conv_1",
			"name":"Tenant A",
			"email":"tenant@example.com",
			"gsm_number":"+905551112233",
			"address":"Istanbul",
			"iban":"TR00",
			"tax_office":"Besiktas",
			"legal_company_title":"Tenant A Ltd",
			"currency_id":"currency_1",
			"sub_merchant_external_id":"tenant-ext-1",
			"identity_number":"",
			"sub_merchant_type":"PRIVATE_COMPANY",
			"tax_number":"1234567890",
			"sub_merchant_key":"",
			"organization_id":"org_1",
			"status":"active",
			"system_time":1710000000,
			"contact_name":"Jane",
			"contact_surname":"Doe"
		}`, string(body))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":200,"message":"created"}`))
	}))
	defer server.Close()

	api := tapsilat.NewCustomAPI(server.URL, "token_sm")
	res, err := api.CreateSubmerchant(context.Background(), tapsilat.SubmerchantCreateRequest{
		Locale:                "tr",
		ConversationID:        "conv_1",
		Name:                  "Tenant A",
		Email:                 "tenant@example.com",
		GsmNumber:             "+905551112233",
		Address:               "Istanbul",
		Iban:                  "TR00",
		TaxOffice:             "Besiktas",
		LegalCompanyTitle:     "Tenant A Ltd",
		CurrencyID:            "currency_1",
		SubmerchantExternalID: "tenant-ext-1",
		IdentityNumber:        "",
		SubmerchantType:       "PRIVATE_COMPANY",
		TaxNumber:             "1234567890",
		SubmerchantKey:        "",
		OrganizationID:        "org_1",
		Status:                "active",
		SystemTime:            1710000000,
		ContactName:           "Jane",
		ContactSurname:        "Doe",
	})
	require.NoError(t, err)
	assert.Equal(t, uint64(200), res.Code)
	assert.Equal(t, "created", res.Message)
}

func TestGetSubmerchant(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/submerchants/sub_1", r.URL.Path)
		assert.Equal(t, "Bearer token_sm", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id":"sub_1",
			"name":"Tenant A",
			"email":"tenant@example.com",
			"sub_merchant_key":"sm_key_1",
			"organization_id":"org_1",
			"status":"active"
		}`))
	}))
	defer server.Close()

	api := tapsilat.NewCustomAPI(server.URL, "token_sm")
	res, err := api.GetSubmerchant(context.Background(), "sub_1")
	require.NoError(t, err)
	assert.Equal(t, "sub_1", res.ID)
	assert.Equal(t, "Tenant A", res.Name)
	assert.Equal(t, "sm_key_1", res.SubmerchantKey)
	assert.Equal(t, "org_1", res.OrganizationID)
}

func TestListSubmerchants(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/submerchants", r.URL.Path)
		assert.Equal(t, "3", r.URL.Query().Get("page"))
		assert.Equal(t, "20", r.URL.Query().Get("per_page"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"page":3,
			"per_page":20,
			"total":1,
			"total_pages":1,
			"row":[{
				"id":"sub_1",
				"name":"Tenant A",
				"email":"tenant@example.com",
				"submerchant_type":"PRIVATE_COMPANY",
				"submerchant_key":"sm_key_1",
				"status":"active"
			}]
		}`))
	}))
	defer server.Close()

	api := tapsilat.NewCustomAPI(server.URL, "token_sm")
	res, err := api.ListSubmerchants(context.Background(), 3, 20)
	require.NoError(t, err)
	assert.Equal(t, int64(3), res.Page)
	assert.Equal(t, int64(20), res.PerPage)
	assert.Len(t, res.Rows, 1)
	assert.Equal(t, "sub_1", res.Rows[0].ID)
	assert.Equal(t, "PRIVATE_COMPANY", res.Rows[0].SubmerchantType)
}

func TestUpdateSubmerchant(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		assert.Equal(t, "/submerchants/sub_1", r.URL.Path)
		assert.Equal(t, "Bearer token_sm", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		assert.Contains(t, string(body), `"status":"passive"`)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":200,"message":"updated"}`))
	}))
	defer server.Close()

	api := tapsilat.NewCustomAPI(server.URL, "token_sm")
	res, err := api.UpdateSubmerchant(context.Background(), "sub_1", tapsilat.SubmerchantUpdateRequest{
		Locale:                "tr",
		ConversationID:        "conv_2",
		Name:                  "Tenant A",
		Email:                 "tenant@example.com",
		GsmNumber:             "+905551112233",
		Address:               "Istanbul",
		Iban:                  "TR00",
		TaxOffice:             "Besiktas",
		LegalCompanyTitle:     "Tenant A Ltd",
		CurrencyID:            "currency_1",
		SubmerchantExternalID: "tenant-ext-1",
		SubmerchantType:       "PRIVATE_COMPANY",
		TaxNumber:             "1234567890",
		OrganizationID:        "org_1",
		Status:                "passive",
		SystemTime:            1710000001,
		ContactName:           "Jane",
		ContactSurname:        "Doe",
	})
	require.NoError(t, err)
	assert.Equal(t, uint64(200), res.Code)
	assert.Equal(t, "updated", res.Message)
}

func TestDeleteSubmerchant(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/submerchants/sub_1", r.URL.Path)
		assert.Equal(t, "Bearer token_sm", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"code":200,"message":"deleted"}`))
	}))
	defer server.Close()

	api := tapsilat.NewCustomAPI(server.URL, "token_sm")
	res, err := api.DeleteSubmerchant(context.Background(), "sub_1")
	require.NoError(t, err)
	assert.Equal(t, uint64(200), res.Code)
	assert.Equal(t, "deleted", res.Message)
}

func TestGetSuborganizations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/organization/suborganizations", r.URL.Path)
		assert.Equal(t, "1", r.URL.Query().Get("page"))
		assert.Equal(t, "15", r.URL.Query().Get("per_page"))
		assert.Equal(t, "Bearer token_org", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"page":1,
			"per_page":15,
			"total":2,
			"total_pages":1,
			"rows":[
				{"id":"org_sub_1","name":"Tenant A Scope"},
				{"id":"org_sub_2","name":"Tenant B Scope"}
			]
		}`))
	}))
	defer server.Close()

	api := tapsilat.NewCustomAPI(server.URL, "token_org")
	res, err := api.GetSuborganizations(context.Background(), 1, 15)
	require.NoError(t, err)
	assert.Equal(t, int64(2), res.Total)
	assert.Len(t, res.Rows, 2)
	assert.Equal(t, "org_sub_1", res.Rows[0].ID)
	assert.Equal(t, "Tenant B Scope", res.Rows[1].Name)
}

func TestGetSuborganization(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/organization/suborganizations/org_sub_1", r.URL.Path)
		assert.Equal(t, "Bearer token_org", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"id":"org_sub_1","name":"Tenant A Scope"}`))
	}))
	defer server.Close()

	api := tapsilat.NewCustomAPI(server.URL, "token_org")
	res, err := api.GetSuborganization(context.Background(), "org_sub_1")
	require.NoError(t, err)
	assert.Equal(t, "org_sub_1", res.ID)
	assert.Equal(t, "Tenant A Scope", res.Name)
}

func TestGetSuborganizationDetail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/organization/suborganizations/org_sub_1", r.URL.Path)
		assert.Equal(t, "Bearer token_org", r.Header.Get("Authorization"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id":"org_sub_1",
			"name":"Tenant A Scope",
			"parent_id":"org_main_1",
			"public_status":1,
			"availability_status":1,
			"created_at":"2026-03-16T11:00:00Z",
			"updated_at":"2026-03-16T11:05:00Z"
		}`))
	}))
	defer server.Close()

	api := tapsilat.NewCustomAPI(server.URL, "token_org")
	res, err := api.GetSuborganizationDetail(context.Background(), "org_sub_1")
	require.NoError(t, err)
	assert.Equal(t, "org_sub_1", res.ID)
	assert.Equal(t, "org_main_1", res.ParentID)
	assert.Equal(t, int64(1), res.PublicStatus)
	assert.Equal(t, "2026-03-16T11:00:00Z", res.CreatedAt)
}

func TestVposMethods(t *testing.T) {
	t.Run("ListVpos", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/vpos", r.URL.Path)
			assert.Equal(t, "2", r.URL.Query().Get("page"))
			assert.Equal(t, "25", r.URL.Query().Get("per_page"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"page":2,
				"per_page":25,
				"total":1,
				"total_pages":1,
				"rows":[{"id":"v_1","name":"Akbank POS","bank_name":"Akbank","env_mode":"test","provider":"akbank","payment_mode":"3d"}]
			}`))
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_vpos")
		res, err := api.ListVpos(context.Background(), 2, 25)
		require.NoError(t, err)
		assert.Equal(t, int64(2), res.Page)
		assert.Len(t, res.Rows, 1)
		assert.Equal(t, "v_1", res.Rows[0].ID)
	})

	t.Run("ListVposWithFilter", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/vpos", r.URL.Path)
			assert.Equal(t, "1", r.URL.Query().Get("page"))
			assert.Equal(t, "10", r.URL.Query().Get("per_page"))
			assert.Equal(t, "org_sub_1", r.URL.Query().Get("suborganization_id"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"page":1,
				"per_page":10,
				"total":1,
				"total_pages":1,
				"rows":[{"id":"v_2","name":"Scoped POS","bank_name":"Garanti","env_mode":"prod","provider":"garanti","payment_mode":"auth"}]
			}`))
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_vpos")
		res, err := api.ListVposWithFilter(context.Background(), 1, 10, tapsilat.VposListFilter{SuborganizationID: "org_sub_1"})
		require.NoError(t, err)
		assert.Equal(t, int64(1), res.Total)
		assert.Len(t, res.Rows, 1)
		assert.Equal(t, "v_2", res.Rows[0].ID)
	})

	t.Run("CreateVpos", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPost, r.Method)
			assert.Equal(t, "/vpos", r.URL.Path)

			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			assert.Contains(t, string(body), `"name":"Akbank POS"`)
			assert.Contains(t, string(body), `"acquirer_id":"acq_1"`)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"code":200,"message":"created"}`))
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_vpos")
		res, err := api.CreateVpos(context.Background(), tapsilat.VposCreateRequest{
			Name:        "Akbank POS",
			BankName:    "Akbank",
			EnvMode:     "test",
			PaymentMode: "3d",
			AcquirerID:  "acq_1",
			CardSchemes: []string{"visa"},
			Currencies:  []string{"try"},
		})
		require.NoError(t, err)
		assert.Equal(t, uint64(200), res.Code)
	})

	t.Run("GetVpos", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/vpos/v_1", r.URL.Path)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"id":"v_1","name":"Akbank POS","bank_name":"Akbank","payment_mode":"3d","currencies":["try"]}`))
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_vpos")
		res, err := api.GetVpos(context.Background(), "v_1")
		require.NoError(t, err)
		assert.Equal(t, "v_1", res.ID)
		assert.Equal(t, "Akbank POS", res.Name)
	})

	t.Run("UpdateVpos", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodPatch, r.Method)
			assert.Equal(t, "/vpos/v_1", r.URL.Path)

			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			assert.Contains(t, string(body), `"payment_mode":"non3d"`)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"code":200,"message":"updated"}`))
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_vpos")
		res, err := api.UpdateVpos(context.Background(), "v_1", tapsilat.VposUpdateRequest{
			Name:        "Akbank POS",
			BankName:    "Akbank",
			EnvMode:     "prod",
			PaymentMode: "non3d",
			AcquirerID:  "acq_1",
			CardSchemes: []string{"visa", "mastercard"},
			Currencies:  []string{"try"},
		})
		require.NoError(t, err)
		assert.Equal(t, uint64(200), res.Code)
	})

	t.Run("DeleteVpos", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodDelete, r.Method)
			assert.Equal(t, "/vpos/v_1", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"code":200,"message":"deleted"}`))
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_vpos")
		res, err := api.DeleteVpos(context.Background(), "v_1")
		require.NoError(t, err)
		assert.Equal(t, uint64(200), res.Code)
	})

	t.Run("ListVposAcquirersAndCardSchemes", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			switch r.URL.Path {
			case "/vpos/acquirers":
				_, _ = w.Write([]byte(`{"items":[{"id":"acq_1","name":"Akbank","prefix":"akbank"}]}`))
			case "/vpos/card-schemes":
				_, _ = w.Write([]byte(`{"items":[{"id":"visa","name":"Visa"}]}`))
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_vpos")
		acqRes, err := api.ListVposAcquirers(context.Background())
		require.NoError(t, err)
		assert.Len(t, acqRes.Items, 1)
		assert.Equal(t, "acq_1", acqRes.Items[0].ID)

		cardRes, err := api.ListCardSchemes(context.Background())
		require.NoError(t, err)
		assert.Len(t, cardRes.Items, 1)
		assert.Equal(t, "visa", cardRes.Items[0].ID)
	})
}

func TestVposSubmerchantMethods(t *testing.T) {
	t.Run("ListVposSubmerchants", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/vpos-submerchant", r.URL.Path)
			assert.Equal(t, "v_1", r.URL.Query().Get("vpos_id"))
			assert.Equal(t, "ext_1", r.URL.Query().Get("external_reference_id"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"page":1,"per_page":10,"total":1,"total_pages":1,"rows":[{"id":"vs_1","external_reference_id":"ext_1","submerchant_id":"sub_1","terminal_no":"1234","vpos_id":"v_1"}]}`))
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_vsm")
		res, err := api.ListVposSubmerchants(context.Background(), 1, 10, "v_1", "ext_1")
		require.NoError(t, err)
		assert.Len(t, res.Rows, 1)
		assert.Equal(t, "vs_1", res.Rows[0].ID)
	})

	t.Run("CreateGetUpdateDeleteVposSubmerchant", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			switch {
			case r.Method == http.MethodPost && r.URL.Path == "/vpos-submerchant":
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				assert.Contains(t, string(body), `"submerchant_id":"sub_1"`)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"code":200,"message":"created"}`))
			case r.Method == http.MethodGet && r.URL.Path == "/vpos-submerchant/vs_1":
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"id":"vs_1","external_reference_id":"ext_1","submerchant_id":"sub_1","terminal_no":"1234","vpos_id":"v_1","mcc":"5411"}`))
			case r.Method == http.MethodPatch && r.URL.Path == "/vpos-submerchant/vs_1":
				body, err := io.ReadAll(r.Body)
				require.NoError(t, err)
				assert.Contains(t, string(body), `"title":"Tenant Updated"`)
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"code":200,"message":"updated"}`))
			case r.Method == http.MethodDelete && r.URL.Path == "/vpos-submerchant/vs_1":
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{"code":200,"message":"deleted"}`))
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_vsm")
		createRes, err := api.CreateVposSubmerchant(context.Background(), tapsilat.VposSubmerchantCreateRequest{
			ExternalReferenceID: "ext_1",
			SubmerchantID:       "sub_1",
			TerminalNo:          "1234",
			VposID:              "v_1",
			Title:               "Tenant Title",
		})
		require.NoError(t, err)
		assert.Equal(t, uint64(200), createRes.Code)

		readRes, err := api.GetVposSubmerchant(context.Background(), "vs_1")
		require.NoError(t, err)
		assert.Equal(t, "vs_1", readRes.ID)
		assert.Equal(t, "5411", readRes.MCC)

		updateRes, err := api.UpdateVposSubmerchant(context.Background(), "vs_1", tapsilat.VposSubmerchantUpdateRequest{
			ExternalReferenceID: "ext_1",
			SubmerchantID:       "sub_1",
			TerminalNo:          "1234",
			Title:               "Tenant Updated",
		})
		require.NoError(t, err)
		assert.Equal(t, uint64(200), updateRes.Code)

		deleteRes, err := api.DeleteVposSubmerchant(context.Background(), "vs_1")
		require.NoError(t, err)
		assert.Equal(t, uint64(200), deleteRes.Code)
	})
}

func TestMappingHelpers(t *testing.T) {
	t.Run("GetSuborganizationBySubmerchant", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/submerchants/sub_1/suborganization", r.URL.Path)
			assert.Equal(t, "Bearer token_map", r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"submerchant_id":"sub_1","suborganization_id":"org_sub_1"}`))
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_map")
		res, err := api.GetSuborganizationBySubmerchant(context.Background(), "sub_1")
		require.NoError(t, err)
		assert.Equal(t, "sub_1", res.SubmerchantID)
		assert.Equal(t, "org_sub_1", res.SuborganizationID)
	})

	t.Run("GetSubmerchantBySuborganization", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/organization/suborganizations/org_sub_1/submerchant", r.URL.Path)
			assert.Equal(t, "Bearer token_map", r.Header.Get("Authorization"))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"suborganization_id":"org_sub_1","submerchant_id":"sub_1"}`))
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_map")
		res, err := api.GetSubmerchantBySuborganization(context.Background(), "org_sub_1")
		require.NoError(t, err)
		assert.Equal(t, "org_sub_1", res.SuborganizationID)
		assert.Equal(t, "sub_1", res.SubmerchantID)
	})
}
