package unit_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

func readContractFixture(t *testing.T, fileName string) []byte {
	t.Helper()
	_, thisFile, _, ok := runtime.Caller(0)
	require.True(t, ok)
	fixturePath := filepath.Join(filepath.Dir(thisFile), "..", "fixtures", "contracts", fileName)
	data, err := os.ReadFile(fixturePath)
	require.NoError(t, err)
	return data
}

func TestContractFixturesDecode(t *testing.T) {
	t.Run("SubmerchantReadFixture", func(t *testing.T) {
		fixture := readContractFixture(t, "submerchant_read.json")
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/submerchants/sub_1", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(fixture)
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_fixture")
		res, err := api.GetSubmerchant(context.Background(), "sub_1")
		require.NoError(t, err)
		assert.Equal(t, "sub_1", res.ID)
		assert.Equal(t, "Tenant A", res.Name)
		assert.Equal(t, "sm_key_1", res.SubmerchantKey)
	})

	t.Run("VposReadFixture", func(t *testing.T) {
		fixture := readContractFixture(t, "vpos_read.json")
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/vpos/v_1", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(fixture)
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_fixture")
		res, err := api.GetVpos(context.Background(), "v_1")
		require.NoError(t, err)
		assert.Equal(t, "v_1", res.ID)
		assert.Equal(t, "Akbank POS", res.Name)
		assert.Equal(t, "acq_1", res.AcquirerID)
	})

	t.Run("VposSubmerchantListFixture", func(t *testing.T) {
		fixture := readContractFixture(t, "vpos_submerchant_list.json")
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/vpos-submerchant", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(fixture)
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_fixture")
		res, err := api.ListVposSubmerchants(context.Background(), 1, 10, "", "")
		require.NoError(t, err)
		assert.Equal(t, int64(1), res.Total)
		assert.Len(t, res.Rows, 1)
		assert.Equal(t, "vs_1", res.Rows[0].ID)
	})

	t.Run("MappingFixtures", func(t *testing.T) {
		submerchantFixture := readContractFixture(t, "submerchant_suborganization_mapping.json")
		suborgFixture := readContractFixture(t, "suborganization_submerchant_mapping.json")
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			switch r.URL.Path {
			case "/submerchants/sub_1/suborganization":
				_, _ = w.Write(submerchantFixture)
			case "/organization/suborganizations/org_sub_1/submerchant":
				_, _ = w.Write(suborgFixture)
			default:
				w.WriteHeader(http.StatusNotFound)
			}
		}))
		defer server.Close()

		api := tapsilat.NewCustomAPI(server.URL, "token_fixture")
		map1, err := api.GetSuborganizationBySubmerchant(context.Background(), "sub_1")
		require.NoError(t, err)
		assert.Equal(t, "org_sub_1", map1.SuborganizationID)

		map2, err := api.GetSubmerchantBySuborganization(context.Background(), "org_sub_1")
		require.NoError(t, err)
		assert.Equal(t, "sub_1", map2.SubmerchantID)
	})
}
