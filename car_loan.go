package tapsilat

import "context"

type CarLoanVehicle struct {
	State           int32   `json:"state,omitempty"`
	FirmID          string  `json:"firm_id,omitempty"`
	TsrsbCode       string  `json:"tsrsb_code,omitempty"`
	Model           string  `json:"model,omitempty"`
	Year            int32   `json:"year,omitempty"`
	InvoiceAmount   float64 `json:"invoice_amount,omitempty"`
	ChassisNumber   string  `json:"chassis_number,omitempty"`
	EngineNumber    string  `json:"engine_number,omitempty"`
	InsuranceAmount float64 `json:"insurance_amount,omitempty"`
	Plate           string  `json:"plate,omitempty"`
	RegistrationID  string  `json:"registration_id,omitempty"`
}

type CarLoanProvider struct {
	ProviderID               string  `json:"provider_id"`
	ProviderName             string  `json:"provider_name"`
	ProviderType             int32   `json:"provider_type"`
	CredentialID             string  `json:"credential_id"`
	Category                 string  `json:"category"`
	IsActive                 bool    `json:"is_active"`
	SortOrder                int32   `json:"sort_order"`
	AllowedInstallmentCounts []int32 `json:"allowed_installment_counts"`
}

type CarLoanProvidersResponse struct {
	ResponseCode    string            `json:"response_code"`
	ResponseMessage string            `json:"response_message"`
	OrganizationID  string            `json:"organization_id"`
	TraceID         string            `json:"trace_id"`
	Providers       []CarLoanProvider `json:"providers"`
}

type CarLoanCheckLimitRequest struct {
	ProviderType               int32           `json:"provider_type,omitempty"`
	ReferenceID                string          `json:"reference_id,omitempty"`
	Tckn                       string          `json:"tckn,omitempty"`
	CustomerNumber             string          `json:"customer_number,omitempty"`
	PhoneNumber                string          `json:"phone_number,omitempty"`
	Dob                        string          `json:"dob,omitempty"`
	RequestedAmount            float64         `json:"requested_amount,omitempty"`
	MaxTerm                    int32           `json:"max_term,omitempty"`
	MonthlyIncome              float64         `json:"monthly_income,omitempty"`
	IncomeType                 string          `json:"income_type,omitempty"`
	Education                  string          `json:"education,omitempty"`
	Occupation                 string          `json:"occupation,omitempty"`
	SellerID                   string          `json:"seller_id,omitempty"`
	EmployeeID                 string          `json:"employee_id,omitempty"`
	OrderID                    string          `json:"order_id,omitempty"`
	Vehicle                    *CarLoanVehicle `json:"vehicle,omitempty"`
	FindeksScore               *int32          `json:"findeks_score,omitempty"`
	FindeksHasDelay            *bool           `json:"findeks_has_delay,omitempty"`
	FindeksHasLegalFollowup    *bool           `json:"findeks_has_legal_followup,omitempty"`
	ExistingMonthlyObligations *float64        `json:"existing_monthly_obligations,omitempty"`
	EmploymentDurationMonths   *int32          `json:"employment_duration_months,omitempty"`
	ResidenceType              *string         `json:"residence_type,omitempty"`
	MaritalStatus              *string         `json:"marital_status,omitempty"`
	DependentCount             *int32          `json:"dependent_count,omitempty"`
	TraceID                    string          `json:"trace_id,omitempty"`
}

type CarLoanOffer struct {
	ID                          string  `json:"id"`
	ProviderID                  string  `json:"provider_id"`
	ProviderName                string  `json:"provider_name"`
	ProviderType                int32   `json:"provider_type"`
	CredentialID                string  `json:"credential_id"`
	SessionID                   string  `json:"session_id"`
	ClientSessionID             string  `json:"client_session_id"`
	BankName                    string  `json:"bank_name"`
	AcquirerID                  string  `json:"acquirer_id"`
	ApplicationID               string  `json:"application_id"`
	AvailableLimit              float64 `json:"available_limit"`
	MinAmount                   float64 `json:"min_amount"`
	MaxTerm                     int32   `json:"max_term"`
	MinInterestRate             float64 `json:"min_interest_rate"`
	MaxInterestRate             float64 `json:"max_interest_rate"`
	InstallmentAmountForMaxTerm float64 `json:"installment_amount_for_max_term"`
	LoanToValueRatio            float64 `json:"loan_to_value_ratio"`
	LogoPath                    string  `json:"logo_path"`
	IsNewCustomer               bool    `json:"is_new_customer"`
	HasError                    bool    `json:"has_error"`
	ErrorMessage                string  `json:"error_message"`
}

type CarLoanBank struct {
	ProviderID      string         `json:"provider_id"`
	ProviderName    string         `json:"provider_name"`
	ProviderType    int32          `json:"provider_type"`
	CredentialID    string         `json:"credential_id"`
	BankName        string         `json:"bank_name"`
	ResponseCode    string         `json:"response_code"`
	ResponseMessage string         `json:"response_message"`
	HasError        bool           `json:"has_error"`
	RejectCode      string         `json:"reject_code"`
	Offers          []CarLoanOffer `json:"offers"`
}

type CarLoanCheckLimitResponse struct {
	ResponseCode    string         `json:"response_code"`
	ResponseMessage string         `json:"response_message"`
	OrganizationID  string         `json:"organization_id"`
	ReferenceID     string         `json:"reference_id"`
	TraceID         string         `json:"trace_id"`
	VehicleValue    float64        `json:"vehicle_value"`
	Banks           []CarLoanBank  `json:"banks"`
	Offers          []CarLoanOffer `json:"offers"`
}

type CarLoanPaymentPlanRequest struct {
	CredentialID    string   `json:"credential_id"`
	OfferID         string   `json:"offer_id"`
	Term            int32    `json:"term"`
	RequestedAmount *float64 `json:"requested_amount,omitempty"`
	InstallmentDate string   `json:"installment_date,omitempty"`
}

type CarLoanPaymentPlanItem struct {
	PaymentPlanID               string  `json:"payment_plan_id"`
	Term                        int64   `json:"term,string"`
	InstallmentAmount           float64 `json:"installment_amount"`
	InterestRate                float64 `json:"interest_rate"`
	TotalPaymentAmount          float64 `json:"total_payment_amount"`
	AnnualEffectiveInterestRate float64 `json:"annual_effective_interest_rate"`
	FirstInstallmentPaymentDate string  `json:"first_installment_payment_date"`
	LoanAmount                  float64 `json:"loan_amount"`
	CreditInsuranceFee          float64 `json:"credit_insurance_fee"`
	CreditUsageFee              float64 `json:"credit_usage_fee"`
}

type CarLoanPaymentPlanResponse struct {
	ResponseCode    string                   `json:"response_code"`
	ResponseMessage string                   `json:"response_message"`
	CredentialID    string                   `json:"credential_id"`
	AcquirerID      string                   `json:"acquirer_id"`
	AcquirerName    string                   `json:"acquirer_name"`
	LogoURL         string                   `json:"logo_url"`
	LoanAmount      float64                  `json:"loan_amount"`
	Plans           []CarLoanPaymentPlanItem `json:"plans"`
}

type CarLoanApplyRequest struct {
	CredentialID    string          `json:"credential_id"`
	OfferID         string          `json:"offer_id"`
	Term            int32           `json:"term"`
	RequestedAmount *float64        `json:"requested_amount,omitempty"`
	PaymentPlanID   string          `json:"payment_plan_id,omitempty"`
	CallbackURL     string          `json:"callback_url,omitempty"`
	SellerID        string          `json:"seller_id,omitempty"`
	EmployeeID      string          `json:"employee_id,omitempty"`
	OrderID         string          `json:"order_id,omitempty"`
	Vehicle         *CarLoanVehicle `json:"vehicle,omitempty"`
}

type CarLoanApplyResponse struct {
	ResponseCode    string            `json:"response_code"`
	ResponseMessage string            `json:"response_message"`
	ApplicationID   string            `json:"application_id"`
	ApplicationDate string            `json:"application_date"`
	RedirectURL     string            `json:"redirect_url"`
	IsApproved      bool              `json:"is_approved"`
	CredentialID    string            `json:"credential_id"`
	BankName        string            `json:"bank_name"`
	FormData        map[string]string `json:"form_data"`
}

type CarLoanStatusRequest struct {
	CredentialID    string `json:"credential_id"`
	SessionID       string `json:"session_id,omitempty"`
	ClientSessionID string `json:"client_session_id,omitempty"`
	ApplicationID   string `json:"application_id"`
	ReferenceID     string `json:"reference_id,omitempty"`
	Tckn            string `json:"tckn,omitempty"`
}

type CarLoanStatusResponse struct {
	ResponseCode    string                  `json:"response_code"`
	ResponseMessage string                  `json:"response_message"`
	ApplicationID   string                  `json:"application_id"`
	Status          string                  `json:"status"`
	IsApproved      bool                    `json:"is_approved"`
	PaymentPlan     *CarLoanPaymentPlanItem `json:"payment_plan"`
	CredentialID    string                  `json:"credential_id"`
}

func (t *API) CarLoanProviders(ctx context.Context) (CarLoanProvidersResponse, error) {
	var response CarLoanProvidersResponse
	err := t.post(ctx, "/loan/car/providers", struct{}{}, &response)
	return response, err
}

func (t *API) CarLoanCheckLimit(ctx context.Context, payload CarLoanCheckLimitRequest) (CarLoanCheckLimitResponse, error) {
	var response CarLoanCheckLimitResponse
	err := t.post(ctx, "/loan/car/limit", payload, &response)
	return response, err
}

func (t *API) CarLoanPaymentPlan(ctx context.Context, payload CarLoanPaymentPlanRequest) (CarLoanPaymentPlanResponse, error) {
	var response CarLoanPaymentPlanResponse
	err := t.post(ctx, "/loan/car/payment-plan", payload, &response)
	return response, err
}

func (t *API) CarLoanApply(ctx context.Context, payload CarLoanApplyRequest) (CarLoanApplyResponse, error) {
	var response CarLoanApplyResponse
	err := t.post(ctx, "/loan/car/apply", payload, &response)
	return response, err
}

func (t *API) CarLoanStatus(ctx context.Context, payload CarLoanStatusRequest) (CarLoanStatusResponse, error) {
	var response CarLoanStatusResponse
	err := t.post(ctx, "/loan/car/status", payload, &response)
	return response, err
}
