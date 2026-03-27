package integration_test

import (
	"context"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89aAbB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)

func requireIntegrationEnv(t *testing.T) (*tapsilat.API, string, string, string, string) {
	t.Helper()

	endpoint := os.Getenv("TAPSILAT_IT_ENDPOINT")
	token := os.Getenv("TAPSILAT_IT_TOKEN")
	submerchantID := os.Getenv("TAPSILAT_IT_SUBMERCHANT_ID")
	suborganizationID := os.Getenv("TAPSILAT_IT_SUBORGANIZATION_ID")
	vposID := os.Getenv("TAPSILAT_IT_VPOS_ID")

	if endpoint == "" {
		endpoint = "https://panel.tapsilat.dev/api/v1"
	}

	if token == "" || submerchantID == "" || suborganizationID == "" {
		t.Skip("set TAPSILAT_IT_TOKEN, TAPSILAT_IT_SUBMERCHANT_ID and TAPSILAT_IT_SUBORGANIZATION_ID for integration tests")
	}

	api := tapsilat.NewCustomAPI(endpoint, token)
	return api, submerchantID, suborganizationID, vposID, endpoint
}

func containsVpos(rows []tapsilat.VposListItem, vposID string) bool {
	for _, row := range rows {
		if row.ID == vposID {
			return true
		}
	}

	return false
}

func TestBackendConsistency_SubmerchantSuborganizationAndScopedVpos(t *testing.T) {
	api, submerchantID, suborganizationID, vposID, endpoint := requireIntegrationEnv(t)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	t.Run("BidirectionalMappingConsistency", func(t *testing.T) {
		currencies, err := api.GetOrganizationCurrencies(ctx)
		require.NoError(t, err, "endpoint: %s", endpoint)
		for _, currency := range currencies.Currencies {
			require.NoError(t, assertUUID(currency.ID), "endpoint: %s", endpoint)
			require.NotEmpty(t, currency.CurrencyUnit, "endpoint: %s", endpoint)
		}

		submerchant, err := api.GetSubmerchant(ctx, submerchantID)
		require.NoError(t, err, "endpoint: %s", endpoint)
		require.Equal(t, submerchantID, submerchant.ID, "endpoint: %s", endpoint)

		suborganization, err := api.GetSuborganizationDetail(ctx, suborganizationID)
		require.NoError(t, err, "endpoint: %s", endpoint)
		require.Equal(t, suborganizationID, suborganization.ID, "endpoint: %s", endpoint)

		forwardMapping, err := api.GetSuborganizationBySubmerchant(ctx, submerchantID)
		require.NoError(t, err, "endpoint: %s", endpoint)
		require.Equal(t, submerchantID, forwardMapping.SubmerchantID, "endpoint: %s", endpoint)
		require.Equal(t, suborganizationID, forwardMapping.SuborganizationID, "endpoint: %s", endpoint)

		reverseMapping, err := api.GetSubmerchantBySuborganization(ctx, suborganizationID)
		require.NoError(t, err, "endpoint: %s", endpoint)
		require.Equal(t, suborganizationID, reverseMapping.SuborganizationID, "endpoint: %s", endpoint)
		require.Equal(t, submerchantID, reverseMapping.SubmerchantID, "endpoint: %s", endpoint)
	})

	t.Run("ScopedVposConsistency", func(t *testing.T) {
		allVpos, err := api.ListVpos(ctx, 1, 100)
		require.NoError(t, err, "endpoint: %s", endpoint)

		scopedVpos, err := api.ListVposWithFilter(ctx, 1, 100, tapsilat.VposListFilter{SuborganizationID: suborganizationID})
		require.NoError(t, err, "endpoint: %s", endpoint)
		require.LessOrEqual(t, scopedVpos.Total, allVpos.Total, "endpoint: %s", endpoint)

		if vposID != "" {
			require.True(t, containsVpos(scopedVpos.Rows, vposID), "expected vpos %s in scoped list, endpoint: %s", vposID, endpoint)

			vpos, err := api.GetVpos(ctx, vposID)
			require.NoError(t, err, "endpoint: %s", endpoint)
			require.Equal(t, vposID, vpos.ID, "endpoint: %s", endpoint)
			require.NotEmpty(t, vpos.Currencies, "endpoint: %s", endpoint)
			for _, currencyID := range vpos.Currencies {
				require.NoError(t, assertUUID(currencyID), "endpoint: %s", endpoint)
			}

			vposSubmerchantRows, err := api.ListVposSubmerchants(ctx, 1, 100, vposID, "")
			require.NoError(t, err, "endpoint: %s", endpoint)
			for _, row := range vposSubmerchantRows.Rows {
				require.Equal(t, vposID, row.VposID, "endpoint: %s", endpoint)
			}
		}
	})
}

func assertUUID(value string) error {
	if uuidRegex.MatchString(value) {
		return nil
	}
	return &tapsilat.ValidationError{
		StatusCode: 400,
		Code:       0,
		Message:    "invalid UUID format",
	}
}
