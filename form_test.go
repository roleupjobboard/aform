package aform_test

import (
	"github.com/roleupjobboard/aform"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"html/template"
	"testing"
)

func TestMust_dontAlterFieldAutoID(t *testing.T) {
	a := assert.New(t)
	boolFld := aform.Must(aform.DefaultBooleanField("OnOff"))
	emailFld := aform.Must(aform.DefaultEmailField("Email"))
	a.NoError(boolFld.SetAutoID(""))
	a.NoError(emailFld.SetAutoID("custom_%s_custom"))
	a.Equal("", boolFld.AutoID())
	a.Equal("custom_%s_custom", emailFld.AutoID())
	a.NotPanics(func() {
		_ = aform.Must(aform.New(aform.WithBooleanField(boolFld), aform.WithEmailField(emailFld)))
	})
	a.Equal("", boolFld.AutoID())
	a.Equal("custom_%s_custom", emailFld.AutoID())
}

func TestForm_DisableAutoID(t *testing.T) {
	a := assert.New(t)
	boolFld := aform.Must(aform.DefaultBooleanField("Test"))
	emailFld := aform.Must(aform.DefaultEmailField("Email"))
	a.NoError(emailFld.SetAutoID("custom_%s_custom"))
	a.Equal(aform.ExportDefaultAutoID, boolFld.AutoID())
	a.Equal("custom_%s_custom", emailFld.AutoID())
	_, err := aform.New(aform.WithBooleanField(boolFld), aform.DisableAutoID(), aform.WithEmailField(emailFld))
	a.NoError(err)
	a.Equal("", boolFld.AutoID())
	a.Equal("custom_%s_custom", emailFld.AutoID())
}

func TestForm_WithAutoID_invalid(t *testing.T) {
	a := assert.New(t)
	_, err := aform.New(aform.WithAutoID("not_empty_without_verb"))
	a.EqualError(err, "autoID must contain one %s verb (e.g. id_%s). To disable auto ID use DisableAutoID(). Given: not_empty_without_verb")
}

func TestForm_WithAutoID(t *testing.T) {
	a := assert.New(t)
	boolFld := aform.Must(aform.DefaultBooleanField("Test"))
	emailFld := aform.Must(aform.DefaultEmailField("Email"))
	a.NoError(emailFld.SetAutoID("custom_%s_custom"))
	a.Equal(aform.ExportDefaultAutoID, boolFld.AutoID())
	a.Equal("custom_%s_custom", emailFld.AutoID())
	_, err := aform.New(aform.WithBooleanField(boolFld), aform.WithAutoID("custom_form_id_%s"), aform.WithEmailField(emailFld))
	a.NoError(err)
	a.Equal("custom_form_id_%s", boolFld.AutoID())
	a.Equal("custom_%s_custom", emailFld.AutoID())
}

func TestForm_New_dontAlterFieldCSSClass(t *testing.T) {
	a := assert.New(t)
	boolFld := aform.Must(aform.DefaultBooleanField("OnOff"))
	emailFld := aform.Must(aform.DefaultEmailField("Email"))
	boolFld.SetRequiredCSSClass("required_bool")
	boolFld.SetErrorCSSClass("error_bool")
	emailFld.SetRequiredCSSClass("required_email")
	emailFld.SetErrorCSSClass("error_email")
	f, err := aform.New(aform.WithBooleanField(boolFld), aform.WithEmailField(emailFld))
	a.NoError(err)
	f.BindData(map[string][]string{})
	a.False(f.IsValid())
	expected := `
<div class="required_bool error_bool"><label class="required_bool" for="id_onoff">OnOff</label>
<ul class="errorlist"><li id="err_0_id_onoff">This field is required</li></ul>
<input type="checkbox" name="onoff" class="error_bool" id="id_onoff" aria-describedby="err_0_id_onoff" aria-invalid="true" required></div>
<div class="required_email error_email"><label class="required_email" for="id_email">Email</label>
<ul class="errorlist"><li id="err_0_id_email">This field is required</li></ul>
<input type="email" name="email" class="error_email" id="id_email" maxlength="254" aria-describedby="err_0_id_email" aria-invalid="true" required></div>`
	a.Equal(expected, string(f.AsDiv()))
}

func TestForm_WithRequiredCSSClass(t *testing.T) {
	a := assert.New(t)
	boolFld := aform.Must(aform.DefaultBooleanField("OnOff"))
	emailFld := aform.Must(aform.DefaultEmailField("Email"))
	boolFld.SetRequiredCSSClass("required_bool")
	emailFld.SetRequiredCSSClass("required_email")
	f, err := aform.New(aform.WithBooleanField(boolFld), aform.WithRequiredCSSClass("formrequired"), aform.WithEmailField(emailFld))
	a.NoError(err)
	f.BindData(map[string][]string{})
	a.False(f.IsValid())
	expected := `
<div class="formrequired"><label class="formrequired" for="id_onoff">OnOff</label>
<ul class="errorlist"><li id="err_0_id_onoff">This field is required</li></ul>
<input type="checkbox" name="onoff" id="id_onoff" aria-describedby="err_0_id_onoff" aria-invalid="true" required></div>
<div class="formrequired"><label class="formrequired" for="id_email">Email</label>
<ul class="errorlist"><li id="err_0_id_email">This field is required</li></ul>
<input type="email" name="email" id="id_email" maxlength="254" aria-describedby="err_0_id_email" aria-invalid="true" required></div>`
	a.Equal(expected, string(f.AsDiv()))
}

func TestForm_WithErrorCSSClass(t *testing.T) {
	a := assert.New(t)
	boolFld := aform.Must(aform.DefaultBooleanField("OnOff"))
	emailFld := aform.Must(aform.DefaultEmailField("Email"))
	boolFld.SetErrorCSSClass("error_bool")
	emailFld.SetErrorCSSClass("error_bool")
	f, err := aform.New(aform.WithBooleanField(boolFld), aform.WithErrorCSSClass("formerror"), aform.WithEmailField(emailFld))
	a.NoError(err)
	f.BindData(map[string][]string{})
	a.False(f.IsValid())
	expected := `
<div class="formerror"><label for="id_onoff">OnOff</label>
<ul class="errorlist"><li id="err_0_id_onoff">This field is required</li></ul>
<input type="checkbox" name="onoff" class="formerror" id="id_onoff" aria-describedby="err_0_id_onoff" aria-invalid="true" required></div>
<div class="formerror"><label for="id_email">Email</label>
<ul class="errorlist"><li id="err_0_id_email">This field is required</li></ul>
<input type="email" name="email" class="formerror" id="id_email" maxlength="254" aria-describedby="err_0_id_email" aria-invalid="true" required></div>`
	a.Equal(expected, string(f.AsDiv()))
}

func TestForm_BindData_withAcceptLanguage(t *testing.T) {
	a := assert.New(t)
	boolFld := aform.Must(aform.DefaultBooleanField("OnOff"))
	emailFld := aform.Must(aform.DefaultEmailField("Email"))
	f := aform.Must(aform.New(aform.WithLocales([]language.Tag{language.English, language.French}), aform.WithBooleanField(boolFld), aform.WithEmailField(emailFld)))
	f.BindData(map[string][]string{}, "fr,en;q=0.9,ru;q=0.8")
	a.False(f.IsValid())
	expected := `
<div><label for="id_onoff">OnOff</label>
<ul class="errorlist"><li id="err_0_id_onoff">Ce champ est obligatoire</li></ul>
<input type="checkbox" name="onoff" id="id_onoff" aria-describedby="err_0_id_onoff" aria-invalid="true" required></div>
<div><label for="id_email">Email</label>
<ul class="errorlist"><li id="err_0_id_email">Ce champ est obligatoire</li></ul>
<input type="email" name="email" id="id_email" maxlength="254" aria-describedby="err_0_id_email" aria-invalid="true" required></div>`
	a.Equal(expected, string(f.AsDiv()))
}

func TestForm_IsBound_afterInit(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First")))))
	a.False(f.IsBound())
}

func TestForm_IsBound_withEmptyData(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First")))))
	f.BindData(map[string][]string{})
	a.True(f.IsBound())
}

func TestForm_IsBound_withDataNotMatchingFields(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First")))))
	f.BindData(map[string][]string{"not_matching": {"whatever"}})
	a.True(f.IsBound())
}

func TestForm_IsBound_withDataMatchingField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First")))))
	f.BindData(map[string][]string{"first": {"sure thing"}})
	a.True(f.IsBound())
}

func TestForm_AsDiv(t *testing.T) {
	tests := []struct {
		name string
		form *aform.Form
		want template.HTML
	}{
		{
			name: "empty",
			form: aform.Must(aform.New()),
			want: "",
		},
		{
			name: "two char fields",
			form: aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First"))), aform.WithCharField(aform.Must(aform.DefaultCharField("Second"))))),
			want: `
<div><label for="id_first">First</label><input type="text" name="first" id="id_first" maxlength="256" required></div>
<div><label for="id_second">Second</label><input type="text" name="second" id="id_second" maxlength="256" required></div>`,
		},
		{
			name: "auto ID deactivated",
			form: aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First"))), aform.DisableAutoID())),
			want: `
<div>First<input type="text" name="first" maxlength="256" required></div>`,
		},
		{
			name: "auto ID customized",
			form: aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First"))), aform.WithAutoID("%s_id"))),
			want: `
<div><label for="first_id">First</label><input type="text" name="first" id="first_id" maxlength="256" required></div>`,
		},
		{
			name: "auto ID customized per fields",
			form: func() *aform.Form {
				fld1 := aform.Must(aform.DefaultCharField("First"))
				_ = fld1.SetAutoID("f_%s_f")
				fld2 := aform.Must(aform.DefaultCharField("Second"))
				_ = fld2.SetAutoID("s_%s_s")
				return aform.Must(aform.New(aform.WithCharField(fld1), aform.WithCharField(fld2)))
			}(),
			want: `
<div><label for="f_first_f">First</label><input type="text" name="first" id="f_first_f" maxlength="256" required></div>
<div><label for="s_second_s">Second</label><input type="text" name="second" id="s_second_s" maxlength="256" required></div>`,
		},
		{
			name: "auto ID customization per fields override disable by the form",
			form: func() *aform.Form {
				fld1 := aform.Must(aform.DefaultCharField("First"))
				_ = fld1.SetAutoID("f_%s_f")
				fld2 := aform.Must(aform.DefaultCharField("Second"))
				_ = fld2.SetAutoID("s_%s_s")
				return aform.Must(aform.New(aform.WithCharField(fld1), aform.WithCharField(fld2), aform.DisableAutoID()))
			}(),
			want: `
<div><label for="f_first_f">First</label><input type="text" name="first" id="f_first_f" maxlength="256" required></div>
<div><label for="s_second_s">Second</label><input type="text" name="second" id="s_second_s" maxlength="256" required></div>`,
		},
		{
			name: "label suffix customized",
			form: aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First"))), aform.WithLabelSuffix(":"), aform.WithCharField(aform.Must(aform.DefaultCharField("Second"))))),
			want: `
<div><label for="id_first">First:</label><input type="text" name="first" id="id_first" maxlength="256" required></div>
<div><label for="id_second">Second:</label><input type="text" name="second" id="id_second" maxlength="256" required></div>`,
		},
		{
			name: "label suffix contains HTML",
			form: aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First"))), aform.WithCharField(aform.Must(aform.DefaultCharField("First", aform.IsSafe()))), aform.WithLabelSuffix("<strong>*</strong>"))),
			want: `
<div><label for="id_first">First&lt;strong&gt;*&lt;/strong&gt;</label><input type="text" name="first" id="id_first" maxlength="256" required></div>
<div><label for="id_first">First<strong>*</strong></label><input type="text" name="first" id="id_first" maxlength="256" required></div>`,
		},
		{
			name: "label suffix customized per fields",
			form: func() *aform.Form {
				fld1 := aform.Must(aform.DefaultCharField("First"))
				fld1.SetLabelSuffix("*")
				fld2 := aform.Must(aform.DefaultCharField("Second"))
				fld2.SetLabelSuffix(":")
				return aform.Must(aform.New(aform.WithCharField(fld1), aform.WithCharField(fld2)))
			}(),
			want: `
<div><label for="id_first">First*</label><input type="text" name="first" id="id_first" maxlength="256" required></div>
<div><label for="id_second">Second:</label><input type="text" name="second" id="id_second" maxlength="256" required></div>`,
		},
		{
			name: "with required class",
			form: aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First"))), aform.WithRequiredCSSClass("required"), aform.WithCharField(aform.Must(aform.DefaultCharField("Second", aform.IsNotRequired()))), aform.WithCharField(aform.Must(aform.DefaultCharField("Third"))))),
			want: `
<div class="required"><label class="required" for="id_first">First</label><input type="text" name="first" id="id_first" maxlength="256" required></div>
<div><label for="id_second">Second</label><input type="text" name="second" id="id_second" maxlength="256"></div>
<div class="required"><label class="required" for="id_third">Third</label><input type="text" name="third" id="id_third" maxlength="256" required></div>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.form.AsDiv(); got != tt.want {
				t.Errorf("AsDiv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestForm_AsDiv_withDisableAutoID(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithBooleanField(aform.Must(aform.DefaultBooleanField("Bool Field"))),
		aform.WithCharField(aform.Must(aform.DefaultCharField("Char Field"))),
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("Select Choice Field", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "option_1", Label: "First Option"}, {Value: "option_2", Label: "Second Option"}})))),
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("Radio Select Choice Field", aform.WithWidget(aform.RadioSelect), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "option_1", Label: "First Option"}, {Value: "option_2", Label: "Second Option"}})))),
		aform.WithEmailField(aform.Must(aform.DefaultEmailField("Email Field"))),
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Multiple Select Choice Field", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "option_1", Label: "First Option"}, {Value: "option_2", Label: "Second Option"}})))),
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Multiple Checkbox Choice Field", aform.WithWidget(aform.CheckboxSelectMultiple), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "option_1", Label: "First Option"}, {Value: "option_2", Label: "Second Option"}})))),
		aform.DisableAutoID(),
	))
	f.BindData(map[string][]string{})
	a.False(f.IsValid())
	want := template.HTML(`
<div>Bool Field
<ul class="errorlist"><li>This field is required</li></ul>
<input type="checkbox" name="bool_field" aria-invalid="true" required></div>
<div>Char Field
<ul class="errorlist"><li>This field is required</li></ul>
<input type="text" name="char_field" maxlength="256" aria-invalid="true" required></div>
<div>Select Choice Field
<ul class="errorlist"><li>This field is required</li></ul>
<select name="select_choice_field" aria-invalid="true" required>
  <option value="option_1">First Option</option>
  <option value="option_2">Second Option</option>
</select></div>
<div>
<fieldset>Radio Select Choice Field
<ul class="errorlist"><li>This field is required</li></ul>
<div>
<label><input type="radio" name="radio_select_choice_field" value="option_1">First Option</label>
<label><input type="radio" name="radio_select_choice_field" value="option_2">Second Option</label>
</div>
</fieldset>
</div>
<div>Email Field
<ul class="errorlist"><li>This field is required</li></ul>
<input type="email" name="email_field" maxlength="254" aria-invalid="true" required></div>
<div>Multiple Select Choice Field
<ul class="errorlist"><li>This field is required</li></ul>
<select name="multiple_select_choice_field" aria-invalid="true" multiple required>
  <option value="option_1">First Option</option>
  <option value="option_2">Second Option</option>
</select></div>
<div>
<fieldset>Multiple Checkbox Choice Field
<ul class="errorlist"><li>This field is required</li></ul>
<div>
<label><input type="checkbox" name="multiple_checkbox_choice_field" value="option_1">First Option</label>
<label><input type="checkbox" name="multiple_checkbox_choice_field" value="option_2">Second Option</label>
</div>
</fieldset>
</div>`)
	a.Equal(want, f.AsDiv())
}

func TestForm_AsDiv_withDisableAutoIDAndRequiredClassAndErrorClass(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithRequiredCSSClass("required_class"),
		aform.WithErrorCSSClass("error_class"),
		aform.WithBooleanField(aform.Must(aform.DefaultBooleanField("Bool Field"))),
		aform.WithCharField(aform.Must(aform.DefaultCharField("Char Field"))),
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("Select Choice Field", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "option_1", Label: "First Option"}, {Value: "option_2", Label: "Second Option"}})))),
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("Radio Select Choice Field", aform.WithWidget(aform.RadioSelect), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "option_1", Label: "First Option"}, {Value: "option_2", Label: "Second Option"}})))),
		aform.WithEmailField(aform.Must(aform.DefaultEmailField("Email Field"))),
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Multiple Select Choice Field", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "option_1", Label: "First Option"}, {Value: "option_2", Label: "Second Option"}})))),
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Multiple Checkbox Choice Field", aform.WithWidget(aform.CheckboxSelectMultiple), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "option_1", Label: "First Option"}, {Value: "option_2", Label: "Second Option"}})))),
		aform.DisableAutoID(),
	))
	f.BindData(map[string][]string{})
	a.False(f.IsValid())
	want := template.HTML(`
<div class="required_class error_class">Bool Field
<ul class="errorlist"><li>This field is required</li></ul>
<input type="checkbox" name="bool_field" class="error_class" aria-invalid="true" required></div>
<div class="required_class error_class">Char Field
<ul class="errorlist"><li>This field is required</li></ul>
<input type="text" name="char_field" class="error_class" maxlength="256" aria-invalid="true" required></div>
<div class="required_class error_class">Select Choice Field
<ul class="errorlist"><li>This field is required</li></ul>
<select name="select_choice_field" class="error_class" aria-invalid="true" required>
  <option value="option_1">First Option</option>
  <option value="option_2">Second Option</option>
</select></div>
<div class="required_class error_class">
<fieldset>Radio Select Choice Field
<ul class="errorlist"><li>This field is required</li></ul>
<div class="error_class">
<label><input type="radio" name="radio_select_choice_field" value="option_1">First Option</label>
<label><input type="radio" name="radio_select_choice_field" value="option_2">Second Option</label>
</div>
</fieldset>
</div>
<div class="required_class error_class">Email Field
<ul class="errorlist"><li>This field is required</li></ul>
<input type="email" name="email_field" class="error_class" maxlength="254" aria-invalid="true" required></div>
<div class="required_class error_class">Multiple Select Choice Field
<ul class="errorlist"><li>This field is required</li></ul>
<select name="multiple_select_choice_field" class="error_class" aria-invalid="true" multiple required>
  <option value="option_1">First Option</option>
  <option value="option_2">Second Option</option>
</select></div>
<div class="required_class error_class">
<fieldset>Multiple Checkbox Choice Field
<ul class="errorlist"><li>This field is required</li></ul>
<div class="error_class">
<label><input type="checkbox" name="multiple_checkbox_choice_field" value="option_1">First Option</label>
<label><input type="checkbox" name="multiple_checkbox_choice_field" value="option_2">Second Option</label>
</div>
</fieldset>
</div>`)
	a.Equal(want, f.AsDiv())
}

func TestForm_AsDiv_withMultipleChoiceFieldAndMultipleValidData(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Select", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "option_1", Label: "First Option"}, {Value: "option_2", Label: "Second Option"}, {Value: "option_3", Label: "Third Option"}})))),
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Checkbox", aform.WithWidget(aform.CheckboxSelectMultiple), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "option_1", Label: "First Option"}, {Value: "option_2", Label: "Second Option"}, {Value: "option_3", Label: "Third Option"}})))),
	))
	f.BindData(map[string][]string{"select": {"option_1", "option_2"}, "checkbox": {"option_1", "option_3"}})
	a.True(f.IsValid())
	want := template.HTML(`
<div><label for="id_select">Select</label><select name="select" id="id_select" multiple required>
  <option value="option_1" id="id_select_0" selected>First Option</option>
  <option value="option_2" id="id_select_1" selected>Second Option</option>
  <option value="option_3" id="id_select_2">Third Option</option>
</select></div>
<div>
<fieldset><legend for="id_checkbox">Checkbox</legend>
<div id="id_checkbox">
<label for="id_checkbox_0"><input type="checkbox" name="checkbox" value="option_1" id="id_checkbox_0" checked>First Option</label>
<label for="id_checkbox_1"><input type="checkbox" name="checkbox" value="option_2" id="id_checkbox_1">Second Option</label>
<label for="id_checkbox_2"><input type="checkbox" name="checkbox" value="option_3" id="id_checkbox_2" checked>Third Option</label>
</div>
</fieldset>
</div>`)
	a.Equal(want, f.AsDiv())
}

func TestForm_AsDiv_withMultipleChoiceFieldAndInvalidData(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Select", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "option_1", Label: "First Option"}, {Value: "option_2", Label: "Second Option"}, {Value: "option_3", Label: "Third Option"}})))),
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Checkbox", aform.WithWidget(aform.CheckboxSelectMultiple), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "option_1", Label: "First Option"}, {Value: "option_2", Label: "Second Option"}, {Value: "option_3", Label: "Third Option"}})))),
	))
	f.BindData(map[string][]string{"select": {"option_1", "option_invalid", "option_2"}, "checkbox": {"option_1", "option_3", "option_invalid"}})
	a.False(f.IsValid())
	want := template.HTML(`
<div><label for="id_select">Select</label>
<ul class="errorlist"><li id="err_0_id_select">Invalid choice</li></ul>
<select name="select" id="id_select" aria-describedby="err_0_id_select" aria-invalid="true" multiple required>
  <option value="option_1" id="id_select_0" selected>First Option</option>
  <option value="option_2" id="id_select_1" selected>Second Option</option>
  <option value="option_3" id="id_select_2">Third Option</option>
</select></div>
<div>
<fieldset><legend for="id_checkbox">Checkbox</legend>
<ul class="errorlist"><li id="err_0_id_checkbox">Invalid choice</li></ul>
<div id="id_checkbox">
<label for="id_checkbox_0"><input type="checkbox" name="checkbox" value="option_1" id="id_checkbox_0" checked>First Option</label>
<label for="id_checkbox_1"><input type="checkbox" name="checkbox" value="option_2" id="id_checkbox_1">Second Option</label>
<label for="id_checkbox_2"><input type="checkbox" name="checkbox" value="option_3" id="id_checkbox_2" checked>Third Option</label>
</div>
</fieldset>
</div>`)
	a.Equal(want, f.AsDiv())
}
