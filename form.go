package aform

import (
	"fmt"
	"golang.org/x/text/language"
	"html/template"
	"net/http"
)

// Form represents a form.
type Form struct {
	fields           []fieldInterface
	autoID           string
	requiredCSSClass string
	errorCSSClass    string
	labelSuffix      string
	bound            bool
	validated        bool
	fieldNames       []string
	boundData        map[string][]string
	cleanedData      map[string][]string
	errors           map[string][]Error
	cleanFunc        func(*Form)
	locales          []language.Tag
}

// FormOption describes a functional option for configuring a Form.
type FormOption func(*Form) error

// FormPointerOrFieldPointer defines a union type to allow the usage of the helper
// function Must with forms and all fields types.
type FormPointerOrFieldPointer interface {
	*Form | *BooleanField | *EmailField | *CharField | *ChoiceField | *MultipleChoiceField
}

// Must is a helper that wraps a call to a function returning (*Form, error)
// or (*Field, error) and panics if the error is non-nil. It is intended
// for use in form creations. e.g.
//	var fld = Must(DefaultCharField("Name"))
//	var f = Must(New(WithCharField(fld)))
func Must[FormPtrFieldPtr FormPointerOrFieldPointer](f FormPtrFieldPtr, err error) FormPtrFieldPtr {
	if err != nil {
		panic(err)
	}
	return f
}

// New returns a Form.
func New(opts ...FormOption) (*Form, error) {
	f := &Form{
		autoID:    defaultAutoID,
		cleanFunc: func(f *Form) {},
	}
	for _, opt := range opts {
		if err := opt(f); err != nil {
			return nil, err
		}
	}
	return f, nil
}

// AsDiv renders the form as a list of <div> tags, with each <div> containing
// one field.
func (f *Form) AsDiv() template.HTML {
	return mustFormAsDivTemplate(f)
}

// BindRequest binds req form data to the Form. After a first binding, following
// bindings are ignored. If you want to bind new data, you should create another
// identical Form to do it. Data is bound but not validated.
// Validation is done when IsValid, CleanedData or Errors are called.
// Error messages are localized according to Accept-Language header. To modify
// this behavior use directly BindData.
func (f *Form) BindRequest(req *http.Request) {
	acceptLanguage := req.Header.Get("Accept-Language")
	f.BindData(req.Form, acceptLanguage)
	return
}

// BindData binds data to the Form. After a first binding, following bindings
// are ignored. If you want to bind new data, you should create another
// identical Form to do it. Data is bound but not validated.
// Validation is done when IsValid, CleanedData or Errors are called.
func (f *Form) BindData(data map[string][]string, langs ...string) {
	if f.bound {
		return
	}
	f.bound = true
	filteredData := map[string][]string{}
	for _, name := range f.fieldNames {
		values, ok := data[name]
		if ok {
			filteredData[name] = values
		}
	}
	f.boundData = filteredData
	propagateLocalesIfNotEmpty(f.fields, []language.Tag{selectLanguage(f.locales, langs...)})
	return
}

// IsBound returns true if the form is already bound to data
// with BindRequest or BindData.
func (f *Form) IsBound() bool {
	return f.bound
}

// WithBooleanField returns a FormOption that adds the BooleanField fld
// to the list of fields.
func WithBooleanField(fld *BooleanField) FormOption {
	return func(f *Form) error {
		return f.addField(fld)
	}
}

// WithCharField returns a FormOption that adds the CharField fld
// to the list of fields.
func WithCharField(fld *CharField) FormOption {
	return func(f *Form) error {
		return f.addField(fld)
	}
}

// WithEmailField returns a FormOption that adds the EmailField fld
// to the list of fields.
func WithEmailField(fld *EmailField) FormOption {
	return func(f *Form) error {
		return f.addField(fld)
	}
}

// WithChoiceField returns a FormOption that adds the ChoiceField fld
// to the list of fields.
func WithChoiceField(fld *ChoiceField) FormOption {
	return func(f *Form) error {
		return f.addField(fld)
	}
}

// WithMultipleChoiceField returns a FormOption that adds the
// MultipleChoiceField fld to the list of fields.
func WithMultipleChoiceField(fld *MultipleChoiceField) FormOption {
	return func(f *Form) error {
		return f.addField(fld)
	}
}

func (f *Form) addField(fld fieldInterface) error {
	f.fields = append(f.fields, fld)
	f.fieldNames = append(f.fieldNames, normalizedNameForField(fld))
	propagateLabelSuffix([]fieldInterface{fld}, f.labelSuffix)
	propagateRequiredCSSClassIfNotEmpty([]fieldInterface{fld}, f.requiredCSSClass)
	propagateErrorCSSClassIfNotEmpty([]fieldInterface{fld}, f.errorCSSClass)
	propagateLocalesIfNotEmpty([]fieldInterface{fld}, f.locales)
	return propagateAutoIDIfNotDefault([]fieldInterface{fld}, f.autoID)
}

// Fields returns the list of fields added to the form. First added comes
// first.
func (f *Form) Fields() []*Field {
	fields := make([]*Field, len(f.fields))
	for i, fld := range f.fields {
		fields[i] = fld.field()
	}
	return fields
}

// FieldByName returns the field with the normalized name field. If there is
// no Field matching, an error is returned.
func (f *Form) FieldByName(field string) (*Field, error) {
	fld, err := f.internalFieldByName(field)
	if err != nil {
		return nil, err
	}
	return fld.field(), nil
}

func (f *Form) internalFieldByName(field string) (fieldInterface, error) {
	nName := normalizedName(field)
	for _, fld := range f.fields {
		if fld.HTMLName() == field || fld.HTMLName() == nName {
			return fld, nil
		}
	}
	return nil, fmt.Errorf("no field with this name %s", field)
}

// WithAutoID returns a FormOption that changes how HTML IDs are automatically
// built for all the fields.
// Correct values are either an empty string or a string containing one verb
// '%s'. An empty string deactivates auto ID. It's equivalent as using
// DisableAutoID. With a string containing one verb '%s' it's possible to
// customize auto ID. e.g. "id_%s"
//
// A non-empty string without a verb '%s' is invalid.
func WithAutoID(autoID string) FormOption {
	return func(f *Form) error {
		if err := validAutoID(autoID); err != nil {
			return err
		}
		return setAndPropagateAutoID(f, autoID)
	}
}

// DisableAutoID returns a FormOption that deactivates the automatic ID
// creation for all the fields. When auto ID is disabled fields don't have
// a label tag unless a field override the auto ID behaviour with
// Field.SetAutoID.
func DisableAutoID() FormOption {
	return func(f *Form) error {
		return setAndPropagateAutoID(f, "")
	}
}

func setAndPropagateAutoID(f *Form, autoID string) error {
	f.autoID = autoID
	return propagateAutoIDIfNotDefault(f.fields, autoID)
}

func propagateAutoIDIfNotDefault(fields []fieldInterface, autoID string) error {
	if autoID == defaultAutoID {
		return nil
	}
	for _, fld := range fields {
		if fld.AutoID() != defaultAutoID {
			continue
		}
		if err := fld.SetAutoID(autoID); err != nil {
			return err
		}
	}
	return nil
}

// WithLabelSuffix returns a FormOption that set a suffix to labels. Label
// suffix is added to all fields.
func WithLabelSuffix(labelSuffix string) FormOption {
	return func(f *Form) error {
		f.labelSuffix = labelSuffix
		propagateLabelSuffix(f.fields, labelSuffix)
		return nil
	}
}

func propagateLabelSuffix(fields []fieldInterface, labelSuffix string) {
	for _, fld := range fields {
		if fld.LabelSuffix() != defaultLabelSuffix {
			continue
		}
		fld.SetLabelSuffix(labelSuffix)
	}
}

// WithRequiredCSSClass returns a FormOption that adds class to CSS classes
// when a field is required. CSS class is added to all required fields.
func WithRequiredCSSClass(class string) FormOption {
	return func(f *Form) error {
		f.requiredCSSClass = class
		propagateRequiredCSSClass(f.fields, class)
		return nil
	}
}

func propagateRequiredCSSClassIfNotEmpty(fields []fieldInterface, class string) {
	if len(class) == 0 {
		return
	}
	propagateRequiredCSSClass(fields, class)
}

func propagateRequiredCSSClass(fields []fieldInterface, class string) {
	for _, fld := range fields {
		fld.SetRequiredCSSClass(class)
	}
}

// WithErrorCSSClass returns a FormOption that adds class to CSS classes
// when a field validation returns an error. CSS class is added to all
// fields that have an error after validation.
func WithErrorCSSClass(class string) FormOption {
	return func(f *Form) error {
		f.errorCSSClass = class
		propagateErrorCSSClass(f.fields, class)
		return nil
	}
}

func propagateErrorCSSClassIfNotEmpty(fields []fieldInterface, class string) {
	if len(class) == 0 {
		return
	}
	propagateErrorCSSClass(fields, class)
}

func propagateErrorCSSClass(fields []fieldInterface, class string) {
	for _, fld := range fields {
		fld.SetErrorCSSClass(class)
	}
}

// WithLocales returns a FormOption that sets the list of locales used to
// translate error messages. Default locale is "en".
func WithLocales(locales []language.Tag) FormOption {
	return func(f *Form) error {
		f.locales = locales
		propagateLocalesIfNotEmpty(f.fields, locales)
		return nil
	}
}

func propagateLocalesIfNotEmpty(fields []fieldInterface, locales []language.Tag) {
	firstLocale := defaultLanguage
	if len(locales) > 0 {
		firstLocale = locales[0]
	}
	for _, fld := range fields {
		fld.SetLocale(firstLocale)
	}
}

// CleanedData maps a field normalized name to the list of bound data after validation
type CleanedData map[string][]string

// Get gets the first cleaned data associated with the given field.
// If there are no cleaned data associated with the field, Get returns
// the empty string. If the field accepts multiple data, use the map
// directly to access all of them.
func (d CleanedData) Get(field string) string {
	if d == nil {
		return ""
	}
	data, ok := d[field]
	if !ok {
		return ""
	}
	if len(data) == 0 {
		return ""
	}
	return data[0]
}

// Has checks whether a given field has at least one cleaned data.
func (d CleanedData) Has(field string) bool {
	_, ok := d[field]
	return ok
}

// FormErrors maps a field normalized name to the list of errors after validation
type FormErrors map[string][]Error

// Get gets the first error associated with the given field.
// If there are no errors associated with the field, Get returns
// an empty Error with its code and message equal empty string.
// To access multiple errors, use the map directly.
func (e FormErrors) Get(field string) Error {
	if e == nil {
		return Error{}
	}
	errs, ok := e[field]
	if !ok {
		return Error{}
	}
	if len(errs) == 0 {
		return Error{}
	}
	return errs[0]
}

// Has checks whether a given field has an error.
func (e FormErrors) Has(field string) bool {
	_, ok := e[field]
	return ok
}
