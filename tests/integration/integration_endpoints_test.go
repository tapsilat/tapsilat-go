package integration_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tapsilat/tapsilat-go"
)

func TestNewOrderMethods(t *testing.T) {
	token := "your_test_token_here"
	if token == "your_test_token_here" {
		t.Skip("Please set a real token for integration tests")
	}

	api := tapsilat.NewAPI(token)
	
	t.Run("GetSystemOrderStatuses", func(t *testing.T) {
		res, err := api.GetSystemOrderStatuses(context.Background())
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("GetOrganizationLimits", func(t *testing.T) {
		res, err := api.GetOrganizationLimits(context.Background())
		require.NoError(t, err)
		require.NotNil(t, res)
	})

	t.Run("GetOrderPaymentDetailsByID", func(t *testing.T) {
		res, err := api.GetOrderPaymentDetailsByID(context.Background(), "invalid_ref")
		// It might fail with APIError but err itself shouldn't be nil
		require.Error(t, err)
		require.NotNil(t, res) // maps are returned even on error
	})
	
	// Adding dummy invocations for all modified and new endpoints to ensure they compile and don't panic
	t.Run("DummyInvocations", func(t *testing.T) {
		api.OrderAccounting(context.Background(), tapsilat.OrderAccountingRequest{})
		api.OrderPostAuth(context.Background(), tapsilat.OrderPostAuthRequest{})
		api.OrderPaymentOptionsUpdate(context.Background(), tapsilat.OrderPaymentOptionsUpdateDTO{})
		api.SplitOrderItemPayment(context.Background(), tapsilat.SplitOrderItemPaymentDTO{})
		api.OrderCallback(context.Background(), "id")
		api.OrderVposQuery(context.Background(), "id")
		api.AddBasketItem(context.Background(), tapsilat.AddBasketItemRequest{})
		api.RemoveBasketItem(context.Background(), tapsilat.RemoveBasketItemRequest{})
		api.UpdateBasketItem(context.Background(), tapsilat.UpdateBasketItemRequest{})
		
		api.GetSystemBasketItemTypes(context.Background())
		api.GetSystemErrorCodes(context.Background())
		api.GetSystemPaymentTermStatuses(context.Background())
		api.GetSystemProductTypes(context.Background())
		api.GetSystemShortcutTypes(context.Background())
		api.GetSystemTransactionPaymentTypes(context.Background())
		api.GetSystemTransactionPurposes(context.Background())
		api.GetSystemTransactionStatuses(context.Background())
		
		api.GetOrganizationCallback(context.Background())
		api.UpdateOrganizationCallback(context.Background(), tapsilat.CallbackURLDTO{})
		api.CreateOrganizationBusiness(context.Background(), tapsilat.OrgCreateBusinessRequest{})
		api.GetOrganizationLimitUser(context.Background(), tapsilat.GetUserLimitRequest{})
		api.SetOrganizationLimitUser(context.Background(), tapsilat.SetLimitUserRequest{})
		api.GetOrganizationMeta(context.Background(), "test")
		api.GetOrganizationScopes(context.Background())
		api.CreateOrganizationUser(context.Background(), tapsilat.OrgCreateUserReq{})
		api.VerifyOrganizationUser(context.Background(), tapsilat.OrgUserVerifyReq{})
		api.VerifyOrganizationUserMobile(context.Background(), tapsilat.OrgUserMobileVerifyReq{})
		
		api.GetOrderPaymentDetails(context.Background(), tapsilat.OrderPaymentDetailDTO{})
		api.OrderManualCallback(context.Background(), tapsilat.OrderManualCallbackDTO{})
		api.DeleteOrderTerm(context.Background(), tapsilat.OrderPaymentTermDeleteDTO{})
		api.UpdateOrderTerm(context.Background(), tapsilat.OrderPaymentTermUpdateDTO{})
		api.GetOrderTerm(context.Background(), "ref")
		api.OrderRelatedUpdate(context.Background(), tapsilat.OrderRelatedReferenceDTO{})
	})
}
