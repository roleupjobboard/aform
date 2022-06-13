package aform

import (
	"golang.org/x/text/language"
	"html/template"
	"net/http"
)

// verify interface compliance
var _ formInterface = (*Form)(nil)

type formInterface interface {
	AsDiv() template.HTML
	BindRequest(req *http.Request)
	BindData(data map[string][]string, langs ...string)
	IsBound() bool
	Fields() []*Field
	FieldByName(field string) (*Field, error)
	IsValid() bool
	CleanedData() CleanedData
	Errors() FormErrors
	SetCleanFunc(clean func(*Form))
	AddError(field string, err error) error
	// CharField(name string) CharField
	// EmailField(name string) EmailField
}

type fieldInterface interface {
	Type() FieldType
	fieldInitializer
	fieldCustomizer
	fieldRenderer
	fieldReader
	field() *Field
}

type fieldInitializer interface {
	SetLabel(label string)
	MarkSafe()
	SetHelpText(help string)
	SetNotRequired()
	SetDisabled()
	AddChoiceOptions(label string, options []ChoiceFieldOption)
	addError(err Error)
}

type fieldCustomizer interface {
	SetLabelSuffix(labelSuffix string)
	CustomizeError(err ErrorCoderTranslator)
	SetAutoID(autoID string) error
	SetRequiredCSSClass(class string)
	SetErrorCSSClass(class string)
	SetAttributes(attrs []Attributable)
	SetWidget(widget Widget)
	SetLocale(locale language.Tag)
	SetSanitizeFunc(update func(current SanitizationFunc) (new SanitizationFunc))
	SetValidateFunc(update func(current ValidationFunc) (new ValidationFunc))
}

type fieldRenderer interface {
	AsDiv() template.HTML
	LabelTag() template.HTML
	LegendTag() template.HTML
	Widget() template.HTML
	Errors() template.HTML
	HelpText() template.HTML
}

type fieldReader interface {
	Name() string
	HTMLName() string
	AutoID() string
	LabelSuffix() string
	UseFieldset() bool
	CSSClasses() string
	Required() bool
	HasHelpText() bool
	HasErrors() bool
}
