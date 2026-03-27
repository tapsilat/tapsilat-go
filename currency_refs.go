package tapsilat

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var uuidRefRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89aAbB][0-9a-fA-F]{3}-[0-9a-fA-F]{12}$`)

func (t *API) normalizeCurrencyID(ctx context.Context, ref string) (string, error) {
	trimmedRef := strings.TrimSpace(ref)
	if trimmedRef == "" {
		return "", &ValidationError{
			StatusCode: 400,
			Code:       0,
			Message:    "currency_id is required",
		}
	}

	if uuidRefRegex.MatchString(trimmedRef) {
		return trimmedRef, nil
	}

	currencyIDsByUnit, err := t.getOrganizationCurrencyIDsByUnit(ctx)
	if err != nil {
		return "", err
	}

	if currencyID, ok := currencyIDsByUnit[strings.ToUpper(trimmedRef)]; ok {
		return currencyID, nil
	}

	return "", &ValidationError{
		StatusCode: 400,
		Code:       0,
		Message:    fmt.Sprintf("unknown currency reference %q; use a currency UUID or organization currency unit", ref),
	}
}

func (t *API) normalizeCurrencyIDs(ctx context.Context, refs []string) ([]string, error) {
	if len(refs) == 0 {
		return nil, &ValidationError{
			StatusCode: 400,
			Code:       0,
			Message:    "currencies must contain at least one currency reference",
		}
	}

	normalized := make([]string, 0, len(refs))
	seen := make(map[string]struct{}, len(refs))

	for _, ref := range refs {
		currencyID, err := t.normalizeCurrencyID(ctx, ref)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[currencyID]; ok {
			continue
		}
		seen[currencyID] = struct{}{}
		normalized = append(normalized, currencyID)
	}

	if len(normalized) == 0 {
		return nil, &ValidationError{
			StatusCode: 400,
			Code:       0,
			Message:    "currencies must contain at least one currency reference",
		}
	}

	return normalized, nil
}

func (t *API) getOrganizationCurrencyIDsByUnit(ctx context.Context) (map[string]string, error) {
	t.currencyRefsMu.RLock()
	if t.currencyCacheReady {
		cached := cloneStringMap(t.currencyIDsByUnit)
		t.currencyRefsMu.RUnlock()
		return cached, nil
	}
	t.currencyRefsMu.RUnlock()

	response, err := t.GetOrganizationCurrencies(ctx)
	if err != nil {
		return nil, err
	}

	resolved := make(map[string]string, len(response.Currencies))
	for _, currency := range response.Currencies {
		unit := strings.ToUpper(strings.TrimSpace(currency.CurrencyUnit))
		if unit == "" || strings.TrimSpace(currency.ID) == "" {
			continue
		}
		resolved[unit] = currency.ID
	}

	t.currencyRefsMu.Lock()
	t.currencyIDsByUnit = cloneStringMap(resolved)
	t.currencyCacheReady = true
	t.currencyRefsMu.Unlock()

	return resolved, nil
}

func cloneStringMap(values map[string]string) map[string]string {
	if len(values) == 0 {
		return map[string]string{}
	}

	cloned := make(map[string]string, len(values))
	for key, value := range values {
		cloned[key] = value
	}
	return cloned
}
