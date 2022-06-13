package aform

import (
	"fmt"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"golang.org/x/text/language"
)

// WithAttributes returns a FieldOption that adds custom attributes to the
// Widget. For all the attributes but the class attribute, values set with this
// function override the values set by any other logic. If the class attribute
// is set, the error class set with the function SetErrorCSSClass is
// concatenated at the beginning. Three attribute names panic: type, name and
// value. To change the type, use a different Field and/or set a different
// Widget on an existing Field. To change the name attribute, set a different
// name to the Field when you create it. To change the value attribute, set an
// initial value when you create the field or bind values with BindRequest or
// BindData functions.
func WithAttributes(attrs []Attributable) FieldOption {
	return func(fld *Field) error {
		fld.SetAttributes(attrs)
		return nil
	}
}

// WithWidget returns a FieldOption that changes the Widget of the Field.
func WithWidget(widget Widget) FieldOption {
	return func(fld *Field) error {
		fld.SetWidget(widget)
		return nil
	}
}

// SetLabelSuffix sets a string appended to the label. To set the same suffix
// to all fields in a form use WithLabelSuffix.
func (fld *Field) SetLabelSuffix(labelSuffix string) {
	fld.labelSuffix = labelSuffix
}

// CustomizeError replaces one built-in Error with err, if err
// ErrorCoderTranslator.Code matches one of the existing Error code. Existing
// Error codes are BooleanErrorCode, EmailErrorCode, ChoiceErrorCode,
// MinLengthErrorCode, MaxLengthErrorCode, RequiredErrorCode and URLErrorCode.
// If err ErrorCoderTranslator.Code is not from this list, it panics.
func (fld *Field) CustomizeError(err ErrorCoderTranslator) {
	e := errorWrapIfNotAsError(err)
	if !slices.Contains(customizableErrors, e.Code()) {
		panic(fmt.Sprintf("CustomizeError called on %s field with a nonexisting Error code: %s", fld.name, e.Code()))
	}
	if fld.customErrors == nil {
		fld.customErrors = map[string]Error{}
	}
	fld.customErrors[e.Code()] = e
}

// SetAutoID sets the string used to build the HTML ID. See WithAutoID for
// details on valid values.
func (fld *Field) SetAutoID(autoID string) error {
	if err := validAutoID(autoID); err != nil {
		return err
	}
	fld.autoID = autoID
	return nil
}

// SetRequiredCSSClass sets a CSS class added to the HTML if the field
// is required. To set the same required class to all fields in a form
// use WithRequiredCSSClass.
func (fld *Field) SetRequiredCSSClass(class string) {
	fld.requiredCSSClass = class
}

// SetErrorCSSClass sets a CSS class added to the HTML when the field has a
// validation error. To set the same error class to all fields in a form use
// WithErrorCSSClass.
func (fld *Field) SetErrorCSSClass(class string) {
	fld.errorCSSClass = class
}

// SetAttributes adds custom attributes to the Widget.
// See WithAttributes for details.
func (fld *Field) SetAttributes(attrs []Attributable) {
	a := attributableListToAttrs(attrs)
	if fld.attrs == nil {
		fld.attrs = a
		return
	}
	maps.Copy(fld.attrs, a)
}

// SetWidget changes the widget to the field. See WithWidget for
// examples.
func (fld *Field) SetWidget(widget Widget) {
	fld.widget = widget
}

// SetLocale changes the locale used by the field to translate error
// messages. To set the same locale to all fields in a form use
// WithLocales.
func (fld *Field) SetLocale(locale language.Tag) {
	fld.locale = locale
}

// SanitizationFunc defines a function to sanitize a Field.
type SanitizationFunc func(string) string

// SetSanitizeFunc sets sanitization function. Parameter update is a
// function that itself has current sanitization function as parameter and must
// return the new sanitization function. Default sanitization function removes
// HTML elements, leading and trailing white characters. It removes as well new
// lines if the widget is not TextArea.
func (fld *Field) SetSanitizeFunc(update func(current SanitizationFunc) (new SanitizationFunc)) {
	fld.sanitizeFunc = update(fld.currentSanitizeFunc())
}

func (fld *Field) currentSanitizeFunc() SanitizationFunc {
	if fld.sanitizeFunc != nil {
		return fld.sanitizeFunc
	}
	return fld.widget.defaultSanitizeFunc()
}

func (fld *Field) sanitize(value string) string {
	return fld.currentSanitizeFunc()(value)
}

// ValidationFunc defines a function to validate a Field.
type ValidationFunc func(string, bool) []Error

// SetValidateFunc sets validation function. Parameter update is a
// function that itself has current validation function as parameter and must
// return the new validation function. Default validation function depends on
// the field type.
func (fld *Field) SetValidateFunc(update func(current ValidationFunc) (new ValidationFunc)) {
	fld.validateFunc = update(fld.validateFunc)
}
