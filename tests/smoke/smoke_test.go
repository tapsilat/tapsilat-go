package smoke_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

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
			_, err := api.GetVpos(ctx, vposID)
			require.NoError(t, err, "endpoint: %s", endpoint)

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
