package tapsilat

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	StatusCode int
	Code       int
	Message    string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Tapsilat Validation Error\nstatus_code:%d\ncode:%d\nerror:%s", e.StatusCode, e.Code, e.Message)
}

// ValidateInstallments validates and parses installments string. Returns slice of valid installments
func ValidateInstallments(installmentsStr string) ([]int, error) {
	if installmentsStr == "" {
		return []int{1}, nil
	}

	parts := strings.Split(installmentsStr, ",")
	var installments []int

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}

		installment, err := strconv.Atoi(trimmed)
		if err != nil {
			return nil, &ValidationError{
				StatusCode: 400,
				Code:       0,
				Message:    "Enabled installments must be comma-separated integers (e.g., 1,2,3 or 2,4,6)",
			}
		}

		if installment < 1 || installment > 12 {
			return nil, &ValidationError{
				StatusCode: 400,
				Code:       0,
				Message:    fmt.Sprintf("Installment value '%d' is invalid. All installment values must be between 1 and 12 (inclusive).", installment),
			}
		}

		installments = append(installments, installment)
	}

	if len(installments) == 0 {
		return []int{1}, nil
	}

	return installments, nil
}

// ValidateGSMNumber validates GSM number format. Returns the cleaned phone if valid.
func ValidateGSMNumber(phone string) (string, error) {
	if phone == "" {
		return phone, nil
	}

	// Remove formatting characters
	cleanPhone := strings.ReplaceAll(phone, " ", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "-", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, "(", "")
	cleanPhone = strings.ReplaceAll(cleanPhone, ")", "")

	// Check if it contains only valid characters (digits, +, 0)
	digitRegex := regexp.MustCompile(`^[\+0-9]+$`)
	if !digitRegex.MatchString(cleanPhone) {
		return "", &ValidationError{
			StatusCode: 400,
			Code:       0,
			Message:    fmt.Sprintf("Invalid phone number format: %s", phone),
		}
	}

	// Remove + signs for length validation but keep original format for actual content check
	contentForValidation := strings.ReplaceAll(cleanPhone, "+", "")
	if len(contentForValidation) == 0 || !regexp.MustCompile(`^[0-9]+$`).MatchString(contentForValidation) {
		return "", &ValidationError{
			StatusCode: 400,
			Code:       0,
			Message:    fmt.Sprintf("Invalid phone number format: %s", phone),
		}
	}

	// Validate length based on format
	if strings.HasPrefix(cleanPhone, "+") {
		if len(cleanPhone) < 8 {
			return "", &ValidationError{
				StatusCode: 400,
				Code:       0,
				Message:    fmt.Sprintf("International phone number too short: %s", phone),
			}
		}
	} else if strings.HasPrefix(cleanPhone, "00") {
		if len(cleanPhone) < 9 {
			return "", &ValidationError{
				StatusCode: 400,
				Code:       0,
				Message:    fmt.Sprintf("International phone number (00 format) too short: %s", phone),
			}
		}
	} else if strings.HasPrefix(cleanPhone, "0") {
		if len(cleanPhone) < 7 {
			return "", &ValidationError{
				StatusCode: 400,
				Code:       0,
				Message:    fmt.Sprintf("National phone number too short: %s", phone),
			}
		}
	} else {
		if len(cleanPhone) < 6 {
			return "", &ValidationError{
				StatusCode: 400,
				Code:       0,
				Message:    fmt.Sprintf("Local phone number too short: %s", phone),
			}
		}
	}

	return cleanPhone, nil
}
