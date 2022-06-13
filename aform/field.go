package aform

import (
	"golang.org/x/text/language"
)

const (
	defaultLabelSuffix = ""
)

// FieldType defines the type of the fields.
type FieldType string

// Field types available. Each of them has its own creation interface.
// For instance to create a BooleanFieldType, you can use NewBooleanField or
// DefaultBooleanField functions.
const (
	BooleanFieldType        = FieldType("BooleanField")
	CharFieldType           = FieldType("CharField")
	EmailFieldType          = FieldType("EmailField")
	URLFieldType            = FieldType("URLField")
	ChoiceFieldType         = FieldType("ChoiceField")
	MultipleChoiceFieldType = FieldType("MultipleChoiceField")
)

// FieldOption describes a functional option for configuring a Field.
type FieldOption func(*Field) error

// WithLabel returns a FieldOption that overrides the default label of the Field.
func WithLabel(label string) FieldOption {
	return func(fld *Field) error {
		fld.SetLabel(label)
		return nil
	}
}

// IsSafe returns a FieldOption that sets the Field as HTML safe. When a field
// is marked as safe, its label and label suffix are not HTML-escaped.
func IsSafe() FieldOption {
	return func(fld *Field) error {
		fld.MarkSafe()
		return nil
	}
}

// WithHelpText returns a FieldOption that adds a help text to the Field.
// Help text is not HTML-escaped.
func WithHelpText(help string) FieldOption {
	return func(fld *Field) error {
		fld.SetHelpText(help)
		return nil
	}
}

// IsNotRequired returns a FieldOption that sets the Field as not required.
// By default, all fields are required.
func IsNotRequired() FieldOption {
	return func(fld *Field) error {
		fld.SetNotRequired()
		return nil
	}
}

// IsDisabled returns a FieldOption that sets the Field as disabled.
func IsDisabled() FieldOption {
	return func(fld *Field) error {
		fld.SetDisabled()
		return nil
	}
}

// WithChoiceOptions returns a FieldOption that adds a list of
// ChoiceFieldOption to the Field. This option is used only by fields
// presenting a list of choices like ChoiceField and MultipleChoiceField.
func WithChoiceOptions(options []ChoiceFieldOption) FieldOption {
	return func(fld *Field) error {
		fld.AddChoiceOptions("", options)
		return nil
	}
}

// WithGroupedChoiceOptions returns a FieldOption that adds a list of
// ChoiceFieldOption grouped together in a <optgroup> tag. Parameter label is
// the value of the <optgroup> attribute label. This option is used only by
// fields presenting a list of choices like ChoiceField and MultipleChoiceField.
func WithGroupedChoiceOptions(label string, options []ChoiceFieldOption) FieldOption {
	return func(fld *Field) error {
		fld.AddChoiceOptions(label, options)
		return nil
	}
}

// Field is the type grouping features shared by all field types.
type Field struct {
	name             string
	boundValues      []string
	errors           []Error
	optionGroups     []choiceFieldOptionGroup
	fieldType        FieldType
	widget           Widget
	autoID           string
	requiredCSSClass string
	errorCSSClass    string
	attrs            tmplAttrs
	label            string
	labelSuffix      string
	isSafe           bool
	helpText         string
	minLength        uint
	maxLength        uint
	notRequired      bool
	disabled         bool
	sanitizeFunc     func(string) string
	validateFunc     func(string, bool) []Error
	customErrors     map[string]Error
	locale           language.Tag
}

// Type returns the Field type.
func (fld *Field) Type() FieldType {
	return fld.fieldType
}
