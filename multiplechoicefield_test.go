package aform_test

import (
	"github.com/roleupjobboard/aform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultipleChoiceField_Clean_withNoOptionAndNoValue(t *testing.T) {
	tests := []struct {
		name  string
		field *aform.MultipleChoiceField
	}{
		{
			name:  "default",
			field: aform.Must(aform.DefaultMultipleChoiceField("test")),
		},
		{
			name:  "widget CheckboxSelectMultiple",
			field: aform.Must(aform.DefaultMultipleChoiceField("test", aform.WithWidget(aform.CheckboxSelectMultiple))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			actual, errs := tt.field.Clean([]string{})
			a.Len(errs, 1)
			a.Equal("This field is required", errs[0].Error())
			a.Equal([]string{}, actual)
		})
	}
}

func TestMultipleChoiceField_Clean_withNoOptionAndOneValue(t *testing.T) {
	tests := []struct {
		name  string
		field *aform.MultipleChoiceField
	}{
		{
			name:  "default",
			field: aform.Must(aform.DefaultMultipleChoiceField("test")),
		},
		{
			name:  "widget CheckboxSelectMultiple",
			field: aform.Must(aform.DefaultMultipleChoiceField("test", aform.WithWidget(aform.CheckboxSelectMultiple))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			actual, errs := tt.field.Clean([]string{"option_1"})
			a.Len(errs, 1)
			a.Equal("Invalid choice", errs[0].Error())
			a.Equal([]string{"option_1"}, actual)
		})
	}
}

func Test_MultipleChoiceField_Clean_withNoOptionAndThreeValues(t *testing.T) {
	tests := []struct {
		name  string
		field *aform.MultipleChoiceField
	}{
		{
			name:  "default",
			field: aform.Must(aform.DefaultMultipleChoiceField("test")),
		},
		{
			name:  "widget CheckboxSelectMultiple",
			field: aform.Must(aform.DefaultMultipleChoiceField("test", aform.WithWidget(aform.CheckboxSelectMultiple))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			actual, errs := tt.field.Clean([]string{"not_option_1", "not_option_2", "not_option_3"})
			a.Len(errs, 3)
			a.Equal("Invalid choice", errs[0].Error())
			a.Equal("Invalid choice", errs[1].Error())
			a.Equal("Invalid choice", errs[2].Error())
			a.Equal([]string{"not_option_1", "not_option_2", "not_option_3"}, actual)
		})
	}
}

func Test_MultipleChoiceField_Clean_withOneValidValue(t *testing.T) {
	tests := []struct {
		name  string
		field *aform.MultipleChoiceField
	}{
		{
			name: "default",
			field: aform.Must(aform.DefaultMultipleChoiceField("test",
				aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "option_1", Label: "Option First"},
				}),
			)),
		},
		{
			name: "widget CheckboxSelectMultiple",
			field: aform.Must(aform.DefaultMultipleChoiceField("test",
				aform.WithWidget(aform.CheckboxSelectMultiple),
				aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "option_1", Label: "Option First"},
				}),
			)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			actual, errs := tt.field.Clean([]string{"option_1"})
			a.Len(errs, 0)
			a.Equal([]string{"option_1"}, actual)
		})
	}
}

func Test_MultipleChoiceField_Clean_withMultipleValidValue(t *testing.T) {
	tests := []struct {
		name  string
		field *aform.MultipleChoiceField
	}{
		{
			name: "default",
			field: aform.Must(aform.DefaultMultipleChoiceField("test",
				aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "option_1", Label: "Option First"},
					{Value: "option_2", Label: "Option Second"},
					{Value: "option_3", Label: "Option Third"},
					{Value: "option_4", Label: "Option Fourth"},
				}),
			)),
		},
		{
			name: "widget CheckboxSelectMultiple",
			field: aform.Must(aform.DefaultMultipleChoiceField("test",
				aform.WithWidget(aform.CheckboxSelectMultiple),
				aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "option_1", Label: "Option First"},
					{Value: "option_2", Label: "Option Second"},
					{Value: "option_3", Label: "Option Third"},
					{Value: "option_4", Label: "Option Fourth"},
				}),
			)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			actual, errs := tt.field.Clean([]string{"option_1", "option_3"})
			a.Len(errs, 0)
			a.Equal([]string{"option_1", "option_3"}, actual)
		})
	}
}

func Test_MultipleChoiceField_Clean_withMultipleValidValueAndOneInvalid(t *testing.T) {
	tests := []struct {
		name  string
		field *aform.MultipleChoiceField
	}{
		{
			name: "default",
			field: aform.Must(aform.DefaultMultipleChoiceField("test",
				aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "option_1", Label: "Option First"},
					{Value: "option_2", Label: "Option Second"},
					{Value: "option_3", Label: "Option Third"},
					{Value: "option_4", Label: "Option Fourth"},
				}),
			)),
		},
		{
			name: "widget CheckboxSelectMultiple",
			field: aform.Must(aform.DefaultMultipleChoiceField("test",
				aform.WithWidget(aform.CheckboxSelectMultiple),
				aform.WithChoiceOptions([]aform.ChoiceFieldOption{
					{Value: "option_1", Label: "Option First"},
					{Value: "option_2", Label: "Option Second"},
					{Value: "option_3", Label: "Option Third"},
					{Value: "option_4", Label: "Option Fourth"},
				}),
			)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			actual, errs := tt.field.Clean([]string{"option_1", "not_an_option", "option_3"})
			a.Len(errs, 1)
			a.Equal("Invalid choice", errs[0].Error())
			a.Equal([]string{"option_1", "not_an_option", "option_3"}, actual)
		})
	}
}
