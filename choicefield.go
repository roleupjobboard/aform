package aform

// ChoiceField is a field type that validates that the given value is one of
// the options listed with WithChoiceOptions or WithGroupedChoiceOptions.
type ChoiceField struct {
	initialValue string
	*Field
}

// verify interface compliance
var _ fieldInterface = (*ChoiceField)(nil)

// NewChoiceField creates a choice field named name. The parameter initial is
// the initial value before data bounding. To add options use WithChoiceOptions
// or WithGroupedChoiceOptions. The default Widget is Select. To change it, use
// WithWidget or Field.SetWidget.
func NewChoiceField(name, initial string, opts ...FieldOption) (*ChoiceField, error) {
	cf := &ChoiceField{
		initial,
		&Field{
			name:         name,
			boundValues:  []string{initial},
			errors:       []Error{},
			optionGroups: []choiceFieldOptionGroup{},
			fieldType:    ChoiceFieldType,
			widget:       Select,
			autoID:       defaultAutoID,
			label:        name,
			labelSuffix:  defaultLabelSuffix,
			locale:       defaultLanguage,
		},
	}
	cf.Field.validateFunc = choiceFieldValidationFunc(cf)
	for _, opt := range opts {
		if err := opt(cf.Field); err != nil {
			return nil, err
		}
	}
	return cf, nil
}

// DefaultChoiceField creates a choice field with reasonable default values.
// initial parameter is the empty string. To add options use WithChoiceOptions
// or WithGroupedChoiceOptions.
func DefaultChoiceField(name string, opts ...FieldOption) (*ChoiceField, error) {
	return NewChoiceField(name, "", opts...)
}

func (fld *ChoiceField) field() *Field {
	return fld.Field
}

// Clean returns the cleaned value. value is first sanitized and
// finally validated. Sanitization can be customized with
// Field.SetSanitizeFunc. Validation can be customized with
// Field.SetValidateFunc.
func (fld *ChoiceField) Clean(value string) (string, []Error) {
	fld.boundValues = []string{value}
	sanitizedValue := fld.sanitize(value)
	if fld.notRequired && len(sanitizedValue) == 0 {
		return fld.EmptyValue(), nil
	}
	fld.errors = customizeErrors(fld.validateFunc(sanitizedValue, !fld.notRequired), fld.customErrors)
	return sanitizedValue, fld.errors
}

// EmptyValue returns the ChoiceField empty value. The empty value is the
// cleaned value returned by Clean when there is no data bound to the field.
// A ChoiceField empty value is always the empty string "".
func (fld *ChoiceField) EmptyValue() string {
	return ""
}

func choiceFieldValidationFunc(fld *ChoiceField) func(string, bool) []Error {
	return func(value string, required bool) []Error {
		notRequiredAndEmpty := !required && len(value) == 0
		if notRequiredAndEmpty {
			return []Error{}
		}
		var choices []string
		for _, group := range fld.optionGroups {
			choices = append(choices, group.values()...)
		}
		return validateValue(value, buildValidationRules(required, buildValidationChoicesRule(choices...)))
	}
}
