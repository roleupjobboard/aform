package aform_test

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"html/template"
	"testing"
)

func TestField_AsDiv_boolean(t *testing.T) {
	tests := []struct {
		name  string
		field *aform.BooleanField
		want  template.HTML
	}{
		{
			name:  "boolean field",
			field: aform.Must(aform.DefaultBooleanField("Remember me")),
			want:  `<div><label for="id_remember_me">Remember me</label><input type="checkbox" name="remember_me" id="id_remember_me" required></div>`,
		},
		{
			// https://html.spec.whatwg.org/multipage/syntax.html#syntax-ambiguous-ampersand
			name:  "boolean field with ambiguous ampersand",
			field: aform.Must(aform.DefaultBooleanField("Remember me&123;")),
			want:  `<div><label for="id_remember_me&123;">Remember me&amp;123;</label><input type="checkbox" name="remember_me&123;" id="id_remember_me&123;" required></div>`,
		},
		{
			// https://html.spec.whatwg.org/multipage/syntax.html#syntax-ambiguous-ampersand
			name:  "boolean field with ambiguous ampersand matching a named character reference",
			field: aform.Must(aform.DefaultBooleanField("Remember me&aacute;")),
			want:  `<div><label for="id_remember_me&aacute;">Remember me&amp;aacute;</label><input type="checkbox" name="remember_me&aacute;" id="id_remember_me&aacute;" required></div>`,
		},
		{
			// https://html.spec.whatwg.org/multipage/syntax.html#syntax-ambiguous-ampersand
			name:  "boolean field with not ambiguous ampersand",
			field: aform.Must(aform.DefaultBooleanField("Remember me&123")),
			want:  `<div><label for="id_remember_me&123">Remember me&amp;123</label><input type="checkbox" name="remember_me&123" id="id_remember_me&123" required></div>`,
		},
		{
			name:  "boolean field with HTML in label",
			field: aform.Must(aform.DefaultBooleanField("Remember me", aform.WithLabel("Remember <strong>me</strong>"))),
			want:  `<div><label for="id_remember_me">Remember &lt;strong&gt;me&lt;/strong&gt;</label><input type="checkbox" name="remember_me" id="id_remember_me" required></div>`,
		},
		{
			name:  "boolean field marked as safe with HTML in label",
			field: aform.Must(aform.DefaultBooleanField("Remember me", aform.WithLabel("Remember <strong>me</strong>"), aform.IsSafe())),
			want:  `<div><label for="id_remember_me">Remember <strong>me</strong></label><input type="checkbox" name="remember_me" id="id_remember_me" required></div>`,
		},
		{
			name: "boolean field required with empty value",
			field: func() *aform.BooleanField {
				fld := aform.Must(aform.DefaultBooleanField("Remember me"))
				fld.Clean("")
				return fld
			}(),
			want: `<div><label for="id_remember_me">Remember me</label>
<ul class="errorlist"><li id="err_0_id_remember_me">This field is required</li></ul>
<input type="checkbox" name="remember_me" id="id_remember_me" aria-describedby="err_0_id_remember_me" aria-invalid="true" required></div>`,
		},
		{
			name: "boolean field not required with empty value",
			field: func() *aform.BooleanField {
				fld := aform.Must(aform.DefaultBooleanField("Remember me", aform.IsNotRequired()))
				fld.Clean("")
				return fld
			}(),
			want: `<div><label for="id_remember_me">Remember me</label><input type="checkbox" name="remember_me" id="id_remember_me"></div>`,
		},
		{
			name: "boolean field required with false value",
			field: func() *aform.BooleanField {
				fld := aform.Must(aform.DefaultBooleanField("Remember me"))
				fld.Clean("off")
				return fld
			}(),
			want: `<div><label for="id_remember_me">Remember me</label><input type="checkbox" name="remember_me" id="id_remember_me" required></div>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.AsDiv(); got != tt.want {
				t.Errorf("AsDiv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_AsDiv_char(t *testing.T) {
	tests := []struct {
		name  string
		field *aform.CharField
		want  template.HTML
	}{
		{
			name:  "char field",
			field: aform.Must(aform.DefaultCharField("Your Name")),
			want:  `<div><label for="id_your_name">Your Name</label><input type="text" name="your_name" id="id_your_name" maxlength="256" required></div>`,
		},
		{
			name: "char field with empty value",
			field: func() *aform.CharField {
				fld := aform.Must(aform.DefaultCharField("Your Name"))
				fld.Clean("")
				return fld
			}(),
			want: `<div><label for="id_your_name">Your Name</label>
<ul class="errorlist"><li id="err_0_id_your_name">This field is required</li></ul>
<input type="text" name="your_name" id="id_your_name" maxlength="256" aria-describedby="err_0_id_your_name" aria-invalid="true" required></div>`,
		},
		{
			name: "char field with help text and empty value",
			field: func() *aform.CharField {
				fld := aform.Must(aform.DefaultCharField("Your Name", aform.WithHelpText("Enter your first name")))
				fld.Clean("")
				return fld
			}(),
			want: `<div><label for="id_your_name">Your Name</label>
<ul class="errorlist"><li id="err_0_id_your_name">This field is required</li></ul>
<input type="text" name="your_name" id="id_your_name" maxlength="256" aria-describedby="helptext_id_your_name err_0_id_your_name" aria-invalid="true" required>
<span class="helptext" id="helptext_id_your_name">Enter your first name</span></div>`,
		},
		{
			name:  "char field with email widget",
			field: aform.Must(aform.DefaultCharField("Your email", aform.WithWidget(aform.EmailInput))),
			want:  `<div><label for="id_your_email">Your email</label><input type="email" name="your_email" id="id_your_email" maxlength="256" required></div>`,
		},
		{
			name: "char field with email widget and empty value",
			field: func() *aform.CharField {
				fld := aform.Must(aform.DefaultCharField("Your email", aform.WithWidget(aform.EmailInput)))
				fld.Clean("")
				return fld
			}(),
			want: `<div><label for="id_your_email">Your email</label>
<ul class="errorlist"><li id="err_0_id_your_email">This field is required</li></ul>
<input type="email" name="your_email" id="id_your_email" maxlength="256" aria-describedby="err_0_id_your_email" aria-invalid="true" required></div>`,
		},
		{
			name:  "char field with url widget",
			field: aform.Must(aform.DefaultCharField("Homepage", aform.WithWidget(aform.URLInput))),
			want:  `<div><label for="id_homepage">Homepage</label><input type="url" name="homepage" id="id_homepage" maxlength="256" required></div>`,
		},
		{
			name:  "char field with password widget",
			field: aform.Must(aform.DefaultCharField("Password", aform.WithWidget(aform.PasswordInput))),
			want:  `<div><label for="id_password">Password</label><input type="password" name="password" id="id_password" maxlength="256" required></div>`,
		},
		{
			name:  "char field with textarea widget",
			field: aform.Must(aform.DefaultCharField("Description", aform.WithWidget(aform.TextArea))),
			want: `<div><label for="id_description">Description</label><textarea name="description" id="id_description" maxlength="256" required>
</textarea></div>`,
		},
		{
			name:  "char field with hidden widget",
			field: aform.Must(aform.DefaultCharField("hidden field", aform.WithWidget(aform.HiddenInput))),
			want:  `<div><label for="id_hidden_field">hidden field</label><input type="hidden" name="hidden_field" id="id_hidden_field" maxlength="256" required></div>`,
		},
		{
			name: "char field with auto ID disabled",
			field: func() *aform.CharField {
				f := aform.Must(aform.DefaultCharField("Your Name"))
				_ = f.SetAutoID("")
				return f
			}(),
			want: `<div>Your Name<input type="text" name="your_name" maxlength="256" required></div>`,
		},
		{
			name: "char field with textarea widget and multiline initial data",
			field: func() *aform.CharField {
				initial := `line 1
line 2

line 4`
				return aform.Must(aform.NewCharField("Description", initial, "", 0, 0, aform.WithWidget(aform.TextArea)))
			}(),
			want: `<div><label for="id_description">Description</label><textarea name="description" id="id_description" required>
line 1
line 2

line 4</textarea></div>`,
		},
		{
			name: "char field with textarea widget and multiline bound data",
			field: func() *aform.CharField {
				f := aform.Must(aform.DefaultCharField("Description", aform.WithWidget(aform.TextArea)))
				f.Clean(`line 1
line 2

line 4`)
				return f
			}(),
			want: `<div><label for="id_description">Description</label><textarea name="description" id="id_description" maxlength="256" required>
line 1
line 2

line 4</textarea></div>`,
		},
		{
			name: "char field with max > 0",
			field: func() *aform.CharField {
				f, _ := aform.NewCharField("Description", "", "", 0, 20)
				return f
			}(),
			want: `<div><label for="id_description">Description</label><input type="text" name="description" id="id_description" maxlength="20" required></div>`,
		},
		{
			name: "char field with min and max > 0",
			field: func() *aform.CharField {
				f, _ := aform.NewCharField("Description", "", "", 3, 99)
				return f
			}(),
			want: `<div><label for="id_description">Description</label><input type="text" name="description" id="id_description" minlength="3" maxlength="99" required></div>`,
		},
		{
			name: "char field with textarea widget, min and max > 0",
			field: func() *aform.CharField {
				f, _ := aform.NewCharField("Description", "", "", 3, 99, aform.WithWidget(aform.TextArea))
				return f
			}(),
			want: `<div><label for="id_description">Description</label><textarea name="description" id="id_description" minlength="3" maxlength="99" required>
</textarea></div>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.AsDiv(); got != tt.want {
				t.Errorf("AsDiv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_AsDiv_email(t *testing.T) {
	tests := []struct {
		name  string
		field *aform.EmailField
		want  template.HTML
	}{
		{
			name: "bug email field with auto ID disabled and invalid input",
			field: func() *aform.EmailField {
				fld := aform.Must(aform.DefaultEmailField("Your Email"))
				_ = fld.SetAutoID("")
				fld.Clean("not an email")
				return fld
			}(),
			want: `<div>Your Email
<ul class="errorlist"><li>Enter a valid email address</li></ul>
<input type="email" name="your_email" value="not an email" maxlength="254" aria-invalid="true" required></div>`,
		},
		{
			name: "bug email field with auto ID disabled and help text",
			field: func() *aform.EmailField {
				fld := aform.Must(aform.DefaultEmailField("Your Email"))
				_ = fld.SetAutoID("")
				fld.SetHelpText("Enter your email")
				return fld
			}(),
			want: `<div>Your Email<input type="email" name="your_email" maxlength="254" required>
<span class="helptext">Enter your email</span></div>`,
		},
		{
			name: "email field not required with empty value",
			field: func() *aform.EmailField {
				fld := aform.Must(aform.DefaultEmailField("Your email", aform.IsNotRequired()))
				fld.Clean("")
				return fld
			}(),
			want: `<div><label for="id_your_email">Your email</label><input type="email" name="your_email" id="id_your_email" maxlength="254"></div>`,
		},
		{
			name: "email field not required with empty value and custom empty value",
			field: func() *aform.EmailField {
				fld := aform.Must(aform.NewEmailField("Your email", "", "catchall@domain.com", 0, 0, aform.IsNotRequired()))
				fld.Clean("")
				return fld
			}(),
			want: `<div><label for="id_your_email">Your email</label><input type="email" name="your_email" id="id_your_email"></div>`,
		},
		{
			name: "email field with max > 0",
			field: func() *aform.EmailField {
				return aform.Must(aform.NewEmailField("Your email", "", "", 0, 20))
			}(),
			want: `<div><label for="id_your_email">Your email</label><input type="email" name="your_email" id="id_your_email" maxlength="20" required></div>`,
		},
		{
			name: "email field with min and max > 0",
			field: func() *aform.EmailField {
				return aform.Must(aform.NewEmailField("Your email", "", "", 3, 99))
			}(),
			want: `<div><label for="id_your_email">Your email</label><input type="email" name="your_email" id="id_your_email" minlength="3" maxlength="99" required></div>`,
		},
		{
			name: "email field with invalid input",
			field: func() *aform.EmailField {
				fld := aform.Must(aform.DefaultEmailField("Your Email"))
				fld.Clean("not an email")
				return fld
			}(),
			want: `<div><label for="id_your_email">Your Email</label>
<ul class="errorlist"><li id="err_0_id_your_email">Enter a valid email address</li></ul>
<input type="email" name="your_email" value="not an email" id="id_your_email" maxlength="254" aria-describedby="err_0_id_your_email" aria-invalid="true" required></div>`,
		},
		{
			name: "email field with fr locale and empty string",
			field: func() *aform.EmailField {
				fld := aform.Must(aform.DefaultEmailField("Your Email"))
				fld.SetLocale(language.French)
				fld.Clean("")
				return fld
			}(),
			want: `<div><label for="id_your_email">Your Email</label>
<ul class="errorlist"><li id="err_0_id_your_email">Ce champ est obligatoire</li></ul>
<input type="email" name="your_email" id="id_your_email" maxlength="254" aria-describedby="err_0_id_your_email" aria-invalid="true" required></div>`,
		},
		{
			name: "email field with fr locale and invalid input",
			field: func() *aform.EmailField {
				fld := aform.Must(aform.DefaultEmailField("Your Email"))
				fld.SetLocale(language.French)
				fld.Clean("not an email")
				return fld
			}(),
			want: `<div><label for="id_your_email">Your Email</label>
<ul class="errorlist"><li id="err_0_id_your_email">Entrez une adresse e-mail valide</li></ul>
<input type="email" name="your_email" value="not an email" id="id_your_email" maxlength="254" aria-describedby="err_0_id_your_email" aria-invalid="true" required></div>`,
		},
		{
			name: "email field with fr and min/max > 0",
			field: func() *aform.EmailField {
				fld, _ := aform.NewEmailField("Your email", "", "", 8, 12)
				fld.SetLocale(language.French)
				fld.Clean("a@g.com")
				return fld
			}(),
			want: `<div><label for="id_your_email">Your email</label>
<ul class="errorlist"><li id="err_0_id_your_email">Assurez-vous que cette valeur fait au minimum 8 caractères</li></ul>
<input type="email" name="your_email" value="a@g.com" id="id_your_email" minlength="8" maxlength="12" aria-describedby="err_0_id_your_email" aria-invalid="true" required></div>`,
		},
		{
			name: "email field with fr and min/max > 0",
			field: func() *aform.EmailField {
				fld, _ := aform.NewEmailField("Your email", "", "", 8, 12)
				fld.SetLocale(language.French)
				fld.Clean("abcdefg@g.com")
				return fld
			}(),
			want: `<div><label for="id_your_email">Your email</label>
<ul class="errorlist"><li id="err_0_id_your_email">Assurez-vous que cette valeur fait au maximum 12 caractères</li></ul>
<input type="email" name="your_email" value="abcdefg@g.com" id="id_your_email" minlength="8" maxlength="12" aria-describedby="err_0_id_your_email" aria-invalid="true" required></div>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.AsDiv(); got != tt.want {
				t.Errorf("AsDiv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_AsDiv_choice(t *testing.T) {
	tests := []struct {
		name  string
		field *aform.ChoiceField
		want  template.HTML
	}{
		{
			name:  "choice field without options",
			field: aform.Must(aform.DefaultChoiceField("Colors")),
			want: `<div><label for="id_colors">Colors</label><select name="colors" id="id_colors" required>
</select></div>`,
		},
		{
			name: "choice field with options",
			field: aform.Must(aform.DefaultChoiceField("Colors", aform.WithChoiceOptions([]aform.ChoiceFieldOption{
				{Value: "red", Label: "Rouge"},
				{Value: "green", Label: "Vert"},
				{Value: "blue", Label: "Bleu"},
			}), aform.WithGroupedChoiceOptions("Fancy", []aform.ChoiceFieldOption{
				{Value: "orange", Label: "Orange"},
				{Value: "purple", Label: "Purple"},
				{Value: "pink", Label: "Pink"},
			}), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
				{Value: "white", Label: "Blanc"},
				{Value: "black", Label: "Noir"},
			}))),
			want: `<div><label for="id_colors">Colors</label><select name="colors" id="id_colors" required>
  <option value="red" id="id_colors_0">Rouge</option>
  <option value="green" id="id_colors_1">Vert</option>
  <option value="blue" id="id_colors_2">Bleu</option>
  <optgroup label="Fancy">
    <option value="orange" id="id_colors_3_0">Orange</option>
    <option value="purple" id="id_colors_3_1">Purple</option>
    <option value="pink" id="id_colors_3_2">Pink</option>
  </optgroup>
  <option value="white" id="id_colors_4">Blanc</option>
  <option value="black" id="id_colors_5">Noir</option>
</select></div>`,
		},
		{
			name: "choice field with three options and first one selected",
			field: func(selected string) *aform.ChoiceField {
				fld := aform.Must(aform.DefaultChoiceField("Colors", aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				fld.Clean(selected)
				return fld
			}("red"),
			want: `<div><label for="id_colors">Colors</label><select name="colors" id="id_colors" required>
  <option value="red" id="id_colors_0" selected>Rouge</option>
  <option value="green" id="id_colors_1">Vert</option>
  <option value="blue" id="id_colors_2">Bleu</option>
</select></div>`,
		},
		{
			name: "choice field with three options and second one selected",
			field: func(selected string) *aform.ChoiceField {
				fld := aform.Must(aform.DefaultChoiceField("Colors", aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				fld.Clean(selected)
				return fld
			}("green"),
			want: `<div><label for="id_colors">Colors</label><select name="colors" id="id_colors" required>
  <option value="red" id="id_colors_0">Rouge</option>
  <option value="green" id="id_colors_1" selected>Vert</option>
  <option value="blue" id="id_colors_2">Bleu</option>
</select></div>`,
		},
		{
			name: "choice field with three options and third one selected",
			field: func(selected string) *aform.ChoiceField {
				fld := aform.Must(aform.DefaultChoiceField("Colors", aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				fld.Clean(selected)
				return fld
			}("blue"),
			want: `<div><label for="id_colors">Colors</label><select name="colors" id="id_colors" required>
  <option value="red" id="id_colors_0">Rouge</option>
  <option value="green" id="id_colors_1">Vert</option>
  <option value="blue" id="id_colors_2" selected>Bleu</option>
</select></div>`,
		},
		{
			name: "choice field with three options and invalid option selected",
			field: func(selected string) *aform.ChoiceField {
				fld := aform.Must(aform.DefaultChoiceField("Colors", aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				fld.Clean(selected)
				return fld
			}("not_an_option"),
			want: `<div><label for="id_colors">Colors</label>
<ul class="errorlist"><li id="err_0_id_colors">Invalid choice</li></ul>
<select name="colors" id="id_colors" aria-describedby="err_0_id_colors" aria-invalid="true" required>
  <option value="red" id="id_colors_0">Rouge</option>
  <option value="green" id="id_colors_1">Vert</option>
  <option value="blue" id="id_colors_2">Bleu</option>
</select></div>`,
		},
		{
			name: "choice field with options and second one selected with initial",
			field: func() *aform.ChoiceField {
				fld := aform.Must(aform.NewChoiceField("Colors", "green", aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				return fld
			}(),
			want: `<div><label for="id_colors">Colors</label><select name="colors" id="id_colors" required>
  <option value="red" id="id_colors_0">Rouge</option>
  <option value="green" id="id_colors_1" selected>Vert</option>
  <option value="blue" id="id_colors_2">Bleu</option>
</select></div>`,
		},
		{
			name: "choice field with options and invalid initial",
			field: func() *aform.ChoiceField {
				fld := aform.Must(aform.NewChoiceField("Colors", "not_an_option", aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				return fld
			}(),
			want: `<div><label for="id_colors">Colors</label><select name="colors" id="id_colors" required>
  <option value="red" id="id_colors_0">Rouge</option>
  <option value="green" id="id_colors_1">Vert</option>
  <option value="blue" id="id_colors_2">Bleu</option>
</select></div>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.AsDiv(); got != tt.want {
				t.Errorf("AsDiv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_AsDiv_choiceRadio(t *testing.T) {
	tests := []struct {
		name  string
		field *aform.ChoiceField
		want  template.HTML
	}{
		{
			name:  "choice field without options",
			field: aform.Must(aform.DefaultChoiceField("Colors", aform.WithWidget(aform.RadioSelect))),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<div id="id_colors">
</div>
</fieldset>
</div>`,
		},
		{
			name: "choice field with options",
			field: aform.Must(aform.DefaultChoiceField("Colors", aform.WithWidget(aform.RadioSelect), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
				{Value: "red", Label: "Rouge"},
				{Value: "green", Label: "Vert"},
				{Value: "blue", Label: "Bleu"},
			}), aform.WithGroupedChoiceOptions("Fancy", []aform.ChoiceFieldOption{
				{Value: "orange", Label: "Orange"},
				{Value: "purple", Label: "Purple"},
				{Value: "pink", Label: "Pink"},
			}), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
				{Value: "white", Label: "Blanc"},
				{Value: "black", Label: "Noir"},
			}))),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<div id="id_colors">
<label for="id_colors_0"><input type="radio" name="colors" value="red" id="id_colors_0">Rouge</label>
<label for="id_colors_1"><input type="radio" name="colors" value="green" id="id_colors_1">Vert</label>
<label for="id_colors_2"><input type="radio" name="colors" value="blue" id="id_colors_2">Bleu</label>
<div><label>Fancy</label>
<label for="id_colors_3_0"><input type="radio" name="colors" value="orange" id="id_colors_3_0">Orange</label>
<label for="id_colors_3_1"><input type="radio" name="colors" value="purple" id="id_colors_3_1">Purple</label>
<label for="id_colors_3_2"><input type="radio" name="colors" value="pink" id="id_colors_3_2">Pink</label>
</div>
<label for="id_colors_4"><input type="radio" name="colors" value="white" id="id_colors_4">Blanc</label>
<label for="id_colors_5"><input type="radio" name="colors" value="black" id="id_colors_5">Noir</label>
</div>
</fieldset>
</div>`,
		},
		{
			name: "choice field with three options and first one selected",
			field: func(selected string) *aform.ChoiceField {
				fld := aform.Must(aform.DefaultChoiceField("Colors", aform.WithWidget(aform.RadioSelect), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				fld.Clean(selected)
				return fld
			}("red"),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<div id="id_colors">
<label for="id_colors_0"><input type="radio" name="colors" value="red" id="id_colors_0" checked>Rouge</label>
<label for="id_colors_1"><input type="radio" name="colors" value="green" id="id_colors_1">Vert</label>
<label for="id_colors_2"><input type="radio" name="colors" value="blue" id="id_colors_2">Bleu</label>
</div>
</fieldset>
</div>`,
		},
		{
			name: "choice field with three options and second one selected",
			field: func(selected string) *aform.ChoiceField {
				fld := aform.Must(aform.DefaultChoiceField("Colors", aform.WithWidget(aform.RadioSelect), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				fld.Clean(selected)
				return fld
			}("green"),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<div id="id_colors">
<label for="id_colors_0"><input type="radio" name="colors" value="red" id="id_colors_0">Rouge</label>
<label for="id_colors_1"><input type="radio" name="colors" value="green" id="id_colors_1" checked>Vert</label>
<label for="id_colors_2"><input type="radio" name="colors" value="blue" id="id_colors_2">Bleu</label>
</div>
</fieldset>
</div>`,
		},
		{
			name: "choice field with three options and third one selected",
			field: func(selected string) *aform.ChoiceField {
				fld := aform.Must(aform.DefaultChoiceField("Colors", aform.WithWidget(aform.RadioSelect), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				fld.Clean(selected)
				return fld
			}("blue"),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<div id="id_colors">
<label for="id_colors_0"><input type="radio" name="colors" value="red" id="id_colors_0">Rouge</label>
<label for="id_colors_1"><input type="radio" name="colors" value="green" id="id_colors_1">Vert</label>
<label for="id_colors_2"><input type="radio" name="colors" value="blue" id="id_colors_2" checked>Bleu</label>
</div>
</fieldset>
</div>`,
		},
		{
			name: "choice field with three options and invalid option selected",
			field: func(selected string) *aform.ChoiceField {
				fld := aform.Must(aform.DefaultChoiceField("Colors", aform.WithWidget(aform.RadioSelect), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				fld.Clean(selected)
				return fld
			}("not_an_option"),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<ul class="errorlist"><li id="err_0_id_colors">Invalid choice</li></ul>
<div id="id_colors">
<label for="id_colors_0"><input type="radio" name="colors" value="red" id="id_colors_0">Rouge</label>
<label for="id_colors_1"><input type="radio" name="colors" value="green" id="id_colors_1">Vert</label>
<label for="id_colors_2"><input type="radio" name="colors" value="blue" id="id_colors_2">Bleu</label>
</div>
</fieldset>
</div>`,
		},
		{
			name: "choice field with options and second one selected with initial",
			field: func() *aform.ChoiceField {
				fld := aform.Must(aform.NewChoiceField("Colors", "green", aform.WithWidget(aform.RadioSelect), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				return fld
			}(),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<div id="id_colors">
<label for="id_colors_0"><input type="radio" name="colors" value="red" id="id_colors_0">Rouge</label>
<label for="id_colors_1"><input type="radio" name="colors" value="green" id="id_colors_1" checked>Vert</label>
<label for="id_colors_2"><input type="radio" name="colors" value="blue" id="id_colors_2">Bleu</label>
</div>
</fieldset>
</div>`,
		},
		{
			name: "choice field with options and invalid initial",
			field: func() *aform.ChoiceField {
				fld := aform.Must(aform.NewChoiceField("Colors", "not_an_option", aform.WithWidget(aform.RadioSelect), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				return fld
			}(),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<div id="id_colors">
<label for="id_colors_0"><input type="radio" name="colors" value="red" id="id_colors_0">Rouge</label>
<label for="id_colors_1"><input type="radio" name="colors" value="green" id="id_colors_1">Vert</label>
<label for="id_colors_2"><input type="radio" name="colors" value="blue" id="id_colors_2">Bleu</label>
</div>
</fieldset>
</div>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.AsDiv(); got != tt.want {
				t.Errorf("AsDiv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_AsDiv_choiceCheckboxSelectMultiple(t *testing.T) {
	tests := []struct {
		name  string
		field *aform.ChoiceField
		want  template.HTML
	}{
		{
			name:  "choice field without options",
			field: aform.Must(aform.DefaultChoiceField("Colors", aform.WithWidget(aform.CheckboxSelectMultiple))),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<div id="id_colors">
</div>
</fieldset>
</div>`,
		},
		{
			name: "choice field with options",
			field: aform.Must(aform.DefaultChoiceField("Colors", aform.WithWidget(aform.CheckboxSelectMultiple), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
				{Value: "red", Label: "Rouge"},
				{Value: "green", Label: "Vert"},
				{Value: "blue", Label: "Bleu"},
			}), aform.WithGroupedChoiceOptions("Fancy", []aform.ChoiceFieldOption{
				{Value: "orange", Label: "Orange"},
				{Value: "purple", Label: "Purple"},
				{Value: "pink", Label: "Pink"},
			}), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
				{Value: "white", Label: "Blanc"},
				{Value: "black", Label: "Noir"},
			}))),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<div id="id_colors">
<label for="id_colors_0"><input type="checkbox" name="colors" value="red" id="id_colors_0">Rouge</label>
<label for="id_colors_1"><input type="checkbox" name="colors" value="green" id="id_colors_1">Vert</label>
<label for="id_colors_2"><input type="checkbox" name="colors" value="blue" id="id_colors_2">Bleu</label>
<div><label>Fancy</label>
<label for="id_colors_3_0"><input type="checkbox" name="colors" value="orange" id="id_colors_3_0">Orange</label>
<label for="id_colors_3_1"><input type="checkbox" name="colors" value="purple" id="id_colors_3_1">Purple</label>
<label for="id_colors_3_2"><input type="checkbox" name="colors" value="pink" id="id_colors_3_2">Pink</label>
</div>
<label for="id_colors_4"><input type="checkbox" name="colors" value="white" id="id_colors_4">Blanc</label>
<label for="id_colors_5"><input type="checkbox" name="colors" value="black" id="id_colors_5">Noir</label>
</div>
</fieldset>
</div>`,
		},
		{
			name: "choice field with three options and first one selected",
			field: func(selected string) *aform.ChoiceField {
				fld := aform.Must(aform.DefaultChoiceField("Colors", aform.WithWidget(aform.CheckboxSelectMultiple), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				fld.Clean(selected)
				return fld
			}("red"),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<div id="id_colors">
<label for="id_colors_0"><input type="checkbox" name="colors" value="red" id="id_colors_0" checked>Rouge</label>
<label for="id_colors_1"><input type="checkbox" name="colors" value="green" id="id_colors_1">Vert</label>
<label for="id_colors_2"><input type="checkbox" name="colors" value="blue" id="id_colors_2">Bleu</label>
</div>
</fieldset>
</div>`,
		},
		{
			name: "choice field with three options and second one selected",
			field: func(selected string) *aform.ChoiceField {
				fld := aform.Must(aform.DefaultChoiceField("Colors", aform.WithWidget(aform.CheckboxSelectMultiple), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				fld.Clean(selected)
				return fld
			}("green"),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<div id="id_colors">
<label for="id_colors_0"><input type="checkbox" name="colors" value="red" id="id_colors_0">Rouge</label>
<label for="id_colors_1"><input type="checkbox" name="colors" value="green" id="id_colors_1" checked>Vert</label>
<label for="id_colors_2"><input type="checkbox" name="colors" value="blue" id="id_colors_2">Bleu</label>
</div>
</fieldset>
</div>`,
		},
		{
			name: "choice field with three options and third one selected",
			field: func(selected string) *aform.ChoiceField {
				fld := aform.Must(aform.DefaultChoiceField("Colors", aform.WithWidget(aform.CheckboxSelectMultiple), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				fld.Clean(selected)
				return fld
			}("blue"),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<div id="id_colors">
<label for="id_colors_0"><input type="checkbox" name="colors" value="red" id="id_colors_0">Rouge</label>
<label for="id_colors_1"><input type="checkbox" name="colors" value="green" id="id_colors_1">Vert</label>
<label for="id_colors_2"><input type="checkbox" name="colors" value="blue" id="id_colors_2" checked>Bleu</label>
</div>
</fieldset>
</div>`,
		},
		{
			name: "choice field with three options and invalid option selected",
			field: func(selected string) *aform.ChoiceField {
				fld := aform.Must(aform.DefaultChoiceField("Colors", aform.WithWidget(aform.CheckboxSelectMultiple), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				fld.Clean(selected)
				return fld
			}("not_an_option"),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<ul class="errorlist"><li id="err_0_id_colors">Invalid choice</li></ul>
<div id="id_colors">
<label for="id_colors_0"><input type="checkbox" name="colors" value="red" id="id_colors_0">Rouge</label>
<label for="id_colors_1"><input type="checkbox" name="colors" value="green" id="id_colors_1">Vert</label>
<label for="id_colors_2"><input type="checkbox" name="colors" value="blue" id="id_colors_2">Bleu</label>
</div>
</fieldset>
</div>`,
		},
		{
			name: "choice field with options and second one selected with initial",
			field: func() *aform.ChoiceField {
				fld := aform.Must(aform.NewChoiceField("Colors", "green", aform.WithWidget(aform.CheckboxSelectMultiple), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				return fld
			}(),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<div id="id_colors">
<label for="id_colors_0"><input type="checkbox" name="colors" value="red" id="id_colors_0">Rouge</label>
<label for="id_colors_1"><input type="checkbox" name="colors" value="green" id="id_colors_1" checked>Vert</label>
<label for="id_colors_2"><input type="checkbox" name="colors" value="blue" id="id_colors_2">Bleu</label>
</div>
</fieldset>
</div>`,
		},
		{
			name: "choice field with options and invalid initial",
			field: func() *aform.ChoiceField {
				fld := aform.Must(aform.NewChoiceField("Colors", "not_an_option", aform.WithWidget(aform.CheckboxSelectMultiple), aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "red", Label: "Rouge"},
					{Value: "green", Label: "Vert"},
					{Value: "blue", Label: "Bleu"},
				})))
				return fld
			}(),
			want: `<div>
<fieldset><legend for="id_colors">Colors</legend>
<div id="id_colors">
<label for="id_colors_0"><input type="checkbox" name="colors" value="red" id="id_colors_0">Rouge</label>
<label for="id_colors_1"><input type="checkbox" name="colors" value="green" id="id_colors_1">Vert</label>
<label for="id_colors_2"><input type="checkbox" name="colors" value="blue" id="id_colors_2">Bleu</label>
</div>
</fieldset>
</div>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.AsDiv(); got != tt.want {
				t.Errorf("AsDiv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_LabelTag(t *testing.T) {
	tests := []struct {
		name string
		fld  *aform.CharField
		want template.HTML
	}{
		{
			name: "name",
			fld:  aform.Must(aform.DefaultCharField("name")),
			want: template.HTML(`<label for="id_name">name</label>`),
		},
		{
			// https://html.spec.whatwg.org/multipage/syntax.html#syntax-ambiguous-ampersand
			name: "with ambiguous ampersand	&123;",
			fld:  aform.Must(aform.DefaultCharField("name&123;")),
			want: template.HTML(`<label for="id_name&123;">name&amp;123;</label>`),
		},
		{
			// https://html.spec.whatwg.org/multipage/syntax.html#syntax-ambiguous-ampersand
			name: "with not ambiguous ampersand	&123",
			fld:  aform.Must(aform.DefaultCharField("name&123")),
			want: template.HTML(`<label for="id_name&123">name&amp;123</label>`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, tt.fld.LabelTag(), "LabelTag()")
		})
	}
}
