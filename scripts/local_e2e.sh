#!/usr/bin/env bash
set -euo pipefail

log() {
  printf '[local-e2e] %s\n' "$*"
}

fail() {
  printf '[local-e2e] ERROR: %s\n' "$*" >&2
  exit 1
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || fail "required command not found: $1"
}

usage() {
  cat <<'USAGE'
Run local end-to-end validation for tapsilat-go (submerchant/vpos scope).

Usage:
  scripts/local_e2e.sh [--start-stack] [--smoke-only] [--integration-only] [--bootstrap-submerchant] [--bootstrap-vpos]

Modes:
  --start-stack       Runs `make compose` in panel/backend before tests.
  --smoke-only        Runs only smoke tests.
  --integration-only  Runs only integration tests.
  --bootstrap-submerchant
                     Creates submerchant before tests and auto-resolves IDs.
  --bootstrap-vpos    Creates vpos and vpos-submerchant mapping before tests.

Environment:
  PANEL_API_BASE_URL            Panel API base URL (default: http://localhost:3001/api/v1)
  PANEL_BACKEND_DIR             panel/backend path (default: ../panel/backend)
  START_STACK                   1 to auto run make compose (same as --start-stack)

  # Auth source (one of these must be provided):
  TAPSILAT_API_TOKEN            Existing API token to use directly
  TAPSILAT_API_KEY              API key from localhost:8080 UI
  TAPSILAT_API_SECRET           API secret from localhost:8080 UI

  # Optional token generation settings
  TOKEN_NAME                    default: tapsilat-go-local-e2e
  TOKEN_EXPIRE_DAYS             default: 1

  # Optional test bootstrap
  AUTO_CREATE_SUBMERCHANT       1 to auto-create submerchant before tests
  AUTO_CREATE_VPOS              1 to auto-create vpos + vpos-submerchant mapping before tests
  BOOTSTRAP_CURRENCY_ID         Primary currency UUID for submerchant create (auto-detected from /wallet/currencies if empty)
  BOOTSTRAP_SUBMERCHANT_TYPE    default: PERSONAL
  BOOTSTRAP_SUBMERCHANT_NAME    default: sdk-e2e-<timestamp>
  BOOTSTRAP_SUBMERCHANT_EMAIL   default: sdk-e2e-<timestamp>@example.com
  BOOTSTRAP_SUBMERCHANT_GSM     default: 555000<last4(timestamp)>
  BOOTSTRAP_SUBMERCHANT_KEY     default: sdk-e2e-key-<timestamp>
  BOOTSTRAP_SUBMERCHANT_EXTERNAL_ID default: sdk-e2e-ext-<timestamp>
  BOOTSTRAP_VPOS_NAME           default: sdk-e2e-vpos-<timestamp>
  BOOTSTRAP_VPOS_EXTERNAL_ID    default: sdk-e2e-vpos-ext-<timestamp>
  BOOTSTRAP_VPOS_TERMINAL_NO    default: 12345678
  BOOTSTRAP_OUTPUT_FILE         default: /tmp/tapsilat_local_e2e_ids.env

  # Optional IDs for deeper checks (recommended)
  TAPSILAT_SMOKE_SUBMERCHANT_ID
  TAPSILAT_SMOKE_SUBORGANIZATION_ID
  TAPSILAT_SMOKE_VPOS_ID
  TAPSILAT_IT_SUBMERCHANT_ID
  TAPSILAT_IT_SUBORGANIZATION_ID
  TAPSILAT_IT_VPOS_ID

Notes:
  - If token is generated from key/secret, User-Agent is forced to Go-http-client/1.1
    to match SDK calls and avoid auth context mismatch.
  - Integration tests require TAPSILAT_IT_SUBMERCHANT_ID and
    TAPSILAT_IT_SUBORGANIZATION_ID; otherwise they are skipped by test code.
USAGE
}

START_STACK_FLAG="${START_STACK:-0}"
RUN_SMOKE=1
RUN_INTEGRATION=1
AUTO_BOOTSTRAP="${AUTO_CREATE_SUBMERCHANT:-0}"
AUTO_BOOTSTRAP_VPOS="${AUTO_CREATE_VPOS:-0}"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --start-stack)
      START_STACK_FLAG=1
      shift
      ;;
    --smoke-only)
      RUN_INTEGRATION=0
      shift
      ;;
    --integration-only)
      RUN_SMOKE=0
      shift
      ;;
    --bootstrap-submerchant)
      AUTO_BOOTSTRAP=1
      shift
      ;;
    --bootstrap-vpos)
      AUTO_BOOTSTRAP_VPOS=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      fail "unknown argument: $1"
      ;;
  esac
done

require_cmd curl
require_cmd jq
require_cmd go

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

PANEL_API_BASE_URL="${PANEL_API_BASE_URL:-http://localhost:3001/api/v1}"
PANEL_BACKEND_DIR="${PANEL_BACKEND_DIR:-$REPO_DIR/../panel/backend}"
TOKEN_NAME="${TOKEN_NAME:-tapsilat-go-local-e2e}"
TOKEN_EXPIRE_DAYS="${TOKEN_EXPIRE_DAYS:-1}"

if [[ "$START_STACK_FLAG" == "1" ]]; then
  [[ -d "$PANEL_BACKEND_DIR" ]] || fail "PANEL_BACKEND_DIR not found: $PANEL_BACKEND_DIR"
  require_cmd make
  log "Starting local stack via make compose in $PANEL_BACKEND_DIR"
  (cd "$PANEL_BACKEND_DIR" && make compose)
fi

resolve_token() {
  if [[ -n "${TAPSILAT_API_TOKEN:-}" ]]; then
    printf '%s' "$TAPSILAT_API_TOKEN"
    return 0
  fi

  [[ -n "${TAPSILAT_API_KEY:-}" ]] || fail "set TAPSILAT_API_TOKEN or both TAPSILAT_API_KEY + TAPSILAT_API_SECRET"
  [[ -n "${TAPSILAT_API_SECRET:-}" ]] || fail "set TAPSILAT_API_TOKEN or both TAPSILAT_API_KEY + TAPSILAT_API_SECRET"

  local token_url="$PANEL_API_BASE_URL/token/generate"
  local token_name payload response token scopes_json

  if [[ -n "${TOKEN_SCOPES_JSON:-}" ]]; then
    scopes_json="$TOKEN_SCOPES_JSON"
  else
    scopes_json='[{"entity":"api","create":true,"read":true,"update":true,"delete":true},{"entity":"submerchant","create":true,"read":true,"update":true,"delete":true},{"entity":"organization","create":false,"read":true,"update":false,"delete":false},{"entity":"vpos","create":true,"read":true,"update":true,"delete":true},{"entity":"vpos_submerchant","create":true,"read":true,"update":true,"delete":true},{"entity":"wallet","create":false,"read":true,"update":false,"delete":false}]'
  fi

  token_name="$TOKEN_NAME"
  payload="$(jq -cn --arg name "$token_name" --argjson exp "$TOKEN_EXPIRE_DAYS" --arg scopes_json "$scopes_json" '{name:$name, expire_time:$exp, scopes: ($scopes_json | fromjson)}')"

  log "Generating API token from key/secret via $token_url"
  response="$(curl -sS -X POST "$token_url" \
    -H "Content-Type: application/json" \
    -H "User-Agent: Go-http-client/1.1" \
    -H "x-mono-key: ${TAPSILAT_API_KEY}" \
    -H "x-mono-secret: ${TAPSILAT_API_SECRET}" \
    --data "$payload")"

  token="$(printf '%s' "$response" | jq -r 'if type == "string" then . else (.token // .data.token // empty) end')"

  if [[ -z "$token" ]] && printf '%s' "$response" | jq -e '.error == "TOKEN_CREATE_API_USER_CREATE_NAME_EXIST"' >/dev/null 2>&1; then
    token_name="${TOKEN_NAME}-$(date +%s)"
    payload="$(jq -cn --arg name "$token_name" --argjson exp "$TOKEN_EXPIRE_DAYS" --arg scopes_json "$scopes_json" '{name:$name, expire_time:$exp, scopes: ($scopes_json | fromjson)}')"
    response="$(curl -sS -X POST "$token_url" \
      -H "Content-Type: application/json" \
      -H "User-Agent: Go-http-client/1.1" \
      -H "x-mono-key: ${TAPSILAT_API_KEY}" \
      -H "x-mono-secret: ${TAPSILAT_API_SECRET}" \
      --data "$payload")"
    token="$(printf '%s' "$response" | jq -r 'if type == "string" then . else (.token // .data.token // empty) end')"
  fi

  [[ -n "$token" ]] || fail "token generation failed. response: $response"

  printf '%s' "$token"
}

API_TOKEN="$(resolve_token)"
log "Token resolved (length=${#API_TOKEN})"

api_get() {
  local path="$1"
  curl -sS -X GET "${PANEL_API_BASE_URL}${path}" \
    -H "Authorization: Bearer ${API_TOKEN}" \
    -H "User-Agent: Go-http-client/1.1"
}

api_post() {
  local path="$1"
  local body="$2"
  curl -sS -X POST "${PANEL_API_BASE_URL}${path}" \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer ${API_TOKEN}" \
    -H "User-Agent: Go-http-client/1.1" \
    --data "$body"
}

write_bootstrap_output() {
  local output_file
  output_file="${BOOTSTRAP_OUTPUT_FILE:-/tmp/tapsilat_local_e2e_ids.env}"

  {
    printf 'BOOTSTRAP_CREATED_SUBMERCHANT_ID=%q\n' "${BOOTSTRAP_CREATED_SUBMERCHANT_ID:-}"
    printf 'BOOTSTRAP_CREATED_SUBORGANIZATION_ID=%q\n' "${BOOTSTRAP_CREATED_SUBORGANIZATION_ID:-}"
    printf 'BOOTSTRAP_SUBMERCHANT_KEY=%q\n' "${BOOTSTRAP_SUBMERCHANT_KEY:-}"
    printf 'BOOTSTRAP_CREATED_VPOS_ID=%q\n' "${BOOTSTRAP_CREATED_VPOS_ID:-}"
    printf 'BOOTSTRAP_CREATED_VPOS_SUBMERCHANT_ID=%q\n' "${BOOTSTRAP_CREATED_VPOS_SUBMERCHANT_ID:-}"
    printf 'BOOTSTRAP_VPOS_EXTERNAL_ID=%q\n' "${BOOTSTRAP_VPOS_EXTERNAL_ID:-}"
    printf 'TAPSILAT_SMOKE_SUBMERCHANT_ID=%q\n' "${TAPSILAT_SMOKE_SUBMERCHANT_ID:-}"
    printf 'TAPSILAT_SMOKE_SUBORGANIZATION_ID=%q\n' "${TAPSILAT_SMOKE_SUBORGANIZATION_ID:-}"
    printf 'TAPSILAT_SMOKE_VPOS_ID=%q\n' "${TAPSILAT_SMOKE_VPOS_ID:-}"
    printf 'TAPSILAT_IT_SUBMERCHANT_ID=%q\n' "${TAPSILAT_IT_SUBMERCHANT_ID:-}"
    printf 'TAPSILAT_IT_SUBORGANIZATION_ID=%q\n' "${TAPSILAT_IT_SUBORGANIZATION_ID:-}"
    printf 'TAPSILAT_IT_VPOS_ID=%q\n' "${TAPSILAT_IT_VPOS_ID:-}"
  } > "$output_file"

  log "Bootstrap IDs written to ${output_file}"
}

pick_currency_id() {
  if [[ -n "${BOOTSTRAP_CURRENCY_ID:-}" ]]; then
    printf '%s' "$BOOTSTRAP_CURRENCY_ID"
    return 0
  fi

  local currencies
  currencies="$(api_get '/wallet/currencies?page=1&per_page=10')"

  local currency_id
  currency_id="$(printf '%s' "$currencies" | jq -r '(.rows // .row // .data // []) | .[0] | .id // .currency_id // .currencyId // empty')"
  [[ -n "$currency_id" ]] || fail "could not resolve currency id from /wallet/currencies. set BOOTSTRAP_CURRENCY_ID"

  printf '%s' "$currency_id"
}

pick_currency_ids_json() {
  if [[ -n "${BOOTSTRAP_CURRENCY_IDS_JSON:-}" ]]; then
    printf '%s' "$BOOTSTRAP_CURRENCY_IDS_JSON"
    return 0
  fi

  local currencies
  currencies="$(api_get '/wallet/currencies?page=1&per_page=10')"

  local currency_ids_json
  currency_ids_json="$(printf '%s' "$currencies" | jq -c '(.rows // .row // .data // []) | map(.id // .currency_id // .currencyId // empty) | map(select(. != "")) | unique | .[:2]')"
  [[ "$currency_ids_json" != "[]" ]] || fail "could not resolve currency ids from /wallet/currencies. set BOOTSTRAP_CURRENCY_IDS_JSON"

  printf '%s' "$currency_ids_json"
}

bootstrap_submerchant_ids() {
  local now suffix default_name default_email default_gsm
  now="$(date +%s)"
  suffix="${now: -4}"

  default_name="sdk-e2e-${now}"
  default_email="sdk-e2e-${now}@example.com"
  default_gsm="555000${suffix}"

  local sm_type sm_name sm_email sm_gsm sm_key sm_ext sm_conversation sm_identity currency_id
  sm_type="${BOOTSTRAP_SUBMERCHANT_TYPE:-PERSONAL}"
  sm_name="${BOOTSTRAP_SUBMERCHANT_NAME:-$default_name}"
  sm_email="${BOOTSTRAP_SUBMERCHANT_EMAIL:-$default_email}"
  sm_gsm="${BOOTSTRAP_SUBMERCHANT_GSM:-$default_gsm}"
  sm_key="${BOOTSTRAP_SUBMERCHANT_KEY:-sdk-e2e-key-${now}}"
  sm_ext="${BOOTSTRAP_SUBMERCHANT_EXTERNAL_ID:-sdk-e2e-ext-${now}}"
  sm_conversation="${BOOTSTRAP_CONVERSATION_ID:-sdk-e2e-conv-${now}}"
  sm_identity="${BOOTSTRAP_IDENTITY_NUMBER:-1000000${suffix}}"
  currency_id="$(pick_currency_id)"

  log "Bootstrapping submerchant fixture"

  local create_payload create_resp
  create_payload="$(jq -cn \
    --arg locale "tr" \
    --arg conversation_id "$sm_conversation" \
    --arg name "$sm_name" \
    --arg email "$sm_email" \
    --arg gsm_number "$sm_gsm" \
    --arg address "Istanbul" \
    --arg iban "TR000000000000000000000000" \
    --arg tax_office "Merter" \
    --arg legal_company_title "$sm_name" \
    --arg currency_id "$currency_id" \
    --arg sub_merchant_external_id "$sm_ext" \
    --arg identity_number "$sm_identity" \
    --arg sub_merchant_type "$sm_type" \
    --arg tax_number "1234567890" \
    --arg sub_merchant_key "$sm_key" \
    --arg organization_id "" \
    --arg status "active" \
    --argjson system_time "$now" \
    --arg contact_name "SDK" \
    --arg contact_surname "E2E" \
    '{
      locale: $locale,
      conversation_id: $conversation_id,
      name: $name,
      email: $email,
      gsm_number: $gsm_number,
      address: $address,
      iban: $iban,
      tax_office: $tax_office,
      legal_company_title: $legal_company_title,
      currency_id: $currency_id,
      sub_merchant_external_id: $sub_merchant_external_id,
      identity_number: $identity_number,
      sub_merchant_type: $sub_merchant_type,
      tax_number: $tax_number,
      sub_merchant_key: $sub_merchant_key,
      organization_id: $organization_id,
      status: $status,
      system_time: $system_time,
      contact_name: $contact_name,
      contact_surname: $contact_surname
    }')"

  create_resp="$(api_post '/submerchants' "$create_payload")"

  if [[ "$(printf '%s' "$create_resp" | jq -r '.error // empty')" != "" ]]; then
    fail "submerchant create returned error: $create_resp"
  fi

  local list_resp submerchant_id
  list_resp="$(api_get '/submerchants?page=1&per_page=100')"
  submerchant_id="$(printf '%s' "$list_resp" | jq -r --arg key "$sm_key" '(.rows // .row // []) | map(select((.submerchant_key // .sub_merchant_key // .subMerchantKey // "") == $key)) | .[0].id // empty')"
  [[ -n "$submerchant_id" ]] || fail "created submerchant id not found by key=$sm_key"

  local suborganization_id mapping_resp

  suborganization_id="$(printf '%s' "$create_resp" | jq -r '.suborganization_id // .sub_organization_id // .suborganizationId // empty')"

  if [[ -z "$suborganization_id" ]]; then
    for _ in $(seq 1 10); do
      mapping_resp="$(api_get "/submerchants/${submerchant_id}/suborganization")"
      suborganization_id="$(printf '%s' "$mapping_resp" | jq -r '.suborganization_id // .sub_organization_id // .suborganizationId // empty')"
      if [[ -n "$suborganization_id" ]]; then
        break
      fi
      sleep 1
    done
  fi

  if [[ -z "$suborganization_id" ]]; then
    local suborg_list
    suborg_list="$(api_get '/organization/suborganizations?page=1&per_page=100')"
    suborganization_id="$(printf '%s' "$suborg_list" | jq -r --arg name "$sm_name" '(.rows // .row // []) | map(select((.name // "") == $name)) | .[0].id // empty')"
  fi

  if [[ -n "$suborganization_id" ]]; then
    local reverse_resp reverse_submerchant_id
    reverse_resp="$(api_get "/organization/suborganizations/${suborganization_id}/submerchant")"
    reverse_submerchant_id="$(printf '%s' "$reverse_resp" | jq -r '.submerchant_id // .sub_merchant_id // .submerchantId // empty')"
    if [[ "$reverse_submerchant_id" != "$submerchant_id" ]]; then
      log "Warning: suborganization mapping not yet consistent (suborg=${suborganization_id}, reverse_submerchant=${reverse_submerchant_id:-<empty>}); skipping suborganization-bound assertions"
      suborganization_id=""
    fi
  else
    log "Warning: suborganization id could not be resolved for submerchant=${submerchant_id}; continuing with submerchant-only assertions"
  fi

  local smoke_submerchant_id smoke_suborganization_id it_submerchant_id it_suborganization_id
  smoke_submerchant_id="$submerchant_id"
  smoke_suborganization_id="$suborganization_id"
  it_submerchant_id="$submerchant_id"
  it_suborganization_id="$suborganization_id"

  if [[ -z "$suborganization_id" ]]; then
    smoke_submerchant_id=""
    it_submerchant_id=""
    it_suborganization_id=""
  fi

  export TAPSILAT_SMOKE_SUBMERCHANT_ID="$smoke_submerchant_id"
  export TAPSILAT_SMOKE_SUBORGANIZATION_ID="$smoke_suborganization_id"
  export TAPSILAT_IT_SUBMERCHANT_ID="$it_submerchant_id"
  export TAPSILAT_IT_SUBORGANIZATION_ID="$it_suborganization_id"

  export BOOTSTRAP_CREATED_SUBMERCHANT_ID="$submerchant_id"
  export BOOTSTRAP_CREATED_SUBORGANIZATION_ID="$suborganization_id"
  export BOOTSTRAP_SUBMERCHANT_KEY="$sm_key"

  log "Bootstrap created submerchant_id=${submerchant_id} suborganization_id=${suborganization_id:-<not-resolved>}"
  write_bootstrap_output
}

bootstrap_vpos_ids() {
  local submerchant_id
  submerchant_id="${TAPSILAT_SMOKE_SUBMERCHANT_ID:-${TAPSILAT_IT_SUBMERCHANT_ID:-${BOOTSTRAP_CREATED_SUBMERCHANT_ID:-}}}"
  [[ -n "$submerchant_id" ]] || fail "vpos bootstrap requires submerchant id. use --bootstrap-submerchant or set TAPSILAT_SMOKE_SUBMERCHANT_ID/TAPSILAT_IT_SUBMERCHANT_ID"

  local now vpos_name vpos_external_id terminal_no acquirer_id currency_id currency_ids_json
  now="$(date +%s)"
  vpos_name="${BOOTSTRAP_VPOS_NAME:-sdk-e2e-vpos-${now}}"
  vpos_external_id="${BOOTSTRAP_VPOS_EXTERNAL_ID:-sdk-e2e-vpos-ext-${now}}"
  terminal_no="${BOOTSTRAP_VPOS_TERMINAL_NO:-12345678}"

  local acq_resp
  acq_resp="$(api_get '/vpos/acquirers')"
  acquirer_id="$(printf '%s' "$acq_resp" | jq -r '.items[0].id // empty')"
  [[ -n "$acquirer_id" ]] || fail "could not resolve acquirer_id from /vpos/acquirers"

  currency_id="$(pick_currency_id)"
  currency_ids_json="$(pick_currency_ids_json)"

  log "Bootstrapping vpos fixture"

  local vpos_create_payload vpos_create_resp
  vpos_create_payload="$(jq -cn \
    --arg name "$vpos_name" \
    --arg env_mode "test" \
    --arg payment_mode "api" \
    --arg acquirer_id "$acquirer_id" \
    --argjson currency_ids "$currency_ids_json" \
    '{
      name: $name,
      env_mode: $env_mode,
      payment_mode: $payment_mode,
      acquirer_id: $acquirer_id,
      currencies: $currency_ids
    }')"

  vpos_create_resp="$(api_post '/vpos' "$vpos_create_payload")"
  if [[ "$(printf '%s' "$vpos_create_resp" | jq -r '.error // empty')" != "" ]]; then
    fail "vpos create returned error: $vpos_create_resp"
  fi

  local vpos_list_resp vpos_id
  vpos_list_resp="$(api_get '/vpos?page=1&per_page=100')"
  vpos_id="$(printf '%s' "$vpos_list_resp" | jq -r --arg n "$vpos_name" '(.rows // .row // []) | map(select((.name // "") == $n)) | .[0].id // empty')"
  [[ -n "$vpos_id" ]] || fail "created vpos id not found by name=$vpos_name"

  local attach_payload attach_resp
  attach_payload="$(jq -cn \
    --arg external_reference_id "$vpos_external_id" \
    --arg submerchant_id "$submerchant_id" \
    --arg terminal_no "$terminal_no" \
    --arg vpos_id "$vpos_id" \
    --arg mcc "5812" \
    --arg tax_id "1234567890" \
    --arg national_id "12345678901" \
    --arg title "SDK E2E" \
    --arg switch_id "SDK-E2E" \
    --arg city "Istanbul" \
    --arg country "Turkey" \
    --arg country_isocode "TR" \
    --arg postal_code "34000" \
    --arg address "Istanbul" \
    --arg submerchant_url "https://example.com" \
    --arg submerchant_nin "12345678901" \
    '{
      external_reference_id: $external_reference_id,
      submerchant_id: $submerchant_id,
      terminal_no: $terminal_no,
      vpos_id: $vpos_id,
      mcc: $mcc,
      tax_id: $tax_id,
      national_id: $national_id,
      title: $title,
      switch_id: $switch_id,
      city: $city,
      country: $country,
      country_isocode: $country_isocode,
      postal_code: $postal_code,
      address: $address,
      submerchant_url: $submerchant_url,
      submerchant_nin: $submerchant_nin
    }')"

  attach_resp="$(api_post '/vpos-submerchant' "$attach_payload")"
  if [[ "$(printf '%s' "$attach_resp" | jq -r '.error // empty')" != "" ]]; then
    fail "vpos-submerchant create returned error: $attach_resp"
  fi

  local mapping_list_resp vpos_submerchant_id mapped_submerchant_id
  mapping_list_resp="$(api_get "/vpos-submerchant?page=1&per_page=100&vpos_id=${vpos_id}")"
  vpos_submerchant_id="$(printf '%s' "$mapping_list_resp" | jq -r --arg ext "$vpos_external_id" '(.rows // .row // []) | map(select((.external_reference_id // "") == $ext)) | .[0].id // empty')"
  mapped_submerchant_id="$(printf '%s' "$mapping_list_resp" | jq -r --arg ext "$vpos_external_id" '(.rows // .row // []) | map(select((.external_reference_id // "") == $ext)) | .[0].submerchant_id // empty')"

  [[ -n "$vpos_submerchant_id" ]] || fail "created vpos-submerchant id not found by external_reference_id=$vpos_external_id"
  [[ "$mapped_submerchant_id" == "$submerchant_id" ]] || fail "vpos-submerchant mapped submerchant mismatch: expected=$submerchant_id actual=${mapped_submerchant_id:-<empty>}"

  local mapping_read_resp mapping_read_vpos_id
  mapping_read_resp="$(api_get "/vpos-submerchant/${vpos_submerchant_id}")"
  mapping_read_vpos_id="$(printf '%s' "$mapping_read_resp" | jq -r '.vpos_id // .vposId // empty')"
  [[ "$mapping_read_vpos_id" == "$vpos_id" ]] || fail "vpos-submerchant read verification failed: expected_vpos=$vpos_id actual=${mapping_read_vpos_id:-<empty>}"

  export BOOTSTRAP_CREATED_VPOS_ID="$vpos_id"
  export BOOTSTRAP_CREATED_VPOS_SUBMERCHANT_ID="$vpos_submerchant_id"
  export BOOTSTRAP_VPOS_EXTERNAL_ID="$vpos_external_id"
  export TAPSILAT_SMOKE_VPOS_ID="$vpos_id"
  export TAPSILAT_IT_VPOS_ID="$vpos_id"

  log "Bootstrap created vpos_id=${vpos_id} vpos_submerchant_id=${vpos_submerchant_id}"
  write_bootstrap_output
}

export TAPSILAT_SMOKE_ENDPOINT="$PANEL_API_BASE_URL"
export TAPSILAT_SMOKE_TOKEN="$API_TOKEN"

export TAPSILAT_IT_ENDPOINT="$PANEL_API_BASE_URL"
export TAPSILAT_IT_TOKEN="$API_TOKEN"

# Forward optional IDs when present
for key in \
  TAPSILAT_SMOKE_SUBMERCHANT_ID \
  TAPSILAT_SMOKE_SUBORGANIZATION_ID \
  TAPSILAT_SMOKE_VPOS_ID \
  TAPSILAT_IT_SUBMERCHANT_ID \
  TAPSILAT_IT_SUBORGANIZATION_ID \
  TAPSILAT_IT_VPOS_ID; do
  export "$key"="${!key:-}"
done

if [[ "$AUTO_BOOTSTRAP" == "1" ]]; then
  bootstrap_submerchant_ids
fi

if [[ "$AUTO_BOOTSTRAP_VPOS" == "1" ]]; then
  bootstrap_vpos_ids
fi

cd "$REPO_DIR"

if [[ "$RUN_SMOKE" == "1" ]]; then
  log "Running smoke tests"
  go test -v ./tests/smoke -run TestSmokeReadAndListFlows
fi

if [[ "$RUN_INTEGRATION" == "1" ]]; then
  log "Running integration tests"
  go test -v ./tests/integration -run TestBackendConsistency_SubmerchantSuborganizationAndScopedVpos
fi

log "Local E2E flow completed"
