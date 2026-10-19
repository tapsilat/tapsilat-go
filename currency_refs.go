package tapsilat

import (
	"context"
	"fmt"
	"maps"
	"regexp"
	"strings"
)

var uuidRefRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// ResolveCurrencyID returns the organization currency UUID for ref. A UUID is returned unchanged;
// a unit such as "TRY" is looked up in GetOrganizationCurrencies, whose result is cached.
// Most endpoints accept units directly, so this is only needed where a UUID is required.
func (t *API) ResolveCurrencyID(ctx context.Context, ref string) (string, error) {
	return t.normalizeCurrencyID(ctx, ref)
}

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

// resolveOptionalCurrencyID resolves ref when it is set and leaves an empty value unchanged.
func (t *API) resolveOptionalCurrencyID(ctx context.Context, ref string) (string, error) {
	if strings.TrimSpace(ref) == "" {
		return ref, nil
	}
	return t.normalizeCurrencyID(ctx, ref)
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
	maps.Copy(cloned, values)
	return cloned
}

func (t *API) invalidateCurrencyCache() {
	t.currencyRefsMu.Lock()
	t.currencyIDsByUnit = nil
	t.currencyCacheReady = false
	t.currencyRefsMu.Unlock()
}
