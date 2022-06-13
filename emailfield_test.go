package aform_test

import (
	"fmt"
	"github.com/roleupjobboard/aform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultEmailField_maxLength(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultEmailField("test")
	a.NoError(err)
	e := "NYvkTgiehbpyjrqrjqssz7BJa2HUoEcVn4UCF8hVY70pGDBaw9h0t@DPcuRWRNloiwAEIKTuuKP8poSnUiqlX8Ms9tpjOt2GxNc9RocglpzV5q64HUuEPf5LAo18EmfyqEC486rzNbI5IyFX0x3VI8nc5MWi01feZcCZxwz1cschhYjodpOjYw8dIczSujE6j5zzGOKhOB9XaPZosVAE3FKaXkTQ67ecZwiK7PmNrmkNWtn0l068Wk0ZXf.com"
	actual, errs := f.Clean(e)
	a.Len(errs, 0)
	a.Equal(e, actual)
	ee := "z" + e
	actual, errs = f.Clean(ee)
	a.Len(errs, 1)
	a.Equal(errs[0].Error(), "Ensure this value has at most 254 characters")
	a.Equal(ee, actual)
}

func TestEmailField_Clean_withValidValue(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultEmailField("test")
	a.NoError(err)
	actual, errs := f.Clean("g@test.com")
	a.Len(errs, 0)
	a.Equal("g@test.com", actual)
}

func TestEmailField_Clean_withValueSurroundingBySpaces(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultEmailField("test")
	a.NoError(err)
	actual, errs := f.Clean("  g@test.com ")
	a.Len(errs, 0)
	a.Equal("g@test.com", actual)
}

func TestEmailField_Clean_withEmptyValue(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultEmailField("test")
	a.NoError(err)
	actual, errs := f.Clean("")
	a.Len(errs, 1)
	a.Equal("", actual)
}

func TestEmailField_Clean_withValueTooShort(t *testing.T) {
	a := assert.New(t)
	f, err := aform.NewEmailField("test", "", "", 8, 99)
	a.NoError(err)
	actual, errs := f.Clean("g@g.com")
	a.Len(errs, 1)
	a.Equal("Ensure this value has at least 8 characters", errs[0].Error())
	a.Equal("g@g.com", actual)
}

func TestEmailField_Clean_withValueTooLong(t *testing.T) {
	a := assert.New(t)
	f, err := aform.NewEmailField("test", "", "", 2, 6)
	a.NoError(err)
	actual, errs := f.Clean("g@g.com")
	a.Len(errs, 1)
	a.Equal("Ensure this value has at most 6 characters", errs[0].Error())
	a.Equal("g@g.com", actual)
}

func TestEmailField_Clean_withInvalidValue(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultEmailField("test")
	a.NoError(err)
	actual, errs := f.Clean("not email")
	a.Len(errs, 1)
	a.Equal("Enter a valid email address", errs[0].Error())
	a.Equal("not email", actual)
}

func TestEmailField_Clean_withInvalidValueSurroundingBySpaces(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultEmailField("test")
	a.NoError(err)
	actual, errs := f.Clean(" not an email  ")
	a.Len(errs, 1)
	a.Equal("not an email", actual)
}

func TestEmailField_Clean_withEmptyValueAndNotRequired(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultEmailField("test", aform.IsNotRequired())
	a.NoError(err)
	actual, errs := f.Clean("")
	a.Len(errs, 0)
	a.Equal("", actual)
}

func TestEmailField_Clean_withEmptyValueNotRequiredAndCustomEmptyValue(t *testing.T) {
	a := assert.New(t)
	f, err := aform.NewEmailField("test", "", "catchall@domain.com", 0, 0, aform.IsNotRequired())
	a.NoError(err)
	actual, errs := f.Clean("")
	a.Len(errs, 0)
	a.Equal("catchall@domain.com", actual)
}

func TestEmailField_CustomizeError_changeMessageEmail(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultEmailField("test")
	a.NoError(err)
	f.CustomizeError(aform.ErrorWrapWithCode(fmt.Errorf("Please, enter your email"), aform.EmailErrorCode))
	_, errs := f.Clean("")
	a.Len(errs, 1)
	a.Equal(aform.RequiredErrorCode, errs[0].Code())
	a.Equal("This field is required", errs[0].Error())
	_, errs = f.Clean("not an email but present")
	a.Len(errs, 1)
	a.Equal(aform.EmailErrorCode, errs[0].Code())
	a.Equal("Please, enter your email", errs[0].Error())
}

func TestEmailField_MustEmail_withValidEmail(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultEmailField("test")
	a.NoError(err)
	a.Equal("g@gmail.com", f.MustEmail("g@gmail.com"))
}

func TestEmailField_MustEmail_withInvalidEmail(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultEmailField("test")
	a.NoError(err)
	a.PanicsWithValue("MustEmail called on test field with an invalid email value: invalid", func() {
		f.MustEmail("invalid")
	})
}
