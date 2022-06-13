package aform

import (
	"fmt"
)

// EmailField is a field type that validates that the given value is a valid
// email address.
type EmailField struct {
	initialValue string
	emptyValue   string
	*Field
}

// verify interface compliance
var _ fieldInterface = (*EmailField)(nil)

// NewEmailField creates an email field named name. The parameter initial is the
// initial value before data bounding. The parameter empty is the cleaned data
// value when there is no data bound to the field. If the parameter min
// (respectively max) is not 0, it validates that the input value is longer or
// equal than min (respectively shorter or equal than max). The default Widget
// is EmailInput. To change it, use WithWidget or SetWidget.
func NewEmailField(name, initial, empty string, min, max uint, opts ...FieldOption) (*EmailField, error) {
	cf := &EmailField{
		initial,
		empty,
		&Field{
			name:        name,
			boundValues: []string{initial},
			errors:      []Error{},
			fieldType:   EmailFieldType,
			widget:      EmailInput,
			autoID:      defaultAutoID,
			label:       name,
			labelSuffix: defaultLabelSuffix,
			minLength:   min,
			maxLength:   max,
			locale:      defaultLanguage,
		},
	}
	cf.Field.validateFunc = emailFieldValidationFunc(cf)
	for _, opt := range opts {
		if err := opt(cf.Field); err != nil {
			return nil, err
		}
	}
	return cf, nil
}

// defaultEmailMaxLength is the default max length allowed for emails
// https://cheatsheetseries.owasp.org/cheatsheets/Input_Validation_Cheat_Sheet.html#syntactic-validation
const defaultEmailMaxLength = 254

// DefaultEmailField creates an email field with reasonable default values.
// initial and empty parameters are the empty string. min length is 0
// and max length is 254.
func DefaultEmailField(name string, opts ...FieldOption) (*EmailField, error) {
	return NewEmailField(name, "", "", 0, defaultEmailMaxLength, opts...)
}

func (fld *EmailField) field() *Field {
	return fld.Field
}

// Clean returns the cleaned value. value is first sanitized and
// finally validated. Sanitization can be customized with
// Field.SetSanitizeFunc. Validation can be customized with
// Field.SetValidateFunc.
func (fld *EmailField) Clean(value string) (string, []Error) {
	fld.boundValues = []string{value}
	sanitizedValue := fld.sanitize(value)
	if fld.notRequired && len(sanitizedValue) == 0 {
		return fld.EmptyValue(), nil
	}
	fld.errors = customizeErrors(fld.validateFunc(sanitizedValue, !fld.notRequired), fld.customErrors)
	return sanitizedValue, fld.errors
}

// EmptyValue returns the EmailField empty value. The empty value is the
// cleaned value returned by Clean when there is no data bound to the field.
// To set a custom empty value use NewEmailField.
func (fld *EmailField) EmptyValue() string {
	return fld.emptyValue
}

// MustEmail returns the clean value if the value provided is a valid email.
// Otherwise, it panics.
func (fld *EmailField) MustEmail(value string) string {
	v, errs := fld.Clean(value)
	if len(errs) > 0 {
		panic(fmt.Sprintf("MustEmail called on %s field with an invalid email value: %s", fld.name, v))

	}
	return v
}

func emailFieldValidationFunc(fld *EmailField) func(string, bool) []Error {
	return func(value string, required bool) []Error {
		var rules []string
		if fld.minLength > 0 {
			rules = append(rules, buildValidationMinRule(fld.minLength))
		}
		if fld.maxLength > 0 {
			rules = append(rules, buildValidationMaxRule(fld.maxLength))
		}
		rules = append(rules, EmailErrorCode)
		return validateValue(value, buildValidationRules(required, rules...))
	}
}
