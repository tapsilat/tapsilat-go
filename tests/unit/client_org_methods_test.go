package unit_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

func TestOrgMethods(t *testing.T) {
	t.Run("GetOrganizationCallback", func(t *testing.T) {
		server := mockServer(t, http.MethodGet, "/organization/callback", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.GetOrganizationCallback(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("UpdateOrganizationCallback", func(t *testing.T) {
		server := mockServer(t, http.MethodPatch, "/organization/callback", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.UpdateOrganizationCallback(context.Background(), tapsilat.CallbackURLDTO{})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("CreateOrganizationBusiness", func(t *testing.T) {
		server := mockServer(t, http.MethodPost, "/organization/business/create", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.CreateOrganizationBusiness(context.Background(), tapsilat.OrgCreateBusinessRequest{})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("GetOrganizationLimitUser", func(t *testing.T) {
		server := mockServer(t, http.MethodGet, "/organization/limit/user", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.GetOrganizationLimitUser(context.Background(), tapsilat.GetUserLimitRequest{})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("SetOrganizationLimitUser", func(t *testing.T) {
		server := mockServer(t, http.MethodPost, "/organization/limit/user", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.SetOrganizationLimitUser(context.Background(), tapsilat.SetLimitUserRequest{})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("GetOrganizationLimits", func(t *testing.T) {
		server := mockServer(t, http.MethodGet, "/organization/limits", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.GetOrganizationLimits(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("GetOrganizationMeta", func(t *testing.T) {
		server := mockServer(t, http.MethodGet, "/organization/meta/test1", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.GetOrganizationMeta(context.Background(), "test1")
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("GetOrganizationScopes", func(t *testing.T) {
		server := mockServer(t, http.MethodGet, "/organization/scopes", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.GetOrganizationScopes(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("CreateOrganizationUser", func(t *testing.T) {
		server := mockServer(t, http.MethodPost, "/organization/user/create", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.CreateOrganizationUser(context.Background(), tapsilat.OrgCreateUserReq{})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("VerifyOrganizationUser", func(t *testing.T) {
		server := mockServer(t, http.MethodPost, "/organization/user/verify", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.VerifyOrganizationUser(context.Background(), tapsilat.OrgUserVerifyReq{})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("VerifyOrganizationUserMobile", func(t *testing.T) {
		server := mockServer(t, http.MethodPost, "/organization/user/verify-mobile", `{"status":"success"}`, nil)
		defer server.Close()
		api := tapsilat.NewCustomAPI(server.URL, "token")
		res, err := api.VerifyOrganizationUserMobile(context.Background(), tapsilat.OrgUserMobileVerifyReq{})
		require.NoError(t, err)
		assert.Equal(t, "success", res["status"])
	})

	t.Run("VerifyWebhook", func(t *testing.T) {
		payload := "hello"
		secret := "secret"
		signature := "sha256=88aab3ede8d3adf94d26ab90d3bafd4a2083070c3bcce9c014ee04a443847c0b"
		isValid := tapsilat.VerifyWebhook(payload, signature, secret)
		assert.True(t, isValid)

		isValidFalse := tapsilat.VerifyWebhook(payload, "sha256=invalid", secret)
		assert.False(t, isValidFalse)
	})
}
