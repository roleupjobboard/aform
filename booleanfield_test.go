package aform_test

import (
	"github.com/roleupjobboard/aform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBooleanField_Clean_withTrueValueAndRequired(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultBooleanField("test")
	a.NoError(err)
	actual, errs := f.Clean("on")
	a.Len(errs, 0)
	a.Equal("on", actual)
}

func TestBooleanField_Clean_withEmptyValueAndRequired(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultBooleanField("test")
	a.NoError(err)
	actual, errs := f.Clean("")
	a.Len(errs, 1)
	a.Equal("This field is required", errs[0].Error())
	a.Equal("", actual)
}

func TestBooleanField_Clean_withEmptyValueAndNotRequired(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultBooleanField("test", aform.IsNotRequired())
	a.NoError(err)
	actual, errs := f.Clean("")
	a.Len(errs, 0)
	a.Equal("off", actual, "it should return the empty value")
}

func TestBooleanField_MustBoolean_withValidBoolean(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultBooleanField("test")
	a.NoError(err)
	a.True(f.MustBoolean("on"))
	a.False(f.MustBoolean("off"))
}

func TestBooleanField_MustBoolean_withInvalidBoolean(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultBooleanField("test")
	a.NoError(err)
	a.PanicsWithValue("MustBoolean called on test field with an invalid boolean value: invalid", func() {
		f.MustBoolean("invalid")
	})
}
