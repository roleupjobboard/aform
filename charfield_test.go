package aform_test

import (
	"fmt"
	"github.com/roleupjobboard/aform"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

func TestDefaultCharField_maxLength(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultCharField("test")
	a.NoError(err)
	s := "NYvkTgiehbpyjrqrjqssz7BJa2HUoEcVn4UCF8hVY70pGDBaw9h0tJDPcuRWRNloiwAEIKTuuKP8poSnUiqlX8Ms9tpjOt2GxNc9RocglpzV5q64HUuEPf5LAo18EmfyqEC486rzNbI5IyFX0x3VI8nc5MWi01feZcCZxwz1cschhYjodpOjYw8dIczSujE6j5zzGOKhOB9XaPZosVAE3FKaXkTQ67ecZwiK7PmNrmkNWtn0l068Wk0ZXf7ioR6o"
	actual, errs := f.Clean(s)
	a.Len(errs, 0)
	a.Equal(s, actual)
	ss := s + "z"
	actual, errs = f.Clean(ss)
	a.Len(errs, 1)
	a.Equal(errs[0].Error(), "Ensure this value has at most 256 characters")
	a.Equal(ss, actual)
}

func TestCharField_Clean_withValidValue(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultCharField("test")
	a.NoError(err)
	actual, errs := f.Clean("boundValue")
	a.Len(errs, 0)
	a.Equal("boundValue", actual)
}

func TestCharField_Clean_withValidValueSurroundingBySpaces(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultCharField("test")
	a.NoError(err)
	actual, errs := f.Clean("  Bound Value ")
	a.Len(errs, 0)
	a.Equal("Bound Value", actual)
}

func TestCharField_Clean_withValidValueIncludingHTML(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultCharField("test")
	a.NoError(err)
	actual, errs := f.Clean(`Bound <a onblur="alert(secret)" href="https://en.wikipedia.org/wiki/Goodbye_Cruel_World_(Pink_Floyd_song)">HTML</a> value`)
	a.Len(errs, 0)
	a.Equal("Bound HTML value", actual)
}

func TestCharField_Clean_withValidMultilineValue(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultCharField("test")
	a.NoError(err)
	actual, errs := f.Clean(`Bound value:
- 1
- 2`)
	a.Len(errs, 0)
	a.Equal("Bound value: - 1 - 2", actual, "inputs remove new lines if the input expected should be on one line")
}

func TestCharField_Clean_withValidMultilineValueAndTextAreaWidget(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultCharField("test", aform.WithWidget(aform.TextArea))
	a.NoError(err)
	expected := `Bound value:
- 1
- 2`
	actual, errs := f.Clean(expected)
	a.Len(errs, 0)
	a.Equal(expected, actual, "inputs with TextArea widget keep new lines")
}

func TestCharField_Clean_withValidMultilineHTMLValueAndTextAreaWidget(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultCharField("test", aform.WithWidget(aform.TextArea))
	a.NoError(err)
	value := `Bound
<a onblur="alert(secret)" href="https://en.wikipedia.org/wiki/Goodbye_Cruel_World_(Pink_Floyd_song)">HTML</a>
value`
	actual, errs := f.Clean(value)
	a.Len(errs, 0)
	expected := `Bound
HTML
value`
	a.Equal(expected, actual, "inputs with TextArea widget keep new lines")
}

func TestCharField_Clean_withEmptyValue(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultCharField("test")
	a.NoError(err)
	actual, errs := f.Clean("")
	a.Len(errs, 1)
	a.Equal("This field is required", errs[0].Error())
	a.Equal("", actual)
}

func TestCharField_Clean_withValueTooShort(t *testing.T) {
	a := assert.New(t)
	f, err := aform.NewCharField("test", "", "", 6, 99)
	a.NoError(err)
	actual, errs := f.Clean("short")
	a.Len(errs, 1)
	a.Equal("Ensure this value has at least 6 characters", errs[0].Error())
	a.Equal("short", actual)
}

func TestCharField_Clean_withValueTooLong(t *testing.T) {
	a := assert.New(t)
	f, err := aform.NewCharField("test", "", "", 2, 7)
	a.NoError(err)
	actual, errs := f.Clean("too long")
	a.Len(errs, 1)
	a.Equal("Ensure this value has at most 7 characters", errs[0].Error())
	a.Equal("too long", actual)
}

func TestCharField_Clean_withEmptyValueAndNotRequired(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultCharField("test", aform.IsNotRequired())
	a.NoError(err)
	actual, errs := f.Clean("")
	a.Len(errs, 0)
	a.Equal("", actual, "it should use the default empty value")
}

func TestCharField_Clean_withEmptyValueNotRequiredAndCustomEmptyValue(t *testing.T) {
	a := assert.New(t)
	f, err := aform.NewCharField("test", "", "empty_custom", 0, 0, aform.IsNotRequired())
	a.NoError(err)
	actual, errs := f.Clean("")
	a.Len(errs, 0)
	a.Equal("empty_custom", actual, "it should use the custom empty value")
}

func TestCharField_CustomizeError_changeMessageRequired(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultCharField("name")
	a.NoError(err)
	f.CustomizeError(aform.ErrorWrapWithCode(fmt.Errorf("Please, enter your name"), aform.RequiredErrorCode))
	_, errs := f.Clean("")
	a.Len(errs, 1)
	a.Equal(aform.RequiredErrorCode, errs[0].Code())
	a.Equal("Please, enter your name", errs[0].Error())
}

func TestCharField_CustomizeError_changeMessageMinLength(t *testing.T) {
	a := assert.New(t)
	f, err := aform.NewCharField("name", "", "", 6, 99)
	a.NoError(err)
	f.CustomizeError(aform.ErrorWrapWithCode(fmt.Errorf("Please, enter enough characters"), aform.MinLengthErrorCode))
	_, errs := f.Clean("short")
	a.Len(errs, 1)
	a.Equal(aform.MinLengthErrorCode, errs[0].Code())
	a.Equal("Please, enter enough characters", errs[0].Error())
}

type onlyPaulError struct{}

func (e onlyPaulError) Error() string {
	return e.Translate("en")
}

func (e onlyPaulError) Code() string {
	return "notpaul"
}

func (e onlyPaulError) Translate(locale string) string {
	switch locale {
	case language.English.String():
		return "Only Pauls are allowed"
	case language.French.String():
		return "Seul les Pauls sont autorisés"
	default:
		return "Only Pauls are allowed"
	}
}

func TestCharField_SetValidateFunc_addNewError(t *testing.T) {
	var e onlyPaulError
	a := assert.New(t)
	f, err := aform.DefaultCharField("name")
	a.NoError(err)
	f.SetValidateFunc(func(current aform.ValidationFunc) (new aform.ValidationFunc) {
		return func(s string, b bool) []aform.Error {
			if s != "Paul" {
				var e onlyPaulError
				return []aform.Error{aform.ErrorWrap(e)}
			}
			return nil
		}
	})
	_, errs := f.Clean("Pierre")
	a.Len(errs, 1)
	a.Equal(e.Code(), errs[0].Code())
	a.Equal("Only Pauls are allowed", errs[0].Translate(language.English.String()))
	a.Equal("Seul les Pauls sont autorisés", errs[0].Translate(language.French.String()))
}
