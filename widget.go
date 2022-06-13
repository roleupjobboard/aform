package aform

import (
	"fmt"
	"golang.org/x/exp/slices"
)

// Widget defines the type of the widgets.
type Widget string

// Widgets are used by the fields to handle HTML rendering.
// Each field as a default widget. It can be changed with WithWidget
// or SetWidget functions.
const (
	// TextInput renders as: <input type="text" ...>
	TextInput = Widget("TextInput")
	// EmailInput renders as: <input type="email" ...>
	EmailInput = Widget("EmailInput")
	// URLInput renders as: <input type="url" ...>
	URLInput = Widget("URLInput")
	// PasswordInput renders as: <input type="password" ...>
	PasswordInput = Widget("PasswordInput")
	// HiddenInput renders as: <input type="hidden" ...>
	HiddenInput = Widget("HiddenInput")
	// TextArea renders as: <textarea>...</textarea>
	TextArea = Widget("TextArea")
	// CheckboxInput renders as: <input type="checkbox" ...>
	CheckboxInput = Widget("CheckboxInput")
	// Select renders as: <select><option ...>...</select>
	Select = Widget("Select")
	// SelectMultiple renders as Select but allow multiple selection: <select multiple>...<	/select>
	SelectMultiple = Widget("SelectMultiple")
	// RadioSelect is similar to Select, but rendered as a list of radio buttons within <div> tags:
	// 	<div>
	//    <div><input type="radio" name="..."></div>
	//    ...
	// 	</div>
	RadioSelect = Widget("RadioSelect")
	// CheckboxSelectMultiple is similar SelectMultiple, but rendered as a list of checkboxes
	// 	<div>
	//    <div><input type="checkbox" name="..." ></div>
	//    ...
	// 	</div>
	CheckboxSelectMultiple = Widget("CheckboxSelectMultiple")
)

func (t Widget) htmlType() string {
	switch t {
	case TextInput:
		return "text"
	case EmailInput:
		return "email"
	case URLInput:
		return "url"
	case PasswordInput:
		return "password"
	case HiddenInput:
		return "hidden"
	case TextArea:
		return "textarea"
	case CheckboxInput:
		return "checkbox"
	case Select:
		return "select"
	case RadioSelect:
		return "radio"
	case SelectMultiple:
		return "select"
	case CheckboxSelectMultiple:
		return "checkbox_select"
	default:
		panic(fmt.Sprintf("%s: unknown widget type", t))
	}
}

func (t Widget) optionWidget() Widget {
	switch t {
	case CheckboxSelectMultiple:
		return CheckboxInput
	default:
		return t
	}
}

func (t Widget) isInput() bool {
	list := []Widget{TextInput, EmailInput, URLInput, PasswordInput, HiddenInput, TextArea, CheckboxInput}
	return slices.Contains(list, t)
}

func (t Widget) isChoice() bool {
	list := []Widget{Select, RadioSelect, SelectMultiple, CheckboxSelectMultiple}
	return slices.Contains(list, t)
}

func (t Widget) isMultiChoice() bool {
	return t == SelectMultiple || t == CheckboxSelectMultiple
}

func (t Widget) selectedAttr(selected bool) (nameValueAttr[string], bool) {
	if !selected {
		return nameValueAttr[string]{}, false
	}
	switch t {
	case TextInput, EmailInput, URLInput, PasswordInput, HiddenInput, TextArea:
		return nameValueAttr[string]{}, false
	case CheckboxInput, CheckboxSelectMultiple:
		return nameValueAttr[string]{n: "checked", v: ""}, true
	case Select, SelectMultiple:
		return nameValueAttr[string]{n: "selected", v: ""}, true
	case RadioSelect:
		return nameValueAttr[string]{n: "checked", v: ""}, true
	default:
		return nameValueAttr[string]{}, false
	}
}

// noAttrValue returns true if the widget never has a value attribute.
func (t Widget) noAttrValue() bool {
	return t == CheckboxInput
}

func (t Widget) defaultSanitizeFunc() SanitizationFunc {
	switch t {
	case TextArea:
		return sanitizeToPlainText
	default:
		return sanitizeToOneLinePlainText
	}
}
