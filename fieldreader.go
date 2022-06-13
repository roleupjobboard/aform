package aform

import (
	"strings"
)

// Name returns the name given to the field.
func (fld *Field) Name() string {
	return fld.name
}

// HTMLName returns the field's name transformed to be the value of the HTML
// name attribute.
func (fld *Field) HTMLName() string {
	return normalizedNameForField(fld)
}

// AutoID returns the auto ID set with SetAutoID.
func (fld *Field) AutoID() string {
	return fld.autoID
}

// LabelSuffix returns the suffix set with SetLabelSuffix.
func (fld *Field) LabelSuffix() string {
	return fld.labelSuffix
}

// UseFieldset returns true if the widget used by the field
func (fld *Field) UseFieldset() bool {
	switch fld.widget {
	case RadioSelect, CheckboxSelectMultiple:
		return true
	default:
		return false
	}
}

func (fld *Field) CSSClasses() string {
	var classes []string
	if !fld.notRequired && len(fld.requiredCSSClass) > 0 {
		classes = append(classes, fld.requiredCSSClass)
	}
	if fld.HasErrors() && len(fld.errorCSSClass) > 0 {
		classes = append(classes, fld.errorCSSClass)
	}
	return strings.Join(classes, " ")
}

// Required returns true if the field is required
func (fld *Field) Required() bool {
	return !fld.notRequired
}

// HasHelpText returns true if a help text has been added to the field.
func (fld *Field) HasHelpText() bool {
	return len(fld.helpText) > 0
}

// HasErrors returns true if an input is bound to the field and there is a
// validation error.
func (fld *Field) HasErrors() bool {
	return len(fld.errors) > 0
}
