package unit_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tapsilat "github.com/tapsilat/tapsilat-go"
)

func TestValidateInstallments(t *testing.T) {
	t.Run("EmptyStringReturnsDefault", func(t *testing.T) {
		result, err := tapsilat.ValidateInstallments("")
		require.NoError(t, err)
		assert.Equal(t, []int{1}, result)
	})

	t.Run("ValidSingleInstallment", func(t *testing.T) {
		result, err := tapsilat.ValidateInstallments("3")
		require.NoError(t, err)
		assert.Equal(t, []int{3}, result)
	})

	t.Run("ValidMultipleInstallments", func(t *testing.T) {
		result, err := tapsilat.ValidateInstallments("1,2,3,6")
		require.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3, 6}, result)
	})

	t.Run("ValidWithSpaces", func(t *testing.T) {
		result, err := tapsilat.ValidateInstallments("1, 2, 3, 6")
		require.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3, 6}, result)
	})

	t.Run("InstallmentValueTooLow", func(t *testing.T) {
		_, err := tapsilat.ValidateInstallments("0,2,3")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Installment value '0' is invalid")
	})

	t.Run("InstallmentValueTooHigh", func(t *testing.T) {
		_, err := tapsilat.ValidateInstallments("1,15,3")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Installment value '15' is invalid")
	})

	t.Run("InvalidFormatLetters", func(t *testing.T) {
		_, err := tapsilat.ValidateInstallments("1,abc,3")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "comma-separated integers")
	})

	t.Run("InvalidFormatMixed", func(t *testing.T) {
		_, err := tapsilat.ValidateInstallments("1,2.5,3")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "comma-separated integers")
	})
}

func TestValidateGSMNumber(t *testing.T) {
	t.Run("EmptyStringReturnsEmpty", func(t *testing.T) {
		result, err := tapsilat.ValidateGSMNumber("")
		require.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("ValidInternationalPlusFormat", func(t *testing.T) {
		result, err := tapsilat.ValidateGSMNumber("+905551234567")
		require.NoError(t, err)
		assert.Equal(t, "+905551234567", result)
	})

	t.Run("ValidInternational00Format", func(t *testing.T) {
		result, err := tapsilat.ValidateGSMNumber("00905551234567")
		require.NoError(t, err)
		assert.Equal(t, "00905551234567", result)
	})

	t.Run("ValidNationalFormat", func(t *testing.T) {
		result, err := tapsilat.ValidateGSMNumber("05551234567")
		require.NoError(t, err)
		assert.Equal(t, "05551234567", result)
	})

	t.Run("ValidLocalFormat", func(t *testing.T) {
		result, err := tapsilat.ValidateGSMNumber("5551234567")
		require.NoError(t, err)
		assert.Equal(t, "5551234567", result)
	})

	t.Run("RemovesFormattingCharacters", func(t *testing.T) {
		result, err := tapsilat.ValidateGSMNumber("+90 555 123-45(67)")
		require.NoError(t, err)
		assert.Equal(t, "+905551234567", result)
	})

	t.Run("InternationalPlusTooShort", func(t *testing.T) {
		_, err := tapsilat.ValidateGSMNumber("+90123")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "too short")
	})

	t.Run("International00TooShort", func(t *testing.T) {
		_, err := tapsilat.ValidateGSMNumber("0090123")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "too short")
	})

	t.Run("NationalTooShort", func(t *testing.T) {
		_, err := tapsilat.ValidateGSMNumber("012345")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "too short")
	})

	t.Run("LocalTooShort", func(t *testing.T) {
		_, err := tapsilat.ValidateGSMNumber("12345")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "too short")
	})

	t.Run("InvalidCharacters", func(t *testing.T) {
		_, err := tapsilat.ValidateGSMNumber("+90abc1234567")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid phone number format")
	})

	t.Run("OnlySpecialCharacters", func(t *testing.T) {
		_, err := tapsilat.ValidateGSMNumber("+++---")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "Invalid phone number format")
	})
}
