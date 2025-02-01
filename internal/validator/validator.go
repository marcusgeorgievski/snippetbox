package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// True if FieldErrors doesn't contan any errors
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// Adds an error message as long as no entry already exists for the given key
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// Adds error only if a validation check is not 'ok'
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// True if value is not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// True if value contains no more than max chars
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// True if a value is in a list of specific permitted values
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
