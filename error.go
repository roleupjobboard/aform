package aform

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"golang.org/x/text/language"
)

// Error codes of the available validations. Each validation has a code and a
// default message. Example: If a BooleanField value is not a boolean, the
// validation returns an error with the code BooleanErrorCode and the default
// message BooleanErrorMessageEn. Error messages can be modified with the
// function Field.CustomizeError.
const (
	BooleanErrorCode   = "boolean"
	EmailErrorCode     = "email"
	ChoiceErrorCode    = "oneof"
	MinLengthErrorCode = "min"
	MaxLengthErrorCode = "max"
	RequiredErrorCode  = "required"
	URLErrorCode       = "url"
)

var customizableErrors = []string{BooleanErrorCode, EmailErrorCode, ChoiceErrorCode, MinLengthErrorCode, MaxLengthErrorCode, RequiredErrorCode, URLErrorCode}

// ErrorCoderTranslator defines the validation errors interface.
type ErrorCoderTranslator interface {
	error
	Code() string
	Translate(locale string) string
}

// Error defines errors returned by the forms validations. It implements
// ErrorCoderTranslator to help the customization of existing errors or
// the addition of application specific errors. To customize existing
// validation error messages use Field.CustomizeError. To add application
// specific errors at the form level use SetCleanFunc and at the field
// level Field.SetValidateFunc.
type Error struct {
	tag string
	err error
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) Error() string {
	if e.err == nil {
		return ""
	}
	return e.err.Error()
}

type errorCoder interface {
	error
	Code() string
}

// Code returns the code of the error.
func (e Error) Code() string {
	var ec errorCoder
	if e.tag != "" {
		return e.tag
	}
	if errors.As(e.err, &ec) {
		return ec.Code()
	}
	return ""
}

type errorTranslator interface {
	error
	Translate(locale string) string
}

// Translate translates error messages. locale is a BCP 47 language tag.
func (e Error) Translate(locale string) string {
	if e.err == nil {
		return e.Error()
	}
	var ec errorTranslator
	if errors.As(e.err, &ec) {
		return ec.Translate(locale)
	}
	return e.Error()
}

// ErrorWrap wraps an error in a validation Error. It makes possible to use any
// error as a validation Error.
func ErrorWrap(err error) Error {
	return Error{err: err}
}

// ErrorWrapWithCode is like ErrorWrap and set or update the Error code.
func ErrorWrapWithCode(err error, code string) Error {
	return Error{tag: code, err: err}
}

// errorWrapIfNotAsError wraps err with ErrorWrap if there is no error of type
// Error in err's chain. If there is one, it returns this error.
func errorWrapIfNotAsError(err error) Error {
	var e Error
	if errors.As(err, &e) {
		return e
	} else {
		return ErrorWrap(err)
	}
}

// customizeErrors customizes errors of a field with the error messages set
// with Field.CustomizeError.
func customizeErrors(errors []Error, customized map[string]Error) []Error {
	if len(errors) == 0 {
		return errors
	}
	updated := make([]Error, len(errors))
	for i, e := range errors {
		ce, ok := customized[e.Code()]
		if ok {
			updated[i] = ce
		} else {
			updated[i] = e
		}
	}
	return updated
}

type errorFromFieldError struct {
	fieldError validator.FieldError
}

func (e errorFromFieldError) Error() string {
	return e.Translate(defaultLanguage.String())
}

func (e errorFromFieldError) Code() string {
	return e.fieldError.Tag()
}

func (e errorFromFieldError) Translate(locale string) string {
	trans, _ := universalTranslator.GetTranslator(locale)
	return e.fieldError.Translate(trans)
}

type simpleError struct {
	code string
	fr   string
	en   string
}

func (e simpleError) Code() string {
	return e.code
}

func (e simpleError) Error() string {
	return RequiredErrorMessageEn
}

func (e simpleError) Translate(locale string) string {
	switch locale {
	case language.French.String():
		return e.fr
	default:
		return e.en
	}
}

var (
	requiredError = ErrorWrap(simpleError{code: RequiredErrorCode, fr: RequiredErrorMessageFr, en: RequiredErrorMessageEn})
)
