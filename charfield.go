package aform

// CharField is a field type that doesn't do semantic validation on the given
// value.
type CharField struct {
	initialValue string
	emptyValue   string
	*Field
}

// verify interface compliance
var _ fieldInterface = (*CharField)(nil)

// NewCharField creates a char field named name. The parameter initial is the
// initial value before data bounding. The parameter empty is the cleaned data
// value when there is no data bound to the field. If the parameter min
// (respectively max) is not 0, it validates that the input value is longer or
// equal than min (respectively shorter or equal than max). Otherwise, all
// inputs are valid. The default Widget is TextInput. To change it, use
// WithWidget or SetWidget.
func NewCharField(name, initial, empty string, min, max uint, opts ...FieldOption) (*CharField, error) {
	cf := &CharField{
		initial,
		empty,
		&Field{
			name:        name,
			boundValues: []string{initial},
			errors:      []Error{},
			fieldType:   CharFieldType,
			widget:      TextInput,
			autoID:      defaultAutoID,
			label:       name,
			labelSuffix: defaultLabelSuffix,
			minLength:   min,
			maxLength:   max,
			locale:      defaultLanguage,
		},
	}
	cf.Field.validateFunc = charFieldValidationFunc(cf)
	for _, opt := range opts {
		if err := opt(cf.Field); err != nil {
			return nil, err
		}
	}
	return cf, nil
}

// DefaultCharField creates a char field with reasonable default values.
// initial and empty parameters are the empty string. min length is 0
// and max length is 256.
func DefaultCharField(name string, opts ...FieldOption) (*CharField, error) {
	return NewCharField(name, "", "", 0, 256, opts...)
}

func (fld *CharField) field() *Field {
	return fld.Field
}

// Clean returns the cleaned value. value is first sanitized and
// finally validated. Sanitization can be customized with
// Field.SetSanitizeFunc. Validation can be customized with
// Field.SetValidateFunc.
func (fld *CharField) Clean(value string) (string, []Error) {
	fld.boundValues = []string{value}
	sanitizedValue := fld.sanitize(value)
	if fld.notRequired && len(sanitizedValue) == 0 {
		return fld.EmptyValue(), nil
	}
	fld.errors = customizeErrors(fld.validateFunc(sanitizedValue, !fld.notRequired), fld.customErrors)
	return sanitizedValue, fld.errors
}

// EmptyValue returns the CharField empty value. The empty value is the
// cleaned value returned by Clean when there is no data bound to the field.
// To set a custom empty value use NewCharField.
func (fld *CharField) EmptyValue() string {
	return fld.emptyValue
}

func charFieldValidationFunc(fld *CharField) func(string, bool) []Error {
	return func(value string, required bool) []Error {
		var rules []string
		if fld.minLength > 0 {
			rules = append(rules, buildValidationMinRule(fld.minLength))
		}
		if fld.maxLength > 0 {
			rules = append(rules, buildValidationMaxRule(fld.maxLength))
		}
		return validateValue(value, buildValidationRules(required, rules...))
	}
}
