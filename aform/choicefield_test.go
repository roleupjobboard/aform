package aform_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChoiceField_Clean_withNoOptionAndNoValue(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultChoiceField("test")
	a.NoError(err)
	actual, errs := f.Clean("")
	a.Len(errs, 1)
	a.Equal("This field is required", errs[0].Error())
	a.Equal("", actual)
}

func TestChoiceField_Clean_withNoOption(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultChoiceField("test")
	a.NoError(err)
	actual, errs := f.Clean("option_first")
	a.Len(errs, 1)
	a.Equal("Invalid choice", errs[0].Error())
	a.Equal("option_first", actual)
}

func TestChoiceField_Clean_withValidValue(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultChoiceField("test", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{
		Value: "option_first",
		Label: "Option First",
	}}))
	a.NoError(err)
	actual, errs := f.Clean("option_first")
	a.Len(errs, 0)
	a.Equal("option_first", actual)
}

func TestChoiceField_Clean_withValidValueContainingSpace(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultChoiceField("test", aform.WithChoiceOptions([]aform.ChoiceFieldOption{
		{Value: "option first", Label: "Option First"},
		{Value: "option second", Label: "Option Second"},
		{Value: "option third", Label: "Option Third"},
	}))
	a.NoError(err)
	actual, errs := f.Clean("second")
	a.Len(errs, 1)
	a.Equal("Invalid choice", errs[0].Error())
	a.Equal("second", actual)
	actual, errs = f.Clean("option second")
	a.Len(errs, 0)
	a.Equal("option second", actual)
}

func TestChoiceField_Clean_withNoValueAndNotRequired(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultChoiceField("test", aform.IsNotRequired(), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{
		Value: "option_first",
		Label: "Option First",
	}}))
	a.NoError(err)
	actual, errs := f.Clean("")
	a.Len(errs, 0)
	a.Equal(f.EmptyValue(), actual, "it should return the EmptyValue()")
}

func TestChoiceField_Clean_withInvalidValue(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultChoiceField("test", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{
		Value: "option_first",
		Label: "Option First",
	}}))
	a.NoError(err)
	actual, errs := f.Clean("not_an_option")
	a.Len(errs, 1)
	a.Equal("Invalid choice", errs[0].Error())
	a.Equal("not_an_option", actual)
}

func TestChoiceField_Clean_withMultipleOptionsAndValidValue(t *testing.T) {
	a := assert.New(t)
	f, err := aform.DefaultChoiceField("test", aform.WithChoiceOptions([]aform.ChoiceFieldOption{
		{Value: "option_first", Label: "Option First"},
		{Value: "option_second", Label: "Option Second"},
		{Value: "option_third", Label: "Option Third"},
	}))
	a.NoError(err)
	actual, errs := f.Clean("option_first")
	a.Len(errs, 0)
	a.Equal("option_first", actual)
	actual, errs = f.Clean("not_option_second")
	a.Len(errs, 1)
	a.Equal("Invalid choice", errs[0].Error())
	a.Equal("not_option_second", actual)
	actual, errs = f.Clean("option_second")
	a.Len(errs, 0)
	a.Equal("option_second", actual)
	actual, errs = f.Clean("option_third")
	a.Len(errs, 0)
	a.Equal("option_third", actual)
}
