package aform

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
)

// AsDiv renders the field in a <div> tag.
func (fld *Field) AsDiv() template.HTML {
	return mustFieldAsDivTemplate(fld)
}

// LabelTag renders the <label> tag.
func (fld *Field) LabelTag() template.HTML {
	return fld.labelOrLegendTag("label")
}

// LegendTag renders the <legend> tag. This tag replaces <label> tag
// when a <fieldset> tag is used to group together options of fields
// using RadioSelect widget or CheckboxSelectMultiple widget.
func (fld *Field) LegendTag() template.HTML {
	return fld.labelOrLegendTag("legend")
}

func (fld *Field) labelOrLegendTag(tag string) template.HTML {
	useTag := hasID(fld)
	attrs := map[string]string{}
	if useTag {
		attrs["for"] = normalizedIDForField(fld)
	}
	classes := fld.labelCSSClassList()
	if len(classes) > 0 {
		attrs["class"] = strings.Join(classes, " ")
	}
	safeLbl := newSafeLabel(useTag, fld.label, fld.labelSuffix, attrs)
	lbl := newLabel(useTag, fld.label, fld.labelSuffix, attrs)
	switch tag {
	case "label":
		if fld.isSafe {
			return mustLabelTemplate(&safeLbl)
		}
		return mustLabelTemplate(&lbl)
	case "legend":
		if fld.isSafe {
			return mustLegendTemplate(&safeLbl)
		}
		return mustLegendTemplate(&lbl)
	default:
		panic(fmt.Sprintf("%T: incompatible label tag %s", fld, tag))
	}
}

// Widget renders the widget.
func (fld *Field) Widget() template.HTML {
	switch fld.widget {
	case TextInput, EmailInput, URLInput, PasswordInput, HiddenInput, TextArea, CheckboxInput:
		return fld.widgetInput(fld.widgetCSSClassList())
	case Select, RadioSelect, SelectMultiple, CheckboxSelectMultiple:
		return fld.widgetChoice(fld.widgetCSSClassList())
	default:
		panic(fmt.Sprintf("%T: incompatible type %s", fld, fld.widget))
	}
}

func (fld *Field) widgetInput(classes []string) template.HTML {
	if !fld.widget.isInput() {
		panic(fmt.Sprintf("%T: incompatible type %s", fld, fld.widget))
	}

	attrs := attributesForField(fld, classes)
	value := ""
	if len(fld.boundValues) > 0 {
		value = fld.boundValues[0]
	}
	if selectedAttr, ok := fld.widget.selectedAttr(valueToBool(value)); ok {
		attrs[selectedAttr.n] = selectedAttr.v
	}
	if fld.widget.noAttrValue() {
		value = ""
	}
	return mustInputTemplate(&widgetInput{
		Type:  fld.widget,
		Name:  normalizedNameForField(fld),
		Value: value,
		Attrs: attrs,
	})
}

func (fld *Field) widgetChoice(classes []string) template.HTML {
	if !fld.widget.isChoice() {
		panic(fmt.Sprintf("%T: incompatible type %s", fld, fld.widget))
	}

	attrs := attributesForField(fld, classes)
	values := func(isMultiChoice bool) []string {
		if isMultiChoice {
			return fld.boundValues
		} else {
			if len(fld.boundValues) > 0 {
				return fld.boundValues[:1]
			}
		}
		return nil
	}(fld.widget.isMultiChoice())
	return mustChoiceTemplate(&widgetChoice{
		Type:   fld.widget,
		Name:   normalizedNameForField(fld),
		Values: values,
		Groups: fld.widgetGroups(values),
		Attrs:  attrs,
	})
}

func (fld *Field) widgetGroups(selected []string) []map[string][]widgetOption {
	return fieldGroupsToWidgetGroups(fld.optionGroups, fld.widget.optionWidget(), normalizedIDForField(fld), normalizedNameForField(fld), selected)
}

func attributesForField(fld *Field, classes []string) tmplAttrs {
	attrs := tmplAttrs{}
	if fld.widget.isMultiChoice() {
		attrs["multiple"] = ""
	}
	if hasID(fld) {
		attrs["id"] = normalizedIDForField(fld)
	}
	if !fld.notRequired {
		attrs["required"] = ""
	}
	if fld.disabled {
		attrs["disabled"] = ""
	}
	if fld.HasErrors() {
		attrs["aria-invalid"] = "true"
	}
	if hasID(fld) {
		var ariaDescribedBy []string
		if len(fld.helpText) > 0 {
			ariaDescribedBy = append(ariaDescribedBy, normalizedDescribedByIDForHelpText(fld))
		}
		if fld.HasErrors() {
			ariaDescribedBy = append(ariaDescribedBy, normalizedDescribedByIDErrList(fld, len(fld.errors))...)
		}
		if len(ariaDescribedBy) > 0 {
			attrs["aria-describedby"] = strings.Join(ariaDescribedBy, " ")
		}
	}
	if fld.minLength > 0 {
		attrs["minlength"] = strconv.FormatUint(uint64(fld.minLength), 10)
	}
	if fld.maxLength > 0 {
		attrs["maxlength"] = strconv.FormatUint(uint64(fld.maxLength), 10)
	}
	for name, value := range fld.attrs {
		attrs[name] = value
	}
	if len(classes) > 0 {
		if classesFromAttributes, ok := attrs["class"]; ok {
			classes = append(classes, classesFromAttributes)
		}
		attrs["class"] = strings.Join(classes, " ")
	}
	return attrs
}

// Errors renders field errors in a <ul> tag with each <li> children tag
// containing one error. CSS class errorlist is set on the <ul> tag. Each <li>
// tag as a generated unique ID based on the error index and the field ID.
func (fld *Field) Errors() template.HTML {
	if !fld.HasErrors() {
		return ""
	}
	list := make([]tmplError, len(fld.errors))
	for i, err := range fld.errors {
		attrs := tmplAttrs{}
		if hasID(fld) {
			attrs["id"] = normalizedDescribedByIDForErr(fld, i)
		}
		list[i] = tmplError{
			Text:  err.Translate(fld.locale.String()),
			Attrs: attrs,
		}
	}
	return mustErrorsTemplate(&tmplErrors{
		List:  list,
		Attrs: map[string]string{"class": "errorlist"},
	})
}

// HelpText renders the help text in a <span> tag. CSS class helptext and an ID
// generated from the field ID are set on the <span> tag.
func (fld *Field) HelpText() template.HTML {
	if len(fld.helpText) == 0 {
		return ""
	}
	attrs := map[string]string{"class": "helptext"}
	if hasID(fld) {
		attrs["id"] = normalizedDescribedByIDForHelpText(fld)
	}
	return mustHelpTextTemplate(&tmplHelpText{
		Text:  template.HTML(fld.helpText),
		Attrs: attrs,
	})
}

func (fld *Field) labelCSSClassList() []string {
	var classes []string
	if !fld.notRequired && len(fld.requiredCSSClass) > 0 {
		classes = append(classes, fld.requiredCSSClass)
	}
	return classes
}

func (fld *Field) widgetCSSClassList() []string {
	var classes []string
	if fld.HasErrors() && len(fld.errorCSSClass) > 0 {
		classes = append(classes, fld.errorCSSClass)
	}
	return classes
}
