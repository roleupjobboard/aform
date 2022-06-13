package aform_test

import (
	"fmt"
	"github.com/roleupjobboard/aform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestForm_IsValid_withoutBinding(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First")))))
	a.False(f.IsValid())
}

func TestForm_IsValid_withEmptyDataAndRequiredField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First")))))
	f.BindData(nil)
	a.False(f.IsValid())
}

func TestForm_IsValid_withValidDataAndRequiredField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First")))))
	f.BindData(map[string][]string{"first": {"good"}})
	a.True(f.IsValid())
}

func TestForm_IsValid_withEmptyDataAndNotRequiredField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First", aform.IsNotRequired())))))
	f.BindData(nil)
	a.True(f.IsValid())
}

func TestForm_IsValid_withValidDataAndNotRequiredField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("First", aform.IsNotRequired())))))
	f.BindData(map[string][]string{"first": {"good"}})
	a.True(f.IsValid())
}

func TestForm_IsValid_withEmptyDataAndMultipleRequiredField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithCharField(aform.Must(aform.DefaultCharField("First"))),
		aform.WithCharField(aform.Must(aform.DefaultCharField("Second"))),
		aform.WithEmailField(aform.Must(aform.DefaultEmailField("Third"))),
		aform.WithBooleanField(aform.Must(aform.DefaultBooleanField("Fourth"))),
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("Fifth", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "fifth_option_1", Label: "First Option"}})))),
	))
	f.BindData(nil)
	a.False(f.IsValid())
}

func TestForm_IsValid_withPartialDataAndMultipleRequiredField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithCharField(aform.Must(aform.DefaultCharField("First"))),
		aform.WithCharField(aform.Must(aform.DefaultCharField("Second"))),
		aform.WithEmailField(aform.Must(aform.DefaultEmailField("Third"))),
		aform.WithBooleanField(aform.Must(aform.DefaultBooleanField("Fourth"))),
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("Fifth", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "fifth_option_1", Label: "First Option"}})))),
	))
	f.BindData(map[string][]string{"second": {"good"}})
	a.False(f.IsValid())
}

func TestForm_IsValid_withAllDataAndMultipleRequiredField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithCharField(aform.Must(aform.DefaultCharField("First"))),
		aform.WithCharField(aform.Must(aform.DefaultCharField("Second"))),
		aform.WithEmailField(aform.Must(aform.DefaultEmailField("Third"))),
		aform.WithBooleanField(aform.Must(aform.DefaultBooleanField("Fourth"))),
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("Fifth", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "fifth_option_1", Label: "First Option"}})))),
	))
	f.BindData(map[string][]string{"first": {"good"}, "second": {"good"}, "third": {"good@email.com"}, "fourth": {"on"}, "fifth": {"fifth_option_1"}})
	a.True(f.IsValid())
}

func TestForm_IsValid_withInvalidEmailAndNotRequiredField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithEmailField(aform.Must(aform.DefaultEmailField("Email", aform.IsNotRequired()))),
	))
	f.BindData(map[string][]string{"email": {"invalid email"}})
	a.False(f.IsValid())
}

func TestForm_IsValid_withNoValueAndNotRequiredChoiceField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("Test", aform.IsNotRequired(), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "test_option_1", Label: "First Test Option"}})))),
	))
	f.BindData(nil)
	a.True(f.IsValid())
}

func TestForm_IsValid_withNoValueAndRequiredMultipleChoiceField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Test", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "test_option_1", Label: "First Test Option"}})))),
	))
	f.BindData(nil)
	a.False(f.IsValid())
}

func TestForm_IsValid_withNoValueAndNotRequiredMultipleChoiceField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Test", aform.IsNotRequired(), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "test_option_1", Label: "First Test Option"}})))),
	))
	f.BindData(nil)
	a.True(f.IsValid())
}

func TestForm_CleanedData_withNoBindingAndEmptyForm(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New())
	a.Len(f.CleanedData(), 0)
}

func TestForm_CleanedData_withNoBindingAndRequiredField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithCharField(aform.Must(aform.DefaultCharField("First"))),
	))
	a.Len(f.CleanedData(), 0)
}

func TestForm_CleanedData_withEmptyBindingAndRequiredField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithCharField(aform.Must(aform.DefaultCharField("First"))),
	))
	f.BindData(nil)
	a.Len(f.CleanedData(), 0)
}

func TestForm_CleanedData_withEmptyBindingAndNotRequiredField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithCharField(aform.Must(aform.DefaultCharField("First", aform.IsNotRequired()))),
	))
	f.BindData(nil)
	a.Equal(aform.CleanedData{"first": {""}}, f.CleanedData())
}

func TestForm_CleanedData_withValidValueBindingAndRequiredField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithCharField(aform.Must(aform.DefaultCharField("First"))),
	))
	f.BindData(map[string][]string{"first": {"hi!"}})
	a.Equal(aform.CleanedData{"first": {"hi!"}}, f.CleanedData())
}

func TestForm_CleanedData_withValidValueBindingAndNotRequiredField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithCharField(aform.Must(aform.DefaultCharField("First", aform.IsNotRequired()))),
	))
	f.BindData(map[string][]string{"first": {"hi!"}})
	a.Equal(aform.CleanedData{"first": {"hi!"}}, f.CleanedData())
}

func TestForm_CleanedData_withEmptyDataAndMultipleNotRequiredField(t *testing.T) {
	a := assert.New(t)
	first, _ := aform.NewCharField("First", "", "empty first", 0, 0, aform.IsNotRequired())
	second, _ := aform.NewCharField("Second", "", "empty second", 0, 0, aform.IsNotRequired())
	third, _ := aform.NewEmailField("Third", "", "empty third", 0, 0, aform.IsNotRequired())
	f := aform.Must(aform.New(
		aform.WithCharField(first),
		aform.WithCharField(second),
		aform.WithEmailField(third),
		aform.WithBooleanField(aform.Must(aform.DefaultBooleanField("Fourth", aform.IsNotRequired()))),
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("Fifth", aform.IsNotRequired(), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "fifth_option_1", Label: "First Option"}})))),
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Sixth", aform.IsNotRequired(), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "sixth_option_1", Label: "First Option"}, {Value: "sixth_option_2", Label: "Second Option"}})))),
	))
	f.BindData(nil)
	a.Len(f.Errors(), 0)
	a.Len(f.CleanedData(), 6)
	a.Equal([]string{"empty first"}, f.CleanedData()["first"])
	a.Equal([]string{"empty second"}, f.CleanedData()["second"])
	a.Equal([]string{"empty third"}, f.CleanedData()["third"])
	a.Equal([]string{"off"}, f.CleanedData()["fourth"])
	a.Equal([]string{""}, f.CleanedData()["fifth"])
	a.Equal([]string{}, f.CleanedData()["sixth"])
}

func TestForm_CleanedData_withEmptyDataAndMultipleRequiredField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithCharField(aform.Must(aform.DefaultCharField("First"))),
		aform.WithCharField(aform.Must(aform.DefaultCharField("Second"))),
		aform.WithEmailField(aform.Must(aform.DefaultEmailField("Third"))),
		aform.WithBooleanField(aform.Must(aform.DefaultBooleanField("Fourth"))),
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("Fifth", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "fifth_option_1", Label: "First Option"}})))),
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Sixth", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "sixth_option_1", Label: "First Option"}, {Value: "sixth_option_2", Label: "Second Option"}})))),
	))
	f.BindData(nil)
	a.Len(f.Errors(), 6)
	a.Len(f.Errors()["first"], 1)
	a.Len(f.Errors()["second"], 1)
	a.Len(f.Errors()["third"], 1)
	a.Len(f.Errors()["fourth"], 1)
	a.Len(f.Errors()["fifth"], 1)
	a.Len(f.Errors()["sixth"], 1)
	a.Equal("This field is required", f.Errors().Get("first").Error())
	a.Equal("This field is required", f.Errors().Get("second").Error())
	a.Equal("This field is required", f.Errors().Get("third").Error())
	a.Equal("This field is required", f.Errors().Get("fourth").Error())
	a.Equal("This field is required", f.Errors().Get("fifth").Error())
	a.Equal("This field is required", f.Errors().Get("sixth").Error())
	a.Len(f.CleanedData(), 0)
}

func TestForm_CleanedData_withEmptyDataAndMultipleRequiredChoiceField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("First", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "first_option_1", Label: "First Option"}, {Value: "first_option_2", Label: "Second Option"}})))),
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("Second", aform.WithWidget(aform.RadioSelect), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "second_option_1", Label: "First Option"}, {Value: "second_option_2", Label: "Second Option"}})))),
	))
	f.BindData(nil)
	a.Len(f.Errors(), 2)
	a.Len(f.Errors()["first"], 1)
	a.Len(f.Errors()["second"], 1)
	a.Equal("This field is required", f.Errors().Get("first").Error())
	a.Equal("This field is required", f.Errors().Get("second").Error())
}

func TestForm_CleanedData_withMultipleRequiredChoiceFieldAndTooManyData(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("First", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "first_option_1", Label: "First Option"}, {Value: "first_option_2", Label: "Second Option"}})))),
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("Second", aform.WithWidget(aform.RadioSelect), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "second_option_1", Label: "First Option"}, {Value: "second_option_2", Label: "Second Option"}})))),
		aform.WithChoiceField(aform.Must(aform.DefaultChoiceField("Third", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "third_option_1", Label: "First Option"}, {Value: "third_option_2", Label: "Second Option"}})))),
	))
	f.BindData(map[string][]string{"first": {"first_option_2", "first_option_1"}, "second": {"second_option_1", "second_option_2"}, "third": {"third_option_2", "ignored_invalid_value"}})
	a.Len(f.Errors(), 0)
	a.Len(f.CleanedData()["first"], 1, "choice field ignore bound data after the first one")
	a.Len(f.CleanedData()["second"], 1, "choice field ignore bound data after the first one")
	a.Len(f.CleanedData()["third"], 1, "choice field ignore bound data after the first one")
	a.Equal("first_option_2", f.CleanedData().Get("first"))
	a.Equal("second_option_1", f.CleanedData().Get("second"))
	a.Equal("third_option_2", f.CleanedData().Get("third"))
}

func TestForm_CleanedData_withEmptyDataAndMultipleRequiredMultipleChoiceField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("First", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "first_option_1", Label: "First Option"}, {Value: "first_option_2", Label: "Second Option"}})))),
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Second", aform.WithWidget(aform.CheckboxSelectMultiple), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "second_option_1", Label: "First Option"}, {Value: "second_option_2", Label: "Second Option"}})))),
	))
	f.BindData(nil)
	a.Len(f.Errors(), 2)
	a.Len(f.Errors()["first"], 1)
	a.Len(f.Errors()["second"], 1)
	a.Equal("This field is required", f.Errors().Get("first").Error())
	a.Equal("This field is required", f.Errors().Get("second").Error())
}

func TestForm_CleanedData_withMultipleRequiredMultipleChoiceFieldAndMultipleData(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("First", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "first_option_1", Label: "First Option"}, {Value: "first_option_2", Label: "Second Option"}})))),
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Second", aform.WithWidget(aform.CheckboxSelectMultiple), aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "second_option_1", Label: "First Option"}, {Value: "second_option_2", Label: "Second Option"}})))),
		aform.WithMultipleChoiceField(aform.Must(aform.DefaultMultipleChoiceField("Third", aform.WithChoiceOptions([]aform.ChoiceFieldOption{{Value: "third_option_1", Label: "First Option"}, {Value: "third_option_2", Label: "Second Option"}})))),
	))
	f.BindData(map[string][]string{"first": {"first_option_2", "first_option_1"}, "second": {"second_option_1", "second_option_2"}, "third": {"third_option_2", "invalid_value"}})
	a.Len(f.Errors(), 1)
	a.Len(f.CleanedData()["first"], 2)
	a.Len(f.CleanedData()["second"], 2)
	a.Len(f.CleanedData()["third"], 0)
	a.Equal([]string{"first_option_2", "first_option_1"}, f.CleanedData()["first"])
	a.Equal([]string{"second_option_1", "second_option_2"}, f.CleanedData()["second"])
}

func TestForm_AddError_toValidForm(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("Name")))))
	f.BindData(map[string][]string{"name": {"aby"}})
	a.True(f.IsValid())
	a.Len(f.Errors()["name"], 0)
	a.Len(f.CleanedData()["name"], 1)
	a.Equal("aby", f.CleanedData().Get("name"))
	anError := fmt.Errorf("any error")
	a.NoError(f.AddError("name", anError))
	a.Len(f.Errors()["name"], 1)
	a.ErrorIs(f.Errors()["name"][0], anError)
	a.Len(f.CleanedData()["name"], 0)
	a.False(f.IsValid())
}

func TestForm_AddError_toNotValidForm(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("Name")))))
	f.BindData(nil)
	a.False(f.IsValid())
	a.Len(f.Errors()["name"], 1)
	a.Len(f.CleanedData()["name"], 0)
	anError := fmt.Errorf("any error")
	a.NoError(f.AddError("name", anError))
	a.Len(f.Errors()["name"], 2)
	a.Equal(f.Errors()["name"][0].Error(), "This field is required")
	a.ErrorIs(f.Errors()["name"][1], anError)
	a.Len(f.CleanedData()["name"], 0)
	a.False(f.IsValid())
}

func TestForm_AddError_toNotValidatedForm(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("Name")))))
	a.EqualError(f.AddError("name", fmt.Errorf("any error")), "you can't add an error to a form not already validated. A form is validated when one of the following method is called: CleanedData(), IsValid() or Errors()")
}

func TestForm_AddError_withInvalidField(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("Name")))))
	f.BindData(map[string][]string{"name": {"aby"}})
	a.True(f.IsValid())
	a.EqualError(f.AddError("unknownfield", fmt.Errorf("any error")), "no field with this name unknownfield")
	a.True(f.IsValid())
}

func TestForm_SetCleanFunc_withValidForm(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("Title"))), aform.WithCharField(aform.Must(aform.DefaultCharField("Name")))))
	calls := 0
	cleanFunc := func(f *aform.Form) {
		calls++
	}
	f.SetCleanFunc(cleanFunc)
	f.BindData(map[string][]string{"title": {"miss"}, "name": {"aby"}})
	a.True(f.IsValid())
	a.Equal(1, calls)
	a.True(f.IsValid())
	a.Equal(1, calls)
	a.True(f.IsValid())
	a.Len(f.CleanedData(), 2)
	a.Len(f.Errors(), 0)
	a.Equal(1, calls)
}

func TestForm_SetCleanFunc_withNotValidForm(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("Title"))), aform.WithCharField(aform.Must(aform.DefaultCharField("Name")))))
	calls := 0
	cleanFunc := func(f *aform.Form) {
		calls++
	}
	f.SetCleanFunc(cleanFunc)
	f.BindData(map[string][]string{"name": {"aby"}})
	a.False(f.IsValid())
	a.Equal(1, calls)
	a.False(f.IsValid())
	a.Equal(1, calls)
	a.False(f.IsValid())
	a.Len(f.CleanedData(), 1)
	a.Len(f.Errors(), 1)
	a.Equal(1, calls)
}

func TestForm_SetCleanFunc_validFormCallValidationFuncsFromCleanFunc(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("Title"))), aform.WithCharField(aform.Must(aform.DefaultCharField("Name")))))
	calls := 0
	cleanFunc := func(f *aform.Form) {
		a.True(f.IsBound())
		a.True(f.IsValid())
		a.Len(f.CleanedData(), 2)
		a.Len(f.Errors(), 0)
		calls++
	}
	f.SetCleanFunc(cleanFunc)
	f.BindData(map[string][]string{"title": {"miss"}, "name": {"aby"}})
	a.True(f.IsValid())
	a.Equal(1, calls)
	a.True(f.IsValid())
	a.Equal(1, calls)
	a.True(f.IsValid())
	a.Len(f.CleanedData(), 2)
	a.Len(f.Errors(), 0)
	a.Equal(1, calls)
}

func TestForm_SetCleanFunc_notValidFormCallValidationFuncsFromCleanFunc(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("Title"))), aform.WithCharField(aform.Must(aform.DefaultCharField("Name")))))
	calls := 0
	cleanFunc := func(f *aform.Form) {
		a.True(f.IsBound())
		a.False(f.IsValid())
		a.Len(f.CleanedData(), 1)
		a.Len(f.Errors(), 1)
		calls++
	}
	f.SetCleanFunc(cleanFunc)
	f.BindData(map[string][]string{"name": {"aby"}})
	a.False(f.IsValid())
	a.Equal(1, calls)
	a.False(f.IsValid())
	a.Equal(1, calls)
	a.False(f.IsValid())
	a.Len(f.CleanedData(), 1)
	a.Len(f.Errors(), 1)
	a.Equal(1, calls)
}

func TestForm_SetCleanFunc_addCleanAfterValidation(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("Title"))), aform.WithCharField(aform.Must(aform.DefaultCharField("Name")))))
	calls := 0
	cleanFunc := func(f *aform.Form) {
		calls++
	}
	f.BindData(map[string][]string{"title": {"miss"}, "name": {"aby"}})
	a.True(f.IsValid())
	a.Equal(0, calls)
	f.SetCleanFunc(cleanFunc)
	a.True(f.IsValid())
	a.Equal(0, calls)
	a.True(f.IsValid())
	a.Len(f.CleanedData(), 2)
	a.Len(f.Errors(), 0)
	a.Equal(0, calls)
}

func TestForm_SetCleanFunc_withAddError(t *testing.T) {
	a := assert.New(t)
	f := aform.Must(aform.New(aform.WithCharField(aform.Must(aform.DefaultCharField("Title"))), aform.WithCharField(aform.Must(aform.DefaultCharField("Name")))))
	cleanFunc := func(f *aform.Form) {
		cleanedData := f.CleanedData()
		if cleanedData.Get("title") == "mr" {
			a.NoError(f.AddError("title", fmt.Errorf("she's a girl")))
		}
	}
	f.SetCleanFunc(cleanFunc)
	f.BindData(map[string][]string{"title": {"mr"}, "name": {"aby"}})
	a.False(f.IsValid())
	a.Len(f.CleanedData(), 1)
	a.Equal("aby", f.CleanedData().Get("name"))
	a.Len(f.Errors(), 1)
	a.Equal("she's a girl", f.Errors().Get("title").Error())
}
