package aform

import (
	"fmt"
)

// BooleanField is a field type that validates that the given value is a valid
// boolean.
type BooleanField struct {
	initialValue string
	*Field
}

// verify interface compliance
var _ fieldInterface = (*BooleanField)(nil)

// NewBooleanField creates a boolean field named name. The parameter initial is
// the initial value before data bounding. The default Widget is CheckboxInput.
// To change it, use WithWidget or SetWidget.
func NewBooleanField(name string, initial bool, opts ...FieldOption) (*BooleanField, error) {
	cf := &BooleanField{
		boolToValue(initial),
		&Field{
			name:         name,
			boundValues:  []string{boolToValue(initial)},
			errors:       []Error{},
			fieldType:    BooleanFieldType,
			widget:       CheckboxInput,
			autoID:       defaultAutoID,
			label:        name,
			labelSuffix:  defaultLabelSuffix,
			validateFunc: booleanFieldValidation,
			locale:       defaultLanguage,
		},
	}
	for _, opt := range opts {
		if err := opt(cf.Field); err != nil {
			return nil, err
		}
	}
	return cf, nil
}

// DefaultBooleanField creates a boolean field with reasonable default values.
// The initial parameter value is false.
func DefaultBooleanField(name string, opts ...FieldOption) (*BooleanField, error) {
	return NewBooleanField(name, false, opts...)
}

func (fld *BooleanField) field() *Field {
	return fld.Field
}

// Clean returns the cleaned value. value is first sanitized and
// finally validated. Sanitization can be customized with
// Field.SetSanitizeFunc. Validation can be customized with
// Field.SetValidateFunc.
func (fld *BooleanField) Clean(value string) (string, []Error) {
	fld.boundValues = []string{value}
	sanitizedValue := fld.sanitize(value)
	if fld.notRequired && len(sanitizedValue) == 0 {
		return fld.EmptyValue(), nil
	}
	fld.errors = customizeErrors(fld.validateFunc(sanitizedValue, !fld.notRequired), fld.customErrors)
	return sanitizedValue, fld.errors
}

// EmptyValue returns the BooleanField empty value. The empty value is the
// cleaned value returned by Clean when there is no data bound to the field.
// A BooleanField empty value is always "off".
func (fld *BooleanField) EmptyValue() string {
	return boolToValue(false)
}

// MustBoolean returns the clean value type cast to bool. It panics
// if the value provided is not a valid boolean input.
func (fld *BooleanField) MustBoolean(value string) bool {
	v, errs := fld.Clean(value)
	if len(errs) > 0 {
		panic(fmt.Sprintf("MustBoolean called on %s field with an invalid boolean value: %s", fld.name, v))
	}
	return valueToBool(v)
}

func booleanFieldValidation(value string, required bool) []Error {
	return validateValue(value, buildValidationRules(required, BooleanErrorCode))
}

func boolToValue(b bool) string {
	if b {
		return "on"
	}
	return "off"
}

func parseBool(v string) (bool, error) {
	switch v {
	case "1", "t", "T", "true", "TRUE", "True", "on", "ON", "On":
		return true, nil
	case "0", "f", "F", "false", "FALSE", "False", "off", "OFF", "Off":
		return false, nil
	}
	return false, fmt.Errorf("not a bool %s", v)
}

func valueToBool(v string) bool {
	b, err := parseBool(v)
	if err != nil {
		return false
	}
	return b
}
