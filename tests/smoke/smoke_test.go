package smoke_test

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

func requireSmokeEnv(t *testing.T) (*tapsilat.API, string, string, string, string) {
	t.Helper()
	endpoint := os.Getenv("TAPSILAT_SMOKE_ENDPOINT")
	token := os.Getenv("TAPSILAT_SMOKE_TOKEN")
	submerchantID := os.Getenv("TAPSILAT_SMOKE_SUBMERCHANT_ID")
	vposID := os.Getenv("TAPSILAT_SMOKE_VPOS_ID")
	suborganizationID := os.Getenv("TAPSILAT_SMOKE_SUBORGANIZATION_ID")

	if endpoint == "" || token == "" {
		t.Skip("set TAPSILAT_SMOKE_ENDPOINT and TAPSILAT_SMOKE_TOKEN for smoke tests")
	}

	api := tapsilat.NewCustomAPI(endpoint, token)
	return api, submerchantID, vposID, suborganizationID, endpoint
}

func TestSmokeReadAndListFlows(t *testing.T) {
	api, submerchantID, vposID, suborganizationID, endpoint := requireSmokeEnv(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	t.Run("OrganizationSettings", func(t *testing.T) {
		_, err := api.GetOrganizationSettings(ctx)
		require.NoError(t, err, "endpoint: %s", endpoint)
	})

	t.Run("OrganizationCurrencies", func(t *testing.T) {
		currencies, err := api.GetOrganizationCurrencies(ctx)
		require.NoError(t, err, "endpoint: %s", endpoint)
		for _, currency := range currencies.Currencies {
			require.NotEmpty(t, currency.ID, "endpoint: %s", endpoint)
			require.NoError(t, validateUUID(currency.ID), "endpoint: %s", endpoint)
			require.NotEmpty(t, currency.CurrencyUnit, "endpoint: %s", endpoint)
		}
	})

	t.Run("SubmerchantAndSuborganizationList", func(t *testing.T) {
		_, err := api.ListSubmerchants(ctx, 1, 10)
		require.NoError(t, err, "endpoint: %s", endpoint)

		_, err = api.GetSuborganizations(ctx, 1, 10)
		require.NoError(t, err, "endpoint: %s", endpoint)

		if suborganizationID != "" {
			_, err = api.GetSuborganization(ctx, suborganizationID)
			require.NoError(t, err, "endpoint: %s", endpoint)
		}
	})

	t.Run("VposListAndCardMetadata", func(t *testing.T) {
		_, err := api.ListVpos(ctx, 1, 10)
		require.NoError(t, err, "endpoint: %s", endpoint)

		if suborganizationID != "" {
			_, err = api.ListVposWithFilter(ctx, 1, 10, tapsilat.VposListFilter{SuborganizationID: suborganizationID})
			require.NoError(t, err, "endpoint: %s", endpoint)
		}

		_, err = api.ListVposAcquirers(ctx)
		require.NoError(t, err, "endpoint: %s", endpoint)

		_, err = api.ListCardSchemes(ctx)
		require.NoError(t, err, "endpoint: %s", endpoint)
	})

	t.Run("OptionalReadById", func(t *testing.T) {
		if submerchantID != "" {
			_, err := api.GetSubmerchant(ctx, submerchantID)
			require.NoError(t, err, "endpoint: %s", endpoint)

			mapping, err := api.GetSuborganizationBySubmerchant(ctx, submerchantID)
			require.NoError(t, err, "endpoint: %s", endpoint)
			if suborganizationID != "" {
				require.Equal(t, suborganizationID, mapping.SuborganizationID, "endpoint: %s", endpoint)
			}
		}
		if vposID != "" {
			vpos, err := api.GetVpos(ctx, vposID)
			require.NoError(t, err, "endpoint: %s", endpoint)
			require.NotEmpty(t, vpos.Currencies, "endpoint: %s", endpoint)
			for _, currencyID := range vpos.Currencies {
				require.NoError(t, validateUUID(currencyID), "endpoint: %s", endpoint)
			}

			_, err = api.ListVposSubmerchants(ctx, 1, 10, vposID, "")
			require.NoError(t, err, "endpoint: %s", endpoint)
		}
		if suborganizationID != "" {
			mapping, err := api.GetSubmerchantBySuborganization(ctx, suborganizationID)
			require.NoError(t, err, "endpoint: %s", endpoint)
			if submerchantID != "" {
				require.Equal(t, submerchantID, mapping.SubmerchantID, "endpoint: %s", endpoint)
			}
		}
	})
}

func validateUUID(value string) error {
	if uuidRegex.MatchString(value) {
		return nil
	}
	return &tapsilat.ValidationError{
		StatusCode: 400,
		Code:       0,
		Message:    "invalid UUID format",
	}
}
