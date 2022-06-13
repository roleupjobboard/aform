package aform

import (
	"fmt"
)

// IsValid returns true if all the fields' validation return no error.
// It does Form validation if it is not already done.
func (f *Form) IsValid() bool {
	if !f.IsBound() {
		return false
	}
	f.doValidationIfNeeded()
	return len(f.Errors()) == 0
}

// CleanedData returns cleaned data validated from Form inputs. It does Form
// validation if it is not already done. If a field validation returns an
// error it doesn't appear in CleanedData.
func (f *Form) CleanedData() CleanedData {
	if !f.IsBound() {
		return map[string][]string{}
	}
	f.doValidationIfNeeded()
	return f.cleanedData
}

// Errors returns errors happened during validation of form inputs. It does Form
// validation if it is not already done. If a field input is valid it doesn't
// appear in FormErrors.
func (f *Form) Errors() FormErrors {
	if !f.IsBound() {
		return map[string][]Error{}
	}
	f.doValidationIfNeeded()
	output := map[string][]Error{}
	for fld, errorList := range f.errors {
		outputErrorList := make([]Error, len(errorList))
		for i, err := range errorList {
			outputErrorList[i] = err
		}
		output[fld] = outputErrorList
	}
	return output
}

// SetCleanFunc sets a clean function to do validation at the form level.
func (f *Form) SetCleanFunc(clean func(*Form)) {
	f.cleanFunc = clean
}

// AddError adds an error to the field named field. An error can be added after
// the form has been validated. AddError can be used in the form clean function
// set with SetCleanFunc to perform Form level validation.
// If err implements ErrorCoderTranslator, error message will be translated
// according to the language automatically identified by BindRequest or
// according to the language given to BindData.
func (f *Form) AddError(field string, fieldErr error) error {
	if !f.validated {
		return fmt.Errorf("you can't add an error to a form not already validated. " +
			"A form is validated when one of the following method is called: " +
			"CleanedData(), IsValid() or Errors()")
	}
	fld, err := f.internalFieldByName(field)
	if err != nil {
		return err
	}
	nName := normalizedNameForField(fld)
	wrappedFieldErr := errorWrapIfNotAsError(fieldErr)
	fld.addError(wrappedFieldErr)
	f.errors[nName] = append(f.errors[nName], wrappedFieldErr)
	delete(f.cleanedData, nName)
	return nil
}

func (f *Form) doValidationIfNeeded() {
	if f.validated {
		return
	}
	f.validated = true
	cleanedData := map[string][]string{}
	errors := map[string][]Error{}
	for _, fld := range f.fields {
		nName := normalizedNameForField(fld)
		values, ok := f.boundData[nName]
		if !ok {
			values = []string{}
		}
		validationFld := disguiseFieldForValidation(fld)
		cleanValues, errs := validationFld.Clean(values)
		if errs != nil {
			errors[nName] = errs
		} else {
			cleanedData[nName] = cleanValues
		}
	}
	f.cleanedData = cleanedData
	f.errors = errors
	f.cleanFunc(f)
}

func disguiseFieldForValidation(fld fieldInterface) multipleValueValidationStateProvider {
	switch fld.Type() {
	case BooleanFieldType, CharFieldType, EmailFieldType, URLFieldType, ChoiceFieldType:
		return singleValueDisguisedInMultipleValueValidationStateProvider{p: fld.(singleValueValidationStateProvider)}
	case MultipleChoiceFieldType:
		return fld.(multipleValueValidationStateProvider)
	default:
		panic(fmt.Sprintf("unknown field type %s", fld.Type()))
	}
}

type singleValueValidationStateProvider interface {
	Required() bool
	Clean(value string) (string, []Error)
	EmptyValue() string
}

type multipleValueValidationStateProvider interface {
	Required() bool
	Clean(values []string) ([]string, []Error)
	EmptyValue() []string
}

type singleValueDisguisedInMultipleValueValidationStateProvider struct {
	p singleValueValidationStateProvider
}

func (p singleValueDisguisedInMultipleValueValidationStateProvider) Required() bool {
	return p.p.Required()
}

func (p singleValueDisguisedInMultipleValueValidationStateProvider) Clean(values []string) ([]string, []Error) {
	v := ""
	if len(values) > 0 {
		v = values[0]
	}
	cleanValue, errs := p.p.Clean(v)
	if errs != nil {
		return []string{}, errs
	} else {
		return []string{cleanValue}, nil
	}
}

func (p singleValueDisguisedInMultipleValueValidationStateProvider) EmptyValue() []string {
	return []string{p.p.EmptyValue()}
}
