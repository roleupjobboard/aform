package aform

// MultipleChoiceField is a field type that validates that the given values
// are contained in options listed with WithChoiceOptions or
// WithGroupedChoiceOptions.
type MultipleChoiceField struct {
	initialValues []string
	*Field
}

// verify interface compliance
var _ fieldInterface = (*MultipleChoiceField)(nil)

// NewMultipleChoiceField creates a choice field named name. It allows multiple
// values to be selected. The parameter initials is the list of initial values
// before data bounding. To add options use WithChoiceOptions
// or WithGroupedChoiceOptions. The default Widget is SelectMultiple. To change
// it, use WithWidget or Field.SetWidget.
func NewMultipleChoiceField(name string, initials []string, opts ...FieldOption) (*MultipleChoiceField, error) {
	cf := &MultipleChoiceField{
		initials,
		&Field{
			name:         name,
			boundValues:  initials,
			errors:       []Error{},
			optionGroups: []choiceFieldOptionGroup{},
			fieldType:    MultipleChoiceFieldType,
			widget:       SelectMultiple,
			autoID:       defaultAutoID,
			label:        name,
			labelSuffix:  defaultLabelSuffix,
			locale:       defaultLanguage,
		},
	}
	cf.Field.validateFunc = multipleChoiceFieldValidationFunc(cf)
	for _, opt := range opts {
		if err := opt(cf.Field); err != nil {
			return nil, err
		}
	}
	return cf, nil
}

// DefaultMultipleChoiceField creates a choice field with reasonable default
// values. initial parameter is an empty slice. To add options use
// WithChoiceOptions or WithGroupedChoiceOptions.
func DefaultMultipleChoiceField(name string, opts ...FieldOption) (*MultipleChoiceField, error) {
	return NewMultipleChoiceField(name, []string{}, opts...)
}

func (fld *MultipleChoiceField) field() *Field {
	return fld.Field
}

// Clean returns the slice of cleaned value. values are first sanitized and
// finally validated. Sanitization can be customized with
// Field.SetSanitizeFunc. Validation can be customized with
// Field.SetValidateFunc.
func (fld *MultipleChoiceField) Clean(values []string) ([]string, []Error) {
	fld.boundValues = values
	sanitizedValues := make([]string, len(values))
	for i, value := range values {
		sanitizedValues[i] = fld.sanitize(value)
	}
	if multipleChoiceCleanShouldReturnEmptyValue(sanitizedValues, !fld.notRequired) {
		return fld.EmptyValue(), nil
	}
	if !multipleChoiceIsRequiredValidation(sanitizedValues, !fld.notRequired) {
		fld.errors = []Error{requiredError}
		return sanitizedValues, fld.errors
	}
	var allErrors []Error
	for _, sanitizedValue := range sanitizedValues {
		errors := fld.validateFunc(sanitizedValue, false)
		if len(errors) > 0 {
			allErrors = append(allErrors, customizeErrors(errors, fld.customErrors)...)
		}
	}
	fld.errors = allErrors
	return sanitizedValues, fld.errors
}

// EmptyValue returns the MultipleChoiceField empty value. The empty value is the
// cleaned value returned by Clean when there is no data bound to the field.
// A MultipleChoiceField empty value is always an empty slice.
func (fld *MultipleChoiceField) EmptyValue() []string {
	return []string{}
}

// multipleChoiceFieldValidationFunc validates one bound value of a multiple choice field.
// A multiple choice field can have multiple values bound.
// The required constraint must be applied on the list of bound values not on each bound value.
// Consequently, the required parameter is not used in the function returned.
func multipleChoiceFieldValidationFunc(fld *MultipleChoiceField) func(string, bool) []Error {
	return func(value string, _ bool) []Error {
		if len(value) == 0 {
			return []Error{}
		}
		var choices []string
		for _, group := range fld.optionGroups {
			choices = append(choices, group.values()...)
		}
		return validateValue(value, buildValidationRules(false, buildValidationChoicesRule(choices...)))
	}
}

func multipleChoiceCleanShouldReturnEmptyValue(values []string, required bool) bool {
	noValue := true
	for _, value := range values {
		if len(value) > 0 {
			noValue = false
		}
	}
	return noValue && !required
}

// multipleChoiceIsRequiredValidation returns true if the values pass the validation, false otherwise.
func multipleChoiceIsRequiredValidation(values []string, required bool) bool {
	if !required {
		return true
	}
	for _, value := range values {
		if len(value) > 0 {
			return true
		}
	}
	return false
}
