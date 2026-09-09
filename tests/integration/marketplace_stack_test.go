package integration_test

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func startMarketplaceStack(t *testing.T, root, packagePath string, extra []string) map[string]string {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "stack-server")
	build := exec.Command("go", "test", "-c", "-o", binary, packagePath)
	build.Dir = root
	output, err := build.CombinedOutput()
	require.NoError(t, err, string(output))
	command := exec.Command(binary, "-test.run=^TestMarketplaceStackServer$", "-test.v")
	command.Dir = root
	command.Env = append(os.Environ(), append([]string{"MARKETPLACE_STACK_SERVER=1", "testing=true"}, extra...)...)
	pipe, err := command.StdoutPipe()
	require.NoError(t, err)
	command.Stderr = command.Stdout
	require.NoError(t, command.Start())
	ready := make(chan map[string]string, 1)
	scanDone := make(chan string, 1)
	go func() {
		scanner := bufio.NewScanner(pipe)
		scanner.Buffer(make([]byte, 4096), 1024*1024)
		var transcript strings.Builder
		for scanner.Scan() {
			line := scanner.Text()
			transcript.WriteString(line + "\n")
			if strings.HasPrefix(line, "MARKETPLACE_STACK_READY:") {
				var value map[string]string
				if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "MARKETPLACE_STACK_READY:")), &value); err == nil {
					ready <- value
				}
			}
		}
		if err := scanner.Err(); err != nil {
			transcript.WriteString(err.Error())
		}
		scanDone <- transcript.String()
	}()
	finished := make(chan error, 1)
	go func() { finished <- command.Wait() }()
	t.Cleanup(func() {
		if err := command.Process.Signal(os.Interrupt); err != nil && !errors.Is(err, os.ErrProcessDone) {
			t.Logf("stack interrupt: %v", err)
		}
		select {
		case err := <-finished:
			if err != nil {
				t.Errorf("stack server exited unsuccessfully: %v", err)
			}
		case <-time.After(10 * time.Second):
			t.Error("stack server shutdown timed out")
			if err := command.Process.Kill(); err != nil {
				t.Logf("stack kill: %v", err)
			}
		}
	})
	select {
	case value := <-ready:
		return value
	case transcript := <-scanDone:
		t.Fatalf("stack server exited before ready:\n%s", transcript)
	case <-time.After(60 * time.Second):
		t.Fatal("stack server startup timed out")
	}
	return nil
}

func TestMarketplaceAdminPanelSDKStack(t *testing.T) {
	admin, panel := os.Getenv("TAPSILAT_STACK_ADMIN_ROOT"), os.Getenv("TAPSILAT_STACK_PANEL_ROOT")
	if admin == "" || panel == "" {
		t.Skip("set TAPSILAT_STACK_ADMIN_ROOT and TAPSILAT_STACK_PANEL_ROOT to backend module paths")
	}
	marker := make([]byte, 24)
	_, err := rand.Read(marker)
	require.NoError(t, err)
	token := hex.EncodeToString(marker)
	a := startMarketplaceStack(t, admin, "./app/api/grpcroutes/submerchant/v1", []string{"MARKETPLACE_STACK_TOKEN=" + token})
	p := startMarketplaceStack(t, panel, "./cmd/routes/v1", []string{"MARKETPLACE_STACK_ADMIN=" + a["address"]})
	api := tapsilat.NewCustomAPI(p["address"], token)
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	base := tapsilat.MarketplaceSubmerchantCreateRequest{Locale: "tr", ConversationID: "stack-create", Name: "SDK Seller",
		Email: "sdk@example.test", GsmNumber: "905551112233", CurrencyID: a["currency_id"], SubmerchantType: "PRIVATE_COMPANY",
		TaxNumber: "1234567890", Iban: "TR180006200119000006672315", ContactName: "Contact", ContactSurname: "Person", IdempotencyKey: "sdk-independent"}
	independent, err := api.CreateMarketplaceSubmerchant(ctx, base)
	require.NoError(t, err)
	require.Empty(t, independent.Provisionings)
	require.NotEmpty(t, independent.RoutingReference)
	replay, err := api.CreateMarketplaceSubmerchant(ctx, base)
	require.NoError(t, err)
	require.Equal(t, independent.ID, replay.ID)
	base.Name, base.Email, base.IdempotencyKey = "SDK selected", "selected@example.test", "sdk-selected"
	base.VposID = a["vpos_id"]
	selected, err := api.CreateMarketplaceSubmerchant(ctx, base)
	require.NoError(t, err)
	require.Len(t, selected.Provisionings, 1)
	require.Equal(t, a["vpos_id"], selected.Provisionings[0].VposID)
	accounts, err := api.GetMarketplaceSellerAccounts(ctx, selected.ID)
	require.NoError(t, err)
	mapped := 0
	for _, account := range accounts.Rows {
		if account.Mapped {
			mapped++
		}
		if account.Environment == "PROD" {
			require.False(t, account.Mapped)
		}
	}
	require.Equal(t, 1, mapped)
	contracts, err := api.GetMarketplaceSellerContracts(ctx, tapsilat.MarketplaceContractFilter{CurrencyID: a["currency_id"], SellerType: "PRIVATE_COMPANY", Operation: "create", VposIDs: []string{a["vpos_id"]}})
	require.NoError(t, err)
	require.True(t, contracts.Contracts[a["vpos_id"]].Fields["contact_name"])
	for _, account := range accounts.Rows {
		if account.VposID != a["vpos_id"] {
			continue
		}
		_, err = api.RunMarketplaceSellerAccountAction(ctx, selected.ID, account.VposID, tapsilat.MarketplaceAccountAction{Action: "retrieve", Environment: "PROD", Revision: account.Revision})
		require.Error(t, err)
		_, err = api.RunMarketplaceSellerAccountAction(ctx, selected.ID, account.VposID, tapsilat.MarketplaceAccountAction{Action: "retrieve", Environment: "TEST", Revision: account.Revision})
		require.NoError(t, err)
	}
	base.Name, base.Email, base.IdempotencyKey = "SDK multiple", "multiple@example.test", "sdk-multiple"
	base.VposID = ""
	targets := []tapsilat.MarketplaceAccountTarget{{VposID: a["vpos_id"], Environment: "TEST"}, {VposID: a["second_vpos_id"], Environment: "TEST"}}
	base.ProviderAccounts = &targets
	multiple, err := api.CreateMarketplaceSubmerchant(ctx, base)
	require.NoError(t, err)
	require.Len(t, multiple.Provisionings, 2)
	for _, item := range multiple.Provisionings {
		require.Equal(t, "completed", item.Status)
	}
	profile, err := api.GetMarketplaceSubmerchant(ctx, independent.ID)
	require.NoError(t, err)
	require.Equal(t, "Contact", profile.ContactName)
	replacement := tapsilat.MarketplaceSubmerchantCreateRequest{Locale: "tr", ConversationID: "edit", Name: "SDK Seller updated", Email: "sdk@example.test",
		CurrencyID: a["currency_id"], SubmerchantType: "PRIVATE_COMPANY", GsmNumber: "905551112233", ContactName: "Updated"}
	_, err = api.UpdateMarketplaceSubmerchant(ctx, independent.ID, tapsilat.MarketplaceProfileUpdate{MarketplaceSubmerchantCreateRequest: replacement, Revision: profile.Revision})
	require.NoError(t, err)
	_, err = api.UpdateMarketplaceSubmerchant(ctx, independent.ID, tapsilat.MarketplaceProfileUpdate{MarketplaceSubmerchantCreateRequest: replacement, Revision: profile.Revision})
	require.Error(t, err)
	rows, err := api.ListMarketplacePayouts(ctx, tapsilat.MarketplacePayoutFilter{OrderID: a["order_id"]})
	require.NoError(t, err)
	require.Len(t, rows.Rows, 1)
	payout := rows.Rows[0]
	require.Equal(t, "waiting", payout.EligibilityStatus)
	require.Equal(t, "service_completed", payout.RequiredEvent)
	_, err = api.ApproveMarketplacePayouts(ctx, tapsilat.MarketplacePayoutApproval{IDs: []string{payout.ID}})
	require.Error(t, err)
	event := tapsilat.SubmerchantPayoutEventRequest{IdempotencyKey: "stay-completed", OrderReference: a["order_reference"], ItemID: "stay", EventType: "service_completed"}
	first, err := api.RecordSubmerchantPayoutEvent(ctx, event)
	require.NoError(t, err)
	again, err := api.RecordSubmerchantPayoutEvent(ctx, event)
	require.NoError(t, err)
	require.Equal(t, first.ID, again.ID)
	wrong := event
	wrong.IdempotencyKey = "wrong-item"
	wrong.ItemID = "another"
	_, err = api.RecordSubmerchantPayoutEvent(ctx, wrong)
	require.Error(t, err)
	payout, err = api.GetMarketplacePayout(ctx, payout.ID)
	require.NoError(t, err)
	require.Equal(t, "eligible", payout.EligibilityStatus)
	accepted, err := api.ApproveMarketplacePayouts(ctx, tapsilat.MarketplacePayoutApproval{IDs: []string{payout.ID}})
	require.NoError(t, err)
	require.Equal(t, "accepted", accepted.Status)
	_, err = api.ApproveMarketplacePayouts(ctx, tapsilat.MarketplacePayoutApproval{IDs: []string{payout.ID}})
	require.Error(t, err)
	policies, err := api.ListMarketplaceReleasePolicies(ctx)
	require.NoError(t, err)
	require.Len(t, policies.Rows, 1)
	require.Equal(t, "service_completed", policies.Rows[0].RequiredExternalEvent)
	for _, target := range []struct{ id, provider, reference string }{{payout.ID, "iyzico", "fixture-transaction"}, {a["paytr_payout_id"], "paytr", a["paytr_transaction"]}} {
		legacy, err := api.DisapproveSubmerchantPayment(ctx, tapsilat.SubmerchantPaymentAction{Locale: "en", ConversationID: "withdraw-" + target.provider, PaymentTransactionID: target.reference, Acquirer: target.provider})
		require.NoError(t, err)
		require.Equal(t, "success", legacy.Status)
		require.NotEmpty(t, legacy.OperationID)
		before, err := api.GetMarketplacePayout(ctx, target.id)
		require.NoError(t, err)
		input := tapsilat.MarketplacePayoutOperationInput{Kind: "update_item", IdempotencyKey: "allocation-" + target.provider, Revision: before.Revision, SellerReference: multiple.RoutingReference, NetAmount: "60.00"}
		op, err := api.RunMarketplacePayoutOperation(ctx, target.id, input)
		require.NoError(t, err)
		require.Equal(t, "succeeded", op.Status)
		replay, err := api.RunMarketplacePayoutOperation(ctx, target.id, input)
		require.NoError(t, err)
		require.Equal(t, op.ID, replay.ID)
		current, err := api.GetMarketplacePayout(ctx, target.id)
		require.NoError(t, err)
		require.Equal(t, multiple.ID, current.SubmerchantID)
		require.Equal(t, "60.00", current.Amount)
		require.Equal(t, "awaiting_approval", current.ReleaseStatus)
		require.Equal(t, before.RequiredEvent, current.RequiredEvent)
		legacyUpdate, err := api.UpdateSubmerchantPaymentItem(ctx, tapsilat.SubmerchantPaymentItemUpdate{Locale: "en", ConversationID: "edit-" + target.provider, PaymentTransactionID: target.reference, Acquirer: target.provider, SubMerchantKey: multiple.RoutingReference, SubMerchantPrice: "62.00"})
		require.NoError(t, err)
		require.Equal(t, "success", legacyUpdate.Status)
		require.NotEmpty(t, legacyUpdate.OperationID)
	}
	payout, err = api.GetMarketplacePayout(ctx, payout.ID)
	require.NoError(t, err)
	lost, err := api.RunMarketplacePayoutOperation(ctx, payout.ID, tapsilat.MarketplacePayoutOperationInput{Kind: "update_item", IdempotencyKey: "lost-response", Revision: payout.Revision, SellerReference: multiple.RoutingReference, NetAmount: "61.00"})
	require.NoError(t, err)
	require.Equal(t, "unknown", lost.Status)
	_, err = api.ApproveMarketplacePayouts(ctx, tapsilat.MarketplacePayoutApproval{IDs: []string{payout.ID}})
	require.Error(t, err)
	resolved, err := api.ReconcileMarketplacePayoutOperation(ctx, lost.ID)
	require.NoError(t, err)
	require.Equal(t, "succeeded", resolved.Status)
	readOperation, err := api.GetMarketplacePayoutOperation(ctx, lost.ID)
	require.NoError(t, err)
	require.Equal(t, "succeeded", readOperation.Status)
	payout, err = api.GetMarketplacePayout(ctx, payout.ID)
	require.NoError(t, err)
	require.Equal(t, "61.00", payout.Amount)
	require.Empty(t, payout.ActiveOperationID)
	_, err = tapsilat.NewCustomAPI(p["address"], "invalid").GetMarketplacePayout(ctx, payout.ID)
	require.Error(t, err)
	t.Log("real SDK HTTP -> Panel gRPC -> Admin domain -> gateway fixture: seller lifecycle, two-provider payout materialization, events, gated approval, legacy disapprove/item update, revisioned allocation, replay and lost-response reconciliation passed; provider and auth fixtures are isolated")
}
