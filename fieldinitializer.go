package aform

// SetLabel overrides the default label of the Field. By default, the label is
// the name of the Field given as parameter to a Field creation function. The
// label is HTML-escaped. To alter more the for= attribute or to completely
// remove the tag <label> use SetAutoID.
func (fld *Field) SetLabel(label string) {
	fld.label = label
}

// MarkSafe marks the field label as safe for HTML. It means the label
// and the label suffix are no more HTML-escaped.
func (fld *Field) MarkSafe() {
	fld.isSafe = true
}

// SetHelpText adds a help text to the Field. Help text is not HTML-escaped.
func (fld *Field) SetHelpText(help string) {
	fld.helpText = help
}

// SetNotRequired sets the Field as not required. By default,
// all fields are required.
func (fld *Field) SetNotRequired() {
	fld.notRequired = true
}

// SetDisabled sets the Field as disabled.
func (fld *Field) SetDisabled() {
	fld.disabled = true
}

// AddChoiceOptions adds a list of ChoiceFieldOption to the Field. if the
// parameter label is the empty string options are not grouped together.
// Example without group label:
//   fld.AddChoiceOptions("", []ChoiceFieldOption{{Value: "red", Label: "Rouge"}, {Value: "green", Label: "Vert"}})
//   // HTML output:
//   // <option value="red" id="id_color_0">Rouge</option>
//   // <option value="green" id="id_color_1">Vert</option>
// Example with group label:
//   fld.AddChoiceOptions("RG", []ChoiceFieldOption{{Value: "red", Label: "Rouge"}, {Value: "green", Label: "Vert"}})
//   // HTML output:
//   //  <optgroup label="RG">
//   //    <option value="red" id="id_color_0_0">Rouge</option>
//   //    <option value="green" id="id_color_0_1">Vert</option>
//   //  </optgroup>
func (fld *Field) AddChoiceOptions(label string, options []ChoiceFieldOption) {
	fld.optionGroups = append(fld.optionGroups, choiceFieldOptionGroup{label: options})
}

func (fld *Field) addError(err Error) {
	fld.errors = append(fld.errors, err)
}
