package aform_test

import (
	"errors"
	"fmt"
	"github.com/roleupjobboard/aform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorWrap(t *testing.T) {
	a := assert.New(t)
	wrappedError := fmt.Errorf("wrapped message")
	err := aform.ErrorWrap(wrappedError)
	a.ErrorIs(err, wrappedError)
	a.Equal("wrapped message", err.Error())
	a.Equal("", err.Code())
	a.Equal("wrapped message", err.Translate("en"))
	a.Equal("wrapped message", err.Translate("fr"))
}

type testDummyTranslatorError struct {
	source string
}

func (d testDummyTranslatorError) Error() string {
	return "Error() string"
}

func (d testDummyTranslatorError) Code() string {
	return "Code() string"
}

func (d testDummyTranslatorError) Translate(locale string) string {
	switch locale {
	case "en":
		return fmt.Sprintf("Hello from %s", d.source)
	case "fr":
		return fmt.Sprintf("Bonjour de %s", d.source)
	default:
		return fmt.Sprintf("Bye %s", d.source)
	}
}

func TestErrorWrap_withTaggerTranslator(t *testing.T) {
	a := assert.New(t)
	err := aform.ErrorWrap(testDummyTranslatorError{source: "Aby"})
	a.Equal("Error() string", err.Error())
	a.Equal("Code() string", err.Code())
	a.Equal("Hello from Aby", err.Translate("en"))
	a.Equal("Bonjour de Aby", err.Translate("fr"))
	var dummy testDummyTranslatorError
	a.ErrorAs(err, &dummy)
}

func TestErrorWrapWithCode(t *testing.T) {
	a := assert.New(t)
	wrappedError := fmt.Errorf("wrapped message")
	err := aform.ErrorWrapWithCode(wrappedError, "testtag")
	a.ErrorIs(err, wrappedError)
	a.Equal("wrapped message", err.Error())
	a.Equal("testtag", err.Code())
	a.Equal("wrapped message", err.Translate("en"))
	a.Equal("wrapped message", err.Translate("fr"))
}

func TestErrorWrapWithCode_withTaggerTranslator(t *testing.T) {
	a := assert.New(t)
	err := aform.ErrorWrapWithCode(testDummyTranslatorError{source: "Aby"}, "testtag")
	a.Equal("Error() string", err.Error())
	a.Equal("testtag", err.Code())
	a.Equal("Hello from Aby", err.Translate("en"))
	a.Equal("Bonjour de Aby", err.Translate("fr"))
	var dummy testDummyTranslatorError
	a.ErrorAs(err, &dummy)
}

func TestError_Unwrap(t *testing.T) {
	a := assert.New(t)
	err := aform.ErrorWrapWithCode(testDummyTranslatorError{source: "Aby"}, "testtag")
	a.Equal("Error() string", err.Error())
	a.Equal("testtag", err.Code())
	var wrappedErr testDummyTranslatorError
	a.True(errors.As(err, &wrappedErr))
	a.Equal("Code() string", wrappedErr.Code())
}

func TestErrorWrap_requiredError(t *testing.T) {
	a := assert.New(t)
	err := aform.ErrorWrap(aform.ExportRequiredError)
	a.Equal(aform.RequiredErrorCode, err.Code())
	a.Equal(aform.RequiredErrorMessageEn, err.Error())
	a.Equal(aform.RequiredErrorMessageEn, err.Translate("en"))
	a.Equal(aform.RequiredErrorMessageFr, err.Translate("fr"))
}
