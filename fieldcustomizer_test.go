package aform_test

import (
	"fmt"
	"github.com/roleupjobboard/aform"
	"github.com/stretchr/testify/assert"
	"html/template"
	"strings"
	"testing"
)

func TestField_WithAttributes(t *testing.T) {
	tests := []struct {
		name  string
		field *aform.CharField
		want  template.HTML
	}{
		{
			name:  "with a boolean attribute",
			field: aform.Must(aform.DefaultCharField("Your Name", aform.WithAttributes([]aform.Attributable{aform.BoolAttr("newattr")}))),
			want:  `<div><label for="id_your_name">Your Name</label><input type="text" name="your_name" id="id_your_name" maxlength="256" newattr required></div>`,
		},
		{
			name:  "with a boolean attribute that collision an attribute it replaces it",
			field: aform.Must(aform.DefaultCharField("Your Name", aform.WithAttributes([]aform.Attributable{aform.BoolAttr("id")}))),
			want:  `<div><label for="id_your_name">Your Name</label><input type="text" name="your_name" id maxlength="256" required></div>`,
		},
		{
			name:  "with multiple boolean attributes",
			field: aform.Must(aform.DefaultCharField("Your Name", aform.WithAttributes([]aform.Attributable{aform.BoolAttr("alpha"), aform.BoolAttr("zeta"), aform.BoolAttr("bravo")}))),
			want:  `<div><label for="id_your_name">Your Name</label><input type="text" name="your_name" id="id_your_name" maxlength="256" alpha bravo required zeta></div>`,
		},
		{
			name:  "with multiple times the same attribute key",
			field: aform.Must(aform.DefaultCharField("Your Name", aform.WithAttributes([]aform.Attributable{aform.BoolAttr("alpha"), aform.Attr("alpha", "alpha value")}))),
			want:  `<div><label for="id_your_name">Your Name</label><input type="text" name="your_name" id="id_your_name" maxlength="256" alpha="alpha value" required></div>`,
		},
		{
			name:  "with multiple times the same attribute key",
			field: aform.Must(aform.DefaultCharField("Your Name", aform.WithAttributes([]aform.Attributable{aform.BoolAttr("alpha"), aform.Attr("alpha", "John Doe"), aform.BoolAttr("alpha")}))),
			want:  `<div><label for="id_your_name">Your Name</label><input type="text" name="your_name" id="id_your_name" maxlength="256" alpha required></div>`,
		},
		{
			name:  "with an attribute",
			field: aform.Must(aform.DefaultCharField("Your Name", aform.WithAttributes([]aform.Attributable{aform.Attr("placeholder", "John Doe")}))),
			want:  `<div><label for="id_your_name">Your Name</label><input type="text" name="your_name" id="id_your_name" maxlength="256" placeholder="John Doe" required></div>`,
		},
		{
			name:  "with an attribute that collision an attribute it replaces it",
			field: aform.Must(aform.DefaultCharField("Your Name", aform.WithAttributes([]aform.Attributable{aform.Attr("id", "id_your_first_name")}))),
			want:  `<div><label for="id_your_name">Your Name</label><input type="text" name="your_name" id="id_your_first_name" maxlength="256" required></div>`,
		},
		{
			name:  "with multiple attributes",
			field: aform.Must(aform.DefaultCharField("Your Name", aform.WithAttributes([]aform.Attributable{aform.Attr("placeholder", "John Doe"), aform.Attr("data-test", "test")}))),
			want:  `<div><label for="id_your_name">Your Name</label><input type="text" name="your_name" id="id_your_name" data-test="test" maxlength="256" placeholder="John Doe" required></div>`,
		},
		{
			name:  "with multiple aform.WithAttributes and multiple attributes",
			field: aform.Must(aform.DefaultCharField("Your Name", aform.WithAttributes([]aform.Attributable{aform.Attr("placeholder", "John Doe"), aform.Attr("data-first", "first")}), aform.WithAttributes([]aform.Attributable{aform.Attr("placeholder", "Jane Doe"), aform.Attr("data-second", "second")}))),
			want:  `<div><label for="id_your_name">Your Name</label><input type="text" name="your_name" id="id_your_name" data-first="first" data-second="second" maxlength="256" placeholder="Jane Doe" required></div>`,
		},
		{
			name: "with class attribute it does not change label classes and it appends to widget error class",
			field: func() *aform.CharField {
				f := aform.Must(aform.DefaultCharField("Your Name", aform.WithAttributes([]aform.Attributable{aform.Attr("class", "big blue")})))
				f.SetRequiredCSSClass("required_class")
				f.SetErrorCSSClass("error_class")
				f.Clean("")
				return f
			}(),
			want: `<div class="required_class error_class"><label class="required_class" for="id_your_name">Your Name</label>
<ul class="errorlist"><li id="err_0_id_your_name">This field is required</li></ul>
<input type="text" name="your_name" class="error_class big blue" id="id_your_name" maxlength="256" aria-describedby="err_0_id_your_name" aria-invalid="true" required></div>`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.AsDiv(); got != tt.want {
				t.Errorf("Field() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_CustomizeError_nonExistingCode(t *testing.T) {
	a := assert.New(t)
	aField := aform.Must(aform.DefaultCharField("Your Field"))
	a.PanicsWithValue(
		"CustomizeError called on Your Field field with a nonexisting Error code: nonexistent",
		func() {
			aField.CustomizeError(aform.ErrorWrapWithCode(fmt.Errorf("anything"), "nonexistent"))
		})
}

type onlyOneError struct{}

func (e onlyOneError) Error() string {
	return "Only one thing accepted"
}

func (e onlyOneError) Code() string {
	return "onlyonething"
}

func (e onlyOneError) Translate(locale string) string {
	return e.Error()
}

func TestField_SetSanitizeFunc(t *testing.T) {
	a := assert.New(t)
	aField := aform.Must(aform.DefaultCharField("Your Field"))
	cleaned, errs := aField.Clean(" some <strong>bold</strong> statement  ")
	a.Len(errs, 0)
	a.Equal(cleaned, "some bold statement")
	aField.SetSanitizeFunc(func(current aform.SanitizationFunc) (new aform.SanitizationFunc) {
		return func(s string) string {
			return strings.TrimSpace(s)
		}
	})
	cleaned, errs = aField.Clean(" some <strong>bold</strong> statement  ")
	a.Len(errs, 0)
	a.Equal(cleaned, "some <strong>bold</strong> statement")
	aField.SetSanitizeFunc(func(current aform.SanitizationFunc) (new aform.SanitizationFunc) {
		return func(s string) string {
			return s
		}
	})
	cleaned, errs = aField.Clean(" some <strong>bold</strong> statement  ")
	a.Len(errs, 0)
	a.Equal(cleaned, " some <strong>bold</strong> statement  ")
}

func TestField_SetSanitizeFunc_withWidgetTypeChanged(t *testing.T) {
	a := assert.New(t)
	aField := aform.Must(aform.DefaultCharField("Your Field"))
	calls := 0
	aField.SetSanitizeFunc(func(current aform.SanitizationFunc) (new aform.SanitizationFunc) {
		return func(s string) string {
			calls++
			return current(s)
		}
	})
	a.Equal(0, calls)
	cleaned, errs := aField.Clean("anything")
	a.Len(errs, 0)
	a.Equal(cleaned, "anything")
	a.Equal(1, calls)
	aField.SetWidget(aform.TextArea)
	cleaned, errs = aField.Clean("anything")
	a.Len(errs, 0)
	a.Equal(cleaned, "anything")
	a.Equal(2, calls)
}

func TestField_SetValidateFunc_withCharField(t *testing.T) {
	a := assert.New(t)
	onlyOneNameField := aform.Must(aform.NewCharField("Your Name", "", "", 5, 0))
	cleaned, errs := onlyOneNameField.Clean("Karine")
	a.Len(errs, 0, "any name longer than 5 characters should work")
	a.Equal(cleaned, "Karine")
	cleaned, errs = onlyOneNameField.Clean("Abygaelle")
	a.Len(errs, 0, "any name longer than 5 characters should work")
	a.Equal(cleaned, "Abygaelle")
	_, errs = onlyOneNameField.Clean("Drew")
	a.Len(errs, 1, "a name shorter than 5 characters should not work")
	onlyOneNameField.SetValidateFunc(func(current aform.ValidationFunc) (new aform.ValidationFunc) {
		return func(s string, b bool) []aform.Error {
			errs := current(s, b)
			if len(errs) > 0 {
				return errs
			}
			if s != "Karine" {
				var e onlyOneError
				return []aform.Error{aform.ErrorWrap(e)}
			}
			return nil
		}
	})
	cleaned, errs = onlyOneNameField.Clean("Karine")
	a.Len(errs, 0, "Karine should work")
	a.Equal(cleaned, "Karine")
	_, errs = onlyOneNameField.Clean("Abygaelle")
	a.Len(errs, 1, "Abygaelle should not work")
	a.Equal("Only one thing accepted", errs[0].Error())
	_, errs = onlyOneNameField.Clean("Drew")
	a.Len(errs, 1, "a name shorter than 5 characters should not work")
	a.Equal("Ensure this value has at least 5 characters", errs[0].Error())
}

func TestField_SetValidateFunc_withEmailField(t *testing.T) {
	a := assert.New(t)
	onlyOneEmailField := aform.Must(aform.DefaultEmailField("Your Email"))
	cleaned, errs := onlyOneEmailField.Clean("x@domain.com")
	a.Len(errs, 0, "any email should work")
	a.Equal(cleaned, "x@domain.com")
	onlyOneEmailField.SetValidateFunc(func(current aform.ValidationFunc) (new aform.ValidationFunc) {
		return func(s string, b bool) []aform.Error {
			if s != "ok@domain.com" {
				var e onlyOneError
				return []aform.Error{aform.ErrorWrap(e)}
			}
			return nil
		}
	})
	_, errs = onlyOneEmailField.Clean("x@domain.com")
	a.Len(errs, 1, "x@domain.com should not work anymore")
	a.Equal("Only one thing accepted", errs[0].Error())
	cleaned, errs = onlyOneEmailField.Clean("ok@domain.com")
	a.Len(errs, 0, "ok@domain.com should work")
	a.Equal(cleaned, "ok@domain.com")
}

func TestField_SetValidateFunc_withWidgetTypeChanged(t *testing.T) {
	a := assert.New(t)
	aField := aform.Must(aform.DefaultCharField("Your Field"))
	calls := 0
	aField.SetValidateFunc(func(current aform.ValidationFunc) (new aform.ValidationFunc) {
		return func(s string, b bool) []aform.Error {
			calls++
			return current(s, b)
		}
	})
	a.Equal(0, calls)
	cleaned, errs := aField.Clean("anything")
	a.Len(errs, 0)
	a.Equal(cleaned, "anything")
	a.Equal(1, calls)
	aField.SetWidget(aform.TextArea)
	cleaned, errs = aField.Clean("anything")
	a.Len(errs, 0)
	a.Equal(cleaned, "anything")
	a.Equal(2, calls)
}
